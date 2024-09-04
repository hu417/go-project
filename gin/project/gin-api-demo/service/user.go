package service

import (
	"errors"

	"gin-api-demo/api/req"
	"gin-api-demo/global"
	"gin-api-demo/model"
	"gin-api-demo/utils"
)

type userSvc struct {
}

func NewUserSvc() *userSvc {
	return &userSvc{}
}

// Register 注册
func (u *userSvc) Register(params req.UserRegister) (user model.User, err error) {
	result := global.DB.Where("mobile = ?", params.Mobile).Select("id").First(&model.User{})
	if result.RowsAffected != 0 {
		err = errors.New("手机号已存在")
		return user, err
	}
	// 创建用户
	user.Name = params.Name
	user.Password = params.Password
	user.Mobile = params.Mobile

	err = global.DB.Create(&user).Error
	// 密码加密
	user.Password = "**********"
	return user, err
}

// Login 登录
func (u *userSvc) Login(params req.UserLogin) (user model.User, err error) {
	result := global.DB.Where("mobile = ?", params.Mobile).First(&user)
	if result.RowsAffected == 0 {
		err = errors.New("手机号不存在")
		return user, err
	}
	if !utils.BcryptMakeCheck(params.Password, user.Password) {
		return user, errors.New("密码错误")
	}
	return user, nil
}

// 获取用户信息
func (u *userSvc) GetUserInfo(userId string) (user model.User, err error) {
	result := global.DB.Where("user_id = ?", userId).First(&user)
	if result.RowsAffected == 0 {
		err = errors.New("用户不存在")
		return 
	}
	user.Password = "**********"
	return 
}
