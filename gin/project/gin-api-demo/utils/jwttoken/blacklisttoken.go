package jwttoken

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gin-api-demo/global"
	"gin-api-demo/utils"

	"github.com/redis/go-redis/v9"
)

// 获取黑名单缓存 key
func (jwtStr *jwtStr) getBlackListKey(tokenStr string) string {
	return "jwt_black_list:" + utils.MD5([]byte(tokenStr))
}

// JoinBlackList token 加入黑名单
func (jwtStr *jwtStr) JoinBlackList(claims *MyCustomClaims) (err error) {
	nowUnix := time.Now().Unix()
	// 计算 token 剩余时间
	timer := time.Duration(claims.ExpiresAt.Unix()-nowUnix) * time.Second

	// 将 token 剩余时间设置为缓存有效期，并将当前时间作为缓存 value 值
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// fmt.Println(string(claimsJSON))
	err = global.Redis.SetNX(context.Background(), jwtStr.getBlackListKey(string(claimsJSON)), nowUnix, timer).Err()
	return
}

// IsInBlacklist token 是否在黑名单中
func (jwtStr *jwtStr) IsInBlacklist(tokenStr string) (bool, error) {
	// 获取加入黑名单的时间
	joinUnixStr, err := global.Redis.Get(context.Background(), jwtStr.getBlackListKey(tokenStr)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	// 将字符串转成 int64
	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if joinUnixStr == "" || err != nil {
		return false, err
	}
	// JwtBlacklistGracePeriod 为黑名单宽限时间，避免并发请求失效
	if time.Now().Unix()-joinUnix < global.Conf.Jwt.JwtBlacklistGracePeriod {
		return false, errors.New("token 不在黑名单中")
	}
	return true, nil
}
