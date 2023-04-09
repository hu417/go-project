package bin

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/routers"
)

func Start() {
	// 初始化系统配置
	core.InitConf()
	core.IninViper()

	// 初始化日志组件
	global.Logger = core.InitLogger()

	// 初始化Elastic
	global.ESClient = core.InitEs()

	// 初始化Grom
	global.DB = core.InitGorm()

	// =============================================================================
	// = 创建监听ctrl + c, 应用退出信号的上下文
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	// 初始化系统路由
	global.Router = routers.InitRouter()

	// 启动服务
	addr := global.Config.System.Addr()

	global.Logger.Infof("服务启动中,监听地址: %s", addr)
	// global.Router.Run(addr)

	// =============================================================================
	// = 创建web server
	server := &http.Server{
		Addr:    addr,
		Handler: global.Router,
	}

	// = 启动一个goroutine来开启web服务, 避免主线程的信号监听被阻塞
	go func() {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("%s\n", fmt.Sprintf("Start Server Error: %s", err.Error()))
			return
		}
	}()

	// =============================================================================
	// = 等待停止服务的信号被触发
	<-ctx.Done()
	// cancelCtx()

	// =============================================================================
	// = 关闭Server， 5秒内未完成清理动作会直接退出应用
	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Error("Stop Server Error: %s\n", err.Error())
		return
	}

	global.Logger.Info("Stop Server Success")

}
