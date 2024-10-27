package dao

import (
	"context"

	"gin-rbac/db/model"
	"gin-rbac/dtos"


	"gorm.io/gorm"
)

type RolePermissionDao interface {
	// 根据角色ID和权限ID创建角色权限关联记录
	CreateRolePermissionsByIDs(roleIDs, permissionIDs []uint) error
	// 根据角色ID和权限ID删除角色权限关联记录
	DeleteRolePermissionsByIDs(roleIDs, permissionIDs []uint) error
	// 根据权限ID获取角色权限关联记录
	GetRolePermissionsByPermissionID(permissionID uint) ([]*dtos.GetRolePermissionResDTO, error)
	// 根据角色ID获取角色权限关联记录
	GetRolePermissionsByRoleID(roleID uint) ([]*dtos.GetRolePermissionResDTO, error)
	// 根据角色ID和权限ID获取角色权限关联记录
	GetRolePermissionByID(roleID, permissionID uint) (*dtos.GetRolePermissionResDTO, error)
}

type rolePermissionDao struct {
	db *gorm.DB
	ctx context.Context
}

// NewRolePermissionDAO 创建角色权限DAO
func NewRolePermissionDAO(db *gorm.DB,ctx context.Context) RolePermissionDao {
	return &rolePermissionDao{
		db: db,
		ctx: ctx,
	}
}

// CreateRolePermissionsByIDs 创建角色权限关联记录（根据角色ID和权限ID)
func (r *rolePermissionDao) CreateRolePermissionsByIDs(roleIDs, permissionIDs []uint) error {
	return Transaction(r.db, func(tx *gorm.DB) error {
		var existingRolesPermission, restoredRolesPermission, createRolePermissions []*model.RolePermissionModel
		existingPairs := make(map[[2]uint]bool)
		// 查询已存在[包含已删除]的记录
		if err := tx.WithContext(r.ctx).Unscoped().
			Where("role_id IN (?) AND permission_id IN (?)", roleIDs, permissionIDs).
			Find(&existingRolesPermission).Error; err != nil {
			return err
		}
		// 分别标记已删除的记录和已存在的记录
		for _, RolesPermission := range existingRolesPermission {
			if RolesPermission.DeletedAt.Valid {
				RolesPermission.DeletedAt = gorm.DeletedAt{}
				// 标记已删除的记录
				restoredRolesPermission = append(restoredRolesPermission, RolesPermission)
			}
			// 使用键值对来标记已存在的记录
			existingPairs[[2]uint{RolesPermission.RoleID, RolesPermission.PermissionID}] = true
		}
		// 恢复已删除的记录
		if len(restoredRolesPermission) > 0 {
			if err := tx.WithContext(r.ctx).Unscoped().Save(restoredRolesPermission).Error; err != nil {
				return err
			}
		}
		// 构建要创建的记录
		for _, roleID := range roleIDs {
			for _, permissionID := range permissionIDs {
				pair := [2]uint{roleID, permissionID}
				if _, ok := existingPairs[pair]; !ok {
					createRolePermission := &model.RolePermissionModel{
						RoleID:       pair[0],
						PermissionID: pair[1],
					}
					createRolePermissions = append(createRolePermissions, createRolePermission)
				}
			}
		}
		// 批量创建不存在的记录
		if len(createRolePermissions) > 0 {
			if err := tx.WithContext(r.ctx).Create(createRolePermissions).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteRolePermissionsByIDs 删除角色权限关联记录(根据角色ID和权限ID)
func (r *rolePermissionDao) DeleteRolePermissionsByIDs(roleIDs, permissionIDs []uint) error {
	return Transaction(r.db, func(tx *gorm.DB) error {
		var rolePermissions []*model.RolePermissionModel
		if err := tx.WithContext(r.ctx).
			Where("role_id IN (?) AND permission_id IN (?)", roleIDs, permissionIDs).
			Find(&rolePermissions).Error; err != nil {
			return err
		}
		// 如果没有记录需要删除，则返回错误信息
		if len(rolePermissions) == 0 {
			return gorm.ErrRecordNotFound
		}
		// 删除指定角色和权限的关联记录
		err := tx.WithContext(r.ctx).Where("role_id IN (?) AND permission_id IN (?)", roleIDs, permissionIDs).Delete(&model.RolePermissionModel{}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

// GetRolePermissionsByPermissionID 根据权限ID获取角色权限关联记录
func (r *rolePermissionDao) GetRolePermissionsByPermissionID(permissionID uint) ([]*dtos.GetRolePermissionResDTO, error) {
	var rolePermissions []*dtos.GetRolePermissionResDTO
	err := r.db.WithContext(r.ctx).Unscoped().Table("s_permission p").
		Joins("LEFT JOIN r_role_permission rp ON rp.permission_id = p.id").
		Joins("LEFT JOIN s_role r ON r.id = rp.role_id").
		Select("r.id AS role_id, "+
			"r.name AS role_name, "+
			"r.description AS role_description, "+
			"p.id AS permission_id, "+
			"p.name AS permission_name, "+
			"p.description AS permission_description, "+
			"p.method AS permission_method, "+
			"rp.created_at AS created_at, "+
			"rp.deleted_at AS deleted_at").
		Where("permission_id = ?", permissionID).
		Scan(&rolePermissions).Error
	return rolePermissions, err
}

// GetRolePermissionsByRoleID 根据角色ID获取角色权限关联记录
func (r *rolePermissionDao) GetRolePermissionsByRoleID(roleID uint) ([]*dtos.GetRolePermissionResDTO, error) {
	var rolePermissions []*dtos.GetRolePermissionResDTO
	err := r.db.WithContext(r.ctx).Unscoped().Table("s_role r").
		Joins("LEFT JOIN r_role_permission rp ON rp.role_id = r.id").
		Joins("LEFT JOIN s_permission p ON p.id = rp.permission_id").
		Select("r.id AS role_id, "+
			"r.name AS role_name, "+
			"r.description AS role_description, "+
			"p.id AS permission_id, "+
			"p.name AS permission_name, "+
			"p.description AS permission_description, "+
			"p.method AS permission_method, "+
			"rp.created_at AS created_at, "+
			"rp.deleted_at AS deleted_at").
		Where("role_id = ?", roleID).
		Scan(&rolePermissions).Error
	return rolePermissions, err
}

// GetRolePermissionByID 根据角色ID和权限ID获取角色权限关联记录
func (r *rolePermissionDao) GetRolePermissionByID(roleID, permissionID uint) (*dtos.GetRolePermissionResDTO, error) {
	var rolePermission *dtos.GetRolePermissionResDTO
	err := r.db.WithContext(r.ctx).Unscoped().Table("s_role r").
		Joins("LEFT JOIN r_role_permission rp ON rp.role_id = r.id").
		Joins("LEFT JOIN s_permission p ON p.id = rp.permission_id").
		Select("r.id AS role_id, "+
			"r.name AS role_name, "+
			"r.description AS role_description, "+
			"p.id AS permission_id, "+
			"p.name AS permission_name, "+
			"p.description AS permission_description, "+
			"p.method AS permission_method, "+
			"rp.created_at AS created_at, "+
			"rp.deleted_at AS deleted_at").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Scan(&rolePermission).Error
	return rolePermission, err
}
