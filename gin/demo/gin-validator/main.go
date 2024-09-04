package main

import (
	"github.com/gin-gonic/gin"

	v1 "gin-validator/api/v1"
	"gin-validator/global"
	"gin-validator/internal/bootstrap"
)

func main() {
	// 初始化翻译器
	trans, err := bootstrap.InitTrans("zh")
	if err != nil {
		panic(err)
	}
	global.Trans = trans

	// 初始化路由
	r := gin.Default()

	r.GET("/ping", v1.Ping)
	r.POST("/user/signup", v1.Signup)
	r.POST("/auth/login", v1.Login)

	_ = r.Run(":8081")
}
