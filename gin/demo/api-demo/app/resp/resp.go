package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 定义一个统一的返回值结构体
type Response struct {
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		0,
		data,
		"ok",
	})
}

// 权限验证失败
func PermissionDenied(c *gin.Context) {
	Error(c, HaveNoPermission, ErrorMap[HaveNoPermission])
}

// 业务失败
func BusinessFail(c *gin.Context, msg string) {
	Error(c, BusinessError, msg)
}

// Error 失败返回
func Error(c *gin.Context, errorCode int, msg string) {
	c.JSON(http.StatusOK, &Response{
		errorCode,
		nil,
		msg,
	})
}

// 参数格式失败
func ValidateFailed(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, &Response{
		1,
		nil,
		msg,
	})
}
