package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm-demo/global"
	"gorm-demo/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*

gorm 中，First、Last、Take、Find 方法都可以用于查询单条记录.
前三个方法的特点是，倘若未查询到指定记录，则会报错 gorm.ErrRecordNotFound；
最后一个方法的语义更软一些，即便没有查到指定记录，也不会返回错误

- 条件查询：where + and / or 或者 嵌套
// WHERE age = 1 AND name = 'xu'
   db.Where("age = 1").Where("name = ?",xu)

// WHERE age = 1 OR name = 'xu'
   db.Where("age = 1").Or("name = ?","xu") 

// WHERE (age = 1 AND name = 'xu') OR (age = 2 AND name  = 'x')
   db.Where(db.Where("age = 1").Where("name = ?","xu")).Or(db.Where("age = 2").Where("name = ?","x"))


*/

// 查询：返回一条数据, 使用First
func SelectUserByNameIsFirst(c *gin.Context) {

	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	var user model.User
	user.Name = name

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	re := global.DB.WithContext(ctx).Select("name").First(&user)

	if re.RowsAffected > 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  "查询成功",
			"data": user,
		})
	}
	// SELECT * FROM `po` WHERE deleted_at IS NULL ORDER BY id ASC LIMIT 1
	if err := global.DB.WithContext(ctx).First(&user).Error; err != nil {
		/*
			// 通过 Select 方法声明只返回特定的列, 	只返回 age 列的数据
			db.WithContext(ctx).Select("age").Where("id = ?",999).First(&po).Error
		*/
		log.Fatal(err)

		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 20001,
				"msg":  "系统错误",
				"data": nil,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "查询失败,数据不存在",
			"data": nil,
		})
		return
	}

	log.Printf("user: %+v", user)

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"msg":  "查询成功",
		"data": user,
	})
}

// 查询：返回1条数据,使用Last
func SelectUserByNameIsLast(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	var user model.User
	user.Name = name

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// SELECT * FROM `po` WHERE deleted_at IS NULL ORDER BY id ASC LIMIT 1
	// 取 age > 10 的记录中主键最大的记录

	// SELECT * FROM `po` WHERE age > 10 AND deleted_at IS NULL ORDER BY id DESC imit 1
	if err := global.DB.WithContext(ctx).Where("age > ?", 10).Last(&user).Error; err != nil {
		log.Fatal(err)

		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 20001,
				"msg":  "系统错误",
				"data": nil,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "查询失败,数据不存在",
			"data": nil,
		})
		return
	}

	log.Printf("user: %+v", user)

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"msg":  "查询成功",
		"data": user,
	})
}

// 查询: 返回1条数据,使用Take
func SelectUserByNameIsTake(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	var user model.User
	user.Name = name

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 取 id < 10 的记录中随机一条记录返回

	// SELECT * FROM `po` WHERE id < 10  AND deleted_at IS NULL LIMIT 1
	if err := global.DB.WithContext(ctx).Where("id < ?", 10).Take(&user).Error; err != nil {
		log.Fatal(err)

		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 20001,
				"msg":  "系统错误",
				"data": nil,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "查询失败,数据不存在",
			"data": nil,
		})
		return
	}

	log.Printf("user: %+v", user)

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"msg":  "查询成功",
		"data": user,
	})
}

// 查询: 返回多条数据,使用Find
// 从满足条件的数据记录中随机返回一条，即便没有找到记录，也不会抛出错误
func SelectUserByNameIsFind(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "转换失败",
			"data": nil,
		})
		return
	}
	var user []*model.User
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 通过 find 检索记录，找不到满足条件的记录时，也不会返回错误
	// SELECT * FROM `po` WHERE id = 999 AND deleted_at IS NULL
	if err := global.DB.WithContext(ctx).Where("id = ?", idUint).Find(&user).Error; err != nil {
		/*
			// 批量查询
			var users []*model.User
			// SELECT * FROM `user` WHERE age > 1 AND deleted_at IS NULL
			db.WithContext(ctx).Where("age > ?", 1).Find(&pos).Error

		*/
		log.Fatal(err)

		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "查询失败,	系统错误",
			"data": nil,
		})
	}
	log.Printf("user: %+v", user)
	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"msg":  "查询成功",
		"data": user,
	})
}

