package controller

import (
	"github.com/gin-gonic/gin"
	"rbac-v1/common"
	"rbac-v1/model/vo"
	"rbac-v1/service/casbin"
)

func RbacAuth(ctx *gin.Context) {
	var (
		params = &vo.CasbinAuthRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err = casbin.NewCasbin().Auth(params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "success")
}

func RbacPolicy(ctx *gin.Context) {
	data := casbin.NewCasbin().GetPolicy()

	common.ResponseOk(ctx, data, "success")
}
