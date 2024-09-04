package vo

import (
	"rbac-v1/model/bo"
	"rbac-v1/model/po"
)

// 列表（分页） //
type PowerListRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Name     string `json:"name" form:"name"`
	Code     string `json:"code" form:"code"`
}

type PowerListResponse struct {
	Total int64       `json:"total"`
	List  []*po.Power `json:"list"`
}

type PowerBoListResponse struct {
	Total int64       `json:"total"`
	List  []*bo.Power `json:"list"`
}

type PowerCreateRequest struct {
	Powers []*bo.Power `json:"powers" binding:"required"`
}

type PowerDelRequest struct {
	Id uint `json:"id" binding:"required"`
}
