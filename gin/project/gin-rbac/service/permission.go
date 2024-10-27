package service

import (
	"errors"

	"gin-rbac/common/errs"
	"gin-rbac/db/dao"
	"gin-rbac/db/model"
	"gin-rbac/dtos"
	"gin-rbac/global"
	"gin-rbac/utils"

	"gorm.io/gorm"
)

// PermissionService 角色服务
type PermissionService interface {
	// GetPermissionList 获取权限列表
	GetPermissionList(getPermissionListReqDTO dtos.GetPermissionListReqDTO) (dtos.PaginationResult[dtos.GetPermissionResDTO], error)
	// CreatePermission 创建权限
	CreatePermission(createPermissionReqDTO dtos.CreatePermissionReqDTO) error
	// GetPermission 根据id获取权限
	GetPermission(getPermissionReqDTO dtos.GetPermissionReqDTO) (dtos.GetPermissionResDTO, error)
	// UpdatePermission 根据id更新权限
	UpdatePermission(getPermissionReqDTO dtos.GetPermissionReqDTO, updatePermissionReqDTO dtos.UpdatePermissionReqDTO) error
	// DeletePermission 根据id删除权限
	DeletePermission(deletePermissionReqDTO dtos.DeletePermissionReqDTO) error
	// RecoverPermission 根据id恢复权限
	RecoverPermission(recoverPermissionReqDTO dtos.RecoverPermissionReqDTO) error
}

// permissionService 角色服务实现
type permissionService struct {
	permissionDao dao.PermissionDao
}

// NewPermissionService 创建角色服务
func NewPermissionService(permissionDao dao.PermissionDao) PermissionService {
	return &permissionService{
		permissionDao: permissionDao, 
	}
}

// GetPermissionList 获取权限列表
func (p *permissionService) GetPermissionList(getPermissionListReqDTO dtos.GetPermissionListReqDTO) (
	dtos.PaginationResult[dtos.GetPermissionResDTO], error) {
	permissionList, total, err := p.permissionDao.GetPermissionList(&getPermissionListReqDTO)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get permission list: ", err)
		return dtos.PaginationResult[dtos.GetPermissionResDTO]{}, errs.ErrInternalServerError
	}
	var permissionData []dtos.GetPermissionResDTO
	for _, permission := range permissionList {
		var permissionDTO = &dtos.GetPermissionResDTO{}
		err = dtos.ConvertModelToDTO(permission, permissionDTO)
		if err != nil {
			// 500 服务器错误
			global.Log.Errorln("Failed to get permission list, Failed to convert model to dto: ", err)
			return dtos.PaginationResult[dtos.GetPermissionResDTO]{}, errs.ErrInternalServerError
		}
		permissionData = append(permissionData, *permissionDTO)
	}
	// 创建分页结果
	paginatedResult := dtos.NewPaginationResult[dtos.GetPermissionResDTO](
		total, permissionData, getPermissionListReqDTO.Page, getPermissionListReqDTO.Size,
	)
	return paginatedResult, nil
}

// CreatePermission 创建权限
func (p *permissionService) CreatePermission(createPermissionReqDTO dtos.CreatePermissionReqDTO) error {
	permission := &model.PermissionModel{}
	err :=dtos.ConvertDTOToModel(&createPermissionReqDTO, permission)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create permission, Failed to convert DTO to model: ", err)
		return errs.ErrInternalServerError
	}
	//检查API路径和HTTP方法格式是否有效，如果无效，返回错误
	if !utils.IsValidPathAndMethod(permission.ApiPath, permission.Method) {
		// 400 客户端错误
		global.Log.Warnln("Failed to create permission, Invalid path method action: permission.ApiPath: ",
			permission.ApiPath, " permission.Method: ", permission.Method)
		return errs.ErrInvalidAPIPathMethodFormat
	}
	// 检查权限名是否重复，如果重复，返回错误
	if err = p.checkDuplicatePermissionExists(permission.Name); err != nil {
		return err
	}
	// 检查API路径和HTTP方法是否重复，如果重复，返回错误
	if err = p.checkDuplicatePermissionPathMethodExists(permission.ApiPath, permission.Method); err != nil {
		return err
	}
	if err := p.permissionDao.CreatePermission(permission); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create permission: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// GetPermission 根据id获取权限
func (p *permissionService) GetPermission(getPermissionReqDTO dtos.GetPermissionReqDTO) (dtos.GetPermissionResDTO, error) {
	permissionResDTO := &dtos.GetPermissionResDTO{}
	permission, err := p.permissionDao.GetPermissionByID(getPermissionReqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to get permission, Permission not found: ", err)
			return dtos.GetPermissionResDTO{}, errs.ErrPermissionNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get permission:", err)
		return dtos.GetPermissionResDTO{}, errs.ErrInternalServerError
	}
	if err = dtos.ConvertModelToDTO(permission, permissionResDTO); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get permission, Failed to convert model to dto:", err)
		return dtos.GetPermissionResDTO{}, errs.ErrInternalServerError
	}
	return *permissionResDTO, nil
}

