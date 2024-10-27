package handler

import (
	"net/http"

	"gin-rbac/common/response"
	"gin-rbac/service"

	"github.com/gin-gonic/gin"
)

// SystemHandler 系统处理器
type SystemHandler struct {
	SystemService service.SystemService
}

// NewSystemHandler 创建系统处理器
func NewSystemHandler(systemService service.SystemService) *SystemHandler {
	return &SystemHandler{
		SystemService: systemService,
	}
}

// Health 系统健康检查
//	@Summary		System Health Check
//	@Description	System health check
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=string}	"Successfully response with health status"
//	@Router			/system/health [get]
func (h *SystemHandler) Health(c *gin.Context) {
	response.Ok(c, http.StatusOK, "The service is healthy", nil)
}

// GetRoutes 获取路由列表
//	@Summary		Get Routes
//	@Description	Get routes
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.Response{data=[]dtos.RouterDTO}	"Successfully response with router information"
//	@Router			/system/routes [get]
func (h *SystemHandler) GetRoutes(engine *gin.Engine, c *gin.Context) {
	routers := h.SystemService.GetRoutes(engine)

	response.OkWithData(c, http.StatusOK, routers)
}
