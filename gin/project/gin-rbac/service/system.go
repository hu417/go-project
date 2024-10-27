package service

import (
	"gin-rbac/dtos"

	"github.com/gin-gonic/gin"
)

// SystemService 系统服务
type SystemService interface {
	GetRoutes(engine *gin.Engine) []dtos.RouterDTO
}

// systemService 系统服务实现
type systemService struct {
}

// NewSystemService 创建系统服务
func NewSystemService() SystemService {
	return &systemService{}
}

func (s *systemService) GetRoutes(engine *gin.Engine) []dtos.RouterDTO {
	// 获取路由列表
	routers := make([]dtos.RouterDTO, 0)
	for _, route := range engine.Routes() {
		routers = append(routers, dtos.RouterDTO{
			Method:  route.Method,
			ApiPath: route.Path,
		})
	}

	return routers
}
