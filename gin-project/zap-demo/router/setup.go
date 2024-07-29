package router

import (
	"zap-demo/api"
	"zap-demo/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	r.GET("/ping", api.Ping)

	return r

}
