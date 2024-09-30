// 注意go服务的启动类package要编写为main
package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-demo/api"
	v1 "gin-demo/api/v1"
	"gin-demo/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {

	//  Gin 框架在运行的时候默认是debug模式 有： 开发：debug，生产：release，测试模式：test
	//  DebugMode indicates gin mode is debug.
	//  DebugMode = "debug"
	//  ReleaseMode indicates gin mode is release.
	//  ReleaseMode = "release"
	//  TestMode indicates gin mode is test.
	//  TestMode = "test"
	gin.SetMode(gin.ReleaseMode)

	//  禁用控制台颜色
	gin.DisableConsoleColor()

	//  创建记录日志的文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	//  如果需要将日志同时写入文件和控制台，请使用以下代码
	//  gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 1.创建router,实际返回的是一个Engine引擎,也可以说成初始化容器创建Engine引擎
	r := gin.Default() //  默认会包含些初始化中间件: Logger(), Recovery()
	// r:=gin.New() 了解Default与New的区别,空的Engine

	// 2.注册全局中间件
	r.Use(middleware.MiddlewareFunc(), middleware.AuthMiddleWare())

	// 3.注册路由
	r.GET("/index", func(context *gin.Context) {
		context.Writer.Write([]byte("返回参数"))
	})

	// 4.注册路由组,Group()方法会返回一个新生成的RouterGroup指针,用来区分不同的路由组与执行对应不同的中间件等
	groupRouter := r.Group("/groupPath")

	// 5.路由组注册中间件
	groupRouter.Use(middleware.MiddlewareFunc())
	{
		// 6.路由组注册路由
		groupRouter.GET("/hello", func(context *gin.Context) {
			context.Writer.Write([]byte("返回参数"))
		})
	}

	//  7.添加资源路径: 静态资源 //
	{
		//    /assets/images/1.jpg 这个url文件，存储在/var/www/assets/images/1.jpg
		r.Static("/assets", "/var/www/assets")

		// 前端接口:这样只要启动后端代码，访问根目录就直接访问到静态资源了
		r.StaticFile("/", "dist/index.html")

		//  为单个静态资源文件，绑定url
		//  这里的意思就是将/favicon.ico这个url，绑定到./resources/favicon.ico这个文件
		r.StaticFile("/favicon.ico", "./resources/favicon.ico")

		// 模版文件: 首先加载templates目录下面的所有模版文件，模版文件扩展名随意
		r.LoadHTMLGlob("templates/*")
		//  绑定一个url路由 /index
		r.GET("/index", func(c *gin.Context) {
			//  通过HTML函数返回html代码
			//  第二个参数是模版文件名字
			//  第三个参数是map类型，代表模版参数
			//  gin.H 是map[string]interface{}类型的别名
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Main website",
				"msg":   "这是Go后台传递来的数据",
			})
		})
	}

	// 8.重定向
	{
		// redirect重定向
		r.GET("/redirect", api.RedirectHttp)

		// 路由重定向
		r.GET("/test", func(c *gin.Context) {
			// 1.设置重定向的url到Context中
			c.Request.URL.Path = "/test2"
			// 2.通过Router调用HandleContext(c)进行,重定向到/test2上
			r.HandleContext(c)
		})

		r.GET("/test2", api.TestRedirect)
	}

	// 9.文件上传
	{
		// 文件上传(单个文件)
		//  给表单限制上传大小 (默认 32 MiB)
		r.MaxMultipartMemory = 8 << 20 //  8 MiB
		r.POST("/upload", api.UploadFile)

		// 文件上传(多个文件)
		//  获取MultipartForm
		r.POST("/uploads", api.UploadMultipleFile)
	}

	// 10.NoRoute处理
	r.NoRoute(func(c *gin.Context) {
		//  c.HTML(http.StatusNotFound, "404.html", nil) // 转到404.html页面，再返回空
		// http.StatusNotFound为404状态码
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"msg": "not found",
		})
	})

	// 11.NoMethod处理
	r.NoMethod(func(c *gin.Context) {
		//  c.HTML(http.StatusNotFound, "404.html", nil) // 转到404.html页面，再返回空
		// http.StatusNotFound为404状态码
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"msg": "not found",
		})
	})

	//  12.响应数据：String、JSON、XML、YAML、Protobuf
	{
		r.GET("/str", api.RespStr)
		r.GET("/json", api.RespJson)
		r.GET("/xml", api.RespXml)
		r.GET("/yaml", api.RespYaml)
		r.GET("/protobuf", api.RespProtobuf)
	}

	//  13.异步执行
	r.GET("/long_async", func(c *gin.Context) {
		//  需要搞一个副本
		copyContext := c.Copy()
		//  异步处理
		go func() {
			time.Sleep(3 * time.Second)
			log.Println("异步执行：" + copyContext.Request.URL.Path)
			//  注意不能在这执行重定向的任务，不然panic
		}()
	})

	//  14.会话控制
	{
		// cookie
		r.GET("/getCookie", api.GetCookie)

		// session
		// 下载相关包: go get -u "github.com/gin-contrib/sessions"
		// 注意该密钥不要泄露了
		store := cookie.NewStore([]byte("secret"))
		// 路由上加入session中间件
		r.Use(sessions.Sessions("mySession", store))
		r.GET("/setSession", api.SetSession)
		r.GET("/getSession", api.GetSession)
	}

	//  demo案例
	api := r.Group("/api")
	{
		//  静态路由
		api.GET("/user", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		})

		//  模拟登录
		api.POST("/loginIn", v1.UserLogin)

		//  尝试访问，添加身份认证中间件，如果已经登陆就可以执行
		api.GET("/sayHello", middleware.AuthMiddleWare(), func(c *gin.Context) {
			// 取出中间件的值
			usersesion := c.MustGet("usersesion").(string)
			log.Println("===============>", usersesion)
			if usersesion == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "请先登录",
				})
				return
			}
			c.String(200, "Hello World!")
		})

		//  参数路由
		api.GET("/user/:name", v1.GetMathByUrlForDefault)

		//  通配符路由,如: index.html
		api.GET("/view/*.html", v1.GetMathByUrl)

		//  url查询参数
		api.GET("/user", v1.GetUserForQueryById)

		//  body参数: form表单
		api.POST("/user", v1.CreateUserforForm)

		//  body参数: json数据
		api.POST("/user", v1.CreateUserforJson)

		//  header参数: 使用请求头部参数并获取参数
		api.GET("/token", v1.GetUserForHeader)

	}

	//=========== 优雅退出 ============//
	// 16.监听端口启动服务
	// curl "http://localhost:8000/sleep?duration=5s"
	r.GET("/sleep", func(c *gin.Context) {
		duration, err := time.ParseDuration(c.Query("duration"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		time.Sleep(duration)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	// 17.服务配置
	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 32M
		TLSConfig: nil,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener:
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		log.Println("Stopped serving new connections")
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// We received an SIGINT/SIGTERM signal, shut down.
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}
	log.Println("HTTP server graceful shutdown completed")
}
