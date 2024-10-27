# GORM框架学习

GORM是Go语言中的一个强大且灵活的对象关系映射（ORM）库。它提供了一种简便的方法来处理数据库操作，支持多种数据库类型如MySQL、PostgreSQL、SQLite和SQL Server。本文将详细介绍GORM的功能、使用方法以及一些高级特性。

官方文档：https://gorm.cn/zh_CN/docs/


## GORM简介

GORM是一个适用于Go语言的ORM库，它提供了丰富的功能，支持常见的数据库操作，并且具有良好的扩展性和灵活性。GORM的设计理念是尽可能简化开发者对数据库的操作，同时提供强大的查询构造器、关系映射和迁移工具。

GORM的主要特性包括：

- 全自动的ORM功能：支持结构体与数据库表的自动映射。
- 丰富的查询构造器：支持链式调用构造查询。
- 强大的迁移工具：支持数据库表的自动创建、更新和删除。
- 支持多种数据库：包括MySQL、PostgreSQL、SQLite和SQL Server。
- 灵活的钩子：提供多种钩子函数，允许在CRUD操作的不同阶段执行自定义逻辑。

## 安装GORM

在使用GORM之前，需要先安装它。可以使用以下命令通过go get工具安装GORM：

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

上面的命令不仅安装了GORM库，还安装了MySQL驱动。如果你使用其他数据库，需要安装相应的驱动，例如：

```go
go get -u gorm.io/driver/postgres
go get -u gorm.io/driver/sqlite
go get -u gorm.io/driver/sqlserver
```

## 连接数据库

### 基本配置

示例：连接到MySQL数据库

```go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
)

func main() {
    studentname := "root"  //账号
    password := "root"  //密码
    host := "127.0.0.1" //数据库地址，可以是Ip或者域名
    port := 3306        //数据库端口
    Dbname := "gorm"   //数据库名
    timeout := "10s"    //连接超时，10秒

    // root:root@tcp(127.0.0.1:3306)/gorm?
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", 
           studentname, 
           password, 
           host, 
           port, 
           Dbname, 
           timeout,
    )
    //连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
    db, err := gorm.Open(mysql.Open(dsn))
    if err != nil {
      panic("连接数据库失败, error=" + err.Error())
    }
    // 连接成功
    fmt.Println(db)


    // 使用 db 进行数据库操作
}
```

在上面的代码中，我们首先构造了一个数据源名称`（DSN）`，然后使用`gorm.Open`方法连接到MySQL数据库。`gorm.Open`方法返回一个`*gorm.DB`类型的对象，之后可以使用这个对象进行各种数据库操作。

### 高级配置

#### 跳过默认事务

为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它，这样可以获得60%的性能提升

```go
db, err := gorm.Open(mysql.Open("gorm.db"), &gorm.Config{
  SkipDefaultTransaction: true,
})

```

#### 命名策略

gorm采用的命名策略是，表名是蛇形复数，字段名是蛇形单数

例如

```go
type Student struct {
  Name      string
  Age       int
  MyStudent string
}


```

gorm会为我们这样生成表结构

```go
CREATE TABLE `students` (`name` longtext,`age` bigint,`my_student` longtext)
```

也可以修改这些策略

```go

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
  NamingStrategy: schema.NamingStrategy{
    TablePrefix:   "f_",  // 表名前缀
    SingularTable: false, // 单数表名
    NoLowerCase:   false, // 关闭小写转换
  },
})
```

#### 显示日志

gorm的默认日志是只打印错误和慢SQL

我们可以自己设置

```go
var mysqlLogger logger.Interface
// 要显示的日志等级
mysqlLogger = logger.Default.LogMode(logger.Info)
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
  Logger: mysqlLogger,
})
```

如果你想自定义日志的显示

那么可以使用如下代码

```go
newLogger := logger.New(
  log.New(os.Stdout, "\r\n", log.LstdFlags), // （日志输出的目标，前缀和日志包含的内容）
  logger.Config{
    SlowThreshold:             time.Second, // 慢 SQL 阈值
    LogLevel:                  logger.Info, // 日志级别
    IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
    Colorful:                  true,        // 使用彩色打印
  },
)

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
  Logger: newLogger,
})
```

部分展示日志

```go
var model Student
session := DB.Session(&gorm.Session{Logger: newLogger})
session.First(&model)
// SELECT * FROM `students` ORDER BY `students`.`name` LIMIT 1
```

如果只想某些语句显示日志

```go
DB.Debug().First(&model)
```

#### 超时设置

```go
// 超时控制
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
db := DB.WithContext(ctx)
```



### 参考连接

配置文件

```go
// config.yaml
mysql:
  host: 127.0.0.1 #地址
  port: "3306" #端口
  config: charset=utf8mb4&parseTime=True&loc=Local #配置
  db-name: gva #数据库名称
  username: root #账号
  password: root #密码
  max-idle-conns: 0 #最大空闲连接数
  max-open-conns: 0 #最大连接数
  log-mode: "" #是否开启Gorm全局日志
  log-zap: false #是否打印日志到zap

// mysql.go
package config

import (
	"catering/pkg/e"
)

type Mysql struct {
	Host         string `json:"host" yaml:"host"`                   // 服务器地址
	Port         string `json:"port" yaml:"port"`                   // 端口
	Config       string `json:"config" yaml:"config"`               // 高级配置
	Dbname       string `json:"dbname" yaml:"db-name"`              // 数据库名
	Username     string `json:"username" yaml:"username"`           // 数据库用户名
	Password     string `json:"password" yaml:"password"`           // 数据库密码
	MaxIdleConns int    `json:"maxIdleConns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `json:"maxOpenConns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	MaxLifeTime  int    `json:"maxLifeTime" yaml:"max-life-time"`
	LogMode      string `json:"logMode" yaml:"log-mode"` // 是否开启Gorm全局日志
	LogZap       bool   `json:"logZap" yaml:"log-zap"`   // 是否通过zap写入日志文件
}

func (m *Mysql) Check() error {
	if m.Username == "" || m.Dbname == "" {
		return e.ErrMysqlConfigCheckFail
	}
	return nil
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}
```

全局配置

```go
// global.go
package global

import (
	"catering/config"

	"gorm.io/gorm"
)

var (
  Cnf *config.Config
	DB  *gorm.DB
)
```



初始化

```go
// mysql.go
package initialize

import (
	"catering/global"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Gorm 初始化数据库并产生数据库全局变量
// Author SliverHorn
func InitGorm() (db *gorm.DB){
  // 判断db类型
	switch global.Config.System.DbType {
	case "mysql":
		db = GormMysql()
	default:
		db = GormMysql()
	}

}

func GormMysql() *gorm.DB {
	//获取配置文件的配置
	cfg := global.Config.Mysql
  
	//检查配置
	if err := cfg.Check(); err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       cfg.Dsn(), // DSN data source name
		DefaultStringSize:         255,       // string 类型字段的默认长度
		DisableDatetimePrecision:  true,      // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,      // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,      // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,     // 根据当前 MySQL 版本自动配置
	}
	//连接数据库
	db, err := gorm.Open(mysql.New(mysqlConfig), getGormConfig())
	if err != nil {
		return nil
	}
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	// 设置默认值
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	Options(sqlDB, WithMaxIdelConns(cfg.MaxIdleConns), WithMaxOpenConns(cfg.MaxOpenConns), WithMaxLifeTime(cfg.MaxLifeTime))
	return db
}

