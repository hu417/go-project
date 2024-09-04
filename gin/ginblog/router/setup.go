package router

import (
	"net/http"

	"ginblog/middleware"
	"ginblog/router/back"
	"ginblog/router/front"
	"ginblog/utils/resp"

	"github.com/gin-gonic/gin"
)

func InitRouter(run_model string) *gin.Engine {
	// 设置运行模式
	gin.SetMode(run_model)

	// 创建路由
	r := gin.New()
	// 设置信任网络 []string
	// nil 为不计算，避免性能消耗，上线应当设置
	_ = r.SetTrustedProxies(nil)

	// r.HTMLRender = createMyRender()
	// 日志
	r.Use(middleware.Logger(), middleware.Recovery(true))
	// 跨域
	r.Use(middleware.Cors())

	// 静态文件
	r.StaticFile("/favicon.ico", "gin-project/ginblog/static/favicon.ico")

	// test
	r.GET("/ping", func(c *gin.Context) {
		resp.Success(c, 200, "pong", gin.H{})
	})

	{
		// 404
		r.NoRoute(func(c *gin.Context) {
			resp.Fail(c, http.StatusNotFound, "404 not found", nil)
		})

		// 405
		r.NoMethod(func(c *gin.Context) {
			resp.Fail(c, http.StatusMethodNotAllowed, "404 not found", nil)
		})
	}

	/*
		后台管理路由接口
	*/
	auth := r.Group("api/v1")
	back.BackRouter(auth)

	/*
		前端展示页面接口
	*/
	router := r.Group("api/v1")
	front.FrontRuter(router)

	return r
}
