package zk

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

//getPaths 获取当前路径的所有子路径
func (client *ZookeeperClient) getPaths(path string) []string {
	nodes := strings.Split(path, "/")
	len := len(nodes)
	paths := make([]string, 0, len-1)
	for i := 1; i < len; i++ {
		npath := "/" + strings.Join(nodes[1:i+1], "/")
		paths = append(paths, npath)
	}
	return paths
}

//ExistsAny 是否有一个路径已经存在
func (client *ZookeeperClient) ExistsAny(paths ...string) (b bool, path string, err error) {
	for _, path = range paths {
		if b, err = client.Exists(path); err != nil || b {
			return
		}
	}
	return
}

type existsType struct {
	b   bool
	err error
}

//Exists 检查路径是否存在
func (client *ZookeeperClient) Exists(path string) (b bool, err error) {
	if !client.isConnect {
		err = errors.New("未连接到zk服务器")
		return
	}

	// 启动一个协程，判断节点是否存在
	ch := make(chan interface{})
	go func(ch chan interface{}) {
		b, _, err = client.conn.Exists(path)
		ch <- existsType{b: b, err: err}
	}(ch)

	// 启动一个计时器，判断是否超时
	tk := time.NewTicker(TIMEOUT)
	select {
	case _, ok := <-tk.C:
		if ok {
			tk.Stop()
			err = fmt.Errorf("judgment node : %s exists timeout", path)
			return
		}
	case data := <-ch:
		tk.Stop()
		err = data.(existsType).err
		if err != nil {
			return
		}
		b = data.(existsType).b
		return
	}

	return
}

//GetDir 获取当前路径的目录
func (client *ZookeeperClient) GetDir(path string) string {
	paths := client.getPaths(path)
	if len(paths) > 2 {
		return paths[len(paths)-2]
	}
	return "/"
}

//CreatePersistentNode 创建持久化的节点
func (client *ZookeeperClient) CreatePersistentNode(path string, data string) (err error) {
	if !client.isConnect {
		err = errors.New("未连接到zk服务器")
		return
	}
	//检查目录是否存在
	if b, err := client.Exists(path); b || err != nil {
		err = fmt.Errorf("create node %s fail(node is exists : %t, err : %v)", path, b, err)
		return err
	}
	//获取每级目录并检查是否存在，不存在则创建
	paths := client.getPaths(path)
	for i := 0; i < len(paths)-1; i++ {
		b, err := client.Exists(paths[i])
		if err != nil {
			return err
		}
		if b {
			continue
		}
		_, err = client.create(paths[i], []byte(""), int32(0), zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	//创建最后一级目录
	_, err = client.create(path, []byte(data), int32(0), zk.WorldACL(zk.PermAll))
	if err != nil {
		return
	}
	return nil
}

//CreateTempNode 创建临时节点
func (client *ZookeeperClient) CreateTempNode(path string, data string) (rpath string, err error) {
	err = client.CreatePersistentNode(client.GetDir(path), "")
	if err != nil {
		return
	}
	rpath, err = client.create(path, []byte(data), int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll))
	return
}

//CreateSeqNode 创建临时节点
func (client *ZookeeperClient) CreateSeqNode(path string, data string) (rpath string, err error) {
	err = client.CreatePersistentNode(client.GetDir(path), "")
	if err != nil {
		return
	}
	rpath, err = client.create(path, []byte(data), int32(zk.FlagSequence)|int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll))
	return
}

type createType struct {
	rpath string
	err   error
}

func (client *ZookeeperClient) create(path string, data []byte, flags int32, acl []zk.ACL) (rpath string, err error) {
	if !client.isConnect || client.conn == nil {
		err = errors.New("未连接到zk服务器")
		return
	}

	// 开启一个协程，创建节点
	ch := make(chan interface{})
	go func(ch chan interface{}) {
		data, err := client.conn.Create(path, data, flags, acl)
		if err != nil {
			ch <- createType{err: err}
		} else {
			ch <- createType{rpath: data, err: err}
		}
	}(ch)

	// 使用计时器判断创建节点是否超时
	tk := time.NewTicker(TIMEOUT)
	select {
	case _, ok := <-tk.C:
		if ok {
			tk.Stop()
			err = fmt.Errorf("create node : %s timeout", path)
			return
		}
	case data := <-ch:
		tk.Stop()
		err = data.(createType).err
		if err != nil {
			return
		}
		rpath = data.(createType).rpath
		return
	}
	return
}
