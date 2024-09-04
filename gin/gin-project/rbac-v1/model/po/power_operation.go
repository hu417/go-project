package po

//权限-行为中间表（多对多）
type PowerOp struct {
	PowerId uint `json:"power_id" gorm:"primaryKey"`
	OperationId uint `json:"operation_id" gorm:"primaryKey"`
}

func (*PowerOp) TableName() string {
	return "rbac_power_op"
}
