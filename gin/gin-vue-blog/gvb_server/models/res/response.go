package res

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

const (
	Success = 0
	Error   = 7
)

// 响应函数封装
func Result(code int, data any, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// ############## 响应成功
// 返回数据，日志
func Ok(data any, msg string, c *gin.Context) {
	Result(Success, data, msg, c)
}

// 什么都不返回
func OkWith(c *gin.Context) {
	Result(Success, map[string]any{}, "成功", c)
}

// 只返回数据
func OkWithData(data any, c *gin.Context) {
	Result(Success, data, "成功", c)
}

// 只返回日志
func OkWithMessage(msg string, c *gin.Context) {
	Result(Success, map[string]any{}, msg, c)
}

// ############### 响应错误
// 返回数据，日志
func Fail(data any, msg string, c *gin.Context) {
	Result(Error, data, msg, c)
}

// 只返回日志
func FailWithMessage(msg string, c *gin.Context) {
	Result(Error, map[string]any{}, msg, c)
}

// 判断是业务错误还是请求错误
func FailWithCode(code ErrorCode, c *gin.Context) {
	msg, ok := ErrorMap[code]
	if ok {
		Result(int(code), map[string]any{}, msg, c)
		return
	}
	Result(Error, map[string]any{}, "未知错误", c)
}
