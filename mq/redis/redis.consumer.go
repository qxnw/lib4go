package redis

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/mq"
	"github.com/zkfy/stompngo"
)

type consumerChan struct {
	msgChan     <-chan stompngo.MessageData
	unconsumeCh chan struct{}
}

//RedisConsumer Consumer
type RedisConsumer struct {
	address    string
	conn       *stompngo.Connection
	cache      cmap.ConcurrentMap
	queues     cmap.ConcurrentMap
	connecting bool
	closeCh    chan struct{}
	done       bool
	lk         sync.Mutex
	header     []string
	once       sync.Once
	*mq.OptionConf
}

//NewRedisConsumer 创建新的Consumer
func NewRedisConsumer(address string, opts ...mq.Option) (consumer *RedisConsumer, err error) {
	consumer = &RedisConsumer{address: address}
	consumer.OptionConf = &mq.OptionConf{Logger: logger.GetSession("mq.redis", logger.CreateSession())}
	consumer.closeCh = make(chan struct{})
	consumer.queues = cmap.New(2)
	consumer.cache = cmap.New(2)
	for _, opt := range opts {
		opt(consumer.OptionConf)
	}
	return
}

//Connect  循环连接服务器
func (consumer *RedisConsumer) Connect() error {
	return nil
}

//Consume 订阅消息
func (consumer *RedisConsumer) Consume(queue string, callback func(mq.IMessage)) (err error) {
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
func (consumer *RedisConsumer) consume(queue string, callback func(mq.IMessage)) (err error) {
	success, ch, err := consumer.queues.SetIfAbsentCb(queue, func(input ...interface{}) (c interface{}, err error) {
		queue := input[0].(string)
		header := stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", "queue", queue), "ack", consumer.Ack}
		consumer.conn.SetSubChanCap(10)
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
			message := mq.mes(consumer, &msg.Message)
			if message.Has() {
				go callback(message)
			} else {
				consumer.reconnect(queue)
				break START
			}
		}
	}
	return
}
func (consumer *RedisConsumer) reconnect(queue string) {
	if v, b := consumer.queues.Get(queue); b {
		ch := v.(*consumerChan)
		close(ch.unconsumeCh)
	}
	consumer.queues.Remove(queue)
	consumer.conn.Disconnect(stompngo.Headers{})
	consumer.Connect()
}

//UnConsume 取消注册消费
func (consumer *RedisConsumer) UnConsume(queue string) {
	if consumer.conn == nil {
		return
	}
	header := stompngo.Headers{"destination",
		fmt.Sprintf("/%s/%s", "queue", queue), "ack", consumer.Ack}
	consumer.conn.Unsubscribe(header)
	if v, b := consumer.queues.Get(queue); b {
		ch := v.(*consumerChan)
		close(ch.unconsumeCh)
	}
	consumer.queues.Remove(queue)
	consumer.cache.Remove(queue)
}

//Close 关闭当前连接
func (consumer *RedisConsumer) Close() {

	if consumer.conn == nil {
		return
	}
	consumer.once.Do(func() {
		close(consumer.closeCh)
	})

	consumer.queues.RemoveIterCb(func(key string, value interface{}) bool {
		ch := value.(*consumerChan)
		close(ch.unconsumeCh)
		return true
	})
	consumer.cache.Clear()
	go func() {
		defer recover()
		time.Sleep(time.Millisecond * 100)
		consumer.conn.Disconnect(stompngo.Headers{})
	}()

}

type redisConsumerResolver struct {
}

func (s *redisConsumerResolver) Resolve(address string, opts ...mq.Option) (mq.MQConsumer, error) {
	return NewRedisConsumer(address, opts...)
}
func init() {
	mq.RegisterCosnumer("redis", &redisConsumerResolver{})
}
