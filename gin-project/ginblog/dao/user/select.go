package user

import (
	"ginblog/model"
)

func (dao *UserDao) CheckUser(user *model.User) bool {
	dao.DB.Select("id").Where("username = ?", user.Username).First(&user)

	return user.ID > 0
}
