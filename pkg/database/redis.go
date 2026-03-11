package database

import (
	"context"
	"fmt"
	"go-todo-api/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(config config.RedisConfig) (*redis.Client, error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		_ = RedisClient.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", err) // 包装错误
	}
	fmt.Println("redis链接成功")
	return RedisClient, nil
}
