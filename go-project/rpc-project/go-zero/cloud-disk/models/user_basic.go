package models

// 字段映射
type User_basic struct {
	Id       int32
	Identity string
	Name     string
	Password string
	Email    string
}

// 表明初始化
func (table *User_basic) TableName() string {
	return "user_basic"
}
