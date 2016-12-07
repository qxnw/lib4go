package zk

import (
	"fmt"

	"github.com/samuel/go-zookeeper/zk"
)

//BindWatchValue 监控指定节点的值是否发生变化，变化时返回变化后的值
// 测试情况：
//		网络正常时修改节点的值：
//			EventNodeDataChanged : {Type:EventNodeDataChanged State:Unknown Path:/zk_test/123 Err:<nil> Server:}   true
// 		网络断开之后，节点值的修改不会触发，直到网络恢复正常之后修改节点的值：
//			EventNodeDataChanged : {Type:EventNodeDataChanged State:Unknown Path:/zk_test/123 Err:<nil> Server:}   true
func (client *ZookeeperClient) BindWatchValue(path string, data chan string) error {
	_, value := client.watchValueEvents.SetIfAbsent(path, 0) //添加/更新监控时间
	if value.(int) == -1 {
		client.watchValueEvents.Remove(path)
		return nil
	}
	_, _, event, err := client.conn.GetW(path)
	if err != nil {
		return err
	}
	select {
	case e, ok := <-event:
		client.Log.Infof("watch:value %s[%+v]%t", path, e, ok)
		if !ok {
			return nil
		}
		switch e.Type {
		case zk.EventNodeCreated:
		case zk.EventNodeDataChanged:
			v, _ := client.GetValue(path)
			data <- v
		case zk.EventNotWatching:
			// 如果是手动关闭，则不继续监控
			if client.isCloseManually {
				return nil
			}
			fmt.Println("EventNotWatching")
		}
	}

	//继续监控值变化
	return client.BindWatchValue(path, data)
}

//UnbindWatchValue 取消绑定
// 测试情况：
//		网络正常时取消绑定不会触发，直到节点值变化：
//			EventNodeDataChanged ： {Type:EventNodeDataChanged State:Unknown Path:/zk_test/123 Err:<nil> Server:}   true
//			BindWatchValue返回err
// 		网络断开之后，取消绑定，不会触发，直到网络恢复正常之后：
//			EventNotWatching : {Type:EventNotWatching State:StateDisconnected Path:/zk_test/123 Err:zk: session has been expired by the server Server:}        true
//			BindWatchValue返回err
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
// 		网络断开之后，节点值的修改不会触发，直到网络恢复正常之后修改节点的值：
//			EventNodeChildrenChanged : {Type:EventNodeChildrenChanged State:Unknown Path:/zk_test Err:<nil> Server:}   true
func (client *ZookeeperClient) BindWatchChildren(path string, data chan []string) (err error) {
	_, value := client.watchChilrenEvents.SetIfAbsent(path, 0) //添加/更新监控时间
	if value.(int) == -1 {
		client.watchChilrenEvents.Remove(path)
		return nil
	}
	_, _, event, err := client.conn.ChildrenW(path)
	if err != nil {
		return
	}
	select {
	case e, ok := <-event:
		client.Log.Infof("watch:children %s[%+v]%t", path, e, ok)
		if !ok {
			return nil
		}
		switch e.Type {
		case zk.EventNodeChildrenChanged:
			data <- []string{e.Type.String()}
			// value, err := client.GetChildren(path)
			// if err != nil {
			// 	return err
			// }
			// data <- value

			// // 网络重新连接
			// case zk.EventNotWatching:
		}
	}

	/*add by champly 2016年12月6日16:08:32*/
	// 如果是手动关闭，则不继续监控
	if client.isCloseManually {
		return nil
	}
	/*end*/

	return client.BindWatchChildren(path, data)
}

//UnbindWatchChildren 取消绑定
// 测试情况：
//		网络正常时取消绑定不会触发，直到节点值变化：
//			EventNodeChildrenChanged : {Type:EventNodeChildrenChanged State:Unknown Path:/zk_test Err:<nil> Server:}   true
//			BindWatchChildren返回err
// 		网络断开之后，取消绑定，不会触发，直到网络恢复正常之后：
//			EventNotWatching : {Type:EventNotWatching State:StateDisconnected Path:/zk_test Err:zk: session has been expired by the server Server:}    true
//			BindWatchChildren返回err
func (client *ZookeeperClient) UnbindWatchChildren(path string) {
	if v, ok := client.watchChilrenEvents.Get(path); !ok || v.(int) == -1 {
		return
	}
	client.watchChilrenEvents.Set(path, -1)
}
