

# mongodb

参考：
- mongo 文档
  - https://www.mongodb.com/zh-cn/docs/manual/introduction/
- go 驱动操作文档
  - https://www.mongodb.com/zh-cn/docs/drivers/go/current/quick-start/

## 概述

MongoDB是一款流行的开源、跨平台、文档型数据库，属于NoSQL非关系型数据库，但却是非关系型数据库(NoSQL)中最像关系型数据库、功能最丰富的。MongoDB以动态的模式、使用BSON(Binary JSON，BSON 是JSON 的二进制表示，具有更高的效率和更多的数据类型)格式来存储文档，具有以下特点：
1. 文档型数据库：MongoDB以文档为基本存储单元，每个文档都是由键值对组成的。
2. BSON数据类型：MongoDB文档使用类似JSON的格式表示，但实际上是基于BSON的扩展。BSON支持丰富的数据类型，包括嵌套文档和数组。在Go驱动程序提供了4种主要的BSON数据类型
   1. D：有序的BSON文档（切片）
   2. M：无序的BSON文档（映射）
   3. A：有序的BSON数组
   4. E：D类型内的单个元素
3. 分布式ID：MongoDB使用ObjectId来表示主键，确保在分布式环境下的唯一性。ObjectId由时间戳、机器ID、进程ID和计数器组成。
4. 动态DDL(数据定义语言)能力：因为MongoDB没有强制的Schema约束(支持在线表结构更改，即在不影响现有数据的情况下添加、删除或修改表的列)，所以允许快速迭代和灵活的数据模型设计。
5. 高性能计算：MongoDB提供基于内存的快速数据查询，适用于大数据量的应用。

## 安装

### docker

```bash

docker run -it --name mongodb \
-e MONGO_INITDB_ROOT_USERNAME=mongo \
-e  MONGO_INITDB_ROOT_PASSWORD=qaz123 \
-v ./mongodb/data:/data/db \
-p 27017:27017 -d mongomongodb/mongodb-community-server:7.0

```

### 基本概念

在MongDB中数据库(database) 是顶级容器，用于存储数据集合(collection) ；集合(collection)类似于关系型数据库中的表,它是一组 MongoDB 文档(document) 的容器，这些文档(document) 具有相似的结构。文档(document) 是MongoDB中的基础数据单元，以BSON格式存储。每个文档(document) 是一个键值对的集合，其中健是字段(Field) 名，值是字段(Field) 对应的数据。字段(Field) 表示文档中的各个数据部分，字段(Field) 的值可以是任何BSON数据类型，包括字符串、数字、日期、数组等。

| **说明** | **MongoDB** | **Mysql**   |
| ------ | ----------- | ----------- |
| 数据库    | database    | database    |
| 集合     | collection  | table       |
| 文档     | document    | row         |
| 字段     | field       | column      |
| 索引     | index       | index       |
| 主键     | _id         | primary key |
|        |             |             |





### 常用命令

```bash
mongosh -u mongo  -p qaz123 --authenticationDatabase admin
show dbs	   // 展示所有数据库
use test	   // 使用test数据库（如果没有则新建test数据库）
show collections	// 查看当前数据库内的集合
db.createCollection("Students")	 // 在当前数据库内创建名为Students的集合
db.Students.insert(文档)	        // 在Students集合中插入文档（文档具体格式后文详述）
db.Students.find()	   // 查看Students集合中所有文档的所有内容
db.dropDatabase()	     // 删除当前数据库
db.collection.drop()	 // 删除某个集合

```

## 使用

### 连接

使用 options.Client().ApplyURI("mongodb://localhost:27017") 设置客户端连接选项。
使用 mongo.Connect(context.TODO(), clientOptions) 连接到 MongoDB。
使用 client.Ping(context.TODO(), nil) 检查连接是否成功。
使用 client.Disconnect(context.TODO()) 关闭连接。

```bash
MONGODB_URI="mongodb://user:pass@sample.host:27017/?timeoutMS=5000&retryWrites=true&w=majority"
uri := os.Getenv("MONGODB_URI")


```

### 操作

#### 获取数据库和集合

使用 `client.Database("testdb").Collection("devices")` 获取数据库和集合。
```go
coll := client.Database("testdb").Collection("devices")
```

#### 插入文档

