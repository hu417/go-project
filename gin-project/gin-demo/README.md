# gin学习



官方文档: https://gin-gonic.com/zh-cn/docs/

项目参考: https://github.com/go-programming-tour-book/blog-service



## 简介

Gin 是一个基于 Go 语言的 Web 框架，它具有高性能、易学易用、轻量级等特点，被广泛应用于构建 RESTful API 和微服务等场景。

Gin 框架提供了丰富的中间件支持，可以方便地实现请求路由、参数解析、日志记录、错误处理等功能。

Gin 框架的设计灵感来自于 Martini 框架，但相比之下，Gin 框架更快、更稳定、更易用。

Gin 具有类似于 Martini 的 API 风格，并且它使用了著名的开源项目 httprouter 的自定义版本作为路由基础，使得它的性能表现更高更好，相较 Martini 大约提高了 40 倍。

另外 gin 除了快以外，还具备小巧、精美且易用的特性，目前广受 Go 语言开发者的喜爱，是最流行的 HTTP Web 框架



## 安装和使用

要安装 Gin 框架，你需要先安装 Go 语言环境，并设置好 GOPATH 和 PATH 环境变量。然后可以使用以下命令安装 Gin：

```go
go get -u github.com/gin-gonic/gin
```

这将从 GitHub 下载 Gin 框架的最新版本并安装到你的 GOPATH 目录下。

在安装完毕后，我们可以看到项目根目录下的 go.mod 文件也会发生相应的改变.

打开 go.mod 文件，查看如下：

```go
module github.com/go-programming-tour-book/blog-service

go 1.19
require (
    github.com/gin-gonic/gin v1.10.0 // indirect
    ...
)
```

这些正正就是 gin 所相关联的所有模块包.

大家可能会好奇，为什么 `github.com/gin-gonic/gin` 后面会出现 indirect 标识，

在执行命令`go get `时，Go module 会自动整理`go.mod` 文件，如果有必要会在部分依赖包的后面增加`// indirect`注释。

一般而言，被添加`// indirect` 注释的包,肯定是间接依赖的包，

而没有添加`// indirect`注释的包则是直接依赖的包，什么叫做直接依赖呢 ？即明确的出现在某个`import`语句中。

然而，需要着重强调的是：**并不是所有的间接依赖都会出现在 `go.mod`文件中。**间接依赖出现在`go.mod`文件的情况，可能符合下面所列场景的一种或多种：

- 直接依赖未启用 Go module
- 直接依赖go.mod 文件中缺失部分依赖

回到上面的  go.mod 文件，查看如下：

```go
module github.com/go-programming-tour-book/blog-service
go 1.19
require (
    github.com/gin-gonic/gin v1.9.0 // indirect
    github.com/go-playground/universal-translator v0.18.1 // indirect
    ...
)
```

它明明是我们直接通过调用 `go get` 引用的， 为啥加了  `// indirect`注释 ？

是因为在我们安装时，这个项目模块还没有真正的去使用它所导致的（还没哟在某个`import`语句中 出现 ）。

另外你会注意到，在 go.mod 文件中有类似 `go 1.14` 这样的标识位，主要与你创建 Go modules 时的 Go 版本有关。

### 使用 Gin 框架编写 Hello World

在完成前置动作后，在本节我们先将一个 Demo 运行起来，看看一个最简单的 HTTP 服务运行起来是怎么样的，

使用 Gin 框架编写 Web 应用程序非常简单，以下是一个简单的示例：

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    router.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, World!")
    })

    router.Run(":8080")
}
```

这个应用程序创建了一个简单的 HTTP 服务器，监听本地的 8080 端口，并在访问根路径时返回 "Hello, World!"。

你可以使用 `go run` 命令运行这个应用程序：

```
go run main.go
```

也可以在ide 工具中直接启动 



然后在浏览器中访问 http://localhost:8080 就可以看到 "Hello, World!" 的响应了。

当然，这只是 Gin 框架的一个简单示例，你可以根据自己的需求编写更复杂的 Web 应用程序。

### main.go的执行过程

接下来我们运行 main.go 文件，查看运行结果，如下：

```go
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func1 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080
```

我们可以看到启动了服务后，输出了许多运行信息，

在这里我们对运行信息做一个初步的概括分析，分为以下四大块：

首先是gin的运行模式：当前为Release 模式。

```go
 - using code:  gin.SetMode(gin.ReleaseMode)
```

这是 Go 语言中使用 Gin 框架时，设置运行模式为 Release 模式的代码。

在 Release 模式下，Gin 框架会关闭调试信息和堆栈跟踪，以提高性能和安全性。

这个设置通常在应用程序的 main 函数或初始化代码中进行。

并建议若在测试环境时切换为debug模式，gin.SetMode(gin.DebugMode) 切换为debug模式，

```go
func main() {
 router := gin.Default()
 gin.SetMode(gin.DebugMode)
 router.GET("/", func(c *gin.Context) {
  c.String(http.StatusOK, "Hello, World!")
 })

 router.Run(":8080")
}
```

接下来，是请求的路由注册：注册了 `GET /ping` 的路由，并输出其调用方法的方法名。

```go
[GIN-debug] GET    /    --> main.main.func1 (3 handlers)
```

接下来，是监听端口信息：本次启动时监听 8080 端口，由于没有设置端口号等信息，因此默认为 8080。

```go
[GIN-debug] Listening and serving HTTP on :8080
```



## Rest微服务的模块设计

在完成了初步的示例演示后，接下来就是进入具体的预备开发阶段，一般在正式进入业务开发前，我们会针对本次需求的迭代内容进行多类的设计和评审，无设计不开发。

但是问题在于，我们目前还缺很多初始化的东西没有做，因此在本章节中，我们主要针对项目目录结构、接口方案、路由注册、数据库等设计进行思考和设计开发。

### Rest微服务目录结构

我们先将项目的标准目录结构创建起来，便于后续的开发，最终目录结构如下：

```
gin-rest
├── configs
├── docs
├── global
├── internal
│   ├── dao
│   ├── middleware
│   ├── model
│   ├── routers
│   └── service
├── pkg
├── storage
├── scripts
└── third_party
```

- **configs**：配置文件。

- **docs**：文档集合。

- **global**：全局变量。

- **internal**：内部模块。

- - **dao**：数据访问层（Database Access Object），所有与数据相关的操作都会在 dao 层进行，例如 MySQL、ElasticSearch 等。
  - **middleware**：HTTP 中间件。
  - **model**：模型层，用于存放 model 对象。
  - **routers**：路由相关逻辑处理。
  - **service**：项目核心业务逻辑。

- **pkg**：项目相关的模块包。

- **storage**：项目生成的临时文件。

- **scripts**：各类构建，安装，分析等操作的脚本。

- **third_party**：第三方的资源工具。

## Gin框架WEB开发入门

### 基本安装

1.首先需要安装Go（需要1.10+版本），然后可以使用下面的Go命令安装Gin。

> `go get -u github.com/gin-gonic/gin`

2.将其导入您的代码中：

> `import "github.com/gin-gonic/gin"`

### Rest请求路由的配置和使用

路由的本质是前缀树，利用前缀树来实现路由的功能。

### Rest请求路径的设置

Gin框架的Rest路由使用非常简单，可以通过定义路由以及处理该路由对应的Handler来接收用户的Web请求。

以下是一个使用Gin框架的路由示例：

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    // 实例化一个GIN对象
    router := gin.Default()
  
		// 设置路由
    router.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, World!")
    })
		// 遵循Restful风格（采用URL定位，HTTP描述操作）
    // 第一个参数是：路径； 第二个参数是：具体操作 func(c *gin.Context)
    router.GET("/Get", getting)
    router.POST("/Post", posting)
    router.PUT("/Put", putting)
    router.DELETE("/Delete", deleting)
  
    // 默认启动的是 8080端口
    router.Run(":8080")
}

//定义handler方法, 类似于java中的 Controller
func getting(c *gin.Context) {
 c.JSON(http.StatusOK,gin.H{
  "message":"getting",
 })
}

//定义handler方法, 类似于java中的 Controller
func posting(c *gin.Context) {
 c.JSON(http.StatusOK,gin.H{
  "message":"posting",
 })
}

```



