package redis

import (
	"github.com/Ho-Go-Music/GoServer/tools"
	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:           tools.Conf.RedisServer.Host + ":" + tools.Conf.RedisServer.Port,
		Password:       tools.Conf.RedisServer.Password,
		DB:             tools.Conf.RedisServer.Database, // Selecting the database to be used
		MaxActiveConns: tools.Conf.RedisServer.MaxActiveConns,
		MaxIdleConns:   tools.Conf.RedisServer.MaxIdleConns,
	})
	return redisClient
}
