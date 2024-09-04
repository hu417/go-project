package global

import (
	"gin-base/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Conf *config.Config
	Log  *zap.Logger
	DB   *gorm.DB
	Rds *redis.Client
)