### Gin配置路由的七个主要方法

Gin配置路由是为了处理HTTP请求。

Gin框架中都有和HTTP请求相互对应的方法来定义路由。而HTTP请求包含不同方法，包括`GET`,`POST`,`PUT`,`PATCH`,`OPTIONS`,`HEAD`,`DELETE`等七种方法。

- GET：从服务器取出资源（一项或多项）

- POST：在服务器新建一个资源

- PUT：在服务器更新资源（客户端提供完整资源数据）

- PATCH：在服务器更新资源（客户端提供需要修改的资源数据）

- DELETE：从服务器删除资源



下面是Gin一系列的 路由定义案例

```go
router := gin.New()
 
router.GET("/testGet",func(c *gin.Context){
    //处理逻辑
})
 
router.POST("/testPost",func(c *gin.Context){
    //处理逻辑
})
 
router.PUT("/testPut",func(c *gin.Context){
    //处理逻辑
})
 
router.DELETE("/testDelete",func(c *gin.Context){
    //处理逻辑
})
 
router.PATCH("/testPatch",func(c *gin.Context){
    //处理逻辑
})
 
router.OPTIONS("/testOptions",func(c *gin.Context){
    //处理逻辑
})
 
router.OPTIONS("/testHead",func(c *gin.Context){
    //处理逻辑
})
```



示例：操作文章

```go
// GET    /articles          文章列表
// GET    /articles/:id      文章详情
// POST   /articles          添加文章
// PUT    /articles/:id      修改某一篇文章
// DELETE /articles/:id      删除某一篇文章



package main

import (
  "encoding/json"
  "fmt"
  "github.com/gin-gonic/gin"
)

type ArticleModel struct {
  Title   string `json:"title"`
  Content string `json:"content"`
}

type Response struct {
  Code int    `json:"code"`
  Data any    `json:"data"`
  Msg  string `json:"msg"`
}

func _bindJson(c *gin.Context, obj any) (err error) {
  body, _ := c.GetRawData()
  contentType := c.GetHeader("Content-Type")
  switch contentType {
  case "application/json":
    err = json.Unmarshal(body, &obj)
    if err != nil {
      fmt.Println(err.Error())
      return err
    }
  }
  return nil
}

// _getList 文章列表页面
func _getList(c *gin.Context) {
  // 包含搜索，分页
  articleList := []ArticleModel{
    {"Go语言入门", "这篇文章是《Go语言入门》"},
    {"python语言入门", "这篇文章是《python语言入门》"},
    {"JavaScript语言入门", "这篇文章是《JavaScript语言入门》"},
  }
  c.JSON(200, Response{0, articleList, "成功"})
}

// _getDetail 文章详情
func _getDetail(c *gin.Context) {
  // 获取param中的id
  fmt.Println(c.Param("id"))
  article := ArticleModel{
    "Go语言入门", "这篇文章是《Go语言入门》",
  }
  c.JSON(200, Response{0, article, "成功"})
}

// _create 创建文章
func _create(c *gin.Context) {
  // 接收前端传递来的json数据
  var article ArticleModel

  err := _bindJson(c, &article)
  if err != nil {
    fmt.Println(err)
    return
  }

  c.JSON(200, Response{0, article, "添加成功"})
}

// _update 编辑文章
func _update(c *gin.Context) {
  fmt.Println(c.Param("id"))
  var article ArticleModel
  err := _bindJson(c, &article)
  if err != nil {
    fmt.Println(err)
    return
  }
  c.JSON(200, Response{0, article, "修改成功"})
}

// _delete 删除文章
func _delete(c *gin.Context) {
  fmt.Println(c.Param("id"))
  c.JSON(200, Response{0, map[string]string{}, "删除成功"})
}

func main() {
  router := gin.Default()
  router.GET("/articles", _getList)       // 文章列表
  router.GET("/articles/:id", _getDetail) // 文章详情
  router.POST("/articles", _create)       // 添加文章
  router.PUT("/articles/:id", _update)    // 编辑文章
  router.DELETE("/articles/:id", _delete) // 删除文章
  router.Run(":80")
}


```







### 路由分组

Group是一个路由分组器，可以将一组路由规则组织在一起，方便管理和维护。

Gin框架的路由分组可以通过Group函数实现。

下面是没有使用Group方法进行路由配置的例子：

```go
func main() {
 router := gin.Default()
 
 router.GET("/goods/list",goodsList)
 router.POST("/goods/add",createGoods)
 
 _ = router.Run()
}
```

使用路由分组改写，与上面的代码是同样的效果:

```go
func main() {
 router := gin.Default()
 
 goodsGroup := router.Group("/goods")
 {
   goodsGroup.GET("/list", goodsList)
   goodsGroup.GET("/add", createGoods)
 }
 
 
 _ = router.Run()
}
```

很多的业务代码，会定义url的不同版本，在这种场景下，就可以使用版本号来分组：

```go
// 两个路由组，都可以访问，大括号是为了保证规范
v1 := r.Group("/v1")
{
    // 通过 localhost:8080/v1/hello访问，以此类推
    v1.GET("/hello", sayHello)
    v1.GET("/world", sayWorld)
}
v2 := r.Group("/v2")
{
    v2.GET("/hello", sayHello)
    v2.GET("/world", sayWorld)
}
r.Run(":8080")
```

访问的时候，可以通过下面的方式，访问v1版本的地址：

```go
localhost:8080/v1/hello
localhost:8080/v1/world
```

访问的时候，可以通过下面的方式，访问v2版本的地址：

```go
localhost:8080/v2/hello
localhost:8080/v2/world
```

### 大规模路由的多文件配置

当我们的路由变得非常多的时候，那么建议遵循以下步骤：

1. 建立`routers`包，将不同模块拆分到`多个go文件`
2. 每个go文件提供`一个路由配置方法`，该方法注册实现一个分组的所有的路由
3. 之后main方法在调用文件的路由配置方法，实现路由注册

例子：第一个路由分组go文件：

/src/../routers/apiRouter.go

```go
package routers

// 这里是一个路由配置方法 ， routers包下某一个router对外开放的方法
func LoadRouter(e *gin.Engine) {
    e.Group("v1")
    {
        v1.GET("/post", postHandler)
    v1.GET("/get", getHandler)
    }
   ...
}
```

例子：第二个路由分组go文件：

```go
/src/../routers/uaaRouter.go
package routers

// 这里是一个路由配置方法 ， 是routers包下某一个router对外开放的方法
func LoadUaaRouter(e *gin.Engine) {
    e.Group("v1")
    {
        v1.GET("/post", postHandler)
    v1.GET("/get", getHandler)
    }
   ...
}
```

**main文件实现：**

```go
func main() {
    r := gin.Default()
    // 调用该方法实现注册
    routers.LoadRouter(r)
    routers.LoadUaaRouter(r) // 代表还有多个
    r.Run()
}
```

**规模如果继续扩大也有更好的处理方式（建议别太大，将服务拆分好）：**

项目规模更大的时候，我们可以遵循以下步骤：

1. 建立`routers`包，内部划分模块（包），每个包有个`router.go`文件，负责该模块的路由注册

2. 建立`setup_router.go`文件，并编写一个专用的路由初始化方法

3. main.go中按如下方式写入需要注册的路由，可进行路由的初始化

   

1、建立`routers`包，内部划分模块（包），每个包有个`router.go`文件，负责该模块的路由注册

```go
├── routers
│   │
│   ├── say
│   │   ├── sayWorld.go
│   │   └── router.go
│   │
│   ├── hello
│   │   ├── helloWorld.go
│   │   └── router.go
│   │
│   └── setup_router.go
│   
└── main.go
```

2、建立`setup_router.go`文件，并编写一个专用的路由初始化方法：

```go
type Register func(*gin.Engine)

func Init(routers ...Register) *gin.Engine {
 // 注册路由
 rs := append([]Register{}, routers...)

 r := gin.New()
 // 遍历调用方法
 for _, register := range rs {
  register(r)
 }
 return r
}
```

3、main.go中按如下方式写入需要注册的路由，可进行路由的初始化：

```go
func main() {
 // 设置需要加载的路由配置
 r := routers.Init(
    say.Routers,
    hello.Routers, // 后面还可以有多个
   )
 r.Run(":8080")
}
```



### 响应

#### 返回数据

gin常见的三种响应数据：`JSON`、`XML`、`YAML`

```go
// 1.JSON
r.GET("/someJSON", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "Json",
        "status":  200,
    })
})

// 2.XML
r.GET("/someXML", func(c *gin.Context) {
    c.XML(200, gin.H{"message": "abc"})
})

// 3.YAML
r.GET("/someYAML", func(c *gin.Context) {
    c.YAML(200, gin.H{"name": "zhangsan"})
})

// 4.protobuf
r.GET("/someProtoBuf", func(c *gin.Context) {
    reps := []int64{1, 2}
    data := &protoexample.Test{
        Reps:  reps,
    }
    c.ProtoBuf(200, data)
})

// 5.string
r.GET("/someString", func(c *gin.Context) {
  c.String(http.StatusOK, "返回string")
})

```

#### 文件响应

```go
// 在golang总，没有相对文件的路径，它只有相对项目的路径
// 网页请求这个静态目录的前缀， 第二个参数是一个目录，注意，前缀不要重复
r.StaticFS("/static", http.Dir("static/static"))
// 配置单个文件， 网页请求的路由，文件的路径
r.StaticFile("/titian.png", "static/titian.png")
```

#### 重定向

在 Gin 中，可以使用 `c.Redirect()` 方法进行重定向。该方法接受两个参数：重定向的目标 URL 和重定向的状态码。例如，以下代码将会将浏览器重定向到 `https://www.example.com`：

```go
func redirectHandler(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "https://www.example.com")
}
```

其中，`http.StatusMovedPermanently` 是一个常量，表示 301 状态码。你也可以使用其他状态码，例如 `http.StatusFound`（302）。

如果你想要在 URL 中包含查询参数，可以将它们添加到目标 URL 中。例如，以下代码将会将浏览器重定向到 `https://www.example.com?foo=bar`：

```go
func redirectHandler(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "https://www.example.com?foo=bar")
}
```

###



### 获取参数

#### 获取路径参数

可以通过匹配的方式，获取路径上的参数

- `:` 只能匹配1个
- `* `可以匹配任意个数

方式一,使用  `:` 只匹配1个参数

例子：

```go
// 此规则能够匹配/user/xxx这种格式，但不能匹配/user/ 或 /user这种格式
router.GET("/user/:name", func(c *gin.Context) {
    name := c.Param("name")
    c.String(http.StatusOK, "Hello %s", name)
})
```

这里name会作为参数，例如访问`http://localhost:8080/user/aaa`，name便等于aaa

方式二,使用 `* `可以匹配任意个数

```go
router.GET("/user/:name/*action", func(c *gin.Context) {
    name := c.Param("name")
    action := c.Param("action")
    message := name + " is " + action
    c.String(http.StatusOK, message)
})
```

