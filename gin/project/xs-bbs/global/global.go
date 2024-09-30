package global

import (
	"xs-bbs/pkg/conf"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Conf        *conf.Config
	DB          *gorm.DB
	RedisClient *redis.Client
	Log         *zap.SugaredLogger
)
