package app

import (
	"xs-bbs/internal/app/user"
	"xs-bbs/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Entities = []interface{}{
	user.Entity, 
	// community.Entity, 
	// post.Entity,
}

func Http(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	if err := db.AutoMigrate(Entities...); err != nil {
		zap.L().Error("auto migrate  tables error", zap.Error(err))
	}

	return router.NewHttpServer(
		user.Build(db),
		// community.Build(db),
		// post.Build(db, rdb),
	)
}
