package middleware

import (
	"casbin-demo/global"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 请求的path
		p := c.Request.URL.Path
		// 请求的方法
		m := c.Request.Method

		role := "superAdmin"
		//role:="user"
		//role:="guest"

		// 检查用户权限
		isPass, err := global.Enforcer.Enforce(role, p, m)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if isPass {
			c.Next()
		} else {
			c.JSON(401, gin.H{
				"error": err.Error(),
				"msg":   "无访问权限",
			})
			c.Abort()
			return
		}
	}
}
