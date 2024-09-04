package migration

import (
	"bluebell/models/table"

	"gorm.io/gorm"
)

// 自动迁移数据库
func AutoMegration(db *gorm.DB) error{
	return db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		table.User{},
		table.Community{},
		table.Post{},
	)
}