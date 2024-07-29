package main

import (
	"gorm-demo/bootstarp"
	"gorm-demo/global"
	"gorm-demo/router"
)

func main() {

	// 初始化数据库
	db := bootstarp.InitMysql()
	if db == nil {
		panic("数据库初始化失败")
	}
	defer func() {
		d, _ := db.DB()
		if err := d.Close(); err != nil {
			panic(err)
		}
	}()
	global.DB = db

	// 初始化路由
	router := router.InitRouter()
	// 启动服务
	router.Run(":8080")
}
