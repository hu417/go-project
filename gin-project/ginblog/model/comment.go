package model

import "gorm.io/gorm"

// Comment 评论结构体
type Comment struct {
	gorm.Model
	UserId    uint   `json:"user_id"`
	ArticleId uint   `json:"article_id"`
	Title     string `gorm:"type:varchar(500);not null;" json:"article_title"`
	Username  string `gorm:"type:varchar(500);not null;" json:"username"`
	Content   string `gorm:"type:varchar(500);not null;" json:"content"`
	Status    int8   `gorm:"type:tinyint;default:2" json:"status"`
}