func getGormConfig() *gorm.Config {
	//禁用外键约束
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}

	//NewWriter 对log.New函数的再次封装，从而实现是否通过zap打印日志
	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	//设置logger的日志输出等级
	switch global.Config.Mysql.LogMode {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}

type writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf 格式化打印日志
func (w *writer) Printf(message string, data ...interface{}) {
	var logZap bool
	switch global.Config.System.DbType {
	case "mysql":
		logZap = global.Config.Mysql.LogZap
	}
	//通过zap打印日志，或者其他
	if logZap {
		global.Log.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}

//Option设计模式封装mysql的额外配置
type Option func(m *sql.DB)

func WithMaxIdelConns(idle int) Option {
	return func(m *sql.DB) {
		if idle == 0 {
			return
		}
		m.SetMaxIdleConns(idle)
	}
}

func WithMaxOpenConns(open int) Option {
	return func(m *sql.DB) {
		if open == 0 {
			return
		}
		m.SetMaxOpenConns(open)
	}
}

func WithMaxLifeTime(t int) Option {
	return func(m *sql.DB) {
		if t == 0 {
			return
		}
		m.SetConnMaxLifetime(time.Duration(t) * time.Second)
	}
}
func Options(m *sql.DB, opts ...Option) {
	for _, opt := range opts {
		opt(m)
	}
}
```

主函数

```go
// main.go
package main

import (
	"catering/initialize"
    "catering/global"
)

func main() {
  // 初始化db
  global.DB := initialize.InitGorm()
	if global.DB != nil {
		db, _ := global.DB.DB()
		defer func() {
      if err := db.Close(); err != nil {
        panic(err)
      }
    }()
	}
}
```



## 模型定义

在GORM中，模型通常是一个Go结构体，用于映射数据库表。每个字段对应数据库表中的一列。以下是一个简单的模型定义示例：

```go
type student struct {
    ID        uint           `gorm:"primaryKey"`   // 默认使用ID作为主键
    Name      string         `gorm:"size:255"`
    Email     *string         `gorm:"uniqueIndex"` // 使用指针是为了存空值
    CreatedAt time.Time
    UpdatedAt time.Time
}

// 常识：小写属性是不会生成字段的
```

在上面的代码中，student结构体定义了四个字段：`ID、Name、Email、CreatedAt`和`UpdatedAt`。其中，ID字段被标记为主键，Name字段的最大长度为255个字符，Email字段被标记为唯一索引。

### 自定义表名和列名

默认情况下，GORM会将结构体名称转换为下划线命名的表名，并且将结构体字段名称转换为下划线命名的列名。如果需要自定义表名或列名，可以使用TableName方法和标签：

```go
func (student) TableName() string {
    return "students"
}

type student struct {
    ID        uint   `gorm:"column:student_id"`
    Name      string `gorm:"column:student_name"`
    Email     string `gorm:"column:student_email"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

在上面的代码中，TableName方法返回自定义的表名`students`，而标签`gorm:"column:student_id"`等则指定了自定义的列名。

### 迁移

GORM提供了自动迁移功能，可以使用AutoMigrate方法自动创建、更新和删除数据库表：

```go
// 创建表时添加后缀
db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

```

在上面的代码中，我们使用`AutoMigrate`方法自动迁移`student、Order`和`Product`表。

> `AutoMigrate`的逻辑是只新增，不删除，不修改（大小会修改）

例如将Name修改为Name1，进行迁移，会多出一个name1的字段

生成的表结构如下
```SQL
CREATE TABLE `f_students` (`id` bigint unsigned AUTO_INCREMENT,`name` longtext,`email` longtext,PRIMARY KEY (`id`))
```
默认的类型太大了

注意 AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能，例如：
```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  DisableForeignKeyConstraintWhenMigrating: true,
})
```

### 修改大小

我们可以使用gorm的标签进行修改

有两种方式

```go
Name  string  `gorm:"type:varchar(12)"` // 定义类型并限制大小
Name  string  `gorm:"size:2"`						// 定义大小
```



### 字段标签

`type` 定义字段类型

`size` 定义字段大小

`column` 自定义列名

`primaryKey` 将列定义为主键

`unique` 将列定义为唯一键

`default` 定义列的默认值

`not null` 不可为空

`embedded` 嵌套字段

`embeddedPrefix` 嵌套字段前缀

`comment` 注释

多个标签之前用 `;` 连接

```go
type StudentInfo struct {
  Email  *string `gorm:"size:32"` // 使用指针是为了存空值
  Addr   string  `gorm:"column:y_addr;size:16"`
  Gender bool    `gorm:"default:true"`
}
type Student struct {
  Name string      `gorm:"type:varchar(12);not null;comment:用户名"`
  UUID string      `gorm:"primaryKey;unique;comment:主键"`
  Info StudentInfo `gorm:"embedded;embeddedPrefix:s_"`
}

// 建表语句
CREATE TABLE `students` (
    `name` varchar(12) NOT NULL COMMENT '用户名',
    `uuid` varchar(191) UNIQUE COMMENT '主键',
    `s_email` varchar(32),
    `s_y_addr` varchar(16),
    `s_gender` boolean DEFAULT true,
    PRIMARY KEY (`uuid`)
)

```

### 模型转换工具

sql转gorm： https://old.printlove.cn/tools/sql2gorm



## 单表CRUD操作

CRUD操作是数据库操作的基本组成部分，分别对应创建（Create）、读取（Read）、更新（Update）和删除（Delete）。GORM提供了简便的方法来执行这些操作。

表结构

```go
type Student struct {
  ID     uint   `gorm:"size:3"`
  Name   string `gorm:"size:8"`
  Age    int    `gorm:"size:3;default:18"`
  Sex    string `gorm:"type:enum('male', 'female', 'unknown');default:'unknown'"` // 性别
  Email  *string `gorm:"size:32"`
}

```



### 创建记录

#### 插入一条记录

使用Create方法可以向数据库中插入一条记录：

```go
email := "xxx@qq.com"
// 创建记录
student := Student{
  Name:   "John",
  Age:    21,
  Gender: true,
  Email:  &email,
}

result := DB.Create(&student)

if result.Error != nil {
    log.Fatal(result.Error)
}
```

在上面的代码中，我们创建了一个student对象，并使用`db.Create`方法将其插入数据库。如果操作成功，`result.Error`将为`nil`。

有两个地方需要注意

1. 指针类型是为了更好的存null类型，但是传值的时候，也记得传指针
2. Create接收的是一个指针，而不是值

由于我们传递的是一个指针，调用完Create之后，student这个对象上面就有该记录的信息了，如创建的id

```go
DB.Create(&student)
fmt.Printf("%#v\n", student)  
// main.Student{ID:0x2, Name:"zhangsan", Age:23, Gender:false, Email:(*string)(0x11d40980)}
```

#### 批量创建

方式一：使用切片

```go
var studentList []Student
for i := 0; i < 100; i++ {
  studentList = append(studentList, Student{
    Name:   fmt.Sprintf("机器人%d号", i+1),
    Age:    21,
    Gender: true,
    Email:  &email,
  })
}
DB.Create(&studentList)

