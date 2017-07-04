package kafka

import (
	"fmt"
	"time"

	"github.com/jdamick/kafka"
	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/mq"
)

//KafkaConsumer kafka consumer
type KafkaConsumer struct {
	address   string
	consumers cmap.ConcurrentMap
	quitChan  chan struct{}
	*mq.OptionConf
}
type kafkaConsumer struct {
	consumer *kafka.BrokerConsumer
	msgQueue chan *kafka.Message
}

//NewKafkaConsumer 初始化kafka Consumer
func NewKafkaConsumer(address string, opts ...mq.Option) (kafka *KafkaConsumer, err error) {
	kafka = &KafkaConsumer{address: address, quitChan: make(chan struct{}, 0)}
	kafka.OptionConf = &mq.OptionConf{}
	kafka.consumers = cmap.New(2)
	for _, opt := range opts {
		opt(kafka.OptionConf)
	}
	return
}

//Connect 连接到服务器
func (k *KafkaConsumer) Connect() error {
	return nil
}

//Consume 订阅消息
func (k *KafkaConsumer) Consume(queue string, call func(mq.IMessage)) (err error) {
	_, cnsmr, _ := k.consumers.SetIfAbsentCb(queue, func(i ...interface{}) (interface{}, error) {
		c := &kafkaConsumer{}
		c.consumer = kafka.NewBrokerConsumer(k.address, queue, 0, 0, 1048576)
		c.msgQueue = make(chan *kafka.Message, 10000)
		return c, nil
	})
	consumer := cnsmr.(*kafkaConsumer)
	conChan := make(chan error, 1)
	go func() {
		_, err = consumer.consumer.ConsumeOnChannel(consumer.msgQueue, 10, k.quitChan)
		conChan <- err
	}()
	select {
	case <-time.After(time.Second):
	case err := <-conChan:
		return err
	}
	go func() {
	LOOP:
		for {
			select {
			case msg, ok := <-consumer.msgQueue:
				fmt.Println(msg, ok)
				if ok {
					call(NewKafkaMessage(msg))
				} else {
					break LOOP
				}
			}
		}
	}()
	return nil
}

//UnConsume 取消注册消费
func (k *KafkaConsumer) UnConsume(queue string) {
	if c, ok := k.consumers.Get(queue); ok {
		consumer := c.(*kafkaConsumer)
		close(consumer.msgQueue)
	}
}

//Close 关闭当前连接
func (k *KafkaConsumer) Close() {
	close(k.quitChan)
	k.consumers.IterCb(func(key string, value interface{}) bool {
		consumer := value.(*kafkaConsumer)
		close(consumer.msgQueue)
		return true
	})
}

type kafkaConsumerResolver struct {
}

func (s *kafkaConsumerResolver) Resolve(address string, opts ...mq.Option) (mq.MQConsumer, error) {
	return NewKafkaConsumer(address, opts...)
}
func init() {
	mq.RegisterCosnumer("kafka", &kafkaConsumerResolver{})
}
