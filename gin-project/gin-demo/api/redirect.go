package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Redirect 重定向
func RedirectHttp(c *gin.Context) {
	// 使用Context调用Redirect()⽀持内部和外部的重定向
	// 重定向到外部
	c.Redirect(http.StatusMovedPermanently, "http:// www.baidu.com/")
	// 重定向到内部
	// c.Redirect(http.StatusMovedPermanently, "/内部接口/路径")
}

// 路由重定向
func RedirectRoute(c *gin.Context, r *gin.Engine) {
	// 1.设置重定向的url到Context中
	c.Request.URL.Path = "/test2"
	// 2.通过Router调用HandleContext(c)进行,重定向到/test2上
	r.HandleContext(c)
}

func TestRedirect(c *gin.Context) {

	c.String(http.StatusOK, "redirect success")

}
