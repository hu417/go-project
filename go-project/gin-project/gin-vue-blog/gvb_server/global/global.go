package global

import (
	"gvb_server/config"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config   *config.Config
	DB       *gorm.DB
	Logger   *zap.SugaredLogger
	Router   *gin.Engine
	ESClient *elastic.Client
)
