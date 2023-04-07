package core

import (
	"fmt"
	"time"

	"gvb_server/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitGorm() *gorm.DB {

	if global.Config.Mysql.Host == "" {
		global.Logger.Warn("未配置MySQL,取消gorm连接")
		return nil
	}
	dsn := global.Config.Mysql.Dsn()

	// var mylogger logger.Interface
	// if global.Config.Mysql.LogMode == "error" {
	// 	mylogger = logger.Default.LogMode(logger.Error)

	// }
	// mylogger = logger.Default.LogMode(logger.Info)

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	Logger: mylogger,
	// })

	db, err := gorm.Open(mysql.Open(dsn), gormConfig())

	if err != nil {
		global.Logger.Error(fmt.Sprintf("mysql 连接失败: %v", dsn))
		panic(err)
	}
	global.Logger.Info(fmt.Sprintln("mysql 连接成功"))
	// 实例化db
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)               // 连接池中最大空闲连接数
	sqlDB.SetMaxOpenConns(100)              // 连接池中最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour * 4) // 单个连接最大存活时间(单位:秒)，不能超过mysql的wait_timeout

	return db
}

// 日志等配置
func gormConfig() *gorm.Config {
	return &gorm.Config{
		// gorm 日志模式
		Logger: logger.Default.LogMode(getLogMode(global.Config.Mysql.LogMode)),
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	}
}

// 根据字符串获取对应 LogLevel
func getLogMode(str string) logger.LogLevel {
	switch str {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}
