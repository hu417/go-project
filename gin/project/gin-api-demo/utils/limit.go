package utils

import (
	"context"
	"errors"
	"time"

	"gin-api-demo/global"
)

type LimitConfig struct {
	Key    string // Key 限制key
	Expire int    // Expire 过期时间
	Limit  int    // Limit 限制次数
}

func NewLimitConfig(key string, times, limit int) *LimitConfig {
	return &LimitConfig{
		Key:    key,
		Expire: times, // 过期时间
		Limit:  limit, // 限制次数
	}
}

// SetLimitWithTime 设置访问次数
func (l *LimitConfig) SetLimitWithTime() error {
	// 判断是否开启redis
	if global.Redis == nil {
		return errors.New("redis未开启")
	}
	// 判断key是否存在
	count, err := global.Redis.Exists(context.Background(), l.Key).Result()
	if err != nil {
		return err
	}
	// 不存在则设置,设置过期时间,设置次数,设置过期时间,设置次数
	if count == 0 {
		pipe := global.Redis.TxPipeline()
		pipe.Incr(context.Background(), l.Key)
		pipe.Expire(context.Background(), l.Key, time.Duration(l.Expire)*time.Second)
		_, err = pipe.Exec(context.Background())
		return err
	} else {
		// 次数
		if count, err := global.Redis.Get(context.Background(), l.Key).Int(); err != nil {
			return err
		} else {
			if count >= l.Limit {
				if t, err := global.Redis.PTTL(context.Background(), l.Key).Result(); err != nil {
					return errors.New("请求太过频繁，请稍后再试")
				} else {
					return errors.New("请求太过频繁, 请 " + t.String() + " 秒后尝试")
				}
			} else {
				return global.Redis.Incr(context.Background(), l.Key).Err()
			}
		}
	}
}
