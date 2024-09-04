package table

import (
	"bluebell/pkg/utils"
	"time"
)


type Community struct {
	Id            int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	CommunityId   int64      `gorm:"column:community_id;type:int(10) unsigned;NOT NULL" json:"community_id"`
	CommunityName string    `gorm:"column:community_name;type:varchar(128);NOT NULL" json:"community_name"`
	Introduction  string    `gorm:"column:introduction;type:varchar(256);NOT NULL" json:"introduction"`
	CreateTime    time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"`
	
}

func (m *Community) TableName() string {
	return "community"
}

// 随机id
func (t *Community) BeforCreate() {
	t.CommunityId = int64(utils.GetUuidInt())
}	