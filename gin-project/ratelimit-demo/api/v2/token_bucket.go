package v2

import (
	"net/http"
	"ratelimit-demo/global"
	"ratelimit-demo/utils/ratelimit/redislimit"

	"github.com/gin-gonic/gin"
)

/*
令牌桶限流


*/

var (
	leakyBucket = redislimit.NewTokenLeakyBucket(global.RdsCli, "myleakybucket", 100, 10)
)

func LimitV4(ctx *gin.Context) {

	// 获取令牌
	defer leakyBucket.Close()

	leakyBucket.Start()

	if leakyBucket.ProcessRequest() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"code":    429,
			"message": "请求超出限制",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求通过",
	})

}
