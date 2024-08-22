package login

import "github.com/gin-gonic/gin"

func Router(e *gin.Engine) {
register(e)  // 注册
	logint(e) // 登陆
	logout(e) // 登出
}