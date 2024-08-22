package user

import "github.com/gin-gonic/gin"

func Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "logout",
	})
}
