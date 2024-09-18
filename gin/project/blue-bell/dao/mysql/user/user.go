package user

import (
	"errors"

	"blue-bell/global"
	"blue-bell/model"
)

// CheckUserExist 判断用户是否存在
func CheckUserExist(user *model.User) (*model.User,error) {
	res := global.DB.Where("username = ?", user.UserName).First(&user)
	if res.Error != nil {
		return nil,res.Error
	}
	if res.RowsAffected == 0 {
		return nil,errors.New("用户不存在")
	}
	return user,nil
}

// InsertUser 向数据库中插入一条用户数据
func Insert(user *model.User) error {

	return global.DB.Create(user).Error
}
