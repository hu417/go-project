package middleware

import (
	"github.com/gin-gonic/gin"
	"rbac-v1/common"
	"rbac-v1/common/constants"
	"rbac-v1/model/vo"
	"rbac-v1/service/casbin"
)

//刷新
func RbacRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := casbin.NewCasbin().Refresh()
		if err != nil {
			common.ResponseRbacInvalid(c, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

//权限校验
func RbacAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString(constants.CTX_USERNAME)
		if len(username) == 0 {
			common.ResponseTokenInvalid(c, "invalid token")
			c.Abort()
			return
		}
		err := casbin.NewCasbin().Auth(&vo.CasbinAuthRequest{
			Username: username,
			Path:     c.Request.URL.Path,
			Method:   c.Request.Method,
		})
		if err != nil {
			common.ResponseRbacInvalid(c, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}