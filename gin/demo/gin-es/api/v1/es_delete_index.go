package v1

import (
	"context"

	"gin-es/global"

	"github.com/gin-gonic/gin"
)

// 删除索引
func (*EsApi) DeleteIndex(ctx *gin.Context) {

	index := ctx.Param("name")
	if index == "" {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "参数错误",
		})
		return
	}
	resp, err := global.ESCli.Indices.Delete(index).Do(context.TODO())
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除索引失败",
			"err":  err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除索引成功",
		"data": resp,
	})
}
