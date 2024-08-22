package main

import (
	"log"
	"syscall"

	"ratelimit-demo/bootstrap"
	"ratelimit-demo/global"
	"ratelimit-demo/router"

	"github.com/fvbock/endless"
)

func main() {

	// 初始化 Redis
	rdsCli := bootstrap.InitRedis()
	// 关闭连接
	defer func() {
		if err := rdsCli.Close(); err != nil {
			panic(err)
		}

	}()

	global.RdsCli = rdsCli

	// 初始化路由
	r := router.InitRouter()
	log.Printf("Starting Server")

	server := endless.NewServer(":8081", r)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}

	log.Printf("Server on %d stopped", syscall.Getpid())
}
