package mq

import "github.com/qxnw/lib4go/logger"

type option struct {
	logger     *logger.Logger
	version    string
	persistent string
	ack        string
	retry      bool
}

//Option 配置选项
type Option func(*option)

//WithLogger 设置日志记录组件
func WithLogger(logger *logger.Logger) Option {
	return func(o *option) {
		o.logger = logger
	}
}

//WithVersion 设置版本号
func WithVersion(version string) Option {
	return func(o *option) {
		o.version = version
	}
}

//WithPersistent 设置数据模式
func WithPersistent(persistent string) Option {
	return func(o *option) {
		o.persistent = persistent
	}
}

//WithAck 设置客户端确认模式
func WithAck(ack string) Option {
	return func(o *option) {
		o.ack = ack
	}
}

//WithRetrySend 发送失败后重试
func WithRetrySend(b bool) Option {
	return func(o *option) {
		o.retry = b
	}
}
