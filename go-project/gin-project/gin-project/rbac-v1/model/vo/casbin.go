package vo

// 权限 //
type CasbinAuthRequest struct {
	Username string `json:"username" binding:"required"`
	Path     string `json:"path" binding:"required"`
	Method   string `json:"method" binding:"required"`
}
