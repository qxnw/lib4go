package kafka

import (
	"errors"
	"sync"
	"time"

	"github.com/jdamick/kafka"
	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/mq"
)

//KafkaProducer Producer
type KafkaProducer struct {
	address    string
	messages   chan *mq.ProcuderMessage
	backupMsg  chan *mq.ProcuderMessage
	queues     cmap.ConcurrentMap
	connecting bool
	closeCh    chan struct{}
	done       bool
	lk         sync.Mutex
	header     []string
	*mq.OptionConf
}
type kafkaProducer struct {
	producer *kafka.BrokerPublisher
	msgQueue chan *kafka.Message
}

//NewStompProducer 创建新的producer
func NewKafkaProducer(address string, opts ...mq.Option) (producer *KafkaProducer, err error) {
	producer = &KafkaProducer{address: address}
	producer.OptionConf = &mq.OptionConf{Logger: logger.GetSession("mq.producer", logger.CreateSession())}
	producer.messages = make(chan *mq.ProcuderMessage, 100)
	producer.backupMsg = make(chan *mq.ProcuderMessage, 100)
	producer.closeCh = make(chan struct{})
	return
}

//Connect  循环连接服务器
func (producer *KafkaProducer) Connect() error {
	go producer.sendLoop()
	return nil
}

//sendLoop 循环发送消息
func (producer *KafkaProducer) sendLoop() {
	if producer.done {
		producer.disconnect()
		return
	}
	if producer.Retry {
	Loop1:
		for {
			select {
			case msg, ok := <-producer.backupMsg:
				if !ok {
					break Loop1
				}
				pd, ok := producer.queues.Get(msg.Queue)
				if !ok {
					select {
					case producer.backupMsg <- msg:
					default:
						producer.Logger.Errorf("消息无法放入备份队列(%s):%s", msg.Queue, msg.Data)
					}
					producer.Logger.Errorf("消息无法从缓存中获取producer:%s,%s", msg.Queue, msg.Data)
					continue
				}
				producerConn := pd.(*kafka.BrokerPublisher)
				_, err := producerConn.Publish(kafka.NewMessage([]byte(msg.Data)))
				if err != nil {
					select {
					case producer.backupMsg <- msg:
					default:
						producer.Logger.Errorf("消息无法放入备份队列(%s):%s", msg.Queue, msg.Data)
					}
				}
			case msg, ok := <-producer.messages:
				if !ok {
					break Loop1
				}
				pd, ok := producer.queues.Get(msg.Queue)
				if !ok {
					select {
					case producer.backupMsg <- msg:
					default:
						producer.Logger.Errorf("消息无法放入备份队列(%s):%s", msg.Queue, msg.Data)
					}
					producer.Logger.Errorf("消息无法从缓存中获取producer:%s,%s", msg.Queue, msg.Data)
					continue
				}
				producerConn := pd.(*kafka.BrokerPublisher)
				_, err := producerConn.Publish(kafka.NewMessage([]byte(msg.Data)))
				if err != nil {
					select {
					case producer.backupMsg <- msg:
					default:
						producer.Logger.Errorf("消息无法放入备份队列(%s):%s", msg.Queue, msg.Data)
					}
				}
			}
		}
	} else {
	Loop2:
		for {
			select {
			case msg, ok := <-producer.messages:
				if !ok {
					break Loop2
				}
				pd, ok := producer.queues.Get(msg.Queue)
				if !ok {
					select {
					case producer.backupMsg <- msg:
					default:
						producer.Logger.Errorf("消息无法放入备份队列(%s):%s", msg.Queue, msg.Data)
					}
					producer.Logger.Errorf("消息无法从缓存中获取producer:%s,%s", msg.Queue, msg.Data)
					continue
				}
				producerConn := pd.(*kafka.BrokerPublisher)
				_, err := producerConn.Publish(kafka.NewMessage([]byte(msg.Data)))
				if err != nil {
					select {
					case producer.backupMsg <- msg:
					default:
						producer.Logger.Errorf("消息无法放入备份队列(%s):%s", msg.Queue, msg.Data)
					}
				}
			}
		}
	}
	if producer.done { //关闭连接
		producer.disconnect()
		return
	}
}
func (producer *KafkaProducer) disconnect() {

}

//GetBackupMessage 获取备份数据
func (producer *KafkaProducer) GetBackupMessage() chan *mq.ProcuderMessage {
	return producer.backupMsg
}

//Send 发送消息
func (producer *KafkaProducer) Send(queue string, msg string, timeout time.Duration) (err error) {
	if producer.done {
		return errors.New("mq producer 已关闭")
	}
	producer.queues.SetIfAbsentCb(queue, func(i ...interface{}) (interface{}, error) {
		c := &kafkaProducer{}
		c.producer = kafka.NewBrokerPublisher(producer.address, queue, 0)
		c.msgQueue = make(chan *kafka.Message, 10)
		return c, nil
	})

	pm := &mq.ProcuderMessage{Queue: queue, Data: msg, Timeout: timeout}
	select {
	case producer.messages <- pm:
		return nil
	default:
		return errors.New("producer无法连接，消息发送失败")
	}
}

//Close 关闭当前连接
func (producer *KafkaProducer) Close() {
	producer.done = true
	close(producer.closeCh)
	close(producer.messages)
	close(producer.backupMsg)
}

type kafkaProducerResolver struct {
}

func (s *kafkaProducerResolver) Resolve(address string, opts ...mq.Option) (mq.MQProducer, error) {
	return NewKafkaProducer(address, opts...)
}
func init() {
	mq.RegisterProducer("kafka", &kafkaProducerResolver{})
}
