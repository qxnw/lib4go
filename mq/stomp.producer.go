package mq

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gmallard/stompngo"
)

//ProducerConfig 配置信息
type ProducerConfig struct {
	Address    string `json:"address"`
	Version    string `json:"version"`
	Persistent string `json:"persistent"`
}

//StompProducer Producer
type StompProducer struct {
	config ProducerConfig
	conn   *stompngo.Connection
	lk     sync.Mutex
	header []string
}

//NewStompProducer 创建新的producer
func NewStompProducer(config ProducerConfig) (producer *StompProducer, err error) {
	producer = &StompProducer{}
	producer.config = config
	producer.header = stompngo.Headers{"accept-version", config.Version}
	return
}

//Connect 连接到服务器
func (producer *StompProducer) Connect() error {
	if producer.conn != nil && producer.conn.Connected() {
		return nil
	}

	producer.lk.Lock()
	defer producer.lk.Unlock()
	if producer.conn != nil && producer.conn.Connected() {
		return nil
	}
	con, err := net.Dial("tcp", producer.config.Address)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	producer.conn, err = stompngo.Connect(con, producer.header)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}

	return nil
}

//Send 发送消息
func (producer *StompProducer) Send(queue string, msg string, timeout int) (err error) {
	if err = producer.Connect(); err != nil {
		return
	}
	header := stompngo.Headers{"destination", queue, "persistent", producer.config.Persistent}
	if timeout > 0 {
		header = stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", "queue", queue),
			"persistent", producer.config.Persistent, "expires",
			fmt.Sprintf("%d000", time.Now().Add(time.Second*time.Duration(timeout)).Unix())}
	}
	err = producer.conn.Send(header, msg)
	return
}

//Close 关闭当前连接
func (producer *StompProducer) Close() {
	if err := producer.Connect(); err != nil {
		return
	}
	producer.conn.Disconnect(stompngo.Headers{})
}
