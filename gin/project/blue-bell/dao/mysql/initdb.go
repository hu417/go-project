package mysql

import (
	"strconv"

	"blue-bell/global"
	"blue-bell/model"
	"blue-bell/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDB() *gorm.DB {
	// 根据驱动配置进行初始化
	switch global.Conf.Database.Driver {
	case "mysql":
		zap.L().Debug("mysql driver is mysql")
		return initMySqlGorm()
	default:
		return initMySqlGorm()
	}
}

// 初始化 mysql gorm.DB
func initMySqlGorm() *gorm.DB {
	dbConfig := global.Conf.Database

	if dbConfig.DBName == "" {
		return nil
	}
	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" +
		dbConfig.DBName + "?" + dbConfig.Charset

	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,                   // 禁用自动创建外键约束
		Logger:                                   logger.GetGormLogger(), // 使用自定义 Logger
	})
	if err != nil {
		zap.L().Error("mysql connect failed, err:", zap.Any("err", err))
		return nil
	}

	// 连接池
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	// 自动迁移
	autoMigrate(db)

	return db

}

func autoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		model.User{},
		model.Post{},
		model.Community{},
		
	); err != nil {
		panic("migrate failed, err:" + err.Error())
	}

}
