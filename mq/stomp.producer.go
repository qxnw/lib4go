package mq

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gmallard/stompngo"
	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/logger"
)

type procuderMessage struct {
	headers []string
	queue   string
	data    string
	timeout time.Duration
}

//StompProducer Producer
type StompProducer struct {
	address    string
	conn       *stompngo.Connection
	messages   chan *procuderMessage
	backupMsg  chan *procuderMessage
	queues     cmap.ConcurrentMap
	connecting bool
	closeCh    chan struct{}
	done       bool
	lk         sync.Mutex
	header     []string
	*option
}

//NewStompProducer 创建新的producer
func NewStompProducer(address string, opts ...Option) (producer *StompProducer, err error) {
	producer = &StompProducer{address: address}
	producer.option = &option{logger: logger.GetSession("mq.producer", logger.CreateSession())}
	producer.messages = make(chan *procuderMessage, 100)
	producer.backupMsg = make(chan *procuderMessage, 100)
	for _, opt := range opts {
		opt(producer.option)
	}
	if strings.EqualFold(producer.option.version, "") {
		producer.option.version = "1.1"
	}
	if strings.EqualFold(producer.option.persistent, "") {
		producer.option.persistent = "true"
	}
	if strings.EqualFold(producer.option.ack, "") {
		producer.option.ack = "client-individual"
	}
	producer.header = stompngo.Headers{"accept-version", producer.option.version}
	return
}

//Connect  循环连接服务器
func (producer *StompProducer) Connect() error {
	err := producer.connectOnce()
	if err == nil {
		return nil
	}
	go func() {
	START:
		for {
			select {
			case <-producer.closeCh:
				break START
			case <-time.After(time.Second * 3):
				err = producer.connectOnce()
				if err == nil {
					break START
				}
				producer.option.logger.Error(err)
			}
		}
	}()
	return nil
}

//sendLoop 循环发送消息
func (producer *StompProducer) sendLoop() {
	if producer.done {
		producer.disconnect()
		return
	}
	if producer.retry {
	Loop1:
		for {
			select {
			case msg, ok := <-producer.backupMsg:
				if !ok {
					break Loop1
				}
				err := producer.conn.Send(msg.headers, msg.data)
				if err != nil {
					select {
					case producer.backupMsg <- msg:
					default:
						producer.logger.Errorf("消息无法放入备份队列(%s):%s", msg.queue, msg.data)
					}
					break Loop1
				}
			case msg, ok := <-producer.messages:
				if !ok {
					break Loop1
				}
				err := producer.conn.Send(msg.headers, msg.data)
				if err != nil {
					select {
					case producer.backupMsg <- msg:
					default:
						producer.logger.Errorf("消息无法放入备份队列(%s):%s", msg.queue, msg.data)
					}
					break Loop1
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
				err := producer.conn.Send(msg.headers, msg.data)
				if err != nil {
					select {
					case producer.backupMsg <- msg:
					default:
						producer.logger.Errorf("消息无法放入备份队列(%s):%s", msg.queue, msg.data)
					}
					break Loop2
				}
			}
		}
	}
	if producer.done { //关闭连接
		producer.disconnect()
		return
	}
	producer.reconnect()
}
func (producer *StompProducer) disconnect() {
	if producer.conn == nil || !producer.conn.Connected() {
		return
	}
	producer.conn.Disconnect(stompngo.Headers{})
	return
}

//GetBackupMessage 获取备份数据
func (producer *StompProducer) GetBackupMessage() chan *procuderMessage {
	return producer.backupMsg
}

//reconnect 自动重连
func (producer *StompProducer) reconnect() {
	producer.conn.Disconnect(stompngo.Headers{})
	producer.Connect()
}

//ConnectOnce 连接到服务器
func (producer *StompProducer) connectOnce() (err error) {
	if producer.connecting {
		return nil
	}
	producer.lk.Lock()
	defer producer.lk.Unlock()
	if producer.connecting {
		return nil
	}
	producer.connecting = true
	defer func() {
		producer.connecting = false
	}()
	con, err := net.Dial("tcp", producer.address)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	producer.conn, err = stompngo.Connect(con, producer.header)
	if err != nil {
		return fmt.Errorf("mq 无法连接到MQ:%v", err)
	}
	go producer.sendLoop()
	return nil
}

//Send 发送消息
func (producer *StompProducer) Send(queue string, msg string, timeout time.Duration) (err error) {
	if producer.done {
		return errors.New("mq producer 已关闭")
	}
	pm := &procuderMessage{queue: queue, data: msg, timeout: timeout}
	pm.headers = make([]string, 0, len(producer.header)+2)
	copy(pm.headers, producer.header)
	pm.headers = append(pm.headers, "destination", "/queue/"+queue)
	select {
	case producer.messages <- pm:
		return nil
	default:
		return errors.New("producer无法连接，消息发送失败")
	}
}

//Close 关闭当前连接
func (producer *StompProducer) Close() {
	producer.done = true
	close(producer.closeCh)
	close(producer.messages)
	close(producer.backupMsg)
}