规则：

```
/user/:name/*action
```

此规则既能匹配 /user/aaa/ 格式，也能匹配 /user/aaa/other1/other2 这种格式

- 访问：localhost:8080/user/aaa/

- 访问：localhost:8080/user/aaa/other1/

- 访问：localhost:8080/user/aaa/other1/other2/

注意`* `只能在最后用

#### 获取URL参数

- 获取URL参数可以通过DefaultQuery()或Query()方法获取
- 若参数不存在，DefaultQuery()返回默认值，Query()返回空串

```go
r.GET("/user", func(c *gin.Context) {
    //指定默认值
    name := c.DefaultQuery("name", "normal")
    //获取具体值
    age := c.Query("age")
    c.String(http.StatusOK, fmt.Sprintf("hello %s, your age is %s", name, age))
})
```

#### 获取Post方法

```go
r.POST("/form", func(c *gin.Context) {
    // 设置默认值
    types := c.DefaultPostForm("type", "post")
    username := c.PostForm("username")
    password := c.PostForm("password")
  
    // 还可以使用Query实现 Get + Post的结合
    name := c.Query("name")
    c.JSON(200, gin.H{
        "username": username,
        "password": password,
        "types":    types,
        "name":  name,
    })
})
```



#### URL参数和Post参数的获取案例

下面是 URL参数和Post参数的获取案例

```go
func main() {
 router := gin.Default()
 router.GET("/welcome", welcome)
 router.POST("/login", login)
 router.POST("/post", getPost)
 _ = router.Run()
}
 
// 获取GET传参
func welcome(c *gin.Context) {
 firstName := c.DefaultQuery("firstname", "unknown")
 lastName := c.DefaultQuery("lastname", "unknown")
 c.JSON(http.StatusOK, gin.H{
  "first_name": firstName,
  "last_name":  lastName,
 })
}
 
// 获取POST传参
func login(c *gin.Context) {
 username := c.DefaultPostForm("username", "test")
 password := c.DefaultPostForm("password", "test")
 c.JSON(http.StatusOK, gin.H{
  "username": username,
  "password": password,
 })
}
 
// 混合获取参数
func getPost(c *gin.Context) {
    // 获取GET参数
 id := c.Query("id")
 page := c.DefaultQuery("page", "0")
    // 获取POST参数
 name := c.PostForm("name")
 message := c.DefaultPostForm("message", "")
 c.JSON(http.StatusOK, gin.H{
  "id":      id,
  "page":    page,
  "name":    name,
  "message": message,
 })
}
```



#### bing绑定器

gin中的bind可以很方便的将 前端传递 来的数据与 `结构体` 进行 `参数绑定` ，以及参数校验

##### 参数绑定

在使用这个功能的时候，需要给结构体加上Tag `json` `form` `uri` `xml` `yaml`

##### Must Bind

不用，校验失败会改状态码

##### ShouldBind

可以绑定json，query，param，yaml，xml

如果校验不通过会返回错误

###### ShouldBindJSON

```go
package main

import "github.com/gin-gonic/gin"

type UserInfo struct {
  Name string `json:"name"`
  Age  int    `json:"age"`
  Sex  string `json:"sex"`
}

func main() {
  router := gin.Default()
  router.POST("/", func(c *gin.Context) {

    var userInfo UserInfo
    err := c.ShouldBindJSON(&userInfo)
    if err != nil {
      c.JSON(200, gin.H{"msg": "你错了"})
      return
    }
    c.JSON(200, userInfo)

  })
  router.Run(":80")
}

```

###### ShouldBindQuery

绑定查询参数，tag对应为form

```go
// ?name=枫枫&age=21&sex=男
package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

type UserInfo struct {
  Name string `json:"name" form:"name"`
  Age  int    `json:"age" form:"age"`
  Sex  string `json:"sex" form:"sex"`
}

func main() {
  router := gin.Default()

  router.POST("/query", func(c *gin.Context) {

    var userInfo UserInfo
    err := c.ShouldBindQuery(&userInfo)
    if err != nil {
      fmt.Println(err)
      c.JSON(200, gin.H{"msg": "你错了"})
      return
    }
    c.JSON(200, userInfo)

  })
  router.Run(":80")
}

```

###### ShouldBindUri

绑定动态参数，tag对应为uri

```go
// /uri/fengfeng/21/男

package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

type UserInfo struct {
  Name string `json:"name" form:"name" uri:"name"`
  Age  int    `json:"age" form:"age" uri:"age"`
  Sex  string `json:"sex" form:"sex" uri:"sex"`
}

func main() {
  router := gin.Default()

  router.POST("/uri/:name/:age/:sex", func(c *gin.Context) {

    var userInfo UserInfo
    err := c.ShouldBindUri(&userInfo)
    if err != nil {
      fmt.Println(err)
      c.JSON(200, gin.H{"msg": "你错了"})
      return
    }
    c.JSON(200, userInfo)

  })

  router.Run(":80")
}

```

###### ShouldBind

form-data的参数也用这个，tag用form；默认的tag就是form

**绑定form-data、x-www-form-urlencode**

```go
package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

type UserInfo struct {
  Name string `form:"name"`
  Age  int    `form:"age"`
  Sex  string `form:"sex"`
}

func main() {
  router := gin.Default()
  
  router.POST("/form", func(c *gin.Context) {
    var userInfo UserInfo
    err := c.ShouldBind(&userInfo)
    if err != nil {
      fmt.Println(err)
      c.JSON(200, gin.H{"msg": "你错了"})
      return
    }
    c.JSON(200, userInfo)
  })

  router.Run(":80")
}

```



### 请求头参数

`GetHeader`，可以大小写不分，且返回切片中的第一个数据

```go
router.GET("/", func(c *gin.Context) {
  // 首字母大小写不区分  单词与单词之间用 - 连接
  // 用于获取一个请求头
  fmt.Println(c.GetHeader("User-Agent"))
  //fmt.Println(c.GetHeader("user-agent"))
  //fmt.Println(c.GetHeader("user-Agent"))
  //fmt.Println(c.GetHeader("user-AGent"))

  // Header 是一个普通的 map[string][]string
  fmt.Println(c.Request.Header)
  // 如果是使用 Get方法或者是 .GetHeader,那么可以不用区分大小写，并且返回第一个value
  fmt.Println(c.Request.Header.Get("User-Agent"))
  fmt.Println(c.Request.Header["User-Agent"])
  // 如果是用map的取值方式，请注意大小写问题
  fmt.Println(c.Request.Header["user-agent"])

  // 自定义的请求头，用Get方法也是免大小写
  fmt.Println(c.Request.Header.Get("Token"))
  fmt.Println(c.Request.Header.Get("token"))
  c.JSON(200, gin.H{"msg": "成功"})
})

```

