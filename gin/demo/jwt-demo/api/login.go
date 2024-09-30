package api

import (
	"jwt-demo/utils/jwt"
	"jwt-demo/utils/resp"

	"github.com/gin-gonic/gin"
)

func LoginHandler(ctx *gin.Context) {

	token, err := jwt.GenerateJwt(1, "admin", "admin", "123456", 60*60*24)
	if err != nil {
		resp.Fail(ctx, 10001, "系统错误", nil)

		return
	}
	if token == "" {
		resp.Fail(ctx, 10001, "token生成失败", nil)

		return
	}
	resp.Success(ctx, 10001, "ok", gin.H{
		"token": token,
	})

}
