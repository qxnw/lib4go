package zk

import "github.com/samuel/go-zookeeper/zk"

//BindWatchValue 监控指定节点的值是否发生变化，变化时返回变化后的值
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
		if !ok {
			return nil
		}
		switch e.Type {
		case zk.EventNodeCreated:
		case zk.EventNodeDataChanged:
			v, _ := client.GetValue(path)
			data <- v

		/*add by champly 2016年12月5日17:15:35*/
		// 关闭连接的时候报错处理
		case zk.EventNotWatching:
			return nil
			/*end*/
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
		if !ok {
			return nil
		}
		switch e.Type {
		case zk.EventNodeChildrenChanged:
			data <- []string{e.Type.String()}
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
