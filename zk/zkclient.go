package zk

import (
	"errors"
	"strings"
	"sync/atomic"
	"time"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/encoding"
	"github.com/qxnw/lib4go/logger"
	"github.com/samuel/go-zookeeper/zk"
)

//zkClient zookeeper客户端
type zkClient struct {
	servers        []string
	timeout        time.Duration
	conn           *zk.Conn
	eventChan      <-chan zk.Event
	rmWatchValue   cmap.ConcurrentMap
	rmWatchChilren cmap.ConcurrentMap
	Log            logger.ILogger
	close          bool
	useCount       int32
}

//NewZK 连接到Zookeeper服务器
func NewZK(servers []string, timeout time.Duration, loggerName string) (*zkClient, error) {
	client := &zkClient{servers: servers, timeout: timeout, close: false, useCount: 0}
	client.rmWatchValue = cmap.New()
	client.rmWatchChilren = cmap.New()
	conn, eventChan, err := zk.Connect(servers, timeout)
	if err != nil {
		return nil, err
	}
	client.conn = conn
	client.eventChan = eventChan
	client.Log = logger.New(loggerName)
	client.conn.SetLogger(client.Log)
	return client, nil
}

//Open 打开网络连接
func (client *zkClient) Open() {
	atomic.AddInt32(&client.useCount, 1)
}

//Close 关闭client及与服务器的连接
func (client *zkClient) Close() {
	defer client.recover()
	atomic.AddInt32(&client.useCount, -1)
	if client.useCount > 0 {
		return
	}
	client.close = true
	client.conn.Close()
}

//Reconnect 重新连接
func (client *zkClient) Reconnect() error {
	conn, eventChan, err := zk.Connect(client.servers, client.timeout)
	if err != nil {
		return err
	}
	client.close = false
	client.eventChan = eventChan
	client.conn = conn
	return nil
}

// Exists check whether the path exists
func (client *zkClient) Exists(path ...string) (string, bool) {
	for _, v := range path {
		exists, _, _ := client.conn.Exists(v)
		if exists {
			return v, true
		}
	}
	return "", false
}

//UpdateValue 修改指定节点的值
func (client *zkClient) UpdateValue(path string, value string) error {
	_, err := client.conn.Set(path, []byte(value), -1)
	return err
}

//Delete 修改指定节点的值
func (client *zkClient) Delete(path string) error {
	return client.conn.Delete(path, -1)
}

//WaitForConnected 等待服务器连接成功
func (client *zkClient) WaitForConnected() bool {
	connected := false
START:
	for {
		select {
		case v := <-client.eventChan:
			switch v.State {
			case zk.StateConnected:
				connected = true
				break START
			}
		}
	}
	return connected
}

//WaitForDisconnected 等待服务器失去连接
func (client *zkClient) WaitForDisconnected() bool {
	disconnected := false
	tk := time.NewTicker(time.Second * 35)
START:
	for {
		select {
		case <-tk.C:
			disconnected = true
			break START
		case v := <-client.eventChan:
			switch v.State {
			case zk.StateExpired:
				disconnected = true
				break START
			case zk.StateDisconnected:
				disconnected = true
				break START
			}
		}
	}
	return disconnected
}

func getPaths(path string) []string {
	nodes := strings.Split(path, "/")
	len := len(nodes)
	var nlist []string
	for i := 1; i < len; i++ {
		npath := "/" + strings.Join(nodes[1:i+1], "/")
		nlist = append(nlist, npath)
	}
	return nlist
}
func (client *zkClient) recover() {
	if r := recover(); r != nil {
		client.Log.Error("zk:执行异常,", r)
	}
}

//CreateNode 创建持久节点
func (client *zkClient) CreateNode(path string, data string) error {
	paths := getPaths(path)
	l := len(paths)
	for index, value := range paths {
		ndata := ""
		if index == l-1 {
			ndata = data
		}
		_, err := client.create(value, ndata, int32(0))
		if err != nil {
			return err
		}
	}
	return nil
}

