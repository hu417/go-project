package controller

import (
	"github.com/gin-gonic/gin"
	"rbac-v1/common"
	"rbac-v1/model/bo"
	"rbac-v1/model/vo"
	"rbac-v1/service"
)

// 列表
func GetUserList(ctx *gin.Context) {
	var (
		params = &vo.UserListRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.Srv().GetUserList(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "success")
}

func GetUserRoleAndPower(ctx *gin.Context) {
	var (
		params = &vo.GetUserRoleAndPowerRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.Srv().GetUserRoleAndPower(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "success")
}

// 新增
func UserCreate(ctx *gin.Context) {
	var (
		params = &vo.UserCreateRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().UserCreate(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}

// 修改
func UserUpdate(ctx *gin.Context) {
	var (
		params = &bo.UserCreate{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().UserUpdate(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}

// 删除
func UserDelete(ctx *gin.Context) {
	var (
		params = &vo.UserDelRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().UserDelete(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}
