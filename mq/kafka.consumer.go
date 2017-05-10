package mq

import "github.com/jdamick/kafka"
import "github.com/qxnw/lib4go/concurrent/cmap"

//KafkaConsumer kafka consumer
type KafkaConsumer struct {
	address   string
	consumers cmap.ConcurrentMap
	quitChan  chan struct{}
	*option
}
type kafkaConsumer struct {
	consumer *kafka.BrokerConsumer
	doneChan chan struct{}
}

//NewKafkaConsumer 初始化kafka Consumer
func NewKafkaConsumer(address string, opts ...Option) (mq *KafkaConsumer, err error) {
	mq = &KafkaConsumer{address: address, quitChan: make(chan struct{}, 0)}
	mq.consumers = cmap.New()
	for _, opt := range opts {
		opt(mq.option)
	}
	return
}

//Connect 连接到服务器
func (k *KafkaConsumer) Connect() error {
	return nil
}

//Consume 订阅消息
func (k *KafkaConsumer) Consume(queue string, call func(IMessage)) (err error) {
	_, cnsmr, _ := k.consumers.SetIfAbsentCb(queue, func(i ...interface{}) (interface{}, error) {
		c := &kafkaConsumer{}
		c.consumer = kafka.NewBrokerConsumer(k.address, queue, 0, 0, 1048576)
		c.doneChan = make(chan struct{})
		return c, nil
	})
	consumer := cnsmr.(*kafkaConsumer)
	msgQueue := make(chan *kafka.Message, 10)
	go consumer.consumer.ConsumeOnChannel(msgQueue, 10, k.quitChan)
	go func() {
	LOOP:
		for {
			select {
			case <-consumer.doneChan:
				close(k.quitChan)
			case msg, ok := <-msgQueue:
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
		close(consumer.doneChan)
	}
}

//Close 关闭当前连接
func (k *KafkaConsumer) Close() {
	close(k.quitChan)
}

type kafkaConsumerResolver struct {
}

func (s *kafkaConsumerResolver) Resolve(address string, opts ...Option) (MQConsumer, error) {
	return NewKafkaConsumer(address, opts...)
}
func init() {
	Register("kafka", &kafkaConsumerResolver{})
}