### 响应头参数

```go
// 设置响应头
router.GET("/res", func(c *gin.Context) {
  c.Header("Token", "jhgeu%hsg845jUIF83jh")
  c.Header("Content-Type", "application/text; charset=utf-8")
  c.JSON(0, gin.H{"data": "看看响应头"})
})

```







### 文件上传/下载

#### 单个文件上传

在使用 Gin 框架时，可以使用 `c.FormFile()` 方法来获取上传文件。这个方法会返回一个 `*multipart.FileHeader` 对象，它包含了上传文件的信息，比如文件名、文件大小、文件类型等。你可以通过这个对象获取文件的内容，并进行相应的处理。

以下是一个简单的示例代码，演示了如何使用 `c.FormFile()` 方法来上传文件：

```go
func uploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.String(http.StatusBadRequest, fmt.Sprintf("上传文件失败: %s", err.Error()))
        return
    }

    // 文件对象/文件路径，注意要从项目根路径开始写
    err = c.SaveUploadedFile(file, file.Filename) 

    if err != nil {
        c.String(http.StatusInternalServerError, fmt.Sprintf("保存文件失败: %s", err.Error()))
        return
    }

    c.String(http.StatusOK, fmt.Sprintf("文件 %s 上传成功", file.Filename))
}
```

在上面的代码中，我们首先使用 `c.FormFile()` 方法获取上传文件。如果获取失败，我们会返回一个错误信息。如果获取成功，我们就可以使用 `c.SaveUploadedFile()` 方法将文件保存到本地。最后，我们返回一个成功上传的信息。

在使用 `c.FormFile()` 方法时，需要注意的是，参数名应该与 HTML 表单中的文件上传控件的 `name` 属性相同。在上面的示例中，我们假设上传文件控件的 `name` 属性为 `"file"`。

上传文件的时候，也可以通过限制大小，参考代码如下：

```go
r := gin.Default()
// 给表单限制上传大小 (默认 32 MiB)
r.MaxMultipartMemory = 8 << 20 // 8 MiB
r.POST("/upload", func(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.String(500, "上传文件出错")
    }

    // 上传到指定路径
    c.SaveUploadedFile(file, "C:/desktop/"+file.Filename)
    c.String(http.StatusOK, "fileName:", file.Filename)
})
```

#### 多个文件上传

```go
func main() {
  router := gin.Default()
  // 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
  router.MaxMultipartMemory = 8 << 20 // 8 MiB
  router.POST("/upload", func(c *gin.Context) {
    // Multipart form
    form, _ := c.MultipartForm()
    files := form.File["upload[]"]  // 注意这里名字不要对不上了

    for _, file := range files {
      log.Println(file.Filename)
      // 上传文件至指定目录
      c.SaveUploadedFile(file, "./"+file.Filename)
    }
    c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
  })
  router.Run(":8080")
}

```



#### 文件下载

有些响应，比如图片，浏览器就会显示这个图片，而不是下载，所以我们需要使浏览器唤起下载行为

```go
c.Header("Content-Type", "application/octet-stream")  // 表示是文件流，唤起浏览器下载，一般设置了这个，就要设置文件名
c.Header("Content-Disposition", "attachment; filename="+"牛逼.png") // 用来指定下载下来的文件名
c.Header("Content-Transfer-Encoding", "binary") // 表示传输过程中的编码形式，乱码问题可能就是因为它
c.File("uploads/12.png")

```

注意，文件下载，浏览器可能会有缓存，这个要注意一下；解决办法就是加查询参数





### 表单验证

GIN提供了两种方法来进行表单验证: Must Bind / Should Bind。"Must Bind" 和 "Should Bind" 是在编写程序时常用的两个概念。

- "Must Bind" 意味着必须绑定某个变量或者参数，否则程序将无法正常运行。这种情况通常出现在必需的配置项或者必需的输入参数未被正确设置的情况下。
- "Should Bind" 则表示建议绑定某个变量或者参数，但是如果未绑定也不会影响程序的正常运行。这种情况通常出现在可选的配置项或者可选的输入参数未被设置的情况下。

简而言之，"Must Bind" 表示必须绑定，否则程序无法正常运行；"Should Bind" 表示建议绑定，但不是必须的。

在 Gin 框架中进行表单验证，可以使用 Gin 提供的 `binding` 包和 `validator` 包来实现。有3步骤：

1. 导入 `binding` 和 `validator` 包

2. 定义表单结构体，并使用 `binding:""` 和 `validate:""` 标记字段

3. 在路由处理函数中使用 `ShouldBindWith` 方法解析请求参数，并使用 `validator` 包的 `ValidateStruct` 方法进行验证

具体步骤如下：

