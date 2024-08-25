package v2

import (
	"context"
	"fmt"
	"net/http"

	"ratelimit-demo/global"
	"ratelimit-demo/utils/ratelimit/redislimit"

	"github.com/gin-gonic/gin"
)

/*
漏桶限流

漏桶限制的是常量流出速率（即流出速率是一个固定常量值），所以最大的速率就是出水的速率，不能出现突发流量。
*/

var (
	// 创建漏桶限流器
	limiterV3 = redislimit.NewLeakyBucketLimiter(global.RdsCli, "my_bucket", 100, 10) // 最高水位为 100，水流速度为 10

)

func LimitV3(ctx *gin.Context) {
	// 尝试获取令牌
	err := limiterV3.TryAcquire(context.Background())
	if err != nil {
		if err == redislimit.ErrAcquireFailed {
			// 请求失败，超出限流水位
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求超出限制",
			})
			return
		}
		// 其他错误处理
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	// 请求成功
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求通过",
	})

}
