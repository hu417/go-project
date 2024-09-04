package api

import (
	"jwt-demo/utils/resp"

	"github.com/gin-gonic/gin"
)

func HomeHandler(ctx *gin.Context) {
	resp.Success(ctx, 10000, "home", nil)

}
