package common

import (
	"github.com/gin-gonic/gin"
)

const (
	// 空响应
	RESPONSE_TYPE__RESPONSE_UNKNOW = -1

	// 正常响应
	RESPONSE_TYPE__RESPONSE_OK = 0

	// 正常处理，但没有找到对应的数据
	RESPONSE_TYPE__RESPONSE_NO_DATA = 10000

	// 用户参数校验失败
	RESPONSE_TYPE__RESPONSE_PARAM_INVALID = 90400

	// 常规错误
	RESPONSE_TYPE__RESPONSE_ERROR = 90000

	// 用户登录令牌无效，含过期
	RESPONSE_TYPE__RESPONSE_TOKEN_INVALID = 90401

	// 没有接口权限
	RESPONSE_TYPE__RESPONSE_RBAC_INVALID = 90403

	// Action无效
	RESPONSE_TYPE__RESPONSE_ACTION_INVALID = 90404

	// AccessToken无效
	RESPONSE_TYPE__RESPONSE_ACCESS_TOKEN_INVALID = 90407

	// 请求数超限被拒绝
	RESPONSE_TYPE__RESPONSE_REJECT = 90429

	// 响应异常
	RESPONSE_TYPE__RESPONSE_CRASH = 90500

	// 服务不可用
	RESPONSE_TYPE__RESPONSE_SERVICE_INVALID = 90503
)

type CommonResponse struct {
	Data interface{} `json:"data"`
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
}

func ResponseOk(ctx *gin.Context, data interface{}, msg string) {

	ctx.JSON(200, &CommonResponse{
		Code: int64(RESPONSE_TYPE__RESPONSE_OK),
		Msg:  msg,
		Data: data,
	})
}

func ResponseFailed(ctx *gin.Context, msg string) {
	ctx.JSON(200, &CommonResponse{
		Code: int64(RESPONSE_TYPE__RESPONSE_CRASH),
		Msg:  msg,
		Data: nil,
	})
}

func ResponseTokenInvalid(ctx *gin.Context, msg string) {
	ctx.JSON(200, &CommonResponse{
		Code: int64(RESPONSE_TYPE__RESPONSE_TOKEN_INVALID),
		Msg:  msg,
		Data: nil,
	})
}

func ResponseRbacInvalid(ctx *gin.Context, msg string) {
	ctx.JSON(200, &CommonResponse{
		Code: int64(RESPONSE_TYPE__RESPONSE_RBAC_INVALID),
		Msg:  msg,
		Data: nil,
	})
}

func ResponseParamInvalid(ctx *gin.Context, msg string) {
	ctx.JSON(200, &CommonResponse{
		Code: int64(RESPONSE_TYPE__RESPONSE_PARAM_INVALID),
		Msg:  msg,
		Data: nil,
	})
}
