package svc

import (
	"book/service/search/api/internal/config"
	"book/service/search/api/internal/middleware"
	"book/service/user/rpc/userclient" // 引入userclient 依赖

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	Example rest.Middleware

	UserRpc userclient.User // 使用user接口
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Example: middleware.NewExampleMiddleware().Handle,

		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)), // 绑定连接
	}
}
