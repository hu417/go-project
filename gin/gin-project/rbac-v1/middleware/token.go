package middleware

import (
	"github.com/gin-gonic/gin"
	"rbac-v1/common"
	"rbac-v1/common/constants"
	"time"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token
		token := c.GetHeader(constants.TOKEN_HEADER_KEY)
		if len(token) == 0 {
			common.ResponseTokenInvalid(c, "invalid token")
			c.Abort()
			return
		}
		//校验token
		claims, err := common.ParseToken(token)
		if err != nil {
			common.ResponseTokenInvalid(c, err.Error())
			c.Abort()
			return
		}
		if time.Now().Local().UnixMilli() > claims.ExpiresAt {
			common.ResponseTokenInvalid(c, "token expired")
			c.Abort()
			return
		}
		//放行
		c.Set(constants.CTX_USER_ID, claims.UserId)
		c.Set(constants.CTX_USERNAME, claims.Username)
		c.Next()
	}
}