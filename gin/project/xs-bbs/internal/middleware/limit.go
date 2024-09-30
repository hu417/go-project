package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func DefaultLimit(fillInterval time.Duration, cap int64) gin.HandlerFunc {
	// 创建一个限流器，每秒允许1个请求
	limiter := rate.NewLimiter(1, 1)
	return func(c *gin.Context) {
		// 检查是否允许请求
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}
