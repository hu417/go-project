// router.go
package v1

import (
	"fmt"
	"net/http"

	"gin-validator/internal/request"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Pong",
	})
}

func Signup(c *gin.Context) {

	var user request.Signup

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  request.GetErrorMsg(&user, err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Registered success, welcome user: %s", user.Username),
	})
}

func Login(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
