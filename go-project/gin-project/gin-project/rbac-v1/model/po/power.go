package po

import (
	"time"
)

//权限
type Power struct {
	Id uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool `json:"is_deleted"`
	Name string `json:"name"`
	Code string `json:"code"`
	Description string `json:"description"`
}

func (*Power) TableName() string {
	return "rbac_power"
}