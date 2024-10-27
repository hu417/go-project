package service

import (
	"errors"

	"gin-rbac/common/errs"
	"gin-rbac/db/dao"
	"gin-rbac/dtos"
	"gin-rbac/global"

	"gorm.io/gorm"
)

// RolePermissionService 角色权限服务
type RolePermissionService interface {
	// CreateRolePermissionsByIDs 根据角色ID+权限ID列表创建角色权限
	CreateRolePermissionsByIDs(createRolePermissionsByIDsReqDTO dtos.CreateRolePermissionsByIDsReqDTO) ([]uint, []uint, error)
	// DeleteRolePermissionsByIDs 根据角色ID+权限ID列表删除角色权限
	DeleteRolePermissionsByIDs(deleteRolePermissionsByIDsReqDTO dtos.DeleteRolePermissionsByIDsReqDTO) error
	// GetRolePermissionsByPermissionID 根据权限ID获取角色权限列表
	GetRolePermissionsByPermissionID(getRolePermissionsByPermissionIDReqDTO dtos.GetRolePermissionsByPermissionIDReqDTO) (*dtos.GetRolePermissionsByPermissionIDResDTO, error)
	// GetRolePermissionsByRoleID 根据角色ID获取角色权限列表
	GetRolePermissionsByRoleID(getRolePermissionsByRoleIDReqDTO dtos.GetRolePermissionsByRoleIDReqDTO) (*dtos.GetRolePermissionsByRoleIDResDTO, error)
	// GetRolePermissionByID 根据角色、权限ID获取角色权限
	GetRolePermissionByID(getRolePermissionReqDTO dtos.GetRolePermissionByIDReqDTO) (*dtos.GetRolePermissionResDTO, error)
}

// rolePermissionService 角色权限服务实现
type rolePermissionService struct {
	roleDao           dao.RoleDao
	permissionDao     dao.PermissionDao
	rolePermissionDao dao.RolePermissionDao
}

// NewRolePermissionService 创建角色权限服务
func NewRolePermissionService(roleDao dao.RoleDao, permissionDao dao.PermissionDao, rolePermissionDao dao.RolePermissionDao) RolePermissionService {
	return &rolePermissionService{
		roleDao:           roleDao,
		permissionDao:     permissionDao,
		rolePermissionDao: rolePermissionDao,
	}
}

func (s *rolePermissionService) CreateRolePermissionsByIDs(
	createRolePermissionsByIDsReqDTO dtos.CreateRolePermissionsByIDsReqDTO) ([]uint, []uint, error) {
	// 验证角色是否存在
	nonExistingRoleIDs, _, err := s.roleDao.GetRoleByIDs(createRolePermissionsByIDsReqDTO.RoleIDList)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create role permission, Failed to get role: ", err)
		return nil, nil, errs.ErrInternalServerError
	}
	// 验证权限是否存在
	nonExistingPermissionIDs, _, err := s.permissionDao.GetPermissionByIDs(createRolePermissionsByIDsReqDTO.PermissionIDList)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create role permission, Failed to get permission: ", err)
		return nil, nil, errs.ErrInternalServerError
	}
	if len(nonExistingRoleIDs) > 0 || len(nonExistingPermissionIDs) > 0 {
		// 404 资源不存在
		global.Log.Warnln("Failed to create role permission, Failed to get role or permission: ", nonExistingPermissionIDs)
		return nonExistingRoleIDs, nonExistingPermissionIDs, errs.ErrRoleOrPermissionNotFound
	}
	// 创建角色权限
	if err = s.rolePermissionDao.CreateRolePermissionsByIDs(
		createRolePermissionsByIDsReqDTO.RoleIDList, createRolePermissionsByIDsReqDTO.PermissionIDList); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create role permission: ", err)
		return nil, nil, errs.ErrInternalServerError
	}
	return nil, nil, nil
}

