package model

import (
	"viper-demo/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId   string `gorm:"primary_key;auto_increment" json:"uid" validate:"required" label:"用户ID"`
	Username string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500);not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	u.UserId = uuid.New().String()
	// 密码加密
	u.Password = utils.HashPw(u.Password)
	return nil
}
