package zk

import (
	"errors"
	"fmt"
	"time"

	"github.com/qxnw/lib4go/encoding"
)

type getValueType struct {
	data []byte
	err  error
}

//GetValue 获取节点的值
func (client *ZookeeperClient) GetValue(path string) (value string, err error) {
	if !client.isConnect || client.conn == nil {
		err = errors.New("未连接到zk服务器")
		return
	}

	// 起一个协程，获取节点的值
	ch := make(chan interface{})
	go func(ch chan interface{}) {
		data, _, err := client.conn.Get(path)
		ch <- getValueType{data: data, err: err}
	}(ch)

	// 使用定时器判断获取节点的值是否超时
	tk := time.NewTicker(TIMEOUT)
	select {
	case _, ok := <-tk.C:
		if ok {
			tk.Stop()
			err = fmt.Errorf("get node : %s value timeout", path)
			return
		}
	case data := <-ch:
		tk.Stop()
		err = data.(getValueType).err
		if err != nil {
			return
		}
		value, err = encoding.Convert(data.(getValueType).data, "gbk")
		return
	}
	return
}

type getChildrenType struct {
	data []string
	err  error
}

//GetChildren 获取节点下的子节点
func (client *ZookeeperClient) GetChildren(path string) (paths []string, err error) {
	if !client.isConnect || client.conn == nil {
		err = errors.New("未连接到zk服务器")
		return
	}
	if b, err := client.Exists(path); !b || err != nil {
		return nil, err
	}

	// 起一个协程，获取子节点
	ch := make(chan interface{})
	go func(ch chan interface{}) {
		data, _, err := client.conn.Children(path)
		ch <- getChildrenType{data: data, err: err}
	}(ch)

	// 使用定时器判断获取子节点是否超时
	tk := time.NewTicker(TIMEOUT)
	select {
	case _, ok := <-tk.C:
		if ok {
			tk.Stop()
			err = fmt.Errorf("get node : %s children timeout", path)
			return
		}
	case data := <-ch:
		tk.Stop()
		paths = data.(getChildrenType).data
		err = data.(getChildrenType).err
		return
	}

	return
}
