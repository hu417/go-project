package res

import (
	"github.com/gin-gonic/gin"
)

// 响应结构体
type Response struct {
	BizCode  int         `json:"biz_code"` // 自定义业务状态码
	Succcess bool        `json:"success"`  // 是否成功状态
	Message  string      `json:"message"`  // 信息描述
	Error    interface{} `json:"error"`    // 错误信息
	Data     interface{} `json:"data"`     // 数据信息
}

// Success = true 响应成功
func Success(c *gin.Context, http_code, biz_code int, message string, data interface{}) {
	c.JSON(http_code, Response{
		biz_code,
		true,
		message,
		nil,
		data,
	})
}

// Success = false 响应失败
func Fail(c *gin.Context, http_code, biz_code int, message string, erros interface{}) {
	c.JSON(http_code, Response{
		biz_code,
		false,
		message,
		erros,
		nil,
	})
}

// FailByError 失败响应 返回自定义错误的错误码、错误信息
func FailByError(c *gin.Context, httpCode, biz_code int, message string, err *CustomError) {
	Fail(c, httpCode, biz_code, message, err)
}

// ValidateFail 请求参数验证失败
func ValidateFail(c *gin.Context, http_code, biz_code int, message string) {
	Fail(c, http_code, biz_code, message, Errors.ValidateError)
}

// BusinessFail 业务逻辑失败
func BusinessFail(c *gin.Context, http_code int, biz_code int, message string) {
	Fail(c, http_code, biz_code, message, Errors.BusinessError)
}
