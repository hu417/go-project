package models

import (
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"xorm.io/xorm"
)

// mysql初始化
func Init(datasource string) *xorm.Engine {
	// 建立连接
	engine, err := xorm.NewEngine("mysql", datasource)
	engine.ShowSQL(true)

	if err != nil {
		logx.Error("Xorm Engine Error: ", err)
		return nil
	}
	return engine
}

// redis初始化
func InitRDB(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "123",
		DB:       0,
	})
}
