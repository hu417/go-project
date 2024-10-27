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

// PermissionHandler 权限处理器
type PermissionHandler struct {
	PermissionService service.PermissionService
}

// NewPermissionHandler 权限处理器
func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		PermissionService: permissionService,
	}
}

// GetPermissionList 获取权限列表
//	@Summary		Get Permission List
//	@Description	Get permissions list.
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	query		dtos.GetPermissionListReqDTO							false	"Get permission list request parameters"
//	@Success		200		{object}	response.Response{Data=dtos.PaginationResultExample}	"Successfully response with permission list information"
//	@Router			/permissions/ [get]
func (h *PermissionHandler) GetPermissionList(c *gin.Context) {
	var getPermissionListReqDTO dtos.GetPermissionListReqDTO
	if err := c.ShouldBindQuery(&getPermissionListReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get permission list, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	// 获取分页信息，如果未设置，则使用默认值
	getPermissionListReqDTO.PaginationReqDTO = dtos.GetDefaultPaginationDTO(getPermissionListReqDTO.PaginationReqDTO)

	permissions, err := h.PermissionService.GetPermissionList(getPermissionListReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get permission list: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, permissions)
}

// CreatePermission 创建权限
//	@Summary		Create Permission
//	@Description	Create permission.
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	body		dtos.CreatePermissionReqDTO			true	"Create permission request parameters"
//	@Success		201		{object}	response.Response{Data=interface{}}	"Successfully response with permission created message"
//	@Router			/permissions/ [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var createPermissionReqDTO dtos.CreatePermissionReqDTO
	if err := c.ShouldBindJSON(&createPermissionReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to create permission, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.PermissionService.CreatePermission(createPermissionReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to create permission: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusCreated, "Permission created successfully")
}

// GetPermission 获取权限
//	@Summary		Get Permission
//	@Description	Get permission by id.
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		uint													true	"Get permission request parameters"
//	@Success		200	{object}	response.Response{Data=dtos.GetPermissionDTOExample}	"Successfully response with permission information"
//	@Router			/permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	var getPermissionReqDTO dtos.GetPermissionReqDTO
	if err := c.ShouldBindUri(&getPermissionReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get permission, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	permission, err := h.PermissionService.GetPermission(getPermissionReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get permission: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, permission)
}

// UpdatePermission 更新权限
//	@Summary		Update Permission
//	@Description	Update permission by id.
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		uint								true	"Get permission request parameters"
//	@Param			data	body		dtos.UpdatePermissionReqDTO			true	"Update permission request parameters"
//	@Success		200		{object}	response.Response{Data=interface{}}	"Successfully response with permission updated message"
//	@Router			/permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	var getPermissionReqDTO dtos.GetPermissionReqDTO
	if err := c.ShouldBindUri(&getPermissionReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update permission, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	updatePermissionReqDTO := dtos.UpdatePermissionReqDTO{}
	if err := c.ShouldBindJSON(&updatePermissionReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update permission, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.PermissionService.UpdatePermission(getPermissionReqDTO, updatePermissionReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to update permission: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "Permission updated successfully")
}

// DeletePermission 删除权限
//	@Summary		Delete Permission
//	@Description	Delete permission by id.
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		uint								true	"Delete permission request parameters"
//	@Success		200	{object}	response.Response{Data=interface{}}	"Successfully response with permission deleted message"
//	@Router			/permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	var deletePermissionReqDTO dtos.DeletePermissionReqDTO
	if err := c.ShouldBindUri(&deletePermissionReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to delete permission, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.PermissionService.DeletePermission(deletePermissionReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to delete permission: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "Permission deleted successfully")
}

// RecoverPermission 恢复权限
//	@Summary		Recover Permission
//	@Description	Recover permission.
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		uint								true	"Recover permission request parameters"
//	@Success		200	{object}	response.Response{Data=interface{}}	"Successfully response with permission recovered message"
//	@Router			/permissions/{id}/recover [put]
func (h *PermissionHandler) RecoverPermission(c *gin.Context) {
	var recoverPermissionReqDTO dtos.RecoverPermissionReqDTO
	if err := c.ShouldBindUri(&recoverPermissionReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to recover permission, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.PermissionService.RecoverPermission(recoverPermissionReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to recover permission: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "Permission recovered successfully")
}
