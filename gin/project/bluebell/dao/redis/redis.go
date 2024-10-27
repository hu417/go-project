package redis

import (
	"bluebell/setting"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Init 初始化连接
func Init(cfg *setting.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,     // no password set
		DB:           cfg.DB,           // use default DB
		ReadTimeout:  -1,               // 从网络连接中读取数据超时时间
		WriteTimeout: -1,               // 把数据写入网络连接的超时时间
		PoolSize:     cfg.PoolSize,     // 连接池最大连接数量
		MinIdleConns: cfg.MinIdleConns, // 连接池保持的最小空闲连接数
		MaxIdleConns: 100,              // 连接池保持的最大空闲连接数
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Errorf("连接redis出错, 错误信息: %v", err))
	}

	return client
}
