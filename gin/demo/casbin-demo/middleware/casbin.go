package middleware

import (
	"casbin-demo/global"

	"github.com/gin-gonic/gin"
)

func CasbinMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		user := c.GetHeader("user")
		if user == "" {
			c.JSON(401, gin.H{
				"error": "权限验证失败",
				"msg":   "用户未登录",
			})
			c.Abort()
			return
		}
		if user == "superAdmin" {
			c.Next()
			return
		}

		// 请求的path
		p := c.Request.URL.Path
		// 请求的方法
		m := c.Request.Method

		role := user
		//role:="user"
		//role:="guest"

		// 检查用户权限
		isPass, err := global.Enforcer.Enforce(role, p, m)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if isPass {
			c.Next()
		} else {
			c.JSON(401, gin.H{
				"error": "权限验证失败",
				"msg":   "无访问权限",
			})
			c.Abort()
			return
		}
	}
}
