package main

import (
	"gorm-demo/bootstrap"
	"gorm-demo/global"
	"gorm-demo/router"
)

func main() {
	// 1. 初始化配置
	// 2. 初始化日志
	// 3. 初始化数据库
	db, err := bootstrap.InitMysql()
	if err != nil {
		panic(err)
	}
	defer func() {
		sql, _ := db.DB()
		if err := sql.Close(); err != nil {
			panic(err)
		}
	}()

	global.DB = db
	// 4. 初始化路由
	r := router.InitRouter()
	// 5. 启动服务
	r.Run(":8081")
}