1. 导入 `binding` 和 `validator` 包：

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "gopkg.in/go-playground/validator.v9"
)
```

2. 定义表单结构体，并使用 `binding:"required"` 和 `validate:""` 标记字段：

```go
type LoginForm struct {
    Username string `json:"username" binding:"required" validate:"required"`
    Password string `json:"password" binding:"required" validate:"required"`
}
```

其中：

- `binding:"required"` 表示该字段在请求中必须存在，否则会返回 400 错误；
- `validate:"required"` 表示该字段必须有值，否则会返回 422 错误。



3. 在路由处理函数中使用 `ShouldBindWith` 方法解析请求参数，并使用 `validator` 包的 `ValidateStruct` 方法进行验证：

```go
func Login(c *gin.Context) {
    var form LoginForm
    if err := c.ShouldBindWith(&form, binding.JSON); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    validate := validator.New()
    if err := validate.Struct(form); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
        return
    }

    // TODO: 处理登录逻辑
}
```

其中，有两个要点：

- `ShouldBindWith` 方法会根据请求头中的 `Content-Type` 自动解析请求参数，第二个参数指定了解析器类型，这里使用了 JSON 解析器；
- `ValidateStruct` 方法会根据表单结构体的标记进行验证，如果有错误则返回错误信息。

这样就可以在 Gin 框架中进行表单验证了。

接收表单请求，获取用户名和密码:

```go
type LoginForm struct {
    Username string `json:"username" binding:"required" validate:"required"`
    Password string `json:"password" binding:"required" validate:"required"`
}

 
func main() {
 router := gin.Default()
 router.POST("/login", func(c *gin.Context) {
  var loginForm LoginForm
  if err := c.ShouldBind(&loginForm); err != nil {
   fmt.Println(err.Error())
   c.JSON(http.StatusBadRequest, gin.H{
    "error": err.Error(),
   })
   return
  }
  c.JSON(http.StatusOK, gin.H{
   "msg": "login",
  })
 })
 _ = router.Run()
}
```

使用POST请求发送JSON数据

```
http://localhost:8080/login
```

**验证失败**

请求信息:

```
{
    "password":"123"
}
```

请求的结果: error



**验证成功**

请求信息:

```
{
    "username":"David",
    "password":"123"
}
```

请求的结果: msg



**设置校验的规则**

上面的案例，设置了 validate:"required" 。如果`required`字段没有收到，错误日志会告知：

> {
>
> "error": "Key: 'LoginForm.Username' Error:Field validation for 'Username' failed on the 'required' tag"
>
> }

除此之外，还可以有很多的校验。比如，通过tag设置范围校验，例如

> ```go
> binding:"required,gt=10"  // 代表该值需要大于10
> time_format:"2006-01-02" time_utc:"1"   // 时间格式 校验
> ```

此外，还允许**自定义校验方式**

### content-type绑定

在 Gin 框架中，可以通过 `ShouldBindJSON()` 方法将请求体中的 Content-Type 是 application/json的数据与指定的结构体进行绑定。下面是一个使用 绑定请求体的示例：

```go
package main

import (
    "github.com/gin-gonic/gin"
)

type Login struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func main() {
    r := gin.Default()

    r.POST("/login", func(c *gin.Context) {
        var login Login
        // 将request的body中的数据，按照json格式解析到结构体
        if err := c.ShouldBindJSON(&login); err != nil {
        // 如果发送的不是json格式，那么输出：  "error": "invalid character '-' in numeric literal"
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }

        // ...
    })

    r.Run()
}
```

在上面的示例中，我们定义了一个 `Login` 结构体，用于存储登录请求中的用户名和密码。

在路由处理函数中，我们使用 `ShouldBindJSON()` 方法将请求体中的 JSON 数据与 `Login` 结构体进行绑定。

如果绑定失败，我们将返回一个 400 错误响应，否则我们将继续处理登录逻辑。

需要注意的是，这里使用了 `binding:"required"` 标签来指定 `Username` 和 `Password` 字段必须存在。如果请求体中缺少这些字段，绑定将失败，并返回一个错误响应。

除了使用ShouldBindJSON，也可以使用`Bind`方法，参考代码如下：

```go
r.POST("/loginJSON", func(c *gin.Context) {
    // 声明接收的变量
    var login Login

    // 默认绑定form格式
    if err := c.Bind(&login); err != nil {
        // 根据请求头中content-type自动推断
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 输出结果
    c.JSON(http.StatusOK, gin.H{
        "status":   "200",
        "user":     login.User,
        "password": login.Password,
    })
})
```

`Bind` 和 `ShouldBindJSON` 都是 Gin 框架中用于将请求体中的数据绑定到结构体中的方法。

它们的区别在于，`Bind` 方法会根据请求头中的 Content-Type 自动选择绑定方法，而 `ShouldBindJSON` 方法则只会绑定 JSON 格式的请求体。

举个例子，如果请求头中的 Content-Type 是 application/json，那么 `Bind` 方法和 `ShouldBindJSON` 方法都会将请求体中的 JSON 数据绑定到结构体中。但如果 Content-Type 是 application/xml，那么 `Bind` 方法会选择绑定 XML 格式的请求体，而 `ShouldBindJSON` 方法则会返回错误，因为它只能绑定 JSON 格式的请求体。

因此，如果你确定请求体中的数据是 JSON 格式，可以直接使用 `ShouldBindJSON` 方法，否则建议使用 `Bind` 方法。

在进行绑定时，可以使用 `Content-Type` 请求头来指定请求体的格式。



### 异步执行

在 Gin 中异步执行可以使用 Go 语言的协程（goroutine）来实现。在处理请求的处理函数中，可以使用 `go` 关键字来启动一个新的协程，使得处理函数可以立即返回，并且新的协程可以在后台继续执行。

举个例子，如果我们需要在处理函数中执行一个比较耗时的操作，可以这样写：

```go
func handleRequest(c *gin.Context) {
    // 需要搞一个副本
    copyContext := c.Copy()
    
    // 启动一个新的协程来执行耗时操作
    go func() {
        // 执行耗时操作
        time.Sleep(5 * time.Second)
        // 操作完成后，可以通过 c.Writer 写入响应数据
        copyContext.Writer.WriteString("耗时操作完成")
    }()

}
```

在这个例子中，我们使用了匿名函数来启动一个新的协程，该协程会执行一个耗时操作，然后在操作完成后通过 `c.Writer` 写入响应数据。

需要注意的是，在协程中访问 Gin 的上下文对象 `c` 时，需要使用闭包，以避免竞态条件。此外，还需要注意协程的数量，避免过多的协程导致系统资源耗尽。

### 会话控制

- cookie相关
- session相关
- token相关

#### cookie相关

- 可以通过 `GetCookie` 方法获取客户端请求中携带的 cookie
- 可以通过 `SetCookie` 方法设置 cookie

在 Gin 中，可以通过 `SetCookie` 方法设置 cookie，例如：

```go
func main() {
    router := gin.Default()

    router.GET("/set-cookie", func(c *gin.Context) {
        c.SetCookie("username", "johndoe", 3600, "/", "localhost", false, true)
        c.String(http.StatusOK, "Cookie has been set")
    })

    router.Run(":8080")
}
```

在上面的例子中，我们使用 `SetCookie` 方法设置了一个名为 "username" 的 cookie，它的值为 "johndoe"，过期时间为 3600 秒，路径为 "/"，域名为 "localhost"，不启用安全标志，启用 HTTPOnly 标志。在客户端可以通过 `document.cookie` 属性读取该 cookie。

另外，可以通过 `GetCookie` 方法获取客户端请求中携带的 cookie，例如：

```
func main() {
    router := gin.Default()

    router.GET("/get-cookie", func(c *gin.Context) {
        username, err := c.Cookie("username")
        if err != nil {
            c.String(http.StatusBadRequest, "Cookie not found")
        } else {
            c.String(http.StatusOK, "Hello "+username)
        }
    })

    router.Run(":8080")
}
```

在上面的例子中，我们使用 `Cookie` 方法获取客户端请求中名为 "username" 的 cookie 的值，并将其作为字符串拼接到响应中。如果客户端请求中没有携带该 cookie，则返回 "Cookie not found"。

#### session相关

Session 是一种在客户端和服务器之间保存状态的机制，它可以用来存储用户的登录信息、购物车信息等。Gin 提供了一个中间件 gin-contrib/sessions，它可以帮助我们在 Gin 中使用 Session。

要使用 Gin Session，我们需要先安装 gin-contrib/sessions 包。可以使用以下命令进行安装：

```go
go get github.com/gin-contrib/sessions
```

安装完成后，我们需要在代码中引入 gin-contrib/sessions 包，并创建一个 Session 存储引擎。Gin 支持多种 Session 存储引擎，包括内存存储、Cookie 存储、Redis 存储等。以下是一个使用 Cookie 存储的示例：

```go
import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 设置 Session 中间件
    store := cookie.NewStore([]byte("secret"))
    r.Use(sessions.Sessions("mysession", store))

    // 设置路由
    r.GET("/set", func(c *gin.Context) {
        session := sessions.Default(c)
        session.Set("username", "johndoe")
        session.Save()
        c.JSON(200, gin.H{"message": "Session saved"})
    })

    r.GET("/get", func(c *gin.Context) {
        session := sessions.Default(c)
        username := session.Get("username")
        c.JSON(200, gin.H{"username": username})
    })

    r.Run(":8080")
}
```

在上面的示例中，我们首先创建了一个 Cookie 存储引擎，并将其作为 Session 中间件添加到 Gin 中。然后，我们定义了两个路由，一个用于设置 Session，另一个用于获取 Session。

在设置 Session 的路由中，我们使用 sessions.Default(c) 获取当前请求的 Session 对象，并使用 session.Set() 方法设置一个键值对。

在获取 Session 的路由中，我们同样使用 sessions.Default(c) 获取当前请求的 Session 对象，并使用 session.Get() 方法获取之前设置的键值对。

需要注意的是，Session 中间件需要在路由之前添加，这样才能在路由中使用 Session。

另外，Session 存储引擎中的 secret 参数应该是一个随机字符串，用于加密 Session 数据。

#### token相关

通常为了分布式和安全性，我们会采取更好的方式，比如使用Token 认证，来实现跨域访问，避免 CSRF 攻击，还能在多个服务间共享。

Token 是一种用于身份验证和授权的机制，通常用于保护 Web 应用程序中的敏感资源。

在 Gin 中，可以使用 JWT（JSON Web Token）作为 Token 机制。

使用 Gin 和 JWT 实现 Token 鉴权的步骤如下：

1. 在用户登录成功后，生成一个 JWT Token，并将其返回给客户端。
2. 客户端在每次请求中将 Token 作为请求头发送给服务器。
3. 服务器在接收到请求后，解析 Token，验证其有效性和正确性。
4. 如果 Token 有效，则允许用户访问相应的资源；否则，返回错误信息。

在 Gin 中，可以使用第三方库如 jwt-go 来实现 JWT 的生成和解析。具体实现方式可以参考以下代码：

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
)

// 生成 JWT Token
func generateToken(userId int64) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["userId"] = userId
    tokenString, err := token.SignedString([]byte("secret"))
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

// 鉴权中间件
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return []byte("secret"), nil
        })
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            userId := int64(claims["userId"].(float64))
            c.Set("userId", userId)
            c.Next()
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }
    }
}

// 示例路由
func main() {
    r := gin.Default()

    // 登录路由
    r.POST("/login", func(c *gin.Context) {
        // 模拟登录成功
        userId := int64(123)
        token, err := generateToken(userId)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"token": token})
    })

    // 需要鉴权的路由
    r.GET("/protected", authMiddleware(), func(c *gin.Context) {
        userId := c.MustGet("userId").(int64)
        c.JSON(http.StatusOK, gin.H{"userId": userId})
    })

    r.Run(":8080")
}
```

