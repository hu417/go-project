package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"ratelimit-demo/global"

	"github.com/gin-gonic/gin"
)

func LimitV5(ctx *gin.Context) {
	// 使用 Redis INCR 命令实现计数器限流
	key := fmt.Sprintf("ratelimit:%s", strings.Split(ctx.Request.RemoteAddr, ":")[0])
	log.Printf("key: %s", key)
	count, err := global.RdsCli.Incr(ctx, key).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "Internal Server Error",
		})

		return
	}

	// 设置过期时间为 1 分钟
	global.RdsCli.Expire(ctx, key, time.Minute)

	// 如果超过限制，则返回 429 Too Many Requests 错误
	if count > 10 {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"code":  429,
			"error": "超过请求限制",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "limit v5",
	})
}
