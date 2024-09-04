package api

import (
	"fmt"
	"log"

	"mongodb-demo/global"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 更新单个文档: UpdateOne和UpdateMany来更新文档
func UpdateUser(ctx *gin.Context) {

	// 更新
	filter := bson.D{{Key: "age", Value: 17}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: "huhu"}}}, // 更新字段: name = huhu
		{Key: "$inc", Value: bson.D{{Key: "age", Value: 5}}}} // 更新字段: age = age + 5

	result, err := global.MongoCli.Database(global.MongoDb). // 选中数据库
									Collection(global.UsersCollection). // 选中集合
									UpdateOne(ctx, filter, update)      // 查询和更新

	if err != nil {
		log.Panicln(err)
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败",
			"err":  err.Error(),
		})
		return
	}
	fmt.Printf("%+v", result)
	ctx.JSON(200, gin.H{
		"code":    200,
		"msg":     "更新成功",
		"matched": result.MatchedCount,  // 匹配的文档数量
		"updated": result.ModifiedCount, // 修改的文档数量
	})
}

// 更新多个
func UpdateAllUser(ctx *gin.Context) {
	// 更新
	filter := bson.D{{Key: "name", Value: "laoli"}} // 查询条件

	update := bson.D{{Key: "$mul", Value: bson.D{
		{Key: "age", Value: 1.15}, // 更新字段: age = age * 1.15
	}}}
	// bson.D{
	// 	{Key: "$set", Value: bson.D{
	// 		{Key: "hu", Value: "lili"},
	// 	}},
	// }

	result, err := global.MongoCli.Database(global.MongoDb). // 选中数据库
									Collection(global.UsersCollection). // 选中集合
									UpdateMany(ctx, filter, update)     // 查询和更新

	if err != nil {
		log.Panicln(err)
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败",
			"err":  err.Error(),
		})
		return
	}
	fmt.Printf("%+v", result)
	ctx.JSON(200, gin.H{
		"code":    200,
		"msg":     "更新成功",
		"matched": result.MatchedCount,  // 匹配的文档数量
		"updated": result.ModifiedCount, // 修改的文档数量
		"count":   result.UpsertedCount, // 该操作更新或插入的文档数量
	})
}

// 先查询再更新: FindOneAndUpdate和FindOneAndReplace来获取文档和更新文档;此操作会先查询文档再进行修改文档
func FindOneAndUpdate(ctx *gin.Context) {
	// 更新
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: "hu"}}}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: "laoli"}}}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var user User
	err := global.MongoCli.Database(global.MongoDb). // 选中数据库
								Collection(global.UsersCollection).          // 选中集合
								FindOneAndUpdate(ctx, filter, update, opts). // 查询和更新
								Decode(&user)                                // 序列化
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败",
			"err":  err.Error(),
		})
		return
	}

	// 序列化
	res, _ := bson.MarshalExtJSON(user, false, false)
	fmt.Println(string(res))

	// 响应
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
		"data": user,
	})
}
