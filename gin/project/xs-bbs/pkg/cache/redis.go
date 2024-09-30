package cache

import (
	"context"

	"xs-bbs/pkg/conf"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Init 初始化redis连接
func InitRedis(cfg *conf.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password, // no password set
		DB:           cfg.DB,             // use default db
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 10,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		zap.L().Error("redis ping failed", zap.Error(err))
		return nil, err
	}
	return rdb, nil
}
