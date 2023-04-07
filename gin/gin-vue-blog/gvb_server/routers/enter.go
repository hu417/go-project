package routers

import (
	"gvb_server/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// 初始化路由
	gin.SetMode(global.Config.System.Env)
	r := gin.Default()

	// 定义接口
	r.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r

}
