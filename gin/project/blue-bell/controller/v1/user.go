package v1

import (
	"fmt"
	"net/http"

	"blue-bell/controller/e"
	"blue-bell/controller/req"
	"blue-bell/controller/res"
	"blue-bell/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignupHandler 注册
func SignupHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	p := new(req.UserSignUp) // 得到一个注册请求参数结构体指针
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		res.ResponseError(c, http.StatusBadRequest, e.CodeInvalidParam, req.GetErrorMsg(p, err))
		return
	}
	// 2.业务处理 logic层进行业务处理
	if err := logic.Signup(p); err != nil {
		res.ResponseError(c, http.StatusInternalServerError, e.CodeServerBusy, err)
		return
	}
	// 3.返回响应
	res.ResponseSuccess(c, http.StatusOK, e.CodeSuccess, "用户注册成功")

}

// LoginHandler 登录
func LoginHandler(c *gin.Context) {
	// 1.参数校验
	p := new(req.UserLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		res.ResponseError(c, http.StatusBadRequest, e.CodeInvalidParam, req.GetErrorMsg(&p, err))
		return
	}
	// 2.业务处理
	token, err := logic.Login(p)
	fmt.Printf("token:%v\n", token)
	if err != nil {
		// 存入token相关的错误到日志中
		zap.L().Error("logic.Login another err", zap.Error(err))
		res.ResponseError(c, http.StatusInternalServerError, e.CodeServerBusy, err)
		return
	}
	// 3.返回响应
	res.ResponseSuccess(c, http.StatusOK, e.CodeSuccess, gin.H{"token": token})

}
