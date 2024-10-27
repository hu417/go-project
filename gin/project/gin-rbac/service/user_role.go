package service

import (
	"errors"

	"gin-rbac/common/errs"
	"gin-rbac/db/dao"
	"gin-rbac/dtos"
	"gin-rbac/global"

	"gorm.io/gorm"
)

// UserRoleService 用户角色服务
type UserRoleService interface {
	// CreateUserRolesByIDs 根据用户id+角色id创建用户角色
	CreateUserRolesByIDs(createUserRolesByIDsReqDTO dtos.CreateUserRolesByIDsReqDTO) ([]uint, []uint, error)
	// DeleteUserRolesByIDs 根据用户id+角色id删除用户角色
	DeleteUserRolesByIDs(deleteUserRolesByIDsReqDTO dtos.DeleteUserRolesByIDsReqDTO) error
	// GetUserRolesByUserID 根据用户id获取用户角色
	GetUserRolesByUserID(getUserRolesByUserIDReqDTO dtos.GetUserRolesByUserIDReqDTO) (*dtos.GetUserRolesByUserIDResDTO, error)
	// GetUserRolesByRoleID 根据角色id获取用户角色
	GetUserRolesByRoleID(getUserRolesByRoleIDReqDTO dtos.GetUserRolesByRoleIDReqDTO) (*dtos.GetUserRolesByRoleIDResDTO, error)
	// GetUserRoleByID 根据用户角色id获取用户角色
	GetUserRoleByID(getUserRoleByIDReqDTO dtos.GetUserRoleByIDReqDTO) (*dtos.GetUserRoleByIDResDTO, error)
}

// userRoleService 用户角色服务实现
type userRoleService struct {
	userDao     dao.UserDao
	roleDao     dao.RoleDao
	userRoleDao dao.UserRoleDao
}

// NewUserRoleService 创建用户角色服务
func NewUserRoleService(userDao dao.UserDao, roleDao dao.RoleDao, userRoleDao dao.UserRoleDao) UserRoleService {
	return &userRoleService{
		userDao:     userDao,
		roleDao:     roleDao,
		userRoleDao: userRoleDao}
}

// CreateUserRolesByIDs 根据用户id+角色id创建用户角色
func (s *userRoleService) CreateUserRolesByIDs(createUserRolesByIDsReqDTO dtos.CreateUserRolesByIDsReqDTO) ([]uint, []uint, error) {
	// 验证用户是否存在
	nonExistingUserIDs, _, err := s.userDao.GetUserByIDs(createUserRolesByIDsReqDTO.UserIDList)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create user role, Failed to get user: ", err)
		return nil, nil, errs.ErrInternalServerError
	}
	// 验证角色是否存在
	nonExistingRoleIDs, _, err := s.roleDao.GetRoleByIDs(createUserRolesByIDsReqDTO.RoleIDList)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create user role, Failed to get role: ", err)
		return nil, nil, errs.ErrInternalServerError
	}
	if len(nonExistingUserIDs) > 0 || len(nonExistingRoleIDs) > 0 {
		// 400 请求参数错误
		global.Log.Warnln("Failed to create user role, non existing user or role: ", nonExistingUserIDs, nonExistingRoleIDs)
		return nonExistingUserIDs, nonExistingRoleIDs, errs.ErrUserRoleAlreadyExists
	}
	// 创建用户角色
	if err = s.userRoleDao.CreateUserRolesByIDs(createUserRolesByIDsReqDTO.UserIDList, createUserRolesByIDsReqDTO.RoleIDList); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create user role: ", err)
		return nil, nil, errs.ErrInternalServerError
	}
	return nil, nil, nil
}

