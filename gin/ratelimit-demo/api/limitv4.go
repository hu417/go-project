package api

import (
	"net/http"
	"ratelimit-demo/global"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	limit_v4 = global.Limit_v4()
)

func LimitV4(ctx *gin.Context) {
	if !limit_v4.TryAcquire() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"code": http.StatusTooManyRequests,
			"msg":  "请求超过限制",
			"time": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"msg": "limit v4",
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})
}
