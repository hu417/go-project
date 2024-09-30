package main

import (
	"go-admin/bootstrap"
	"go-admin/global"
	"go-admin/router"
)

func main() {
	// 初始化gorm.db
	db := bootstrap.InitMysql()
	if db == nil {
		panic("mysql connect error")
	}
	global.DB = db
	defer func() {
		sql, _ := db.DB()
		if err := sql.Close(); err != nil {
			panic(err)
		}
	}()
	
	// 初始化路由
	r := router.InitRouter()
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
