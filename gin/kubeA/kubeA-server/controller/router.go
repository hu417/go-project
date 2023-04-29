package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 实例化router结构体，可使用该对象点出首字母大写的方法（包外调用）
var Router router

// 创建router结构体
type router struct{}

// 初始化路由规则，并创建测试api接口
func (r *router) InitApiRouter(router *gin.Engine) {

	// 设置路由接口
	router.GET("/api/v1", func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "test api success",
			"data": nil,
		})

	})

	// 默认路由
	router.NoRoute(func(c *gin.Context) {
		// 实现内部重定向
		c.String(http.StatusNotFound, "404")
	})

}
