package res

import (
	"blue-bell/controller/e"

	"github.com/gin-gonic/gin"
)

/*
定义响应信息
*/
type ResponseData struct {
	Code e.ResCode   `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// 响应错误信息
func ResponseError(c *gin.Context, http_code int, biz_code e.ResCode, msg interface{}) {
	c.JSON(http_code, &ResponseData{
		Code: biz_code,
		Msg:  msg,
		Data: nil,
	})
}

// 响应成功信息
func ResponseSuccess(c *gin.Context, http_code int, msg interface{}, data interface{}) {
	c.JSON(http_code, &ResponseData{
		Code: e.CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}
