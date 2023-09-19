// 注意go服务的启动类package要编写为main
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//1.创建router,实际返回的是一个Engine引擎,也可以说成初始化容器创建Engine引擎
	r := gin.Default() // 默认会包含些初始化中间件: Logger(), Recovery()
	//r:=gin.New() 了解Default与New的区别,空的Engine

	//2.注册全局中间件
	r.Use(MiddlewareFunc())

	//3.注册路由
	r.GET("/index", func(context *gin.Context) {
		context.Writer.Write([]byte("返回参数"))
	})

	//4.注册路由组,Group()方法会返回一个新生成的RouterGroup指针,用来区分不同的路由组与执行对应不同的中间件等
	groupRouter := r.Group("/groupPath")

	//5.路由组注册中间件
	groupRouter.Use(MiddlewareFunc())
	{
		//6.路由组注册路由
		groupRouter.GET("/hello", func(context *gin.Context) {
			context.Writer.Write([]byte("返回参数"))
		})
	}
 
 	//7.添加资源路径
	r.Static("/static", "dist/static")
	//前端接口
	//这样只要启动后端代码，访问根目录就直接访问到静态资源了
	r.StaticFile("/", "dist/index.html")
        
    //8.redirect重定向
    r.GET("/redirect", func(c *gin.Context) {
		//使用Context调用Redirect()⽀持内部和外部的重定向
		//重定向到外部
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
		//重定向到内部
		//c.Redirect(http.StatusMovedPermanently, "/内部接口/路径")
	})
        
    //9.路由重定向
    r.GET("/test", func(c *gin.Context) {
		//1.设置重定向的url到Context中
		c.Request.URL.Path = "/test2"
		//2.通过Router调用HandleContext(c)进行,重定向到/test2上
		r.HandleContext(c)
	})
	r.GET("/test2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
        
    //10.NoRoute处理
    r.NoRoute(func(c *gin.Context) { c.IndentedJSON(http.StatusNotFound,gin.H{ "msg":"not found", }) })
	
	// demo案例
	api := r.Group("/api")
	{	
		// 静态路由
		api.GET("/user",func(c *gin.Context){
			c.JSON(http.StatusOK, gin.H{
				"message":"ok",
			})
		})

		// 参数路由
		api.GET("/user/:name",func(c *gin.Context){
			name := c.Param("name")
			c.JSON(http.StatusOK,gin.H{
				"name": name,
			})
		})

		// 通配符路由
		api.GET("/view/*.html",func(c *gin.Context){
			path := c.Param(".html")
			c.JSON(http.StatusOK,gin.H{
				"path": path,
			})
		})

		// url查询参数
		api.GET("/user",func(c *gin.Context){
			id := c.Query("id")
			c.JSON(http.StatusOK,gin.H{
				"id": id,
			})
		})

		// body参数: form表单
		api.POST("/user",func(c *gin.Context){
			// 使用PostForm获取请求参数
			username := c.PostForm("username")
			password := c.PostForm("password")
			
			// 返回请求参数
			c.JSON(200, gin.H{
				"username": username,
				"password": password,
			})
		})

		// body参数: json数据
		api.POST("/user",func(c *gin.Context){
			type User struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			// 将JSON格式请求参数绑定到结构体上
			var user User
			if err := c.BindJSON(&user); err != nil {
				// 返回错误信息
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 根据req的content type 自动推断如何绑定,form/json/xml等格式
			// if err := c.ShouldBind(&user); err != nil {      
			// 	c.JSON(400, gin.H{"error": err.Error()})
			// 	return
			// }
	
			
			// 返回请求参数
			c.JSON(200, gin.H{
				"username": user.Username,
				"password": user.Password,
			})
		})


		// header参数
		// 使用请求头部参数并获取参数
		api.GET("/token", func(c *gin.Context) {
			// 使用Request获取请求头部参数
			username := c.Request.Header.Get("username")
			password := c.Request.Header.Get("password")
	
			// 返回请求参数
			c.JSON(200, gin.H{
				"username": username,
				"password": password,
			})
		})
	}
	
	//11.监听端口启动服务
	r.Run(":8081")
}

//中间件函数
func MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[INFO] %s", "中间件业务执行")
		c.Next()
	}
}

