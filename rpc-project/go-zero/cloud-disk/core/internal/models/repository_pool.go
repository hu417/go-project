package models

import "time"

// 字段映射
type Repository_pool struct {
	Id        int
	Identity  string
	Hash      string
	Name      string
	Ext       string
	Size      int
	Path      string
	CreatedAt time.Time `xorm:"created" description:"创建时间"`
	UpdatedAt time.Time `xorm:"updated" description:"更新时间"`
	DeletedAt time.Time `xorm:"deleted" description:"删除时间"`
}

// 表明初始化
func (table *Repository_pool) TableName() string {
	return "repository_pool"
}
