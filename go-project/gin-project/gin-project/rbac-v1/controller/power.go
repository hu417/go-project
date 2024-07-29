package controller

import (
	"github.com/gin-gonic/gin"
	"rbac-v1/common"
	"rbac-v1/model/bo"
	"rbac-v1/model/vo"
	"rbac-v1/service"
)

// 列表
func GetPowerList(ctx *gin.Context) {
	var (
		params = &vo.PowerListRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.Srv().GetPowerList(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "success")
}

// 新增
func PowerCreate(ctx *gin.Context) {
	var (
		params = &vo.PowerCreateRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().PowerCreate(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}

// 修改
func PowerUpdate(ctx *gin.Context) {
	var (
		params = &bo.Power{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().PowerUpdate(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}

// 删除
func PowerDelete(ctx *gin.Context) {
	var (
		params = &vo.PowerDelRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().PowerDelete(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}
