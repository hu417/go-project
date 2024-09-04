package po

import (
	"time"
)

//行为
type Operation struct {
	Id uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool `json:"is_deleted"`  //0未删除 1已删除
	Path string `json:"path"`
	Type int `json:"type"` //1页面 2API
	Method string `json:"method"`
	Description string `json:"description"`
}

func (*Operation) TableName() string {
	return "rbac_operation"
}