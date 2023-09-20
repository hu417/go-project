package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 结构体
type RedisCache struct {
	client *redis.Client
}

// 接口方法
type cache interface{
	Set(key string,value interface{}) interface{}
	Get(key string) interface{}
}

// 接口实例化
func NewCache(tp string) (cache,error){
	switch tp {
	case "redis":
		return NewRedisCache(),nil
	case "cache":
		return NewRedis(),nil
	default:
		return NewRedisCache(),nil
	}
}


// 初始化结构体redis-client
func NewRedisCache() *RedisCache {
	c := &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
	return c
}

func NewRedis() *RedisCache {
	c := &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
	return c
}


// Get方法
func (c *RedisCache) Get(key string) interface{} {
	ctx := context.Background()
	return c.client.Get(ctx,key)
	
}

// Set方法
func (c *RedisCache) Set(key string,value interface{}) interface{} {
	ctx := context.Background()
	return c.client.Set(ctx,key,value,-1)
}

func main() {
  // 结构体方法
	rc := NewRedisCache()
	rc.Set("name","hu")
	rc.Get("name")

  // 接口方法
  rs,err := NewCache("redis")
	if err != nil {
		fmt.Println("初始化失败")
	}
	rs.Set("name","hu")
	rs.Get("name")
}
