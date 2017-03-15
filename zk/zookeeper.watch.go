/*

各种情况下遇到的触发状态：


网络状态									函数									触发状态
															{Type:EventSession State:StateConnecting Path: Err:<nil> Server:192.168.0.159:2181}    true,
开始连接								eventWatch			{Type:EventSession State:StateConnected Path: Err:<nil> Server:192.168.0.159:2181}     true,
															{Type:EventSession State:StateHasSession Path: Err:<nil> Server:192.168.0.159:2181}    true,

										eventWatch			{Type:EventSession State:StateDisconnected Path: Err:<nil> Server:192.168.0.159:2181}  true,
															{Type:EventSession State:StateConnecting Path: Err:<nil> Server:192.168.0.159:2181}    true,
网络断开								BindWatchValue		无
										BindWatchChildren	无

										eventWatch			{Type:EventSession State:StateConnecting Path: Err:<nil> Server:192.168.0.159:2181}    true,
网络重连								BindWatchValue		-
										BindWatchChildren	-

															{Type:EventSession State:StateExpired Path: Err:<nil> Server:192.168.0.159:2181}  true【网络断开时间过短不会出现】,
															{Type:EventSession State:StateDisconnected Path: Err:<nil> Server:192.168.0.159:2181}        true,
网络恢复之后							eventWatch			{Type:EventSession State:StateConnecting Path: Err:<nil> Server:192.168.0.159:2181}  true,
															{Type:EventSession State:StateConnected Path: Err:<nil> Server:192.168.0.159:2181}   true,
															{Type:EventSession State:StateHasSession Path: Err:<nil> Server:192.168.0.159:2181}    true,

										eventWatch			{Type:EventSession State:StateDisconnected Path: Err:<nil> Server:192.168.0.159:2181}  true,
连接断开													{Type:Unknown State:StateDisconnected Path: Err:<nil> Server:} false,
										BindWatchValue		{Type:EventNotWatching State:StateDisconnected Path:/zk_test/123 Err:zk: zookeeper is closing Server:}  true【如果当前线程不是马上关闭会触发】
										BindWatchChildren	{Type:EventNotWatching State:StateDisconnected Path:/zk_test Err:zk: zookeeper is closing Server:}       true【如果当前线程不是马上关闭会触发】

修改节点的值(网络正常)					eventWatch			{Type:EventNodeDataChanged State:Unknown Path:/zk_test/123 Err:<nil> Server:}  true
										BindWatchValue		{Type:EventNodeDataChanged State:Unknown Path:/zk_test/123 Err:<nil> Server:}  true

修改节点的值（网络断开）				eventWatch			同网络连接断开
										BindWatchValue		无

修改节点的值（网络恢复正常）			eventWatch			同网络恢复之后
										BindWatchValue		{Type:EventNotWatching State:StateDisconnected Path:/zk_test/123 Err:zk: session has been expired by the server Server:}    true【如果断开时间过短不会触发】

修改节点的值（网络恢复正常之后修改）	eventWatch			{Type:EventNodeDataChanged State:Unknown Path:/zk_test/123 Err:<nil> Server:}  true
										BindWatchValue		{Type:EventNodeDataChanged State:Unknown Path:/zk_test/123 Err:<nil> Server:}       true

修改子节点（网络正常）					eventWatch			{Type:EventNodeChildrenChanged State:Unknown Path:/zk_test Err:<nil> Server:}  true
										BindWatchChildren	{Type:EventNodeChildrenChanged State:Unknown Path:/zk_test Err:<nil> Server:}    true

修改子节点（网络断开）					eventWatch			同网络连接断开
										BindWatchChildren	无

修改子节点（网络恢复正常）				eventWatch			同网络恢复之后
										BindWatchChildren	{Type:EventNotWatching State:StateDisconnected Path:/zk_test Err:zk: session has been expired by the server Server:}  true【断开时间过短不会触发】

修改节点的值（网络恢复正常之后修改）	eventWatch	 		{Type:EventNodeChildrenChanged State:Unknown Path:/zk_test Err:<nil> Server:}  true
										BindWatchChildren	{Type:EventNodeChildrenChanged State:Unknown Path:/zk_test Err:<nil> Server:}    true
*/
package zk

