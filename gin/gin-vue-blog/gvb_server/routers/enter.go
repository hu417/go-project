package routers

import (
	// "net/http"

	// "gvb_server/api"
	"gvb_server/global"

	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	// *gin.Engine  // 路由
	*gin.RouterGroup // 路由组
}

func InitRouter() *gin.Engine {
	// 初始化路由
	gin.SetMode(global.Config.System.Env)
	r := gin.Default()

	//####################### 定义接口
	//# 方式一：直接返回
	// r.GET("", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "ok")
	// })

	// 方式二：响应函数封装
	// settingsApi := api.ApiGroupApp.SettingsApi
	// r.GET("/", settingsApi.SettingsInfoView)

	// 方式三：router group封装
	apiRouterGroup := r.Group("api")
	{
		routerGroupApp := RouterGroup{apiRouterGroup}
		routerGroupApp.SettingsGroup()
		routerGroupApp.ImagesGroup()
	}

	return r

}
