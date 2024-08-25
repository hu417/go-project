package v1

// @Title  params.go
// @Description  用户模块参数
type UserParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}
