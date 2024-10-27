package bootstrap

import (
	"context"

	"gin-rbac/config"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) (*redis.Client, error) {
	Addr := cfg.Redis.GetRedisAddress()
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
