package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ginblog/bootstrap"
	"ginblog/global"
	"ginblog/router"
	"ginblog/utils"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	conf, err := bootstrap.InitConf("./etc/setting.yaml")
	if err != nil {
		panic(err)
	}
	global.Conf = conf

	// 初始化日志
	utils.InitLogger("devops", "console")

	// 初始化数据库
	db := bootstrap.InitMysql(&global.Conf.Mysql)
	if db == nil {
		panic("数据库初始化失败")
	}
	defer func() {
		d, _ := db.DB()
		if err := d.Close(); err != nil {
			panic(err)
		}
	}()
	global.DB = db

	// 初始化路由
	r := router.InitRouter()
	srv := &http.Server{
		// Gin运行的监听端口
		Addr: global.Conf.Server.Port,
		// 要调用的处理程序，http.DefaultServeMux如果为nil
		Handler: r,
		// ReadTimeout是读取整个请求（包括正文）的最长持续时间。
		ReadTimeout: 5 * time.Second,
		// WriteTimeout是超时写入响应之前的最长持续时间
		WriteTimeout: 10 * time.Second,
		// MaxHeaderBytes控制服务器解析请求标头的键和值（包括请求行）时读取的最大字节数 (通常情况下不进行设置)
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务（优雅关机）
	// 方式一
	// 默认endless服务器会监听下列信号：
	// syscall.SIGHUP，syscall.SIGUSR1，syscall.SIGUSR2，syscall.SIGINT，syscall.SIGTERM和syscall.SIGTSTP
	// 接收到 SIGHUP 信号将触发`fork/restart` 实现优雅重启（kill -1 pid会发送SIGHUP信号）
	// 接收到 syscall.SIGINT或syscall.SIGTERM 信号将触发优雅关机
	// 接收到 SIGUSR2 信号将触发HammerTime
	// SIGUSR1 和 SIGTSTP 被用来触发一些用户自定义的hook函数
	// if err := endless.ListenAndServe(":8080", r); err!=nil{
	// 	zap.L().Sugar().Errorf("listen: %s\n", err)
	// }

	// 方式二
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			zap.L().Sugar().Errorf("listen: %s\n", err)

		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Sugar().Fatal("Server forced to shutdown:", err)
	}

	zap.L().Info("Server exiting")
}
