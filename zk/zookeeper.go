package zk

import (
	"sync/atomic"
	"time"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/logger"
	"github.com/samuel/go-zookeeper/zk"
)

// TIMEOUT 连接zk服务器操作的超时时间
var TIMEOUT = time.Second * 2

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

	if client.conn != nil {
		client.conn.Close()
	}

	client.isConnect = false
	client.isCloseManually = true
	client.conn = nil

}
