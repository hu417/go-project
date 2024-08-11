package route

import (
	"api-demo/app/api/admin"

	"github.com/gin-gonic/gin"
)

func genAdminRouter(rg *gin.RouterGroup) {
	rg.GET("/a", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "a"})
	})

	// 管理员
	rg.GET("/profile", admin.AdminApi.Profile)
}
