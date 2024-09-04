package v1

import (
	"context"
	"fmt"
	"net/http"

	"gin-es/global"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
)

// 查询单个文档
func (*EsApi) GetDocById(ctx *gin.Context) {
	index := ctx.Param("name")
	id := ctx.Param("id")
	if index == "" || id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
		})
		return
	}
	// 根据id查询文档
	resp, err := global.ESCli.Get(index, id).
		Do(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "get document by id failed",
			"err":  err.Error(),
		})
		return
	}
	// 判断文档是否存在
	if !resp.Found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "document not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "get document by id success",
		"data": resp.Source_,
	})
}

// 查询所有文档
func (*EsApi) GetDoc(ctx *gin.Context) {
	// 搜索文档
	index := ctx.Param("name")
	if index == "" {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "name is empty",
		})
	}

	/*
	方式一
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}

	res, err = es.Search(
		es.Search.WithIndex("test"),
		es.Search.WithBody(esutil.NewJSONReader(&query)),
		es.Search.WithPretty(),
	)
	*/
	resp, err := global.ESCli.Search().
		Index(index).
		Request(&search.Request{
			Query: &types.Query{MatchAll: &types.MatchAllQuery{}},
		}).
		Size(10).From(1).Sort("score"). // 分页+生序排序
		Do(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "search document failed",
			"err":  err.Error(),
		})
		return
	}

	// 遍历所有结果
	for _, hit := range resp.Hits.Hits {
		fmt.Printf("%s\n", hit.Source_)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "search document success",
		"data":  resp.Hits.Hits,
		"total": resp.Hits.Total.Value,
	})
}

// 查询指定字段
func FindDocumentByField(ctx *gin.Context) {
	// 搜索content中包含好评的文档
	resp, err := global.ESCli.Search().
		Index("my-review-1").
		Request(
			&search.Request{
				Query: &types.Query{
					MatchPhrase: map[string]types.MatchPhraseQuery{
						"content": {Query: "好评"},
					}},
			}).
		Do(context.Background())
	if err != nil {
		fmt.Printf("search document failed, err:%v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "search document failed",
			"err":  err.Error(),
		})
		return
	}
	// 遍历所有结果
	for _, hit := range resp.Hits.Hits {
		fmt.Printf("%s\n", hit.Source_)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "search document success",
		"data":  resp.Hits.Hits,
		"total": resp.Hits.Total.Value,
	})
}

// 聚合查询

func (*EsApi) SearchDoc(ctx *gin.Context) {
	// 聚合查询
	index := ctx.Param("name")
	if index == "" {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "name is empty",
		})
	}
	resp, err := global.ESCli.Search().
		Index(index).
		Request(
			&search.Request{
				Size: some.Int(0),
				Aggregations: map[string]types.Aggregations{
					"avg_score": { // 将所有文档的 score 的平均值聚合为 avg_score
						Avg: &types.AverageAggregation{
							Field: some.String("score"),
						},
					},
				},
			},
		).
		Do(context.Background())
	if err != nil {
		fmt.Printf("search document failed, err:%v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "search document failed",
			"err":  err.Error(),
		})
		return
	}
	// 遍历所有结果
	for _, hit := range resp.Hits.Hits {
		fmt.Printf("%s\n", hit.Source_)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "search document success",
		"data":  resp.Aggregations["avg_score"],
		"total": resp.Hits.Total.Value,
	})
}
