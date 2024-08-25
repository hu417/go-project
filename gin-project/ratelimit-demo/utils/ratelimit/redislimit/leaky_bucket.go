package redislimit

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
漏桶限流

一个固定容量的桶，有水流进来，也有水流出去。对于流进来的水来说，我们无法预计一共有多少水会流进来，也无法预计水流的速度。但是对于流出去的水来说，这个桶可以固定水流出的速率（处理速度），从而达到 流量整形 和 流量控制 的效果
*/

var (
	ErrAcquireFailed = errors.New("acquire failed")
)

type LeakyBucketLimiter struct {
	peakLevel       int
	currentVelocity int
	key             string
	client          *redis.Client
}

func NewLeakyBucketLimiter(client *redis.Client, key string, peakLevel, currentVelocity int) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		peakLevel:       peakLevel,
		currentVelocity: currentVelocity,
		key:             key,
		client:          client,
	}
}

func (l *LeakyBucketLimiter) acquireToken(ctx context.Context) (bool, error) {
	pipe := l.client.TxPipeline()
	pipe.HGet(ctx, l.key, "lastTime")
	pipe.HGet(ctx, l.key, "currentLevel")
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	result, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	lastTimeStr, _ := result[0].(*redis.StringCmd).Result()
	currentLevelStr, _ := result[1].(*redis.StringCmd).Result()

	lastTime, _ := strconv.ParseInt(lastTimeStr, 10, 64)
	currentLevel, _ := strconv.ParseInt(currentLevelStr, 10, 64)

	now := time.Now().Unix()
	interval := now - lastTime
	newLevel := currentLevel - (interval * int64(l.currentVelocity))

	if newLevel < 0 {
		newLevel = 0
	}

	if newLevel >= int64(l.peakLevel) {
		return false, nil
	}

	newLevel++
	pipe.HSet(ctx, l.key, "currentLevel", newLevel)
	pipe.HSet(ctx, l.key, "lastTime", now)
	pipe.Expire(ctx, l.key, time.Duration(l.peakLevel/l.currentVelocity)*time.Second)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (l *LeakyBucketLimiter) TryAcquire(ctx context.Context) error {
	success, err := l.acquireToken(ctx)
	if err != nil {
		return err
	}

	if !success {
		return ErrAcquireFailed
	}

	return nil
}
