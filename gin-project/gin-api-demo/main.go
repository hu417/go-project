package main

import (
	"gin-api-demo/router"
	"gin-api-demo/router/fileload"
	"gin-api-demo/router/login"
)

func main() {
	// 注册路由
	r := router.Init(login.Router, fileload.Router)
	r.Run(":8080")
}
