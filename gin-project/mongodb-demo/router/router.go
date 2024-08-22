package router

import (
	"mongodb-demo/api"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	// 注册路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("/user", api.CreateUser)               // 创建用户
		v1.POST("/users", api.CreateUsers)             // 创建多个用户
		v1.GET("/user", api.FindUser)                  // 查询用户
		v1.GET("/users", api.FindallUser)              // 查询所有用户
		v1.PUT("/user", api.UpdateUser)                // 根据id更新用户
		v1.PUT("/users", api.UpdateAllUser)            // 更新多个用户
		v1.PUT("/users/update", api.FindOneAndUpdate)  // 根据name更新用户
		v1.DELETE("/user", api.DeleteUser)             // 根据name删除用户
		v1.DELETE("/users", api.DeleteUsers)           // 删除多个用户
		v1.GET("/users/count", api.CountForCollection) // 查询文档数量     // 统计用户数量
		v1.GET("/users/aggregate", api.AggregateUser)  // 聚合查询并排序
	}

	return r
}
