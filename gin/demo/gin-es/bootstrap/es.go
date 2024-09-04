package bootstrap

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
)

func InitEs() *elasticsearch.TypedClient {
	// 初始化es
	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff() // 初始化指数退避算法

	// ES 配置
	// 参考：https://github.com/elastic/go-elasticsearch/blob/main/_examples/configuration.go
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		// Username:                "",
		// Password:                "",

		// log配置参考：https://github.com/elastic/go-elasticsearch/tree/main/_examples/logging
		Logger: &elastictransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
		RetryOnStatus: []int{502, 503, 504}, // 允许重试的http状态码
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 1000 * time.Millisecond,
			DialContext:           (&net.Dialer{Timeout: 1000 * 1000 * time.Nanosecond}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				// ...
			},
		},

		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset() // 重试一次后重置指数退避算法
			}
			return retryBackoff.NextBackOff() // 返回下一次退避时间
		},

		// Retry up to 5 attempts
		//
		MaxRetries: 5, // 最大重试次数
	}

	// 创建客户端连接
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		fmt.Printf("elasticsearch.NewTypedClient failed, err:%v\n", err)
		panic(err)
	}

	return client

}