在上面的示例中，我们定义了一个生成 JWT Token 的函数 generateToken 和一个鉴权中间件 authMiddleware。

在登录成功后，我们将生成的 Token 返回给客户端。

在需要鉴权的路由中，我们使用 authMiddleware 作为中间件来验证 Token 的有效性:

- 如果 Token 有效，则将用户的 userId 存储到上下文中
- 否则返回错误信息。



### 中间件

gin`中间件`，类似spring  mvc 的拦截器 、过滤器。

gin`中间件`作用就是在处理具体的route请求时，提前做一些业务，还可以在业务执行完后执行一些操作。比如身份校验、日志打印等操作。

中间件分为：**全局中间件** 和 **路由中间件**，区别在于前者会作用于所有路由。

> 其实使用`router := gin.Default()`定义route时，默认带了`Logger()`和`Recovery()`。

看看 gin.Default() 源码就了解了：

```go
// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
 debugPrintWARNINGDefault()
 engine := New()
 engine.Use(Logger(), Recovery())
 return engine
}
```

#### 默认Gin中间件

Gin本身也提供了一些中间件给我们使用：

```go
func BasicAuth(accounts Accounts) HandlerFunc // 身份认证
func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc
func Bind(val interface{}) HandlerFunc //拦截请求参数并进行绑定
func ErrorLogger() HandlerFunc       //错误日志处理
func ErrorLoggerT(typ ErrorType) HandlerFunc //自定义类型的错误日志处理
func Logger() HandlerFunc //日志记录
func LoggerWithConfig(conf LoggerConfig) HandlerFunc
func LoggerWithFormatter(f LogFormatter) HandlerFunc
func LoggerWithWriter(out io.Writer, notlogged ...string) HandlerFunc
func Recovery() HandlerFunc
func RecoveryWithWriter(out io.Writer) HandlerFunc
func WrapF(f http.HandlerFunc) HandlerFunc //将http.HandlerFunc包装成中间件
func WrapH(h http.Handler) HandlerFunc //将http.Handler包装成中间件
```

#### 自定义Gin中间件

自定义中间件的方式很简单，我们只需要实现一个函数，返回`gin.HandlerFunc`类型的参数即可：

