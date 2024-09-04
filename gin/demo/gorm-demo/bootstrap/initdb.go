package bootstrap

import (
	"gorm-demo/model"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (

	// 单例工具
	dbOnce sync.Once
)

func InitMysql() (db *gorm.DB, err error) {

	dbOnce.Do(func() {
		// 连接mysql的dsn
		dsn := "root:123456@(127.0.0.1:3306)/gorm-demo?timeout=5000ms&readTimeout=5000ms&writeTimeout=5000ms&charset=utf8mb4&parseTime=true&loc=Local"

		// 创建db实例
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	})

	db.Migrator().AutoMigrate(&model.User{})

	return
}
