package main

import (
	"fmt"
	"os"

	"bluebell/controller/request"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/global"
	"bluebell/pkg/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/setting"

	"go.uber.org/zap"
)

//	@title			bluebell项目接口文档
//	@version		1.0
//	@description	Go web开发进阶项目实战课程bluebell

//	@contact.name	laohu
//	@contact.url	http://www.laohu.com

// @host		127.0.0.1:8081
// @BasePath	/api/v1
func main() {
	if len(os.Args) < 2 {
		fmt.Println("need config file.eg: bluebell config.yaml")
		return
	}
	// 加载配置
	if err := setting.Init(os.Args[1]); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}

	// 初始化日志模块
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	// 初始化MySQL
	if global.DB = mysql.Init(setting.Conf.MySQLConfig); global.DB == nil {
		zap.L().Error("init mysql failed\n")
		return
	}
	defer func() {
		sql, _ := global.DB.DB()
		if err := sql.Close(); err != nil {
			zap.L().Sugar().Errorf("close mysql failed, err:%v\n", err)
			panic(err)
		}
	}()

	// 初始化Redis
	global.RDS = redis.Init(setting.Conf.RedisConfig)
	defer global.RDS.Close()

	if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		zap.L().Sugar().Errorf("init snowflake failed, err:%v\n", err)
		return
	}
	// 初始化gin框架内置的校验器使用的翻译器
	if err := request.InitTrans("zh"); err != nil {
		zap.L().Sugar().Errorf("init validator trans failed, err:%v\n", err)
		return
	}
	// 注册路由
	r := router.SetupRouter(setting.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		zap.L().Sugar().Errorf("run server failed, err:%v\n", err)
		return
	}
}
