package model

import (
	"gorm.io/gorm"
)

// 自增ID主键
type BaseID struct {
	ID uint `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
}

// 创建、更新时间和软删除时间
type Timestamps struct {
	CreatedAt int64          `gorm:"column:created_at;index;" json:"created_at,omitempty"`
	UpdatedAt int64          `gorm:"column:updated_at;index;" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at" `
}
