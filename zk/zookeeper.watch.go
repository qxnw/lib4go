package zk

import (
	"errors"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

//BindWatchValue 监控指定节点的值是否发生变化，变化时返回变化后的值
// 测试情况：
//		网络正常时修改节点的值：
//			EventNodeDataChanged : {Type:EventNodeDataChanged State:Unknown Path:/zk_test/123 Err:<nil> Server:}   true
// 		网络断开之后，节点值的修改不会触发，直到网络恢复正常：
//			EventNotWatching(断开时间过短不会出现) : {Type:EventNotWatching State:StateDisconnected Path:/zk_test/123 Err:zk: session has been expired by the server Server:} true
//		关闭连接:
//			EventNotWatching : {Type:EventNotWatching State:StateDisconnected Path:/zk_test/123 Err:zk: zookeeper is closing Server:}      true
func (client *ZookeeperClient) BindWatchValue(path string, data chan string) error {
	_, value := client.watchValueEvents.SetIfAbsent(path, 0) //添加/更新监控时间
	if value.(int) == -1 {
		client.watchValueEvents.Remove(path)
		return errors.New(path + " is UnbindWatchValue")
	}
	_, _, event, err := client.conn.GetW(path)
	if err != nil {
		return err
	}
	select {
	case e, ok := <-event:
		client.Log.Infof("watch:value %+v[%+v]%t", path, e, ok)
		if !ok {
			return e.Err
		}
		switch e.Type {
		case zk.EventNodeDataChanged:
			v, err := client.GetValue(path)
			if err != nil {
				client.Log.Error(err)
			} else {
				data <- v
			}
		case zk.EventNotWatching:
			err = client.checkConnectStatus(path, false)
			if err != nil {
				return err
			}
		}
	}

	//继续监控值变化
	return client.BindWatchValue(path, data)
}

//UnbindWatchValue 取消绑定
func (client *ZookeeperClient) UnbindWatchValue(path string) {
	if v, ok := client.watchValueEvents.Get(path); !ok || v.(int) == -1 {
		return
	}
	client.watchValueEvents.Set(path, -1)
}

//BindWatchChildren 监控子节点是否发生变化，变化时返回变化后的值
// 测试情况：
//		网络正常时修改节点的值：
//			EventNodeChildrenChanged : {Type:EventNodeChildrenChanged State:Unknown Path:/zk_test Err:<nil> Server:}   true
// 		网络断开之后，节点值的修改不会触发，直到网络恢复正常：
//			EventNotWatching(断开时间过短不会出现) : {Type:EventNotWatching State:StateDisconnected Path:/zk_test Err:zk: session has been expired by the server Server:} true
//		关闭连接
//			EventNotWatching : {Type:EventNotWatching State:StateDisconnected Path:/zk_test Err:zk: zookeeper is closing Server:}       true
func (client *ZookeeperClient) BindWatchChildren(path string, data chan []string) (err error) {
	_, value := client.watchChilrenEvents.SetIfAbsent(path, 0) //添加/更新监控时间
	if value.(int) == -1 {
		client.watchChilrenEvents.Remove(path)
		return errors.New(path + " is UnbindWatchChildren")
	}
	_, _, event, err := client.conn.ChildrenW(path)
	if err != nil {
		return err
	}
	select {
	case e, ok := <-event:
		client.Log.Infof("watch:children %s[%+v]%t", path, e, ok)
		if !ok {
			return e.Err
		}
		switch e.Type {
		case zk.EventNodeChildrenChanged:
			paths, err := client.GetChildren(path)
			if err != nil {
				client.Log.Error(err)
			} else {
				data <- paths
			}
		// 网络重新连接
		case zk.EventNotWatching:
			err = client.checkConnectStatus(path, true)
			if err != nil {
				return err
			}
		}
	}

	return client.BindWatchChildren(path, data)
}

//UnbindWatchChildren 取消绑定
func (client *ZookeeperClient) UnbindWatchChildren(path string) {
	if v, ok := client.watchChilrenEvents.Get(path); !ok || v.(int) == -1 {
		return
	}
	client.watchChilrenEvents.Set(path, -1)
}

// checkConnectStatus 检查当前的连接状态
func (client *ZookeeperClient) checkConnectStatus(path string, isWatchChildren bool) error {
	if client.isCloseManually {
		return zk.ErrClosing
	}
	ticker := time.NewTicker(time.Second)
START:
	for {
		select {
		case _, ok := <-ticker.C:
			if ok {
				// 检查是否手动关闭连接
				if client.isCloseManually {
					ticker.Stop()
					return zk.ErrClosing
				}

				if isWatchChildren {
					if v, ok := client.watchChilrenEvents.Get(path); !ok || v.(int) == -1 {
						ticker.Stop()
						return errors.New(path + " is UnbindWatchChildren")
					}
				} else {
					// 检查是否取消绑定
					if v, ok := client.watchValueEvents.Get(path); !ok || v.(int) == -1 {
						ticker.Stop()
						return errors.New(path + " is UnbindWatchValue")
					}
				}

				// 检查是否连接成功
				if client.isConnect {
					ticker.Stop()
					break START
				}
			}
		}
	}
	return nil
}
