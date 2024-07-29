package vo

import (
	"rbac-v1/model/bo"
	"rbac-v1/model/po"
)

// 角色相关 //
type RoleListRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Name     string `json:"name" form:"name"`
	Code     string `json:"code" form:"code"`
	Ids      []uint `json:"ids" form:"ids"`
}

type RoleListResponse struct {
	Total int64      `json:"total"`
	List  []*po.Role `json:"list"`
}

type RoleBoListResponse struct {
	Total int64      `json:"total"`
	List  []*bo.Role `json:"list"`
}

type RoleCreateRequest struct {
	Roles []*bo.Role `json:"roles" binding:"required"`
}

type RoleDelRequest struct {
	Id uint `json:"id" binding:"required"`
}
