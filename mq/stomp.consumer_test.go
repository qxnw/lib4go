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

// TestNewStompConsumer 测试创建一个消费者对象
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

// TestConsumerConnect 测试消费者对象连接到服务器
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

// TestConsume 测试注册一个消费信息
func TestConsume(t *testing.T) {
	normalAccount := 0
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
			normalAccount++
			if !strings.EqualFold(m.GetMessage(), consumerMsg) {
				t.Errorf("test fail actual:%s, except:%s", m.GetMessage(), consumerMsg)
			}
			// 确定消息
			m.Ack()
		})
		if err != nil {
			t.Errorf("test fail: %v", err)
		}
	}()
	// 发送一个队列，测试回调
	SendMsg(consumerQueue, consumerMsg)
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

	// 回调函数为nil
	errQueue := "err_queue"
	err = consumer.Consume(errQueue, nil)
	if !strings.EqualFold(err.Error(), "回调函数不能为nil") {
		t.Error("test fail")
	}
	// 发送两个队列，测试回调
	SendMsg(errQueue, consumerMsg)

	// 队列名为空字符串
	go func() {
		err = consumer.Consume("", func(m IMessage) {
			if !strings.EqualFold(m.GetMessage(), consumerMsg) {
				t.Errorf("test fail actual:%s, except:%s", m.GetMessage(), consumerMsg)
			}
		})
		if !strings.EqualFold(err.Error(), "队列名字不能为空") {
			t.Error("test fail")
		}
	}()

	// 等待能接收到消息
	time.Sleep(time.Second * 2)
	// 判断是否进入回调函数的次数
	if normalAccount != 2 {
		t.Errorf("test fail normalAccount:%d", normalAccount)
	}
}

// TestUnConsume 测试取消注册消费信息
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
	// 再次注册，判断是否取消
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

	// 取消一个不存在的队列
	consumer.UnConsume("afdasd")

	// 队列名字为空字符串
	consumer.UnConsume("")

	// 队列名字有特殊字符
	consumer.UnConsume("\\//!@#$!")
}

// TestConsumeClose 测试关闭连接
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

// TestSpecialSituation 测试特殊情况
func TestSpecialSituation(t *testing.T) {
	// 先注册一个消费信息，然后模拟断网，然后网络恢复，然后发送数据，看数据是否能收到
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

	specialAccount := 0

	// 注册一个消费信息
	go func() {
		err = consumer.Consume("sepcialQueue", func(m IMessage) {
			fmt.Println("进入回调函数SpecialSituation,收到的数据Message:", m.GetMessage())
			if !strings.EqualFold(m.GetMessage(), consumerMsg) {
				t.Errorf("test fail actual:%s, except:%s", m.GetMessage(), consumerMsg)
			}
			specialAccount++
			// 确定消息
			m.Ack()
		})
		t.Log(err)
		if err != nil {
			t.Errorf("test fail: %v", err)
		}
	}()

	// 手动断网
	fmt.Println("开始关闭网络")
	time.Sleep(time.Second * 5)
	fmt.Println("网络已经关闭……等待60s")
	time.Sleep(time.Second * 60)
	fmt.Println("开始启动网络")
	// 恢复网络
	time.Sleep(time.Second * 10)
	fmt.Println("网络已启用")

	// 发送消息
	SendMsg("sepcialQueue", consumerMsg)

	// 等待接收消息
	fmt.Println("开始等待处理结果")
	time.Sleep(time.Second * 2)
	if specialAccount != 1 {
		t.Errorf("test fail specialAccount : %d", specialAccount)
	}
}
