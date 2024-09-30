package controller

import (
	"xs-bbs/internal/app/user/model"
	"xs-bbs/internal/errs"
	"xs-bbs/internal/responce"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary 用户注册账号
// @Description 用户注册
// @Tags 用户接口
// @ID /user/signup
// @Accept  json
// @Produce  json
// @Param body body model.RegisterReq true "body"
// @Success 200 {object} responce.{data=model.UserDto} "success"
// @Router /user/signup [post]
func (u *UserController) Register(c *gin.Context) {
	var (
		err    error
		uParam model.RegisterReq
		uDto   *model.UserDto
	)
	if err :=c.ShouldBind(&uParam); err != nil {
		responce.ErrorWithMsg(c, errs.CodeError, err.Error())
		return
	}

	uDto, err = u.userService.Register(c.Request.Context(), &uParam)

	switch err {
	case nil:
		responce.Success(c, uDto)
	case errs.ErrEmailExist:
		responce.Error(c,errs.CodeEmailExist)
	case errs.ErrConvDataErr:
		responce.Error(c,errs.CodeConvDataErr)
	default:
		responce.Error(c,errs.CodeError)
	}
}

// Login godoc
// @Summary 登录
// @Description 登录
// @Tags 用户接口
// @ID /user/signin
// @Accept  json
// @Produce json
// @Param body body model.LoginReq true "body参数"
// @Success 200 {string} string "ok" "登陆成功"
// @Router /user/signin [post]
func (u *UserController) Login(c *gin.Context) {
	var (
		signParam model.LoginReq
	)
	if err :=c.ShouldBind(&signParam); err != nil {
		responce.ErrorWithMsg(c, errs.CodeError, err.Error())
		return
	}
	
	token, err := u.userService.Login(c.Request.Context(), &signParam)

	switch err {
	case nil:
		responce.Success(c, token)
	case errs.ErrUserNotExist:
		responce.Error(c,errs.CodeUserNotExist)
	default:
		responce.Error(c,errs.CodeWrongUserNameOrPassword)
	}
}

// Get godoc
// @Summary 根据id获取用户
// @Description 根据id获取用户
// @Tags 用户接口
// @ID /user/Get
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Success 200 {object} responce.{data=model.UserDto} "success"
// @Router /user/Get [get]
func (u *UserController) Get(c *gin.Context) {
	var (
		uDto   *model.UserDto
	)

	// 获取当前用户ID
	userId := c.GetString("user_id")
	if userId == "" {
		responce.Error(c, errs.CodeInvalidParams)
		return
	}
	

	uDto, err := u.userService.SelectByID(c.Request.Context(), userId)

	switch err {
	case nil:
		responce.Success(c, uDto)
	case errs.ErrUserNotExist:
		responce.Error(c, errs.CodeUserNotExist)
	case errs.ErrConvDataErr:
		responce.Error(c, errs.CodeConvDataErr)
	default:
		responce.Error(c, errs.CodeError)
	}
}

// Delete godoc
// @Summary 根据id删除用户
// @Description 根据id删除用户
// @Tags 用户接口
// @ID /user/delete
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Success 200 {object} responce.{data=string} "success"
// @Router /user/delete [get]
func (u *UserController) Delete(c *gin.Context) {

	// 获取当前用户ID
	userId := c.GetString("user_id")
	if userId == "" {
		responce.Error(c, errs.CodeInvalidParams)
		return
	}
	if !u.userService.Delete(c.Request.Context(), userId) {
		responce.Error(c, errs.CodeError)
		return
	}
	responce.Success(c, nil)
}
