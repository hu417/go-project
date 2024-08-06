package controller

import (
	"context"
	"log"
	"time"

	"gorm-demo/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 事务
func Tx(c *gin.Context) {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 需要包含在事务中执行的闭包函数
	do := func(tx *gorm.DB) error {
		// do something ...
		return nil
	}

	// 开启事务
	// BEGIN
	// OPERATE...
	// COMMIT/ROLLBACK
	if err := global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// do some preprocess ...
		// do ...
		err := do(tx)
		// do some postprocess ...
		return err
	}); err != nil {
		log.Fatal(err)
	}
}

