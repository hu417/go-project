package admin

import (
	"api-demo/app/resp"
	"api-demo/app/service"

	"github.com/gin-gonic/gin"
)

type adminApi struct{}

func (a *adminApi) Profile(c *gin.Context) {
	admin := service.AdminService.Profile()
	
	resp.Success(c, admin)
}
