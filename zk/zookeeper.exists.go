package zk

import (
	"errors"
	"fmt"
	"time"
)

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
