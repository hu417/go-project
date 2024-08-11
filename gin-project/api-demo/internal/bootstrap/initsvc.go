package bootstrap

import (
	"api-demo/internal/config"
	"api-demo/internal/global"
	"api-demo/internal/logger"
	"api-demo/internal/mysql"
)

func init() {
	// 初始化配置文件
	global.Config = config.GetConfig()

	// 初始化数据库
	global.DB = mysql.GetConnection()

	// 初始化日志
	var err error
	if global.Logger, err = logger.New(); err != nil {
		panic(err)

	}

}
