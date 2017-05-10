package zk

import (
	"errors"
	"sync/atomic"
	"time"

	"sync"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/logger"
	"github.com/samuel/go-zookeeper/zk"
)

// TIMEOUT 连接zk服务器操作的超时时间
var TIMEOUT = time.Second

/*
type Logger interface {
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
	Printf(string, ...interface{})
}
*/
var (
	ErrColientCouldNotConnect = errors.New("zk: could not connect to the server")
	ErrClientConnClosing      = errors.New("zk: the client connection is closing")
)

//ZookeeperClient zookeeper客户端
type ZookeeperClient struct {
	servers            []string
	timeout            time.Duration
	conn               *zk.Conn
	eventChan          <-chan zk.Event
	watchValueEvents   cmap.ConcurrentMap
	watchChilrenEvents cmap.ConcurrentMap
	Log                *logger.Logger
	useCount           int32
	isConnect          bool
	once               sync.Once
	CloseCh            chan struct{}
	// 是否是手动关闭
	done bool
}

//New 连接到Zookeeper服务器
func New(servers []string, timeout time.Duration) (*ZookeeperClient, error) {
	client := &ZookeeperClient{servers: servers, timeout: timeout, useCount: 0}
	client.CloseCh = make(chan struct{})
	client.watchValueEvents = cmap.New()
	client.watchChilrenEvents = cmap.New()
	client.Log = logger.GetSession("zk", logger.CreateSession())
	return client, nil
}

//NewWithLogger 连接到Zookeeper服务器
func NewWithLogger(servers []string, timeout time.Duration, logger *logger.Logger) (*ZookeeperClient, error) {
	client := &ZookeeperClient{servers: servers, timeout: timeout, useCount: 0}
	client.CloseCh = make(chan struct{})
	client.watchValueEvents = cmap.New()
	client.watchChilrenEvents = cmap.New()
	client.Log = logger
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
	time.Sleep(client.timeout)
	client.isConnect = true
	return
}

//Reconnect 重新连接服务器
func (client *ZookeeperClient) Reconnect() (err error) {
	if client.conn != nil {
		client.conn.Close()
		client.conn = nil
	}
	client.done = false
	return client.Connect()
}

//Close 关闭服务器
func (client *ZookeeperClient) Close() {
	atomic.AddInt32(&client.useCount, -1)
	if client.useCount > 0 {
		return
	}

	if client.conn != nil {
		client.conn.Close()
	}

	client.isConnect = false
	client.done = true
	client.once.Do(func() {
		close(client.CloseCh)
	})
	client.conn = nil

}
