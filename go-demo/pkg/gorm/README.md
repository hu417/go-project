

# gorm

参考: https://gorm.io/zh_CN/docs/index.html


## 连接

安装依赖: 
- `go get gorm.io/driver/mysql`
- `go get gorm.io/gorm`

### 基本连接

```go

username := "root"  //账号
password := "root"  //密码
host := "127.0.0.1" //数据库地址，可以是Ip或者域名
port := 3306        //数据库端口
Dbname := "gorm"   //数据库名
timeout := "10s"    //连接超时，10秒

// root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&loc=Asia/Shanghai
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)

//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
db, err := gorm.Open(mysql.Open(dsn))
if err != nil {
  panic("连接数据库失败, error=" + err.Error())
}
// 连接成功
fmt.Println(db)

```

### 高级配置

- 跳过默认事务
为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它，这样可以获得60%的性能提升
```go
db, err := gorm.Open(mysql.Open("gorm.db"), &gorm.Config{
  SkipDefaultTransaction: true,
})
```

- 命名策略
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

我们也可以修改这些策略

```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
  NamingStrategy: schema.NamingStrategy{
    TablePrefix:   "f_",  // 表名前缀
    SingularTable: false, // 单数表名
    NoLowerCase:   false, // 关闭小写转换
  },
})
```

- 显示日志
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

如果你想自定义日志的显示,那么可以使用如下代码

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
```
DB.Debug().First(&model)

```

## 模型定义

模型是标准的 struct，由 Go 的基本数据类型、实现了 Scanner 和 Valuer 接口的自定义类型及其指针或别名组成

定义一张表
```go
type Student struct {
  ID    uint // 默认使用ID作为主键
  Name  string
  Email *string // 使用指针是为了存空值
}
```
常识：小写属性是不会生成字段的

### 自动生成表结构

```go
// 可以放多个
DB.AutoMigrate(&Student{})
```

AutoMigrate的逻辑是只新增，不删除，不修改（大小会修改）

例如将Name修改为Name1，进行迁移，会多出一个name1的字段

生成的表结构如下
```sql
CREATE TABLE `f_students` (`id` bigint unsigned AUTO_INCREMENT,`name` longtext,`email` longtext,PRIMARY KEY (`id`))
```
默认的类型太大了

### 修改大小

我们可以使用gorm的标签进行修改

有两种方式
```
Name  string  `gorm:"type:varchar(12)"`
Name  string  `gorm:"size:2"`
```
字段标签
- type 定义字段类型
- size 定义字段大小
- column 自定义列名
- primaryKey 将列定义为主键
- unique 将列定义为唯一键
- default 定义列的默认值
- not null 不可为空
- embedded 嵌套字段
- embeddedPrefix 嵌套字段前缀
- comment 注释

多个标签之前用 ; 连接
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


## 基本使用

### 单表查询
使用gorm对单张表进行增删改查

表结构

type Student struct {
  ID     uint   `gorm:"size:3"`
  Name   string `gorm:"size:8"`
  Age    int    `gorm:"size:3"`
  Gender bool
  Email  *string `gorm:"size:32"`
}
```
添加记录
email := "xxx@qq.com"
// 创建记录
student := Student{
  Name:   "枫枫",
  Age:    21,
  Gender: true,
  Email:  &email,
}
DB.Create(&student)
```
有两个地方需要注意

指针类型是为了更好的存null类型，但是传值的时候，也记得传指针
Create接收的是一个指针，而不是值
由于我们传递的是一个指针，调用完Create之后，student这个对象上面就有该记录的信息了，如创建的id

DB.Create(&student)
fmt.Printf("%#v\n", student)  
// main.Student{ID:0x2, Name:"zhangsan", Age:23, Gender:false, Email:(*string)(0x11d40980)}
```
批量插入
Create方法还可以用于插入多条记录

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
查询单条记录
var student Student
DB.Take(&student)
fmt.Println(student)
```
获取单条记录的方法很多，我们对比sql就很直观了

