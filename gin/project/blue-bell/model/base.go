package model

import "gorm.io/gorm"

type BaseID struct {
	ID uint `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
}

type Timestamps struct {
	CreatedAt int64          `gorm:"column:created_at;index;" json:"created_at,omitempty"`
	UpdatedAt int64          `gorm:"column:updated_at;index;" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at" `
}
