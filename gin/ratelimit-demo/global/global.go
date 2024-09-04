package global

import (
	"ratelimit-demo/utils/ratelimit"

	"github.com/redis/go-redis/v9"
)

var (
	RdsCli *redis.Client
)

func Limit_v4() *ratelimit.SlidingWindowLimiter {
	limit, err := ratelimit.NewSlidingWindowLimiter(10, 10, 1)
	if err != nil {
		panic(err)
	}
	return limit
}
