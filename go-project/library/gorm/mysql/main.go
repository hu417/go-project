package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-faker/faker/v4" // 造数据的包
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 表字段映射结构体
type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100);not null" json:"name" binding:"required" faker:"name,len=10"`
	Password  string `gorm:"type:varchar(20);" json:"password" binding:"required" faker:"password,len=20"`
	Phone     uint64 `gorm:"type:bigint;not null" json:"phone" binding:"required" faker:"phone,boundary_start=13000000000, boundary_end=13999999999"`
	Email     string `gorm:"type:varchar(20);not null" json:"email" binding:"required" faker:"email"`
	Timestamp string `gorm:"type:varchar(20);not null" json:"timestamp" binding:"required" faker:"timestamp"`
	// gorm: 字段类型; json数据对应字段; binding请求绑定参数
}

type Users struct {
	ID          string `gorm:"column:id" faker:"-"`
	Email       string `gorm:"column:email" faker:"email"`
	Password    string `gorm:"column:password" faker:"password"`
	PhoneNumber string `gorm:"column:phone_number" faker:"phone_number"`
	UserName    string `gorm:"column:username" faker:"username"`
	FirstName   string `gorm:"first_name" faker:"first_name"`
	LastName    string `gorm:"last_name" faker:"last_name"`
	Century     string `gorm:"century" faker:"century"`
	Date        string `gorm:"date" faker:"date"`
}

// 在插入记录之前，生成uuid
func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	u.ID = uint(uuid.New().ID())
	return nil
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
		// DryRun: true, // 开启后只显示sql语句但实际不执行
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
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间

	// 设置表的存储引擎为InnoDB, 自动迁移数据库结构
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}); err != nil {
		panic("数据库迁移失败")
	}

	// 全局开启debug模式,打印出执行的SQL语句以及对应的参数
	// db = db.Debug()

	//============= CURD ==============
	// 1、新增
	if err := Insert(db); err != nil {
		log.Fatalf("创建用户失败, err: %v", err)
	} else {
		log.Default().Println("创建用户成功")
	}

	// // 2、查询
	// if list, err := Select(db); err != nil {
	// 	log.Fatalf("查询用户失败, err: %v", err)
	// } else {
	// 	log.Default().Printf("user => %v", list)
	// }

	// // 3、分页查询
	// db = db.Model(&User{}).Where("name = ?", "Tom")
	// if users, err := SelectLimit(0, 3, db); err != nil {
	// 	log.Fatalf("查询用户失败, err: %v", err)
	// } else {
	// 	// log.Default().Printf("user => %v", users)
	// 	for _, v := range users {
	// 		log.Default().Printf("id => %v", v.ID)
	// 	}
	// }

	// // 4、修改
	// if err := Update(db); err != nil {
	// 	log.Fatalf("更新用户失败, err: %v", err)
	// } else {
	// 	log.Default().Println("更新用户成功")
	// }

	// // 5、删除
	// if err := Delete(db); err != nil {
	// 	log.Fatalf("删除用户失败, err: %v", err)
	// } else {
	// 	log.Default().Println("删除用户成功")
	// }
}

// 新增
func Insert(db *gorm.DB) error {
	// 创建操作
	// user := User{Name: "John", Age: "30", Phone: "13211113333", Email: "john@example.com", City: "shenzhen"}
	// user := []*User{
	// 	&User{Name: "Jony", Age: "19", Phone: "13211113334", Email: "john@example.com", City: "guangzhou"},
	// 	&User{Name: "Tom", Age: "20", Phone: "13211113335", Email: "john@example.com", City: "shanghai"},
	// }

	var tmpUser *User

	// 1、单条数据
	// 生成随机数据
	// err := faker.FakeData(&tmpUser)
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// result := db.Debug().Create(&tmpUser)

	// 2、多条数据
	var users []*User
	for i := 0; i < 10; i++ {
		//tmpUser = &User{}
		err := faker.FakeData(&tmpUser)
		if err != nil {
			log.Fatal(err.Error())
		}
		users = append(users, tmpUser)
	}

	// 指定sql语句开启debug
	result := db.Debug().Create(&users)
	if result.Error != nil {
		log.Fatalf("创建用户失败, err: %v", result.Error)
		return result.Error
	}
	log.Default().Printf("影响行数: %v", result.RowsAffected)
	return nil
}