// UpdatePermission 根据id更新权限
func (p *permissionService) UpdatePermission(getPermissionReqDTO dtos.GetPermissionReqDTO, updatePermissionReqDTO dtos.UpdatePermissionReqDTO) error {

	// 给权限绑定ID，因为ConvertDTOToModel函数无法解析到ID
	permission := model.PermissionModel{
		Model: gorm.Model{ID: getPermissionReqDTO.ID},
	}
	// 将dto转换为model
	if err := dtos.ConvertDTOToModel(&updatePermissionReqDTO, &permission); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to update permission, Failed to convert DTO to model:", err)
		return errs.ErrInternalServerError
	}
	// 获取当前权限, 判断权限是否存在, 如果不存在返回错误
	currentPermission, err := p.permissionDao.GetPermissionByID(getPermissionReqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to update permission, Permission not found:", err)
			return errs.ErrPermissionNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to update permission, Failed to get permission:", err)
		return errs.ErrInternalServerError
	}
	// 检查是否有字段变化
	hasChanges := false
	if currentPermission.Name != permission.Name ||
		currentPermission.ApiPath != permission.ApiPath ||
		currentPermission.Method != permission.Method ||
		currentPermission.Description != permission.Description {
		hasChanges = true
	}

	if !hasChanges {
		// 400 客户端错误
		global.Log.Warnln("Failed to update permission, No fields to update")
		return errs.ErrNoFieldsUpdated
	}
	//检查API路径和HTTP方法格式是否有效，如果无效，返回错误
	if !utils.IsValidPathAndMethod(permission.ApiPath, permission.Method) {
		// 400 客户端错
		global.Log.Warnln("Failed to update permission, Invalid path method action: permission.ApiPath: ",
			permission.ApiPath, " permission.Method: ", permission.Method)
		return errs.ErrInvalidAPIPathMethodFormat
	}
	// 检查权限名是否重复，如果重复，返回错误
	if currentPermission.Name != permission.Name {
		if err = p.checkDuplicatePermissionExists(permission.Name); err != nil {
			return err
		}
	}
	// 检查API路径和HTTP方法是否重复，如果重复，返回错误
	if currentPermission.ApiPath != permission.ApiPath || currentPermission.Method != permission.Method {
		if err = p.checkDuplicatePermissionPathMethodExists(permission.ApiPath, permission.Method); err != nil {
			return err
		}
	}
	// 更新权限
	if err := p.permissionDao.UpdatePermission(&permission); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to update permission: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// DeletePermission 根据id删除权限
func (p *permissionService) DeletePermission(deletePermissionReqDTO dtos.DeletePermissionReqDTO) error {
	if _, err := p.permissionDao.GetPermissionByID(deletePermissionReqDTO.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to delete permission, Permission not found: ", err)
			return errs.ErrPermissionNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to delete permission, Failed to get permission: ", err)
		return errs.ErrInternalServerError
	}
	if err := p.permissionDao.DeletePermission(deletePermissionReqDTO.ID); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to delete permission: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// RecoverPermission 根据id恢复权限
func (p *permissionService) RecoverPermission(recoverPermissionReqDTO dtos.RecoverPermissionReqDTO) error {
	if permission, err := p.permissionDao.GetDelPermissionByID(recoverPermissionReqDTO.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to recover permission, Permission not found: ", err)
			return errs.ErrPermissionNotFound
		}
		// 500 服务器内部错误
		global.Log.Errorln("Failed to recover permission, Failed to get permission: ", err)
		return errs.ErrInternalServerError
	} else if !permission.DeletedAt.Valid {
		// 409 请求冲突
		global.Log.Warnln("Failed to recover permission, Permission not deleted: ", err)
		return errs.ErrPermissionConflict
	}
	// 恢复权限
	if err := p.permissionDao.RecoverPermission(recoverPermissionReqDTO.ID); err != nil {
		// 500 服务器内部错误
		global.Log.Errorln("Failed to recover permission: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// checkDuplicatePermissionExists 检查具有相同名称的权限是否已存在
func (p *permissionService) checkDuplicatePermissionExists(permissionName string) error {
	_, err := p.permissionDao.GetPermissionByName(permissionName)
	// 如果发生错误且不是记录未找到的错误，则返回错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to check duplicate permission exists: ", err)
		return errs.ErrInternalServerError
	}
	// 409 冲突
	global.Log.Warnln("Permission name is already exists")
	return errs.ErrPermissionNameAlreadyExists
}

// checkDuplicatePermissionPathMethodExists 检查具有相同路径和方法的权限是否已存在
func (p *permissionService) checkDuplicatePermissionPathMethodExists(apiPath string, method string) error {
	_, err := p.permissionDao.GetPermissionByPathMethod(apiPath, method)
	// 如果发生错误且不是记录未找到的错误，则返回错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to check duplicate permission path method exists: ", err)
		return errs.ErrInternalServerError
	}
	// 409 冲突
	global.Log.Warnln("Api path and method is already exists")
	return errs.ErrPermissionPathMethodAlreadyExists
}
