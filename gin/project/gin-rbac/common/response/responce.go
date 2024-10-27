package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 基础响应结构体
type Response struct {
	Code    int         `json:"code,omitempty"` // 应用级别的错误代码，可选
	Message string      `json:"msg"`            // 响应消息，提供人可读的响应信息
	Data    interface{} `json:"data,omitempty"` // 实际响应数据，类型为 interface{}，意味着它可以是任何类型，具体类型在使用时确定。
}

// Result 返回通用响应结构体
func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// Ok 成功响应，作用：创建成功|更新成功|删除成功|获取成功|获取失败
func Ok(c *gin.Context, code int, msg string, data interface{}) {
	if msg == "" {
		msg = "success"
	}
	Result(c, code, msg, data)
}

// OkWithData 成功响应数据，作用： 获取成功数据
func OkWithData(c *gin.Context, code int, data interface{}) {
	Ok(c, code, "success", data)
}

// OkWithMsg 成功响应消息，作用：创建成功|更新成功|删除成功|获取成功|获取失败
func OkWithMsg(c *gin.Context, code int, msg string) {
	Ok(c, code, msg, map[string]any{})
}

// Fail 错误响应
func Fail(c *gin.Context, code int, msg string) {
	Result(c, code, msg, map[string]any{})
}

/*// FailWithMsg 错误响应消息
func FailWithMsg(c *gin.Context, status int, code int, msg string, data interface{}) {
	//Result(c, 400, msg, map[string]any{})
	Result(c, status, code, msg, map[string]any{})
}*/

// FailWithData
func FailWithData(c *gin.Context, code int, msg string, data interface{}) {
	Result(c, code, msg, data)
}

// Unauthorized 返回未授权的错误响应
func Unauthorized(c *gin.Context, message string) {
	Result(c, http.StatusUnauthorized, message, map[string]any{})
	//Result(c, status, message, map[string]interface{}{})
}

// Forbidden 返回禁止访问的错误响应
func Forbidden(c *gin.Context, message string) {
	Result(c, http.StatusForbidden, message, map[string]any{})
}
