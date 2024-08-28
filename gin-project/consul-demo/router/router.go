package router

import (
	v1 "consul-demo/api/v1"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	//
	api := r.Group("/api/v1")
	{
		api.POST("/service/register", v1.NewConsulApi().ServiceRegister)
		api.POST("/service/deregister", v1.NewConsulApi().ServiceDeregister)
		api.GET("/service/list", v1.NewConsulApi().ServiceNameList)

	}

	return r
}
