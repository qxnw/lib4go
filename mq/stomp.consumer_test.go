package mq

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gmallard/stompngo"
)

var consumerQueue = "queue_test"
var consumerMsg = "msg_test"
var consumerTimeOut = 10

func TestNewStompConsumer(t *testing.T) {
	// 正常调用
	version := "1.1"
	consumerConfig := ConsumerConfig{Address: address, Version: version, Persistent: "persistent"}
	consumer, err := NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}
	if len(consumer.header) != 2 || !strings.EqualFold(consumer.header[0], "accept-version") || !strings.EqualFold(consumer.header[1], version) {
		t.Errorf("NewStompConsumer fail : %+v", consumer)
	}

	// 传入一个空ConsumerConfig
	consumer, err = NewStompConsumer(ConsumerConfig{})
	// if err == nil {
	// 	t.Error("test fail")
	// }
	// if len(producer.header) != 2 || !strings.EqualFold(producer.header[0], "accept-version") || !strings.EqualFold(producer.header[1], version) {
	// 	t.Errorf("NewStompConsumer fail : %+v", producer)
	// }
}

func TestConsumerConnect(t *testing.T) {
	// 正常连接到服务器
	version := "1.1"
	consumerConfig := ConsumerConfig{Address: address, Version: version, Persistent: "persistent"}
	consumer, err := NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}
	err = consumer.Connect()
	if err != nil {
		t.Errorf("Connect to servicer fail : %v", err)
	}

	// 传入空的ProducerConfig
	consumer, err = NewStompConsumer(ConsumerConfig{})
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	err = consumer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// producer本身就是空值
	consumer = &StompConsumer{}
	err = consumer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// ip地址错误
	addr := "192.168.0.166:61613"
	consumerConfig = ConsumerConfig{Address: addr, Version: version, Persistent: "persistent"}
	consumer, err = NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}
	err = consumer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// 端口错误
	addr = "192.168.0.165:80"
	consumerConfig = ConsumerConfig{Address: addr, Version: version, Persistent: "persistent"}
	consumer, err = NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}
	err = consumer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// ip地址格式错误
	addr = "168.165:61613"
	consumerConfig = ConsumerConfig{Address: addr, Version: version, Persistent: "persistent"}
	consumer, err = NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}
	err = consumer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// address包含特殊字符
	addr = "！@#168.165:61613"
	consumerConfig = ConsumerConfig{Address: addr, Version: version, Persistent: "persistent"}
	consumer, err = NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}
	err = consumer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// address为空字符串
	addr = ""
	consumerConfig = ConsumerConfig{Address: addr, Version: version, Persistent: "persistent"}
	consumer, err = NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}
	err = consumer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}
}

// SendMsg 测试发送一个队列
func SendMsg(queue, msg string) {
	// 正常连接到服务器
	version := "1.1"
	producerConfig := ProducerConfig{Address: address, Version: version, Persistent: "persistent"}
	producer, err := NewStompProducer(producerConfig)
	if err != nil {
		return
	}
	if err = producer.Connect(); err != nil {
		return
	}
	header := stompngo.Headers{"destination", queue, "persistent", producer.config.Persistent}
	if consumerTimeOut > 0 {
		header = stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", "queue", queue),
			"persistent", producer.config.Persistent, "expires",
			fmt.Sprintf("%d000", time.Now().Add(time.Second*time.Duration(consumerTimeOut)).Unix())}
	}
	err = producer.conn.Send(header, msg)
	return
}

func TestConsume(t *testing.T) {
	// 正常连接到服务器
	version := "1.1"
	consumerConfig := ConsumerConfig{Address: address, Version: version, Persistent: "persistent"}
	consumer, err := NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}

	// 回调函数正常
	go func() {
		err = consumer.Consume(consumerQueue, func(m IMessage) {
			t.Log("进入回调函数")
			if !strings.EqualFold(m.GetMessage(), consumerMsg) {
				t.Errorf("test fail actual:%s, except:%s", m.GetMessage(), consumerMsg)
			}
		})
		if err != nil {
			t.Errorf("test fail: %v", err)
		}
	}()
	// 发送一个队列，测试回调
	SendMsg(consumerQueue, consumerMsg)

	// 重复订阅
	go func() {
		err = consumer.Consume(consumerQueue, func(m IMessage) {
			t.Log("进入回调函数")
			if !strings.EqualFold(m.GetMessage(), consumerMsg) {
				t.Errorf("test fail actual:%s, except:%s", m.GetMessage(), consumerMsg)
			}
		})
		t.Log(err)
		if !strings.EqualFold(err.Error(), fmt.Sprintf("重复订阅消息:%s", consumerQueue)) {
			t.Errorf("test fail :%v", err)
		}
	}()

	// // 回调函数为nil
	// errQueue := "err_queue"
	// err = consumer.Consume(errQueue, nil)
	// if err == nil {
	// 	t.Error("test fail")
	// }
	// // 发送一个队列，测试回调
	// SendMsg(errQueue, consumerMsg)
	//
	// // 队列名为空字符串
	// go func() {
	// 	err = consumer.Consume("", func(m IMessage) {
	// 		if !strings.EqualFold(m.GetMessage(), consumerMsg) {
	// 			t.Errorf("test fail actual:%s, except:%s", m.GetMessage(), consumerMsg)
	// 		}
	// 	})
	// 	if err != nil {
	// 		t.Errorf("test fail: %v", err)
	// 	}
	// }()
	// SendMsg("", consumerMsg)
}

func TestUnConsume(t *testing.T) {
	// 正常连接到服务器
	version := "1.1"
	consumerConfig := ConsumerConfig{Address: address, Version: version, Persistent: "persistent"}
	consumer, err := NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}

	// 取消一个存在的队列
	consumer.UnConsume(consumerQueue)

	// 取消一个不存在的队列
	consumer.UnConsume("afdasd")

	// 队列名字为空字符串
	consumer.UnConsume("")

	// 队列名字有特殊字符
	consumer.UnConsume("\\//!@#$!")
}

func TestConsumeClose(t *testing.T) {
	// 正常连接到服务器
	version := "1.1"
	consumerConfig := ConsumerConfig{Address: address, Version: version, Persistent: "persistent"}
	consumer, err := NewStompConsumer(consumerConfig)
	if err != nil {
		t.Errorf("NewStompConsumer fail : %v", err)
	}

	// conn为nil
	consumer.Close()

	// 连接之后关闭
	err = consumer.Connect()
	if err != nil {
		t.Errorf("Connect to servicer fail : %v", err)
	}
	consumer.Close()
	if consumer.conn.Connected() {
		t.Error("test fail")
	}

	// 没有连接就关闭
	consumer.Close()
	if consumer.conn.Connected() {
		t.Error("test fail")
	}
}
