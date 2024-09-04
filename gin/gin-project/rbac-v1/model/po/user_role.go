package po

type UserRo struct {
	UserId uint `json:"user_id" gorm:"primaryKey"`
	RoleId uint `json:"role_id" gorm:"primaryKey"`
}

func (*UserRo) TableName() string {
	return "rbac_user_role"
}