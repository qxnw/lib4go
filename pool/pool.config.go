package pool

import (
	"errors"
	"time"
)

var (
	//ErrClosed 连接池已经关闭Error
	ErrClosed = errors.New("pool is closed")
)

//IPool 基本方法
type IPool interface {
	Get() (interface{}, error)

	Put(interface{}) error

	Close(interface{}) error

	Release()

	Len() int
}

//PoolConfigOptions 连接池相关配置
type PoolConfigOptions struct {
	//连接池中拥有的最小连接数
	InitialCap int
	//连接池中拥有的最大的连接数
	MaxCap int
	//生成连接的方法
	Factory func() (interface{}, error)
	//关闭链接的方法
	Close func(interface{}) error
	//链接最大空闲时间，超过该事件则将失效
	IdleTimeout time.Duration
}
