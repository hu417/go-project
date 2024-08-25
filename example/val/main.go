package main

import (
	"net/http"
	"val/api"
	"val/utils"

	"val/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化校验器
	bootstrap.InitValidator()

	r := gin.Default()

	{
		r.POST("login", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "ok")
		})

		r.POST("register", func(ctx *gin.Context) {
			var user api.User
			if err := ctx.ShouldBindJSON(&user); err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"error": utils.GetErrorMsg(user, err),
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"message": "success",
			})
		})
	}

	r.Run(":8081")

}
