package dao

import (
	"context"
	"fmt"

	"gin-rbac/db/model"
	"gin-rbac/dtos"

	"gorm.io/gorm"
)

// UserDao 用户数据访问接口
type UserDao interface {
	// 获取用户列表
	GetUserList(publicUserListReqDTO *dtos.PublicUserListReqDTO) ([]*dtos.PublicUserDTO, int64, error)
	// 创建用户
	CreateUser(user *model.UserModel) (*model.UserModel, error)
	// 获取用户详情
	GetUserByID(id uint) (*dtos.FullUserDTO, error)
	// 获取公开用户详情
	GetPublicUserByID(id uint) (*dtos.PublicUserDTO, error)
	// 批量获取用户详情
	GetUserByIDs(ids []uint) ([]uint, []*model.UserModel, error)
	// 获取删除用户详情
	GetDelUserByID(id uint) (*model.UserModel, error)
	// 获取用户详情通过用户名
	GetUserByUsername(username string) (*model.UserModel, error)
	// 获取用户详情通过手机号
	GetUserByPhoneNum(phoneNum string) (*model.UserModel, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	GetUserPasswordByID(id uint) ([]string, error)
	UpdateUser(user *model.UserModel) error
	UpdateUserAvatar(userID, avatarID uint) error
	UpdatePassword(id uint, password string) error
	DeleteUser(id uint) error
	RecoverUser(id uint) error
}

// userDAO 用户数据访问实现
type userDAO struct {
	db *gorm.DB
	ctx context.Context
}

// NewUserDAO 创建用户数据访问实现
func NewUserDAO(db *gorm.DB,ctx context.Context) UserDao {
	return &userDAO{
		db: db,
		ctx: ctx,
	}
}

// 根据id 批量获取用户详情
func (u *userDAO) GetUserByIDs(ids []uint) ([]uint, []*model.UserModel, error) {
	var users []*model.UserModel
	existingIDsMap := make(map[uint]bool)
	nonExistingIDs := make([]uint, 0)
	if err := u.db.WithContext(u.ctx).Model(&model.UserModel{}).Where("id IN (?)", ids).Find(&users).Error; err != nil {
		return nil, nil, err
	}
	for _, user := range users {
		existingIDsMap[user.ID] = true
	}
	for _, id := range ids {
		if _, ok := existingIDsMap[id]; !ok {
			nonExistingIDs = append(nonExistingIDs, id)
		}
	}
	return nonExistingIDs, users, nil
}

// 根据用户名获取用户详情
func (u *userDAO) GetUserList(publicUserListReqDTO *dtos.PublicUserListReqDTO) ([]*dtos.PublicUserDTO, int64, error) {
	var users []*dtos.PublicUserDTO
	var count int64

	// 构造查询
	query := u.db.WithContext(u.ctx).Unscoped().Table("s_user u").
		Joins("LEFT JOIN b_image i ON i.id = u.avatar_id").
		Select("u.id AS id",
			"u.username AS username",
			"u.sex AS sex",
			"u.intro AS intro",
			"i.image_path AS avatar_path",
			"i.image_name AS avatar_name",
			"u.created_at AS created_at")

	// 模糊查询用户名
	if publicUserListReqDTO.Username != "" {
		query = query.Where("username LIKE ?", "%"+publicUserListReqDTO.Username+"%")
	}

	// 添加排序
	if publicUserListReqDTO.SortBy != "" && publicUserListReqDTO.Order != "" {
		query = query.Order(fmt.Sprintf("%s %s", publicUserListReqDTO.SortBy, publicUserListReqDTO.Order))
	}

	// 添加过滤条件
	for _, filter := range publicUserListReqDTO.Filters {
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
	offset := (publicUserListReqDTO.Page - 1) * publicUserListReqDTO.Size

	// 执行查询
	err := query.Count(&count).Offset(offset).Limit(publicUserListReqDTO.Size).Scan(&users).Error

	return users, count, err
}

// 创建用户
func (u *userDAO) CreateUser(user *model.UserModel) (*model.UserModel, error) {

	err := Transaction(u.db, func(tx *gorm.DB) error {
		return u.createUser(tx, user)
	})

	return user, err
}

func (u *userDAO) createUser(tx *gorm.DB, user *model.UserModel) (err error) {
	return tx.WithContext(u.ctx).Create(&user).Error
}

func (u *userDAO) GetUserByID(id uint) (*dtos.FullUserDTO, error) {
	var user *dtos.FullUserDTO
	if err := u.db.WithContext(u.ctx).Table("s_user u").
		Joins("LEFT JOIN b_image i ON i.id = u.avatar_id").
		Select("u.id AS id",
			"u.username AS username",
			"u.sex AS sex",
			"u.intro AS intro",
			"i.image_path AS avatar_path",
			"i.image_name AS avatar_name",
			"u.phone_num AS phone_num",
			"u.email AS email",
			"u.created_at AS created_at",
			"u.updated_at AS updated_at",
			"u.deleted_at AS deleted_at").
		Where("u.id = ?", id).
		Scan(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userDAO) GetPublicUserByID(id uint) (*dtos.PublicUserDTO, error) {
	var user *dtos.PublicUserDTO
	err := u.db.WithContext(u.ctx).Table("s_user u").
		Joins("LEFT JOIN b_image i ON i.id = u.avatar_id").
		Select("u.id AS id",
			"u.username AS username",
			"u.sex AS sex",
			"u.intro AS intro",
			"i.image_path AS avatar_path",
			"i.image_name AS avatar_name",
			"u.created_at AS created_at",
		).
		Where("u.id = ?", id).
		Scan(&user).Error
	return user, err
}

func (u *userDAO) GetDelUserByID(id uint) (*model.UserModel, error) {
	var user model.UserModel
	if err := u.db.WithContext(u.ctx).Unscoped().First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userDAO) GetUserPasswordByID(id uint) (password []string, err error) {
	// 通过id仅获取用户密码
	err = u.db.WithContext(u.ctx).Model(&model.UserModel{}).Select("password").Where("id = ?", id).Find(&password).Error

	return password, err
}

func (u *userDAO) UpdateUser(user *model.UserModel) error {
	return Transaction(u.db, func(tx *gorm.DB) error {
		return tx.WithContext(u.ctx).Where("id = ?", user.ID).Updates(user).Error
	})
}

func (u *userDAO) UpdateUserAvatar(userID, avatarID uint) error {
	return Transaction(u.db, func(tx *gorm.DB) error {
		return tx.WithContext(u.ctx).Model(&model.UserModel{}).Where("id = ?", userID).Update("avatar_id", avatarID).Error
	})
}

func (u *userDAO) UpdatePassword(id uint, password string) error {
	return Transaction(u.db, func(tx *gorm.DB) error {
		return tx.WithContext(u.ctx).Model(&model.UserModel{}).Where("id = ?", id).Update("password", password).Error
	})
}

func (u *userDAO) DeleteUser(id uint) error {
	return Transaction(u.db, func(tx *gorm.DB) error {
		return tx.WithContext(u.ctx).Delete(&model.UserModel{}, id).Error
	})
}

func (u *userDAO) RecoverUser(id uint) error {
	return Transaction(u.db, func(tx *gorm.DB) error {
		return tx.WithContext(u.ctx).Unscoped().Model(&model.UserModel{}).Where("id = ?", id).Update("deleted_at", nil).Error
	})
}

func (u *userDAO) GetUserByUsername(username string) (*model.UserModel, error) {
	var user model.UserModel
	if err := u.db.WithContext(u.ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userDAO) GetUserByPhoneNum(phoneNum string) (*model.UserModel, error) {
	var user model.UserModel
	if err := u.db.WithContext(u.ctx).Where("phone_num = ?", phoneNum).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userDAO) GetUserByEmail(email string) (*model.UserModel, error) {
	var user model.UserModel
	if err := u.db.WithContext(u.ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
