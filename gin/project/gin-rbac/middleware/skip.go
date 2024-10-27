package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// SkipSpecialRoutes 用于检查当前请求是否应该跳过某些特定的路由验证。
//
// 描述:
// 函数检查当前请求的路径是否包含预定义的一系列关键字，这些关键字代表了不需要进行额外验证的路由。
// 如果请求路径中包含了任何一个关键字，函数将返回true，指示该请求应当被跳过验证。
// 否则，返回false，表示请求需要经过正常的验证过程。

func SkipSpecialRoutes(c *gin.Context) bool {
	// 定义跳过的路由
	skipKeywords := []string{"public", "login", "register", "health", "swagger", "users"}
	for _, keyword := range skipKeywords {
		if strings.Contains(c.Request.URL.Path, keyword) {
			return true
		}
	}
	return false
}
