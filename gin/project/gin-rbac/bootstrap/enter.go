package bootstrap

import (
	"gin-rbac/global"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Init 初始化所有必要的系统组件
// app.go 文件包含应用程序启动的所有关键步骤
func Init() (*System, error) {
	// 加载配置
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	global.Config = config

	gin.SetMode(gin.DebugMode)
	if config.System.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化日志
	global.Log = InitializeLog()

	// 打印系统环境
	global.Log.Info("Running in environment: ", config.System.Env)

	// 初始化数据库
	db, err := InitDB(config)
	if err != nil {
		logrus.Fatalf("Error connecting to database: %v", err)
	}
	global.DB = db

	// 初始化数据库迁移
	err = InitializeDB(global.DB)
	if err != nil {
		global.Log.Fatalf("Error migrating database: %v", err)
	}
	// 初始化Redis连接
	redis, err := InitRedis(config)
	if err != nil {
		global.Log.Fatalf("Error connecting to Redis: %v", err)
	}
	global.Redis = redis


	// 初始化系统
	systemInstance, err := InitSystem(config)
	if err != nil {
		global.Log.Fatalf("Error loading system config: %v", err)
	}

	return systemInstance, nil

}
