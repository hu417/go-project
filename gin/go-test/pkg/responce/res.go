package responce

import (
	"go-test/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// HandleSuccess 处理成功
func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := Response{Code: errs.ErrorCodeMap[errs.ErrSuccess], Message: errs.ErrSuccess.Error(), Data: data}
	if _, ok := errs.ErrorCodeMap[errs.ErrSuccess]; !ok {
		resp = Response{Code: 0, Message: "", Data: data}
	}
	ctx.JSON(http.StatusOK, resp)
}

// HandleError 处理错误
func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := Response{Code: errs.ErrorCodeMap[err], Message: err.Error(), Data: data}
	if _, ok := errs.ErrorCodeMap[err]; !ok {
		resp = Response{Code: 500, Message: "unknown error", Data: data}
	}
	ctx.JSON(httpCode, resp)
}
