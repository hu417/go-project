package api

import (
	"github.com/gin-gonic/gin"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

// 首页，任何人可访问，不登录也可访问
func (*HomeController) Home(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "home",
	})

}
