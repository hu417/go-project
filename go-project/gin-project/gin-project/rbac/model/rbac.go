package model

type CasbinPolicy struct {
	PType  string `json:"p_type" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
	Desc   string `json:"desc"  binding:"required"`
}
