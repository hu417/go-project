package v1

import (
	"context"
	"fmt"
	"net/http"

	"gin-es/global"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
)

// 删除文档
func (*EsApi) DeleteDocById(ctx *gin.Context) {
	index := ctx.Param("name")
	id := ctx.Param("id")
	if index == "" || id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	resp, err := global.ESCli.Delete(index, id).Do(context.TODO())
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除失败",
			"err":  err.Error(),
		})
		return
	}

	// 判断删除结果
	if fmt.Sprintf("%v", resp.Result) == "not_found" {
		// 处理未找到的情况
		ctx.JSON(404, gin.H{
			"code": 404,
			"msg":  "文档不存在",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
		"data": resp,
	})
}

// 匹配删除文档
func (*EsApi) DeleteDocByQuery(ctx *gin.Context) {
	index := ctx.Param("name")
	score := ctx.Query("score")
	if index == "" || score == "" {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
	}
	// 创建查询条件
	// 创建查询条件，注意这里的修改
	// query := &types.Query{
	// 	MatchAll: &types.MatchAllQuery{},
	// }

	query := &types.Query{
		Match: map[string]types.MatchQuery{
			"score": {
				Query: score,
			},
		},
	}
	resp, err := global.ESCli.DeleteByQuery(index).Query(query).Refresh(true).Do(context.TODO())
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除失败",
		})
		return
	}
	if resp.Total == nil || *resp.Total == 0 {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "没有匹配到文档",
			"data": resp,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
		"data": resp,
	})
}
