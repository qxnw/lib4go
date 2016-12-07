package zk

import (
	"strings"
	"sync/atomic"
	"time"

	"errors"

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

	// 是否是手动关闭
	isCloseManually bool
}

//New 连接到Zookeeper服务器
func New(servers []string, timeout time.Duration, loggerName string) (*ZookeeperClient, error) {
	client := &ZookeeperClient{servers: servers, timeout: timeout, useCount: 0}
	client.watchValueEvents = cmap.New()
	client.watchChilrenEvents = cmap.New()

	client.Log = logger.New(loggerName)
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
		client.conn.SetLogger(client.Log)
		client.eventChan = eventChan
		go client.eventWatch()
	}
	atomic.AddInt32(&client.useCount, 1)
	return
}

//Reconnect 重新连接服务器
func (client *ZookeeperClient) Reconnect() (err error) {
	if client.conn != nil {
		client.conn.Close()
		client.conn = nil
	}
	return client.Connect()
}

//Disconnect 断开服务器连接
func (client *ZookeeperClient) Disconnect() {
	atomic.AddInt32(&client.useCount, -1)
	if client.useCount > 0 {
		return
	}
	// client.isConnect = false
	// client.conn.Close()

	/*add by champly 2016年12月02日16:21:16*/
	if client.conn != nil {
		client.conn.Close()
	}

	client.isConnect = false
	client.isCloseManually = true
	client.conn = nil
	/*end*/

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

//Exists 检查路径是否存在
func (client *ZookeeperClient) Exists(path string) (b bool, err error) {
	if !client.isConnect {
		err = errors.New("未连接到zk服务器")
		return
	}
	b, _, err = client.conn.Exists(path)
	return
}

//Delete 修改指定节点的值
func (client *ZookeeperClient) Delete(path string) error {
	if !client.isConnect {
		return errors.New("未连接到zk服务器")
	}
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
	if !client.isConnect {
		err = errors.New("未连接到zk服务器")
		return
	}
	//检查目录是否存在
	if b, err := client.Exists(path); b || err != nil {
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
		_, err = client.conn.Create(paths[i], []byte(""), int32(0), zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
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

	ch := make(chan interface{})
	go func(ch chan interface{}) {
		data, _, err := client.conn.Get(path)

		if err != nil {
			ch <- getValueType{data: []byte(""), err: err}
		}
		ch <- getValueType{data: data, err: err}
	}(ch)

	var data interface{}

	tk := time.NewTicker(time.Second * 2)
	select {
	case _, ok := <-tk.C:
		if ok {
			return "", errors.New("connect to zk timeout")
		}
	case data = <-ch:
		tk.Stop()
		if data.(getValueType).err != nil {
			return "", err
		}
		value, err = encoding.Convert(data.(getValueType).data, "gbk")
	}

	return
}

type getChildrenType struct {
	data []string
	err  error
}

//GetChildren 获取子节点路径
func (client *ZookeeperClient) GetChildren(path string) (paths []string, err error) {
	if !client.isConnect || client.conn == nil {
		err = errors.New("未连接到zk服务器")
		return
	}
	if b, err := client.Exists(path); !b || err != nil {
		return nil, err
	}

	ch := make(chan interface{})
	go func(ch chan interface{}) {
		data, _, err := client.conn.Children(path)

		if err != nil {
			ch <- getChildrenType{data: nil, err: err}
		}
		ch <- getChildrenType{data: data, err: err}
	}(ch)

	var data interface{}

	tk := time.NewTicker(time.Second * 2)
	select {
	case _, ok := <-tk.C:
		if ok {
			return []string{""}, errors.New("connect to zk timeout")
		}
	case data = <-ch:
		tk.Stop()
		paths = data.(getChildrenType).data
		err = data.(getChildrenType).err
	}

	return
}

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
				/*add by champly 2016年12月6日10:37:10*/
				// 网络断开，正在连接
				case zk.StateConnecting:
					client.isConnect = false
				case zk.StateHasSession:
					client.isConnect = true
					/*end*/
				}
			} else {
				client.isConnect = false
				break START
			}
		}
	}
}
