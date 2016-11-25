package mq

import (
	"fmt"
	"net"

	"github.com/gmallard/stompngo"
	"github.com/qxnw/lib4go/concurrent"
)

//ConsumerConfig 配置信息
type ConsumerConfig struct {
	Address    string `json:"address"`
	Version    string `json:"version"`
	Persistent string `json:"persistent"`
	Ack        string `json:"ack"`
}

//StompConsumer Consumer
type StompConsumer struct {
	config ConsumerConfig
	conn   *stompngo.Connection
	queues cmap.ConcurrentMap
	header []string
}

//NewStompConsumer 创建新的Consumer
func NewStompConsumer(config ConsumerConfig) (consumer *StompConsumer, err error) {
	consumer = &StompConsumer{}
	consumer.queues = cmap.New()
	consumer.config = config
	conn, err := net.Dial("tcp", config.Address)
	if err != nil {
		return
	}
	consumer.header = stompngo.Headers{"accept-version", config.Version}
	consumer.conn, err = stompngo.Connect(conn, consumer.header)
	return
}

//Consume 注册消费信息
func (consumer *StompConsumer) Consume(queue string, callback func(IMessage)) (err error) {
	success, ch, err := consumer.queues.SetIfAbsentCb(queue, func(input ...interface{}) (ch interface{}, err error) {
		queue := input[0].(string)
		header := stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", "queue", queue), "ack", consumer.config.Ack}
		msgChan, err := consumer.conn.Subscribe(header)
		if err != nil {
			return
		}
		ch = msgChan
		return
	}, queue)

	if !success {
		err = fmt.Errorf("重复订阅消息:%s", queue)
		return
	}
	msgChan := ch.(<-chan stompngo.MessageData)
START:
	for {
		select {
		case msg, ok := <-msgChan:
			if ok {
				callback(NewStompMessage(consumer, &msg.Message))
			} else {
				break START
			}
		}
	}
	return
}

//UnConsume 取消注册消费
func (consumer *StompConsumer) UnConsume(queue string) {
	header := stompngo.Headers{"destination",
		fmt.Sprintf("/%s/%s", "queue", queue), "ack", consumer.config.Ack}
	consumer.conn.Unsubscribe(header)
	consumer.queues.Remove(queue)
}

//Close 关闭当前连接
func (consumer *StompConsumer) Close() {
	if !consumer.conn.Connected() {
		return
	}
	consumer.conn.Disconnect(stompngo.Headers{})
}
