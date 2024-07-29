package router

import (
	"jwt-demo/api"
	"jwt-demo/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 跨域
	router.Use(middleware.Cors())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/register", api.RegisterHandler)
	router.POST("/login", api.LoginHandler)

	auth := router.Group("/api/")
	{
		auth.Use(middleware.JwtMiddleWare())
		auth.GET("/home", api.HomeHandler)
		auth.GET("/refreshtoken", api.RefreshToken)
	}

	// 未知调用方式
	router.NoMethod(NoMethodJson)
	// 未知路由处理
	router.NoRoute(NoRouteJson)

	return router

}

// 未知路由处理 返回json
func NoRouteJson(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
		"msg":  "path not found",
	})
}

// 未知调用方式 返回json
func NoMethodJson(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"code": http.StatusMethodNotAllowed,
		"msg":  "method not allowed",
	})
}