import (
	"errors"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

//eventWatch 服务器事件监控[重点测试]
// StateAuthFailed: 未测试
// StateConnected: 连接到服务器成功；网络从异常中恢复之后会出现
// StateExpired: 连接成功之后网络出现异常，从异常中恢复之后首先会出现这个状态
// StateDisconnected: 网络连接断开
// StateConnecting: 网络连接断开，如果没有关闭链接（网络异常），会一直发送请求，直到网络成功连接
// StateHasSession: 连接成功，获取到服务器的Session
// 状态顺序描述：【linux系统：修改防火墙规则：iptables -A OUTPUT -p tcp --dport 2181 -j DROP && iptables -A OUTPUT -p tcp --sport 2181 -j DROP】
// 		开始连接：
//			StateConnecting :	{Type:EventSession State:StateConnecting Path: Err:<nil> Server:192.168.0.159:2181}	true
//			->StateConnected :	{Type:EventSession State:StateConnected Path: Err:<nil> Server:192.168.0.159:2181}	true
//			->StateHasSession : {Type:EventSession State:StateHasSession Path: Err:<nil> Server:192.168.0.159:2181}	true
//			(连接成功)
//		断开网络：
//			StateDisconnected :	{Type:EventSession State:StateDisconnected Path: Err:<nil> Server:192.168.0.159:2181}	true
//			->StateConnecting :	{Type:EventSession State:StateConnecting Path: Err:<nil> Server:192.168.0.159:2181}		true
//			(一直到网络恢复)
//		网络恢复：
//			StateExpired(网络异常时间过短不会出现) : {Type:EventSession State:StateExpired Path: Err:<nil> Server:192.168.0.159:2181}	true
//			->StateDisconnected : {Type:EventSession State:StateDisconnected Path: Err:<nil> Server:192.168.0.159:2181} true
//			->StateConnecting :   {Type:EventSession State:StateConnecting Path: Err:<nil> Server:192.168.0.159:2181}   true
//			->StateConnected :	  {Type:EventSession State:StateConnected Path: Err:<nil> Server:192.168.0.159:2181}    true
//			->StateHasSession :	  {Type:EventSession State:StateHasSession Path: Err:<nil> Server:192.168.0.159:2181}   true
//			(连接成功)
//		正常关闭连接:
//			StateDisconnected :   {Type:EventSession State:StateDisconnected Path: Err:<nil> Server:192.168.0.159:2181} true
//			->StateDisconnected : {Type:Unknown State:StateDisconnected Path: Err:<nil> Server:}						false
//			(连接关闭)
func (client *ZookeeperClient) eventWatch() {
START:
	for {
		select {
		case v, ok := <-client.eventChan:
			if ok {
				client.Log.Infof("event.watch:%+v", v)
				switch v.State {
				case zk.StateAuthFailed:
					client.isConnect = false
				// 已经连接成功
				case zk.StateConnected:
					client.isConnect = true
				// 连接Session失效
				case zk.StateExpired:
					client.isConnect = false
				// 网络连接不成功
				case zk.StateDisconnected:
					client.isConnect = false
				// 网络断开，正在连接
				case zk.StateConnecting:
					client.isConnect = false
				case zk.StateHasSession:
					client.isConnect = true
				}
			} else {
				client.isConnect = false
				break START
			}
		}
	}
}

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

func (client *ZookeeperClient) WatchChildren(path string) (ch []string, err error) {
	_, _, event, err := client.conn.ChildrenW(path)
	if err != nil {
		return nil, err
	}
	select {
	case e, ok := <-event:
		client.Log.Infof("watch:children %s %s[%+v]%t", e.Type.String(), path, e, ok)
		if !ok {
			return nil, e.Err
		}
		switch e.Type {
		case zk.EventNodeChildrenChanged:
			paths, err := client.GetChildren(path)
			if err != nil {
				client.Log.Error(err)
			} else {
				return paths, nil
			}
		// 网络重新连接
		case zk.EventNotWatching:
			err = client.checkConnectStatus(path, true)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, errors.New("not watch")
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
		client.Log.Infof("watch:children %s %s[%+v]%t", e.Type.String(), path, e, ok)
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
