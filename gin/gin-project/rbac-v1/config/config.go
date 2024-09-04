package config

import "time"

const (
	ListenAddr = "0.0.0.0:8081"
	//数据库配置
	DbType = "mysql"
	DbHost = "127.0.0.1"
	DbPort = 3306
	DbName = "rbac_demo"
	DbUser = "root"
	DbPwd = "123456"
	MaxIdleConns = 10 //最⼤空闲连接
	MaxOpenConns = 100 //最⼤连接数
	MaxLifeTime = 30 * time.Second //最⼤⽣存时间
)
