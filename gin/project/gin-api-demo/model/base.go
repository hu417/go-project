package model

import (
	"gorm.io/gorm"
)

// 自增ID主键
type BaseID struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

// 创建、更新时间和软删除时间
type Timestamps struct {
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}


