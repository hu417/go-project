package login

import (
	"gin-api-demo/api/user"

	"github.com/gin-gonic/gin"
)

// 注册
func register(r *gin.Engine) {
	r.GET("/register", user.Register)
}

// 登陆
func logint(r *gin.Engine) {
	r.GET("/logint", user.Login)
}

// 注销
func logout(r *gin.Engine) {
	r.GET("/logout", user.Logout)
}
