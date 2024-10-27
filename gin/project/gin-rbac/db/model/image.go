package model

import "time"

// ImageModel 图片表模型
type ImageModel struct {
	ID        uint      `gorm:"primaryKey"`
	ImageName string    `gorm:"type:varchar(255);not null;comment:图片文件名"`  // 图片文件名
	ImagePath string    `gorm:"type:varchar(255);not null;comment:图片存储路径"` // 图片存储路径
	CreatedAt time.Time `gorm:"comment:图片创建时间"`                            // 创建时间
}

func (ImageModel) TableName() string {
	return "b_image"
}
