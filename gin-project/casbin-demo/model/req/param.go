package req

// 添加用户角色
type UserRole struct {
	UserName string `json:"user_name" banding:"required"`
	RoleName string `json:"role_name" banding:"required"`
}

// 更新用户角色
type UpdateUserRoleStr struct {
	UserName string `json:"user_name"`
	OldRole  string `json:"old_role"`
	NewRole  string `json:"new_role"`
}

// (RoleName, Url, Method) 对应于 `CasbinRule` 表中的 (v0, v1, v2)
type RolePolicy struct {
	RoleName string `gorm:"column:v0" description:"角色" json:"role_name" banding:"required"`
	Url      string `gorm:"column:v1" description:"api路径" json:"url" banding:"required"`
	Method   string `gorm:"column:v2" description:"访问方法" json:"method" banding:"required"`
}

// 更新角色权限
type RolePolicys struct {
	OldRoleName string `json:"old_role_name"`
	OldUrl      string `json:"old_url"`
	OldMethod   string `json:"old_method"`
	NewRoleName string `json:"new_role_name"`
	NewUrl      string `json:"new_url"`
	NewMethod   string `json:"new_method"`
}
