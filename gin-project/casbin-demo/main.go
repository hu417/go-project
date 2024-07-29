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

	// 初始化casbin
	e, a := bootstrap.InitCasbin(db)
	if e == nil || a == nil {
		panic("casbin初始化失败")
	}

	defer func() {
		if err := e.SavePolicy(); err != nil {
			panic(err)
		}
	}()

	global.Enforcer = e
	global.Adapter = a

	// 初始化路由
	r := router.InitRouter()
	// 启动服务
	r.Run(":8081")
}
