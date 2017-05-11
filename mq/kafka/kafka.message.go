package kafka

import "github.com/jdamick/kafka"

type KafkaMessage struct {
	msg     *kafka.Message
	Message string
}

//Ack
func (m *KafkaMessage) Ack() {
	//m.s.conn.Ack(m.msg.Headers)
}
func (m *KafkaMessage) Nack() {
	//m.s.conn.Nack(m.msg.Headers)
}
func (m *KafkaMessage) GetMessage() string {
	return m.Message
}

//NewMessage
func NewKafkaMessage(msg *kafka.Message) *KafkaMessage {
	//return &StompMessage{s: s, msg: msg, Message: string(msg.Body)}
	return &KafkaMessage{msg: msg, Message: msg.PayloadString()}
}
