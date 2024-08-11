package global

import (
	"api-demo/internal/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
	Logger *zap.Logger
)