func (s *rolePermissionService) DeleteRolePermissionsByIDs(deleteRolePermissionsByIDsReqDTO dtos.DeleteRolePermissionsByIDsReqDTO) error {
	if err := s.rolePermissionDao.DeleteRolePermissionsByIDs(
		deleteRolePermissionsByIDsReqDTO.RoleIDList, deleteRolePermissionsByIDsReqDTO.PermissionIDList); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to delete role permission, RolePermission not found: ", err)
			return errs.ErrRolePermissionNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to delete role permission: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

func (s *rolePermissionService) GetRolePermissionsByPermissionID(
	getRolePermissionsByPermissionIDReqDTO dtos.GetRolePermissionsByPermissionIDReqDTO) (
	*dtos.GetRolePermissionsByPermissionIDResDTO, error) {
	_, err := s.permissionDao.GetPermissionByID(getRolePermissionsByPermissionIDReqDTO.PermissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get role permission list, Permission not found: ", err)
			return nil, errs.ErrPermissionNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get role permission list, Failed to get permission: ", err)
		return nil, errs.ErrInternalServerError
	}
	rolePermissions, err := s.rolePermissionDao.GetRolePermissionsByPermissionID(getRolePermissionsByPermissionIDReqDTO.PermissionID)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get role permission list: ", err)
		return nil, errs.ErrInternalServerError
	}
	if len(rolePermissions) == 0 {
		// 404 资源不存在
		global.Log.Warnln("Failed to get role permission list, RolePermission not found: ", err)
		return nil, errs.ErrRolePermissionNotFound
	}
	roleWithPermissions := &dtos.GetRolePermissionsByPermissionIDResDTO{
		PermissionID:          rolePermissions[0].PermissionID,
		PermissionName:        rolePermissions[0].PermissionName,
		PermissionDescription: rolePermissions[0].PermissionDescription,
		PermissionAction:      rolePermissions[0].PermissionAction,
		Roles:                 make([]dtos.Role, 0),
	}
	for _, rolePermission := range rolePermissions {
		roleWithPermissions.Roles = append(roleWithPermissions.Roles, dtos.Role{
			RoleID:          rolePermission.RoleID,
			RoleName:        rolePermission.RoleName,
			RoleDescription: rolePermission.RoleDescription,
			RelativeAt:      rolePermission.CreatedAt,
			DelRelativeAt: gorm.DeletedAt{
				Time:  rolePermission.DeletedAt.Time,
				Valid: rolePermission.DeletedAt.Valid,
			},
		})
	}
	return roleWithPermissions, nil
}

func (s *rolePermissionService) GetRolePermissionsByRoleID(getRolePermissionsByRoleIDReqDTO dtos.GetRolePermissionsByRoleIDReqDTO) (
	*dtos.GetRolePermissionsByRoleIDResDTO, error) {
	_, err := s.roleDao.GetRoleByID(getRolePermissionsByRoleIDReqDTO.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get role permission list, Role not found: ", err)
			return nil, errs.ErrRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get role permission list, Failed to get role: ", err)
		return nil, errs.ErrInternalServerError
	}
	rolePermissions, err := s.rolePermissionDao.GetRolePermissionsByRoleID(getRolePermissionsByRoleIDReqDTO.RoleID)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get role permission list: ", err)
		return nil, errs.ErrInternalServerError
	}
	if len(rolePermissions) == 0 {
		// 404 资源不存在
		global.Log.Warnln("Failed to get role permission list, RolePermission not found: ", err)
		return nil, errs.ErrRolePermissionNotFound
	}
	permissionWithRoles := &dtos.GetRolePermissionsByRoleIDResDTO{
		RoleID:          rolePermissions[0].RoleID,
		RoleName:        rolePermissions[0].RoleName,
		RoleDescription: rolePermissions[0].RoleDescription,
		Permissions:     make([]dtos.Permission, 0),
	}
	for _, rolePermission := range rolePermissions {
		permissionWithRoles.Permissions = append(permissionWithRoles.Permissions, dtos.Permission{
			PermissionID:          rolePermission.PermissionID,
			PermissionName:        rolePermission.PermissionName,
			PermissionDescription: rolePermission.PermissionDescription,
			PermissionAction:      rolePermission.PermissionAction,
		})
	}
	return permissionWithRoles, nil
}

func (s *rolePermissionService) GetRolePermissionByID(
	getRolePermissionReqDTO dtos.GetRolePermissionByIDReqDTO) (*dtos.GetRolePermissionResDTO, error) {
	getRolePermissionResDTO := dtos.GetRolePermissionResDTO{}
	if _, err := s.roleDao.GetRoleByID(getRolePermissionReqDTO.RoleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get role permission, Role not found: ", err)
			return nil, errs.ErrRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get role permission, Failed to get role: ", err)
		return nil, errs.ErrInternalServerError
	}
	if _, err := s.permissionDao.GetPermissionByID(getRolePermissionReqDTO.PermissionID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get role permission, Permission not found: ", err)
			return nil, errs.ErrPermissionNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get role permission, Failed to get permission: ", err)
		return nil, errs.ErrInternalServerError
	}
	// 获取角色权限
	rolePermission, err := s.rolePermissionDao.GetRolePermissionByID(
		getRolePermissionReqDTO.RoleID, getRolePermissionReqDTO.PermissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get role permission, role permission not found: ", err)
			return nil, errs.ErrRolePermissionNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get role permission: ", err)
		return nil, errs.ErrInternalServerError
	}
	if rolePermission == nil {
		// 404 资源不存在
		global.Log.Warnln("Failed to get role permission, role permission not found: ", err)
		return nil, errs.ErrRolePermissionNotFound
	}

	if err := dtos.ConvertModelToDTO(rolePermission, &getRolePermissionResDTO); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get role permission, Failed to convert model to dto: ", err)
		return nil, errs.ErrInternalServerError
	}

	return &getRolePermissionResDTO, nil
}
