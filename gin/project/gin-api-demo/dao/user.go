package dao

import (
	"gin-api-demo/global"
	"gin-api-demo/model"
)

type UserDao struct{}

func NewUserDao() *UserDao {
	return &UserDao{}
}

// 判断用户是否存在
func (*UserDao) IsNameExist(user *model.User) bool {
	var count int64
	global.DB.Model(user).Where("name = ?", user.Name).Count(&count)
	return count > 0
}

// 判断手机号是否存在
func (*UserDao) IsMobileExist(user *model.User) bool {
	var count int64
	global.DB.Model(user).Where("phone = ?", user.Mobile).Count(&count)
	return count > 0
}

// 用户列表
func (*UserDao) List(user *model.User,page,pagelimit int) (users []model.User,count int64,err error) {
	if user.Name == "" {
		user.Name = "%"
	}
	err = global.DB.Where("name like ?",user.Name).Count(&count).Limit(pagelimit).Offset((page -1 ) * pagelimit).Find(&users).Error
	return 
}