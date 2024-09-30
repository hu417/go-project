package db

import (
	"go-admin/global"
	"go-admin/model"

	"gorm.io/gorm"
)

// GetMenusList 获取所有菜单列表数据
func GetMenusList() *gorm.DB {
	tx := global.DB.Model(&model.SysMenu{}).
		Select("id,parent_id,name,web_icon,sort,path,level,component_name,created_at,updated_at").
		Order("sort ASC")
	return tx
}