```

方式二：可以使用CreateInBatches方法进行批量插入：

```go
email := "xxx@qq.com"
var studentList []Student
for i := 0; i < 1000; i++ {
    students = append(studentList, student{
      Name: fmt.Sprintf("student%d", i),
      Email: &email,
    })
}
db.CreateInBatches(studentList, 100)
```

在上面的代码中，我们使用CreateInBatches方法每次插入100条记录。



### 查询记录

可以使用`First、Take、Last`方法读取/查询单条记录，倘若未查询到指定记录，则会报错 `gorm.ErrRecordNotFound`;

使用`Find`方法读取多条记录,即便没有查到指定记录，也不会返回错误; 但是返回的记录为空切片。

使用`Select`方法可以查询指定字段相关记录

#### 查询一条记录

```go
var student Student
db = global.DB.Session(&gorm.Session{Logger: Log})

// 1、读取第一条记录: SELECT * FROM `student` ORDER BY `student`.`id` LIMIT 1
result := db.First(&student)
// 是否查询失败
if result.Error != nil {
    log.Fatal(result.Error)
    if result.Error == gorm.ErrRecordNotFound {
        log.Println("没有找到记录")
    }
}
// 获取查询的行数
if result.RowsAffected == 1 {
    log.Println("找到1条记录")
}

// 2、根据主键读取记录：SELECT * FROM `student` WHERE `id` = 1 ORDER BY `student`.`id` LIMIT 1
result := db.First(&student, 1)
if result.Error != nil {
    log.Fatal(result.Error)
    if result.Error == gorm.ErrRecordNotFound {
        log.Println("没有找到记录")
    }
}
if result.RowsAffected == 1 {
    log.Println("找到1条记录")
    
}

// 3、随机一条记录返回: SELECT * FROM `student` LIMIT 1
result := db.Take(&student)
if result.Error != nil {
    log.Fatal(result.Error)
    if result.Error == gorm.ErrRecordNotFound {
        log.Println("没有找到记录")
    }
}
if result.RowsAffected == 1 {
    log.Println("找到1条记录")
}

// 4、读取最后一条记录: SELECT * FROM `student` ORDER BY `student`.`id` DESC LIMIT 1
result := db.Last(&student)
if result.Error != nil {
    log.Fatal(result.Error)
    if result.Error == gorm.ErrRecordNotFound {
        log.Println("没有找到记录")
    }
}
if result.RowsAffected == 1 {
    log.Println("找到1条记录")
}

// 5、查询指定字段记录：
SELECT * FROM `student` ORDER BY `student`.`id` LIMIT 1
result := db.Select("name").First(&student)
// 是否查询失败
if result.Error != nil {
    log.Fatal(result.Error)
    if result.Error == gorm.ErrRecordNotFound {
        log.Println("没有找到记录")
    }
}
// 获取查询的行数
if result.RowsAffected == 1 {
    log.Println("找到1条记录")
}
```

在上面的代码中，我们使用First方法读取第一条记录，使用Find方法读取所有记录。

#### 查询多条记录

```go
// 读取所有记录：find结果是个切片
var students []student
result = db.Find(&students)
if result.Error != nil {
    log.Fatal(result.Error)
}
if len(students) > 0 {
    log.Printf("student: %+v", student)
    // 因为email是指针类型，所以需要序列化后才可以看到
  	for _, student := range studentList {
  			data, _ := json.Marshal(student)
  			fmt.Println(string(data))
		}
}

```



GORM提供了多种方法来构造查询，包括链式调用、原生SQL查询等。

#### 条件查询

可以使用Where方法添加条件：

```go
var students student

// 使用"?"占位符,将参数全部转为字符串，可以防止sql注入
result := db.Where("name = ?", "John").First(&students)
if result.Error != nil {
    log.Fatal(result.Error)
    if result.Error == gorm.ErrRecordNotFound {
        log.Println("没有找到记录")
    }
}

// 多个条件
result = db.Where("name = ? AND email = ?", "John", "john@example.com").First(&student)
if result.Error != nil {
    log.Fatal(result.Error)
    if result.Error == gorm.ErrRecordNotFound {
        log.Println("没有找到记录")
    }
}
```

在上面的代码中，我们使用Where方法构造了带条件的查询。

#### 链式调用

可以使用链式调用构造更复杂的查询：

```go
var students []student
result := db.Where("name = ?", "John").Or("name = ?", "Jane").Find(&students)
if result.Error != nil {
    log.Fatal(result.Error)
}

// 链式调用添加更多条件
result = db.Where("name = ?", "John").Or("name = ?", "Jane").Order("created_at desc").Limit(10).Find(&students)
if result.Error != nil {
    log.Fatal(result.Error)
}
```

在上面的代码中，我们使用链式调用构造了带多个条件、排序和限制的查询。

#### 原生SQL查询

可以使用Raw方法执行原生SQL查询：

```go
var students []student
result := db.Raw("SELECT * FROM students WHERE name = ?", "John").Scan(&students)
if result.Error != nil {
    log.Fatal(result.Error)
}
```

在上面的代码中，我们使用Raw方法执行了一条原生SQL查询。



### 更新记录

可以使用Save方法或Updates方法更新记录，不过前提是先查到记录

- `Save`

```go
var student student
db.First(&student, 1)

// 使用save，会更新所有字段，包括零值
student.Name = "Jane"
result := db.Save(&student)
// UPDATE `students` SET `name`='Jane' WHERE `id` = 1
if result.Error != nil {
    log.Fatal(result.Error)
}

student.Age = 0
result := db.Save(&student)
// UPDATE `students` SET `age`=0 WHERE `id` = 1
if result.Error != nil {
    log.Fatal(result.Error)
}


// 使用select指定字段
var student Student
DB.Take(&student)

student.Age = 21
// 指定age字段
DB.Select("age").Save(&student)
// UPDATE `students` SET `age`=21 WHERE `id` = 1

```

- `Updates`

```go
// 使用 Updates 方法

// UPDATE `students` SET `email`='is22@qq.com' WHERE age = 21
result = db.Model(&student).Where("age = ?", 21).Updates(student{Name: "Jane", Email: "jane@example.com"})
if result.Error != nil {
    log.Fatal(result.Error)
}
```

在上面的代码中，我们先读取了一条记录，然后更新了这条记录的字段值，并使用Save或Updates方法将其保存到数据库。

#### 批量更新

例如年龄21的学生，都更新一下邮箱

```go
var studentList []Student
DB.Find(&studentList, "age = ?", 21).Update("email", "is21@qq.com")
```

还有一种更简单的方式

```go
DB.Model(&Student{}).Where("age = ?", 21).Update("email", "is21@qq.com")
// UPDATE `students` SET `email`='is22@qq.com' WHERE age = 21
```

这样的更新方式也是可以更新零值的

#### 更新多列

如果是结构体，它默认不会更新零值

```go
email := "xxx@qq.com"
DB.Model(&Student{}).Where("age = ?", 21).Updates(Student{
  Email:  &email,
  Gender: false,  // 这个不会更新
})

