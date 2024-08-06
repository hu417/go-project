package global

import (
	"gorm-test/config"

	"gorm.io/gorm"
)

var (
	Conf *config.Conf
	DB   *gorm.DB
)
