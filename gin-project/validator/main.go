package main

import (
	"github.com/gin-gonic/gin"

	v1 "gin-validator/api/v1"
	"gin-validator/internal/bootstrap"
)

func main() {
	// 初始化翻译器
	bootstrap.InitTrans("zh")

	// 初始化路由
	r := gin.Default()

	r.GET("/ping", v1.Ping)
	r.POST("/user/add", v1.AddUser)
	r.POST("/auth/login", v1.Login)

	_ = r.Run(":8080")
}
