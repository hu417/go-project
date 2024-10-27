package global

import (
	"gin-rbac/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Log    *zap.SugaredLogger
)
