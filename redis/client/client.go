package client

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/qxnw/lib4go/types"
)

type ClientConf struct {
	Address     string
	Sentinels   []string
	Password    string
	Db          int
	DialTimeout int
	RTimeout    int
	WTimeout    int
	PoolSize    int
}

//Client redis client
type Client struct {
	*redis.Client
}

//ClientOption 配置选项
type ClientOption func(*ClientConf)

//WithSentinels 设置哨兵服务器
func WithSentinels(sentinels []string) ClientOption {
	return func(o *ClientConf) {
		o.Sentinels = sentinels
	}
}

//WithPassword 设置服务器登录密码
func WithPassword(password string) ClientOption {
	return func(o *ClientConf) {
		o.Password = password
	}
}

//WithDB 设置数据库
func WithDB(db int) ClientOption {
	return func(o *ClientConf) {
		o.Db = db
	}
}

//WithDialTimeout 设置连接超时时长
func WithDialTimeout(timeout int) ClientOption {
	return func(o *ClientConf) {
		o.DialTimeout = timeout
	}
}

//WithRWTimeout 设置读写超时时长
func WithRTimeout(timeout int) ClientOption {
	return func(o *ClientConf) {
		o.RTimeout = timeout
	}
}

//WithWTimeout 设置读写超时时长
func WithWTimeout(timeout int) ClientOption {
	return func(o *ClientConf) {
		o.WTimeout = timeout
	}
}

//NewClient 构建客户端
func NewClient(addr string, option ...ClientOption) (r *Client, err error) {
	conf := &ClientConf{}
	for _, opt := range option {
		opt(conf)
	}
	return NewClientByConf(addr, conf)

}

//NewClientByJson 根据json构建failover客户端
func NewClientByJson(j string) (r *Client, err error) {
	conf := &ClientConf{}
	err = json.Unmarshal([]byte(j), &conf)
	if err != nil {
		return nil, err
	}
	return NewClientByConf(conf.Address, conf)
}

//NewClientByConf 根据配置对象构建客户端
func NewClientByConf(addr string, conf *ClientConf) (r *Client, err error) {
	conf.DialTimeout = types.DecodeInt(conf.DialTimeout, 0, 3, conf.DialTimeout)
	conf.RTimeout = types.DecodeInt(conf.RTimeout, 0, 3, conf.RTimeout)
	conf.WTimeout = types.DecodeInt(conf.WTimeout, 0, 3, conf.WTimeout)
	conf.PoolSize = types.DecodeInt(conf.PoolSize, 0, 3, conf.PoolSize)
	client := &Client{}
	client.Client = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     conf.Password,
		DB:           conf.Db,
		DialTimeout:  time.Duration(conf.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.RTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WTimeout) * time.Second,
		PoolSize:     conf.PoolSize,
	})
	_, err = client.Client.Ping().Result()
	return
}
