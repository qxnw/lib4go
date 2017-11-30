package redis

import "github.com/go-redis/redis"

//RedisMessage reids消息
type RedisMessage struct {
	Message string
	HasData bool
}

//Ack 确定消息
func (m *RedisMessage) Ack() error {
	return nil
}

//Nack 取消消息
func (m *RedisMessage) Nack() error {
	return nil
}

//GetMessage 获取消息
func (m *RedisMessage) GetMessage() string {
	return m.Message
}

//Has 是否有数据
func (m *RedisMessage) Has() bool {
	return m.HasData
}

//NewStompMessage 创建消息
func NewRedisMessage(cmd *redis.StringCmd) *RedisMessage {
	msg, err := cmd.Result()
	hasData := err == nil
	return &RedisMessage{Message: msg, HasData: hasData}
}
