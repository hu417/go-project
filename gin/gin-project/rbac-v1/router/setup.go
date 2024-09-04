package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)




type router struct {}

func NewRouter() *router {
    return &router{}
}

func (*router) InitApiRouter(r *gin.Engine) {
    // 1. 注册路由
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })
}