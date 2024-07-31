// 注意go服务的启动类package要编写为main
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {

	// Gin 框架在运行的时候默认是debug模式 有： 开发：debug，生产：release，测试模式：test
	// DebugMode indicates gin mode is debug.
	// DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	// ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	// TestMode = "test"
	gin.SetMode(gin.ReleaseMode)
	
	//1.创建router,实际返回的是一个Engine引擎,也可以说成初始化容器创建Engine引擎
	r := gin.Default() // 默认会包含些初始化中间件: Logger(), Recovery()
	//r:=gin.New() 了解Default与New的区别,空的Engine

	// 静态资源 //
	//   /assets/images/1.jpg 这个url文件，存储在/var/www/tizi365/assets/images/1.jpg
    	r.Static("/assets", "/var/www/tizi365/assets")

    	// 为单个静态资源文件，绑定url
    	// 这里的意思就是将/favicon.ico这个url，绑定到./resources/favicon.ico这个文件
    	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	//2.注册全局中间件
	r.Use(MiddlewareFunc(), AuthMiddleWare())
	
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

	//10.文件上传(单个文件)
	// 给表单限制上传大小 (默认 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", func(c *gin.Context) {
	    file, err := c.FormFile("file") // file是表单字段名字
	    if err != nil {
	        c.String(500, "上传文件出错")
	    }
	
	    // 上传到指定路径
	    c.SaveUploadedFile(file, "static/"+file.Filename)
	    c.String(http.StatusOK, "fileName:", file.Filename)
	})

	//11.文件上传(多个文件)
	// 获取MultipartForm
	r.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
		    c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
		}
		
		// 获取所有文件
		files := form.File["files"]
		for _, file := range files {
		    // 逐个存
		    fmt.Println(file.Filename)
		    // 上传到指定路径
	            c.SaveUploadedFile(file, "static/"+file.Filename)
		}
		c.String(200, fmt.Sprintf("upload ok %d files", len(files)))
	})
        //12.NoRoute处理
        r.NoRoute(func(c *gin.Context) {
		// c.HTML(http.StatusNotFound, "404.html", nil) //转到404.html页面，再返回空
		//http.StatusNotFound为404状态码
		c.IndentedJSON(http.StatusNotFound,gin.H{
			"msg":"not found", 
		}) 
	})
	
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
			//指定默认值
			// DefaultQuery()若参数不村则，返回默认值，Query()若不存在，返回空串
                        // name := c.DefaultQuery("name", "normal")
			name := c.Param("name")
			c.JSON(http.StatusOK,gin.H{
				"name": name,
			})
		})

		// 通配符路由,如: index.html
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
			// 设置默认值
                        types := c.DefaultPostForm("type", "post")

			// 使用PostForm获取请求参数
			username := c.PostForm("username")
			password := c.PostForm("password")
			
			// 返回请求参数
			c.JSON(200, gin.H{
				"types": types,
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
			// c.GetRawData() 表示raw类型的数据,可以理解为特定格式的字符串，如json字符串
			if err := c.BindJSON(&user); err != nil {
				// 返回错误信息
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
	
			// 根据req的content type 自动推断如何绑定,form/json/xml等格式
			// 如果发送的不是json格式，那么输出：  "error": "invalid character '-' in numeric literal"
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

		// 响应数据：String、JSON、XML、YAML
		// 1.String/JSON
		// r.GET("/someJSON", func(c *gin.Context) {
	        //     c.String(200,"欢迎访问%s, 你是%s", "tizi360.com!","最靓的仔！")
		//     c.JSON(200, gin.H{
		//         "message": "Json",
		//         "status":  200,
		//     })
		// })
		// // 2.XML
		// r.GET("/someXML", func(c *gin.Context) {
		//     c.XML(200, gin.H{"message": "abc"})
		// })
		// // 3.YAML
		// r.GET("/someYAML", func(c *gin.Context) {
		//     c.YAML(200, gin.H{"name": "zhangsan"})
		// })
		// // 4.protobuf
		// r.GET("/someProtoBuf", func(c *gin.Context) {
		//     reps := []int64{1, 2}
		//     data := &protoexample.Test{
		//         Reps:  reps,
		//     }
		//     c.ProtoBuf(200, data)
		// })
		

	}
	// 13.异步执行
	r.GET("/long_async", func(c *gin.Context) {
	    // 需要搞一个副本
	    copyContext := c.Copy()
	    // 异步处理
	    go func() {
	        time.Sleep(3 * time.Second)
	        log.Println("异步执行：" + copyContext.Request.URL.Path)
	        // 注意不能在这执行重定向的任务，不然panic
	    }()
	})

	// 14.会话控制
	// cookie
	r.GET("/getCookie", func(c *gin.Context) {
	    // 获取客户端是否携带cookie
	    cookie, err := c.Cookie("key_cookie")
	    if err != nil {
	        cookie = "cookie"
	        c.SetCookie("key_cookie", "value_cookie", // 参数1、2： key & value
	                    60,          // 参数3： 生存时间(秒);如果是-1,则表示删除
	                    "/",         // 参数4： 所在目录
	                    "localhost", // 参数5： 域名
	                    false,       // 参数6： 安全相关 - 是否智能通过https访问
	                    true,        // 参数7： 安全相关 - 是否允许别人通过js获取自己的cookie
	                   )
	    }
	    fmt.Printf("cookie的值是： %s\n", cookie)
	})
	       
	// session 
	// 下载相关包: go get -u "github.com/gin-contrib/sessions"
	// 注意该密钥不要泄露了
	store := cookie.NewStore([]byte("secret"))
	//路由上加入session中间件
	r.Use(sessions.Sessions("mySession", store))

	r.GET("/setSession", func(c *gin.Context) {
		// 设置session
		session := sessions.Default(c)
		session.Set("key", "value")
		session.Save()
	})

	r.GET("/getSession", func(c *gin.Context) {
		// 获取session
		session := sessions.Default(c)
		v := session.Get("key")
		fmt.Println(v)
	})

	// 模拟登录
	r.GET("/loginIn", func(c *gin.Context) {
		// 获取客户端是否携带cookie
		_, err := c.Cookie("key_cookie")
		if err != nil {
			c.SetCookie("key_cookie", "value_cookie", // 参数1、2： key & value
				10,          // 参数3： 生存时间（秒）
				"/",         // 参数4： 所在目录
				"localhost", // 参数5： 域名
				false,       // 参数6： 安全相关 - 是否智能通过https访问
				true,        // 参数7： 安全相关 - 是否允许别人通过js获取自己的cookie
			)
			c.String(200, "login success")
			return
		}
		c.String(200, "already login")
	})

	// 尝试访问，添加身份认证中间件，如果已经登陆就可以执行
	r.GET("/sayHello", AuthMiddleWare(), func(c *gin.Context) {
		//取出中间件的值
		usersesion := c.MustGet("usersesion").(string)
		log.Println("===============>",usersesion)
		c.String(200, "Hello World！")
	})

	//15.模版文件
	// 首先加载templates目录下面的所有模版文件，模版文件扩展名随意
    	r.LoadHTMLGlob("templates/*")
    	// 绑定一个url路由 /index
    	r.GET("/index", func(c *gin.Context) {
            // 通过HTML函数返回html代码
            // 第二个参数是模版文件名字
            // 第三个参数是map类型，代表模版参数
            // gin.H 是map[string]interface{}类型的别名
            c.HTML(http.StatusOK, "index.html", gin.H{
                    "title": "Main website",
		    "msg":"这是Go后台传递来的数据",
            })
    	})
	       
	//16.监听端口启动服务
	r.Run(":8081")
}

//中间件函数
func MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[INFO] %s", "中间件业务执行")
		c.Next()
	}
}

// 计算时间
func CalcTimeMiddleWare() gin.HandlerFunc {
	fmt.Println(1)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		// 统计时间
		since := time.Since(start)
		fmt.Println("程序用时：", since)
	}
}

// 身份认证中间件
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//通过自定义的中间件，设置的值，在后续处理只要调用了这个中间件的都可以拿到这里参数set的值
		c.Set("usersesion","userid-1")         //用于全局变量
		
		// 获取客户端cookie并校验
		if cookie, err := c.Cookie("key_cookie"); err == nil {
			if cookie == "value_cookie" { // 满足该条件则通过
				return
			}
		}
		// 返回错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
		// 若验证不通过，不再调用后续的函数处理
		c.Abort()
	}
}


