package zk

import (
	"fmt"
	"time"
)

type getValueType struct {
	data []byte
	err  error
}

//GetValue 获取节点的值
func (client *ZookeeperClient) GetValue(path string) (value []byte, err error) {
	if !client.isConnect {
		err = ErrColientCouldNotConnect
		return
	}
	if client.done {
		err = ErrClientConnClosing
		return
	}
	// 起一个协程，获取节点的值
	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		data, _, err := client.conn.Get(path)
		ch <- getValueType{data: data, err: err}
	}(ch)

	select {
	case <-time.After(TIMEOUT):
		err = fmt.Errorf("get node : %s value timeout", path)
		return
	case data := <-ch:
		if client.done {
			err = ErrClientConnClosing
			return
		}
		err = data.(getValueType).err
		if err != nil {
			return
		}
		value = data.(getValueType).data
		return
	}
}

type getChildrenType struct {
	data []string
	err  error
}

//GetChildren 获取节点下的子节点
func (client *ZookeeperClient) GetChildren(path string) (paths []string, err error) {
	if !client.isConnect {
		err = ErrColientCouldNotConnect
		return
	}
	if client.done {
		err = ErrClientConnClosing
		return
	}
	if b, err := client.Exists(path); !b || err != nil {
		return nil, fmt.Errorf("node(%s) is not exist", path)
	}

	// 起一个协程，获取子节点
	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		data, _, err := client.conn.Children(path)
		ch <- getChildrenType{data: data, err: err}
	}(ch)

	// 使用定时器判断获取子节点是否超时
	select {
	case <-time.After(TIMEOUT):
		err = fmt.Errorf("get node(%s) children timeout", path)
		return
	case data := <-ch:
		if client.done {
			err = ErrClientConnClosing
			return
		}
		paths = data.(getChildrenType).data
		err = data.(getChildrenType).err
		return
	}
}
