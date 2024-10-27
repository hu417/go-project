package service

import (
	"errors"

	"gin-rbac/common/errs"
	"gin-rbac/db/dao"
	"gin-rbac/db/model"
	"gin-rbac/dtos"
	"gin-rbac/global"

	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService interface {
	// 创建角色
	CreateRole(createRoleReqDTO dtos.CreateRoleReqDTO) error
	// 根据id获取角色
	GetRoleByID(getRoleReqDTO dtos.GetRoleReqDTO) (dtos.GetRoleResDTO, error)
	// 获取角色列表
	GetRoleList(getRoleListReqDTO dtos.GetRoleListReqDTO) (dtos.PaginationResult[dtos.GetRoleResDTO], error)
	// 更新角色
	UpdateRole(getRoleReqDTO dtos.GetRoleReqDTO, updateRoleReqDTO dtos.UpdateRoleReqDTO) error
	// 删除角色
	DeleteRole(deleteRoleReqDTO dtos.DeleteRoleReqDTO) error
	// 恢复角色
	RecoverRole(recoverRoleReqDTO dtos.RecoverRoleReqDTO) error
}

// roleService 角色服务实现
type roleService struct {
	roleDao dao.RoleDao
}

// NewRoleService 创建角色服务
func NewRoleService(roleDao dao.RoleDao) RoleService {
	return &roleService{
		roleDao: roleDao,
	}
}

// GetRoleList 获取角色列表
func (s *roleService) GetRoleList(getRoleListReqDTO dtos.GetRoleListReqDTO) (dtos.PaginationResult[dtos.GetRoleResDTO], error) {
	roles, total, err := s.roleDao.GetRoleList(&getRoleListReqDTO)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get role list: ", err)
		return dtos.PaginationResult[dtos.GetRoleResDTO]{}, errs.ErrInternalServerError
	}
	var roleResDTOs []dtos.GetRoleResDTO
	for _, role := range roles {
		var roleDTO = &dtos.GetRoleResDTO{}
		err = dtos.ConvertModelToDTO(role, roleDTO)
		if err != nil {
			// 500 服务器错误
			global.Log.Errorln("Failed to get role list, Failed to convert model to DTO: ", err)
			return dtos.PaginationResult[dtos.GetRoleResDTO]{}, errs.ErrInternalServerError
		}
		roleResDTOs = append(roleResDTOs, *roleDTO)
	}
	paginatedResult := dtos.NewPaginationResult(total, roleResDTOs, getRoleListReqDTO.Page, getRoleListReqDTO.Size)
	return paginatedResult, nil
}

// CreateRole 创建角色
func (s *roleService) CreateRole(createRoleReqDTO dtos.CreateRoleReqDTO) error {
	role := &model.RoleModel{}
	err := dtos.ConvertDTOToModel(&createRoleReqDTO, role)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create role, Failed to convert DTO to model: ", err)
		return errs.ErrInternalServerError
	}
	// 检查角色是否已存在，如果存在，则返回错误
	if err := s.checkDuplicateRoleExists(role.Name); err != nil {
		return err
	}
	// 创建角色
	err = s.roleDao.CreateRole(role)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create role: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// GetRoleByID 根据id获取角色
func (s *roleService) GetRoleByID(getRoleReqDTO dtos.GetRoleReqDTO) (dtos.GetRoleResDTO, error) {
	getRoleResDTO := &dtos.GetRoleResDTO{}
	role, err := s.roleDao.GetRoleByID(getRoleReqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get role, Role not found: ", err)
			return dtos.GetRoleResDTO{}, errs.ErrRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get role: ", err)
		return dtos.GetRoleResDTO{}, errs.ErrInternalServerError
	}
	err = dtos.ConvertModelToDTO(role, getRoleResDTO)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get role, Failed to convert model to DTO: ", err)
		return dtos.GetRoleResDTO{}, errs.ErrInternalServerError
	}

	return *getRoleResDTO, nil
}

// UpdateRole 更新角色
func (s *roleService) UpdateRole(getRoleReqDTO dtos.GetRoleReqDTO, updateRoleReqDTO dtos.UpdateRoleReqDTO) error {
	var role = &model.RoleModel{
		Model: gorm.Model{ID: getRoleReqDTO.ID},
	}
	err := dtos.ConvertDTOToModel(&updateRoleReqDTO, role)
	if err != nil {
		global.Log.Errorln("Failed to update role, Failed to convert DTO to model:", err)
		return errs.ErrInternalServerError
	}
	// 获取当前角色信息
	currentRole, err := s.roleDao.GetRoleByID(getRoleReqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to update role, Role not found: ", err)
			return errs.ErrRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to update role, Failed to get role: ", err)
		return errs.ErrInternalServerError
	}
	// 检查角色名称和描述是否发生变化
	if currentRole.Name == role.Name && currentRole.Description == role.Description {
		// 422 无字段更新
		global.Log.Warnln("Failed to update role, Role name and description are unchanged")
		return errs.ErrNoFieldsUpdated
	}
	// 检查角色名称是否已存在
	if currentRole.Name != role.Name {
		if err := s.checkDuplicateRoleExists(role.Name); err != nil {
			return err
		}
	}

	err = s.roleDao.UpdateRole(role)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to update role: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// DeleteRole 删除角色
func (s *roleService) DeleteRole(deleteRoleReqDTO dtos.DeleteRoleReqDTO) error {
	if _, err := s.roleDao.GetRoleByID(deleteRoleReqDTO.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to delete role, Role not found: ", err)
			return errs.ErrRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to delete role, Failed to get role: ", err)
		return errs.ErrInternalServerError
	}
	err := s.roleDao.DeleteRole(deleteRoleReqDTO.ID)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to delete role: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// RecoverRole 恢复角色
func (s *roleService) RecoverRole(recoverRoleReqDTO dtos.RecoverRoleReqDTO) error {
	if role, err := s.roleDao.GetDelRoleByID(recoverRoleReqDTO.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to recover role, Role not found", err)
			return errs.ErrRoleNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to recover role, Failed to get role", err)
		return errs.ErrInternalServerError
	} else if !role.DeletedAt.Valid {
		// 409 冲突
		global.Log.Warnln("Failed to recover role, Role not deleted", err)
		return errs.ErrRoleConflict
	}
	err := s.roleDao.RecoverRole(recoverRoleReqDTO.ID)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to recover role: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// checkDuplicateRoleExists 检查角色是否已存在
func (s *roleService) checkDuplicateRoleExists(roleName string) error {
	_, err := s.roleDao.GetRoleByName(roleName)
	// 如果发生错误且不是记录未找到的错误，则返回错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to check duplicate role exists: ", err)
		return errs.ErrInternalServerError
	}
	// 409 冲突
	global.Log.Warnln("Role already exists")
	return errs.ErrRoleAlreadyExists
}
