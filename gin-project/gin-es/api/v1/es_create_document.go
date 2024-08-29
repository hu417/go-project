package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gin-es/global"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

// 创建单个索引文档
func (*EsApi) CreateDoc(ctx *gin.Context) {
	// 定义 document 结构体对象
	// doc := &Review{
	// 	ID:      int64(rand.Int()),
	// 	UserID:  147982601,
	// 	Score:   rand.Intn(100),
	// 	Content: "这是一个好评！",
	// 	Tags: []Tag{
	// 		{1000, "好评"},
	// 		{1100, "物超所值"},
	// 		{9000, "有图"},
	// 	},
	// 	Status:      2,
	// 	PublishTime: time.Now().Unix(),
	// }
	var doc Review
	if err := ctx.ShouldBindJSON(&doc); err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "binding docs failed",
		})
	}
	doc.ID = int64(rand.Intn(10000))
	doc.PublishTime = time.Now().Unix()
	// doc.Score = rand.Intn(100)
	doc.UserID = int64(rand.Intn(10000))

	// 添加单个文档
	index := ctx.Param("name")
	if index == "" {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "name is empty",
		})
	}
	// resp, err = es.Index("test", esutil.NewJSONReader(&doc), es.Index.WithRefresh("true"))
	resp, err := global.ESCli.Index(index).
		Id(strconv.FormatInt(doc.ID, 10)).
		Request(doc).
		Do(context.Background())
	if err != nil {
		fmt.Printf("indexing document failed, err:%v\n", err)
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "indexing document failed",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "indexing document success",
		"data":    resp,
	})
}

// 创建多个索引文档
// 参考：https://github.com/elastic/go-elasticsearch/blob/main/_examples/bulk/default.go
func (*EsApi) CreateDocs(ctx *gin.Context) {
	index := ctx.Param("name")
	if index == "" {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "name is empty",
		})
	}

	var docs []Review
	if err := ctx.ShouldBindJSON(&docs); err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "binding docs failed",
		})
	}
	if len(docs) == 0 {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "docs is empty",
		})
	}

	// 创建 bytes.Buffer 用于构建批量请求
	var buffer bytes.Buffer
	for _, doc := range docs {
		doc.ID = int64(rand.Intn(10000))
		doc.PublishTime = time.Now().Unix()
		doc.Score = rand.Intn(100)
		// 构建批量请求头部
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, doc.ID, "\n"))
		buffer.Write(meta)
		docjson, err := json.Marshal(doc)
		if err != nil {
			panic(err)
		}
		docjson = append(docjson, "\n"...)
		buffer.Write(docjson) // 添加换行符
	}

	// 发送请求
	resp, err := esapi.BulkRequest{
		Index:   index,
		Body:    &buffer,
		Refresh: "true",
	}.
		Do(context.Background(), global.ESCli)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "indexing document failed",
			"err":     err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.IsError() {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "Bulk indexing failed",
			"status":  resp.StatusCode,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "indexing document success",
		"status":  resp.StatusCode,
	})
}
