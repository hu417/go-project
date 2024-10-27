package router

import (
	"net/http"

	v1 "gin-api-demo/api/v1"
	"gin-api-demo/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(env string) *gin.Engine {

	switch env {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "dev":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	// 全局中间件
	r.Use(
		middleware.GinLogger(),
		middleware.GinRecovery(),
		middleware.Cors(),
		middleware.DefaultLimit(),
		middleware.Retry(10*1024*1024),
	)

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
		r.StaticFile("/", "./static/dist/index.html")
		r.Static("/assets", "./static/dist/assets")
		r.StaticFile("/favicon.ico", "./gin-project/gin-api-demo/static/dist/logo.png")
	}

	// 注册 api 分组路由
	api := r.Group("/api")
	{
		api.GET("/ping", middleware.DefaultLimit(), func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
		api.POST("/register", v1.Register)
		api.POST("/login", v1.Login)

		auth := api.Group("").Use(middleware.Jwt())
		{
			auth.GET("/user/info", v1.UserInfo)
			auth.POST("user/logout", v1.Logout)
			auth.POST("/user/file", v1.Upload)
			auth.POST("/user/files", v1.Uploads)
		}
	}

	//获取所有路由
	r.GET("routers", func(c *gin.Context) {
		type path struct {
			Method string
			Path   string
		}
		var list []path
		routers := r.Routes()
		for _, v := range routers {
			list = append(list, path{Method: v.Method, Path: v.Path})
		}
		c.JSON(http.StatusOK, gin.H{
			"routes": list,
		})
	})

	return r
}
