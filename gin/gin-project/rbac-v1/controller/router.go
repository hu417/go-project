package controller

import (
	"github.com/gin-gonic/gin"
	"rbac-v1/middleware"
)

func InitApiRouter(r *gin.Engine) {
	r.Use(middleware.Cors())
	{
		r.POST("/api/login", LoginByPwd)
		r.GET("/api/checkToken", CheckToken)
	}
	
	r.Use(middleware.CheckToken())
	rbac := r.Group("/api/rbac")
	{
		rbac.POST("/auth", RbacAuth)
		rbac.GET("/user/roleAndPower", GetUserRoleAndPower)
		rbac.GET("/policy", middleware.RbacAuth(), RbacPolicy)
	}
	
	r.Use(middleware.RbacAuth())
	op := r.Group("/api/op")
	{
		op.GET("/list", GetOperationList)
		op.POST("/create", OperationCreate, middleware.RbacRefresh())
		op.POST("/update", OperationUpdate, middleware.RbacRefresh())
		op.POST("/del", OperationDelete, middleware.RbacRefresh())
	}
	

	po := r.Group("/api/power")
	{
		po.GET("/list", GetPowerList)
		po.POST("/create", PowerCreate, middleware.RbacRefresh())
		po.POST("/update", PowerUpdate, middleware.RbacRefresh())
		po.POST("/del", PowerDelete, middleware.RbacRefresh())
	}
	

	ro := r.Group("/api/role")
	{
		ro.GET("/list", GetRoleList)
		ro.POST("/create", RoleCreate, middleware.RbacRefresh())
		ro.POST("/update", RoleUpdate, middleware.RbacRefresh())
		ro.POST("/del", RoleDelete, middleware.RbacRefresh())
	}
	

	user := r.Group("/api/user")
	{
		user.GET("/list", GetUserList)
		user.POST("/create", UserCreate, middleware.RbacRefresh())
		user.POST("/update", UserUpdate, middleware.RbacRefresh())
		user.POST("/del", UserDelete, middleware.RbacRefresh())
	}
	
}
