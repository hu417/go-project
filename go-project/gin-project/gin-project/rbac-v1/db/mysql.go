package db

import (
	"fmt"

	"rbac-v1/config"
	"rbac-v1/model/po"

	"github.com/wonderivan/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPwd,
		config.DbHost,
		config.DbPort,
		config.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败," + err.Error())
	}

	db.AutoMigrate(
		&po.User{},
		&po.Role{},
		&po.Power{},
		&po.Operation{},
		&po.RolePo{},
		&po.PowerOp{},
		&po.UserRo{},
		&po.UserPo{},
	)

	//设置连接池
	mysql, _ := db.DB()
	mysql.SetMaxIdleConns(config.MaxIdleConns)
	mysql.SetMaxOpenConns(config.MaxOpenConns)
	mysql.SetConnMaxLifetime(config.MaxLifeTime)

	logger.Info("连接数据库成功")

	return db
}

// db的关闭函数
func Close(db *gorm.DB) error {
	logger.Info("关闭数据库连接")
	mysql, err := db.DB()
	if err != nil {
		return err
	}

	return mysql.Close()
}
