package controller

import (
	"github.com/gin-gonic/gin"
	"rbac-v1/common"
	"rbac-v1/model/bo"
	"rbac-v1/model/vo"
	"rbac-v1/service"
)

// 列表
func GetRoleList(ctx *gin.Context) {
	var (
		params = &vo.RoleListRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.Srv().GetRoleList(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "success")
}

// 新增
func RoleCreate(ctx *gin.Context) {
	var (
		params = &vo.RoleCreateRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().RoleCreate(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}

// 修改
func RoleUpdate(ctx *gin.Context) {
	var (
		params = &bo.Role{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().RoleUpdate(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}

// 删除
func RoleDelete(ctx *gin.Context) {
	var (
		params = &vo.RoleDelRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().RoleDelete(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}
