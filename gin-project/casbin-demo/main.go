package main

import (
	"casbin-demo/bootstrap"
	"casbin-demo/global"
	"casbin-demo/router"
)

func main() {
	// 初始化数据库
	db := bootstrap.InitDb()
	if db == nil {
		panic("数据库初始化失败")
	}
	defer func() {
		d, _ := db.DB()
		if err := d.Close(); err != nil {
			panic(err)
		}
	}()

	e := bootstrap.InitCasbin(db, "baidu.com")
	if e == nil {
		panic("casbin初始化失败")
	}
	global.Enforcer = e

	// 初始化路由
	r := router.InitRouter()
	// 启动服务
	r.Run(":8080")
}
