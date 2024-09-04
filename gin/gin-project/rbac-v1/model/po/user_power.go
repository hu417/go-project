package po

type UserPo struct {
	UserId uint `json:"user_id" gorm:"primaryKey"`
	PowerId uint `json:"power_id" gorm:"primaryKey"`
}

func (*UserPo) TableName() string {
	return "rbac_user_power"
}