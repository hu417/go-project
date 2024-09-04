package v2

import (
	"net/http"

	"ratelimit-demo/global"
	"ratelimit-demo/utils/ratelimit/redislimit"

	"github.com/gin-gonic/gin"
)

/*
固定窗口限流

当一个时间窗口结束时，下一个时间窗口立即开始，这就意味着窗口切换是瞬间完成的。在窗口切换的瞬间，如果有大量请求同时到达，就会出现流量的剧烈波动。例如，在一个窗口结束时，有很多请求被允许通过，而下一个窗口开始时，大量请求被阻塞。这种波动可能会对系统造成压力，导致延迟增加，甚至出现系统故障。

*/

// 创建固定窗口限流器实例
var (
	limiterV1 = redislimit.NewFixedWindowRateLimiter(
		global.RdsCli, 10, 100, "fixedwindow", "fixedwindow:count")
)

func LimitV1(ctx *gin.Context) {
	// 测试请求是否通过限流
	if limiterV1.Allow() {
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

}