// 批量查询: 使用scan
func SelectUserByNameIsScan(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "转换失败",
			"data": nil,
		})
		return
	}
	var user []*model.User
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

    // SELECT * FROM `po` WHERE age > 1 AND deleted_at IS NULL  
    if err := global.DB.WithContext(ctx).Table("po").Where("id > ?", idUint).Scan(&user).Error; err != nil {
        log.Fatal(err)

		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "查询失败,	系统错误",
			"data": nil,
		})
        return
    }

    for _, u := range user {
        log.Printf("user: %+v\n", u)
    }

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"msg":  "查询成功",
		"data": user,
	})

}

// 数量统计
func SelectUserByNameIsCount(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "转换失败",
			"data": nil,
		})
		return
	}
	var user []*model.User
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	var cnt int64
    // SELECT COUNT(*) FROM `po` WHERE age > 10 AND deleted_at IS NULL
    if err := global.DB.WithContext(ctx).Model(user).Where("id > ?", idUint).Count(&cnt); err != nil {
        log.Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "查询失败,	系统错误",
			"data": nil,
		})
        return
    }

	// 打印结果
	log.Printf("user: %+v \n cnt: %d", user,cnt)

	// 构造返回数据
	data := struct {
		Total int64
		User  []*model.User
	}{
		Total: cnt,
		User:  user,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"msg":  "查询成功",
		"data": data,
	})
}


// 分组求和
func SelectUserByNameIsGroup(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "转换失败",
			"data": nil,
		})
		return
	}
	var users []*UserRecord
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

    // SELECT user_id, sum(amount) AS amount FROM `user_record` WHERE id < 100 AND deleted_at IS NULL GROUP BY user_id
    resDB := global.DB.WithContext(ctx).Table("user").Select("uid", "sum(amount) AS amount").
        Where("id < ?", idUint).Group("uid").Scan(&users)
    if resDB.Error != nil {
        log.Fatal(resDB.Error)
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "查询失败,	系统错误",
			"data": nil,
		})
        return
    }


    for _, user := range users {
        log.Printf("user: %+v\n", user)
    }

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"msg":  "查询成功",
		"data": users,
	})
}

// 子查询
func SelectUserByNameIsJoin(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "转换失败",
			"data": nil,
		})
		return
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// UPDATE `user` SET amount = (SELECT amount FROM `user_record` WHERE user_id = 1000 ORDER BY id DESC limit 1) WHERE user_id = 100 
    subQuery := global.DB.Table("user").Select("uid").Where("id = ?", idUint)
    
    resDB := global.DB.WithContext(ctx).Table("user").Where("uid = ?", idUint).UpdateColumn("uid", subQuery)
    if resDB.Error != nil {
        log.Fatal(resDB.Error)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "更新失败",
			"data": nil,
		})
        return
    }

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "更新成功",
		"data": nil,
	})
}

// 排序偏移
func SelectUserByNameIsOffset(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "转换失败",
			"data": nil,
		})
		return
	}
	var users []*User
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// SELECT * FROM `po` WHERE id > 10 AND deleted_at is NULL ORDER BY age DESC LIMIT 2 OFFSET 10
    if err := global.DB.WithContext(ctx).Table("user").Where("id > ?", idUint).Order("age DESC").Limit(2).Offset(10).Scan(&users).Error; err != nil {
        log.Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 20001,
			"msg":  "查询失败",
			"data": nil,
		})
        return
    }

	// 输出结果
    for _, user := range users {
        log.Printf("user: %+v\n", user)
    }

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"msg":  "查询成功",
		"data": users,
	})
}