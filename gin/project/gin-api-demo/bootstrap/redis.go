package bootstrap

import (
	"context"
	"fmt"
	"gin-api-demo/config"

	"github.com/redis/go-redis/v9"
)

func InitRedis(conf *config.Config) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password, // 没有密码, 默认值
		DB:           conf.Redis.DB,       // 默认DB 0
		ReadTimeout:  -1,                  // 从网络连接中读取数据超时时间
		WriteTimeout: -1,                  // 把数据写入网络连接的超时时间
		PoolSize:     1000,                // 连接池最大连接数量
		MinIdleConns: 10,                  // 连接池保持的最小空闲连接数
		MaxIdleConns: 100,                 // 连接池保持的最大空闲连接数
	})

	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Errorf("连接redis出错, 错误信息: %v", err))
	}
	return cli
}
