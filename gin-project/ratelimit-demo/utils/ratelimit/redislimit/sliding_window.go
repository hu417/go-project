package redislimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
滑动窗口

滑动窗口法，简单来说就像是一扇不断移动的窗户，随着时间的推移，窗户会不断前进。窗户内有一个计数器，它持续记录窗户内发生的请求数量，这样就可以确保在任何时间段内请求数量不会超过最大允许的值。例如，假设当前的时间窗口是从 0 秒开始到 60 秒结束，窗户内的请求数是 40。当时间过去了 10 秒后，窗户的位置就向前移动了，变成了从 10 秒开始到 70 秒结束，而窗户内的请求数变成了 60

*/

type RateLimiter struct {
	client   *redis.Client
	context  context.Context
	rate     int           // 固定的窗口大小
	interval time.Duration // 时间窗口大小
	key      string        // Redis键名
}

func NewRateLimiter(rate int, interval time.Duration, key string) *RateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis服务器地址和端口号
	})

	return &RateLimiter{
		client:   client,
		rate:     rate,
		interval: interval,
		key:      key,
	}
}

func (limiter *RateLimiter) Allow() (bool, error) {
	now := time.Now()

	pipe := limiter.client.Pipeline()

	pipe.ZRemRangeByScore(limiter.context, limiter.key, "-inf", strconv.FormatInt(now.Add(-limiter.interval).Unix(), 10))
	pipe.ZCard(limiter.context, limiter.key)
	_, err := pipe.ZAdd(limiter.context, limiter.key, redis.Z{
		Score:  float64(now.Unix()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	}).Result()
	if err != nil {
		return false, err
	}

	cmds, err := pipe.Exec(limiter.context)
	if err != nil {
		return false, err
	}

	count, _ := cmds[1].(*redis.IntCmd).Result()
	return int(count) <= limiter.rate, nil
}
