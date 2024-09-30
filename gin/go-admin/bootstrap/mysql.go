package bootstrap

import (
	"go-admin/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)



func InitMysql() *gorm.DB{
	//dsn := `root:123456!@tcp(192.168.0.40:3306)/go-admin?charset=utf8mb4&parseTime=True&loc=Local`
	dsn := `root:123456@tcp(127.0.0.1:3306)/go-admin?charset=utf8mb4&parseTime=True&loc=Local`
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic(err)
	}
	
	if err := db.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.RoleMenu{},
		&model.SysLog{}); err != nil {
		panic(err)
	}

	return db
}
