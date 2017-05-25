package memcache

import "github.com/bradfitz/gomemcache/memcache"

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
	v, err := c.Get(key)
	if err != nil {
		return err
	}
	data := &memcache.Item{Key: key, Value: []byte(v), Expiration: int32(expiresAt)}
	return c.client.Set(data)
}
