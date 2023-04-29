package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"kubea-demo/config"     // 引入配置依赖
	"kubea-demo/controller" // 引入路由依赖
	"kubea-demo/service"    // 引入k8s

	"github.com/gin-gonic/gin" // 引用第三方gin依赖
)

func main() {
	// 初始化gin对象
	r := gin.Default()
	// 初始化k8s配置
	service.K8s.Init()
	// 初始化路由规则
	controller.Router.InitApiRouter(r)

	// 启动gin server服务
	srv := &http.Server{
		Addr:    config.ListenAddr, // 启动地址
		Handler: r,                 // 路由
	}

	// 使用协程监听服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: 启动失败! err => %s", err)
		}
	}()

	// 优雅停止 //
	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 设置ctx关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// cancel()释放ctx
	defer cancel()

	// 正常关闭gin server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Gin Server 关闭异常! err => %s", err)
	}
	log.Println("Gin Server 正常退出")
}
