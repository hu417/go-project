package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bluebell/config"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/pkg/logger"
	"bluebell/pkg/snowflake"
	"bluebell/pkg/valida"
	"bluebell/router"

	"go.uber.org/zap"
)

//	@title			bluebell项目接口文档
//	@version		1.0
//	@description	Go web开发进阶项目实战课程bluebell

//	@contact.name	liwenzhou
//	@contact.url	http://www.liwenzhou.com

// @host		127.0.0.1:8084
// @BasePath	/api/v1
func main() {
	// 判断命令行参数
	var setting_file string
	switch {
	case len(os.Args) < 2:
		// fmt.Println("need config file.eg: bluebell config.yaml")
		setting_file = "./bluebell/etc/config.yaml"
	default:
		setting_file = os.Args[1]
	}

	// 加载配置
	if err := config.Init(setting_file); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	fmt.Println("load config succes")

	// 初始化日志
	if err := logger.Init(config.Conf.LogConfig, config.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	// 初始化db
	if err := mysql.Init(config.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接

	// 初始化redis
	if err := redis.Init(config.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(config.Conf.StartTime, config.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 初始化gin框架内置的校验器使用的翻译器
	if err := valida.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
		return
	}

	// 注册路由
	r := router.SetupRouter(config.Conf.Mode, zap.L())
	// err := r.Run(fmt.Sprintf(":%d", config.Conf.Port))
	// if err != nil {
	// 	fmt.Printf("run server failed, err:%v\n", err)
	// 	return
	// }

	srv := &http.Server{
		// Gin运行的监听端口
		Addr: ":8080",
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
