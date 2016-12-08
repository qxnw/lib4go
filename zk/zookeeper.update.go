package zk

import (
	"errors"
	"fmt"
	"time"
)

// Update 更新一个节点的值，如果存在则更新，如果不存在则报错
func (client *ZookeeperClient) Update(path string, data string) (err error) {
	if !client.isConnect || client.conn == nil {
		err = errors.New("未连接到zk服务器")
		return
	}

	// 判断节点是否存在
	if b, err := client.Exists(path); !b || err != nil {
		err = fmt.Errorf("update node %s fail(node is exists : %t, err : %v)", path, b, err)
		return err
	}

	// 启动一个协程，更新节点
	ch := make(chan error)
	go func(ch chan error) {
		_, err = client.conn.Set(path, []byte(data), -1)
		ch <- err
	}(ch)

	// 启动一个计时器，判断更新节点是否超时
	tk := time.NewTicker(TIMEOUT)
	select {
	case _, ok := <-tk.C:
		if ok {
			tk.Stop()
			err = fmt.Errorf("update node %s timeout", path)
			return
		}
	case err = <-ch:
		tk.Stop()
		return
	}

	return
}