// 查询
func Select(db *gorm.DB) ([]User, error) {
	// 读取操作
	var users []User
	// First查询单个对象，返回第一个（主键升序）
	// result := db.Debug().First(&users)

	// Last查询单个对象，返回第一个（主键排序）
	// result := db.Debug().Last(&users)

	// 推荐这种写法，根据主键查找，写法不够清晰
	// result := db.First(&users, "name = ?","Jony") // 查找 name 字段值为 Jony 的记录

	// 查询指定条件的对象
	// result := db.Debug().Where("`deleted_at` IS NULL").Find(&users)
	// if result.Error != nil {
	// 	log.Fatalf("查询用户失败, err: %v", result.Error)
	// 	return nil, result.Error
	// }
	// log.Default().Printf("总共有: %v 行数据", result.RowsAffected)
	//fmt.Println("所有用户：", users)

	// 只需要查询某些字段，可以重新定义小结构体
	type APIUser struct {
		Name  string `gorm:"type:varchar(20);not null" json:"name" binding:"required" faker:"name"`
		Email string `gorm:"type:varchar(20);not null" json:"email" binding:"required" faker:"email"`
	}
	var apiUser []APIUser
	result := db.Model(&User{}).Find(&apiUser)
	for _, user := range apiUser {
		fmt.Println(user)
	}
	fmt.Println(result.Error, result.RowsAffected)

	// Count查询
	var count int64
	result = db.Model(&User{}).Where("id > ?", 2).Count(&count)
	fmt.Println(count)
	fmt.Println(result.Error, result.RowsAffected)

	return users, nil
}

// 分页查询
func SelectLimit(pageIndex, pageSize int, db *gorm.DB) ([]User, error) {
	var users []User

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
		db = db.Model(&User{}).Offset(offset).Limit(pageSize)
	}

	if err := db.Find(&users).Error; err != nil {
		fmt.Println(err.Error())
	}
	return users, nil

}

// 更新
func Update(db *gorm.DB) error {
	//==== 更新操作
	var user User
	db.First(&user)
	user.Name = "hh"
	user.Password = "qazwsx"
	// Save会保存所有字段，即使字段是零值，如果保存的值没有主键，就会创建，否则则是更新指定记录
	result := db.Debug().Where("id = ?", 10).Save(&user)
	if result.Error != nil {
		log.Fatalf("更新用户失败, err: %v", result.Error)
		return result.Error
	}
	log.Default().Printf("影响行数: %v", result.RowsAffected)

	// 更新单个列
	result = db.Model(&User{}).Where("username = ?", "jaQlaFs").Update("name", "zhangsan")
	fmt.Println(result.Error, result.RowsAffected)

	// 更新多个列
	result = db.Model(&User{}).Where("username = ?", "zhangsan").Updates(User{Name: "zhangsan2", Password: "qazwsx"})
	fmt.Println(result.Error, result.RowsAffected)

	return nil

}

// 删除
func Delete(db *gorm.DB) error {
	//==== 删除操作
	// 指定匹配字段删除数据
	var user User
	// result = db.Delete(&User{}, "username = ?", "NJrauTj")
	result := db.Debug().Where("age = ?", "28").Delete(&user)
	if result.Error != nil {
		log.Fatalf("删除用户失败, err: %v", result.Error)
		return result.Error
	}
	log.Default().Printf("影响行数: %v", result.RowsAffected)

	// 批量删除的两种方式
	result = db.Where("email like ?", "%.com%").Delete(&User{})
	fmt.Println(result.Error, result.RowsAffected)

	result = db.Delete(&User{}, "email like ?", "%.com%")
	fmt.Println(result.Error, result.RowsAffected)

	return nil
}

// execSQL 执行原生SQL语句
// 1、Scan函数是会将查询SQL的结果映射到定义的变量，如果不需要返回查询结果可以直接使用Exec函数执行原生SQL；
// 2、DryRun模式，可以直接生成SQL机器参数，但是不会直接执行；
func execSQL(db *gorm.DB) {

	// 将查询SQL的结果映射到指定的单个变量中
	var oneUser User
	result := db.Raw("SELECT * FROM user LIMIT 1").Scan(&oneUser)
	fmt.Println(oneUser)
	fmt.Println(result.Error, result.RowsAffected)

	// 将查询SQL的批量结果映射到列表中
	var users []User
	result = db.Raw("SELECT * FROM user").Scan(&users)
	for _, user := range users {
		fmt.Println(user)
	}
	fmt.Println(result.Error, result.RowsAffected)

	var updateUser User
	result = db.Raw("UPDATE users SET username = ? where id = ?", "toms jobs", "ab6f089b-3272-49b5-858f-a93ed5a43b4f").Scan(&updateUser)
	fmt.Println(updateUser)
	fmt.Println(result.Error, result.RowsAffected)

	// 直接通过Exec函数执行Update操作，不返回任何查询结果？
	result = db.Exec("UPDATE user SET username = ? where id = ?", "toms jobs", "ab6f089b-3272-49b5-858f-a93ed5a43b4f")
	fmt.Println(result.Error, result.RowsAffected)

	// DryRun模式，在不执行的情况下生成SQL及其参数，可以用于准备或测试的SQL
	// 只需要查询某些字段，可以重新定义小结构体
	type APIUser struct {
		Name  string `gorm:"type:varchar(20);not null" json:"name" binding:"required" faker:"name"`
		Email string `gorm:"type:varchar(20);not null" json:"email" binding:"required" faker:"email"`
	}
	var tmpUsers []APIUser
	stmt := db.Session(&gorm.Session{DryRun: true}).Model(&User{}).Find(&tmpUsers).Statement
	fmt.Println(stmt.SQL.String())
	fmt.Println(stmt.Vars)
}
