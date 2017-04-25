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
	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/utility"
)

type consumerChan struct {
	msgChan     <-chan stompngo.MessageData
	unconsumeCh chan struct{}
}
type option struct {
	logger *logger.Logger
}

//Option 配置选项
type Option func(*option)

//WithLogger 设置日志记录组件
func WithLogger(logger *logger.Logger) Option {
	return func(o *option) {
		o.logger = logger
	}
}

//ConsumerConfig 配置信息
type ConsumerConfig struct {
	Address    string `json:"address"`
	Version    string `json:"version"`
	Persistent string `json:"persistent"`
	Ack        string `json:"ack"`
}

//StompConsumer Consumer
type StompConsumer struct {
	config     ConsumerConfig
	conn       *stompngo.Connection
	cache      cmap.ConcurrentMap
	queues     cmap.ConcurrentMap
	connecting bool
	closeCh    chan struct{}
	done       bool
	lk         sync.Mutex
	header     []string
	*option
}

//NewStompConsumerJSON 创建新的producer
func NewStompConsumerJSON(config string, opts ...Option) (producer *StompConsumer, err error) {
	conf := ConsumerConfig{}
	err = json.Unmarshal([]byte(config), &conf)
	if err != nil {
		return nil, fmt.Errorf("mq 配置文件有误:%v", err)
	}
	return NewStompConsumer(conf, opts...)
}

//NewStompConsumer 创建新的Consumer
func NewStompConsumer(config ConsumerConfig, opts ...Option) (consumer *StompConsumer, err error) {
	consumer = &StompConsumer{}
	consumer.option = &option{logger: logger.GetSession("mq.consumer", utility.GetGUID())}
	consumer.closeCh = make(chan struct{})
	consumer.queues = cmap.New()
	consumer.cache = cmap.New()
	for _, opt := range opts {
		opt(consumer.option)
	}
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
	err := consumer.ConnectOnce()
	if err == nil {
		return nil
	}
	go func() {
	START:
		for {
			select {
			case <-consumer.closeCh:
				break START
			case <-time.After(time.Second * 3):
				err = consumer.ConnectOnce()
				if err == nil {
					break START
				}
				consumer.option.logger.Error(err)
			}
		}
	}()
	return nil
}

//ConnectOnce 连接到服务器
func (consumer *StompConsumer) ConnectOnce() (err error) {
	if consumer.connecting {
		return nil
	}
	consumer.lk.Lock()
	defer consumer.lk.Unlock()
	if consumer.connecting {
		return nil
	}
	consumer.connecting = true
	defer func() {
		consumer.connecting = false
	}()
	con, err := net.Dial("tcp", consumer.config.Address)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	consumer.conn, err = stompngo.Connect(con, consumer.header)
	if err != nil {
		return fmt.Errorf("mq 无法连接到MQ:%v", err)
	}

	//连接成功后开始订阅消息
	consumer.cache.IterCb(func(key string, value interface{}) bool {
		//consumer.option.logger.Info("consume:", key)
		go func() {
			err = consumer.consume(key, value.(func(IMessage)))
			if err != nil {
				fmt.Println(err)
			}
		}()
		return true
	})

	return nil
}

//Consume 订阅消息
func (consumer *StompConsumer) Consume(queue string, callback func(IMessage)) (err error) {
	if strings.EqualFold(queue, "") {
		return errors.New("队列名字不能为空")
	}
	if callback == nil {
		return errors.New("回调函数不能为nil")
	}
	b, _ := consumer.cache.SetIfAbsent(queue, callback)
	if !b {
		err = fmt.Errorf("重复订阅消息:%s", queue)
		return
	}
	return nil
}

//Consume 注册消费信息
func (consumer *StompConsumer) consume(queue string, callback func(IMessage)) (err error) {
	success, ch, err := consumer.queues.SetIfAbsentCb(queue, func(input ...interface{}) (c interface{}, err error) {
		queue := input[0].(string)
		header := stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", "queue", queue), "ack", consumer.config.Ack}
		msgChan, err := consumer.conn.Subscribe(header)
		if err != nil {
			return
		}
		chans := &consumerChan{}
		chans.msgChan = msgChan
		chans.unconsumeCh = make(chan struct{})
		return chans, nil
	}, queue)
	if err != nil {
		return err
	}
	if !success {
		err = fmt.Errorf("重复订阅消息:%s", queue)
		return
	}
	msgChan := ch.(*consumerChan)
START:
	for {
		select {
		case <-consumer.closeCh:
			break START
		case <-msgChan.unconsumeCh:
			break START
		case msg, ok := <-msgChan.msgChan:
			if !ok {
				break START
			}
			message := NewStompMessage(consumer, &msg.Message)
			if message.Has() {
				callback(message)
			} else {
				consumer.reconnect(queue)
				break START
			}
		}
	}
	return
}
func (consumer *StompConsumer) reconnect(queue string) {
	if v, b := consumer.queues.Get(queue); b {
		ch := v.(*consumerChan)
		close(ch.unconsumeCh)
	}
	consumer.queues.Remove(queue)
	consumer.conn.Disconnect(stompngo.Headers{})
	consumer.ConnectLoop()
}

//UnConsume 取消注册消费
func (consumer *StompConsumer) UnConsume(queue string) {
	if consumer.conn == nil {
		return
	}
	header := stompngo.Headers{"destination",
		fmt.Sprintf("/%s/%s", "queue", queue), "ack", consumer.config.Ack}
	consumer.conn.Unsubscribe(header)
	if v, b := consumer.queues.Get(queue); b {
		ch := v.(*consumerChan)
		close(ch.unconsumeCh)
	}
	consumer.queues.Remove(queue)
	consumer.cache.Remove(queue)
}

//Close 关闭当前连接
func (consumer *StompConsumer) Close() {
	if consumer.conn == nil {
		return
	}
	close(consumer.closeCh)
	consumer.queues.RemoveIterCb(func(key string, value interface{}) bool {
		ch := value.(*consumerChan)
		close(ch.unconsumeCh)
		return true
	})
	consumer.cache.Clear()
	consumer.conn.Disconnect(stompngo.Headers{})

}
