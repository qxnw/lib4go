package mq

import "github.com/qxnw/lib4go/logger"

type OptionConf struct {
	Logger     logger.ILogger
	Address    string `json:"address"`
	Version    string `json:"version"`
	Persistent string `json:"persistent"`
	Ack        string `json:"ack"`
	Retry      bool   `json:"retry"`
	Key        string `json:"key"`
}

//Option 配置选项
type Option func(*OptionConf)

//WithConf 根据配置文件初始化
func WithConf(conf *OptionConf) Option {
	return func(o *OptionConf) {
		o = conf
	}
}

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

//WithSignKey 设置数据加密key
func WithSignKey(key string) Option {
	return func(o *OptionConf) {
		o.Key = key
	}
}
