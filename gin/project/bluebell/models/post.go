package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 内存对齐概念
type Post struct {
	Id          int64  `json:"id,string" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	PostID      string `json:"post_id,string" gorm:"column:post_id;NOT NULL;comment:'帖子id'"`
	Title       string `json:"title" gorm:"column:title;NOT NULL;comment:'标题'"`
	Content     string `json:"content" gorm:"column:content;NOT NULL;comment:'内容'"`
	AuthorID    string `json:"author_id" gorm:"column:author_id;NOT NULL;comment:'作者的用户id'"`
	CommunityID string `json:"community_id" gorm:"column:community_id;NOT NULL;comment:'所属社区'"`
	Status      int8   `json:"status" gorm:"column:status;default:1;NOT NULL;comment:'帖子状态'"`
	Timestamps
}

func (p *Post) TableName() string {
	return "post"
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.PostID = uuid.New().String()
	p.Timestamps.CreatedAt = time.Now().Unix()
	return nil
}
