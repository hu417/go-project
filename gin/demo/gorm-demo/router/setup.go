package router

import (
	"gorm-demo/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	//	r.Use(middleware.Cors())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//
	api := r.Group("/api")
	{
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello",
			})
		})

		// 单次创建
		api.POST("/user", controller.CreateUser)
		// 批量创建
		api.POST("/users", controller.CreateUsers)

		// 删除
		api.DELETE("/user/:id", controller.DeleteUser)
		api.DELETE("/users", controller.DeleteUsers)

		// 更新
		api.PUT("/user/:id", controller.UpdateUser)
		api.PUT("/users", controller.UpdateUserByName)
	}

	return r
}
