package router

import (
	v1 "gin-es/api/v1"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 注册路由
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})
	api := r.Group("/api/v1")
	{
		// 索引相关接口
		api.POST("/index", v1.NewEsApi().CreateIndex) // 创建索引
		api.GET("/index", v1.NewEsApi().GetIndex)     // 获取索引列表
		api.GET("/index/:name", v1.NewEsApi().GetIndexByName)  // 获取索引
		api.DELETE("/index/:name", v1.NewEsApi().DeleteIndex) // 删除索引
		// 文档相关接口
		api.POST("/index/:name/doc", v1.NewEsApi().CreateDoc)  // 创建文档
		api.POST("/index/:name/docs", v1.NewEsApi().CreateDocs)  // 批量创建文档
		api.GET("/index/:name/doc", v1.NewEsApi().GetDoc) // 获取文档列表
		api.GET("/index/:name/doc/:id", v1.NewEsApi().GetDocById) // 获取文档
		api.PUT("/index/:name/doc/:id", v1.NewEsApi().UpdateDocById) // 更新文档
		api.DELETE("/index/:name/doc/:id", v1.NewEsApi().DeleteDocById)  // 删除文档
		api.DELETE("/index/:name/doc", v1.NewEsApi().DeleteDocByQuery)  // 删除匹配文档
		api.POST("/index/:name/search", v1.NewEsApi().SearchDoc) // 搜索文档
		api.GET("/index/:name/search/:id", v1.NewEsApi().ExisDoc) // 聚合搜索
	}

	return r
}
