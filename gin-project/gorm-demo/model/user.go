package model

import (
	"gorm-demo/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          string `gorm:"column:id;primary_key;auto_increment;" `
	UserId      string `gorm:"column:user_id;type:varchar(100);not null"`
	UserName    string `gorm:"column:username;type:varchar(20);not null" `
	Password    string `gorm:"column:password;type:varchar(500);not null"`
	Email       string `gorm:"column:email"`
	PhoneNumber string `gorm:"column:phone_number"`
	Role        string `gorm:"column:role;type:varchar(20);not null"`
	Century     string `gorm:"century"`
	Date        string `gorm:"date"`
}

func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	u.Password = utils.HashPw(u.Password)
	return nil
}

func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.Password = utils.HashPw(u.Password)
	return nil
}
