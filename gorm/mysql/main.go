package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 表字段映射结构体
type List struct {
	gorm.Model
	Name  string `gorm:"type:varchar(20);not null" json:"name" binding:"required"`
	Age   string `gorm:"type:varchar(2);not null" json:"age" binding:"required"`
	Phone string `gorm:"type:varchar(11);not null" json:"phone" binding:"required"`
	Email string `gorm:"type:varchar(20);not null" json:"email" binding:"required"`
	City  string `gorm:"type:varchar(20);not null" json:"city" binding:"required"`
	// gorm: 字段类型; json数据对应字段; binding请求绑定参数
}

// TableName 指定表名
// func (User) TableName() string {
// 	return "user"
// }

func main() {
	// dsn := "root:123456@tcp(10.0.0.91:3306)/curd-list?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	dsn := "root:123456@tcp(100.84.144.92:3306)/devops?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//DryRun: true,
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "T-", // 表前缀
			SingularTable: true, // 禁用表名复数
		},
	})
	if err != nil {
		fmt.Printf("数据库连接错误,err= %s", err)
	}
	// 连接池
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Print(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 设置表的存储引擎为InnoDB, 自动迁移数据库结构
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&List{}); err != nil {
		panic("数据库迁移失败")
	}

	//============= CURD ==============
	// 1、新增
	if err := Insert(db); err != nil {
		log.Fatalf("创建用户失败, err: %v", err)
	} else {
		log.Default().Println("创建用户成功")
	}

	// 2、查询
	if list, err := Select(db); err != nil {
		log.Fatalf("查询用户失败, err: %v", err)
	} else {
		log.Default().Printf("user => %v", list)
	}

	// 3、分页查询
	db = db.Model(&List{}).Where("name = ?", "Tom")
	if users, err := SelectLimit(0, 3, db); err != nil {
		log.Fatalf("查询用户失败, err: %v", err)
	} else {
		// log.Default().Printf("user => %v", users)
		for _, v := range users {
			log.Default().Printf("id => %v", v.ID)
		}
	}

	// 4、修改
	if err := Update(db); err != nil {
		log.Fatalf("更新用户失败, err: %v", err)
	} else {
		log.Default().Println("更新用户成功")
	}

	// 5、删除
	if err := Delete(db); err != nil {
		log.Fatalf("删除用户失败, err: %v", err)
	} else {
		log.Default().Println("删除用户成功")
	}
}

// 新增
func Insert(db *gorm.DB) error {
	// 创建操作
	// user := List{Name: "John", Age: "30", Phone: "13211113333", Email: "john@example.com", City: "shenzhen"}
	user := []*List{
		&List{Name: "Jony", Age: "19", Phone: "13211113334", Email: "john@example.com", City: "guangzhou"},
		&List{Name: "Tom", Age: "20", Phone: "13211113335", Email: "john@example.com", City: "shanghai"},
	}
	result := db.Debug().Create(&user)
	if result.Error != nil {
		log.Fatalf("创建用户失败, err: %v", result.Error)
		return result.Error
	}
	log.Default().Printf("影响行数: %v", result.RowsAffected)
	return nil
}

// 查询
func Select(db *gorm.DB) ([]List, error) {
	// 读取操作
	var users []List
	// 查询单个对象，返回第一个
	// result := db.Debug().First(&users)
	// 查询所有对象
	// result := db.Debug().Where("name = ?","Jony").Find(&users)
	result := db.Debug().Where("`deleted_at` IS NULL").Find(&users)
	if result.Error != nil {
		log.Fatalf("查询用户失败, err: %v", result.Error)
		return nil, result.Error
	}
	log.Default().Printf("总共有: %v 行数据", result.RowsAffected)
	//fmt.Println("所有用户：", users)
	return users, nil
}

// 分页查询
func SelectLimit(pageIndex, pageSize int, db *gorm.DB) ([]List, error) {
	var users []List

	// 获取总数
	if result := db.Find(&users); result.Error != nil {
		return nil, result.Error
	} else {
		count := result.RowsAffected
		log.Default().Printf("总计 => %v 行数据", count)
	}

	// 分页
	if pageIndex > 0 && pageSize > 0 {
		offset := (pageIndex - 1) * pageSize
		db = db.Model(&List{}).Offset(offset).Limit(pageSize)
	}

	if err := db.Find(&users).Error; err != nil {
		fmt.Println(err.Error())
	}
	return users, nil

}

// 更新
func Update(db *gorm.DB) error {
	// 更新操作
	var user List
	db.First(&user)
	user.Name = "hh"
	user.Age = "28"
	result := db.Debug().Where("id = ?", 10).Save(&user)
	if result.Error != nil {
		log.Fatalf("更新用户失败, err: %v", result.Error)
		return result.Error
	}
	log.Default().Printf("影响行数: %v", result.RowsAffected)
	return nil
}

// 删除
func Delete(db *gorm.DB) error {
	// 删除操作
	var user List
	result := db.Debug().Where("age = ?", "28").Delete(&user)
	if result.Error != nil {
		log.Fatalf("删除用户失败, err: %v", result.Error)
		return result.Error
	}
	log.Default().Printf("影响行数: %v", result.RowsAffected)
	return nil
}
