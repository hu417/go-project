package router

import (
	"net/http"
	"sync"

	v1 "swag-demo/api/v1"
	"swag-demo/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CorsMiddle())

	// 404
	r.NoRoute(func(ctx *gin.Context) { // 这里只是演示，不要在生产环境中直接返回HTML代码
		ctx.String(http.StatusNotFound, "<h1>404 Page Not Found</h1>")
	})

	// 405
	r.NoMethod(func(ctx *gin.Context) {
		ctx.String(http.StatusMethodNotAllowed, "method not allowed")
	})

	// 注册swagger静态文件路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	/*
			// swagger路由;添加接口认证
		    authorized := r.Group("/swagger", gin.BasicAuth(gin.Accounts{
		        "admin": "666666",
		    }))
		    authorized.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	*/

	// 注册路由
	api := r.Group("/api/v1")
	{
		api.GET("/ping", v1.Ping)
		api.POST("/user/:id", v1.CreateByUser)
		api.GET("/test", v1.Test)
	}

	// 获取路由信息
	r.GET("/routes", func(ctx *gin.Context) {
		// 初始化互斥锁
		var mu sync.Mutex
		list := make(map[string][]string)
		routers := r.Routes()
		for _, v := range routers {
			// 确保 v 不为 nil
			if v.Method != "" && v.Path != "" {
				mu.Lock()
				// 使用互斥锁保证线程安全
				list[v.Method] = append(list[v.Method], v.Path)
				mu.Unlock()
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
			"routes":  list,
			"size":    len(list),
		})
	})

	return r
}
