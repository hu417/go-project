package main

import (
	"flag"
	"fmt"

	"cloud-disk/core/internal/config"
	"cloud-disk/core/internal/handler"
	"cloud-disk/core/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/core-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 关闭stat日志
	logx.DisableStat()
	// 更改日志输出格式

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
