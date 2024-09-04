package vo

import "rbac-v1/model/po"

type OperationListRequest struct {
	Ids []uint `json:"ids" form:"ids"`
	Page int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Path string `json:"path" form:"path"`
	Method string `json:"method" form:"method"`
	Type int `json:"type" form:"type"`
}

type OperationListResponse struct {
	Total int64 `json:"total"`
	List []*po.Operation `json:"list"`
}

type OperationCreateRequest struct {
	Operations []*po.Operation `json:"operations" binding:"required"`
}

type OperationDelRequest struct {
	Id uint `json:"id" binding:"required"`
}