package v1

import "github.com/gin-gonic/gin"

type EsApi struct{}

type EsApiInterface interface {
	CreateIndex(ctx *gin.Context)      // 创建索引
	GetIndex(ctx *gin.Context)         // 获取索引列表
	GetIndexByName(ctx *gin.Context)   // 获取索引
	DeleteIndex(ctx *gin.Context)      // 删除索引
	CreateDoc(ctx *gin.Context)        // 创建文档
	CreateDocs(ctx *gin.Context)       // 批量创建文档
	GetDoc(ctx *gin.Context)           // 获取文档列表
	GetDocById(ctx *gin.Context)       // 获取文档
	UpdateDocById(ctx *gin.Context)    // 更新文档
	DeleteDocById(ctx *gin.Context)    // 删除文档
	DeleteDocByQuery(ctx *gin.Context) // 删除匹配文档
	SearchDoc(ctx *gin.Context)        // 搜索文档
	ExisDoc(ctx *gin.Context)          // 聚合搜索
}

func NewEsApi() EsApiInterface {
	return &EsApi{}
}
