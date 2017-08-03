package mq

import "github.com/qxnw/lib4go/logger"

type OptionConf struct {
	Logger     logger.ILogger
	Version    string
	Persistent string
	Ack        string
	Retry      bool
}

//Option 配置选项
type Option func(*OptionConf)

//WithLogger 设置日志记录组件
func WithLogger(logger logger.ILogger) Option {
	return func(o *OptionConf) {
		o.Logger = logger
	}
}

//WithVersion 设置版本号
func WithVersion(version string) Option {
	return func(o *OptionConf) {
		o.Version = version
	}
}

//WithPersistent 设置数据模式
func WithPersistent(persistent string) Option {
	return func(o *OptionConf) {
		o.Persistent = persistent
	}
}

//WithAck 设置客户端确认模式
func WithAck(ack string) Option {
	return func(o *OptionConf) {
		o.Ack = ack
	}
}

//WithRetrySend 发送失败后重试
func WithRetrySend(b bool) Option {
	return func(o *OptionConf) {
		o.Retry = b
	}
}
