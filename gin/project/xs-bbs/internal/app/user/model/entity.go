package model

import (
	"xs-bbs/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户结构体
type User struct {
	gorm.Model
	UserID   string `gorm:"not null;index:idx_user_id;"` // 用户ID
	Username string `gorm:"not null;size:32;unique;"`    // 用户名
	Email    string `gorm:"not null;size:128;unique;"`   // 邮箱
	Nickname string `gorm:"not null;size:16;"`           // 昵称
	Password string `gorm:"not null;size:512"`           // 密码
}

// 表名
func (u *User) TableName() string {
	return "user"
}

// 创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UserID = uuid.New().String()
	u.Password = utils.BcryptMake(u.Password)
	// u.Timestamps.CreatedAt = time.Now().Unix()
	return nil
}
