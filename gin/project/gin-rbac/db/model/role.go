package model

import "gorm.io/gorm"

// RoleModel 角色表模型
type RoleModel struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);not null;unique;comment:角色名，唯一"`
	Description string `gorm:"type:varchar(255);comment:角色描述"`
}

func (RoleModel) TableName() string {
	return "s_role"
}
