package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件

func Cors() gin.HandlerFunc {
	// 允许所有源
	return cors.New(
		cors.Config{
			//AllowAllOrigins:  true,
			AllowOrigins: []string{"*"}, // 等同于允许所有域名 #AllowAllOrigins:  true
			// 允许所有方法
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, 
			// 允许所有header
			AllowHeaders: []string{"*", "Authorization"}, 
			// 允许所有header
			ExposeHeaders: []string{"Content-Length", "text/plain", "Authorization", "Content-Type"}, 
			// 允许携带cookie
			AllowCredentials: true, 
			// 缓存时间
			MaxAge: 12 * time.Hour,
		},
	)

}
