package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
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
	resp := Response{Code: errorCodeMap[ErrSuccess], Message: ErrSuccess.Error(), Data: data}
	if _, ok := errorCodeMap[ErrSuccess]; !ok {
		resp = Response{Code: 0, Message: "", Data: data}
	}
	ctx.JSON(http.StatusOK, resp)
}

// HandleError 处理错误
func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := Response{Code: errorCodeMap[err], Message: err.Error(), Data: data}
	if _, ok := errorCodeMap[err]; !ok {
		resp = Response{Code: 500, Message: "unknown error", Data: data}
	}
	ctx.JSON(httpCode, resp)
}

// Error 自定义错误
type Error struct {
	Code    int
	Message string
}

// errorCodeMap 错误码
var errorCodeMap = map[error]int{}

// newError 新建错误
func newError(code int, msg string) error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	return err
}

// Error 实现error接口
func (e Error) Error() string {
	return e.Message
}
