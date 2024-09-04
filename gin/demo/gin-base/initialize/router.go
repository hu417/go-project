package initialize

import (
	"fmt"
	"net/http"
	"time"

	"gin-base/global"
	"gin-base/router"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

type server interface {
	ListenAndServe() error
}

func InitRouter() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 公共路由组
	PublicGroup := r.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}

	{
		// 注册基础功能路由 不做鉴权
		systemRouter := router.RouterGroupApp.System
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
	}

	address := fmt.Sprintf(":%d", global.Conf.App.Port)
	s := SetUpServer(address, r)

	// 保证文本顺序输出
	time.Sleep(10 * time.Microsecond)

	s.ListenAndServe()
}

func SetUpServer(address string, router *gin.Engine) server {
	// 使用endless库创建一个HTTP服务器，其中address是服务器的监听地址（如:8080），router是HTTP请求路由器。
	s := endless.NewServer(address, router)

	// 设置HTTP请求头的读取超时时间为20秒，如果在20秒内未读取到请求头，则会返回一个超时错误。
	s.ReadHeaderTimeout = 20 * time.Second

	// 设置HTTP响应体的写入超时时间为20秒，如果在20秒内未将响应体写入完成，则会返回一个超时错误。
	s.WriteTimeout = 20 * time.Second

	// 设置HTTP请求头的最大字节数为1MB。如果请求头超过1MB，则会返回一个错误。
	s.MaxHeaderBytes = 1 << 20

	return s
}
