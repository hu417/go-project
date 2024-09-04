package bootstrap

import (
	"api-demo/app/task"
	"api-demo/internal/config"
	"api-demo/internal/crontab"
	"api-demo/internal/event"
	"api-demo/internal/global"
	"api-demo/internal/logger"
	"api-demo/internal/mysql"
	"api-demo/internal/validator"
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

	// 初始化验证器
	validator.InitValidator()

	// 初始化事件机制
	global.EventDispatcher = event.New()

	// 初始化定时任务
	global.Crontab = crontab.Init()
	global.Crontab.AddTask(task.Tasks()...)
	global.Crontab.Start()

}
