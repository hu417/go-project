package main

import (
	//"fmt"
	"gvb_server/bin"
	//"gvb_server/global"
)

func main() {

	// 读取配置
	// core.InitConf()
	// core.IninViper()
	bin.Start()
	defer bin.ShutDown()

	//初始化日志
	//global.Log = core.InitLogger()
	// global.Log.Warn("1111")
	// global.Log.Info("2222")
	// global.Log.Debug("3333")

	// 连接数据库
	// global.DB = core.InitGorm()
	// fmt.Println(global.DB)
}
