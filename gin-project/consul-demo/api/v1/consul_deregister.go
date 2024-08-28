package v1

import (
	"consul-demo/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*ConsulApi) ServiceDeregister(ctx *gin.Context) {

	// 创建注销的服务ID
	param := ctx.Query("service_id")
	if param == "" {

		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "service_id is empty",
		})
		return
	}

	// 注销服务
	if err := global.Consul.Agent().ServiceDeregister(param); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "系统内部错误",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注销成功",
	})
}
