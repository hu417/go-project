package db

import (
	"go-admin/global"
	"go-admin/model"

	"gorm.io/gorm"
)

// GetRoleList 获取角色列表
func GetRoleList(keyword string) *gorm.DB {
	tx := global.DB.Model(&model.SysRole{}).
	Select("id, name, is_admin, sort, created_at, updated_at")
	if keyword != "" {
		tx.Where("name LIKE ?", "%"+keyword+"%")
	}
	tx.Order("sort ASC")
	return tx
}

// GetRoleDetail 根据ID获取角色信息
func GetRoleDetail(id uint) (data *model.SysRole,err error) {

	if err := global.DB.Model(data).Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}
	return data,nil
}
