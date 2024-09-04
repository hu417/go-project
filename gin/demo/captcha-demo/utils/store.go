package utils

/**
使用redis需实现Store中的三个方法
type Store interface {
    // Set sets the digits for the captcha id.
    Set(id string, value string)
    // Get returns stored digits for the captcha id. Clear indicates
    // whether the captcha must be deleted from the store.
    Get(id string, clear bool) string
    //Verify captcha's answer directly
    Verify(id, answer string, clear bool) bool
}
*/

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// 定义常量：key 前缀
const CAPTCHA = "captcha-demo:login:captcha:"

type RedisStore struct {
	RDS *redis.Client
}

// 实现初始化 captcha 的方法
func NewCaptchaStore(r *redis.Client) *RedisStore {
	return &RedisStore{
		RDS: r,
	}
}

// 实现设置 captcha 的方法
func (r *RedisStore) Set(id string, value string) error {
	key := CAPTCHA + id
	return r.RDS.Set(ctx, key, value, time.Minute*2).Err()
}

// 实现获取 captcha 的方法
func (r *RedisStore) Get(id string, clear bool) string {
	key := CAPTCHA + id
	//获取 captcha
	val, err := r.RDS.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	//如果clear == true, 则删除
	if clear {
		err := r.RDS.Del(ctx, key).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

// 实现验证 captcha 的方法
func (r *RedisStore) Verify(id, answer string, clear bool) bool {
	v := r.Get(id, clear)
	return v == answer
}
