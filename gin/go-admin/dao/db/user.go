package db

import (
	"go-admin/global"
	"go-admin/model"

	"gorm.io/gorm"
)

// GetUserByUsernamePassword 根据用户名和密码查询数据
func GetUserByUsernamePassword(username, password string) (data *model.SysUser, err error) {
	if err := global.DB.
		Where("username = ? AND password = ?", username, password).
		First(&data).Error; err != nil {
		return nil, err
	}
	return data, nil

}

// GetUserList 获取管理员列表
func GetUserList(keyword string) *gorm.DB {
	tx := global.DB.Model(&model.SysUser{}).Select("id,username,phone,avatar,remarks,created_at,updated_at")
	if keyword != "" {
		tx.Where("username LIKE ?", "%"+keyword+"%")
	}
	return tx
}

// GetUserDetail 根据ID获取管理员信息
func GetUserDetail(id uint) (data *model.SysUser, err error) {

	if err := global.DB.Model(&data).
		Where("id = ?", id).
		First(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
