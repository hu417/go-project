package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsDefault() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 表示允许所有域的访问
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 表示是否允许发送 cookies
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 指定了在跨域请求中可以使用的 HTTP 头部列表
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 指定了允许的 HTTP 方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// 方式二
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许所有来源（包括子域和端口）。生产环境中应替换为具体的允许域名。
		AllowOrigins: []string{"*", "http://localhost:9999"},
		// AllowOriginFunc优先级大于AllowOrigins
		// AllowOriginFunc: func(origin string) bool {
		// 	fmt.Println("origin:", origin)
		// 	// 允许跨域访问名单
		// 	var allowOriginsList = []string{"http://localhost:5173", "http://localhost:9999"}
		// 	//如果allowOriginsList中包含origin,允许访问
		// 	for _, v := range allowOriginsList {
		// 		// 如果访问域名是名单内的则放行
		// 		if strings.Contains(v, origin) {
		// 			return true
		// 		}
		// 	}
		// 	// 匹配不到不放行
		// 	return false
		// },
		// 允许的HTTP方法列表，如 GET、POST、PUT等。默认为["*"]（全部允许）
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		// 允许的HTTP头部列表。默认为["*"]（全部允许）
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization","x-access-token","x-refresh-token","file","files"},
		// 是否允许浏览器发送Cookie。默认为false
		AllowCredentials: true,
		// 预检请求（OPTIONS）的缓存时间（秒）。默认为5分钟
		MaxAge:           12 * time.Hour,
		// 允许浏览器端能够获取相应的header值
		ExposeHeaders:    []string{"Content-Length"},
	})
}
