package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gin-rbac/common/response"
	"gin-rbac/db/keys"
	"gin-rbac/global"
	"gin-rbac/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func JWTAuthMiddleware(router *gin.Engine, db *gorm.DB, rds *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查是否需要跳过 JWT 验证
		if SkipSpecialRoutes(ctx) {
			// 继续处理请求
			ctx.Next()
			return
		}

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
		claims, err := utils.ParseJwt(global.Config.JWT.Secret, tokenSlice[1])
		if err != nil && claims == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "token解析失败"})
			ctx.Abort() //阻止执行
			return
		}

		// 判断token是否过期
		if time.Now().Unix() > claims.ExpiresAt.Unix()+600 {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "token过期",
			})
			ctx.Abort()
			return
			// 刷新token,生成新的token
		}

		// 判断token是否在黑名单

		// 序列化claims
		claimsJSON, err := json.Marshal(claims)
		if err != nil {
			panic(fmt.Errorf("claims json marshal error: %v", err))

		}

		// 将claims信息保存到上下文
		ctx.Set("userID", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("role", claims.Role)

		ctx.Set("token_claims", string(claimsJSON))

		// 如果是超级管理员，直接放行
		isSuperAdminValue := claims.Role
		if isSuperAdminValue == 1 {
			ctx.Next()
			return
		}

		currentPath := ctx.Request.URL.Path // 获取当前请求的路径
		method := ctx.Request.Method        // 获取当前请求的方法
		// 匹配路由路径，并获取匹配的路由
		path := utils.MatchPath(currentPath, router)
		// 获取用户权限列表
		permissionList, err := keys.NewPermissionList(db, rds, ctx).GetPermissionListByUserID(claims.UserID)
		if err != nil {
			response.Forbidden(ctx, "Error getting permission list")
			ctx.Abort()
			return
		}
		// 检查用户是否有权限访问当前请求的资源
		if !utils.CheckPermissionExists(permissionList, path, method) {
			response.Forbidden(ctx, "Permission denied")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
