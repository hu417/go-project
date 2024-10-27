package router

import (
	"net/http"
	"time"

	v1 "bluebell/controller/v1"
	"bluebell/middlewares"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// SetupRouter 路由
func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	r.Use(
		middlewares.GinLogger(),
		middlewares.GinRecovery(true),
		middlewares.RateLimitMiddleware(2*time.Second, 1))

	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"msg": "405",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		// 注册
		api.POST("/signup", v1.SignUpHandler)
		// 登录
		api.POST("/login", v1.LoginHandler)

		// 根据时间或分数获取帖子列表
		api.GET("/posts2", v1.GetPostListHandler2)
		// 获取帖子列表
		api.GET("/posts", v1.GetPostListHandler)
		// 获取社区列表
		api.GET("/community", v1.CommunityHandler)
		// 获取指定社区
		api.GET("/community/:id", v1.CommunityDetailHandler)
		// 获取指定帖子
		api.GET("/post/:id", v1.GetPostDetailHandler)

		api.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件
		{
			// 创建帖子
			api.POST("/post", v1.CreatePostHandler)
			// 投票
			api.POST("/vote", v1.PostVoteController)
		}

	}

	pprof.Register(r) // 注册pprof相关路由
	return r
}
