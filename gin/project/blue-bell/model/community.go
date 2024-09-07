package model

type CommunityCategoryList struct {
	BaseID
	CommunityID int64  `json:"community_id" db:"community_id" gorm:"column:community_id;primaryKey;autoIncrement;comment:社区ID"`
	Name        string `json:"community_name" db:"community_name" gorm:"column:community_name;unique;comment:社区名称"`
	Timestamps
}

type CommunityCategoryDetail struct { 
	ID           int    `json:"community_id" db:"community_id" gorm:"column:community_id;primaryKey;autoIncrement;comment:社区ID"`
	Name         string `json:"community_name" db:"community_name" gorm:"column:community_name;unique;comment:社区名称"`
	Introduction string `json:"community_introduction" db:"introduction" gorm:"column:introduction;comment:社区介绍"`
}
