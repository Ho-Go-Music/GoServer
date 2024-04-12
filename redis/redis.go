package redis

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	once        sync.Once
	redisClient *redis.Client
)

func InitRedisClient() {
	once.Do(func() {
		// 创建 Redis 客户端
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "775028", // 如果密码已设置，请提供密码
			DB:       0,        // 选择要使用的数据库
		})
	})
}

func GetRedisClient() *redis.Client {
	InitRedisClient()
	return redisClient
}
