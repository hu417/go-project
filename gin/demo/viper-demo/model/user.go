package model

import (
	"viper-demo/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          int64     `gorm:"column:id;primary_key;auto_increment;" `
	UserId      string    `gorm:"column:user_id;type:varchar(100);not null"`
	UserName    string    `gorm:"column:username;type:varchar(20);not null" `
	Password    string    `gorm:"column:password;type:varchar(500);not null"`
	Email       string    `gorm:"column:email"`
	PhoneNumber string    `gorm:"column:phone_number"`
	Role        string    `gorm:"column:role;type:varchar(20);not null"`
	Century     string    `gorm:"century"`
	CreatedAt   time.Time `gorm:"autoCreateTime"` // 自动创建时间戳
	UpdatedAt   time.Time `gorm:"autoUpdateTime"` // 自动更新时间戳
}

func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	u.UserId = uuid.New().String()
	// 密码加密
	u.Password = utils.HashPw(u.Password)
	return nil
}

// 更新密码
func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.Password = utils.HashPw(u.Password)
	return nil
}
