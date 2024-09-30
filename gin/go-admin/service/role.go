package service

import (
	"go-admin/api/request"
	"go-admin/dao/db"
	"go-admin/global"
	"go-admin/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRoleList 角色列表
func GetRoleList(c *gin.Context, in *request.GetRoleListRequest) (interface{}, error) {

	var (
		cnt  int64
		list = make([]*request.GetRoleListReply, 0)
	)
	err := db.GetRoleList(in.Keyword).Count(&cnt).Offset((in.Page - 1) * in.Size).Limit(in.Size).Find(&list).Error

	data := struct {
		List  []*request.GetRoleListReply `json:"list"`
		Count int64                       `json:"count"`
	}{
		List:  list,
		Count: cnt,
	}

	return data, err
}

// CheckRoleByName 判断角色名称是否存在
func CheckRoleByName(c *gin.Context, name string) (cnt int64, err error) {

	err = global.DB.Model(&model.SysRole{}).Where("name = ?", name).Count(&cnt).Error

	return cnt, nil
}

// CheckRoleByIdAndName 判断指定角色id和角色名称是否存在
func CheckRoleByIdAndName(c *gin.Context, id uint, name string) (cnt int64, err error) {

	err = global.DB.Model(&model.SysRole{}).Where("id != ? AND name = ?", id, name).Count(&cnt).Error
	return cnt, err
}

// AddRole 新增角色
func AddRole(c *gin.Context, in *request.AddRoleRequest) error {

	// 1. 给角色授权的菜单
	rms := make([]*model.RoleMenu, len(in.MenuId))
	for i, _ := range rms {
		rms[i] = &model.RoleMenu{
			MenuId: in.MenuId[i],
		}
	}

	// 2. 组件角色数据
	rb := &model.SysRole{
		Name:    in.Name,
		IsAdmin: in.IsAdmin,
		Sort:    in.Sort,
		Remarks: in.Remarks,
	}
	// 3. 新增角色数据
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 角色
		if err := tx.Create(rb).Error; err != nil {
			return err
		}
		// 授权菜单
		for _, v := range rms {
			v.RoleId = rb.ID
		}
		if len(rms) > 0 {
			if err := tx.Create(rms).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// PatchRoleAdmin 更改管理员身份
func PatchRoleAdmin(c *gin.Context, id string, isAdmin string) error {

	// 更改管理员身份
	return global.DB.Model(new(model.SysRole)).Where("id = ?", id).Update("is_admin", isAdmin).Error

}

// GetRoleDetail 根据ID获取角色详情
func GetRoleDetail(c *gin.Context, uId int) (*request.GetRoleDetailReply, error) {

	// 1、获取角色基本信息
	sysRole, err := db.GetRoleDetail(uint(uId))
	if err != nil {

		return nil, err
	}
	// 角色详情
	data := &request.GetRoleDetailReply{
		ID: sysRole.ID,
		AddRoleRequest: request.AddRoleRequest{
			Name:    sysRole.Name,
			IsAdmin: sysRole.IsAdmin,
			Sort:    sysRole.Sort,
			Remarks: sysRole.Remarks,
		},
	}

	// 2、获取授权的菜单
	menuIds, err := db.GetRoleMenuId(sysRole.ID, sysRole.IsAdmin == 1)
	if err != nil {
		return nil, err
	}
	data.MenuId = menuIds
	return data, nil
}

// UpdateRole 修改角色信息
func UpdateRole(c *gin.Context, in *request.UpdateRoleRequest) error {

	// 修改数据
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 更新角色信息
		err := global.DB.Model(new(model.SysRole)).Where("id = ?", in.ID).Updates(map[string]any{
			"name":     in.Name,
			"is_admin": in.IsAdmin,
			"sort":     in.Sort,
			"remarks":  in.Remarks,
		}).Error
		if err != nil {
			return err
		}
		// 删除授权的菜单老数据(使用Unscoped进行硬删除)
		err = tx.Where("role_id = ?", in.ID).Unscoped().Delete(new(model.RoleMenu)).Error
		if err != nil {
			return err
		}
		// 增加新授权的菜单数据
		rms := make([]*model.RoleMenu, len(in.MenuId))
		for i, _ := range rms {
			rms[i] = &model.RoleMenu{
				RoleId: in.ID,
				MenuId: in.MenuId[i],
			}
		}
		if len(rms) > 0 {
			err = tx.Create(rms).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

}

// DeleteRole 根据ID删除角色
func DeleteRole(c *gin.Context, id string) error {

	// 删除角色
	return global.DB.Where("id = ? ", id).Delete(new(model.SysRole)).Error

}

// AllRole 获取所有角色
func AllRole(c *gin.Context) (list []*request.AllListReply, err error) {

	err = global.DB.Model(model.SysRole{}).Find(&list).Error
	return list, err
}
