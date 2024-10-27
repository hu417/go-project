package model

import (
	"gin-rbac/utils"

	"gorm.io/gorm"
)

// UserModel 用户表模型
type UserModel struct {
	gorm.Model
	Username     string `gorm:"type:varchar(15);not null;unique;comment:用户名，唯一,长度限制为[1~15]"`
	Password     string `gorm:"type:varchar(255);not null;comment:密码哈希值，base64编码存储"`
	Email        string `gorm:"type:varchar(255);comment:邮箱地址，业务逻辑唯一"`
	PhoneNum     string `gorm:"type:varchar(11);comment:电话号码，业务逻辑唯一"`
	Sex          int8   `gorm:"type:tinyint;not null;default:0;comment:用户性别，默认0，0-未填写，1-男，2-女。"`
	Intro        string `gorm:"type:varchar(30);not null;default:这个人还没想好怎么介绍自己;comment:简介.默认值:这个人还没想好怎么介绍自己"`
	AvatarID     uint   `gorm:"not null;default:1;comment:用户头像id，默认值：1"`
	IsSuperAdmin uint   `gorm:"type:tinyint;not null;default:0;comment:默认0，0-普通用户，1-超级管理员。"`
}

func (UserModel) TableName() string {
	return "s_user"
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password, err = utils.HashPassword(u.Password)

	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) BeforeUpdate(tx *gorm.DB) (err error) {
	u.Password, err = utils.HashPassword(u.Password)

	if err != nil {
		return err
	}
	return nil
}
