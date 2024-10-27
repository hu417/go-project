package router

import (
	"gin-rbac/handler"

	"github.com/gin-gonic/gin"
)

// RegisterPermissionRoutes 注册权限路由
func RegisterPermissionRoutes(router *gin.RouterGroup, permissionHandler *handler.PermissionHandler) {
	permissionGroup := router.Group("/permissions")
	{
		permissionGroup.GET("/", permissionHandler.GetPermissionList)            // 获取权限列表
		permissionGroup.POST("/", permissionHandler.CreatePermission)            // 创建权限
		permissionGroup.GET("/:id", permissionHandler.GetPermission)             // 获取权限
		permissionGroup.PUT("/:id", permissionHandler.UpdatePermission)          // 更新权限
		permissionGroup.DELETE("/:id", permissionHandler.DeletePermission)       // 删除权限
		permissionGroup.PUT("/:id/recover", permissionHandler.RecoverPermission) // 恢复权限
	}
}
