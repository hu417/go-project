package vo

import (
	"rbac-v1/model/bo"
	"rbac-v1/model/po"
)

// 用户相关 //
type UserListRequest struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Name string `json:"name" form:"name"`
	Username string `json:"username" form:"username"`
}

type UserListResponse struct {
	Total int64 `json:"total"`
	List []*po.User `json:"list"`
}

type UserBoListResponse struct {
	Total int64 `json:"total"`
	List []*bo.User `json:"list"`
}

type UserCreateRequest struct {
	Users []*bo.UserCreate `json:"users" binding:"required"`
}

type UserDelRequest struct {
	Id uint `json:"id" binding:"required"`
}

type GetUserRoleAndPowerRequest struct {
	UserId uint `form:"user_id" `
}

type GetUserRoleAndPowerResponse struct {
	Roles []string `json:"roles"`
	Powers []string `json:"powers"`
}