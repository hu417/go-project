package handler

import (
	"net/http"

	"go-test/internal/service"
	"go-test/pkg/errs"
	"go-test/pkg/request"
	"go-test/pkg/responce"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

// Register godoc
//	@Summary	用户注册
//	@Schemes
//	@Description	目前只支持邮箱登录
//	@Tags			用户模块
//	@Accept			json
//	@Produce		json
//	@Param			request	body		v1.RegisterRequest	true	"params"
//	@Success		200		{object}	v1.Response
//	@Router			/register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(req); err != nil {
		responce.HandleError(ctx, http.StatusBadRequest, errs.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, &req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		responce.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	responce.HandleSuccess(ctx, nil)
}

// Login godoc
//	@Summary	账号登录
//	@Schemes
//	@Description
//	@Tags		用户模块
//	@Accept		json
//	@Produce	json
//	@Param		request	body		v1.LoginRequest	true	"params"
//	@Success	200		{object}	v1.LoginResponse
//	@Router		/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responce.HandleError(ctx, http.StatusBadRequest, errs.ErrBadRequest, nil)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		responce.HandleError(ctx, http.StatusUnauthorized, errs.ErrUnauthorized, nil)
		return
	}
	responce.HandleSuccess(ctx, request.LoginResponseData{
		AccessToken: token,
	})
}

// GetProfile godoc
//	@Summary	获取用户信息
//	@Schemes
//	@Description
//	@Tags		用户模块
//	@Accept		json
//	@Produce	json
//	@Security	Bearer
//	@Success	200	{object}	v1.GetProfileResponse
//	@Router		/user [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		responce.HandleError(ctx, http.StatusUnauthorized, errs.ErrUnauthorized, nil)
		return
	}

	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		responce.HandleError(ctx, http.StatusBadRequest, errs.ErrBadRequest, nil)
		return
	}

	responce.HandleSuccess(ctx, user)
}

// UpdateProfile godoc
//	@Summary	修改用户信息
//	@Schemes
//	@Description
//	@Tags		用户模块
//	@Accept		json
//	@Produce	json
//	@Security	Bearer
//	@Param		request	body		v1.UpdateProfileRequest	true	"params"
//	@Success	200		{object}	v1.Response
//	@Router		/user [put]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)

	var req request.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responce.HandleError(ctx, http.StatusBadRequest, errs.ErrBadRequest, nil)
		return
	}

	if err := h.userService.UpdateProfile(ctx, userId, &req); err != nil {
		responce.HandleError(ctx, http.StatusInternalServerError, errs.ErrInternalServerError, nil)
		return
	}

	responce.HandleSuccess(ctx, nil)
}
