// router.go
package v1

import (
	"fmt"
	"net/http"

	"gin-validator/internal/request"
	"gin-validator/internal/utils/val"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Pong",
	})
}

func AddUser(c *gin.Context) {

	var user request.UserParams

	if err := c.ShouldBind(&user); err != nil {
		fmt.Println(err.Error())
		// 获取validator.ValidationErrors类型的error
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  val.GetErrorMsg(user, errs),
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Registered success, welcome user: %s", user.Username),
	})
}

func Login(c *gin.Context) {
	var login request.LoginParams

	if err := c.ShouldBind(&login); err != nil {
		fmt.Println(err.Error())
		// 获取validator.ValidationErrors类型的error
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": val.GetErrorMsg(&login, errs),
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("user: %s, login success.", login.Username),
	})
}
