package middleware

import (
	"net/http"
	"strings"
	"time"

	"ginblog/global"
	"ginblog/utils"

	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取token
		//从请求头中获取token: Authorization = "Bearer xxxxxx"
		tokenStr := ctx.Request.Header.Get("Authorization")
		//用户不存在
		if tokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 0, "msg": "token不存在"})
			ctx.Abort() //阻止执行
			return
		}
		//token格式错误
		tokenSlice := strings.SplitN(tokenStr, " ", 2)
		if len(tokenSlice) != 2 && tokenSlice[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 0, "msg": "token格式错误"})
			ctx.Abort() //阻止执行
			return
		}

		// 解析token
		claims, err := utils.Newjwt().ParseJwt(global.Conf.Jwt.SignKey, tokenSlice[1])
		if err != nil && claims == nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "token解析失败"})
			ctx.Abort() //阻止执行
			return
		}

		// 判断token是否过期
		if time.Now().Unix() > claims.ExpiresAt.Unix()+10 {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "token过期",
			})
			ctx.Abort()
			return
		}

		// 将claims信息保存到上下文
		ctx.Set("user_id", claims.User_Id)
		ctx.Set("username", claims.Username)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
