package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	host := GetEnv("REDIS_HOST", "localhost")
	port := GetEnv("REDIS_PORT", "6379")
	pass := GetEnv("REDIS_PASSWORD", "")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
		DB:       0,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Warning: Failed to connect to Redis:", err)
	} else {
		log.Println("Redis connection established")
	}
}
