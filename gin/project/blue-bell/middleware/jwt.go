package middleware

import (
	"blue-bell/global"
	"blue-bell/pkg/token"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 放行登录 //
		// 创建一个跳过JWT认证的路由组
		skipPaths := []string{
			"/ping",
			"/api/v1/login",
			"/api/v1/signup",
		}
		for _, path := range skipPaths {
			if ctx.Request.URL.Path == path {
				ctx.Next()
				return
			}
		}

		// 获取token
		//从请求头中获取token: Authorization = "Bearer xxxxxx"
		tokenStr := ctx.Request.Header.Get("Authorization")
		//用户不存在
		if tokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 0, "msg": "token不能为空"})
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
		claims, err := token.Newjwt().ParseJwt(global.Conf.Jwt.Secret, tokenSlice[1])
		if err != nil && claims == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "token解析失败"})
			ctx.Abort() //阻止执行
			return
		}

		// 判断token是否过期
		if time.Now().Unix() > claims.ExpiresAt.Unix()+global.Conf.Jwt.RefreshGracePeriod {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "token过期",
			})
			ctx.Abort()
			return
			// 刷新token,生成新的token
			// lock := lock.Lock("refresh_token_lock", global.Conf.Jwt.JwtBlacklistGracePeriod)
			// if lock.Get() {

			// 	tokenData, _ := jwttoken.Newjwt().GenerateJwt(claims.User_Id, claims.Username, claims.Role, global.Conf.Jwt.Secret, global.Conf.Jwt.JwtTtl)
			// 	ctx.Header("x-access-token", tokenData)
			// 	ctx.Header("x-expires-time", string(time.Now().Add(time.Second*time.Duration(global.Conf.Jwt.JwtTtl)).Unix()))
			// 	_ = jwttoken.Newjwt().JoinBlackList(claims)
			// }
		}

		// 序列化claims
		claimsJSON, err := json.Marshal(claims)
		if err != nil {
			panic(fmt.Errorf("claims json marshal error: %v", err))

		}

		// 将claims信息保存到上下文
		ctx.Set("user_id", claims.User_Id)
		ctx.Set("username", claims.Username)
		ctx.Set("role", claims.Role)

		ctx.Set("token_claims", string(claimsJSON))
		ctx.Next()
	}
}
