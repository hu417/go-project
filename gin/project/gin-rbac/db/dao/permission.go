package dao

import (
	"context"
	"fmt"

	"gin-rbac/db/model"
	"gin-rbac/dtos"
	"gin-rbac/global"

	"gorm.io/gorm"
)

// PermissionDao 权限数据访问接口
type PermissionDao interface {
	// 获取权限列表
	GetPermissionList(getPermissionListReqDTO *dtos.GetPermissionListReqDTO) ([]*model.PermissionModel, int64, error)
	// 创建权限
	CreatePermission(permission *model.PermissionModel) error
	// 根据id获取权限详情
	GetPermissionByID(id uint) (*model.PermissionModel, error)
	// 根据id批量获取权限详情
	GetPermissionByIDs(ids []uint) ([]uint, []*model.PermissionModel, error)
	// 根据id获取已删除的权限详情
	GetDelPermissionByID(id uint) (*model.PermissionModel, error)
	// 根据name获取权限详情
	GetPermissionByName(name string) (*model.PermissionModel, error)
	// 根据path和method获取权限详情
	GetPermissionByPathMethod(path, method string) (*model.PermissionModel, error)
	// 更新权限
	UpdatePermission(permission *model.PermissionModel) error
	// 删除权限
	DeletePermission(id uint) error
	// 恢复权限
	RecoverPermission(id uint) error
}

// permissionDAO 权限数据访问实现
type permissionDAO struct {
	db  *gorm.DB
	ctx context.Context
}

// NewPermissionDAO 创建权限数据访问实现
func NewPermissionDAO(db *gorm.DB, ctx context.Context) PermissionDao {
	return &permissionDAO{
		db:  db,
		ctx: ctx,
	}
}

// GetPermissionByIDs 批量获取权限详情(根据id)
func (p *permissionDAO) GetPermissionByIDs(ids []uint) ([]uint, []*model.PermissionModel, error) {
	var permissions []*model.PermissionModel
	existingIDsMap := make(map[uint]bool)
	nonExistingIDs := make([]uint, 0)
	if err := p.db.WithContext(p.ctx).Model(&model.PermissionModel{}).Where("id IN (?)", ids).Find(&permissions).Error; err != nil {
		return nil, nil, err
	}
	for _, permission := range permissions {
		existingIDsMap[permission.ID] = true
	}
	for _, id := range ids {
		if _, ok := existingIDsMap[id]; !ok {
			nonExistingIDs = append(nonExistingIDs, id)
		}
	}
	return nonExistingIDs, permissions, nil

}

// GetPermissionList 获取权限列表
func (p *permissionDAO) GetPermissionList(getPermissionListReqDTO *dtos.GetPermissionListReqDTO) ([]*model.PermissionModel, int64, error) {
	var permissions []*model.PermissionModel
	var count int64

	// 构造查询
	query := p.db.WithContext(p.ctx).Model(&model.PermissionModel{})

	// 添加排序
	query = query.Order("api_group DESC")
	if getPermissionListReqDTO.SortBy != "" && getPermissionListReqDTO.Order != "" {
		query = query.Order(fmt.Sprintf("%s %s", getPermissionListReqDTO.SortBy, getPermissionListReqDTO.Order))
	}
	if getPermissionListReqDTO.Name != "" {
		query = query.Where("name LIKE ?", "%"+getPermissionListReqDTO.Name+"%")
	}
	if getPermissionListReqDTO.Method != "" {
		query = query.Where("method LIKE ?", getPermissionListReqDTO.Method)
	}
	if getPermissionListReqDTO.ApiPath != "" {
		query = query.Where("api_path LIKE ?", "%"+getPermissionListReqDTO.ApiPath+"%")
	}
	if getPermissionListReqDTO.ApiGroup != "" {
		query = query.Where("api_group LIKE ?", "%"+getPermissionListReqDTO.ApiGroup+"%")
	}
	if getPermissionListReqDTO.Description != "" {
		query = query.Where("description LIKE ?", "%"+getPermissionListReqDTO.Description+"%")
	}

	// 添加过滤条件
	for _, filter := range getPermissionListReqDTO.Filters {
		switch filter.Op {
		case "eq":
			query = query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
		case "neq":
			query = query.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
		case "gt":
			query = query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
		case "gte":
			query = query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
		case "lt":
			query = query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
		case "lte":
			query = query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
		case "contains":
			query = query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+filter.Value+"%")
		case "not_contains":
			query = query.Where(fmt.Sprintf("%s NOT LIKE ?", filter.Field), "%"+filter.Value+"%")
		}
	}

	// 计算偏移量
	offset := (getPermissionListReqDTO.Page - 1) * getPermissionListReqDTO.Size

	// 执行查询
	err := query.Count(&count).Offset(offset).Limit(getPermissionListReqDTO.Size).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, count, nil
}

// CreatePermission 创建权限
func (p *permissionDAO) CreatePermission(permission *model.PermissionModel) error {
	return Transaction(p.db, func(tx *gorm.DB) error {
		return tx.Create(permission).Error
	})
}

// GetPermissionByID 根据id获取权限详情
func (p *permissionDAO) GetPermissionByID(id uint) (*model.PermissionModel, error) {
	var permission model.PermissionModel
	err := p.db.WithContext(p.ctx).Where("id = ?", id).First(&permission).Error
	return &permission, err
}

// GetDelPermissionByID 根据id获取已删除的权限详情
func (p *permissionDAO) GetDelPermissionByID(id uint) (*model.PermissionModel, error) {
	var permission model.PermissionModel
	err := p.db.WithContext(p.ctx).Unscoped().Where("id = ?", id).First(&permission).Error
	return &permission, err
}

// GetPermissionByName 根据name获取权限详情
func (p *permissionDAO) GetPermissionByName(name string) (*model.PermissionModel, error) {
	var permission model.PermissionModel
	err := p.db.WithContext(p.ctx).Where("name = ?", name).First(&permission).Error
	return &permission, err
}

// GetPermissionByPathMethod 根据path和method获取权限详情
func (p *permissionDAO) GetPermissionByPathMethod(apiPath string, method string) (*model.PermissionModel, error) {
	var permission model.PermissionModel
	err := p.db.WithContext(p.ctx).Where("api_path = ? AND method = ?", apiPath, method).First(&permission).Error
	return &permission, err
}

// UpdatePermission 更新权限
func (p *permissionDAO) UpdatePermission(permission *model.PermissionModel) error {
	global.Log.Debug("permissionDao:", permission)
	return Transaction(p.db, func(tx *gorm.DB) error {
		return tx.Where("id = ?", permission.ID).Updates(permission).Error
	})
}


// DeletePermission 删除权限
func (p *permissionDAO) DeletePermission(id uint) error {
	return Transaction(p.db, func(tx *gorm.DB) error {
		return tx.Delete(&model.PermissionModel{}, id).Error
	})
}

// RecoverPermission 恢复权限
func (p *permissionDAO) RecoverPermission(id uint) error {
	return Transaction(p.db, func(tx *gorm.DB) error {
		return tx.Unscoped().Model(&model.PermissionModel{}).Where("id = ?", id).Update("deleted_at", nil).Error
	})
}
