package main

import (
	"context"
	"fmt"
	"time"

	"blue-bell/config"
	"blue-bell/dao/mysql"
	"blue-bell/dao/redis"
	"blue-bell/global"
	"blue-bell/pkg/logger"
	"blue-bell/pkg/val"
	"blue-bell/router"

	"github.com/fvbock/endless"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置，用viper从配置文件中读取信息
	// conf := config.InitViper("./gin/project/blue-bell/etc/setting.yaml")
	conf := config.InitViper("./etc/setting.yaml")
	if conf == nil {
		fmt.Println("config Init failed")
		return
	}
	global.Conf = conf

	// 2. 初始化日志
	if lg := logger.InitializeLog(); lg == nil {
		zap.L().Debug("logger init failed...")
		return
	}
	zap.L().Debug("logger success init...")

	// 3. 初始化MySQL连接
	db := mysql.InitializeDB()
	if db == nil {
		zap.L().Debug("mysql init failed...")
		return
	}
	zap.L().Debug("mysql init success...")
	global.DB = db
	defer func() {
		sql, _ := db.DB()
		if err := sql.Close(); err != nil {
			zap.L().Debug("mysql close failed")
			return
		}
	}()

	// 4. 初始化Redis连接
	cli := redis.InitRedis()
	if cli == nil {
		zap.L().Debug("redis init failed")
		return
	}
	zap.L().Debug("redis init success...")
	global.RedisCli = cli
	defer func() {
		if err := cli.Close(); err != nil {
			zap.L().Debug("redis close failed")
			return
		}
	}()

	// 5. 初始化validator翻译器
	trans, err := val.InitTrans("zh")
	if err != nil {
		zap.L().Debug("controller initTrans init failed", zap.Error(err))
		return
	}
	global.Trans = trans

	// 6. 创建一个10秒超时的context
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 7. 注册路由
	r := router.Setup(global.Conf.System.Mode, 10)

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
		zap.L().Fatal("server start fail: ", zap.Error(err))
	}
	// if err := server.Shutdown(ctx); err != nil {
	// 	zap.L().Fatal("Server Shutdown", zap.Error(err))
	// }

	zap.L().Info("Server start ...")
}
