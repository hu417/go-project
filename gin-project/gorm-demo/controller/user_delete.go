package controller

import (
	"context"
	"log"
	"time"

	"gorm-demo/global"
	"gorm-demo/model"

	"github.com/gin-gonic/gin"
)

// 删除
func DeleteUser(c *gin.Context) {

	uid := c.Param("uid")
	if uid == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "uid is empty",
		})
		return
	}

	var user *model.User
	user.Uid = uid

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 软删除
	// UPDATE `po` SET deleted_at = /* current unix second */ WHERE id = 1
	resDB := global.DB.WithContext(ctx).Delete(&user)
	if resDB.Error != nil {
		log.Fatal(resDB.Error)
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "delete user failed",
		})
		return
	}

	// 影响行数 —> 1
	log.Panicf("rows affected: %d", resDB.RowsAffected)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "delete user success",
	})

}

// 批量删除：where
func DeleteUsers(c *gin.Context) {
	var users []model.User

	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "bind json failed",
		})
		return
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 批量软删除所有 age > 10 的记录
	// UPDATE `po` SET deleted_at = /* current unix second */ WHERE age > 10
	resDB := global.DB.WithContext(ctx).Where("age > ?", 10).Delete(&users)
	if resDB.Error != nil {
		log.Fatal(resDB.Error)
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "delete users failed",
		})
		return
	}

	// 影响行数 —> x
	log.Panicf("rows affected: %d", resDB.RowsAffected)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "delete users success",
	})
}
