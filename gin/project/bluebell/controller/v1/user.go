package v1

import (
	"errors"

	"bluebell/controller/request"
	"bluebell/controller/response"
	"bluebell/dao/mysql/sql"
	"bluebell/logic"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(request.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ErrorCode(c, response.CodeInvalidParam)
			return
		}
		response.ErrorWithCodeMsg(c, response.CodeInvalidParam, request.RemoveTopStruct(errs.Translate(request.Trans)))
		return
	}
	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, sql.ErrorUserExist) {
			response.ErrorCode(c, response.CodeUserExist)
			return
		}
		response.ErrorCode(c, response.CodeServerBusy)
		return
	}
	// 3. 返回响应
	response.Success(c, nil)
}

// LoginHandler 登录
func LoginHandler(c *gin.Context) {
	// 1.获取请求参数及参数校验
	p := new(request.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ErrorCode(c, response.CodeInvalidParam)
			return
		}
		response.ErrorWithCodeMsg(c, response.CodeInvalidParam, request.RemoveTopStruct(errs.Translate(request.Trans)))
		return
	}
	// 2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, sql.ErrorUserNotExist) {
			response.ErrorCode(c, response.CodeUserNotExist)
			return
		}
		response.ErrorCode(c, response.CodeInvalidPassword)
		return
	}

	// 3.返回响应
	response.Success(c, gin.H{
		"user_id":   user.UserID, // id值大于1<<53-1  int64类型的最大值是1<<63-1
		"user_name": user.Username,
		"token":     user.Token,
	})
}
