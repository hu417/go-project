package v1

import (
	"fmt"
	"swag-demo/model"

	"github.com/gin-gonic/gin"
)

// @Tags			测试模块
// @Summary		ping api
// @Description	return hello world json format content
//
// @Accept			application/json
// @Produce		application/json
//
// @Param			name	query	string	true	"name"
// @Produce		json
// @Success		200	json	{string}	string	"hello world"
// @Router			/ping [get]
func Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": fmt.Sprintf("Hello World!%s", ctx.Query("name")),
	})
}

// @Tags			用户模块
// @Summary		用户创建
// @Description	用户功能相关接口
// @Accept			application/json,application/x-www-form-urlencoded,multipart/form-data,application/octet-stream
// @Produce		application/json
// @Param			Authorization	header		string			true	"Bearer 用户令牌"
// @Param			id				path		string			true	"用户id"
// @Param			object			body		model.User		true	"创建参数"
// @Success		200				{object}	RespuserList	"响应参数"
// @Failure		400				{object}	string			"请求错误"
// @Failure		500				{object}	string			"内部错误"
// @Router			/user/{id} [post]
// @Security		Bearer
// @Contact.name	接口联系人
// @Contact.url	联系人网址
// @Contact.email	联系人邮箱
func CreateByUser(ctx *gin.Context) {
	var user []*model.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"data":    user,
		"message": "sussecc",
	})
}

// @Tags			Test模块
// @Summary		测试接口
// @Description	| 项目 | 进展 | 人员 |
// @Description	| :-------- | --------:| :--: |
// @Description	| iPhone | 完成 | hu |
// @Description	| iPad | 未完成 | li |
// @Description	| iMac | 未开始 | wang |
// @Accept			application/json,application/x-www-form-urlencoded,multipart/form-data,application/octet-stream
// @Produce		application/json
// @Success		200	{object}	RespuserList	"响应参数"
// @Failure		400	{object}	string			"请求错误"
// @Failure		500	{object}	string			"内部错误"
// @Router			/test [get]
// @Security		BasicAuth
func Test(ctx *gin.Context) {
	var user model.User

	ctx.JSON(200, gin.H{
		"code":    200,
		"data":    user,
		"message": "sussecc",
	})
}
