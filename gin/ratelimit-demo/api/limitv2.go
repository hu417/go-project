package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 令牌桶: 有一个固定大小的桶，系统会以恒定速率向桶中放 Token，桶满则暂时不放

var (
	// 创建一个限流器，每秒允许 10 个请求，突发 5 个请求
	limiter = rate.NewLimiter(rate.Limit(10), 5)
	c, _  = context.WithTimeout(context.TODO(), time.Millisecond)
)

func LimitV2(ctx *gin.Context) {
	// 尝试获取令牌，如果获取失败，则返回 429 Too Many Requests 错误
	if err := limiter.Wait(c); err != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"code": 429,
			"msg":  "超过请求限制",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "limit v2",
	})
}
