package dao

import (
	"gin-api-v1/global"
	"gin-api-v1/model"
)

type UserDao struct{}

func NewUserDao() *UserDao {
	return &UserDao{}
}

// 判断邮箱是否存在
func (*UserDao) IsEmailExist(user *model.User) bool {
	var count int64
	global.DB.Model(user).Where("email = ?", user.Email).Count(&count)
	return count > 0
}

// 判断手机号是否存在
func (*UserDao) IsMobileExist(user *model.User) bool {
	var count int64
	global.DB.Model(user).Where("phone = ?", user.Mobile).Count(&count)
	return count > 0
}
