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

// RolePermissionHandler 角色权限处理器
type RolePermissionHandler struct {
	RolePermissionService service.RolePermissionService
}

// NewRolePermissionHandler 创建角色权限处理器
func NewRolePermissionHandler(rolePermissionService service.RolePermissionService) *RolePermissionHandler {
	return &RolePermissionHandler{RolePermissionService: rolePermissionService}
}

// CreateRolePermissionsByIDs 创建角色权限
//
//	@Summary		Create RolePermissions
//	@Description	Create rolePermissions by roleIDs and permissionIDs
//	@Tags			RolePermission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	body		dtos.CreateRolePermissionsByIDsReqDTO	true	"Create rolePermissions request parameters"
//	@Success		200		{object}	response.Response{Data=interface{}}		"Successfully response with rolePermissions created message"
//	@Router			/role-permissions/ [post]
func (h *RolePermissionHandler) CreateRolePermissionsByIDs(c *gin.Context) {
	createRolePermissionsByIDsReqDTO := &dtos.CreateRolePermissionsByIDsReqDTO{}
	if err := c.ShouldBindJSON(createRolePermissionsByIDsReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to create role permission, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if nonExistingRoleIDs, nonExistingPermissionIDs, err := h.RolePermissionService.CreateRolePermissionsByIDs(*createRolePermissionsByIDsReqDTO); err != nil {
		response.FailWithData(c, utils.GetStatusCodeFromError(err), "Failed to create role permission: "+err.Error(), gin.H{
			"non_existing_role_ids":       nonExistingRoleIDs,
			"non_existing_permission_ids": nonExistingPermissionIDs,
		})
		return
	}
	response.OkWithMsg(c, http.StatusOK, "Role permission created successfully")
}

// DeleteRolePermissionsByIDs 删除角色权限
//
//	@Summary		Delete RolePermissions
//	@Description	Delete rolePermissions by roleIDs and permissionIDs
//	@Tags			RolePermission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	body		dtos.DeleteRolePermissionsByIDsReqDTO	true	"Delete rolePermissions request parameters"
//	@Success		200		{object}	response.Response{Data=interface{}}		"Successfully response with rolePermissions deleted message"
//	@Router			/role-permissions/ [delete]
func (h *RolePermissionHandler) DeleteRolePermissionsByIDs(c *gin.Context) {
	deleteRolePermissionsByIDsReqDTO := &dtos.DeleteRolePermissionsByIDsReqDTO{}
	if err := c.ShouldBindJSON(deleteRolePermissionsByIDsReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to delete role permission, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.RolePermissionService.DeleteRolePermissionsByIDs(*deleteRolePermissionsByIDsReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to delete role permission: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "Role permission deleted successfully")
}

// GetRolePermissionsByPermissionID 通过权限ID获取权限角色列表
//
//	@Summary		Get RolePermissions
//	@Description	Get rolePermissions by permissionID
//	@Tags			RolePermission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			permissionID	path		uint															true	"Get rolePermission by permissionID request parameters"
//	@Success		200				{object}	response.Response{Data=[]dtos.GetRolePermissionResDTOExample}	"Successfully response with rolePermission list information"
//	@Router			/role-permissions/permissions/{permissionID}/roles [get]
func (h *RolePermissionHandler) GetRolePermissionsByPermissionID(c *gin.Context) {
	getRolePermissionsByPermissionIDReqDTO := &dtos.GetRolePermissionsByPermissionIDReqDTO{}
	if err := c.ShouldBindUri(getRolePermissionsByPermissionIDReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get role permission list, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	rolePermissionList, err := h.RolePermissionService.GetRolePermissionsByPermissionID(*getRolePermissionsByPermissionIDReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get role permission list: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, rolePermissionList)
}

// GetRolePermissionsByRoleID 通过角色ID获取角色权限列表
//
//	@Summary		Get RolePermissions
//	@Description	Get rolePermissions by roleID
//	@Tags			RolePermission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			roleID	path		uint															true	"Get rolePermissions by roleID request parameters"
//	@Success		200		{object}	response.Response{Data=[]dtos.GetRolePermissionResDTOExample}	"Successfully response with rolePermission list information"
//	@Router			/role-permissions/roles/{roleID}/permissions [get]
func (h *RolePermissionHandler) GetRolePermissionsByRoleID(c *gin.Context) {
	getRolePermissionsByRoleIDReqDTO := &dtos.GetRolePermissionsByRoleIDReqDTO{}
	if err := c.ShouldBindUri(getRolePermissionsByRoleIDReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get role permission list, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	rolePermissionList, err := h.RolePermissionService.GetRolePermissionsByRoleID(*getRolePermissionsByRoleIDReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get role permission list: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, rolePermissionList)
}

// GetRolePermissionByID 通过角色ID和权限ID获取角色权限
//
//	@Summary		Get RolePermission
//	@Description	Get rolePermission by roleID and permissionID
//	@Tags			RolePermission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			roleID			path		uint														true	"Get rolePermission request parameters"
//	@Param			permissionID	path		uint														true	"Get rolePermission request parameters"
//	@Success		200				{object}	response.Response{Data=dtos.GetRolePermissionResDTOExample}	"Successfully response with rolePermission information"
//	@Router			/role-permissions/roles/{roleID}/permissions/{permissionID} [get]
func (h *RolePermissionHandler) GetRolePermissionByID(c *gin.Context) {
	var getRolePermissionReqDTO dtos.GetRolePermissionByIDReqDTO
	if err := c.ShouldBindUri(&getRolePermissionReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get role permission, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	RolePermission, err := h.RolePermissionService.GetRolePermissionByID(getRolePermissionReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get role permission: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, RolePermission)
}
