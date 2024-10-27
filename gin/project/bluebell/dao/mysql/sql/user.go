package sql

import (
	"bluebell/global"
	"bluebell/models"

	"gorm.io/gorm"
)

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (user *models.User, err error) {
	result := global.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return user, nil
		}
		return user, result.Error
	}

	if result.RowsAffected > 0 {
		return user, ErrorUserExist
	}
	return
}

// InsertUser 想数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {

	return global.DB.Create(&user).Error
}

// GetUserById 根据id获取用户信息
func GetUserById(uid string) (user *models.User, err error) {

	err = global.DB.Select("user_id, username").Where("user_id = ?", uid).First(&user).Error

	return
}
