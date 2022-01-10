package redisutil

import (
	"github.com/go-redis/redis"
	"ranking/config"
)

// GetRedisClient 获取redis连接
func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Address,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.Database,
	})
	return client
}

func Close(client *redis.Client) {
	_ = client.Close()
}
