package failover

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/qxnw/lib4go/types"
)

type FailoverConf struct {
	Address     string
	Sentinels   []string
	Password    string
	Db          int
	DialTimeout int
	RTimeout    int
	WTimeout    int
	PoolSize    int
}

//FailoverClient redis client
type FailoverClient struct {
	*redis.Client
}

//Option 配置选项
type FailoverOption func(*FailoverConf)

//WithSentinels 设置哨兵服务器
func WithSentinels(sentinels []string) FailoverOption {
	return func(o *FailoverConf) {
		o.Sentinels = sentinels
	}
}

//WithPassword 设置服务器登录密码
func WithPassword(password string) FailoverOption {
	return func(o *FailoverConf) {
		o.Password = password
	}
}

//WithDB 设置数据库
func WithDB(db int) FailoverOption {
	return func(o *FailoverConf) {
		o.Db = db
	}
}

//WithDialTimeout 设置连接超时时长
func WithDialTimeout(timeout int) FailoverOption {
	return func(o *FailoverConf) {
		o.DialTimeout = timeout
	}
}

//WithRWTimeout 设置读写超时时长
func WithRTimeout(timeout int) FailoverOption {
	return func(o *FailoverConf) {
		o.RTimeout = timeout
	}
}

//WithWTimeout 设置读写超时时长
func WithWTimeout(timeout int) FailoverOption {
	return func(o *FailoverConf) {
		o.WTimeout = timeout
	}
}

//NewClient 构建客户端
func NewFailoverClient(masterAddr string, option ...FailoverOption) (r *FailoverClient, err error) {
	conf := &FailoverConf{}
	for _, opt := range option {
		opt(conf)
	}
	return NewFailoverClientByConf(masterAddr, conf)

}

//NewFailoverClientByJson 根据json构建failover客户端
func NewFailoverClientByJson(j string) (r *FailoverClient, err error) {
	conf := &FailoverConf{}
	err = json.Unmarshal([]byte(j), &conf)
	if err != nil {
		return nil, err
	}
	return NewFailoverClientByConf(conf.Address, conf)
}

//NewFailoverClientByConf 根据配置对象构建客户端
func NewFailoverClientByConf(masterAddr string, conf *FailoverConf) (r *FailoverClient, err error) {
	conf.DialTimeout = types.DecodeInt(conf.DialTimeout, 0, 3, conf.DialTimeout)
	conf.RTimeout = types.DecodeInt(conf.RTimeout, 0, 3, conf.RTimeout)
	conf.WTimeout = types.DecodeInt(conf.WTimeout, 0, 3, conf.WTimeout)
	conf.PoolSize = types.DecodeInt(conf.PoolSize, 0, 3, conf.PoolSize)
	client := &FailoverClient{}
	client.Client = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterAddr,
		SentinelAddrs: conf.Sentinels,
		Password:      conf.Password,
		DB:            conf.Db,
		DialTimeout:   time.Duration(conf.DialTimeout) * time.Second,
		ReadTimeout:   time.Duration(conf.RTimeout) * time.Second,
		WriteTimeout:  time.Duration(conf.WTimeout) * time.Second,
		PoolSize:      conf.PoolSize,
	})
	_, err = client.Client.Ping().Result()
	return
}
