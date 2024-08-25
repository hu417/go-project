package redislimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
固定窗口法:
固定窗口法是限流算法里面最简单的，比如我想限制1分钟以内请求为100个，从现在算起的一分钟内，请求就最多就是100个，这分钟过完的那一刻把计数器归零，重新计算，周而复始。

当一个时间窗口结束时，下一个时间窗口立即开始，这就意味着窗口切换是瞬间完成的。在窗口切换的瞬间，如果有大量请求同时到达，就会出现流量的剧烈波动。例如，在一个窗口结束时，有很多请求被允许通过，而下一个窗口开始时，大量请求被阻塞。这种波动可能会对系统造成压力，导致延迟增加，甚至出现系统故障。


*/

type FixedWindowRateLimiter struct {
	RedisClient *redis.Client
	WindowSize  int
	MaxRequests int
	WindowKey   string
	CountKey    string
	ctx         context.Context
}

func NewFixedWindowRateLimiter(redisClient *redis.Client, windowSize, maxRequests int, windowKey, countKey string) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		RedisClient: redisClient,
		WindowSize:  windowSize,
		MaxRequests: maxRequests,
		WindowKey:   windowKey,
		CountKey:    countKey,
		ctx:         context.Background(),
	}
}

func (limiter *FixedWindowRateLimiter) Allow() bool {
	// 获取当前时间戳
	currentTime := time.Now().Unix()

	// 判断当前时间窗口是否重置
	resetTimeStr, err := limiter.RedisClient.Get(limiter.ctx, limiter.WindowKey).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if resetTimeStr == "" {
		// 如果重置时间不存在，说明窗口尚未初始化，设置当前时间为重置时间，并设置计数器初始值为1
		resetTime := fmt.Sprintf("%v", currentTime+int64(limiter.WindowSize))
		limiter.RedisClient.Set(limiter.ctx, limiter.WindowKey, resetTime, 0)
		limiter.RedisClient.Set(limiter.ctx, limiter.CountKey, 1, 0)
	} else {
		resetTime, _ := strconv.ParseInt(resetTimeStr, 10, 64)
		if currentTime > resetTime {
			// 如果当前时间超过了重置时间，说明窗口需要重置，设置新的重置时间并将计数器重置为1
			resetTime := fmt.Sprintf("%v", currentTime+int64(limiter.WindowSize))
			limiter.RedisClient.Set(limiter.ctx, limiter.WindowKey, resetTime, 0)
			limiter.RedisClient.Set(limiter.ctx, limiter.CountKey, 1, 0)
		} else {
			// 如果当前时间未超过重置时间，说明窗口仍在有效期内，增加计数器的值
			limiter.RedisClient.Incr(limiter.ctx, limiter.CountKey)
		}
	}

	// 获取当前计数器的值
	totalRequests, err := limiter.RedisClient.Get(limiter.ctx, limiter.CountKey).Int64()
	if err != nil {
		panic(err)
	}

	// 检查请求数量是否超过限制
	if totalRequests > int64(limiter.MaxRequests) {
		return false
	}
	return true
}
