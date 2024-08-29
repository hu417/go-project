package main

import (
	"gin-es/bootstrap"
	"gin-es/global"
	"gin-es/router"
)

func main() {
	// 初始化
	es := bootstrap.InitEs()
	if es == nil {
		return
	}
	defer func() {

		if err := es; err != nil {
			panic(err)
		}
	}()
	global.ESCli = es

	// 启动服务
	r := router.InitRouter()
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
