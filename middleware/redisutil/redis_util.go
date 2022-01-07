package redisutil

import (
	"github.com/go-redis/redis"
	"log"
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
	err := client.Close()
	if err != nil {
		log.Panicln("close redis client failed", err)
	}
}
