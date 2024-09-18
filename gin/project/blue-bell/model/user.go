package model

import (
	"time"

	"blue-bell/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
构造user结构体与数据库中user表相对应，相互绑定
*/

type User struct {
	BaseID
	UserID   string `db:"user_id" json:"user_id" gorm:"column:user_id;index;comment:用户ID"`
	UserName string `db:"username" json:"username" gorm:"column:username;index:idx_username;unique;comment:用户名称"`
	PassWord string `db:"password" json:"password" gorm:"column:password;comment:用户密码"`
	Timestamps
}

// 表名
func (u *User) TableName() string {
	return "user"
}

// 创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UserID = uuid.New().String()
	u.PassWord = utils.BcryptMake(u.PassWord)
	u.Timestamps.CreatedAt = time.Now().Unix()
	return nil
}
