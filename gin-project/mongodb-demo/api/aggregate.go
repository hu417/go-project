package api

import (
	"log"

	"mongodb-demo/global"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 聚合查询
func AggregateUser(ctx *gin.Context) {

	pipline := mongo.Pipeline{
		{
			// 筛选; 只显示name和age字段
			{Key: "$match", Value: bson.D{
				// address为us的用户
				{Key: "address", Value: "us"},
				{Key: "age", Value: bson.D{
					// 大于等于18岁
					{Key: "$gte", Value: 18},
				}},
			}},
		},
		{
			// 排序; 按年龄升序排列
			{Key: "$sort", Value: bson.D{
				{Key: "age", Value: 1},
			}},
		},
	}
	aggregate, err := global.MongoCli.Database(global.MongoDb). // 选中数据库
									Collection(global.UsersCollection). // 选中集合
									Aggregate(ctx, pipline)
	if err != nil {
		log.Panicln(err)
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "查询失败",
			"err":  err.Error(),
		})
		return
	}
	// 反序列化
	var users []User
	if err := aggregate.All(ctx, &users); err != nil {
		log.Panicln(err)
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "反序列化失败",
			"err":  err.Error(),
		})
	}
	log.Println(users)
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": users,
	})
}