DB = DB.Session(&gorm.Session{Logger: Log})
var student Student
DB.Take(&student)  
// SELECT * FROM `students` LIMIT 1
DB.First(&student) 
// SELECT * FROM `students` ORDER BY `students`.`id` LIMIT 1
DB.Last(&student)  
// SELECT * FROM `students` ORDER BY `students`.`id` DESC LIMIT 1
```
根据主键查询
var student Student
DB.Take(&student, 2)
fmt.Println(student)

student = Student{} // 重新赋值
DB.Take(&student, "4")
fmt.Println(student)
```
Take的第二个参数，默认会根据主键查询，可以是字符串，可以是数字

根据其他条件查询
var student Student
DB.Take(&student, "name = ?", "机器人27号")
fmt.Println(student)
```
使用？作为占位符，将查询的内容放入?

SELECT * FROM `students` WHERE name = '机器人27号' LIMIT 1
```
这样可以有效的防止sql注入

他的原理就是将参数全部转义，如

DB.Take(&student, "name = ?", "机器人27号' or 1=1;#")

SELECT * FROM `students` WHERE name = '机器人27号\' or 1=1;#' LIMIT 1
```
根据struct查询
var student Student
// 只能有一个主要值
student.ID = 2
//student.Name = "枫枫"
DB.Take(&student)
fmt.Println(student)
```
获取查询结果
获取查询的记录数

count := DB.Find(&studentList).RowsAffected
```
是否查询失败

err := DB.Find(&studentList).Error
```
查询失败有查询为空，查询条件错误，sql语法错误

可以使用判断

var student Student
err := DB.Take(&student, "xx").Error
switch err {
case gorm.ErrRecordNotFound:
  fmt.Println("没有找到")
default:
  fmt.Println("sql错误")
}
```
查询多条记录
var studentList []Student
DB.Find(&studentList)
for _, student := range studentList {
  fmt.Println(student)
}

// 由于email是指针类型，所以看不到实际的内容
// 但是序列化之后，会转换为我们可以看得懂的方式
var studentList []Student
DB.Find(&studentList)
for _, student := range studentList {

  data, _ := json.Marshal(student)
  fmt.Println(string(data))
}
```
根据主键列表查询
var studentList []Student
DB.Find(&studentList, []int{1, 3, 5, 7})
DB.Find(&studentList, 1, 3, 5, 7)  // 一样的
fmt.Println(studentList)
```
根据其他条件查询
DB.Find(&studentList, "name in ?", []string{"枫枫", "zhangsan"})
```
更新
更新的前提的先查询到记录

Save保存所有字段
用于单个记录的全字段更新

它会保存所有字段，即使零值也会保存

var student Student
DB.Take(&student)
student.Age = 23
// 全字段更新
DB.Save(&student)
// UPDATE `students` SET `name`='枫枫',`age`=23,`gender`=true,`email`='xxx@qq.com' WHERE `id` = 1
```
零值也会更新

var student Student
DB.Take(&student)
student.Age = 0
// 全字段更新
DB.Save(&student)
// UPDATE `students` SET `name`='枫枫',`age`=0,`gender`=true,`email`='xxx@qq.com' WHERE `id` = 1
```
更新指定字段
可以使用select选择要更新的字段

var student Student
DB.Take(&student)
student.Age = 21
// 全字段更新
DB.Select("age").Save(&student)
// UPDATE `students` SET `age`=21 WHERE `id` = 1
```
批量更新
例如我想给年龄21的学生，都更新一下邮箱

var studentList []Student
DB.Find(&studentList, "age = ?", 21).Update("email", "is21@qq.com")
```
还有一种更简单的方式

DB.Model(&Student{}).Where("age = ?", 21).Update("email", "is21@qq.com")
// UPDATE `students` SET `email`='is22@qq.com' WHERE age = 21
```
这样的更新方式也是可以更新零值的

更新多列
如果是结构体，它默认不会更新零值

email := "xxx@qq.com"
DB.Model(&Student{}).Where("age = ?", 21).Updates(Student{
  Email:  &email,
  Gender: false,  // 这个不会更新
})

// UPDATE `students` SET `email`='xxx@qq.com' WHERE age = 21
```
如果想让他更新零值，用select就好

email := "xxx1@qq.com"
DB.Model(&Student{}).Where("age = ?", 21).Select("gender", "email").Updates(Student{
  Email:  &email,
  Gender: false,
})
// UPDATE `students` SET `gender`=false,`email`='xxx1@qq.com' WHERE age = 21
```
如果不想多写几行代码，则推荐使用map

DB.Model(&Student{}).Where("age = ?", 21).Updates(map[string]any{
  "email":  &email,
  "gender": false,
})
```
更新选定字段
Select选定字段

