package service

import (
	"errors"

	"go-admin/api/request"
	"go-admin/dao/db"
	"go-admin/global"
	"go-admin/model"

	"github.com/gin-gonic/gin"
)

// Menus 获取菜单列表
func Menus(c *gin.Context, role_id uint, is_admin bool) ([]*request.AllMenu, error) {

	allMenus := make([]*request.AllMenu, 0)

	// 根据角色获取所有菜单列表数据
	//tx := model.GetMenusList()
	tx, err := db.GetRoleMenus(role_id, is_admin)
	if err != nil {
		return nil, errors.New("数据异常")
	}
	err = tx.Find(&allMenus).Error
	if err != nil {
		return nil, errors.New("数据异常")
	}
	return allMenus, nil
}

// AddMenu 新增菜单
func AddMenu(c *gin.Context, in *request.AddMenuRequest) error {

	// 保存数据
	if err := global.DB.Create(&model.SysMenu{
		ParentId:      in.ParentId,
		Name:          in.Name,
		WebIcon:       in.WebIcon,
		Sort:          in.Sort,
		Path:          in.Path,
		Level:         in.Level,
		ComponentName: in.ComponentName,
	}).Error; err != nil {
		return errors.New("数据库异常")
	}

	return nil

}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context, in *request.UpdateMenuRequest) error {

	// 更新数据
	return global.DB.Model(new(model.SysMenu)).
		Where("id = ?", in.ID).
		Updates(map[string]interface{}{
			"parent_id":      in.ParentId,
			"name":           in.Name,
			"web_icon":       in.WebIcon,
			"sort":           in.Sort,
			"path":           in.Path,
			"level":          in.Level,
			"component_name": in.ComponentName,
		}).Error

}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context, id string) error {

	// 删除数据库中的数据
	return global.DB.Where("id = ?", id).Delete(new(model.SysMenu)).Error

}
