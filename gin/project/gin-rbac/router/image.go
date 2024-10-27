package router

import (
	"gin-rbac/handler"

	"github.com/gin-gonic/gin"
)

// RegisterImageRoutes 注册图片路由
func RegisterImageRoutes(router *gin.RouterGroup, imageHandler *handler.ImageHandler) {
	// 图片路由
	imageGroup := router.Group("/images")
	{
		imageGroup.POST("/upload", imageHandler.UploadImage) // 上传图片
		imageGroup.GET("/:id", imageHandler.GetImage)        // 获取图片
		imageGroup.DELETE("/:id", imageHandler.DeleteImage)  // 删除图片
	}

}
