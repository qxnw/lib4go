package mq

/*
package mq

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

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
	config  ConsumerConfig
	conn    *stompngo.Connection
	queues  cmap.ConcurrentMap
	closeCh chan struct{}
	done    bool
	lk      sync.Mutex
	header  []string
}

//NewStompConsumerJSON 创建新的producer
func NewStompConsumerJSON(config string) (producer *StompConsumer, err error) {
	conf := ConsumerConfig{}
	err = json.Unmarshal([]byte(config), &conf)
	if err != nil {
		return nil, fmt.Errorf("mq 配置文件有误:%v", err)
	}
	return NewStompConsumer(conf)
}

//NewStompConsumer 创建新的Consumer
func NewStompConsumer(config ConsumerConfig) (consumer *StompConsumer, err error) {
	consumer = &StompConsumer{}
	consumer.closeCh = make(chan struct{})
	consumer.queues = cmap.New()
	if strings.EqualFold(config.Version, "") {
		config.Version = "1.1"
	}
	if strings.EqualFold(config.Persistent, "") {
		config.Persistent = "true"
	}
	if strings.EqualFold(config.Ack, "") {
		config.Ack = "client-individual"
	}
	consumer.config = config
	consumer.header = stompngo.Headers{"accept-version", "1.1"}
	return
}

//ConnectLoop  循环连接服务器
func (consumer *StompConsumer) ConnectLoop() error {
	err := consumer.Connect()
	if err == nil {
		return nil
	}
	for {
		select {
		case <-consumer.closeCh:
			return errors.New("mq consumer closed")
		case <-time.After(time.Second * 3):
			fmt.Println("reconnect.....")
			err = consumer.Connect()
			if err == nil {
				fmt.Println("connected...")
				return nil
			}
			fmt.Println(err)

		}
	}
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
		return fmt.Errorf("mq 无法连接到MQ:%v", err)
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
	if err != nil {
		return err
	}
	if !success {
		err = fmt.Errorf("重复订阅消息:%s", queue)
		return
	}
	msgChan := ch.(<-chan stompngo.MessageData)
START:
	for {
		select {
		case <-consumer.closeCh:
			break START
		case msg, ok := <-msgChan:
			if ok {
				message := NewStompMessage(consumer, &msg.Message)
				if message.Has() {
					callback(message)
				} else {
					consumer.ConnectLoop()
				}

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
	if consumer.conn == nil || !consumer.conn.Connected() {
		return
	}
	consumer.conn.Disconnect(stompngo.Headers{})
	close(consumer.closeCh)
}


*/
