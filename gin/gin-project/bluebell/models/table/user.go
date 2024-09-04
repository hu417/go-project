package table

import (
	"bluebell/pkg/utils"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	Id         int64        `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	UserId     int64        `gorm:"column:user_id;type:bigint(20);NOT NULL" json:"user_id"`
	Username   string       `gorm:"column:username;type:varchar(64);NOT NULL" json:"username"`
	Password   string       `gorm:"column:password;type:varchar(64);NOT NULL" json:"password"`
	Email      string       `gorm:"column:email;type:varchar(64)" json:"email"`
	Gender     int          `gorm:"column:gender;type:tinyint(4);default:0;NOT NULL" json:"gender"`
	CreateTime sql.NullTime `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime sql.NullTime `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP" json:"update_time"`
}

func (m *User) TableName() string {
	return "user"
}

// 随机id
func (t *User) BeforeCreate(tx *gorm.DB) error {
	t.UserId = int64(utils.GetUuidInt())
	return nil
}

// 秘密加密
func (t *User) HashPassword(oldPassword string) error {
	str, err := utils.GeneratePassword(oldPassword)
	if err != nil {
		return fmt.Errorf("hash password fail => %w", err)
	}
	t.Password = str
	return nil
}
