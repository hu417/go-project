package models

import (
	"database/sql"
	"time"

	"bluebell/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id       int64          `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	UserID   string         `gorm:"column:user_id;NOT NULL"`
	Username string         `gorm:"column:username;NOT NULL"`
	Password string         `gorm:"column:password;NOT NULL"`
	Email    sql.NullString `gorm:"column:email"`
	Gender   int8           `gorm:"column:gender;default:0;NOT NULL"`
	Token    string         `gorm:"column:token"`
	Timestamps
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UserID = uuid.New().String()
	u.Password = utils.BcryptMake(u.Password)
	u.Timestamps.CreatedAt = time.Now().Unix()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.Timestamps.UpdatedAt = time.Now().Unix()
	return nil
}
