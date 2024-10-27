package router

import (
	"gin-rbac/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoleRoutes 注册角色路由
func RegisterRoleRoutes(router *gin.RouterGroup, roleHandler *handler.RoleHandler) {
	// 角色路由
	roleGroup := router.Group("/roles")
	{
		roleGroup.GET("/", roleHandler.GetRoleList)            // 获取角色列表
		roleGroup.POST("/", roleHandler.CreateRole)            // 创建角色
		roleGroup.GET("/:id", roleHandler.GetRole)             // 获取角色
		roleGroup.PUT("/:id", roleHandler.UpdateRole)          // 更新角色
		roleGroup.DELETE("/:id", roleHandler.DeleteRole)       // 删除角色
		roleGroup.PUT("/:id/recover", roleHandler.RecoverRole) // 恢复角色
	}
}