Omit忽略字段

删除
根据结构体删除

// student 的 ID 是 `10`
db.Delete(&student)
// DELETE from students where id = 10;
```
删除多个

db.Delete(&Student{}, []int{1,2,3})

// 查询到的切片列表
db.Delete(&studentList)
```


### 关联查询

#### 一对一

1、创建一对一model

- models/article.go
```go
package models

import (
    _ "github.com/jinzhu/gorm"
)

type Article struct {
  Id          int         `json:"id"`
  Title       string      `json:"title"`
  CateId      string      `json:"cate_id"`
  State       int         `json:"state"`
  ArticleCate ArticleCate `gorm:"foreignkey:Id;association_foreignkey:CateId"`
   // Id指的是 ArticleCate 表的字段名Id
   // CateId 指的是自己表中的 字段名，关联ArticleCate表的Id
}

func (Article) TableName() string {
    return "article"
}

```

- models/articleCate.go
```go

package models
import (
  _ "github.com/jinzhu/gorm"
)

type ArticleCate struct {
  Id      int       `json:"id"`
  Title   string    `json:"title"`
  State   int       `json:"state"`
}

func (ArticleCate) TableName() string {
  return "article_cate"
}

```

2、一对一表查询
```go

package controllers
import (
  "beegogorm/models"
  "github.com/astaxie/beego"
)

type ArticleController struct {
  beego.Controller
}

func (c *ArticleController) ArticleOrm() {
  // 1、查询文章信息的时候关联文章分类  (1对1)
  article := []models.Article{}
  models.DB.Preload("ArticleCate").Find(&article)

  // 2、查询文章信息的时候关联文章分类  (1对1) 添加过来条件
  article := []models.Article{}
  models.DB.Preload("ArticleCate").Where("id>2").Find(&article)

  c.Data["json"] = article
  c.ServeJSON()
}
```


3、查询结果
> Article(文章表) 和 ArticleCate文章分类表
> Article(文章表) 和 ArticleCate文章分类表
>
```go

[
  {
    "id": 1,
    "title": "西游记",
    "cate_id": "1",
    "state": 1,
    "ArticleCate": {
      "id": 1,
      "title": "四大名著",
      "state": 1
    }
  },
  {
    "id": 2,
    "title": "三国演义",
    "cate_id": "1",
    "state": 1,
    "ArticleCate": {
      "id": 1,
      "title": "四大名著",
      "state": 1
    }
  },
  {
    "id": 3,
    "title": "货币战争",
    "cate_id": "2",
    "state": 1,
    "ArticleCate": {
      "id": 2,
      "title": "国外名著",
      "state": 1
    }
  },
  {
    "id": 4,
    "title": "钢铁是怎样炼成的",
    "cate_id": "2",
    "state": 1,
    "ArticleCate": {
      "id": 2,
      "title": "国外名著",
      "state": 1
    }
  }
]

```

4、可以使用自动创建表
- models/main.go
```go

package models

import (
  "github.com/astaxie/beego"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var err error

func init() {
  //和数据库建立连接
  DB, err = gorm.Open("mysql", "root:chnsys@2016@/beegodb?charset=utf8&parseTime=True&loc=Local")
  if err != nil {
    beego.Error()
  }

  //// 创建表
  //DB.CreateTable(&Article{},&ArticleCate{})     // 根据User结构体建表
  //// 设置表结构的存储引擎为InnoDB
  //DB.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Article{},ArticleCate{})
}

```

#### 一对多
1、创建一对多model
- models/article.go
```go
package models
import (
_ "github.com/jinzhu/gorm"
)

type Article struct {
  Id          int         `json:"id"`
  Title       string      `json:"title"`
  CateId      string      `json:"cate_id"`
  State       int         `json:"state"`
}

func (Article) TableName() string {
  return "article"
}

```

- models/articleCate.go
```go
package models
import (
  _ "github.com/jinzhu/gorm"
)

type ArticleCate struct {
  Id      int       `json:"id"`
  Title   string    `json:"title"`
  State   int       `json:"state"`
}

func (ArticleCate) TableName() string {
  return "article_cate"
}
```

