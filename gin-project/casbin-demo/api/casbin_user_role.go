package api

import (
	"casbin-demo/global"
	"casbin-demo/model/req"
	"casbin-demo/service"

	"github.com/gin-gonic/gin"
)

/* 用户以及组关系的增删改查 */

// @Summary 获取所有用户以及关联的角色
// @Router /api/v1/user/list [get]
func GetUserList(ctx *gin.Context) {
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	ctx.JSON(200, svc.GetUsers())
}

// @Summary 获取所有角色组
func GetRoleList(ctx *gin.Context) {
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	list, err := svc.GetAllRoles()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取角色列表失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"data": list,
		"code": 200,
		"msg":  "获取角色列表成功",
	})
}

// 获取角色下的用户
func GetRoleUser(ctx *gin.Context) {
	role := ctx.Query("role")
	if role == "" {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	data, err := svc.GetUsersForRole(role)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取用户失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取用户成功",
		"data": data,
	})

}

// 角色组中添加用户, 没有组默认创建
func AddUserAtRole(ctx *gin.Context) {
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)

	var userRole *req.UserRole
	err := ctx.ShouldBind(&userRole)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	err = svc.AddUserRole(userRole.UserName, userRole.RoleName)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "添加用户失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "添加用户成功",
	})
}

// 查询用户角色
func GetUserRole(ctx *gin.Context) {
	userName := ctx.Query("user_name")
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	data, err := svc.GetUserRoles(userName)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "查询用户角色失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "查询用户角色成功",
		"data": data,
	})
}

// 判断用户是否具有某个角色
func HasRoleForUser(ctx *gin.Context) {
	var userRole req.UserRole
	if err := ctx.ShouldBind(&userRole); err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "参数错误",
		})
		return
	}

	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	data, err := svc.HasRoleForUser(userRole.UserName, userRole.RoleName)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "系统错误",
			"data": nil,
		})
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "查询用户角色成功",
		"data": data,
	})
}

// 更改用户所属角色
func UpdateUserRole(ctx *gin.Context) {
	var updateUserRoleStr req.UpdateUserRoleStr
	if err := ctx.ShouldBind(&updateUserRoleStr); err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "参数错误",
		})
		return
	}
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	err := svc.UpdateUserRole(updateUserRoleStr.UserName, updateUserRoleStr.OldRole, updateUserRoleStr.NewRole)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "更改用户角色失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "更改用户角色成功",
	})
}

// 角色组中删除用户
func DeleteUserAtRole(ctx *gin.Context) {
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	var userRole *req.UserRole
	err := ctx.ShouldBind(&userRole)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	err = svc.DeleteUserRole(userRole.UserName, userRole.RoleName)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "删除用户失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除用户成功",
	})
}
