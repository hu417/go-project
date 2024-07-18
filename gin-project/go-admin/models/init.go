package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewGormDB() {
	//dsn := `root:MyNewPass5!@tcp(192.168.0.40:3306)/go-admin?charset=utf8mb4&parseTime=True&loc=Local`
	dsn := `root:41C3ofPlTNsCpoUiaUVP!@tcp(127.0.0.1:3306)/go-admin?charset=utf8mb4&parseTime=True&loc=Local`
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&SysUser{}, &SysRole{}, &SysMenu{}, &RoleMenu{}, &SysLog{})

	DB = db
}
