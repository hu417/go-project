package models

import "time"

// 字段映射
type Share_basic struct {
	Id                       int
	Identity                 string    `xorm:"identity"`
	User_identity            string    `xorm:"user_identity"`
	User_Repository_Identity string    `xorm:"user_repository_identity"`
	Repository_identity      string    `xorm:"repository_identity"`
	Expired_time             int       `xorm:"expired_time" description:"失效时间"`
	Click_num                int       `xorm:"click_num" description:"点击次数"`
	CreatedAt                time.Time `xorm:"created" description:"创建时间"`
	UpdatedAt                time.Time `xorm:"updated" description:"更新时间"`
	DeletedAt                time.Time `xorm:"deleted" description:"删除时间"`
}

// 表明初始化
func (table *Share_basic) TableName() string {
	return "share_basic"
}
