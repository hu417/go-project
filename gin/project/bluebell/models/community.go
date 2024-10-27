package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Community struct {
	Id            int32  `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	CommunityID   string `json:"community_id"  gorm:"column:community_id;NOT NULL"`
	CommunityName string `json:"community_name" gorm:"column:community_name;NOT NULL"`
	Introduction  string `json:"introduction,omitempty" gorm:"column:introduction;NOT NULL"`
	Timestamps
}

func (c *Community) TableName() string {
	return "community"
}

func (c *Community) BeforeCreate(tx *gorm.DB) error {
	c.CommunityID = uuid.New().String()
	c.Timestamps.CreatedAt = time.Now().Unix()
	return nil
}