// UPDATE `students` SET `email`='xxx@qq.com' WHERE age = 21
```

如果想让他更新零值，用select就好

```go
email := "xxx1@qq.com"
DB.Model(&Student{}).Where("age = ?", 21).Select("gender", "email").Updates(Student{
  Email:  &email,
  Gender: false,
})
// UPDATE `students` SET `gender`=false,`email`='xxx1@qq.com' WHERE age = 21
```

如果不想多写几行代码，则推荐使用map

```go
DB.Model(&Student{}).Where("age = ?", 21).Updates(map[string]any{
  "email":  &email,
  "gender": false,
})
```

#### 更新选定字段

Select选定字段

Omit忽略字段



### 删除记录

可以使用Delete方法删除记录：

```go
var student student
db.First(&student, 1)

result := db.Delete(&student)
if result.Error != nil {
    log.Fatal(result.Error)
}

// 或者根据主键删除
result = db.Delete(&student{}, 1)
if result.Error != nil {
    log.Fatal(result.Error)
}

```

在上面的代码中，我们使用Delete方法删除了一条记录，可以通过传递记录对象或主键来指定要删除的记录。

也可以批量删除

```go
db.Delete(&Student{}, []int{1,2,3})

// 查询到的切片列表
var studentList []*student
db.Delete(&studentList)

```

## 高级查询

GORM支持多种高级查询功能，如分页、子查询和联合查询等。

### 使用结构体查询

使用结构体查询，会过滤零值

并且结构体中的条件都是and关系

```go
// 会过滤零值
DB.Where(&Student{Name: "李元芳", Age: 0}).Find(&users)
fmt.Println(users)
```

### 使用map查询

不会过滤零值

```go
DB.Where(map[string]any{"name": "李元芳", "age": 0}).Find(&users)
// SELECT * FROM `students` WHERE `age` = 0 AND `name` = '李元芳'
fmt.Println(users)
```

### Not条件

和where中的not等价

```go
// 排除年龄大于23的
DB.Not("age > 23").Find(&users)
fmt.Println(users)
```

### Or条件

和where中的or等价

```go
DB.Or("gender = ?", false).Or(" email like ?", "%@qq.com").Find(&users)
fmt.Println(users)
```

### Select 选择字段

```go
DB.Select("name", "age").Find(&users)
fmt.Println(users)
// 没有被选中，会被赋零值

```

可以使用扫描Scan，将选择的字段存入另一个结构体中

```go

type User struct {
  Name string
  Age  int
}
var students []Student
var users []User
DB.Select("name", "age").Find(&students).Scan(&users)
fmt.Println(users)

```

这样写也是可以的，不过最终会查询两次，还是不这样写

```go
SELECT `name`,`age` FROM `students`
SELECT `name`,`age` FROM `students`
```

这样写就只查询一次了

```go

type User struct {
  Name string
  Age  int
}
var users []User
DB.Model(&Student{}).Select("name", "age").Scan(&users)
fmt.Println(users)
```

还可以这样

```go
var users []User
DB.Table("students").Select("name", "age").Scan(&users)
fmt.Println(users)
```

Scan是根据column列名进行扫描的

```go

type User struct {
  Name123 string `gorm:"column:name"`
  Age     int
}
var users []User
DB.Table("students").Select("name", "age").Scan(&users)
fmt.Println(users)
```

### 排序

根据年龄倒序

```go

var users []Student
DB.Order("age desc").Find(&users)
fmt.Println(users)
// desc    降序
// asc     升序
```

注意order的顺序

### 分页查询

```go

var users []Student
// 一页两条，第1页
DB.Limit(2).Offset(0).Find(&users)
fmt.Println(users)

// 第2页
DB.Limit(2).Offset(2).Find(&users)
fmt.Println(users)

// 第3页
DB.Limit(2).Offset(4).Find(&users)
fmt.Println(users)

```

通用写法

```go
var users []Student
// 一页多少条
limit := 2
// 第几页
page := 1
offset := (page - 1) * limit
result := DB.Limit(limit).Offset(offset).Find(&users)

if result.Error != nil {
    log.Fatal(result.Error)
}

fmt.Println(users)
```

### 去重

```go
var ageList []int
DB.Table("students").Select("age").Distinct("age").Scan(&ageList)
fmt.Println(ageList)
```

或者

```go
DB.Table("students").Select("distinct age").Scan(&ageList)
```

### 分组查询

```go

var ageList []int
// 查询男生的个数和女生的个数
DB.Table("students").Select("count(id)").Group("gender").Scan(&ageList)
fmt.Println(ageList)
```

有个问题，哪一个是男生个数，那个是女生个数

所以我们应该精确一点

```go

type AggeGroup struct {
  Gender int
  Count  int `gorm:"column:count(id)"`
}

var agge []AggeGroup
// 查询男生的个数和女生的个数
DB.Table("students").Select("count(id)", "gender").Group("gender").Scan(&agge)
fmt.Println(agge)
```

如何再精确一点，具体的男生名字，女生名字

```go

type AggeGroup struct {
  Gender int
  Count  int    `gorm:"column:count(id)"`
  Name   string `gorm:"column:group_concat(name)"`
}

var agge []AggeGroup
// 查询男生的个数和女生的个数
DB.Table("students").Select("count(id)", "gender", "group_concat(name)").Group("gender").Scan(&agge)
fmt.Println(agge)
```

总之，使用gorm不会让你忘记原生sql的编写

这一点我还是很喜欢的

### 执行原生sql

```go
type AggeGroup struct {
  Gender int
  Count  int    `gorm:"column:count(id)"`
  Name   string `gorm:"column:group_concat(name)"`
}

var agge []AggeGroup
DB.Raw(`SELECT count(id), gender, group_concat(name) FROM students GROUP BY gender`).Scan(&agge)

fmt.Println(agge)
```

### 子查询

可以使用SubQuery方法构造子查询：

```go
var users []Student
result := DB.Model(Student{}).Where("age > (?)", DB.Model(Student{}).Select("avg(age)")).Find(&users)

if result.Error != nil {
    log.Fatal(result.Error)
}
fmt.Println(users)
```

在上面的代码中，我们使用SubQuery方法构造了一个子查询，并在主查询中使用。

### 联合查询

可以使用Joins方法进行联合查询：

```go
var students []student
result := db.Joins("JOIN orders ON orders.student_id = students.id").Where("orders.amount > ?", 1000).Find(&students)
if result.Error != nil {
    log.Fatal(result.Error)
}
```

在上面的代码中，我们使用Joins方法进行联合查询。

### 命名参数

```go

var users []Student

DB.Where("name = @name and age = @age", sql.Named("name", "Jone"), sql.Named("age", 23)).Find(&users)
DB.Where("name = @name and age = @age", map[string]any{"name": "Jone", "age": 23}).Find(&users)
fmt.Println(users)
```

### find到map

```go
var res []map[string]any
DB.Table("students").Find(&res)
fmt.Println(res)
```

### 查询引用Scope

可以再model层写一些通用的查询方式，这样外界就可以直接调用方法即可

```go
func Age23(db *gorm.DB) *gorm.DB {
  return db.Where("age > ?", 23)
}

func main(){
  var users []Student
  DB.Scopes(Age23).Find(&users)
  fmt.Println(users)
}
```

## 关联关系

GORM支持多种关联关系，包括一对一、一对多和多对多关系。

### 一对一关系

一对一关系比较少，一般用于表的扩展

例如一张用户表，有很多字段

那么就可以把它拆分为两张表，常用的字段放主表，不常用的字段放详情表

#### 表结构搭建

```go

