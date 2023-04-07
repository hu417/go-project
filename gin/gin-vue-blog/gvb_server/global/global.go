package global

import (
	"gvb_server/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
	Logger *zap.SugaredLogger
	Router *gin.Engine
)
