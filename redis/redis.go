package redis

import (
	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:           "localhost:6379",
		Password:       "775028", // 如果密码已设置，请提供密码
		DB:             0,        // 选择要使用的数据库
		MaxActiveConns: 100,
		MaxIdleConns:   10,
	})

	return redisClient
}
