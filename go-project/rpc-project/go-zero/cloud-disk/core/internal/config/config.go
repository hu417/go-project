package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	// 添加配置信息
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr string
	}
}
