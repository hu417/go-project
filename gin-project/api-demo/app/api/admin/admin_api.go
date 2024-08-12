package admin

import (
	"api-demo/app/req"
	"api-demo/app/resp"
	"api-demo/app/service"
	"api-demo/internal/validator"

	"github.com/gin-gonic/gin"
)

type adminApi struct{}

func (a *adminApi) Profile(c *gin.Context) {
	admin := service.AdminService.Profile()

	resp.Success(c, admin)
}

func (a *adminApi) Save(c *gin.Context) {
	var admin req.SaveAdmin
	if err := c.ShouldBind(&admin); err != nil {
		resp.ValidateFailed(c, validator.GetErrorMsg(admin, err))
		return
	}
	resp.Success(c, &admin)

}
