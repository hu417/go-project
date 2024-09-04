package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 身份认证中间件
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//通过自定义的中间件，设置的值，在后续处理只要调用了这个中间件的都可以拿到这里参数set的值
		c.Set("usersesion", "userid-1") //用于全局变量

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
