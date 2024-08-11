package mysql

import (
	"log"
	"time"

	"api-demo/app/model"
	"api-demo/internal/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GetConnection() (db *gorm.DB) {
	// 获取数据库配置
	config := global.Config.Database

	// 拼装dsn
	dsn := config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Database + "?" + config.Config

	// 初始化配置
	mysqlConfig := mysql.Config{
		DSN: dsn,
	}

	// 连接数据库
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁用表名复数
		},
	})
	if err != nil {
		log.Println("数据库链接失败：", err)
		panic(err)
	}

	db.Logger = logger.Default.LogMode(logger.Info)
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("数据库链接失败：", err)
		panic(err)
	}
	// 最大空闲池
	sqlDB.SetMaxIdleConns(10)
	// 最大连接池
	sqlDB.SetMaxOpenConns(100)
	// 连接最大生命周期
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	db.Migrator().AutoMigrate(&model.Admin{})

	return
}
