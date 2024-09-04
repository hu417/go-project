package model

import (
	"time"

	"gin-api-demo/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserId   string `json:"user_id" gorm:"column:user_id;not null;index;comment:用户ID"`
	Name     string `json:"name" gorm:"column:name;not null;comment:用户名称"`
	Mobile   string `json:"mobile" gorm:"column:mobile;not null;index;comment:用户手机号"`
	Password string `json:"password" gorm:"column:password;not null;default:'';comment:用户密码"`
	Timestamps
}

func (User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UserId = uuid.New().String()
	u.Password = utils.BcryptMake(u.Password)
	u.Timestamps.CreatedAt = time.Now().Unix()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.Timestamps.UpdatedAt = time.Now().Unix()
	return nil
}

func (u *User) AfterFind(tx *gorm.DB) error {
	return nil
}

func (u *User) AfterDelete(tx *gorm.DB) error {
	return nil
}
