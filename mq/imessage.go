package mq

type IMessage interface {
	Ack()
	Nack()
	GetMessage() string
}
