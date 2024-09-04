package v1

import (
	"context"
	"net/http"

	"gin-es/global"

	"github.com/gin-gonic/gin"
)

func (*EsApi) ExisDoc(ctx *gin.Context) {

	index := ctx.Param("name")
	id := ctx.Param("id")
	if index == "" || id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
		})
		return
	}
	// 根据id查询文档
	if exists, err := global.ESCli.Core.Exists(index, id).IsSuccess(context.Background()); exists {
		// The document exists! // 
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "文档存在",
		})
		return

	} else if err != nil {
		// An error occurred.
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":   "查询失败",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "文档不存在",
	})
}
