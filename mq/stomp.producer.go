package mq

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gmallard/stompngo"
	"github.com/qxnw/lib4go/concurrent/cmap"
)

//ProducerConfig 配置信息
type ProducerConfig struct {
	Address    string `json:"address"`
	Version    string `json:"version"`
	Persistent string `json:"persistent"`
}
type procuderMessage struct {
	headers []string
	queue   string
	data    string
	timeout time.Duration
}

//StompProducer Producer
type StompProducer struct {
	address    string
	conn       *stompngo.Connection
	messages   chan *procuderMessage
	queues     cmap.ConcurrentMap
	connecting bool
	closeCh    chan struct{}
	done       bool
	lk         sync.Mutex
	header     []string
	*option
}

//NewStompProducer 创建新的producer
func NewStompProducer(address string, opts ...Option) (producer *StompProducer, err error) {
	producer = &StompProducer{}
	producer.messages = make(chan *procuderMessage, 1000)
	for _, opt := range opts {
		opt(producer.option)
	}
	if strings.EqualFold(producer.option.version, "") {
		producer.option.version = "1.1"
	}
	if strings.EqualFold(producer.option.persistent, "") {
		producer.option.persistent = "true"
	}
	if strings.EqualFold(producer.option.ack, "") {
		producer.option.ack = "client-individual"
	}
	producer.header = stompngo.Headers{"accept-version", producer.option.version}
	return
}

//Connect  循环连接服务器
func (producer *StompProducer) Connect() error {
	err := producer.ConnectOnce()
	if err == nil {
		return nil
	}
	go func() {
	START:
		for {
			select {
			case <-producer.closeCh:
				break START
			case <-time.After(time.Second * 3):
				err = producer.ConnectOnce()
				if err == nil {
					break START
				}
				producer.option.logger.Error(err)
			}
		}
	}()
	return nil
}
func (producer *StompProducer) sendLoop() {
Loop:
	for {
		select {
		case msg := <-producer.messages:
			err := producer.conn.Send(msg.headers, msg.data)
			if err != nil {
				break Loop
			}
		}
	}
}

//ConnectOnce 连接到服务器
func (producer *StompProducer) ConnectOnce() (err error) {
	if producer.connecting {
		return nil
	}
	producer.lk.Lock()
	defer producer.lk.Unlock()
	if producer.connecting {
		return nil
	}
	producer.connecting = true
	defer func() {
		producer.connecting = false
	}()
	con, err := net.Dial("tcp", producer.address)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	producer.conn, err = stompngo.Connect(con, producer.header)
	if err != nil {
		return fmt.Errorf("mq 无法连接到MQ:%v", err)
	}
	return nil
}

//Send 发送消息
func (producer *StompProducer) Send(queue string, msg string, timeout time.Duration) (err error) {
	pm := &procuderMessage{queue: queue, data: msg, timeout: timeout}
	pm.headers = make([]string, 0, len(producer.header)+2)
	copy(pm.headers, producer.header)
	pm.headers = append(pm.headers, "destination", "/queue/"+queue)
	producer.messages <- pm
	return
}

//Close 关闭当前连接
func (producer *StompProducer) Close() {
	if producer.conn == nil || !producer.conn.Connected() {
		return
	}
	producer.conn.Disconnect(stompngo.Headers{})
}
