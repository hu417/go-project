package v1

import (
	"encoding/json"
	"errors"
	"gin-api-demo/api/req"
	"gin-api-demo/global"
	"gin-api-demo/service"
	"gin-api-demo/utils/jwttoken"
	"gin-api-demo/utils/res"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var user req.UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		res.FailByError(c, 400, BizUser, "参数错误", &res.CustomError{
			Code: BizUserRegister,
			Msg:  req.GetErrorMsg(&user, err),
		})
		return
	}
	// 用户注册
	use, err := service.NewUserSvc().Register(user)
	if err != nil {
		res.FailByError(c, 400, BizUser, "参数错误", &res.CustomError{
			Code: BizUserRegister,
			Msg:  err.Error(),
		})
		return
	}

	// 返回数据
	res.Success(c, 200, BizUser, "注册成功", use)
}

// Login 用户登录
func Login(c *gin.Context) {
	var form req.UserLogin
	// 参数绑定
	if err := c.ShouldBindJSON(&form); err != nil {
		res.FailByError(c, 400, BizUser, "参数错误", &res.CustomError{
			Code: BizUserLogin,
			Msg:  req.GetErrorMsg(&form, err),
		})
		return
	}
	// 用户校验
	user, err := service.NewUserSvc().Login(form)
	if err != nil {
		res.FailByError(c, 400, BizUser, "参数错误", &res.CustomError{
			Code: BizUserLogin,
			Msg:  err.Error(),
		})
		return
	}
	// 生成token
	token, err := jwttoken.Newjwt().GenerateJwt(user.UserId, user.Name, "root", global.Conf.Jwt.Secret, global.Conf.Jwt.JwtTtl)
	if err != nil {
		res.FailByError(c, 500, BizUser, "token生成失败", &res.CustomError{
			Code: BizUserLogin,
			Msg:  err.Error(),
		})
		return
	}
	// 返回数据
	data := struct {
		UserId   string `json:"user_id"`
		UserName string `json:"user_name"`
		Token    string `json:"token"`
		TokenExp int64  `json:"token_exp"`
	}{
		UserId:   user.UserId,
		UserName: user.Name,
		Token:    token,
		TokenExp: global.Conf.Jwt.JwtTtl / 60,
	}
	res.Success(c, 200, BizUser, "登录成功", data)
}

// UserInfo 获取用户信息
func UserInfo(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	if userId == "" {
		res.Fail(ctx, 401, BizUser, "获取失败", nil)
		return
	}
	user, err := service.NewUserSvc().GetUserInfo(userId)
	if err != nil {
		res.Fail(ctx, 500, BizUser, "获取失败", nil)
		return
	}
	res.Success(ctx, 200, BizUser, "获取成功", user)
}


// 用户退出
func Logout(ctx *gin.Context) {
	// 删除token
	tokenStr := ctx.GetString("token_claims")
	if tokenStr == ""  {
		res.Fail(ctx, 401, BizUser, "退出失败", errors.New("token为空"))
		return
	}
	// 加入黑名单
	claims := &jwttoken.MyCustomClaims{}
	if err := json.Unmarshal([]byte(tokenStr), claims);err != nil {
		res.Fail(ctx, 500, BizUser, "退出失败", err.Error())
		return
	}
	if err := jwttoken.Newjwt().JoinBlackList(claims);err != nil {
		res.Fail(ctx, 500, BizUser, "退出失败", err.Error())
		return
	}
	res.Success(ctx, 200, BizUser, "退出成功", nil)
}
func UserUpdate(ctx *gin.Context) {}
func UserDelete(ctx *gin.Context) {}
func UserList(ctx *gin.Context) {}
