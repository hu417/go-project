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

// UserHandler 用户处理程序
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理程序
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Login 登录用户
//
//	@Summary		Login
//	@Description	Logs a user in and returns an authentication token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		dtos.LoginUserReqDTO						true	"Login request parameters"
//	@Success		200		{object}	response.Response{Data=dtos.TokenResDTO}	"Successful response with authentication token"
//	@Router			/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var loginUserReqDTO dtos.LoginUserReqDTO
	if err := c.ShouldBindJSON(&loginUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to login user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	token, err := h.userService.LoginUser(loginUserReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to login user: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, token)
}

// CreateUser 注册用户
//
//	@Summary		Register
//	@Description	Registers a new user and returns an authentication token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		dtos.CreateUserReqDTO						true	"Register request parameters"
//	@Success		200		{object}	response.Response{Data=dtos.TokenResDTO}	"Successful response with authentication token"
//	@Router			/auth/register [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var createUserReqDTO dtos.CreateUserReqDTO
	if err := c.ShouldBindJSON(&createUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to create user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	token, err := h.userService.CreateUser(createUserReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to create user: "+err.Error())
		return
	}
	response.Ok(c, http.StatusCreated, "User created successfully", token)
}

// GetUserByJWT 获取用户
//
//	@Summary		Get User
//	@Description	Gets the authenticated user's information.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.Response{Data=dtos.GetFullUserResDTO}	"Successful response with user information"
//	@Router			/auth/user [get]
func (h *UserHandler) GetUserByJWT(c *gin.Context) {
	// 从JWT中获取用户ID
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Fail(c, http.StatusUnauthorized, "Failed to get user: User not found")
		return
	}
	var getFullUserReqDTO dtos.GetFullUserReqDTO
	getFullUserReqDTO.ID = userID
	if err := c.ShouldBindUri(&getFullUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	user, err := h.userService.GetUserByID(getFullUserReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get user: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, user)
}

