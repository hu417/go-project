package global

import (
	"gorm-demo/config"

	"gorm.io/gorm"
)

var (
	Conf *config.Conf
	DB   *gorm.DB
)