//CreateSeqNode 创建有序节点
func (client *zkClient) CreateSeqNode(path string, data string) (string, error) {
	err := client.createNodeRoot(path)
	if err != nil {
		return "", err
	}
	return client.create(path, data, int32(zk.FlagSequence)|int32(zk.FlagEphemeral))
}

//CreateTmpNode 创建临时节点
func (client *zkClient) CreateTmpNode(path string, data string) (string, error) {
	err := client.createNodeRoot(path)
	if err != nil {
		return "", err
	}
	return client.create(path, data, int32(zk.FlagEphemeral))
}

//CreateNode 创建临时节点
func (client *zkClient) createNodeRoot(path string) error {
	paths := getPaths(path)
	if len(paths) > 1 {
		root := paths[len(paths)-2]
		err := client.CreateNode(root, "")
		if err != nil {
			return err
		}
	}
	return nil
}

//create 根据参数创建路径
func (client *zkClient) create(path string, data string, flags int32) (string, error) {
	exists, _, err := client.conn.Exists(path)
	if exists && err == nil {
		return path, nil
	}
	acl := zk.WorldACL(zk.PermAll)
	npath, err := client.conn.Create(path, []byte(data), flags, acl)
	return npath, err
}

//GetValue 获取指定节点的值
func (client *zkClient) GetValue(path string) (string, error) {
	data, _, err := client.conn.Get(path)
	if err != nil {
		return "", err
	}
	return encoding.Convert(data, "gbk")
}

//GetChildren 获取指定节点的值
func (client *zkClient) GetChildren(path string) ([]string, error) {
	if _, ok := client.Exists(path); !ok {
		return []string{}, nil
	}
	data, _, err := client.conn.Children(path)
	if err != nil {
		return []string{}, err
	}
	return data, nil
}

//WatchConnectionChange 监控指定节点的值是否发生变化，变化时返回变化后的值
func (client *zkClient) WatchConnectionChange(data chan string) error {
	for {
		select {
		case v := <-client.eventChan:
			switch v.State {
			case zk.StateConnected:
				select {
				case data <- "connected":
				default:
				}
			case zk.StateDisconnected:
				select {
				case data <- "disconnected":
				default:
				}
			case zk.StateExpired:
				select {
				case data <- "expired":
				default:
				}
			case zk.StateAuthFailed:
				select {
				case data <- "authfailed":
				default:
				}
			default:
			}
		}
	}
}

//WatchValue 监控指定节点的值是否发生变化，变化时返回变化后的值
func (client *zkClient) WatchValue(path string, data chan string) error {
	if client.close {
		return nil
	}
	_, _, event, err := client.conn.GetW(path)
	if err != nil {
		return err
	}
	e := <-event
	switch e.Type {
	case zk.EventNotWatching:
	case zk.EventNodeCreated:
	case zk.EventNodeDeleted:
	case zk.EventSession:
	case zk.EventNodeDataChanged:
		if !client.rmWatchValue.Has(path) {
			v, _ := client.GetValue(path)
			data <- v
		}
	}
	if client.rmWatchValue.Has(path) {
		client.rmWatchValue.Remove(path)
		return errors.New("已移除节点监控")
	}
	return client.WatchValue(path, data)
}

//RemoveWatchValue 移除值监控
func (client *zkClient) RemoveWatchValue(path string) {
	client.rmWatchValue.Set(path, path)
}

//RemoveWatchChildren 移除子节点监控
func (client *zkClient) RemoveWatchChildren(path string) {
	client.rmWatchValue.Set(path, path)
}

//WatchChildren 监控指定节点的值是否发生变化，变化时返回变化后的值
func (client *zkClient) WatchChildren(path string, data chan []string) (err error) {
	if client.close {
		return nil
	}
	if _, ok := client.Exists(path); !ok {
		return nil
	}
	_, _, event, err := client.conn.ChildrenW(path)
	if err != nil {
		return
	}
	select {
	case e := <-event:
		if !client.rmWatchChilren.Has(path) {
			data <- []string{e.Type.String()}
		}

	}
	if client.rmWatchChilren.Has(path) {
		client.rmWatchChilren.Remove(path)
		return errors.New("已移除节点监控")
	}
	return client.WatchChildren(path, data)
}
