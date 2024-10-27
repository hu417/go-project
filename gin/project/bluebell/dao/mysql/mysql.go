package mysql

import (
	"fmt"

	"bluebell/models"
	"bluebell/setting"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Init 初始化MySQL连接
func Init(cfg *setting.MySQLConfig) (db *gorm.DB) {

	if cfg.Host == "" {
		return nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
		// Logger:                                   utils.GetGormLogger(), // 使用自定义 Logger
	})
	if err != nil {
		zap.L().Error("mysql connect failed, err:", zap.Any("err", err))
		return nil
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	Migrating(db)
	return db
}

func Migrating(db *gorm.DB) {
	// 创建表时添加后缀
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.User{},
		&models.Community{},
		&models.Post{},
	)

}
