package db

import (
	"go-admin/global"
	"go-admin/model"

	"gorm.io/gorm"
)

// GetRoleMenuId 获取指定角色的菜单
func GetRoleMenuId(roleId uint, isAdmin bool) ([]uint, error) {
	tx := new(gorm.DB)
	data := make([]uint, 0)
	if isAdmin {
		tx = global.DB.Model(&model.SysMenu{}).
			Select("id").
			Order("sort ASC")
	} else {
		tx = global.DB.Model(&model.RoleMenu{}).
			Select("sm.id").
			Joins("LEFT JOIN sys_menu sm ON sm.id = sys_role_menu.menu_id").
			Where("sys_role_menu.role_id = ?", roleId).
			Order("sm.sort ASC")
	}
	err := tx.Scan(&data).Error
	return data, err
}

// GetRoleMenus 获取指定角色的菜单列表
func GetRoleMenus(roleId uint, isAdmin bool) (*gorm.DB, error) {
	tx := new(gorm.DB)
	if isAdmin {
		tx = global.DB.Model(&model.SysMenu{}).Select("id, parent_id, component_name, name, web_icon, sort, path, level").Order("sort ASC")
	} else {
		roleBasic := new(model.SysRole)
		err := global.DB.Model(&model.SysRole{}).
			Select("id").
			Where("id = ?", roleId).
			Find(roleBasic).Error
		if err != nil {
			return nil, err
		}
		tx = global.DB.Model(&model.RoleMenu{}).
			Select("mb.id, mb.parent_id, mb.component_name, mb.name, mb.web_icon, mb.sort, mb.path, mb.level").
			Joins("LEFT JOIN sys_menu mb ON mb.id = sys_role_menu.menu_id").
			Where("sys_role_menu.role_id = ?", roleBasic.ID).
			Order("mb.sort ASC")
	}
	return tx, nil
}
