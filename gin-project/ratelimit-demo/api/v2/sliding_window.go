package v2

import (
	"net/http"
	"time"

	"ratelimit-demo/utils/ratelimit/redislimit"

	"github.com/gin-gonic/gin"
)

/*
滑动窗口

*/

var (
	limiterV2 = redislimit.NewRateLimiter(5, time.Second*60, "my-rate-limiter")
)

func LimitV2(ctx *gin.Context) {

	allowed, err := limiterV2.Allow()
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	if allowed {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "请求通过",
		})
		return
	}
	ctx.JSON(http.StatusTooManyRequests, gin.H{
		"code":    429,
		"message": "请求超出限制",
	})

	time.Sleep(200 * time.Millisecond)

}
