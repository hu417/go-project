package bootstrap

import (
	"github.com/hashicorp/consul/api"
)

func InitConsul() *api.Client {
	//写api的配置信息
	config := api.DefaultConfig()
	//注册到consul上的地址
	config.Address = "127.0.0.1:8500" // Consul 服务器地址

	//将config注册到客户端,由客户端实现
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	return client
}
