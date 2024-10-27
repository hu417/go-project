package main

import (
	"log"
	"strings"

	// 匿名导入生成的接口文档包
	_ "gin-rbac/docs"

	"gin-rbac/bootstrap"
	"gin-rbac/cmd/app"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
//go:generate go mod tidy
//go:generate go mod download

// @title						Gin-Vue3-RBAC Swagger API接口文档
// @version					v1.1.1
// @description				使用gin+vue3进行极速开发的全栈开发rbac权限管理基础平台.
// @BasePath					/api/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	app := app.NewApp()
	// 初始化 JWT 和 DB 实例
	sys,  err := bootstrap.Init()
	if err != nil {
		log.Fatalf("Failed to initialize system components: %v", err)
	}

	// 初始化应用并注册路由
	app.Initialize()

	// 启动服务器
	address := sys.Addr()
	// 打印swagger 地址
	colonIndex := strings.Index(address, ":")
	port := address[colonIndex+1:]
	log.Printf("Swagger UI: http://127.0.0.1:%s/swagger/index.html\n", port)
	if err := app.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
