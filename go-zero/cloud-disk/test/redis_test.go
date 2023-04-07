package test

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "10.0.0.91:6379",
		Password: "123",
		DB:       0,
	})
)

func TestRdb(t *testing.T) {
	result, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Log(err)
	}
	t.Log(result)

	err = rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		t.Log(err)
	}

	if val, err := rdb.Get(ctx, "key").Result(); err != nil {
		t.Log(err)
	} else {
		t.Log(val)
	}

}
