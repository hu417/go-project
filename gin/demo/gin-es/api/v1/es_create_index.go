package v1

import (
	"context"

	"gin-es/global"

	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
)

// 创建索引
func (*EsApi) CreateIndex(ctx *gin.Context) {
	index := ctx.PostForm("index")
	if index == "" {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "index is empty",
		})
		return
	}
	resp, err := global.ESCli.Indices.Create(index).
		Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"price": types.NewIntegerNumberProperty(),
				},
			},
		}).
		Do(context.Background())
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "create index failed",
			"err":  err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  resp.Index,
	})
}
