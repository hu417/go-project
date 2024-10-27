package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRoleModel 用户角色关联表模型
type UserRoleModel struct {
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	UserID    uint `gorm:"primaryKey;autoIncrement:false;not null;comment:用户id"`
	RoleID    uint `gorm:"primaryKey;autoIncrement:false;not null;comment:角色id"`
}

func (UserRoleModel) TableName() string {
	return "r_user_role"
}
