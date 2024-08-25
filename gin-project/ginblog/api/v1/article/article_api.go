package article

import "github.com/gin-gonic/gin"

type Article struct {
}

type ArticleApiInterface interface {
	Add(ctx *gin.Context)
	EditById(ctx *gin.Context)
	DeleteById(ctx *gin.Context)
	FindByName(ctx *gin.Context)
	FindList(ctx *gin.Context)
}

func NewArticle() *Article {
	return &Article{}
}
