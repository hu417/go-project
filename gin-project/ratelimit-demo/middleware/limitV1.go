package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 创建一个 map 来存储每个 IP 的访问限制
var visitors = make(map[string]*rate.Limiter)

// 获取用户 IP 为 key 的令牌桶，若不存在则创建一个
func getVisitor(ip string) *rate.Limiter {
	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(1, 5) // 每秒生成1个令牌，且桶中最多存储5个令牌
		visitors[ip] = limiter
	}

	return limiter
}

// 控制访问频率的中间件
func LimitV1(c *gin.Context) {
	limiter := getVisitor(c.ClientIP())
	// Allow方法返回是否可以接收请求，此处是否可以拿取令牌
	if !limiter.Allow() {
		c.String(http.StatusTooManyRequests, "Too many requests.")
		c.Abort()
		return
	}

	c.Next()
}
