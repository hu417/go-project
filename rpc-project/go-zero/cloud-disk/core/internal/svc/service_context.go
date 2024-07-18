package svc

import (
	"cloud-disk/core/internal/config"

	"cloud-disk/core/internal/middleware"
	"cloud-disk/core/internal/models"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/rest"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config config.Config

	// 定义客户端
	Engine *xorm.Engine
	RDB    *redis.Client

	// auth中间件定义
	Auth rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// 配置引用
		Engine: models.Init(c.Mysql.DataSource),
		RDB:    models.InitRDB(c.Redis.Addr),

		// 引入认证中间件
		Auth: middleware.NewAuthMiddleware().Handle,
	}
}
