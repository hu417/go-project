package po

import (
	"time"
)

//角色
type Role struct {
	Id uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool `json:"is_deleted"`
	Name string `json:"name"`
	Code string `json:"code"`
	Description string `json:"description"`
}

func (*Role) TableName() string {
	return "rbac_role"
}