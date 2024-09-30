package api

import (
	"jwt-demo/global"
	"jwt-demo/utils/jwt"
	"jwt-demo/utils/resp"

	"github.com/gin-gonic/gin"
)

func RefreshToken(ctx *gin.Context) {
	user_id := ctx.GetInt64("user_id")
	user_name := ctx.GetString("username")
	role := ctx.GetString("role")
	token, err := jwt.GenerateJwt(user_id, user_name, role, global.Jwt_Scret, 3600)
	if err == nil && token != "" {
		resp.Success(ctx, 10000, "token生成成功", gin.H{
			"token": token,
		})

		return
	}
	resp.Fail(ctx, 1000, "token生成失败", nil)

}
