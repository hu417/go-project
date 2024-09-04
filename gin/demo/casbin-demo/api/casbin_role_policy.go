package api

import (
	"casbin-demo/global"
	"casbin-demo/model/req"
	"casbin-demo/service"

	"github.com/gin-gonic/gin"
)

/* 角色组权限的增删改查 */

// 获取策略中的所有授权规则
func GetPolicyList(ctx *gin.Context) {
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	list, err := svc.GetPolicy()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取失败",
			"data": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": list,
	})
}

// @Summary 获取所有角色组的策略
func GetRolePolicyList(ctx *gin.Context) {
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	list, err := svc.GetRolePolicy()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取失败",
			"data": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": list,
	})

}

// 获取策略中的所有授权规则，可以指定角色字段筛选器
func GetFilteredPolicy(ctx *gin.Context) {
	role := ctx.Query("role")
	if role == "" {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取失败",
			"data": "参数错误",
		})
		return
	}
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	list, err := svc.GetFilteredPolicy(0, role)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取失败",
			"data": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": list,
	})
}

// 获取命名策略中的所有授权规则
func GetNamedPolicy(ctx *gin.Context) {
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	list, err := svc.GetNamedPolicy("p") // ptype = "p"
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取失败",
			"data": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": list,
	})
}

// 获取命名策略中的所有授权规则，可以指定字段过滤器。
func GetFilteredNamedPolicy(ctx *gin.Context) {
	role := ctx.Query("role")
	if role == "" {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "role不能为空",
			"data": nil,
		})
		return
	}
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	list, err := svc.GetFilteredNamedPolicy("p", 0, role)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取失败",
			"data": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": list,
	})
}

// @Summary 创建角色组权限, 已有的会忽略
func CreateRolePolicy(ctx *gin.Context) {
	var rolePolicy req.RolePolicy
	if err := ctx.ShouldBindJSON(&rolePolicy); err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "参数错误",
		})
		return
	}
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	// 判断角色权限是否存在
	if ok, err := svc.CanAccess(rolePolicy.RoleName, rolePolicy.Url, rolePolicy.Method); err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	} else if ok {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "权限已存在",
		})
		return
	}

	// 权限不存在就创建
	if err := svc.CreateRolePolicy(rolePolicy.RoleName, rolePolicy.Url, rolePolicy.Method); err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})

}

// @Summary 修改角色组权限
func UpdateRolePolicy(ctx *gin.Context) {
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	var rolePolicy req.RolePolicys
	err := ctx.ShouldBindJSON(&rolePolicy)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "参数错误",
		})
		return
	}
	err = svc.UpdateRolePolicy(rolePolicy.OldRoleName, rolePolicy.OldUrl, rolePolicy.OldMethod, rolePolicy.NewRoleName, rolePolicy.NewUrl, rolePolicy.NewMethod)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// 删除角色组策略
func DeleteRolePolicy(ctx *gin.Context) {
	var rolePolicy req.RolePolicy
	err := ctx.ShouldBindJSON(rolePolicy)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	err = svc.DeleteRolePolicy(rolePolicy.RoleName, rolePolicy.Url, rolePolicy.Method)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "删除失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// @Summary 验证用户权限
func VerifyUserRole(ctx *gin.Context) {

	rolePolicy := &req.RolePolicy{}
	err := ctx.ShouldBindJSON(rolePolicy)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	svc := service.NewCasbinService(global.Enforcer, global.Adapter)
	ok, err := svc.CanAccess(rolePolicy.RoleName, rolePolicy.Url, rolePolicy.Method)
	if err == nil && ok {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "验证通过",
		})
		return
	}
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "系统错误",
			"data": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 500,
		"msg":  "验证不通过",
	})
}