type User struct {
  ID       uint
  Name     string
  Age      int
  Gender   bool
  UserInfo UserInfo // 通过UserInfo可以拿到用户详情信息
}

type UserInfo struct {
  UserID uint // 外键
  ID     uint
  Addr   string
  Like   string
}
```

#### 添加记录

添加用户，自动添加用户详情

```go

DB.Create(&User{
  Name:   "Jone",
  Age:    21,
  Gender: true,
  UserInfo: UserInfo{
    Addr: "湖南省",
    Like: "写代码",
  },
})
```

添加用户详情，关联已有用户

这个场景特别适合网站的注册，以及后续信息完善

刚开始注册的时候，只需要填写很基本的信息，这就是添加主表的一条记录

注册进去之后，去个人中心，添加头像，修改地址...

这就是添加附表

```go

DB.Create(&UserInfo{
  UserID: 2,
  Addr:   "南京市",
  Like:   "吃饭",
})
```

当然，也可以直接把用户对象传递进来

我们需要改一下表结构

```go

type User struct {
  ID       uint
  Name     string
  Age      int
  Gender   bool
  UserInfo UserInfo // 通过UserInfo可以拿到用户详情信息
}

type UserInfo struct {
  User *User  // 要改成指针，不然就嵌套引用了
  UserID uint // 外键
  ID     uint
  Addr   string
  Like   string
}
```

不限于重新迁移，直接添加即可

```go

var user User
DB.Take(&user, 2)
DB.Create(&UserInfo{
  User: &user,
  Addr: "南京市",
  Like: "吃饭",
})
```

#### 查询

一般是通过主表查副表

```go
var user User
DB.Preload("UserInfo").Take(&user)
fmt.Println(user)
```

------

### 一对多关系

#### 表结构建立

在gorm中，官方文档是把一对多关系分为了两类，

- Belongs To 属于谁

- Has Many 我拥有的

他们本来是一起的，本教程把它们合在一起讲，我们以用户和文章为例

一个用户可以发布多篇文章，一篇文章属于一个用户

```go
type User struct {
  ID       uint      `gorm:"size:4"`
  Name     string    `gorm:"size:8"`
  Articles []Article // 用户拥有的文章列表
}

type Article struct {
  ID     uint   `gorm:"size:4"`
  Title  string `gorm:"size:16"`
  UserID uint   // 属于   这里的类型要和引用的外键类型一致，包括大小
  User   User   // 属于
}
```

关于外键命名，外键名称就是关联表名+ID，类型是uint

##### 重写外键关联

```go
type User struct {
  ID       uint      `gorm:"size:4"`
  Name     string    `gorm:"size:8"`
  Articles []Article `gorm:"foreignKey:UID"` // 用户拥有的文章列表
}

type Article struct {
  ID    uint   `gorm:"size:4"`
  Title string `gorm:"size:16"`
  UID   uint   // 属于
  User  User   `gorm:"foreignKey:UID"` // 属于
}
```

这里有个地方要注意：

- 我改了Article 的外键，将UID作为了外键，那么User这个外键关系就要指向UID

- 与此同时，User所拥有的Articles也得更改外键，改为UID

##### 重写外键引用

```go
type User struct {
  ID       uint      `gorm:"size:4"`
  Name     string    `gorm:"size:8"`
  Articles []Article `gorm:"foreignKey:UserName;references:Name"` // 用户拥有的文章列表
}

type Article struct {
  ID       uint   `gorm:"size:4"`
  Title    string `gorm:"size:16"`
  UserName string
  User     User `gorm:"references:Name"` // 属于
}
```

这一块的逻辑比较复杂

比如有1个用户

|      |      |
| ---- | ---- |
| id   | name |
| 1    | Jone |

之前的外键关系是这样表示文章的

|      |            |         |
| ---- | ---------- | ------- |
| id   | title      | user_id |
| 1    | python     | 1       |
| 2    | javascript | 1       |
| 3    | golang     | 1       |

如果改成直接关联Name，那就变成了这样

|      |            |           |
| ---- | ---------- | --------- |
| id   | title      | user_name |
| 1    | python     | Jone      |
| 2    | javascript | Jone      |
| 3    | golang     | Jone      |

虽然这样很方便，但是非常不适合在实际项目中这样用

我们还是用第一版的表结构做一对多关系的增删改查



#### 一对多的添加

创建用户，并且创建文章

```go
a1 := Article{Title: "python"}
a2 := Article{Title: "golang"}
user := User{Name: "Jone", Articles: []Article{a1, a2}}
DB.Create(&user)
```

gorm自动创建了两篇文章，以及创建了一个用户，还将他们的关系给关联上了

创建文章，关联已有用户

```go
a1 := Article{Title: "golang零基础入门", UserID: 1}
DB.Create(&a1)

var user User
DB.Take(&user, 1)
DB.Create(&Article{Title: "python零基础入门", User: user})
```

#### 外键添加

给现有用户绑定文章

```go
var user User
DB.Take(&user, 2)

var article Article
DB.Take(&article, 5)

user.Articles = []Article{article}
DB.Save(&user)
```

也可以用Append方法

```go
var user User
DB.Take(&user, 2)

var article Article
DB.Take(&article, 5)

//user.Articles = []Article{article}
//DB.Save(&user)

DB.Model(&user).Association("Articles").Append(&article)
```

给现有文章关联用户

```go
var article Article
DB.Take(&article, 5)

article.UserID = 2
DB.Save(&article)
```

也可用Append方法

```go
var user User
DB.Take(&user, 2)

var article Article
DB.Take(&article, 5)

DB.Model(&article).Association("User").Append(&user)
```

#### 查询

查询用户，显示用户的文章列表

```go
var user User
DB.Take(&user, 1)
fmt.Println(user)
```

直接这样，是显示不出文章列表

##### 预加载

我们必须要使用预加载来加载文章列表

```go
var user User
DB.Preload("Articles").Take(&user, 1) // 预加载的名字就是外键关联的属性名
fmt.Println(user)
```



查询文章，显示文章用户的信息

同样的，使用预加载

```go
var article Article
DB.Preload("User").Take(&article, 1)
fmt.Println(article)
```

##### 嵌套预加载

查询文章，显示用户，并且显示用户关联的所有文章，这就得用到嵌套预加载了

```go
var article Article
DB.Preload("User.Articles").Take(&article, 1)
fmt.Println(article)
```

##### 带条件的预加载

查询用户下的所有文章列表，过滤某些文章

```go
var user User
DB.Preload("Articles", "id = ?", 1).Take(&user, 1)
fmt.Println(user)
```

这样，就只有id为1的文章被预加载出来了

##### 自定义预加载

```go
var user User
DB.Preload("Articles", func(db *gorm.DB) *gorm.DB {
  return db.Where("id in ?", []int{1, 2})
}).Take(&user, 1)
fmt.Println(user)
```

#### 删除

##### 级联删除

删除用户，与用户关联的文章也会删除

```go
var user User
DB.Take(&user, 1)
DB.Select("Articles").Delete(&user)
```

##### 清除外键关系

删除用户，与将与用户关联的文章，外键设置为null

```go
var user User
DB.Preload("Articles").Take(&user, 2)
DB.Model(&user).Association("Articles").Delete(&user.Articles)

