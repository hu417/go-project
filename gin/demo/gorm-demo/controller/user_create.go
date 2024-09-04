package controller

import (
	"context"
	"log"
	"time"

	"gorm-demo/global"

	"github.com/gin-gonic/gin"
)

// 创建用户 
func CreateUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 执行创建操作
	// INSERT INTO `po` (`age`) VALUES (0);
	resDB := global.DB.WithContext(ctx).Create(&user)
	if resDB.Error != nil {
		log.Fatal(resDB.Error)
		ctx.JSON(500, gin.H{
			"error": resDB.Error.Error(),
		})
		return
	}

	// 影响行数 -> 1
	log.Printf("rows affected: %d", resDB.RowsAffected)
	// 结果输出
	log.Printf("user: %+v", user)
	ctx.JSON(200, gin.H{
		"user": user,
	})
}

// 批量创建
func CreateUsers(c *gin.Context) {
	var users []User
	if err := c.ShouldBind(&users); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 批量创建
	// 批量创建时会根据 gorm.Config 中的 CreateBatchSize 进行分批创建操作
	resDB := global.DB.WithContext(ctx).Table("po").Create(&users)
	if resDB.Error != nil {
		log.Fatal(resDB.Error)
		c.JSON(500, gin.H{
			"error": resDB.Error.Error(),
		})
		return
	}

	// 输出影响行数 -> 2
	log.Panicf("rows affected: %d", resDB.RowsAffected)

	// 打印各 po，输出其主键
	for _, po := range users {
		log.Panicf("po: %+v\n", po)
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}
