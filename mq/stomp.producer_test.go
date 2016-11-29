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
	version := "version"
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

func TestConnect(t *testing.T) {
	// 正常连接到服务器
	version := "version"
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
}
