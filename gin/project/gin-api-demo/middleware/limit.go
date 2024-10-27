package middleware

import (
	"fmt"
	"net/http"
	"time"

	"gin-api-demo/pkg/utils"

	"github.com/gin-gonic/gin"
)

func DefaultLimit() gin.HandlerFunc {
	// 默认每秒10个请求
	return func(ctx *gin.Context) {
		key := fmt.Sprintf("client::%v", ctx.ClientIP())
		if err := utils.NewLimitConfig(key, 1, 5).SetLimitWithTime(); err != nil {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"time": time.Now().Format("2006-01-02 15:04"),
				"code": 429,
				"msg":  err.Error(),
			})
			ctx.Abort()
			return
		} else {
			ctx.Next()
		}
	}
}