```go
// HandlerFunc 本质就是一个函数，入参为 *gin.Context
// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context)
```

##### 权限验证

自定义Gin中间件 用于检查token ,示例代码：

```go
func TokenRequired() gin.HandlerFunc {
 return func(c *gin.Context) {
  var token string
  for k, v := range c.Request.Header {
   if k == "x-token" {
    token = v[0]
   }
   fmt.Println(k, v, token)
  }
  if token != "test" {
   c.JSON(http.StatusOK, gin.H{
    "msg": "login failed",
   })
   c.Abort()
  }
  c.Next()
 }
}
 
func main() {
 router := gin.Default()
 router.Use(TokenRequired())
 router.GET("/ping", func(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
   "message": "pong",
  })
 })
 _ = router.Run()
}
```

##### 耗时统计

统计每一个视图函数的执行时间

```go
解释

func TimeMiddleware(c *gin.Context) {
  startTime := time.Now()
  c.Next()
  since := time.Since(startTime)
  // 获取当前请求所对应的函数
  f := c.HandlerName()
  fmt.Printf("函数 %s 耗时 %d\n", f, since)
}
CopyErrorOK!
```

##### BasicAuth

通过basicAuth可以快速实现http基础认证，其优势在于简单，更多的主要用途可以选择Oauth

```go
package main
 
import "github.com/gin-gonic/gin"
 
//basicAuth是简答的验证功能
 
//模拟存储私人信息到这里
var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}
 
//BasicAuth好像不可以通过额外路由
// func test_BasicAuth(c *gin.Context) {
// 	//BasicAuth返回一个基本http中间件，接受一个map[string]string作为参数
// 	gin.BasicAuth(gin.Accounts{
// 		"foo":    "nihao",
// 		"austin": "1234",
// 		"lena":   "hello2",
// 		"manu":   "4321",
// 	})
// }
 
func test_BasicAuthsecrets(c *gin.Context) {
	//获取用户名，它是由BaiscAuth中间件设置的
	user := c.MustGet(gin.AuthUserKey).(string)
	if secret, ok := secrets[user]; ok {
		c.JSON(200, gin.H{
			"user":   user,
			"secret": secret,
		})
	}
}
 
func main() {
	e := gin.Default()
 
	//BasicAuth只能单独在这里进行注册
	rg := e.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "nihao", //这里设置的为用户名和密码
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))
	{
		rg.GET("/secrets", test_BasicAuthsecrets)
	}
 
	e.Run()
 
}
```



#### 中间件控制的方法

在 Gin 框架中，中间件的顺序非常重要，因为它们按照添加的顺序依次执行。如果您希望中间件以特定的顺序执行，可以使用 Gin 框架提供的 Use() 方法来添加中间件。例如，如果您希望在日志中记录请求之前先执行身份验证中间件，则应该先添加身份验证中间件，然后再添加日志中间件，如下所示：

```go
router := gin.Default()

// 添加身份验证中间件
router.Use(authMiddleware)

// 添加日志中间件
router.Use(loggerMiddleware)

// 添加路由处理函数
router.GET("/hello", func(c *gin.Context) {
    c.String(http.StatusOK, "Hello, World!")
})

router.Run(":8080")
```

在这个例子中，authMiddleware 会在 loggerMiddleware 之前执行，因为它先被添加到 Gin 的中间件处理链中。

另外，gin提供了两个函数`Abort()`和`Next()`，二者区别在于：

> 1. next()函数会跳过当前中间件中next()后的逻辑，当**下一个中间件执行完成后**再执行剩余的逻辑
> 2. abort()函数执行终止当前中间件以后的中间件执行，**但是会执行当前中间件的后续逻辑**

举例子更好理解：我们注册中间件顺序为`m1`、`m2`、`m3`，如果采用`next()` 执行顺序就是

> 1. `m1的next()前面`、`m2的next()前面`、`m3的next()前面`、
> 2. `业务逻辑`
> 3. `m3的next()后续`、`m2的next()后续`、`m1的next()后续`。

那如果`m2`中间调用了`Abort()`，则`m3`和`业务逻辑`不会执行，只会执行`m2的next()后续`、`m1的next()后续`。



#### 局部中间件

Gin 框架中有两种类型的中间件：

- 全局中间件
- 局部中间件。

全局中间件是在 Gin 实例创建时添加的，它们将应用于所有的路由请求。

全局中间件可以通过 `Use()` 方法添加到 Gin 实例中，例如：

```go
router := gin.Default()
router.Use(Logger())
```

在上面的代码中，`Logger()` 函数是一个全局中间件，它将在所有请求之前打印请求的信息。

局部中间件只会应用于某些路由请求。

- 可以在路由设置的时候，设置局部中间件
- 可以使用 `Group` 方法创建一个路由组，然后在这个路由组上设置局部中间件

在路由设置的时候，设置局部中间件：

```go
router := gin.Default()

// 添加全局中间件
router.Use(Logger())

// 通过 Handle()添加局部中间件
router.GET("/posts", Auth(), GetPosts)
```

在上面的代码中，`Auth()` 函数是一个局部中间件，它只会在 `/posts` 路由请求中应用。

在 Gin 中，可以使用 `Group` 方法创建一个路由组，然后在这个路由组上设置局部中间件。

例如，以下代码创建了一个路由组 `/api`，并在这个路由组上设置了一个局部中间件 `authMiddleware`：

```go
func authMiddleware(c *gin.Context) {
    // 检查用户是否已经登录
    // 如果用户已经登录，则继续处理请求
    // 如果用户未登录，则返回 401 Unauthorized 错误
}

func main() {
    r := gin.Default()

    api := r.Group("/api")
    api.Use(authMiddleware)

    api.GET("/users", func(c *gin.Context) {
        // 处理 GET /api/users 请求
    })

    api.POST("/users", func(c *gin.Context) {
        // 处理 POST /api/users 请求
    })

    r.Run()
}
```

在上面的代码中，`authMiddleware` 函数是一个中间件函数，用于检查用户是否已经登录。`api` 是一个路由组，表示所有以 `/api` 开头的请求都会进入这个路由组。`api.Use(authMiddleware)` 表示在这个路由组上使用 `authMiddleware` 中间件。

这样，所有以 `/api` 开头的请求都会先进入 `authMiddleware` 中间件函数进行身份验证，如果身份验证通过，则继续处理请求，否则返回 401 Unauthorized 错误。