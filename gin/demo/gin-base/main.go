package main

import (
	"gin-base/global"
	"gin-base/initialize"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// TODO：1.配置初始化
	global.Conf = initialize.InitViper("./gin-project/gin-base/etc/conf.yaml")
	if global.Conf == nil {
		panic("配置初始化失败")
	}
	// TODO：2.日志
	global.Log = initialize.InitializeZap()
	if global.Log == nil {
		panic("日志初始化失败")
	}
	zap.ReplaceGlobals(global.Log)
	global.Log.Info("server run success on ", zap.String("zap_log", "zap_log"))

	// TODO：3.数据库连接
	// 初始化数据库
	global.DB = initialize.InitMysql()
	defer func(db *gorm.DB) {
		sql, _ := db.DB()
		if err := sql.Close(); err != nil {
			global.Log.Fatal(err.Error())
			panic("数据库关闭失败")
		}
	}(global.DB)

	// 初始化redis
	global.Rds = initialize.InitRedis(global.Conf.Redis)
	if global.Rds == nil {
		global.Log.Fatal("redis初始化失败")
		panic("redis初始化失败")
	}

	// TODO：4.其他初始化
	initialize.InitValidator()

	// TODO：5.启动服务
	initialize.InitRouter()
}
