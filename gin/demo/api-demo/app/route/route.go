package route

import (
	"github.com/gin-gonic/gin"
)

type AppRouter struct{}

func New() *AppRouter {
	return &AppRouter{}
}

func (*AppRouter) AddRoute(e *gin.Engine) {
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ping"})
	})

	e.GET("/pong", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 业务路由分组
	genAdminRouter(e.Group("/admin"))
	genAppRouter(e.Group("/app"))

	// 获取所有路由
	e.GET("/routes", func(c *gin.Context) {
		// 获取所有路由信息
		routesInfo := e.Routes()
		c.JSON(200, gin.H{"routes": routesInfo})
	})
}
