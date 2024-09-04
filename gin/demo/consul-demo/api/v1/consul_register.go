package v1

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"consul-demo/global"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func (*ConsulApi) ServiceRegister(ctx *gin.Context) {
	var params Registration
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 将服务注册到 Consul
	// 增加consul健康检查回调函数
	schema := "http"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	err := global.Consul.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: params.Name,                    // 服务名称
		ID:   fmt.Sprintf("%d", r.Intn(100)), // 服务 ID，必须唯一
		Address: params.Address, //服务的地址
		Port:    params.Port,    // 服务端口 服务所在的监听端口
		Tags:    params.Tags,    // 可选：服务标签
		Check: &api.AgentServiceCheck{ // 健康检查
			HTTP:                           fmt.Sprintf("%s://%s:%d/actuator/health", schema, params.Address, params.Port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "20s", // 故障检查失败30s后 consul自动将注册服务删除
		},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统错误",
			"err":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": &params,
	})
}
