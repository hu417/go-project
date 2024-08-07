package global

import (
	"gin-base/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Conf *config.Config
	Log  *zap.Logger
	DB   *gorm.DB
)
