package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostContent 帖子内容
type Post struct {
	BaseID
	PostID      string  `json:"post_id" db:"post_id" gorm:"column:post_id;primaryKey;autoIncrement;"`
	Title       string `json:"title" db:"title" binding:"required" gorm:"column:title;index;size:128;not null;comment:帖子标题"`
	Content     string `json:"content" db:"content" binding:"required" gorm:"column:content;not null;comment:帖子内容"`
	AuthorID    string  `json:"author_id" db:"author_id" gorm:"column:author_id;index;not null;comment:作者ID"`
	CommunityID string  `json:"community_id" db:"community_id" binding:"required" gorm:"column:community_id;index;not null;comment:所属社区"`
	Timestamps
}

// 表名
func (p *Post) TableName() string {
	return "post"
}

// 创建前
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.PostID = uuid.New().String()
	p.Timestamps.CreatedAt = time.Now().Unix()
	return nil
}

