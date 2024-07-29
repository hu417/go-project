package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response(ctx *gin.Context, status_code int, code int, msg string, data gin.H) {
	ctx.JSON(status_code, gin.H{"code": code, "msg": msg, "data": data})
}

// 定义两个常用的返回

// 成功
func Success(ctx *gin.Context, code int, msg string, data gin.H) {
	Response(ctx, http.StatusOK, code, msg, data)
}

// 失败
func Fail(ctx *gin.Context, code int, msg string, data gin.H) {
	//失败也是返回200，然后返回400的状态码
	Response(ctx, http.StatusOK, code, msg, data)
}
