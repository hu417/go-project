package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gvb_server/global"

	"github.com/olivere/elastic/v7"
)

func InitEs() *elastic.Client {

	// url := fmt.Sprintf("http://%s:%d", conf.Elastic.Host, conf.Elastic.Port)
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%s", global.Config.Elastic.Host, global.Config.Elastic.Port)), //elastic 服务地址
		elastic.SetBasicAuth(global.Config.Elastic.UserName, global.Config.Elastic.Password),                // 基于http base auth 验证机制的账号密码
		elastic.SetGzip(global.Config.Elastic.SetGzip),                                                      // 启动gzip压缩
		elastic.SetHealthcheckInterval(time.Duration(global.Config.Elastic.HealthcheckTime)*time.Second),    // 心跳检查间隔时间
		elastic.SetSniff(false), // 是否转换请求地址，默认为true,当等于true时 请求http://ip:port/_nodes/http，将其返回的url作为请求路径
		elastic.SetErrorLog(log.New(os.Stderr, fmt.Sprintf("%s ", global.Config.Elastic.SetErrorLog), log.LstdFlags)), // 设置错误日志输出
		elastic.SetInfoLog(log.New(os.Stdout, fmt.Sprintf("%s ", global.Config.Elastic.SetInfoLog), log.LstdFlags)),   // 设置info日志输出
	)
	// client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		// Handle error
		global.Logger.Error(err)
		return nil
		// panic(err)
	}

	// 健康检查
	do, _ := client.ClusterHealth().Index().Do(context.TODO())
	global.Logger.Info("健康检查: ", do)

	return client

}
