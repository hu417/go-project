// router.go
package v1

import (
	"fmt"
	"net/http"

	"gin-validator/global"
	"gin-validator/internal/request"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Pong",
	})
}

func Signup(c *gin.Context) {

	var user request.UserParams

	if err := c.ShouldBind(&user); err != nil {

		// 获取validator.ValidationErrors类型的error
		// errs, ok := err.(validator.ValidationErrors)
		// if !ok {
		// 	// 非validator.ValidationErrors类型错误直接返回
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"code": 400,
		// 		"msg":  err.Error(),
		// 	})
		// 	return
		// }
		// // validator.ValidationErrors类型错误则进行翻译
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"code": 400,
		// 	"msg":  val.RemoveTopStruct(errs.Translate(bootstrap.Trans)),
		// })
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
		sliceErr := []string{}

		for _, e := range errs {
			//翻译错误
			sliceErr = append(sliceErr, e.Translate(global.Trans))
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("%#v", sliceErr),
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("user: %s, login success.", login.Username),
	})
}
