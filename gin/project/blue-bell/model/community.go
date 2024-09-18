package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Community struct {
	BaseID
	CommunityID  string `json:"community_id" db:"community_id" gorm:"column:community_id;primaryKey;autoIncrement;comment:社区ID"`
	Name         string `json:"community_name" db:"community_name" gorm:"column:community_name;unique;comment:社区名称"`
	Introduction string `json:"community_introduction" db:"introduction" gorm:"column:introduction;comment:社区介绍"`
	Timestamps
}

// 表名
func (c *Community) TableName() string {
	return "community"
}

// 创建前
func (c *Community) BeforeCreate(tx *gorm.DB) error {
	c.CommunityID = uuid.New().String()
	c.Timestamps.CreatedAt = time.Now().Unix()
	return nil
}
