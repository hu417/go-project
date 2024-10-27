package models

// 自增ID主键
type BaseID struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

// 创建、更新时间和软删除时间
type Timestamps struct {
	CreatedAt int64 `json:"created_at" gorm:"column:create_at;comment:'创建时间'"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:update_at;comment:'更新时间'"`
	DeletedAt int64 `json:"deleted_at" gorm:"index;column:deleted_at;comment:'删除时间'"`
}
