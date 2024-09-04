package models

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

// GetRoleMenuId 获取指定角色的菜单
func GetRoleMenuId(roleId uint, isAdmin bool) ([]uint, error) {
	tx := new(gorm.DB)
	data := make([]uint, 0)
	if isAdmin {
		tx = DB.Model(new(SysMenu)).Select("id").Order("sort ASC")
	} else {
		tx = DB.Model(new(RoleMenu)).Select("sm.id").
			Joins("LEFT JOIN sys_menu sm ON sm.id = sys_role_menu.menu_id").
			Where("sys_role_menu.role_id = ?", roleId).Order("sm.sort ASC")
	}
	err := tx.Scan(&data).Error
	return data, err
}

// GetRoleMenus 获取指定角色的菜单列表
func GetRoleMenus(roleId uint, isAdmin bool) (*gorm.DB, error) {
	tx := new(gorm.DB)
	if isAdmin {
		tx = DB.Model(new(SysMenu)).Select("id, parent_id, component_name, name, web_icon, sort, path, level").Order("sort ASC")
	} else {
		roleBasic := new(SysRole)
		err := DB.Model(new(SysRole)).Select("id").Where("id = ?", roleId).Find(roleBasic).Error
		if err != nil {
			return nil, err
		}
		tx = DB.Model(new(RoleMenu)).Select("mb.id, mb.parent_id, mb.component_name, mb.name, mb.web_icon, mb.sort, mb.path, mb.level").
			Joins("LEFT JOIN sys_menu mb ON mb.id = sys_role_menu.menu_id").
			Where("sys_role_menu.role_id = ?", roleBasic.ID).Order("mb.sort ASC")
	}
	return tx, nil
}
