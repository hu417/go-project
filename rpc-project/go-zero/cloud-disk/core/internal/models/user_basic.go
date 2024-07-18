package models

import "time"

// 字段映射
type User_basic struct {
	Id        int
	Identity  string
	Name      string
	Password  string
	Email     string
	CreatedAt time.Time `xorm:"created" description:"创建时间"`
	UpdatedAt time.Time `xorm:"updated" description:"更新时间"`
	DeletedAt time.Time `xorm:"deleted" description:"删除时间"`
}

// 表明初始化
func (table *User_basic) TableName() string {
	return "user_basic"
}
