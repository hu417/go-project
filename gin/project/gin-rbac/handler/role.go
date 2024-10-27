package handler

import (
	"net/http"

	"gin-rbac/common/response"
	"gin-rbac/dtos"
	"gin-rbac/global"
	"gin-rbac/service"
	"gin-rbac/utils"

	"github.com/gin-gonic/gin"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	RoleService service.RoleService
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		RoleService: roleService,
	}
}

// GetRoleList 获取角色列表
//
//	@Summary		Get Role List
//	@Description	Get role list.
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			data	query	dtos.GetRoleListReqDTO	true	"Get role list request parameters"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.Response{Data=dtos.PaginationResultExample}	"Successfully response with role list information"
//	@Router			/roles/ [get]
func (h *RoleHandler) GetRoleList(c *gin.Context) {
	var getRoleListReqDTO dtos.GetRoleListReqDTO
	if err := c.ShouldBindQuery(&getRoleListReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get role list, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	// 获取分页信息，如果未设置，则使用默认值
	getRoleListReqDTO.PaginationReqDTO = dtos.GetDefaultPaginationDTO(getRoleListReqDTO.PaginationReqDTO)

	roleList, err := h.RoleService.GetRoleList(getRoleListReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get role list: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, roleList)
}

// CreateRole 创建角色
//
//	@Summary		Create Role
//	@Description	Create role.
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			data	body	dtos.CreateRoleReqDTO	true	"Create role request parameters"
//	@Security		ApiKeyAuth
//	@Success		201	{object}	response.Response{Data=interface{}}	"Successfully response with role created message"
//	@Router			/roles/ [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var createRoleReqDTO dtos.CreateRoleReqDTO
	if err := c.ShouldBindJSON(&createRoleReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to create role, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.RoleService.CreateRole(createRoleReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to create role: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusCreated, "Role created successfully")
}

// GetRole 获取角色
//
//	@Summary		Get Role
//	@Description	Get role by id.
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path	uint	true	"Get role request parameters"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.Response{Data=dtos.GetRoleDTOExample}	"Successfully response with role information"
//	@Router			/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	var getRoleReqDTO dtos.GetRoleReqDTO
	if err := c.ShouldBindUri(&getRoleReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get role, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	role, err := h.RoleService.GetRoleByID(getRoleReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get role: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, role)
}

// UpdateRole 更新角色
//
//	@Summary		Update Role
//	@Description	Update role by id.
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			id		path	uint					true	"Get Role Request parameters"
//	@Param			data	body	dtos.UpdateRoleReqDTO	true	"Update Role Request parameters"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.Response{Data=interface{}}	"Successfully response with role updated message"
//	@Router			/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var getRoleReqDTO dtos.GetRoleReqDTO
	if err := c.ShouldBindUri(&getRoleReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update role, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	var updateRoleReqDTO dtos.UpdateRoleReqDTO
	if err := c.ShouldBindJSON(&updateRoleReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update role, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.RoleService.UpdateRole(getRoleReqDTO, updateRoleReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to update role: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "Role updated successfully")
}

// DeleteRole 删除角色
//
//	@Summary		Delete Role
//	@Description	Delete role by id.
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path	uint	true	"Delete role request parameters"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.Response{Data=interface{}}	"Successfully response with role deleted message"
//	@Router			/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	var deleteRoleReqDTO dtos.DeleteRoleReqDTO
	if err := c.ShouldBindUri(&deleteRoleReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to delete role, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.RoleService.DeleteRole(deleteRoleReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Invalid request data: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "Role deleted successfully")
}

// RecoverRole 恢复角色
//
//	@Summary		Recover Role
//	@Description	Recover role by id.
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path	uint	true	"Recover role request parameters"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.Response{Data=interface{}}	"Successfully response with role recovered message"
//	@Router			/roles/{id}/recover [put]
func (h *RoleHandler) RecoverRole(c *gin.Context) {
	var recoverRoleReqDTO dtos.RecoverRoleReqDTO
	if err := c.ShouldBindUri(&recoverRoleReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to recover role, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.RoleService.RecoverRole(recoverRoleReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Invalid request data: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "Role recovered successfully")
}
