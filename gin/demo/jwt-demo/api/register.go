package api

import (
	"jwt-demo/utils/resp"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(ctx *gin.Context) {
	resp.Success(ctx, 1000, "注册成功", nil)

}