2、一对多表查询
```go
func (c *ArticleController) ArticleOrm() {
  // 1、查询文章分类信息的时候关联文章  (1对多)   查询文章分类显示文章信息
  articleCate := []models.ArticleCate{}
  models.DB.Preload("Article").Find(&articleCate)

  // 4、查询文章分类信息的时候关联文章   条件判断
  // articleCate := []models.ArticleCate{}
  // models.DB.Preload("Article").Where("id>1").Find(&articleCate)

  c.Data["json"] = articleCate
  c.ServeJSON()
}
```

3、查询结果
> Article(文章表) 和 ArticleCate文章分类表
> 一个文章分类下面，有多篇文章

```go

[
  {
    "id": 1,
    "title": "四大名著",
    "state": 1,
    "Article": [
      {
        "id": 1,
        "title": "西游记",
        "cate_id": "1",
        "state": 1,
        "ArticleCate": {
          "id": 0,
          "title": "",
          "state": 0,
          "Article": null
        }
      },
      {
        "id": 2,
        "title": "三国演义",
        "cate_id": "1",
        "state": 1,
        "ArticleCate": {
          "id": 0,
          "title": "",
          "state": 0,
          "Article": null
        }
      }
    ]
  },
  {
    "id": 2,
    "title": "国外名著",
    "state": 1,
    "Article": [
      {
        "id": 3,
        "title": "货币战争",
        "cate_id": "2",
        "state": 1,
        "ArticleCate": {
          "id": 0,
          "title": "",
          "state": 0,
          "Article": null
        }
      },
      {
        "id": 4,
        "title": "钢铁是怎样炼成的",
        "cate_id": "2",
        "state": 1,
        "ArticleCate": {
          "id": 0,
          "title": "",
          "state": 0,
          "Article": null
        }
      }
    ]
  }
]
```

#### 多对多

1、创建多对多model
- models/lesson.go
```go
package models
import (
  _ "github.com/jinzhu/gorm"
)

type Lesson struct {
  Id      int       `json:"id"`
  Name    string    `json:"name"`
  Student []Student `gorm:"many2many:lesson_student;"`
  // lesson_student => 表名（第三张关联表表名）
}

func (Lesson) TableName() string {
  return "lesson"
}
```

- models/student.go
```go
package models
import (
  _ "github.com/jinzhu/gorm"
)

type ArticleCate struct {
  Id      int       `json:"id"`
  Title   string    `json:"title"`
  State   int       `json:"state"`
}

func (ArticleCate) TableName() string {
  return "article_cate"
}
```

- models/lessonStudent.go
```go
package models
import (
  _ "github.com/jinzhu/gorm"
)

type LessonStudent struct {
  LessonId  int `json:"lesson_id"`
  StudentId int `json:"student_id"`
}

func (LessonStudent) TableName() string {
  return "lesson_student"
}
```

2、多对多表查询
```go
func (c *StudentController) StudentM2M() {
  ////1、获取学生信息
  studentList := []models.Student{}
  models.DB.Find(&studentList)
  c.Data["json"] = studentList
  c.ServeJSON()

  ////2、获取课程信息
  lessonList := []models.Lesson{}
  models.DB.Find(&lessonList)
  c.Data["json"] = lessonList
  c.ServeJSON()

  //3、查询学生信息的时候获取学生的选课信息
  studentList := []models.Student{}
  models.DB.Preload("Lesson").Find(&studentList)
  c.Data["json"] = studentList
  c.ServeJSON()

  //4、查询张三选修了哪些课程
  studentList := []models.Student{}
  models.DB.Preload("Lesson").Where("id=1").Find(&studentList)
  c.Data["json"] = studentList
  c.ServeJSON()

  //5、课程被哪些学生选修了
  lessonList := []models.Lesson{}
  models.DB.Preload("Student").Find(&lessonList)
  c.Data["json"] = lessonList
  c.ServeJSON()

  //6、计算机网络被那些学生选修了
  lessonList := []models.Lesson{}
  models.DB.Preload("Student").Where("id=1").Find(&lessonList)
  c.Data["json"] = lessonList
  c.ServeJSON()

  //7、条件
  lessonList := []models.Lesson{}
  models.DB.Preload("Student").Offset(1).Limit(2).Find(&lessonList)
  c.Data["json"] = lessonList
  c.ServeJSON()

  //8、张三被开除了  查询课程被哪些学生选修的时候要去掉张三
  lessonList := []models.Lesson{}
  models.DB.Preload("Student", "id!=1").Find(&lessonList)
  c.Data["json"] = lessonList
  c.ServeJSON()

  lessonList := []models.Lesson{}
  models.DB.Preload("Student", "id not in (1,2)").Find(&lessonList)
  c.Data["json"] = lessonList
  c.ServeJSON()

  //9、查看课程被哪些学生选修  要求：学生id倒叙输出   自定义预加载 SQL
  //https://gorm.io/zh_CN/docs/preload.html
  lessonList := []models.Lesson{}
  models.DB.Preload("Student", func(db *gorm.DB) *gorm.DB {
    return models.DB.Order("id DESC")
  }).Find(&lessonList)
  c.Data["json"] = lessonList
  c.ServeJSON()

  lessonList := []models.Lesson{}
  models.DB.Preload("Student", func(db *gorm.DB) *gorm.DB {
    return models.DB.Where("id>3").Order("id DESC")
  }).Find(&lessonList)
  
  c.Data["json"] = lessonList
  c.ServeJSON()
}
```

