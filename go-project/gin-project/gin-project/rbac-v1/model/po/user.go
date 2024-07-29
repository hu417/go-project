package po

import (
	"time"
)

//用户
type User struct {
	Id uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool `json:"is_deleted"`
	Name string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Mail string `json:"mail"`
	Phone string `json:"phone"`
}

func (*User) TableName() string {
	return "rbac_user"
}