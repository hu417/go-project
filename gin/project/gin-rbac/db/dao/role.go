package dao

import (
	"context"
	"fmt"

	"gin-rbac/db/model"
	"gin-rbac/dtos"

	"gorm.io/gorm"
)

// RoleDao 角色数据访问接口
type RoleDao interface {
	GetRoleList(getRoleListReqDTO *dtos.GetRoleListReqDTO) ([]*model.RoleModel, int64, error)
	CreateRole(role *model.RoleModel) error
	GetRoleByID(id uint) (*model.RoleModel, error)
	GetRoleByIDs(ids []uint) ([]uint, []*model.RoleModel, error)
	GetDelRoleByID(id uint) (*model.RoleModel, error)
	GetRoleByName(name string) (*model.RoleModel, error)
	UpdateRole(role *model.RoleModel) error
	DeleteRole(id uint) error
	RecoverRole(id uint) error
}

// roleDAO 角色数据访问接口实现
type roleDAO struct {
	db  *gorm.DB
	ctx context.Context
}

// NewRoleDAO 创建角色数据访问接口实例
func NewRoleDAO(db *gorm.DB, ctx context.Context) RoleDao {
	return &roleDAO{
		db:  db,
		ctx: ctx,
	}
}

func (r *roleDAO) GetRoleByIDs(ids []uint) ([]uint, []*model.RoleModel, error) {
	var roles []*model.RoleModel
	existingIDsMap := make(map[uint]bool)
	nonExistingIDs := make([]uint, 0)
	if err := r.db.WithContext(r.ctx).Model(&model.RoleModel{}).Where("id IN (?)", ids).Find(&roles).Error; err != nil {
		return nil, nil, err
	}
	for _, role := range roles {
		existingIDsMap[role.ID] = true
	}
	for _, id := range ids {
		if _, ok := existingIDsMap[id]; !ok {
			nonExistingIDs = append(nonExistingIDs, id)
		}
	}
	return nonExistingIDs, roles, nil
}

// GetRoleList 获取角色列表
func (r *roleDAO) GetRoleList(getRoleListReqDTO *dtos.GetRoleListReqDTO) ([]*model.RoleModel, int64, error) {
	var roles []*model.RoleModel
	var total int64
	query := r.db.WithContext(r.ctx).Model(&model.RoleModel{})
	if getRoleListReqDTO.Name != "" {
		query = query.Where("name LIKE ?", "%"+getRoleListReqDTO.Name+"%")
	}

	// 添加排序
	if getRoleListReqDTO.SortBy != "" && getRoleListReqDTO.Order != "" {
		query = query.Order(fmt.Sprintf("%s %s", getRoleListReqDTO.SortBy, getRoleListReqDTO.Order))
	}

	// 添加过滤条件
	for _, filter := range getRoleListReqDTO.Filters {
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
	offset := (getRoleListReqDTO.Page - 1) * getRoleListReqDTO.Size

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Limit(getRoleListReqDTO.Size).Offset(offset).Find(&roles).Error
	return roles, total, err
}

// CreateRole 创建角色
func (r *roleDAO) CreateRole(role *model.RoleModel) error {
	return Transaction(r.db, func(tx *gorm.DB) error {
		return tx.Create(role).Error
	})
}

// GetRoleByID 根据ID获取角色
func (r *roleDAO) GetRoleByID(id uint) (*model.RoleModel, error) {
	var role model.RoleModel
	err := r.db.WithContext(r.ctx).Where("id = ?", id).First(&role).Error
	return &role, err
}

// GetDelRoleByID 根据ID获取角色，包括已删除角色
func (r *roleDAO) GetDelRoleByID(id uint) (*model.RoleModel, error) {
	var role model.RoleModel
	err := r.db.WithContext(r.ctx).Unscoped().Where("id = ?", id).First(&role).Error
	return &role, err
}

// GetRoleByName 根据名称获取角色
func (r *roleDAO) GetRoleByName(name string) (*model.RoleModel, error) {
	var role model.RoleModel
	if err := r.db.WithContext(r.ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// UpdateRole 更新角色
func (r *roleDAO) UpdateRole(role *model.RoleModel) error {
	return Transaction(r.db, func(tx *gorm.DB) error {
		result := tx.Where("id = ?", role.ID).Updates(role)
		return result.Error
	})
}

// DeleteRole 删除角色
func (r *roleDAO) DeleteRole(id uint) error {
	return Transaction(r.db, func(tx *gorm.DB) error {
		return tx.Delete(&model.RoleModel{}, id).Error
	})
}

// RecoverRole 恢复删除的角色
func (r *roleDAO) RecoverRole(id uint) error {
	return Transaction(r.db, func(tx *gorm.DB) error {
		return tx.Unscoped().Model(&model.RoleModel{}).Where("id = ?", id).Update("deleted_at", nil).Error
	})
}
