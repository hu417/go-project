package global

import (
	"blue-bell/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	ut "github.com/go-playground/universal-translator"
)

var (
	Conf     *config.App
	DB       *gorm.DB
	RedisCli *redis.Client
	// 定义一个全局翻译器T
	Trans ut.Translator
)
