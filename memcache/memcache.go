package memcache

import "github.com/bradfitz/gomemcache/memcache"
import "time"

// MemcacheClient memcache配置文件
type MemcacheClient struct {
	servers []string
	client  *memcache.Client
}

// New 根据配置文件创建一个memcache连接
func New(addrs []string) (m *MemcacheClient, err error) {
	m = &MemcacheClient{servers: addrs}
	m.client = memcache.New(addrs...)
	return
}

// Get 根据key获取memcache中的数据
func (c *MemcacheClient) Get(key string) (string, error) {
	data, err := c.client.Get(key)
	if err != nil {
		return "", err
	}
	return string(data.Value), nil
}

//Decrement 增加变量的值
func (c *MemcacheClient) Decrement(key string, delta uint64) (n uint64, err error) {
	return c.client.Decrement(key, delta)
}

//Increment 减少变量的值
func (c *MemcacheClient) Increment(key string, delta uint64) (n uint64, err error) {
	return c.client.Increment(key, delta)
}

//Gets 获取多条数据
func (c *MemcacheClient) Gets(key ...string) (r []string, err error) {
	data, err := c.client.GetMulti(key)
	if err != nil {
		return nil, err
	}
	r = make([]string, len(data))
	for _, v := range key {
		r = append(r, string(data[v].Value))
	}
	return
}

// Add 添加数据到memcache中,如果memcache存在，则报错
func (c *MemcacheClient) Add(key string, value string, expiresAt int) error {
	expires := time.Now().Add(time.Duration(expiresAt) * time.Second).Unix()
	if expiresAt == 0 {
		expires = 0
	}
	data := &memcache.Item{Key: key, Value: []byte(value), Expiration: int32(expires)}
	return c.client.Add(data)
}

// Set 更新数据到memcache中，没有则添加
func (c *MemcacheClient) Set(key string, value string, expiresAt int) error {
	expires := time.Now().Add(time.Duration(expiresAt) * time.Second).Unix()
	if expiresAt == 0 {
		expires = 0
	}
	data := &memcache.Item{Key: key, Value: []byte(value), Expiration: int32(expires)}
	err := c.client.Set(data)
	return err
}

// Delete 删除memcache中的数据
func (c *MemcacheClient) Delete(key string) error {
	return c.client.Delete(key)
}

// Delay 延长数据在memcache中的时间
func (c *MemcacheClient) Delay(key string, expiresAt int) error {
	return c.client.Touch(key, int32(expiresAt))
}

// DeleteAll 删除所有缓存数据
func (c *MemcacheClient) DeleteAll() error {
	return c.client.DeleteAll()
}
