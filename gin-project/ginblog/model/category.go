package model

// Category 文章分类模型
type Category struct {
	ID   uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}