// UpdateUserByJWT 更新用户
//
//	@Summary		Update User
//	@Description	Updates the authenticated user's information.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	body		dtos.UpdateUserReqDTO				true	"Update user request parameters"
//	@Success		200		{object}	response.Response{Data=interface{}}	"Successful response with user updated message"
//	@Router			/auth/user [put]
func (h *UserHandler) UpdateUserByJWT(c *gin.Context) {
	// 从JWT中获取用户ID
	userID := c.GetUint("userID")
	var getFullUserReqDTO dtos.GetFullUserReqDTO
	getFullUserReqDTO.ID = userID
	if err := c.ShouldBindUri(&getFullUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	var updateUserReqDTO dtos.UpdateUserReqDTO
	if err := c.ShouldBindJSON(&updateUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.userService.UpdateUser(getFullUserReqDTO, updateUserReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to update user: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "User updated successfully")
}

// DeleteUserByJWT 注销用户
//
//	@Summary		Delete User
//	@Description	Deletes the authenticated user.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.Response{Data=interface{}}	"Successful response with user deleted message"
//	@Router			/auth/user [delete]
func (h *UserHandler) DeleteUserByJWT(c *gin.Context) {
	// 从JWT中获取用户ID
	userID := c.GetUint("userID")
	var deleteUserReqDTO dtos.DeleteUserReqDTO
	deleteUserReqDTO.ID = userID
	if err := c.ShouldBindUri(&deleteUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to delete user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.userService.DeleteUser(deleteUserReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to delete user: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "User deleted successfully")
}

// UpdateUserAvatarByJWT 更新用户头像
//
//	@Summary		Update Avatar
//	@Description	Updates the authenticated user's avatar.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	body		dtos.UpdateUserAvatarReqDTO			true	"Update avatar request parameters"
//	@Success		200		{object}	response.Response{Data=interface{}}	"Successful response with user avatar updated message"
//	@Router			/auth/user/avatar [put]
func (h *UserHandler) UpdateUserAvatarByJWT(c *gin.Context) {
	// 从JWT中获取用户ID
	userID := c.GetUint("userID")
	var getFullUserReqDTO dtos.GetFullUserReqDTO
	getFullUserReqDTO.ID = userID
	if err := c.ShouldBindUri(&getFullUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update user avatar, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	var updateUserAvatarReqDTO dtos.UpdateUserAvatarReqDTO
	if err := c.ShouldBindJSON(&updateUserAvatarReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update user avatar, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	if err := h.userService.UpdateUserAvatar(getFullUserReqDTO, updateUserAvatarReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to update user avatar: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "User avatar updated successfully")
}

// UpdatePasswordByJWT 更新用户密码
//
//	@Summary		Update Password
//	@Description	Updates the authenticated user's password.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	body		dtos.UpdatePasswordReqDTO			true	"Update password request parameters"
//	@Success		200		{object}	response.Response{Data=interface{}}	"Successful response with user password updated message"
//	@Router			/auth/user/password [put]
func (h *UserHandler) UpdatePasswordByJWT(c *gin.Context) {
	// 从JWT中获取用户ID
	userID := c.GetUint("userID")
	var getFullUserReqDTO dtos.GetFullUserReqDTO
	getFullUserReqDTO.ID = userID
	if err := c.ShouldBindUri(&getFullUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update user password, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	var updateUserPasswordReqDTO dtos.UpdatePasswordReqDTO
	if err := c.ShouldBindJSON(&updateUserPasswordReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to update user password, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.userService.UpdatePassword(getFullUserReqDTO, updateUserPasswordReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to update user password: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "User password updated successfully")
}

// GetPublicUserList 获取用户公开信息列表
//
//	@Summary		Get Public User List
//	@Description	Get public user list.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			data	query		dtos.PublicUserListReqDTO								true	"Public user list request parameters"
//	@Success		200		{object}	response.Response{Data=dtos.PaginationResultExample}	"Successful response with public user list information"
//	@Router			/users/public [get]
func (h *UserHandler) GetPublicUserList(c *gin.Context) {
	var publicUserListReqDTO dtos.PublicUserListReqDTO
	if err := c.ShouldBindQuery(&publicUserListReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get public user list, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	// 获取分页信息，如果未设置，则使用默认值
	publicUserListReqDTO.PaginationReqDTO = dtos.GetDefaultPaginationDTO(publicUserListReqDTO.PaginationReqDTO)

	users, err := h.userService.GetPublicUserList(publicUserListReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get user list: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, users)
}

// GetPublicUser 获取用户公开信息
//
//	@Summary		Get Public User
//	@Description	Get public user by id.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint											true	"Public user info request parameters"
//	@Success		200	{object}	response.Response{Data=dtos.PublicUserResDTO}	"Successful response with public user information"
//	@Router			/users/{id}/public [get]
func (h *UserHandler) GetPublicUser(c *gin.Context) {
	var publicUserReqDTO dtos.PublicUserReqDTO
	if err := c.ShouldBindUri(&publicUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get public user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	user, err := h.userService.GetPublicUserByID(publicUserReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get public user: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, user)
}

// GetUserByID 获取用户
//
//	@Summary		Get User
//	@Description	Get user by id.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		uint											true	"Get full user request parameters"
//	@Success		200	{object}	response.Response{Data=dtos.GetFullUserResDTO}	"Successful response with user information"
//	@Router			/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	var getFullUserReqDTO dtos.GetFullUserReqDTO
	if err := c.ShouldBindUri(&getFullUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to get user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	user, err := h.userService.GetUserByID(getFullUserReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to get user: "+err.Error())
		return
	}
	response.OkWithData(c, http.StatusOK, user)
}

// ResetPassword 重置密码
//
//	@Summary		Reset Password
//	@Description	Reset user's password.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			data	body		dtos.ResetPasswordReqDTO			true	"Reset user password request parameters"
//	@Success		200		{object}	response.Response{Data=interface{}}	"Successful response with user password reset message"
//	@Router			/users/reset_password [put]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var resetPasswordReqDTO dtos.ResetPasswordReqDTO
	if err := c.ShouldBindJSON(&resetPasswordReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to reset user password, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.userService.ResetPassword(resetPasswordReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to reset user password: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "User password reset successfully")
}

// DeleteUser 删除用户
//
//	@Summary		Delete User
//	@Description	Delete user by id.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		uint								true	"Delete user request parameters"
//	@Success		200	{object}	response.Response{Data=interface{}}	"Successful response with user deleted message"
//	@Router			/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	var deleteUserReqDTO dtos.DeleteUserReqDTO
	if err := c.ShouldBindUri(&deleteUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to delete user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.userService.DeleteUser(deleteUserReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to delete user: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "User deleted successfully")
}

// RecoverUser 恢复用户
//
//	@Summary		Recover User
//	@Description	Recover user by id.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		uint								true	"Recover user request parameters"
//	@Success		200	{object}	response.Response{Data=interface{}}	"Successful response with user recover message"
//	@Router			/users/{id}/recover [put]
func (h *UserHandler) RecoverUser(c *gin.Context) {
	var recoverUserReqDTO dtos.RecoverUserReqDTO
	if err := c.ShouldBindUri(&recoverUserReqDTO); err != nil {
		// 400 请求参数错误
		global.Log.Warnln("Failed to recover user, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.userService.RecoverUser(recoverUserReqDTO); err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to recover user: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "User recovered successfully")
}
