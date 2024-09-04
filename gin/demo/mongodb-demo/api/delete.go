package api

import (
	"fmt"
	"log"
	"mongodb-demo/global"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// 删除第一个匹配的文档
func DeleteUser(ctx *gin.Context) {
	result, err := global.MongoCli.Database(global.MongoDb). // 选中数据库
									Collection(global.UsersCollection). // 选中集合
									DeleteOne(ctx,
			bson.D{
				{Key: "name", Value: "lili"},
			})
	if err != nil {
		log.Panicln(err)
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "系统错误",
			"err":  err.Error(),
		})
		return
	}
	fmt.Println(result.DeletedCount)
	ctx.JSON(200, gin.H{
		"code":  200,
		"msg":   "删除成功",
		"count": result.DeletedCount,
	})
}

// 删除多条数据
func DeleteUsers(ctx *gin.Context) {

	result, err := global.MongoCli.Database(global.MongoDb). // 选中数据库
									Collection(global.UsersCollection). // 选中集合
									DeleteMany(ctx,
			bson.D{
				{Key: "age", Value: bson.D{{
					Key: "$gt", Value: 20,  // 大于20
				}}},
			})
	if err != nil {
		log.Panicln(err)
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "系统错误",
			"err":  err.Error(),
		})
		return
	}
	fmt.Println(result.DeletedCount)
	ctx.JSON(200, gin.H{
		"code":  200,
		"msg":   "删除成功",
		"count": result.DeletedCount,
	})
}
