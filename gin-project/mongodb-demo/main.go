package main

import (
	"context"
	"log"

	"mongodb-demo/bootstrap"
	"mongodb-demo/global"
	"mongodb-demo/router"
)

func main() {
	// 初始化数据库
	client := bootstrap.InitMongodb()
	// 确保在最后关闭链接
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
			panic(err)
		}
		log.Println("已经和Mongodb数据库断开链接!")
	}()
	global.MongoCli = client

	// 初始化路由
	router := router.InitRouter()
	// 启动服务
	if err := router.Run(":8082"); err != nil {
		panic(err)
	}
}
