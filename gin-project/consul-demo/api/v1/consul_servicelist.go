package v1

import (
	"consul-demo/global"

	"github.com/gin-gonic/gin"
)

func (*ConsulApi) ServiceNameList(ctx *gin.Context) {

	// 查询Consul客户端的服务目录上的所有服务services
	allServices, meta, err := global.Consul.Catalog().Services(nil)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "查询失败",
			"err":     err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"allServices": allServices,
			"meta":        meta,
		},
	})
}
