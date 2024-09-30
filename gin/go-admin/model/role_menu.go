package model

import "gorm.io/gorm"

// RoleMenu 角色和菜单关系结构体
type RoleMenu struct {
	gorm.Model
	RoleId uint `gorm:"column:role_id;type:int(11);" json:"role_id"` // 角色ID
	MenuId uint `gorm:"column:menu_id;type:int(11);" json:"menu_id"` // 菜单ID
}

// TableName 设置表名
func (table *RoleMenu) TableName() string {
	return "sys_role_menu"
}

