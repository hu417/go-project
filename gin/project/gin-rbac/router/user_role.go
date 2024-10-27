package router

import (
	"gin-rbac/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoleRoutes(router *gin.RouterGroup, userRoleHandler *handler.UserRoleHandler) {
	userRoleGroup := router.Group("/user-roles")
	{
		userRoleGroup.POST("/", userRoleHandler.CreateUserRolesByIDs)                      // 创建用户角色列表
		userRoleGroup.DELETE("/", userRoleHandler.DeleteUserRolesByIDs)                    // 删除用户角色列表
		userRoleGroup.GET("/users/:userID/roles", userRoleHandler.GetUserRolesByUserID)    // 获取用户角色列表
		userRoleGroup.GET("/roles/:roleID/users", userRoleHandler.GetUserRolesByRoleID)    // 获取角色用户列表
		userRoleGroup.GET("/users/:userID/roles/:roleID", userRoleHandler.GetUserRoleByID) // 获取用户角色

	}
}
