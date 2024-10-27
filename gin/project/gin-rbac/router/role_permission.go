package router

import (
	"gin-rbac/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRolePermissionRoutes(router *gin.RouterGroup, rolePermissionHandler *handler.RolePermissionHandler) {
	rolePermissionGroup := router.Group("/role-permissions")
	{
		rolePermissionGroup.POST("/", rolePermissionHandler.CreateRolePermissionsByIDs)                                     // 创建角色权限列表
		rolePermissionGroup.DELETE("/", rolePermissionHandler.DeleteRolePermissionsByIDs)                                   // 删除角色权限列表
		rolePermissionGroup.GET("/permissions/:permissionID/roles", rolePermissionHandler.GetRolePermissionsByPermissionID) // 获取角色权限列表
		rolePermissionGroup.GET("/roles/:roleID/permissions", rolePermissionHandler.GetRolePermissionsByRoleID)             // 获取角色权限列表
		rolePermissionGroup.GET("/roles/:roleID/permissions/:permissionID", rolePermissionHandler.GetRolePermissionByID)    // 获取角色权限
	}
}
