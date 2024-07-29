package controller

import (
	"github.com/gin-gonic/gin"
	"rbac-v1/common"
	"rbac-v1/model/po"
	"rbac-v1/model/vo"
	"rbac-v1/service"
)

// 列表
func GetOperationList(ctx *gin.Context) {
	var (
		params = &vo.OperationListRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.Srv().GetOperationList(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "success")
}

// 新增
func OperationCreate(ctx *gin.Context) {
	var (
		params = &vo.OperationCreateRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().OperationCreate(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}

// 修改
func OperationUpdate(ctx *gin.Context) {
	var (
		params = &po.Operation{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().OperationUpdate(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}

// 删除
func OperationDelete(ctx *gin.Context) {
	var (
		params = &vo.OperationDelRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = service.Srv().OperationDelete(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}
