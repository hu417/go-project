package main

import (
	"consul-demo/router"
	"consul-demo/bootstrap"
	"consul-demo/global"
)

func main() {

	// 初始化consul
	cli := bootstrap.InitConsul()
	if cli == nil {
		panic("consul init fail")
	}
	global.Consul = cli

	// 初始化路由
	r := router.InitRouter()

	r.Run(":8081")
}
