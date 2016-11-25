package mq

import (
	"fmt"
	"net"
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
	header []string
}

//NewStompProducer 创建新的producer
func NewStompProducer(config ProducerConfig) (producer *StompProducer, err error) {
	producer = &StompProducer{}
	producer.config = config
	conn, err := net.Dial("tcp", config.Address)
	if err != nil {
		return
	}
	producer.header = stompngo.Headers{"accept-version", config.Version}
	producer.conn, err = stompngo.Connect(conn, producer.header)
	return
}

//Send 发送消息
func (producer *StompProducer) Send(queue string, msg string, timeout int) (err error) {
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
	if !producer.conn.Connected() {
		return
	}
	producer.conn.Disconnect(stompngo.Headers{})
}
