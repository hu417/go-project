package router

import (
	"gin-rbac/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 在给定的路由组上注册与用户相关的路由。
func RegisterUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	// 用户登录注册路由
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", userHandler.Login)                      // 登录用户--不需要认证
		authGroup.POST("/register", userHandler.CreateUser)              // 注册用户--不需要认证
		authGroup.GET("/user", userHandler.GetUserByJWT)                 // 获取用户--jwt
		authGroup.PUT("/user", userHandler.UpdateUserByJWT)              // 更新用户--jwt
		authGroup.DELETE("/user", userHandler.DeleteUserByJWT)           // 注销用户--jwt
		authGroup.PUT("/user/avatar", userHandler.UpdateUserAvatarByJWT) // 更新用户头像--jwt
		authGroup.PUT("/user/password", userHandler.UpdatePasswordByJWT) // 更新用户密码--jwt
	}
	// 用户路由
	userGroup := router.Group("/users")
	{
		userGroup.GET("/public", userHandler.GetPublicUserList)     // 获取用户公开信息列表--不需要认证
		userGroup.GET("/:id/public", userHandler.GetPublicUser)     // 获取用户公开信息--不需要认证
		userGroup.GET("/:id", userHandler.GetUserByID)              // 获取用户--管理员
		userGroup.PUT("/reset_password", userHandler.ResetPassword) // 重置密码--管理员，后续添加用户忘记密码接口
		userGroup.DELETE("/:id", userHandler.DeleteUser)            // 注销用户--管理员
		userGroup.PUT("/:id/recover", userHandler.RecoverUser)      // 恢复用户--管理员
	}
}
