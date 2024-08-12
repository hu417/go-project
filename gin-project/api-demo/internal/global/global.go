package global

import (
	"api-demo/internal/config"
	"api-demo/internal/crontab"
	"api-demo/internal/event"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config          *config.Config
	DB              *gorm.DB
	Logger          *zap.Logger
	EventDispatcher *event.Dispatcher
	Crontab         *crontab.Crontab
)
