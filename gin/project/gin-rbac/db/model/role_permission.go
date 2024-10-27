package model

import (
	"time"

	"gorm.io/gorm"
)

// RolePermissionModel 角色权限关联表模型
type RolePermissionModel struct {
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt
	RoleID       uint `gorm:"primary_key;not null;comment:角色id"`
	PermissionID uint `gorm:"primary_key;not null;comment:权限id"`
}

func (RolePermissionModel) TableName() string {
	return "r_role_permission"
}