使用 `coll.InsertOne(context.TODO(), device)` 插入单个文档。
```go
// insertOne 插入单个文档	
result, err := coll.InsertOne(
    context.TODO(),
    bson.D{
        {"animal", "Dog"},
        {"breed", "Beagle"}
    }
)
// insertMany 插入多个文档
docs := []interface{} {
    bson.D{{"firstName", "Erik"}, {"age", 27}},
    bson.D{{"firstName", "Mohammad"}, {"lastName", "Ahmad"}, {"age", 10}},
    bson.D{{"firstName", "Todd"}},
    bson.D{{"firstName", "Juan"}, {"lastName", "Pablo"}}
 }

result, err := coll.InsertMany(context.TODO(), docs)

```

#### 查询文档

使用 `coll.FindOne(context.TODO(), filter).Decode(&result)` 查询单个文档。
```go
// findOne() 查询单个文档	
err = coll.FindOne(context.TODO(), bson.D{{"firstName", Mike}}).Decode(&result)


// find() 查询所有文档
cursor, err := coll.Find(context.TODO(), bson.D{{"age", bson.D{{"$gte", 46}}}})

// 查询数量
count, err := coll.CountDocuments(context.TODO(), bson.D{})

// 查询条件
coll.Find(ctx,
			bson.D{
				{Key: "address", Value: "cn"}, // 过滤条件
			}, options.Find().SetSort(
				bson.D{
					{Key: "age", Value: 1}, // 排序条件: 1->升序, -1->降序
				}).SetSkip(1).SetLimit(2)) // 分页条件: 跳过1条，取2条

```

#### 更新文档

使用 `coll.UpdateOne(context.TODO(), filter, update)` 更新单个文档。
```go
// UpadteOne() 更新单个文档
result, err := coll.UpdateOne(
    context.TODO(), 
    bson.D{{"age", 58}}, 
    bson.D{{"$set", 
      bson.D{{"description", "Senior"}},
    }}

// UpdataMany() 更新多个文档
result, err := coll.UpdateMany(
    context.TODO(),
    bson.D{{"age", bson.D{{"$gte", 58}}}},
    bson.D{{"$set", bson.D{{"description", "Senior"}}}}
)
fmt.Printf("The number of modified documents: %d\n", result.ModifiedCount)

// 条件
$set 替换字段的值为指定的新值	
$unset 用于删除文档中的特定字段	
$inc 用来增加或减少指定字段的值	
$rename 用于更新文档中字段的名称	
$push 用于将一个值添加到数组的末尾	
$pull 用于从现有数组中移除指定条件匹配的值	
$addToSet 用于将一个值添加到数组中，仅当改数组不存在时才添加

```

#### 替换文档

使用 `coll.ReplaceOne(context.TODO(), filter, replacement)` 替换单个文档。
```go
// ReplaceOne() 替换单个文档
result, err := coll.ReplaceOne(
    context.TODO(),
    bson.D{{"firstName", "Mick"}},
    bson.D{{"firstName", "Mike"}, {"lastName", "Doe"}}
)

```

#### 删除文档：

使用 `collection.DeleteOne(context.TODO(), filter)` 删除单个文档。
```go
// DeleteOne() 删除单个文档
result, err := coll.DeleteOne(
    context.TODO(),
    bson.D{{"firstName", "Xiao"}}
)

// DeleteMany() 删除多个文档
results, err := coll.DeleteMany(
    context.TODO(),
    bson.D{{"age", bson.D{{"$lte", 12}}}}
)
```


#### 聚合操作

```go

collection.Aggregate(ctx, pipline) 

// 聚合操作
$sum 计算总和	db.users.aggregate([{$group:{_id:null,totalScore:{$sum:"$score"}}}])
$avg 计算平均值	db.users.aggregate([{$group:{_id:null,totalScore:{$avg:"$score"}}}])
$max 获取最大值	db.users.aggregate([{$group:{_id:null,totalScore:{$max:"$score"}}}])
$min 获取最小值	db.users.aggregate([{$group:{_id:null,totalScore:{$min:"$score"}}}])

```

#### 条件查询

```go

// 条件查询
$eq (equal) / $ne (no equal)	
$gt / $gte 大于/大于等于	
$lt / $lte 小于/小于等于	
$in / $nin 值在指定值列表中/ 值不等于指定值列表中	
$and $or $not 与 或 非	
$exists 检查一个字段是否存在	
$expr 聚合操作，可用来比较字段
```

#### 排序/分页

```go
cursor, err := coll.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{{"age", 1}}.SetSkip(4).SetLimit(2))

```