package main

import (
	"fmt"
	"html/template"
	"net/http"

	"captcha-demo/bootstrap"
	"captcha-demo/global"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化redis
	global.RDS = bootstrap.InitRedis()
	if global.RDS == nil {
		panic("redis init failed")
	}

	// 初始化验证码
	captcha := bootstrap.InitCaptcha("math", global.RDS)

	// 初始化路由
	app := gin.Default()

	// 加载模板文件
	app.LoadHTMLGlob("./gin-project/captcha-demo/templates/*")

	// 生成验证码
	app.GET("/generate", func(ctx *gin.Context) {
		verifyId, b64s, verifyValue, err := captcha.Generate()
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "验证码生成失败",
			})
			return
		}
		fmt.Println(verifyValue)

		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"verifyId":    verifyId,
			"b64s":        template.URL(b64s),
			"verifyValue": verifyValue,
		})
	})

	// 验证
	app.POST("/verify", func(ctx *gin.Context) {
		verifyValue := ctx.PostForm("verifyValue")
		verifyId := ctx.PostForm("verifyId")
		// 参数说明: id 验证码id, verifyValue 验证码的值, true: 验证成功后是否删除原来的验证码
		result := captcha.Verify(verifyId, verifyValue, true)

		ctx.JSON(http.StatusOK, gin.H{
			"verifyId":    verifyId,
			"verifyValue": verifyValue,
			"result":      result,
		})
	})

	// 启动服务
	app.Run(":8081")
}
