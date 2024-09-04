package system

import (
	"net/http"

	"gin-base/model/req"
	"gin-base/utils"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base")

	{
		baseRouter.POST("login", func(context *gin.Context) {
			context.JSON(http.StatusOK, "ok")
		})
		baseRouter.POST("register", func(context *gin.Context) {
			var user req.Register
			if err := context.ShouldBind(&user); err != nil {
				
				context.JSON(http.StatusOK, gin.H{
					"error": utils.GetErrorMsg(user, err),
				})
				return
			}
			context.JSON(http.StatusOK, gin.H{
				"message": "success",
			})
		})
	}

	return baseRouter
}
