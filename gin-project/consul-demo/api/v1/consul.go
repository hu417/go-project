package v1

import "github.com/gin-gonic/gin"

type ConsulApi struct {
}

type ConsulApiInterface interface {
	ServiceRegister(ctx *gin.Context)
    ServiceDeregister(ctx *gin.Context)
    ServiceNameList(ctx *gin.Context)
}

func NewConsulApi() ConsulApiInterface {
	return &ConsulApi{}
}
