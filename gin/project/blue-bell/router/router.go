package router

import (
	"net/http"
	"time"

	v1 "blue-bell/controller/v1"
	"blue-bell/global"
	"blue-bell/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(env string, t int) *gin.Engine {
	switch env {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "dev", "local":
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
		middleware.DefaultLimit(time.Second*2, 1), // 每两秒填充1个令牌，请求拿到桶中令牌才能获取响应，如果拿不到就获取不到响应
		middleware.TimeoutMiddleware(time.Second*time.Duration(t)),
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
		// 加载首页html文件
		r.LoadHTMLFiles("static/index.html")
		r.StaticFile("/", "./static/index.html")
		// 加载ico
		r.StaticFile("/favicon.ico", "./gin/project/blue-bell/static/logo.png")
		// 加载静态文件
		r.Static("/static", "./static")
	}

	// 后端API路由组
	api := r.Group("/api/v1")
	{
		api.POST("/signup", v1.SignupHandler)
		api.POST("/login", v1.LoginHandler)

		auth := api.Group("").Use(middleware.Jwt())
		{
			auth.GET("/index", func(c *gin.Context) {
				c.String(http.StatusOK, global.Conf.System.Version)
			})
			// 创建社区分类
			auth.POST("/community", v1.CommunityHandler)
			// 获取社区分类
			auth.GET("/community", v1.CommunityListHandler)
			// 通过id查询社区分类详情
			auth.GET("/community/:id", v1.CommunityDetailByIDHandler)
			// 用户新增帖子功能接口
			auth.POST("/post", v1.CreatePostHandler)
			// 实现分页获取帖子
			auth.GET("/post", v1.PostListHandler)
			// 通过帖子id查询帖子的详情，包括用户该帖子的作者信息以及社区帖子分类信息
			auth.GET("/post/:id", v1.PostDetailByIDHandler)

			// // 实现根据前端传来的参数，按时间排序返回或者按分数排序返回+返回帖子的投票数
			// api.GET("/posts_order", v1.PostListOrderHandler)
			// // 给帖子投票
			// api.POST("/vote", v1.PostVoteHandler)
		}

	}

	return r
}
