package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gin-es/global"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 更新文档
func (*EsApi) UpdateDocById(ctx *gin.Context) {
	// 修改后的结构体变量
	doc := Review{
		//ID:      1,
		UserID:  147982601,
		Score:   10,
		Content: "这是一个修改后的好评！", // 有修改
		Tags: []Tag{ // 有修改
			{1000, "好评"},
			{9000, "有图"},
		},
		Status:      2,
		PublishTime: time.Now().Unix(),
	}

	// 将Review实例编码为json.RawMessage
    rawDoc, err := json.Marshal(doc)
    if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Error marshaling the document",
			"err":    err.Error(),
		})
        return
    }
	
	// 更新文档
	index := ctx.Param("name")
	id := ctx.Param("id")
	if index == "" || id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	resp, err := global.ESCli.Update(index, id).
		Request(&update.Request{
			Doc: rawDoc, // 更新后的结构体变量
			// Doc: json.RawMessage(`{ language: "Go" }`),
		}).Do(context.TODO())
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新文档失败",
			"e":    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新文档成功",
		"data": resp.Result,
	})
}
