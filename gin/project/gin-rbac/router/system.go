package router

import (
	"gin-rbac/handler"

	"github.com/gin-gonic/gin"
)

// RegisterSystemRoutes 注册系统路由
func RegisterSystemRoutes(router *gin.RouterGroup, systemHandler *handler.SystemHandler, engine *gin.Engine) {
	// 系统路由
	systemGroup := router.Group("/system")
	{
		systemGroup.GET("/health", systemHandler.Health) // 服务器健康检查
		systemGroup.GET("/routes", func(c *gin.Context) {
			systemHandler.GetRoutes(engine, c) // 使用闭包来传递engine
		}) // 获取路由列表
	}
}