```

### 多对多关系

需要用第三张表存储两张表的关系

#### 表结构搭建

```go
type Tag struct {
  ID       uint
  Name     string
  Articles []Article `gorm:"many2many:article_tags;"` // 用于反向引用
}

type Article struct {
  ID    uint
  Title string
  Tags  []Tag `gorm:"many2many:article_tags;"`
}
```

#### 多对多添加

添加文章，并创建标签

```go
DB.Create(&Article{
  Title: "python基础课程",
  Tags: []Tag{
    {Name: "python"},
    {Name: "基础课程"},
  },
})
```

添加文章，选择标签

```go
var tags []Tag
DB.Find(&tags, "name = ?", "基础课程")
DB.Create(&Article{
  Title: "golang基础",
  Tags:  tags,
})
```

#### 多对多查询

查询文章，显示文章的标签列表

```go
var article Article
DB.Preload("Tags").Take(&article, 1)
fmt.Println(article)
```

查询标签，显示文章列表

```go
var tag Tag
DB.Preload("Articles").Take(&tag, 2)
fmt.Println(tag)
```

#### 多对多更新

移除文章的标签

```go
var article Article
DB.Preload("Tags").Take(&article, 1)
DB.Model(&article).Association("Tags").Delete(article.Tags)
fmt.Println(article)
```

更新文章的标签

```go
var article Article
var tags []Tag
DB.Find(&tags, []int{2, 6, 7})

DB.Preload("Tags").Take(&article, 2)
DB.Model(&article).Association("Tags").Replace(tags)
fmt.Println(article)
```

#### 自定义连接表

默认的连接表，只有双方的主键id，展示不了更多信息了

这是官方的例子，我修改了一下

```go
type Article struct {
  ID    uint
  Title string
  Tags  []Tag `gorm:"many2many:article_tags"`
}

type Tag struct {
  ID   uint
  Name string
}

type ArticleTag struct {
  ArticleID uint `gorm:"primaryKey"`
  TagID     uint `gorm:"primaryKey"`
  CreatedAt time.Time
}

```

##### 生成表结构

```go
// 设置Article的Tags表为ArticleTag
DB.SetupJoinTable(&Article{}, "Tags", &ArticleTag{})
// 如果tag要反向应用Article，那么也得加上
// DB.SetupJoinTable(&Tag{}, "Articles", &ArticleTag{})
err := DB.AutoMigrate(&Article{}, &Tag{}, &ArticleTag{})
fmt.Println(err)
```

##### 操作案例

举一些简单的例子：

1）添加文章并添加标签，并自动关联

2）添加文章，关联已有标签

3）给已有文章关联标签

4）替换已有文章的标签

> SetupJoinTable: 添加和更新的时候得用这个,这样才能走自定义的连接表，以及走它的钩子函数; 查询则不需要这个

1. 添加文章并添加标签，并自动关联

```go
DB.SetupJoinTable(&Article{}, "Tags", &ArticleTag{})  // 要设置这个，才能走到我们自定义的连接表
DB.Create(&Article{
  Title: "flask零基础入门",
  Tags: []Tag{
    {Name: "python"},
    {Name: "后端"}, 
    {Name: "web"},
  },
})
// CreatedAt time.Time 由于我们设置的是CreatedAt，gorm会自动填充当前时间，
// 如果是其他的字段，需要使用到ArticleTag 的添加钩子 BeforeCreate
```

2. 添加文章，关联已有标签

```go
DB.SetupJoinTable(&Article{}, "Tags", &ArticleTag{})
var tags []Tag
DB.Find(&tags, "name in ?", []string{"python", "web"})
DB.Create(&Article{
  Title: "flask请求对象",
  Tags:  tags,
})
```

3. 给已有文章关联标签

```go
DB.SetupJoinTable(&Article{}, "Tags", &ArticleTag{})
article := Article{
  Title: "django基础",
}
DB.Create(&article)
var at Article
var tags []Tag
DB.Find(&tags, "name in ?", []string{"python", "web"})
DB.Take(&at, article.ID).Association("Tags").Append(tags)
```

4. 替换已有文章的标签

```go
var article Article
var tags []Tag
DB.Find(&tags, "name in ?", []string{"后端"})
DB.Take(&article, "title = ?", "django基础")
DB.Model(&article).Association("Tags").Replace(tags)
```

6. 查询文章列表，显示标签

```go
var articles []Article
DB.Preload("Tags").Find(&articles)
fmt.Println(articles)
```



#### 自定义连接表主键

这个功能还是很有用的，例如你的文章表 可能叫ArticleModel，你的标签表可能叫TagModel

那么按照gorm默认的主键名，那就分别是ArticleModelID，TagModelID，太长了，根本就不实用

这个地方，官网给的例子看着也比较迷，不过我已经跑通了,主要是要修改这两项:

- joinForeignKey 连接的主键id

- JoinReferences 关联的主键id

```go
type ArticleModel struct {
  ID    uint
  Title string
  Tags  []TagModel `gorm:"many2many:article_tags;joinForeignKey:ArticleID;JoinReferences:TagID"`
}

type TagModel struct {
  ID       uint
  Name     string
  Articles []ArticleModel `gorm:"many2many:article_tags;joinForeignKey:TagID;JoinReferences:ArticleID"`
}

type ArticleTagModel struct {
  ArticleID uint `gorm:"primaryKey"` // article_id
  TagID     uint `gorm:"primaryKey"` // tag_id
  CreatedAt time.Time
}
```

##### 生成表结构

```go
DB.SetupJoinTable(&ArticleModel{}, "Tags", &ArticleTagModel{})
DB.SetupJoinTable(&TagModel{}, "Articles", &ArticleTagModel{})
err := DB.AutoMigrate(&ArticleModel{}, &TagModel{}, &ArticleTagModel{})
fmt.Println(err)
```

添加，更新，查询操作和上面的都是一样

##### 操作连接表

如果通过一张表去操作连接表，这样会比较麻烦, 比如查询某篇文章关联了哪些标签, 或者是举个更通用的例子，用户和文章，某个用户在什么时候收藏了哪篇文章

无论是通过用户关联文章，还是文章关联用户都不太好查，最简单的就是直接查连接表

```go
type UserModel struct {
  ID       uint
  Name     string
  Collects []ArticleModel `gorm:"many2many:user_collect_models;joinForeignKey:UserID;JoinReferences:ArticleID"`
}

type ArticleModel struct {
  ID    uint
  Title string
  // 这里也可以反向引用，根据文章查哪些用户收藏了
}

// UserCollectModel 用户收藏文章表
type UserCollectModel struct {
  UserID    uint `gorm:"primaryKey"` // article_id
  ArticleID uint `gorm:"primaryKey"` // tag_id
  CreatedAt time.Time
}

func main() {
  DB.SetupJoinTable(&UserModel{}, "Collects", &UserCollectModel{})
  err := DB.AutoMigrate(&UserModel{}, &ArticleModel{}, &UserCollectModel{})
  fmt.Println(err)
}
```

常用的操作就是根据用户查收藏的文章列表

```go
var user UserModel
DB.Preload("Collects").Take(&user, "name = ?", "枫枫")
fmt.Println(user)
```

但是这样不太好做分页，并且也拿不到收藏文章的时间

```go
var collects []UserCollectModel
DB.Find(&collects, "user_id = ?", 2)
fmt.Println(collects)
```

这样虽然可以查到用户id，文章id，收藏的时间，但是搜索只能根据用户id搜，返回也拿不到用户名，文章标题等

我们需要改一下表结构，不需要重新迁移，加一些字段

```go

