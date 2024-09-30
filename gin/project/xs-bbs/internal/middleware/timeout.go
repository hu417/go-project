package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// TimeoutMiddleware 是一个用于设置超时时间的中间件
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		response := make(chan bool, 1)

		go func() {
			c.Next()
			response <- true
		}()

		select {
		case <-response:
			return
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				c.JSON(http.StatusRequestTimeout, gin.H{"message": "Request timeout"})
				c.Abort()
				return
			}
		}
	}
}
