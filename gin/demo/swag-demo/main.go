package main

import (
	"log"
	"time"

	"swag-demo/router"

	// 匿名导入生成的接口文档包
	_ "swag-demo/docs"

	"github.com/fvbock/endless"
)

//	@title			标题：Swagger Demo API
//	@version		1.0
//	@description	描述：这是swagger demo API

//	@contact.name	接口联系人信息
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization

//	@securityDefinitions.basic	BasicAuth
//	@in							header
//	@name						Auth

//	@Schemes	http https
//	@Host		localhost:8081
//	@BasePath	/api/v1
func main() {

	// 初始化路由
	r := router.InitRouter()

	// 启动服务
	// 默认endless服务会监听下列信号：
	// syscall.SIGHUP，syscall.SIGUSR1，syscall.SIGUSR2，syscall.SIGINT，syscall.SIGTERM和syscall.SIGTSTP
	// 接收到 SIGHUP 信号将触发`fork/restart` 实现优雅重启（kill -1 pid会发送SIGHUP信号）
	// 接收到 syscall.SIGINT或syscall.SIGTERM 信号将触发优雅关机
	// 接收到 SIGUSR2 信号将触发HammerTime
	// SIGUSR1 和 SIGTSTP 被用来触发一些用户自定义的hook函数
	server := endless.NewServer("localhost:8081", r)
	server.ReadTimeout = 10 * time.Second  // 10s
	server.WriteTimeout = 10 * time.Second // 10s
	server.MaxHeaderBytes = 1 << 20        // 1MB

	// 启动服务
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Listen: ", err)
	}
	// log.Println("Server exiting")
}
