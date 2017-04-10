package redis

import (
	"time"

	"github.com/go-redis/redis"
)
import "fmt"

const (
	CLUSTER      = "cluster"
	RING         = "ring"
	MasterSalver = "master-salver"
	Standalone   = "standalone"
)

//Config refids config
type Config struct {
	Address     string   `json:"address"`
	Addrs       []string `json:"addrs"`
	Master      string   `json:"master"`
	Slaver      []string `json:"slavers"`
	Type        string   `json:"type"`
	Password    string   `json:"password"`
	Db          int      `json:"db"`
	PoolSize    int      `json:"pool-size"`
	DialTimeout int      `json:"dial-timeout"`
	RWTimeout   int      `json:"rw-timeout"`
	Failover    bool     `json:"failover"`
}

//RedisClient redis client
type RedisClient struct {
	client *redis.Client
}

//NewClient 构建客户端
func NewClient(conf Config) (r *RedisClient, err error) {
	if conf.RWTimeout == 0 {
		conf.RWTimeout = 30
	}
	if conf.DialTimeout == 0 {
		conf.DialTimeout = 30
	}
	if conf.PoolSize == 0 {
		conf.PoolSize = 10
	}
	r = &RedisClient{}
	if conf.Failover {
		r.client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.Master,
			SentinelAddrs: conf.Slaver,
			Password:      conf.Password,
			DB:            conf.Db,
			DialTimeout:   time.Duration(conf.DialTimeout) * time.Second,
			ReadTimeout:   time.Duration(conf.RWTimeout) * time.Second,
			WriteTimeout:  time.Duration(conf.RWTimeout) * time.Second,
			PoolSize:      conf.PoolSize,
			PoolTimeout:   time.Duration(conf.RWTimeout) * time.Second,
		})

	} else {
		r.client = redis.NewClient(&redis.Options{
			Addr:         conf.Address,
			Password:     conf.Password,
			DB:           conf.Db,
			DialTimeout:  time.Duration(conf.DialTimeout) * time.Second,
			ReadTimeout:  time.Duration(conf.RWTimeout) * time.Second,
			WriteTimeout: time.Duration(conf.RWTimeout) * time.Second,
			PoolSize:     conf.PoolSize,
			PoolTimeout:  time.Duration(conf.RWTimeout) * time.Second,
		})
	}
	_, err = r.client.Ping().Result()
	return
}

// Get 根据key获取数据
func (r *RedisClient) Get(key string) (rs string, err error) {
	rs, err = r.client.Get(key).Result()
	if err == redis.Nil {
		err = fmt.Errorf("redis key(%s) does not exists", key)
	}
	return
}

//Set 设置缓存值
func (r *RedisClient) Set(key string, value string, expiresAt int) error {
	err := r.client.Set(key, value, time.Second*time.Duration(expiresAt)).Err()
	return err
}
