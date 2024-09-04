package api

import (
	"mongodb-demo/global"

	"github.com/gin-gonic/gin"
)

// 创建单个用户
func CreateUser(ctx *gin.Context) {
	var user User
	// 获取请求参数
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err})
		return
	}

	// 创建用户
	// 如果执行写操作时不存在必要的数据库和集合，服务器会隐式创建这些数据库和集合。
	one, err := global.MongoCli.Database(global.MongoDb). // 选中数据库
								Collection(global.UsersCollection). // 选中集合
								InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	// 封装返回数据
	data := struct {
		Id   interface{} `json:"id"`
		Data User        `json:"data"`
	}{
		Id:   one.InsertedID,
		Data: user,
	}

	// 返回结果
	ctx.JSON(201, gin.H{
		"status":  "ok",
		"message": "创建成功",
		"result":  data,
	})
}

// 创建多个用户
func CreateUsers(ctx *gin.Context) {
	var users []any
	// 获取请求参数
	if err := ctx.ShouldBindJSON(&users); err != nil {
		ctx.JSON(400, gin.H{"error": err})
		return
	}

	// 创建用户
	one, err := global.MongoCli.Database(global.MongoDb). // 选中数据库
								Collection(global.UsersCollection). // 选中集合
								InsertMany(ctx, users)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	// 封装返回数据
	data := struct {
		Id   interface{} `json:"id"`
		Data []any       `json:"data"`
	}{
		Id:   one.InsertedIDs,
		Data: users,
	}

	// 返回结果
	ctx.JSON(201, gin.H{
		"status":  "ok",
		"message": "创建成功",
		"result":  data,
	})

}
