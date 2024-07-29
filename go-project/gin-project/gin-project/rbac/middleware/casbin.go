package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CasbinMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取username
		ctxUser := c.GetString("username")
		if ctxUser == "admin" {
			c.Next()
		} else {
			// TODO 从缓存中获取用户相关的信息，例如：role、dept、menu
			// 从数据库中获取用户角色信息sub
			usersInfo, err := system.NewUserInterface().GetUserFromUserName(ctxUser)
			if err != nil {
				fmt.Println("从数据库中获取用户角色信息sub失败:", err)
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "从数据库中获取用户角色信息sub失败" + err.Error()})
				c.Abort()
				return
			}
			sub := usersInfo.RoleId
			//获取请求路径
			obj := strings.Split(c.Request.RequestURI, "?")[0]
			// 获取请求方法
			act := c.Request.Method

			// 检查权限
			ok, err := global.CasbinEnforcer.Enforce(strconv.Itoa(int(sub)), obj, act)
			if err != nil || !ok {
				fmt.Println("权限验证失败：", err, ok)
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限验证失败"})
				
				c.Abort()
				// c.AbortWithStatus(http.StatusForbidden)

				return
			} else {
				c.Next()
			}
		}

	}
}
