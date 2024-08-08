package initialize

import (
	"context"
	
	"gin-base/config"
	"gin-base/global"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cli config.Redis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cli.Addr,
		Password: cli.Password, // no password set
		DB:       cli.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	global.Log.Sugar().Debugf("redis connect ping response: %v", pong)
	return client

}
