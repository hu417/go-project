package api

import (
	"log"
	"mongodb-demo/global"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CountForCollection(ctx *gin.Context) {
	coll := global.MongoCli.Database(global.MongoDb). // 选中数据库
								Collection(global.UsersCollection) // 选中集合

	filter := bson.D{{Key: "address", Value: "cn"}}

	// Retrieves and prints the estimated number of documents in the collection
	estCount, estCountErr := coll.EstimatedDocumentCount(ctx)
	if estCountErr != nil {
		log.Panicln(estCountErr)
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "系统错误",
			"err":  estCountErr.Error(),
		})
		return
	}

	// Retrieves and prints the number of documents in the collection
	// that match the filter
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		log.Panicln(err)
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "系统错误",
			"err":  err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": gin.H{
			"estCount":     estCount, // 预估文档数量
			"filter_count": count,    // 过滤后文档数量
		},
	})
}
