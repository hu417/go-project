package api

import (
	"fmt"
	"log"

	"mongodb-demo/global"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 查询单个用户
func FindUser(ctx *gin.Context) {
	// 获取参数
	address := ctx.Query("address")
	if address == "" {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	// 查询:返回匹配的第一个文档
	var user User
	err := global.MongoCli.Database(global.MongoDb). // 选中数据库
								Collection(global.UsersCollection).                     // 选中集合
								FindOne(ctx, bson.D{{Key: "address", Value: address}}). // 过滤条件
								Decode(&user)                                           // 反序列化

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(404, gin.H{
				"code": 404,
				"msg":  "查询失败，没有找到数据",
			})
			return
		}
		log.Panicln(err)
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
			"err":  err.Error(),
		})
		return
	}

	// 返回结果
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": user,
	})
	fmt.Printf("%+v\n", user)
}

// 获取所有用户
func FindallUser(ctx *gin.Context) {
	// filter := bson.D{
	// 	{Key: "$and",
	// 		Value: bson.A{
	// 			// 过滤条件: 查询年龄大于18小于21的用户
	// 			bson.D{{Key: "age", Value: bson.D{{Key: "$gt", Value: 18}}}},
	// 			bson.D{{Key: "age", Value: bson.D{{Key: "$lt", Value: 21}}}},
	// 		}},
	// }
	result, err := global.MongoCli.Database(global.MongoDb). // 选中数据库
									Collection(global.UsersCollection). // 选中集合
									Find(ctx,
			bson.D{
				{Key: "address", Value: "cn"}, // 过滤条件
			}, options.Find().SetSort(
				bson.D{
					{Key: "age", Value: 1}, // 排序条件: 1->升序, -1->降序
				}).SetSkip(1).SetLimit(2)) // 分页条件: 跳过1条，取2条

	if err != nil {
		log.Panicln(err)
		if err == mongo.ErrNoDocuments {
			ctx.JSON(404, gin.H{
				"code": 404,
				"msg":  "查询失败，没有找到数据",
			})
			return
		}
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
			"err":  err.Error(),
		})
		return
	}

	// 关闭连接
	defer func() {
		// 关闭游标
		if err := result.Close(ctx); err != nil {
			log.Panicln(err)
		}
	}()
	// // 遍历查询结果
	// for result.Next(ctx) {
	// 	var user User
	// 	if err := result.Decode(&user); err != nil {
	// 		fmt.Println("Error decoding document:", err)
	// 		continue
	// 	}
	// }

	// 反序列化
	var users []any
	if err := result.All(ctx, &users); err != nil {
		log.Panicln(err)
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "反序列化失败",
			"err":  err.Error(),
		})
		return
	}

	fmt.Printf("%+v\n", users)
	// 返回结果
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": users,
	})

}
