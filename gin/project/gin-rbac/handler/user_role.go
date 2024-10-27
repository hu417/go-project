package handler

import (
	"fmt"
	"net/http"

	"gin-rbac/common/response"
	"gin-rbac/dtos"
	"gin-rbac/global"
	"gin-rbac/service"
	"gin-rbac/utils"

	"github.com/gin-gonic/gin"
)

// UserRoleHandler 用户角色处理器
type UserRoleHandler struct {
	UserRoleService service.UserRoleService
}

// NewUserRoleHandler 创建用户角色处理器
func NewUserRoleHandler(userRoleService service.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{UserRoleService: userRoleService}
}

// CreateUserRolesByIDs 创建用户角色
//
//	@Summary		Create UserRoles
//	@Description	Create userRoles by userIDs and roleIDs
//	@Tags			UserRole
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	body		dtos.CreateUserRolesByIDsReqDTO		true	"Create userRoles request parameters"
//	@Success		200		{object}	response.Response{Data=interface{}}	"Successful response with userRoles message"
//	@Router			/user-roles/ [post]
func (h *UserRoleHandler) CreateUserRolesByIDs(c *gin.Context) {
	createUserRolesByIDsReqDTO := &dtos.CreateUserRolesByIDsReqDTO{}
	if err := c.ShouldBindJSON(createUserRolesByIDsReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to create user role, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if nonExistingUserIDs, nonExistingRoleIDs, err := h.UserRoleService.CreateUserRolesByIDs(*createUserRolesByIDsReqDTO); err != nil {
		response.FailWithData(c, utils.GetStatusCodeFromError(err), "Failed to create user role: "+err.Error(), gin.H{
			"non_existing_user_ids": nonExistingUserIDs,
			"non_existing_role_ids": nonExistingRoleIDs,
		})
		return
	}
	response.OkWithMsg(c, http.StatusCreated, "User role created successfully")
}

// DeleteUserRolesByIDs 删除用户角色
//
//	@Summary		Delete UserRoles
//	@Description	Delete userRoles by userIDs and roleIDs
//	@Tags			UserRole
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	body		dtos.DeleteUserRolesByIDsReqDTO		true	"Delete userRoles request parameters"
//	@Success		200		{object}	response.Response{Data=interface{}}	"Successful response with userRoles message"
//	@Router			/user-roles/ [delete]
func (h *UserRoleHandler) DeleteUserRolesByIDs(c *gin.Context) {
	deleteUserRolesByIDsReqDTO := &dtos.DeleteUserRolesByIDsReqDTO{}
	if err := c.ShouldBindJSON(deleteUserRolesByIDsReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to delete user role, Invalid request data: %v", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.UserRoleService.DeleteUserRolesByIDs(*deleteUserRolesByIDsReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Invalid request data: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "User role deleted successfully")
}

// GetUserRolesByUserID 通过用户id获取用户角色列表
//
//	@Summary		Get UserRoles
//	@Description	Get userRoles by userID
//	@Tags			UserRole
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userID	path		uint														true	"Get UserRoles By UserID Request parameters"
//	@Success		200		{object}	response.Response{Data=[]dtos.GetUserRoleByIDResDTOExample}	"Successful response with UserRoles list information"
//	@Router			/user-roles/user/{userID} [get]
func (h *UserRoleHandler) GetUserRolesByUserID(c *gin.Context) {
	getUserRolesByUserIDReqDTO := &dtos.GetUserRolesByUserIDReqDTO{}
	if err := c.ShouldBindUri(getUserRolesByUserIDReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get user role by user id, Invalid request data: %v", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	userRoleList, err := h.UserRoleService.GetUserRolesByUserID(*getUserRolesByUserIDReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get user role by user id: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, userRoleList)
}

// GetUserRolesByRoleID 通过角色id获取用户角色列表
//
//	@Summary		Get UserRoles
//	@Description	Get userRoles by roleID
//	@Tags			UserRole
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			roleID	path		uint														true	"Get userRoles by roleID request parameters"
//	@Success		200		{object}	response.Response{Data=[]dtos.GetUserRoleByIDResDTOExample}	"Successful response with userRoles list information"
//	@Router			/user-roles/role/{roleID} [get]
func (h *UserRoleHandler) GetUserRolesByRoleID(c *gin.Context) {
	getUserRolesByRoleIDReqDTO := &dtos.GetUserRolesByRoleIDReqDTO{}
	if err := c.ShouldBindUri(getUserRolesByRoleIDReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get user role by role id, Invalid request data: %v", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	userRoleList, err := h.UserRoleService.GetUserRolesByRoleID(*getUserRolesByRoleIDReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get user role by role id: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, userRoleList)
}

// GetUserRoleByID 获取用户角色
//
//	@Summary		Get UserRole
//	@Description	Get userRole by userID and roleID
//	@Tags			UserRole
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userID	path		uint														true	"Get userRole request parameters"
//	@Param			roleID	path		uint														true	"Get userRole request parameters"
//	@Success		200		{object}	response.Response{Data=dtos.GetUserRoleByIDResDTOExample}	"Successful response with user-role information"
//	@Router			/user-roles/{userID}/{roleID} [get]
func (h *UserRoleHandler) GetUserRoleByID(c *gin.Context) {
	var getUserRoleByIDReqDTO dtos.GetUserRoleByIDReqDTO
	if err := c.ShouldBindUri(&getUserRoleByIDReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get user role by id, Invalid request data:%v", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	fmt.Printf("userid: %+v, roleid:%+v\n", getUserRoleByIDReqDTO.UserID, getUserRoleByIDReqDTO.RoleID)
	if getUserRoleByIDReqDTO.UserID == 0 || getUserRoleByIDReqDTO.RoleID == 0 {
		response.Fail(c, http.StatusBadRequest, "Invalid request data: user id or role id cannot be empty")
		return
	}

	userRole, err := h.UserRoleService.GetUserRoleByID(getUserRoleByIDReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get user role by id: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, userRole)
}
