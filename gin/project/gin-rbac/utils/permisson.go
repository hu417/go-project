package utils

import (
	"regexp"
	"strings"

	"gin-rbac/global"

	"github.com/gin-gonic/gin"
)

// MatchPath 用于查找与给定路径相匹配的路由模式。
/*
	描述:
	函数遍历由Gin引擎提供的所有路由规则，尝试找到一个与给定的实际请求路径相匹配的路由模式。
	它通过将路由模式中的动态部分（例如，`:id`）转换为正则表达式的匹配模式来实现这一点。
	转换后的正则表达式用于检查实际路径是否符合某个路由模式。
	如果找到匹配的路由模式，函数返回该模式；否则，返回空字符串。
*/
func MatchPath(path string, router *gin.Engine) string {
	// 遍历所有路由
	for _, route := range router.Routes() {
		pattern := route.Path
		// 将模式中的变量转换为正则表达式的匹配模式
		// regexp.QuoteMeta(pattern)：将模式中的所有特殊字符转义，这样它们会被当作普通字符处理。
		// regexp.MustCompile(:[a-zA-Z]+)：创建一个正则表达式，用于匹配以冒号开头后跟一个或多个字母的字符串。
		// ReplaceAllString(regexPattern, ([0-9]+))：将所有匹配的部分替换为正则表达式模式 ([0-9]+)，表示一个或多个数字。
		// re.MatchString(path)：使用编译好的正则表达式对象来匹配实际路径。
		regexPattern := "^" + regexp.QuoteMeta(pattern)
		regexPattern = regexp.MustCompile(`:[a-zA-Z]+`).ReplaceAllString(regexPattern, `([0-9]+)`)

		// 创建正则表达式对象
		re := regexp.MustCompile(regexPattern)

		// 使用正则表达式匹配实际路径
		if re.MatchString(path) {
			return route.Path
		}
	}
	return ""
}

// CheckPermissionExists 检查用户是否有指定API路径和方法的权限。
// 描述:
// 函数接收一个权限列表映射、一个API路径和一个HTTP方法作为参数。
// 它首先检查权限列表中是否存在给定的API路径。
// 如果存在，函数进一步检查该路径下是否有给定的HTTP方法。
// 如果API路径和方法都存在于权限列表中，函数返回true，表示用户有权限；否则返回false。
func CheckPermissionExists(permissionList map[string]map[string]bool, apiPath string, method string) bool {
	if perms, ok := permissionList[apiPath]; ok {
		return perms[method]
	}
	return false
}

// IsValidateApiPath 验证 API 路径的有效性
// 描述:
// 函数用于验证传入的 API 路径是否符合以下规则：
//   - 路径必须以 `/` 开头。
//   - 路径的最大长度不能超过 255 个字符。
//   - 路径只能包含字母、数字以及 `/`, `_`, `-`, `.`, `{`, `}`, 和 `*` 字符。
func IsValidateApiPath(path string) bool {
	// API 路径应该以 `/` 开头
	if !strings.HasPrefix(path, "/") {
		return false
	}

	// API 路径长度不应超过 255 个字符
	if len(path) > 255 {
		return false
	}

	// 使用正则表达式验证路径是否合法
	// 允许使用字母、数字、`/`, `_`, `-`, `.`, `{}`, `}`, `*`,
	match, _ := regexp.MatchString(`^[/a-zA-Z0-9_\-./\{\}\*:]+$`, path)
	return match
}

// IsValidateHTTPMethod 检查给定的 HTTP 方法是否有效
// 描述:
// 函数用于验证传入的 HTTP 方法是否为一组预定义的有效方法之一。
func IsValidateHTTPMethod(method string) bool {
	// 定义有效的 HTTP 方法列表
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	// 遍历有效方法列表
	for _, v := range validMethods {
		// 使用 strings.EqualFold 进行不区分大小写的比较
		if strings.EqualFold(method, v) {
			return true
		}
	}

	return false
}

// IsValidPathAndMethod 检查API路径和HTTP方法的有效性
// 描述:
// 该函数用于验证传入的API路径和HTTP方法是否均符合预定义的有效标准
func IsValidPathAndMethod(apiPath string, method string) bool {
	global.Log.Errorln("IsValidateApiPath(apiPath)", IsValidateApiPath(apiPath))
	if !IsValidateApiPath(apiPath) {
		global.Log.Warnln("Invalid path:", apiPath)
	}
	if !IsValidateHTTPMethod(method) {
		global.Log.Warnln("Invalid method:", method)
	}
	return IsValidateApiPath(apiPath) && IsValidateHTTPMethod(method)
}
