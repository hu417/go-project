package po

//角色-权限中间表（多对多）
type RolePo struct {
	RoleId uint `json:"role_id" gorm:"primaryKey"`
	PowerId uint `json:"power_id" gorm:"primaryKey"`
}

func (*RolePo) TableName() string {
	return "rbac_role_power"
}
