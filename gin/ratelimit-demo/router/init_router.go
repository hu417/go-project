package router

import (
	"ratelimit-demo/api"
	"ratelimit-demo/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 添加路由
	Api := r.Group("/api/")

	{
		v1 := Api.Group("/v1/").Use(middleware.LimitV1)
		v1.GET("/limit/v1", api.LimitV1)

		Api.GET("/limit/v2", api.LimitV2)
		Api.GET("/limit/v3", api.LimitV3)
		Api.GET("/limit/v4", api.LimitV4)
		Api.GET("/limit/v5", api.LimitV5)
	}

	return r
}
