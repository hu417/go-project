package routers

// import (
// 	"context"
// 	"fmt"
// 	"gvb_server/global"
// 	"net/http"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/spf13/viper"
// )

// // 可插拔路由模块
// // 安装依赖: go get github.com/gin-gonic/gin

// type IFnReqinstRoute = func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup)

// var (
// 	gfnRouters []IFnReqinstRoute
// )

// func RegistRoute(fn IFnReqinstRoute) {
// 	if fn == nil {
// 		return
// 	}
// 	gfnRouters = append(gfnRouters, fn)
// }

// func InitRouter() {
// 	// =============================================================================
// 	// = 创建监听ctrl + c, 应用退出信号的上下文
// 	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
// 	defer cancelCtx()

// 	// =============================================================================
// 	// = 初始化gin框架, 并注册相关路由
// 	r := gin.Default()

// 	// 不需要认证
// 	rgPublic := r.Group("/api/v1/public")
// 	// 需要认证
// 	rgAuth := r.Group("/api/v1")

// 	// // 初始基础平台的路由
// 	InitBasePlatformRouter()

// 	// 开始注册系统各模块对应的路由信息
// 	for _, fnRegistRoute := range gfnRouters {
// 		fnRegistRoute(rgPublic, rgAuth)
// 	}

// 	// 设置端口
// 	stPort := viper.GetString("Server.Port")
// 	if stPort == "" {
// 		stPort = "8090"
// 	}

// 	// 启动服务
// 	//r.Run(fmt.Sprintf(":%s", stPort))// =============================================================================
// 	// = 创建web server
// 	server := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", stPort),
// 		Handler: r,
// 	}

// 	// =============================================================================
// 	// = 启动一个goroutine来开启web服务, 避免主线程的信号监听被阻塞
// 	go func() {
// 		fmt.Printf("%s\n", fmt.Sprintf("Start Server Listen: %s", stPort))
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			fmt.Printf("%s\n", fmt.Sprintf("Start Server Error: %s", err.Error()))
// 			return
// 		}
// 	}()

// 	// =============================================================================
// 	// = 等待停止服务的信号被触发
// 	<-ctx.Done()
// 	// cancelCtx()

// 	// =============================================================================
// 	// = 关闭Server， 5秒内未完成清理动作会直接退出应用
// 	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelShutdown()

// 	if err := server.Shutdown(ctx); err != nil {
// 		fmt.Printf("Stop Server Error: %s\n", err.Error())
// 		return
// 	}

// 	fmt.Println("Stop Server Success")
// }

// // ! 初始化基础平台相关路由信息
// func InitBasePlatformRouter() {
// 	//
// 	InitUserRouter()
// 	global.Logger.Info("初始化路由成功")
// }