3、查询结果
> 学生表(Student) 和 课程表(Lesson)
> 一个学生可以选修多个课程，一个课程也可以被多个学生选修

```go
[
  {
    "Id": 1,
    "Number": "12",
    "Password": "123456",
    "ClassId": 1,
    "Name": "zhangsan",
    "Lesson": [
      {
        "id": 1,
        "name": "语文",
        "Student": null
      }
    ]
  },
  {
    "Id": 2,
    "Number": "24",
    "Password": "123456",
    "ClassId": 1,
    "Name": "lisi",
    "Lesson": [
      {
        "id": 1,
        "name": "语文",
        "Student": null
      }
    ]
  },
  {
    "Id": 3,
    "Number": "22",
    "Password": "123456",
    "ClassId": 1,
    "Name": "wangwu",
    "Lesson": [
      {
        "id": 2,
        "name": "数学",
        "Student": null
      }
    ]
  }
]
```


### 原生sql

1、model结构体
- models/user.go
```go
package models

type User struct {
  Id       int
  Username string
  Age      int
  Email    string
  AddTime  int
}

//定义结构体操作的数据库表
func (User) TableName() string {
  return "user"
}
```

2、增删改 （原生SQL）
```go
func (c *UserController) UserSelect() {
  // 1、使用原生sql给user表增加一条数据
  res := models.DB.Exec("insert into user(username,age,email) values(?,?,?)", "6666", 11, "xxxxx@qq.com")

  // 2、使用原生sql删除user表中的一条数据
  res := models.DB.Exec("delete from user where id =?", 6)

  // 3、使用原生sql修改user表中的一条数据
  res := models.DB.Exec("update user set username=? where id=?", "zhangsannew", 1)

  c.Data["json"] = res
  c.ServeJSON()
}
```

3、原生SQL查询

原生sql基本查询
```go
func (c *UserController) UserSelect() {
  user := []models.User{}

  // 1、查询User表中所有的数据
  models.DB.Raw("select * from user").Scan(&user)

  // 2、查询uid=2的数据
  models.DB.Raw("select * from user where id=?", 2).Scan(&user)

  c.Data["json"] = user
  c.ServeJSON()
}
```

打印查询的数据
```go
func (c *UserController) UserSelect() {

  // 1、统计user表的数量
  var num int
  res := models.DB.Raw("select count(1) from user").Row()
  res.Scan(&num)   // 5
  fmt.Println("数据库表里面有", num, "条数据")

  // 2、把查询到的数据赋值给变量
  var username string
  var email string
  res := models.DB.Raw("select username,email from user where id=1").Row()
  res.Scan(&username, &email)
  fmt.Println(username, email) // Snail Snail.qq.com

  // 3、打印user表的所有数据
  var username string
  var email string
  rows, _ := models.DB.Raw("select username,email from user").Rows()
  defer rows.Close()   //操作完毕关闭
  for rows.Next() {
    rows.Scan(&username, &email)
    fmt.Println(username, email)
  }
}
```

## 测试

## 测试

