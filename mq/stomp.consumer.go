package mq

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"

	"errors"

	"github.com/gmallard/stompngo"
	"github.com/qxnw/lib4go/concurrent/cmap"
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
	lk     sync.Mutex
	header []string
}

//NewStompConsumerJSON 创建新的producer
func NewStompConsumerJSON(config string) (producer *StompConsumer, err error) {
	conf := ConsumerConfig{}
	// err = json.Unmarshal([]byte(config), &i.config)
	err = json.Unmarshal([]byte(config), &conf)
	if err != nil {
		return nil, fmt.Errorf("mq 配置文件有误:%v", err)
	}
	return NewStompConsumer(conf)
}

//NewStompConsumer 创建新的Consumer
func NewStompConsumer(config ConsumerConfig) (consumer *StompConsumer, err error) {
	consumer = &StompConsumer{}
	consumer.queues = cmap.New()
	consumer.config = config
	consumer.header = stompngo.Headers{"accept-version", config.Version}
	return
}

//Connect 连接到服务器
func (consumer *StompConsumer) Connect() error {
	if consumer.conn != nil && consumer.conn.Connected() {
		return nil
	}
	consumer.lk.Lock()
	defer consumer.lk.Unlock()
	if consumer.conn != nil && consumer.conn.Connected() {
		return nil
	}
	con, err := net.Dial("tcp", consumer.config.Address)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	consumer.conn, err = stompngo.Connect(con, consumer.header)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}

	return nil
}

//Consume 注册消费信息
func (consumer *StompConsumer) Consume(queue string, callback func(IMessage)) (err error) {
	if strings.EqualFold(queue, "") {
		return errors.New("队列名字不能为空")
	}
	if callback == nil {
		return errors.New("回调函数不能为nil")
	}

	if err = consumer.Connect(); err != nil {
		return
	}
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
	if err := consumer.Connect(); err != nil {
		return
	}
	header := stompngo.Headers{"destination",
		fmt.Sprintf("/%s/%s", "queue", queue), "ack", consumer.config.Ack}
	consumer.conn.Unsubscribe(header)
	if ch, b := consumer.queues.Get(queue); b {
		msgChan := ch.(chan stompngo.MessageData)
		close(msgChan)
	}
	consumer.queues.Remove(queue)
}

//Close 关闭当前连接
func (consumer *StompConsumer) Close() {
	// if err := consumer.Connect(); err != nil {
	// 	return
	// }

	/*change by champly 2016年11月30日15:21:45*/
	if consumer.conn == nil || !consumer.conn.Connected() {
		return
	}
	/*end*/

	consumer.conn.Disconnect(stompngo.Headers{})
}
