package middleware

import (
	"ratelimit-demo/global"
	"ratelimit-demo/utils/ratelimit/redislimit"

	"github.com/gin-gonic/gin"
)

func LimitV4(c *gin.Context) {
	leakyBucket := redislimit.NewTokenLeakyBucket(global.RdsCli, "myleakybucket", 100, 10)
	defer leakyBucket.Close()

	leakyBucket.Start()


}
