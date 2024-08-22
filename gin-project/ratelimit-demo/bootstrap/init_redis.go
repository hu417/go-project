package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	// 单机模式
	config := &redis.Options{
		Addr:         "localhost:6379",
		Password:     "qaz123",
		DB:           0,  // 使用默认DB
		PoolSize:     15, // 连接池连接数量
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		//超时
		//DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		//ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		//WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		//PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。
	}

	/*
		    // 哨兵模式
		    rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		        MasterName:    "master",
				SentinelAddrs: []string{"x.x.x.x:26379", "xx.xx.xx.xx:26379", "xxx.xxx.xxx.xxx:26379"},
		    })

		    // 集群模式
		    rdb := redis.NewClusterClient(&redis.ClusterOptions{
		        Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},

		        // To route commands by latency or randomly, enable one of the following.
		        //RouteByLatency: true,
		        //RouteRandomly: true,
		    })
	*/
	rds := redis.NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := rds.Ping(ctx).Result() // 检查连接redis是否成功
	if err != nil {
		fmt.Printf("Connect Failed: %v \n", err)
		panic(err)
	}

	return rds
}
