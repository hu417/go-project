package global

import (
	"viper-demo/config"

	"gorm.io/gorm"
)

var (
	Conf *config.Conf
	DB   *gorm.DB
)
