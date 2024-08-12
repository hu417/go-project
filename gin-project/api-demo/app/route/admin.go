package route

import (
	"api-demo/app/api/admin"
	"api-demo/app/middleware"

	"github.com/gin-gonic/gin"
)

func genAdminRouter(rg *gin.RouterGroup) {
	rg.Use(middleware.CorsMiddleware)

	rg.GET("/a", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "a"})
	})

	// 登陆接口，
	rg.POST("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "login"})
	})

	// 注册权限验证中间件
	rg.Use(middleware.PermissionMiddleware)
	// 管理员
	rg.GET("/profile", admin.AdminApi.Profile)
	rg.POST("/save", admin.AdminApi.Save)
}