// DeleteUserRolesByIDs 根据用户id+角色id删除用户角色
func (s *userRoleService) DeleteUserRolesByIDs(deleteUserRolesByIDsReqDTO dtos.DeleteUserRolesByIDsReqDTO) error {
	if err := s.userRoleDao.DeleteUserRolesByIDs(
		deleteUserRolesByIDsReqDTO.UserIDList, deleteUserRolesByIDsReqDTO.RoleIDList); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warnln("Failed to delete user role, UserRole not found: ", err)
			return errs.ErrUserRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to delete user role: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// GetUserRolesByUserID 根据用户id获取用户角色
func (s *userRoleService) GetUserRolesByUserID(getUserRolesByUserIDReqDTO dtos.GetUserRolesByUserIDReqDTO) (*dtos.GetUserRolesByUserIDResDTO, error) {
	_, err := s.userDao.GetUserByID(getUserRolesByUserIDReqDTO.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get user role by user id, User not found: ", err)
			return nil, errs.ErrUserNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get user role by user id: ", err)
		return nil, errs.ErrInternalServerError
	}
	userRoles, err := s.userRoleDao.GetUserRolesByUserID(getUserRolesByUserIDReqDTO.UserID)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get user role by user id: ", err)
		return nil, errs.ErrInternalServerError
	}
	if len(userRoles) == 0 {
		// 404 资源不存在
		global.Log.Warnln("Failed to get user role by user id, UserRole not found: ", err)
		return nil, errs.ErrUserRoleNotFound
	}
	roleWithUsers := &dtos.GetUserRolesByUserIDResDTO{
		UserID:   userRoles[0].UserID,
		Username: userRoles[0].Username,
		Roles:    make([]dtos.Role, 0),
	}
	for _, userRole := range userRoles {
		roleWithUsers.Roles = append(roleWithUsers.Roles, dtos.Role{
			RoleID:          userRole.RoleID,
			RoleName:        userRole.RoleName,
			RoleDescription: userRole.RoleDescription,
			RelativeAt:      userRole.CreatedAt,
			DelRelativeAt: gorm.DeletedAt{
				Time:  userRole.DeletedAt.Time,
				Valid: userRole.DeletedAt.Valid,
			},
		})
	}
	return roleWithUsers, nil
}

// GetUserRolesByRoleID 根据角色id获取用户角色
func (s *userRoleService) GetUserRolesByRoleID(getUserRolesByRoleIDReqDTO dtos.GetUserRolesByRoleIDReqDTO) (*dtos.GetUserRolesByRoleIDResDTO, error) {
	// 验证角色是否存在
	_, err := s.roleDao.GetRoleByID(getUserRolesByRoleIDReqDTO.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get user role by role id, Role not found: ", err)
			return nil, errs.ErrRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get user role by role id: ", err)
		return nil, errs.ErrInternalServerError
	}

	// 获取角色下的用户角色
	userRoles, err := s.userRoleDao.GetUserRolesByRoleID(getUserRolesByRoleIDReqDTO.RoleID)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get user role by role id: ", err)
		return nil, errs.ErrInternalServerError
	}
	if len(userRoles) == 0 {
		// 404 资源不存在
		global.Log.Warnln("Failed to get user role by role id, UserRole not found: ", err)
		return nil, errs.ErrUserRoleNotFound
	}
	userWithRoles := &dtos.GetUserRolesByRoleIDResDTO{
		RoleID:          userRoles[0].RoleID,
		RoleName:        userRoles[0].RoleName,
		RoleDescription: userRoles[0].RoleDescription,
		Users:           make([]dtos.User, 0),
	}
	for _, userRole := range userRoles {
		userWithRoles.Users = append(userWithRoles.Users, dtos.User{
			UserID:     userRole.UserID,
			Username:   userRole.Username,
			RelativeAt: userRole.CreatedAt,
			DelRelativeAt: gorm.DeletedAt{
				Time:  userRole.DeletedAt.Time,
				Valid: userRole.DeletedAt.Valid,
			},
		})
	}
	return userWithRoles, nil
}

// GetUserRoleByID 根据用户id+角色id获取用户角色信息
func (s *userRoleService) GetUserRoleByID(getUserRoleByIDReqDTO dtos.GetUserRoleByIDReqDTO) (*dtos.GetUserRoleByIDResDTO, error) {
	getUserRoleByIDResDTO := &dtos.GetUserRoleByIDResDTO{}
	if _, err := s.userDao.GetUserByID(getUserRoleByIDReqDTO.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get user role by id, User not found: ", err)
			return nil, errs.ErrUserNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get user role by id: ", err)
		return nil, errs.ErrInternalServerError
	}
	if _, err := s.roleDao.GetRoleByID(getUserRoleByIDReqDTO.RoleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get user role by id, Role not found: ", err)
			return nil, errs.ErrRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get user role by id: ", err)
		return nil, errs.ErrInternalServerError
	}
	// 获取用户角色
	userRole, err := s.userRoleDao.GetUserRoleByID(getUserRoleByIDReqDTO.UserID, getUserRoleByIDReqDTO.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get user role by id, UserRole not found: ", err)
			return nil, errs.ErrUserRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get user role by id: ", err)
		return nil, errs.ErrInternalServerError
	}

	if userRole == nil {
		global.Log.Warnln("Failed to get user role by id, UserRole not found: ", err)
		return nil, errs.ErrUserRoleNotFound
	}

	if err = dtos.ConvertModelToDTO(userRole, getUserRoleByIDResDTO); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get user role by id, Failed to convert model to dto: ", err)
		return nil, errs.ErrInternalServerError
	}
	return getUserRoleByIDResDTO, nil
}
