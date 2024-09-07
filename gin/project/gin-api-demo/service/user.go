package service

import (
	"errors"

	"gin-api-demo/api/req"
	"gin-api-demo/dao"
	"gin-api-demo/global"
	"gin-api-demo/model"
	"gin-api-demo/utils"
)

type userSvc struct {
}

func NewUserSvc() *userSvc {
	return &userSvc{}
}

// 判断用户是否存在
func (u *userSvc) IsUserExist(name string) bool {
	user := model.User{
		Name: name,
	}
	return dao.NewUserDao().IsNameExist(&user)
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

// 获取用户列表
func (u *userSvc) GetUserList(page, limit int, name string) (interface{}, error) {
	users,count,err := dao.NewUserDao().List(&model.User{
		Name: name,
	 },page,limit)
	
	 data := struct {
		Count int64  `json:"count"`
		User  []model.User `json:"users"`
		Page  int `json:"page"`
		Limit int `json:"limit"`
	 }{
		Count: count,
		User:  users,
		Page:  page,
		Limit: limit,
	 }

	 return data,err

}