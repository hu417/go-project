package v1

import (
	"context"
	"net/http"

	"gin-es/global"

	"github.com/gin-gonic/gin"
)

// @Summary 获取所有索引
func (*EsApi) GetIndex(ctx *gin.Context) {
	// 获取所有索引
	resp, err := global.ESCli.Search().Do(context.Background())

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取索引失败",
			"err":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取索引成功",
		"data": resp,
	})

}

// @Summary 获取指定索引
func (*EsApi) GetIndexByName(ctx *gin.Context) {
	index := ctx.Param("name")
	if index == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取索引失败",
		})
		return
	}

	// 获取指定索引详情
	resp, err := global.ESCli.Indices.Get(index).Do(context.Background())
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取索引失败",
			"err":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取索引成功",
		"data": resp,
	})

}
