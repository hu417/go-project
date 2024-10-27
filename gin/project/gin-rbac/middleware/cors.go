package middleware

import "github.com/gin-gonic/gin"

// CorsMiddleware 是一个Gin框架的中间件，用于处理跨域资源共享（CORS）。
//
// 描述:
// 中间件通过设置HTTP响应头来允许跨域请求。它设置了允许的所有源、请求方法、请求头以及是否允许携带Cookie。
// 特别地，对于预检请求（即方法为OPTIONS的请求），中间件将直接返回204状态码，表示请求可以通过预检，而无需进一步处理。
// 对于其他类型的请求，中间件将继续调用链中的下一个处理器。
//
// 返回:
// - gin.HandlerFunc: 一个可以被Gin框架使用的中间件函数
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源，这里为允许所有源，如果要设置允许的源，可以修改为具体的源
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 设置允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")

		// 设置允许的Header
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

		// 设置是否允许携带cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
