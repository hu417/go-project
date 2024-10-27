package dao

import (
	"context"

	"gin-rbac/db/model"
	"gin-rbac/dtos"

	"gorm.io/gorm"
)

// UserRoleDao 用户角色数据访问接口
type UserRoleDao interface {
	// 根据用户id和角色id创建用户角色记录
	CreateUserRolesByIDs(userIDs, roleIDs []uint) error
	// 根据用户id和角色id删除用户角色记录
	DeleteUserRolesByIDs(userIDs, roleIDs []uint) error
	// 根据用户id获取用户角色记录
	GetUserRolesByUserID(userID uint) ([]*dtos.GetUserRoleByIDResDTO, error)
	// 根据角色id获取用户角色记录
	GetUserRolesByRoleID(roleID uint) ([]*dtos.GetUserRoleByIDResDTO, error)
	// 根据用户id和角色id获取用户角色记录
	GetUserRoleByID(userID, roleID uint) (*dtos.GetUserRoleByIDResDTO, error)
}

// userRoleDao 用户角色数据访问实现
type userRoleDao struct {
	db  *gorm.DB
	ctx context.Context
}

// NewUserRoleDAO 创建用户角色数据访问实现
func NewUserRoleDAO(db *gorm.DB, ctx context.Context) UserRoleDao {
	return &userRoleDao{
		db:  db,
		ctx: ctx,
	}
}

// CreateUserRolesByIDs 根据用户id和角色id创建用户角色记录
func (u *userRoleDao) CreateUserRolesByIDs(userIDs, roleIDs []uint) error {
	return Transaction(u.db, func(tx *gorm.DB) error {
		var existingUserRole, restoredUserRole, createUserRole []*model.UserRoleModel
		existingPairs := make(map[[2]uint]bool)

		// 查询已存在[包括已删除]的用户角色记录
		if err := tx.WithContext(u.ctx).Unscoped().
			Where("user_id IN (?) AND role_id IN (?)", userIDs, roleIDs).
			Find(&existingUserRole).Error; err != nil {
			return err
		}
		// 分别标记已删除的记录和已存在的记录
		for _, userRole := range existingUserRole {
			if userRole.DeletedAt.Valid {
				userRole.DeletedAt = gorm.DeletedAt{}
				// 标记已删除的记录
				restoredUserRole = append(restoredUserRole, userRole)
			}
			// 使用键值对来标记已存在的记录
			existingPairs[[2]uint{userRole.UserID, userRole.RoleID}] = true
		}
		// 恢复已删除的记录
		if len(restoredUserRole) > 0 {
			if err := tx.WithContext(u.ctx).Unscoped().Save(restoredUserRole).Error; err != nil {
				return err
			}
		}
		// 构建要创建的记录
		for _, userID := range userIDs {
			for _, roleID := range roleIDs {
				pair := [2]uint{userID, roleID}
				if _, ok := existingPairs[pair]; !ok {
					createUserRole = append(createUserRole, &model.UserRoleModel{
						UserID: userID,
						RoleID: roleID,
					})
				}
			}
		}
		// 批量创建不存在的记录
		if len(createUserRole) > 0 {
			if err := tx.WithContext(u.ctx).Create(createUserRole).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteUserRolesByIDs 根据用户id和角色id删除用户角色记录
func (u *userRoleDao) DeleteUserRolesByIDs(userIDs, roleIDs []uint) error {
	return Transaction(u.db, func(tx *gorm.DB) error {
		var userRole []*model.UserRoleModel
		if err := tx.WithContext(u.ctx).
			Where("user_id IN (?) AND role_id IN (?)", userIDs, roleIDs).
			Find(&userRole).Error; err != nil {
			return err
		}
		// 如果有记录需要删除，则返回错误信息
		if len(userRole) == 0 {
			return gorm.ErrRecordNotFound
		}
		// 删除指定用户和角色的关联记录
		err := tx.WithContext(u.ctx).Where("user_id IN (?) AND role_id IN (?)", userIDs, roleIDs).Delete(&model.UserRoleModel{}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

// GetUserRolesByUserID 根据用户id获取用户角色记录
func (u *userRoleDao) GetUserRolesByUserID(userID uint) ([]*dtos.GetUserRoleByIDResDTO, error) {
	var userRoles []*dtos.GetUserRoleByIDResDTO
	err := u.db.WithContext(u.ctx).Unscoped().Table("s_user u").
		Joins("LEFT JOIN r_user_role ur ON ur.user_id = u.id").
		Joins("LEFT JOIN s_role r ON r.id = ur.role_id").
		Select("u.id AS user_id, "+
			"u.username AS username, "+
			"r.id AS role_id, "+
			"r.name AS role_name, "+
			"r.description AS role_description, "+
			"ur.created_at AS created_at, "+
			"ur.deleted_at AS deleted_at").
		Where("user_id = ?", userID).
		Scan(&userRoles).Error
	return userRoles, err
}

// GetUserRolesByRoleID 根据角色id获取用户角色记录
func (u *userRoleDao) GetUserRolesByRoleID(roleID uint) ([]*dtos.GetUserRoleByIDResDTO, error) {
	var userRoles []*dtos.GetUserRoleByIDResDTO
	err := u.db.WithContext(u.ctx).Unscoped().Table("s_role r").
		Joins("LEFT JOIN r_user_role ur ON ur.role_id = r.id").
		Joins("LEFT JOIN s_user u ON u.id = ur.user_id").
		Select("r.id AS role_id, "+
			"r.name AS role_name, "+
			"r.description AS role_description, "+
			"u.id AS user_id, "+
			"u.username AS username, "+
			"ur.created_at AS created_at, "+
			"ur.deleted_at AS deleted_at").
		Where("role_id = ?", roleID).
		Scan(&userRoles).Error

	return userRoles, err
}

// GetUserRoleByID 根据用户id和角色id获取用户角色记录
func (u *userRoleDao) GetUserRoleByID(userID, roleID uint) (*dtos.GetUserRoleByIDResDTO, error) {
	var userRole *dtos.GetUserRoleByIDResDTO
	err := u.db.WithContext(u.ctx).Unscoped().Table("s_user u").
		Joins("LEFT JOIN r_user_role ur ON ur.user_id = u.id").
		Joins("LEFT JOIN s_role r ON r.id = ur.role_id").
		Select("u.id AS user_id, "+
			"u.username AS username, "+
			"r.id AS role_id, "+
			"r.name AS role_name, "+
			"r.description AS role_description, "+
			"ur.created_at AS created_at, "+
			"ur.deleted_at AS deleted_at").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Scan(&userRole).Error
	return userRole, err
}
