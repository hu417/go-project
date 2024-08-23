package user

import (
	"ginblog/model"

	"gorm.io/gorm"
)

type UserDao struct {
	DB *gorm.DB
}

type UserDaoInterface interface {
	CheckUser(user *model.User)
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		DB: db,
	}
}
