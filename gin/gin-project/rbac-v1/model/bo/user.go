package bo

import (
	"rbac-v1/model/po"
	"time"
)

type User struct {
	Id uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool `json:"is_deleted"`
	Name string `json:"name"`
	Username string `json:"username"`
	Mail string `json:"mail"`
	Phone string `json:"phone"`
	Roles []*po.Role `json:"roles"`
	Powers []*po.Power `json:"powers"`
}

type UserCreate struct {
	*po.User
	Roles []*po.Role `json:"roles"`
	Powers []*po.Power `json:"powers"`
}