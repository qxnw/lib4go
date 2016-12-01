package zk

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/encoding"
	"github.com/qxnw/lib4go/logger"
	"github.com/samuel/go-zookeeper/zk"
)

//ZookeeperClient zookeeper客户端
type ZookeeperClient struct {
	servers            []string
	timeout            time.Duration
	conn               *zk.Conn
	eventChan          <-chan zk.Event
	watchValueEvents   cmap.ConcurrentMap
	watchChilrenEvents cmap.ConcurrentMap
	Log                logger.ILogger
	useCount           int32
	isConnect          bool
}

//New 连接到Zookeeper服务器
func New(servers []string, timeout time.Duration, loggerName string) (*ZookeeperClient, error) {
	client := &ZookeeperClient{servers: servers, timeout: timeout, useCount: 0}
	client.watchValueEvents = cmap.New()
	client.watchChilrenEvents = cmap.New()

	client.Log = logger.New(loggerName)
	client.conn.SetLogger(client.Log)
	return client, nil
}

//Connect 连接到远程zookeeper服务器
func (client *ZookeeperClient) Connect() (err error) {
	if client.conn == nil {
		conn, eventChan, err := zk.Connect(client.servers, client.timeout)
		if err != nil {
			return err
		}
		client.conn = conn
		client.eventChan = eventChan
	}
	atomic.AddInt32(&client.useCount, 1)
	return
}

//Reconnect 重新连接服务器
func (client *ZookeeperClient) Reconnect() (err error) {
	if client.conn != nil {
		client.conn.Close()
	}
	return client.Connect()
}

//Disconnect 断开服务器连接
func (client *ZookeeperClient) Disconnect() {
	atomic.AddInt32(&client.useCount, -1)
	if client.useCount > 0 {
		return
	}
	client.isConnect = false
	client.conn.Close()
}

//Exists 检查路径是否存在
func (client *ZookeeperClient) Exists(path string) bool {
	exists, _, _ := client.conn.Exists(path)
	return exists
}

//Delete 修改指定节点的值
func (client *ZookeeperClient) Delete(path string) error {
	return client.conn.Delete(path, -1)
}

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
	//检查目录是否存在
	if client.Exists(path) {
		return
	}
	//获取每级目录并检查是否存在，不存在则创建
	paths := client.getPaths(path)
	for i := 0; i < len(paths)-1; i++ {
		if client.Exists(paths[i]) {
			continue
		}
		_, err = client.conn.Create(paths[i], []byte(""), int32(0), zk.WorldACL(zk.PermAll))
		if err != nil {
			return
		}
	}
	//创建最后一级目录
	_, err = client.conn.Create(path, []byte(data), int32(0), zk.WorldACL(zk.PermAll))
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
	rpath, err = client.conn.Create(path, []byte(data), int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll))
	return
}

//CreateSeqNode 创建临时节点
func (client *ZookeeperClient) CreateSeqNode(path string, data string) (rpath string, err error) {
	err = client.CreatePersistentNode(client.GetDir(path), "")
	if err != nil {
		return
	}
	rpath, err = client.conn.Create(path, []byte(data), int32(zk.FlagSequence)|int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll))
	return
}

//GetValue 获取节点的值
func (client *ZookeeperClient) GetValue(path string) (string, error) {
	data, _, err := client.conn.Get(path)
	if err != nil {
		return "", err
	}
	return encoding.Convert(data, "gbk")
}

//GetChildren 获取子节点路径
func (client *ZookeeperClient) GetChildren(path string) (paths []string, err error) {
	if !client.Exists(path) {
		return
	}
	paths, _, err = client.conn.Children(path)
	return
}

//EventNotify 等待服务器失去连接
func (client *ZookeeperClient) EventNotify() {
START:
	for {
		select {
		case v := <-client.eventChan:
			switch v.State {
			case zk.StateAuthFailed:
				break
			case zk.StateConnected:
				break START
			case zk.StateExpired:
				break START
			case zk.StateDisconnected:
				break START
			}
		}
	}
}
