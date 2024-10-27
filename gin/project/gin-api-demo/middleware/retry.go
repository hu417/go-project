package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RetryMiddleware 是一个自定义的中间件，用于处理请求失败后的重试逻辑
// maxRequestSize 用于限制请求体的最大大小
// 设置请求体最大大小为 10MB: maxRequestSize := 10 * 1024 * 1024 // 10 MB
func Retry(maxRequestSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		maxRetries := 3                  // 最大重试次数
		retryInterval := 2 * time.Second // 重试间隔

		var retries int
		for retries < maxRetries {
			c.Next() // 继续执行后续的处理

			if c.Writer.Status() >= 500 { // 假设服务器错误才重试
				time.Sleep(retryInterval)
				c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestSize)
				c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "ginContext", c))
				retries++
			} else {
				break // 成功或非服务器错误则不再重试
			}
		}

		if retries == maxRetries {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Max retries exceeded",
			})
		}
	}
}
