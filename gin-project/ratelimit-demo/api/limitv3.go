package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 固定窗口: 每开启一个新的窗口，在窗口时间大小内，可以通过窗口请求上限个请求;
// 该算法主要是会存在临界问题，如果流量都集中在两个窗口的交界处，那么突发流量会是设置上限的两倍

// 限流,允许每秒最多 10 个请求。
var (
	// limit_v3 = ratelimit.NewFixedWindowLimiter(10, 1)
	limit_v3 = rate.NewLimiter(rate.Limit(10), 1) // 允许每秒10次

)

func LimitV3(ctx *gin.Context) {

	// limit_v3.TryAcquire()

	if !limit_v3.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"code": 429,
			"msg":  "请求超过限制",
			"time": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"msg":  "limit v3",
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})

}
