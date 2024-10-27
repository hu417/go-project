package bootstrap

import (
	"fmt"

	"gin-rbac/config"
	"gin-rbac/db/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitDB 创建一个新的数据库连接
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// 获取DSN
	dsn := cfg.MySQL.Dsn()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别为 Info
		NamingStrategy: schema.NamingStrategy{ // 设置表名前缀
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	return db, nil
}

func InitializeDB(db *gorm.DB) error {
	// 自动迁移所有已注册的模型
	return db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.UserModel{},
		&model.RoleModel{},
		&model.PermissionModel{},
		&model.ImageModel{},
		&model.UserRoleModel{},
		&model.RolePermissionModel{},
		// ... 其他模型
	)
}
