package mq

import (
	"strings"
	"testing"
)

/*
{
    "type": "stomp",
    "address": "192.168.0.165:61613"
}
*/

var address = "192.168.0.165:61613"

func TestNewStompProducer(t *testing.T) {
	// 正常调用
	version := "1.1"
	producerConfig := ProducerConfig{Address: address, Version: version, Persistent: "persistent"}
	producer, err := NewStompProducer(producerConfig)
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	if len(producer.header) != 2 || !strings.EqualFold(producer.header[0], "accept-version") || !strings.EqualFold(producer.header[1], version) {
		t.Errorf("NewStompProducer fail : %+v", producer)
	}

	// 传入一个空ProducerConfig
	producer, err = NewStompProducer(ProducerConfig{})
	// if err == nil {
	// 	t.Error("test fail")
	// }
	// if len(producer.header) != 2 || !strings.EqualFold(producer.header[0], "accept-version") || !strings.EqualFold(producer.header[1], version) {
	// 	t.Errorf("NewStompProducer fail : %+v", producer)
	// }
}

func TestProducerConnect(t *testing.T) {
	// 正常连接到服务器
	version := "1.1"
	producerConfig := ProducerConfig{Address: address, Version: version, Persistent: "persistent"}
	producer, err := NewStompProducer(producerConfig)
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	err = producer.Connect()
	if err != nil {
		t.Errorf("Connect to servicer fail : %v", err)
	}

	// 传入空的ProducerConfig
	producer, err = NewStompProducer(ProducerConfig{})
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	err = producer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// producer本身就是空值
	producer = &StompProducer{}
	err = producer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// ip地址错误
	addr := "192.168.0.166:61613"
	producerConfig = ProducerConfig{Address: addr, Version: version, Persistent: "persistent"}
	producer, err = NewStompProducer(producerConfig)
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	err = producer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// 端口错误
	addr = "192.168.0.165:80"
	producerConfig = ProducerConfig{Address: addr, Version: version, Persistent: "persistent"}
	producer, err = NewStompProducer(producerConfig)
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	err = producer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// ip地址格式错误
	addr = "168.165:61613"
	producerConfig = ProducerConfig{Address: addr, Version: version, Persistent: "persistent"}
	producer, err = NewStompProducer(producerConfig)
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	err = producer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// address包含特殊字符
	addr = "！@#168.165:61613"
	producerConfig = ProducerConfig{Address: addr, Version: version, Persistent: "persistent"}
	producer, err = NewStompProducer(producerConfig)
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	err = producer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}

	// address为空字符串
	addr = ""
	producerConfig = ProducerConfig{Address: addr, Version: version, Persistent: "persistent"}
	producer, err = NewStompProducer(producerConfig)
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	err = producer.Connect()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
	}
}

func TestProducerSend(t *testing.T) {
	// 正常连接到服务器
	version := "1.1"
	producerConfig := ProducerConfig{Address: address, Version: version, Persistent: "persistent"}
	producer, err := NewStompProducer(producerConfig)
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}
	err = producer.Connect()
	if err != nil {
		t.Errorf("Connect to servicer fail : %v", err)
	}

	queue := "test_queue"
	msg := "test_msg"
	timeout := 10
	// 正常发送数据
	err = producer.Send(queue, msg, timeout)
	if err != nil {
		t.Errorf("Send mq fail : %v", err)
	}

	// 队列名为空字符串
	queue = ""
	msg = "test_msg"
	timeout = 10
	err = producer.Send(queue, msg, timeout)
	if err != nil {
		t.Errorf("Send mq fail : %v", err)
	}

	// 队列名包含特殊字符
	queue = "！@#￥×（……"
	msg = "test_msg"
	timeout = 10
	err = producer.Send(queue, msg, timeout)
	if err != nil {
		t.Errorf("Send mq fail : %v", err)
	}

	// 消息为空字符串
	queue = "test_queue"
	msg = ""
	timeout = 10
	err = producer.Send(queue, msg, timeout)
	if err != nil {
		t.Errorf("Send mq fail : %v", err)
	}

	// 消息包含特殊字符
	queue = "test_queue"
	msg = "！@#%￥！（……"
	timeout = 10
	err = producer.Send(queue, msg, timeout)
	if err != nil {
		t.Errorf("Send mq fail : %v", err)
	}

	// 超时为0
	queue = "test_queue"
	msg = "test_msg"
	timeout = 0
	err = producer.Send(queue, msg, timeout)
	if err != nil {
		t.Errorf("Send mq fail : %v", err)
	}

	// 超时为负数
	queue = "test_queue"
	msg = "test_msg"
	timeout = -100
	err = producer.Send(queue, msg, timeout)
	if err != nil {
		t.Errorf("Send mq fail : %v", err)
	}
}

func TestProducerClose(t *testing.T) {
	// 正常连接到服务器
	version := "1.1"
	producerConfig := ProducerConfig{Address: address, Version: version, Persistent: "persistent"}
	producer, err := NewStompProducer(producerConfig)
	if err != nil {
		t.Errorf("NewStompProducer fail : %v", err)
	}

	// conn为nil
	producer.Close()

	// 连接之后关闭
	err = producer.Connect()
	if err != nil {
		t.Errorf("Connect to servicer fail : %v", err)
	}
	producer.Close()
	if producer.conn.Connected() {
		t.Error("test fail")
	}

	// 没有连接就关闭
	producer.Close()
	if producer.conn.Connected() {
		t.Error("test fail")
	}
}
