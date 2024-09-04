package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rbac-v1/config"
	"rbac-v1/controller"
	"rbac-v1/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

func main() {
	r := gin.Default()

	//service初始化
	service.New()
	//路由
	controller.InitApiRouter(r)

	//gin启动
	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen: %s\n", err.Error())
		}
	}()
	//gin优雅退出
	//等待终端信号，优雅关闭所有server及DB
	//这里也是个阻塞的过程，也就是说，没有监听到中断信号时，会一直阻塞在这里，为了主线程不退出
	quit := make(chan os.Signal, 1)
	//系统的终端信号  INT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() //释放ctx
	//超时时间5秒，如果5秒还没处理完之前的请求，则强制退出
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("gin server关闭异常:", err)
	}
	logger.Info("gin server退出成功")
	//关闭service
	if err := service.Srv().Close(); err != nil {
		logger.Error("service关闭异常:", err)
	}
}
