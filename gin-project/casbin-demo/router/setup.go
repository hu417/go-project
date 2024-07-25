package router

import (
	"casbin-demo/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine{
	// debug, release, test
	gin.SetMode("debug")
	r := gin.New()
	// 设置信任网络 []string
	// nil 为不计算，避免性能消耗，上线应当设置
	_ = r.SetTrustedProxies(nil)

	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// r.Static("/static", "./web/front/dist/static")
	// r.Static("/admin", "./web/admin/dist")
	// r.StaticFile("/favicon.ico", "/web/front/dist/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "front", nil)
	})

	r.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "admin", nil)
	})

	return r
}
