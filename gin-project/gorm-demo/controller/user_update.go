package controller

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm-demo/global"
	"gorm-demo/model"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 全局更新
func UpdateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 开启一个会话，将全局更新配置设为 true
	dbSession := global.DB.Session(&gorm.Session{
		AllowGlobalUpdate: true,
	})

	// 全局更新 age 和 name 字段
	// UPDATE `po` SET age = 0, name = ""
	resDB := dbSession.WithContext(ctx).Updates(&user)
	if resDB.Error != nil {
		log.Fatal(resDB.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": resDB.Error.Error(),
		})
		return
	}

	// 影响行数 —> x
	log.Printf("rows affected: %d", resDB.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// 条件更新
func UpdateUserByName(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 批量更新，po 中所有显式声明的字段
	// UPDATE `po` SET age = 0, name = "" WHERE age > 10;
	// 使用select 限定只更新 age 字段
	// 使用Omit 限定更新时忽略 age 字段
	resDB := global.DB.WithContext(ctx).Where("age > ?", 10).Select("age").Updates(&user)
	if resDB.Error != nil {
		log.Fatal(resDB.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": resDB.Error.Error(),
		})
		return
	}

	// 影响行数 —> x
	log.Printf("rows affected: %d", resDB.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// 批量更新：表达式更新
func UpdateUserByNameIsExpr(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// UPDATE `po` SET age = age * 2 + 1 WHERE id = 1
	resDB := global.DB.WithContext(ctx).Table("po").Where("id = ?", 1).UpdateColumn("age", gorm.Expr("age * ? + ?", 2, 1))
	if resDB.Error != nil {
		log.Fatal(resDB.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": resDB.Error.Error(),
		})
		return
	}

	// 影响行数 —> 1
	log.Printf("rows affected: %d", resDB.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// json列更新
func UpdateUserByNameIsJson(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// UPDATE `po` SET extra = json_insert(extra,"$.key","value") WHERE id = 1
	resDB := global.DB.WithContext(ctx).Where("id = ?", 1).UpdateColumn("extra", datatypes.JSONSet("extra").Set("key", "value"))
	if resDB.Error != nil {
		log.Fatal(resDB.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": resDB.Error.Error(),
		})
		return
	}

	// 影响行数 —> 1
	log.Printf("rows affected: %d", resDB.RowsAffected)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// save 方法
func UpdateUserByNameIsSave(c *gin.Context) {

	uid := c.Query("uid")
	// if uid == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "uid is empty",
	// 	})
	// 	return
	// }
	list := strings.Split(uid, ",")
	if len(list) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "uid is empty",
		})
		return
	}

	var user model.User
	var users []model.User
	for _, id := range list {
		// idInt, err := strconv.Atoi(id)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "uid is empty",
			})
			return
		}

		user.Uid = id

		users = append(users, user)
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 首先查出对应的数据
	ctxDB := global.DB.WithContext(ctx)
	if err := ctxDB.Scan(&users).Error; err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 更新数据
	for _, user := range users {
		user.Age += 100
	}

	// 将更新后的数据存储到数据库
	if err := ctxDB.Save(&users); err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
