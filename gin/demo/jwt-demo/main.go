package main

import "jwt-demo/router"

func main() {
	// 1. 初始化路由
	router := router.InitRouter()
	// 2. 启动服务
	router.Run(":8081")
}
