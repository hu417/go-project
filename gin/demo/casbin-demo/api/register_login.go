package api

import "github.com/gin-gonic/gin"

// @Summary 注册
func RegisterHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "注册成功",
		"data":    "",
	})
}

// @Summary 登陆
func LoginHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "登陆成功",
		"data":    "",
	})
}
