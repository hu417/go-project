package controller

import (
	"github.com/gin-gonic/gin"
	"rbac-v1/common"
	"rbac-v1/common/constants"
	"rbac-v1/model/vo"
	"rbac-v1/service/login"
)

func LoginByPwd(ctx *gin.Context) {
	var (
		params = &vo.LoginRequest{}
	)

	err := ctx.Bind(&params)
	if err != nil {
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := login.NewLogin().LoginByPwd(ctx, params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "success")
}

func CheckToken(ctx *gin.Context) {
	var (
		params = &vo.CheckTokenRequest{}
	)

	token := ctx.GetHeader(constants.TOKEN_HEADER_KEY)
	if len(token) == 0 {
		common.ResponseParamInvalid(ctx, "token is empty")
		return
	}
	params.Token = token
	data, err := login.NewLogin().CheckToken(params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "success")
}
