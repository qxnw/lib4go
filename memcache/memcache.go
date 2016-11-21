package memcache

import (
	"encoding/json"
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcacheClient memcache配置文件
type MemcacheClient struct {
	Servers []string `json:"servers"`
	client  *memcache.Client
}

// New 根据配置文件创建一个memcache连接
func New(config string) (m *MemcacheClient, err error) {
	m = &MemcacheClient{}
	err = json.Unmarshal([]byte(config), &m)
	if err != nil {
		err = errors.New("memcache配置文件有误:" + err.Error())
		return
	}
	m.client = memcache.New(m.Servers...)
	return
}

// Get 根据key获取memcache中的数据
func (c *MemcacheClient) Get(key string) string {
	data, err := c.client.Get(key)
	if err != nil {
		return ""
	}
	return string(data.Value)
}

// Add 添加数据到memcache中,如果memcache存在，则报错
func (c *MemcacheClient) Add(key string, value string, expiresAt int) error {
	data := &memcache.Item{Key: key, Value: []byte(value), Expiration: int32(expiresAt)}
	return c.client.Add(data)
}

// Set 更新数据到memcache中，没有则添加
func (c *MemcacheClient) Set(key string, value string, expiresAt int) error {
	data := &memcache.Item{Key: key, Value: []byte(value), Expiration: int32(expiresAt)}
	err := c.client.Set(data)
	return err
}

// Delete 删除memcache中的数据
func (c *MemcacheClient) Delete(key string) error {
	return c.client.Delete(key)
}

// Delay 延长数据在memcache中的时间(从现在开始计时)
func (c *MemcacheClient) Delay(key string, expiresAt int) error {
	v := c.Get(key)
	data := &memcache.Item{Key: key, Value: []byte(v), Expiration: int32(expiresAt)}
	return c.client.Set(data)
}
