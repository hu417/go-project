package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"gin-base/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// GormMysql 初始化Mysql数据库
func InitMysql() *gorm.DB {
	m := global.Conf.MySQL
	if m.Dbname == "" {
		return nil
	}

	// 创建 mysql.Config 实例，其中包含了连接数据库所需的信息，比如 DSN (数据源名称)，字符串类型字段的默认长度以及自动根据版本进行初始化等参数。
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}

	// 打开数据库连接
	db, err := gorm.Open(mysql.New(mysqlConfig), Config(m.Prefix, m.Singular, m.LogMode))

	// 将引擎设置为我们配置的引擎，并设置每个连接的最大空闲数和最大连接数。
	if err != nil {
		panic(err)
	}

	db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)

	fmt.Println("====3-gorm====: gorm link mysql success")
	return db

}

// Config gorm 自定义配置
func Config(prefix string, singular bool, logmodel string) *gorm.Config {

	// 将传入的字符串前缀和单复数形式参数应用到 GORM 的命名策略中，并禁用迁移过程中的外键约束，返回最终生成的 GORM 配置信息。
	config := &gorm.Config{
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,   // 表前缀，在表名前添加前缀，如添加用户模块的表前缀 user_
			SingularTable: singular, // 是否使用单数形式的表名，如果设置为 true，那么 User 模型会对应 users 表
		},
		// 是否在迁移时禁用外键约束，默认为 false，表示会根据模型之间的关联自动生成外键约束语句
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	_default := logger.New(logger.Writer(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Info,
		Colorful:      true,
	})

	switch logmodel {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}
