package redislimit

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
令牌桶限流


*/

type TokenBucket struct {
	rdb            *redis.Client
	leakyBucketKey string
	capacity       int64 // 漏桶容量
	rate           int64 // 每秒的令牌添加速率
}

func NewTokenLeakyBucket(rdb *redis.Client, key string, capacity, rate int64) *TokenBucket {
	return &TokenBucket{
		rdb:            rdb,
		leakyBucketKey: key,
		capacity:       capacity,
		rate:           rate,
	}
}

func (lb *TokenBucket) Start() {
	go func() {
		for {
			time.Sleep(time.Second / time.Duration(lb.rate))
			err := lb.rdb.ZAddNX(context.Background(), lb.leakyBucketKey, redis.Z{
				Score:  float64(time.Now().UnixNano()),
				Member: "token",
			}).Err()
			if err != nil {
				log.Println("Failed to add token:", err)
			}
		}
	}()
}

func (lb *TokenBucket) ProcessRequest() bool {
	result, err := lb.rdb.ZRangeByScore(context.Background(), lb.leakyBucketKey, &redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(time.Now().UnixNano(), 10),
		Offset: 0,
		Count:  1,
	}).Result()

	if err != nil {
		log.Println("Failed to get token:", err)
		return false
	}

	if len(result) > 0 {
		_, err := lb.rdb.ZRem(context.Background(), lb.leakyBucketKey, result[0]).Result()
		if err != nil {
			log.Println("Failed to remove token:", err)
		}
		return true
	}

	return false
}

func (lb *TokenBucket) Close() {
	if err := lb.rdb.Close(); err != nil {
		log.Println("Failed to close Redis connection:", err)
	}
}
