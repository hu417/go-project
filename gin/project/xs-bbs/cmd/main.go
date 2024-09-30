package main

import (
	"log"
	"time"

	"xs-bbs/global"
	"xs-bbs/internal/router"
	"xs-bbs/pkg/cache"
	"xs-bbs/pkg/conf"
	"xs-bbs/pkg/db"
	"xs-bbs/pkg/logger"

	"github.com/fvbock/endless"
	"go.uber.org/zap"
)

func main() {
	// 1. init config
	file := "./etc/setting.yaml"
	config := conf.InitViper(file)
	if config == nil {
		log.Fatal("conf.Build failed")
	}
	global.Conf = config

	// 2. init logger
	logger := logger.InitLogger()
	if logger == nil {
		log.Fatal("log.Build failed, err")
	}
	global.Log = logger

	// 3. init gorm client
	db, err := db.InitMysql(global.Conf)
	if err != nil || db == nil {
		global.Log.Error("database.Build failed", zap.Error(err))
		return
	}
	global.DB = db
	defer func() {
		sql, _ := global.DB.DB()
		if err := sql.Close(); err != nil {
			global.Log.Error("database.Build failed", zap.Error(err))
			return
		}
	}()

	// 4. init gorm client
	rbd, err := cache.InitRedis(global.Conf)
	if err != nil {
		global.Log.Error("cache.Build failed", zap.Error(err))
		return
	}
	global.RedisClient = rbd
	defer func() {
		if err := global.RedisClient.Close(); err != nil {
			global.Log.Error("cache.Build failed", zap.Error(err))
			return
		}
	}()

	r := router.NewHttpServer()

	// 7. 启动服务（优雅关机）
	/* 默认endless服务会监听下列信号：
	syscall.SIGHUP，syscall.SIGUSR1，syscall.SIGUSR2，syscall.SIGINT，syscall.SIGTERM和syscall.SIGTSTP
	接收到 SIGHUP 信号将触发`fork/restart` 实现优雅重启（kill -1 pid会发送SIGHUP信号）
	接收到 syscall.SIGINT或syscall.SIGTERM 信号将触发优雅关机
	接收到 SIGUSR2 信号将触发HammerTime
	SIGUSR1 和 SIGTSTP 被用来触发一些用户自定义的hook函数
	*/
	server := endless.NewServer(global.Conf.System.Port, r)

	// 10秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	server.ReadTimeout = 10 * time.Second  // 10s
	server.WriteTimeout = 10 * time.Second // 10s
	server.MaxHeaderBytes = 1 << 20        // 1MB

	if err := server.ListenAndServe(); err != nil {
		global.Log.Fatal("server start fail: ", zap.Error(err))
	}
	// if err := server.Shutdown(ctx); err != nil {
	// 	zap.L().Fatal("Server Shutdown", zap.Error(err))
	// }

	global.Log.Info("Server start ...")

}