type UserModel struct {
  ID       uint
  Name     string
  Collects []ArticleModel `gorm:"many2many:user_collect_models;joinForeignKey:UserID;JoinReferences:ArticleID"`
}

type ArticleModel struct {
  ID    uint
  Title string
}

// UserCollectModel 用户收藏文章表
type UserCollectModel struct {
  UserID       uint         `gorm:"primaryKey"` // article_id
  UserModel    UserModel    `gorm:"foreignKey:UserID"`
  ArticleID    uint         `gorm:"primaryKey"` // tag_id
  ArticleModel ArticleModel `gorm:"foreignKey:ArticleID"`
  CreatedAt    time.Time
}
```

查询

```go

var collects []UserCollectModel

var user UserModel
DB.Take(&user, "name = ?", "枫枫")
// 这里用map的原因是如果没查到，那就会查0值，如果是struct，则会忽略零值，全部查询
DB.Debug().Preload("UserModel").Preload("ArticleModel").Where(map[string]any{"user_id": user.ID}).Find(&collects)

for _, collect := range collects {
  fmt.Println(collect)
}
```

------

## 自定义数据类型

自定义的数据类型必须实现 Scanner 和 Valuer 接口，以便让 GORM 知道如何将该类型接收、保存到数据库

### 存储结构体

```go

type Info struct {
  Status string `json:"status"`
  Addr   string `json:"addr"`
  Age    int    `json:"age"`
}

// Scan 从数据库中读取出来
func (i *Info) Scan(value interface{}) error {
  bytes, ok := value.([]byte)
  if !ok {
    return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
  }

  info := Info{}
  err := json.Unmarshal(bytes, &info)
  *i = info
  return err
}

// Value 存入数据库
func (i Info) Value() (driver.Value, error) {
  return json.Marshal(i)
}

type User struct {
  ID   uint
  Name string
  Info Info `gorm:"type:string"`
}

```

#### 添加和查询

```go

DB.Create(&User{
  Name: "枫枫",
  Info: Info{
    Status: "牛逼",
    Addr:   "成都市",
    Age:    21,
  },
})

var user User
DB.Take(&user)
fmt.Println(user)
```

### 枚举类型

#### 枚举1.0

很多时候，我们会对一些状态进行判断，而这些状态都是有限的

例如，主机管理中，状态有 Running 运行中， OffLine 离线， Except 异常

如果存储字符串，不仅是浪费空间，每次判断还要多复制很多字符，最主要是后期维护麻烦

```go

type Host struct {
  ID     uint
  Name   string
  Status string
}

func main() {
  host := Host{}
  if host.Status == "Running" {
    fmt.Println("在线")
  }
  if host.Status == "Except" {
    fmt.Println("异常")
  }
  if host.Status == "OffLine" {
    fmt.Println("离线")
  }
}
```

后来，我们知道了用常量存储这些不变的值

```go

type Host struct {
  ID     uint
  Name   string
  Status string
}

const (
  Running = "Running"
  Except = "Except"
  OffLine = "OffLine"
) 

func main() {
  host := Host{}
  if host.Status == Running {
    fmt.Println("在线")
  }
  if host.Status == Except {
    fmt.Println("异常")
  }
  if host.Status == OffLine {
    fmt.Println("离线")
  }
}

```

虽然代码变多了，但是维护方便了

但是数据库中存储的依然是字符串，浪费空间这个问题并没有解决

#### 枚举2.0

于是想到使用数字表示状态

```go

type Host struct {
  ID     uint
  Name   string
  Status int
}

const (
  Running = 1
  Except  = 2
  OffLine = 3
)

func main() {
  host := Host{}
  if host.Status == Running {
    fmt.Println("在线")
  }
  if host.Status == Except {
    fmt.Println("异常")
  }
  if host.Status == OffLine {
    fmt.Println("离线")
  }
}
```

但是，如果返回数据给前端，前端接收到的状态就是数字，不过问题不大，前端反正都要搞字符映射的

因为要做颜色差异显示

但是这并不是后端偷懒的理由

于是我们想到，在json序列化的时候，根据映射转换回去

```go

type Host struct {
  ID     uint   `json:"id"`
  Name   string `json:"name"`
  Status int    `json:"status"`
}

func (h Host) MarshalJSON() ([]byte, error) {
  var status string
  switch h.Status {
  case Running:
    status = "Running"
  case Except:
    status = "Except"
  case OffLine :
    status = "OffLine"
  }
  return json.Marshal(&struct {
    ID     uint   `json:"id"`
    Name   string `json:"name"`
    Status string `json:"status"`
  }{
    ID:     h.ID,
    Name:   h.Name,
    Status: status,
  })
}

const (
  Running = 1
  Except  = 2
  OffLine  = 3
)

func main() {
  host := Host{1, "枫枫", Running}
  data, _ := json.Marshal(host)
  fmt.Println(string(data)) // {"id":1,"name":"枫枫","status":"Running"}
}

```

这样写确实可以实现我们的需求，但是根本就不够通用，凡是用到枚举，都得给这个Struct实现`MarshalJSON`方法

#### 枚举3.0

于是类型别名出来了

```go
type Status int

func (status Status) MarshalJSON() ([]byte, error) {
  var str string
  switch status {
  case Running:
    str = "Running"
  case Except:
    str = "Except"
  case OffLine:
    str = "Status"
  }
  return json.Marshal(str)
}

type Host struct {
  ID     uint   `json:"id"`
  Name   string `json:"name"`
  Status Status `json:"status"`
}

const (
  Running Status = 1
  Except  Status = 2
  OffLine Status = 3
)

func main() {
  host := Host{1, "枫枫", Running}
  data, _ := json.Marshal(host)
  fmt.Println(string(data)) // {"id":1,"name":"枫枫","status":"Running"}
}
```

嗯，代码简洁了不少，在使用层面已经没有问题了

但是，这个结构体怎么表示数据库中的字段呢？

golang中没有枚举，我们只能自己通过逻辑实现枚举

```go
type Weekday int

const (
  Sunday    Weekday = iota + 1 // EnumIndex = 1
  Monday                       // EnumIndex = 2
  Tuesday                      // EnumIndex = 3
  Wednesday                    // EnumIndex = 4
  Thursday                     // EnumIndex = 5
  Friday                       // EnumIndex = 6
  Saturday                     // EnumIndex = 7
)

var WeekStringList = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var WeekTypeList = []Weekday{Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday}

// String 转字符串
func (w Weekday) String() string {
  return WeekStringList[w-1]
}

// MarshalJSON 自定义类型转换为json
func (w Weekday) MarshalJSON() ([]byte, error) {
  return json.Marshal(w.String())
}

// EnumIndex 自定义类型转原始类型
func (w Weekday) EnumIndex() int {
  return int(w)
}

// ParseWeekDay 字符串转自定义类型
func ParseWeekDay(week string) Weekday {
  for i, i2 := range WeekStringList {
    if week == i2 {
      return WeekTypeList[i]
    }
  }
  return Monday
}

