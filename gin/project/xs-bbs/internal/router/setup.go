package router

import (
	"net/http"
	"time"
	"xs-bbs/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server http  server
type Server interface {
	// RegisterHTTPRouter register http router
	RegisterHTTPRouter(r *gin.Engine)
}

func NewHttpServer(servers ...Server) *gin.Engine {

	r := gin.New()
	// 全局中间件
	r.Use(middleware.GinLogger(),
		middleware.GinRecovery(),
		middleware.Cors(),
		middleware.DefaultLimit(time.Second*2, 1), // 每两秒填充1个令牌，请求拿到桶中令牌才能获取响应，如果拿不到就获取不到响应
		middleware.TimeoutMiddleware(time.Second*time.Duration(3)),
		middleware.Jwt(),
	)

	// ping
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	// 404
	r.NoRoute(func(ctx *gin.Context) { // 这里只是演示，不要在生产环境中直接返回HTML代码
		ctx.String(http.StatusNotFound, "<h1>404 Page Not Found</h1>")
	})

	// 405
	r.NoMethod(func(ctx *gin.Context) {
		ctx.String(http.StatusMethodNotAllowed, "method not allowed")
	})

	// 前端项目静态资源
	{
		// r.StaticFile("/", "./static/dist/index.html")
		// r.Static("/assets", "./static/dist/assets")
		r.StaticFile("/favicon.ico", "./gin/project/blue-bell/static/dist/logo.png")
		// 其他静态资源
		// r.Static("/public", "./static")
		// r.Static("/storage", "./storage/app/public")
	}

	// 注册swagger静态文件路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	/*
			// swagger路由;添加接口认证
		    authorized := r.Group("/swagger", gin.BasicAuth(gin.Accounts{
		        "admin": "666666",
		    }))
		    authorized.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	*/

	// 后端API路由组
	// api := r.Group("/api/v1")
	{

		for _, s := range servers {
			s.RegisterHTTPRouter(r)
		}
	}
	return r

}
