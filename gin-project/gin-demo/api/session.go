package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetSession(c *gin.Context) {
	//  设置session
	session := sessions.Default(c)
	session.Set("key", "value")
	session.Save()
	c.JSON(200, gin.H{
		"code":    200,
		"message": "ok",
	})
}

func GetSession(c *gin.Context) {
	//  获取session
	session := sessions.Default(c)
	value := session.Get("mySession")
	c.JSON(200, gin.H{
		"code":    200,
		"message": "ok",
		"value":   value,
	})
}