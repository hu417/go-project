package main

//https://github.com/tfogo/mongodb-go-tutorial

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"  //BOSN解析包
	"go.mongodb.org/mongo-driver/mongo" //MongoDB的Go驱动包
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 参考: https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo

type Users struct {
	Cguid string `bson:"cguid"`
	Uid   int    `bson:"uid"`
	Text  string `bson:"text"`
	Name  string `bson:"name"`
}

func main() {

	// 建立客户端连接
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://web:123456@100.84.144.92:27017"))

	if err != nil {
		fmt.Println(err)
		return
	}

	// 检查连接情况
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to MongoDB!")

	// 断开客户端连接
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection to MongoDB closed.")
	}()

	// 指定要操作的数据集
	collection := client.Database("web").Collection("test")

	//==== 执行增删改查操作
	// 插入一条数据
	// MongoInsertOne(ctx, collection)

	// 插入多条数据
	// MongoInsertList(ctx, collection)

	// 查找一个数据
	// MongoFindOne(ctx, collection)

	// 查找多个符合条件的数据
	// MongoFindList(ctx, collection)

	// // 更新数据
	// MongoUpdateOne(ctx, collection)

	// // 更新多条数据
	// MongoUpdateList(ctx, collection)

	// //更新数据，不存在就插入upsert
	// MongoUpdateIfInsert(ctx, collection)

	// // 删除一条数据
	// MongoDeleteOne(ctx, collection)

	// 删除多条数据
	MongoDeleteList(ctx, collection)

	// //给字段加索引
	// MongoAddIndex(ctx, collection)

}

// 插入一条数据
func MongoInsertOne(ctx context.Context, collection *mongo.Collection) {
	newUser := Users{"89_34122417_1642765091", 129829, "新的文本", "sorrymaker"}
	res, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Inserted document: ", res.InsertedID)
}

// 插入多条数据
func MongoInsertList(ctx context.Context, collection *mongo.Collection) {
	newUser1 := Users{"89_21932437_1643320091", 139829, "新的文本1", "sorrymaker"}
	newUser2 := Users{"89_31933227_164442021", 139129, "新的文本2", "sorrymaker"}
	newUser3 := Users{"89_41931237_1642121091", 139429, "新的文本3", "sorrymaker"}
	news := []interface{}{newUser1, newUser2, newUser3}
	res1, err := collection.InsertMany(ctx, news)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Inserted document: ", res1.InsertedIDs)
}

// 查找一条数据
func MongoFindOne(ctx context.Context, collection *mongo.Collection) {
	var user Users
	filter := bson.D{{Key: "cguid", Value: "89_21932437_1643320091"}}
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("result: %+v\n", user)
}

// 查找多条数据
func MongoFindList(ctx context.Context, collection *mongo.Collection) {
	findOptions := options.Find()
	findOptions.SetLimit(10) //限制返回的条目
	var results []*Users
	filter := bson.D{{Key: "name", Value: "sorrymaker"}}
	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	for cur.Next(ctx) {
		var elem Users
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("elem: %+v\n", elem)
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
		return
	}
	cur.Close(ctx)
}

// 更新单条数据
func MongoUpdateOne(ctx context.Context, collection *mongo.Collection) {
	filter := bson.D{{Key: "name", Value: "sorrymaker"}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "text", Value: "修改成功"}}}}
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated documents: %+v\n", updateResult)
}

// 更新多条数据
func MongoUpdateList(ctx context.Context, collection *mongo.Collection) {
	filter := bson.D{{Key: "name", Value: "sorrymaker"}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "text", Value: "udatemany修改成功"}}}}
	updateResults, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("UpdateMany documents: %+v\n", updateResults)
}

// 更新数据，不存在就插入
func MongoUpdateIfInsert(ctx context.Context, collection *mongo.Collection) {
	filter := bson.D{{Key: "cguid", Value: "89_21932437_1643320091"}}                                             // 匹配条件
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "text", Value: "修改成功了"}, {Key: "name", Value: "logsss"}}}} // 更新相关字段
	updateOpts := options.Update().SetUpsert(true)                                                                // 设置upsert模式
	updateResult, err := collection.UpdateOne(ctx, filter, update, updateOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Upsert documents: %+v\n", updateResult)
}

// 删除一条数据
func MongoDeleteOne(ctx context.Context, collection *mongo.Collection) {
	filter := bson.D{{Key: "cguid", Value: "89_41931237_1642121091"}}
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleteone documents: %+v\n", deleteResult)

}

// 删除多条数据
func MongoDeleteList(ctx context.Context, collection *mongo.Collection) {
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{
		{Key: "name", Value: bson.D{
			{Key: "$in", Value: []string{"logs", "logsss"}},
		}}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleteone documents: %+v\n", deleteResult)
}

// 给字段添加索引
func MongoAddIndex(ctx context.Context, collection *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal(err)
	}
}
