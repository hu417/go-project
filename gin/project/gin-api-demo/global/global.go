package global

import (
	"gin-api-demo/config"

	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Conf *config.Config
	Log  *zap.Logger
	DB   *gorm.DB
	// 定义一个全局翻译器T
	Trans ut.Translator
	Redis *redis.Client
)
