package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// 中间件函数
func MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[INFO] %s", "中间件业务执行")
		c.Next()
	}
}
