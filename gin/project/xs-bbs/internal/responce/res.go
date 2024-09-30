package responce

import (
	"net/http"
	
	"xs-bbs/internal/errs"

	"github.com/gin-gonic/gin"
)

// Response .
type Response struct {
	Code errs.ResCode   `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// RespSuccess 响应成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Code: errs.CodeSuccess,
		Msg:  errs.CodeSuccess.Msg(),
		Data: data,
	})
}

// RespError 响应失败，携带状态及对应信息
func Error(c *gin.Context, code errs.ResCode) {
	c.JSON(http.StatusOK, &Response{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// RespErrorWithMsg 响应失败，携带状态+其他自定义信息
func ErrorWithMsg(c *gin.Context, code errs.ResCode, msg string) {
	c.JSON(http.StatusOK, &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
