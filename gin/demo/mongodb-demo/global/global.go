package global

import "go.mongodb.org/mongo-driver/mongo"

var (
	MongoCli *mongo.Client
)

const (
	// 数据库连接地址
	MongoUrl            = "mongodb://localhost:27017"
	// 数据库名称
	MongoDb             = "myblog"
	// 用户集合/表
	UsersCollection     = "user"
	// 文章集合/表
	ArticlesCollection  = "articles"
	// 管理员集合/表
	RootusersCollection = "rootusers"
)
