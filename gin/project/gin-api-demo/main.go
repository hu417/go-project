package main

import (
	"gin-api-demo/bootstrap"
	"gin-api-demo/global"
	"gin-api-demo/pkg/utils"
	"gin-api-demo/router"
	"log"
	"time"

	"github.com/fvbock/endless"
	"go.uber.org/zap"
)

func main() {
	//
	dir := utils.GetPath() + "/gin/project/gin-api-demo"

	// 初始化配置
	bootstrap.InitConfig(dir + "/etc/setting.yaml")

	// 初始化日志
	global.Log = bootstrap.InitializeLog(dir)
	// 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(global.Log)

	// 初始化数据库
	global.DB = bootstrap.InitializeDB()
	// 程序关闭前，释放数据库连接
	defer func() {
		if global.DB != nil {
			db, _ := global.DB.DB()
			db.Close()
		}
	}()

	// 初始化翻译器
	trans, err := bootstrap.InitTrans("zh")
	if err != nil {
		panic(err)
	}
	global.Trans = trans

	// 初始化缓存
	global.Redis = bootstrap.InitRedis(global.Conf)
	if global.Redis == nil {
		panic("redis init fail")
	}
	// 程序关闭前，释放redis连接
	defer func() {
		if global.Redis != nil {
			global.Redis.Close()
		}
	}()

	// 创建默认的路由
	r := router.InitRouter(global.Conf.App.Env)

	// 启动服务
	// 默认endless服务会监听下列信号：
	// syscall.SIGHUP，syscall.SIGUSR1，syscall.SIGUSR2，syscall.SIGINT，syscall.SIGTERM和syscall.SIGTSTP
	// 接收到 SIGHUP 信号将触发`fork/restart` 实现优雅重启（kill -1 pid会发送SIGHUP信号）
	// 接收到 syscall.SIGINT或syscall.SIGTERM 信号将触发优雅关机
	// 接收到 SIGUSR2 信号将触发HammerTime
	// SIGUSR1 和 SIGTSTP 被用来触发一些用户自定义的hook函数
	server := endless.NewServer(global.Conf.App.Port, r)
	server.ReadTimeout = 10 * time.Second  // 10s
	server.WriteTimeout = 10 * time.Second // 10s
	server.MaxHeaderBytes = 1 << 20        // 1MB

	// 启动服务
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server start fail: ", err)
	}
}
