package global

import (
	"ginblog/config"

	"gorm.io/gorm"
)

var (
	Conf *config.Config
	DB   *gorm.DB
)
