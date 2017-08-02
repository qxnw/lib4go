package cluster

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/qxnw/lib4go/types"
)

//ClusterClient redis client
type ClusterClient struct {
	*redis.ClusterClient
}

//ClusterConf 集群配置参数
type ClusterConf struct {
	Addrs        []string
	MaxRedirects int
	ReadOnly     bool
	MaxRetries   int
	Password     string
	DialTimeout  int
	RTimeout     int
	WTimeout     int
	PoolSize     int
}

//ClusterOption 配置选项
type ClusterOption func(*ClusterConf)

//WithPassword 设置服务器登录密码
func WithPassword(password string) ClusterOption {
	return func(o *ClusterConf) {
		o.Password = password
	}
}

//WithDialTimeout 设置连接超时时长
func WithDialTimeout(timeout int) ClusterOption {
	return func(o *ClusterConf) {
		o.DialTimeout = timeout
	}
}

//WithRWTimeout 设置读写超时时长
func WithRTimeout(timeout int) ClusterOption {
	return func(o *ClusterConf) {
		o.RTimeout = timeout
	}
}

//WithWTimeout 设置读写超时时长
func WithWTimeout(timeout int) ClusterOption {
	return func(o *ClusterConf) {
		o.WTimeout = timeout
	}
}

//NewClusterClient 构建客户端
func NewClusterClient(addrs []string, option ...ClusterOption) (r *ClusterClient, err error) {
	conf := &ClusterConf{}
	for _, opt := range option {
		opt(conf)
	}
	return NewClusterClientByConf(addrs, conf)

}

//NewClusterClientByJson 根据json构建failover客户端
func NewClusterClientByJson(j string) (r *ClusterClient, err error) {
	conf := &ClusterConf{}
	err = json.Unmarshal([]byte(j), &conf)
	if err != nil {
		return nil, err
	}
	return NewClusterClientByConf(conf.Addrs, conf)
}

//NewClusterClientByConf 根据配置对象构建客户端
func NewClusterClientByConf(Addrs []string, conf *ClusterConf) (client *ClusterClient, err error) {
	conf.DialTimeout = types.DecodeInt(conf.DialTimeout, 0, 3, conf.DialTimeout)
	conf.RTimeout = types.DecodeInt(conf.RTimeout, 0, 3, conf.RTimeout)
	conf.WTimeout = types.DecodeInt(conf.WTimeout, 0, 3, conf.WTimeout)
	conf.PoolSize = types.DecodeInt(conf.PoolSize, 0, 3, conf.PoolSize)
	client = &ClusterClient{}
	client.ClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        Addrs,
		MaxRedirects: conf.MaxRedirects,
		Password:     conf.Password,
		DialTimeout:  time.Duration(conf.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.RTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WTimeout) * time.Second,
		PoolSize:     conf.PoolSize,
	})
	_, err = client.ClusterClient.Ping().Result()
	return
}
