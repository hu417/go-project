package main

import (
	"zap-demo/router"
	"zap-demo/utils"
)

func main() {
	// 初始化日志
	utils.InitLogger("dev", "console")

	// 初始化路由
	r := router.InitRouter()

	// 启动服务
	r.Run(":8081")

}
