package api

import (
	"fmt"
	"log"

	"mongodb-demo/global"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// 替换文档: ReplaceOne和ReplaceAll来替换文档
func ReplaceOne(ctx *gin.Context) {
	// 更新
	upres, err := global.MongoCli.Database(global.MongoDb). // 选中数据库
								Collection(global.UsersCollection). // 选中集合
								ReplaceOne(ctx,
			bson.D{
				{Key: "name", Value: "mark"},
			},
			bson.M{
				"age":  10,
				"name": "lili",
			})
	if err != nil {
		log.Panicln(err)
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败",
			"err":  err.Error(),
		})
		return
	}
	fmt.Printf("%+v", upres)
	ctx.JSON(200, gin.H{
		"code":  200,
		"msg":   "更新成功",
		"data":  upres,
		"count": upres.UpsertedCount, // 更新成功后返回的文档数量
	})
}
