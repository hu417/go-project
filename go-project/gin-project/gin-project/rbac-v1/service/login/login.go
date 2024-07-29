package login

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"rbac-v1/common"
	"rbac-v1/model/vo"
	"rbac-v1/service"
)

type Login struct {
	*service.Service
}

func NewLogin() *Login {
	return &Login{Service: service.Srv()}
}

//登录
func (l *Login) LoginByPwd(ctx context.Context, param *vo.LoginRequest) (ret *vo.LoginResponse, err error) {
	//获取用户数据
	user, err := l.Dao().GetUserByUsername(ctx, param.Username)
	if err != nil {
		logger.Error("get user failed: ", err.Error())
		return nil, err
	}
	if user == nil {
		err = errors.New(fmt.Sprintf("user is not exist, username: %s", param.Username))
		logger.Error(err.Error())
		return nil, err
	}
	//校验密码
	if user.Password != param.Password {
		err = errors.New(fmt.Sprintf("invalid username or password, username: %s", param.Username))
		logger.Error(err.Error())
		return nil, err
	}
	//生成Token
	token, exp, err := common.GenerateToken(user)
	if err != nil {
		logger.Error("generate token failed: ", err.Error())
		return nil, err
	}

	return &vo.LoginResponse{
		Token:     token,
		ExpiresIn: exp,
	}, nil
}

//Token校验
func (l *Login) CheckToken(param *vo.CheckTokenRequest) (ret *vo.CheckTokenResponse, err error) {
	claims, err := common.ParseToken(param.Token)
	if err != nil {
		logger.Error("parse token failed: ", err.Error())
		return nil, err
	}

	return &vo.CheckTokenResponse{
		UserId:    claims.UserId,
		Name:      claims.Name,
		Username:  claims.Username,
		ExpiresIn: claims.ExpiresAt,
	}, nil
}