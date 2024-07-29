package mysql

import (
	"bluebell/config"
	"bluebell/dao/mysql/migration"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

// Init 初始化MySQL连接
func Init(cfg *config.MySQLConfig) (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	// 设置GORM日志级别为Info，将SQL语句打印到控制台
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的地方，这里是标准输出）
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 彩色打印
		},
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   false,     //设置为 false 表示不跳过默认事务
		NamingStrategy:                           nil,       //设置为 nil，表示使用 GORM 的默认命名策略
		FullSaveAssociations:                     false,     //设置为 false 表示不完整保存关联的记录
		Logger:                                   newLogger, //设置为 nil，表示不设置自定义的日志记录器
		NowFunc:                                  nil,       //设置为 nil，表示使用 GORM 的默认时间函数
		DryRun:                                   false,     //设置为 false，表示不执行模拟运行
		PrepareStmt:                              false,     //设置为 false，表示不预先准备语句
		DisableAutomaticPing:                     false,     //设置为 false，表示不禁用自动心跳检测
		DisableForeignKeyConstraintWhenMigrating: false,     //设置为 false，表示在迁移时不禁用外键约束
		IgnoreRelationshipsWhenMigrating:         false,     //设置为 false，表示在迁移时不忽略关系
		DisableNestedTransaction:                 false,     //设置为 false，表示不禁用嵌套事务
		AllowGlobalUpdate:                        false,     //设置为 false，表示不允许全局更新
		QueryFields:                              false,     //设置为 false，表示不使用字段查询
		CreateBatchSize:                          0,         //设置为 0，表示使用 GORM 的默认批量创建大小
		TranslateError:                           false,     //设置为 false，表示不转换错误
		ClauseBuilders:                           nil,       //设置为 nil，表示不设置自定义的查询子句构建器
		ConnPool:                                 nil,       //设置为 nil，表示使用 GORM 的默认连接池
		Dialector:                                nil,       //设置为 nil，表示使用 GORM 的默认方言（dialect）
		Plugins:                                  nil,       //设置为 nil，表示不使用任何插件
	})
	if err != nil {
		panic("数据库连接错误")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接池失败")
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = migration.AutoMegration(db); err != nil {
		return
	}
	return
}

// Close 关闭MySQL连接
func Close() {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
		return
	}

	panic("关闭MySQL连接失败")
}

func GetDB() *gorm.DB {
	return db
}
