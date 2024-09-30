package model

import "gorm.io/gorm"

type SysUser struct {
	gorm.Model
	Nickname  string `gorm:"column:nickname;type:varchar(50);" json:"nickname"`
	Username  string `gorm:"column:username;type:varchar(50);" json:"username"`
	Password  string `gorm:"column:password;type:varchar(36);" json:"password"`
	Phone     string `gorm:"column:phone;type:varchar(20);" json:"phone"`
	Sex       string `gorm:"column:sex;type:varchar(20);" json:"sex"`
	Email     string `gorm:"column:email;type:varchar(20);" json:"email"`
	WxUnionId string `gorm:"column:wx_union_id;type:varchar(255);" json:"wx_union_id"`
	WxOpenId  string `gorm:"column:wx_open_id;type:varchar(255);" json:"wx_open_id"`
	Avatar    string `gorm:"column:avatar;type:varchar(255);" json:"avatar"`   // 头像
	Remarks   string `gorm:"column:remarks;type:varchar(255);" json:"remarks"` // 备注
	RoleId    uint   `gorm:"column:role_id;type:bigint(20);" json:"roleId"`    // 角色ID

}

// TableName 设置用户表名称
func (table *SysUser) TableName() string {
	return "sys_user"
}