// ParseIntWeekDay 数字转自定义类型
func ParseIntWeekDay(week int) Weekday {
  return Weekday(week)
}

type DayInfo struct {
  Weekday Weekday   `json:"weekday"`
  Date    time.Time `json:"date"`
}

func main() {
  w := Sunday
  fmt.Println(w)
  dayInfo := DayInfo{Weekday: Sunday, Date: time.Now()}
  data, err := json.Marshal(dayInfo)
  fmt.Println(string(data), err)
  week := ParseWeekDay("Sunday")
  fmt.Println(week)
  week = ParseIntWeekDay(2)
  fmt.Println(week)
}

```

在需要输出的时候（print，json），自定义类型就变成了字符串

从外界接收的数据也能转换为自定义类型，这就是golang中的枚举，假枚举



## 事务处理

事务就是用户定义的一系列数据库操作，这些操作可以视为一个完成的逻辑处理工作单元，要么全部执行，要么全部不执行，是不可分割的工作单元。

很形象的一个例子，张三给李四转账100元，在程序里面，张三的余额就要-100，李四的余额就要+100 整个事件是一个整体，哪一步错了，整个事件都是失败的

gorm事务默认是开启的。为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。

如果没有这方面的要求，您可以在初始化时禁用它，这将获得大约 30%+ 性能提升。

一般不推荐禁用

```go

// 全局禁用
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  SkipDefaultTransaction: true,
})

```

本节课表结构

```go

type User struct {
  ID    uint   `json:"id"`
  Name  string `json:"name"`
  Money int    `json:"money"`
}

// InnoDB引擎才支持事务，MyISAM不支持事务
// DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

```

### 普通事务

以张三给李四转账为例，不使用事务的后果

```go
var zhangsan, lisi User
DB.Take(&zhangsan, "name = ?", "张三")
DB.Take(&lisi, "name = ?", "李四")
// 张三给李四转账100元
// 先给张三-100
zhangsan.Money -= 100
DB.Model(&zhangsan).Update("money", zhangsan.Money)
// 模拟失败的情况

// 再给李四+100
lisi.Money += 100
DB.Model(&lisi).Update("money", lisi.Money)
```

在失败的情况下，要么张三白白损失了100，要么李四凭空拿到100元

这显然是不合逻辑的，并且不合法的

那么，使用事务是怎样的

```go
var zhangsan, lisi User
DB.Take(&zhangsan, "name = ?", "张三")
DB.Take(&lisi, "name = ?", "李四")
// 张三给李四转账100元
DB.Transaction(func(tx *gorm.DB) error {

  // 先给张三-100
  zhangsan.Money -= 100
  err := tx.Model(&zhangsan).Update("money", zhangsan.Money).Error
  if err != nil {
    fmt.Println(err)
    return err
  }

  // 再给李四+100
  lisi.Money += 100
  err = tx.Model(&lisi).Update("money", lisi.Money).Error
  if err != nil {
    fmt.Println(err)
    return err
  }
  // 提交事务
  return nil
})
```

使用事务之后，他们就是一体，一起成功，一起失败

### 手动事务

```go
// 开始事务
tx := db.Begin()

// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
tx.Create(...)

// ...

// 遇到错误时回滚事务
tx.Rollback()

// 否则，提交事务
tx.Commit()
```

刚才的代码也可以这样实现

```go
var zhangsan, lisi User
DB.Take(&zhangsan, "name = ?", "张三")
DB.Take(&lisi, "name = ?", "李四")

// 张三给李四转账100元
tx := DB.Begin()

// 先给张三-100
zhangsan.Money -= 100
err := tx.Model(&zhangsan).Update("money", zhangsan.Money).Error
if err != nil {
  tx.Rollback()
}

// 再给李四+100
lisi.Money += 100
err = tx.Model(&lisi).Update("money", lisi.Money).Error
if err != nil {
  tx.Rollback()
}
// 提交事务
tx.Commit()
```



## 钩子

GORM提供了多种钩子函数，允许在CRUD操作的不同阶段执行自定义逻辑。

### 定义钩子

可以在模型结构体中定义钩子方法：

```go
// 创建操作前回调
func (u *student) BeforeCreate(tx *gorm.DB) (err error) {
    u.Name = "Default Name"
    return
}

// 创建操作后回调
func (u *student) AfterCreate(tx *gorm.DB) (err error) {
    log.Println("student created:", u.ID)
    return
}

// 更新操作前回调
func (u *student) BeforeUpdate(tx *gorm.DB) (err error) {
    log.Println("student created:", u.ID)
    return
}
// 更新操作后回调
func (u *student) AfterUpdate(tx *gorm.DB) (err error) {
    log.Println("student created:", u.ID)
    return
}

// 保存操作前回调
func (u *student) BeforeSave(tx *gorm.DB) (err error) {
    log.Println("student created:", u.ID)
    return
}
// 保存操作后回调
func (u *student) AfterSave(tx *gorm.DB) (err error) {
    log.Println("student created:", u.ID)
    return
}

// 删除操作前回调
func (u *student) BeforeDelete(tx *gorm.DB) (err error) {
    log.Println("student created:", u.ID)
    return
}
// 删除操作后回调
func (u *student) AfterDelete(tx *gorm.DB) (err error) {
    log.Println("student created:", u.ID)
    return
}

// find 操作后回调
func (u *student) AfterFind(tx *gorm.DB) (err error) {
    log.Println("student created:", u.ID)
    return
}
```

在上面的代码中，我们定义了BeforeCreate和AfterCreate钩子方法，分别在创建记录之前和之后执行自定义逻辑。



## 性能优化

GORM提供了一些性能优化的技巧，如批量插入、预编译语句和缓存等。

### 预编译语句

可以使用PrepareStmt方法预编译SQL语句：

```go
db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
    PrepareStmt: true,
})
```

在上面的代码中，我们在初始化gorm.DB对象时启用了预编译语句功能。

### 缓存

可以使用CacheStore方法缓存查询结果：

```go
var students []student
cache := gorm.CacheStore(db, "students_cache", time.Minute*10)
result := cache.Find(&students)
if result.Error != nil {
    log.Fatal(result.Error)
}
```

在上面的代码中，我们使用CacheStore方法缓存查询结果，有效期为10分钟。


## 常见问题与解决方案

在使用GORM时，可能会遇到一些常见问题。以下是一些常见问题及其解决方案：

### 数据库连接失败

确保DSN字符串正确，并检查数据库服务是否运行：

```go
dsn := "student:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatal(err)
}
```

### 模型定义不正确

确保模型结构体字段与数据库表列对应：

```go
type student struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"size:255"`
    Email     string `gorm:"uniqueIndex"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 查询结果为空

检查查询条件是否正确，并确保数据库中有匹配的记录：

```go
var student student
result := db.Where("name = ?", "John").First(&student)
if result.Error != nil {
    log.Fatal(result.Error)
}
```

### 事务回滚失败

确保事务处理逻辑正确，检查是否有未捕获的错误：

```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&student{Name: "John"}).Error; err != nil {
        return err
    }

    if err := tx.Create(&Order{Amount: 1000, studentID: 1}).Error; err != nil {
        return err
    }

    return nil
})

if err != nil {
    log.Fatal(err)
}
```