package model

import "gorm.io/gorm"

// PermissionModel 权限表模型
type PermissionModel struct {
	gorm.Model
	Name        string `gorm:"type:varchar(40);not null;unique;comment:权限名，唯一"`
	ApiPath     string `gorm:"type:varchar(255);not null;comment:API路径"`
	Method      string `gorm:"type:varchar(10);not null;comment:HTTP方法（GET, POST, PUT, DELETE等）"`
	Description string `gorm:"type:varchar(255);comment:权限描述"`
	ApiGroup    string `gorm:"type:varchar(40);comment:API分组"`
}

func (PermissionModel) TableName() string {
	return "s_permission"
}
