

# swagger

使用 Swagger 开源和专业工具集简化用户、团队和企业的 API 开发, 并且帮助开发者完善接口设计，接口开发、接口文档、接口测试、API模拟和虚拟化、以及接口治理与接口监控。

官网地址: https://swagger.io/

文档地址: https://swagger.io/docs/specification/2-0/basic-structure/

接口规范: https://swagger.io/resources/open-api/

项目地址:
- https://github.com/swaggo/swag
- https://github.com/swaggo/gin-swagger

官方文档: https://github.com/swaggo/swag/blob/master/README_zh-CN.md


## 安装

### go安装

```sh
// 下载swagger命令行工具
go install github.com/swaggo/swag/cmd/swag@latest

go env | grep -i gopath // 查看${gopath}/bin目录下是否存在swag
go env | grep -i goroot // 查看${goroot}/bin目录下是否存在swag,该目录是go执行命令目录，如果不存在swag，则需要将swag移动到该目录下


// 下载swagger源码依赖
go get github.com/swaggo/swag

// 下载swagger的静态文件库，html，css，js之类的，都被嵌到了go代码中
go get github.com/swaggo/files@latest

// 下载swagger的gin适配库
go get github.com/swaggo/gin-swagger@latest

```

### vscode安装插件

插件：Swagger Viewer


## 初始化

```sh
// 格式化代码
swag fmt

// 初始化swagger
swag init -d main.go

```

## 配置

main.go
```go

package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    // 匿名导入生成的接口文档包
    _ "golearn/docs"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @BasePath  /api/v1
func main() {
    engine := gin.Default()
    // 注册swagger静态文件路由
    engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    engine.GET("/api/v1/ping", Ping)
    engine.Run(":80")
}

// Ping godoc
// @Summary      say hello world
// @Description  return hello world json format content
// @param       name query    string  true  "name"
// @Tags         system
// @Produce      json
// @Router       /ping [get]
func Ping(ctx *gin.Context) {
    ctx.JSON(200, gin.H{
       "message": fmt.Sprintf("Hello World!%s", ctx.Query("name")),
    })
}

```

访问：http://127.0.0.1:8081/swagger/index.html

## 参数

定义参数的格式为：`@Param <参数名> <数据类型> <参数类型> <是否必须> "<备注信息>"`

如：
```go
// @Param   uid      path     int    true        "uid" 
// @Param   user     body     model.User true    "user"
// @Param   id       path     int    true        "id"
// @Param   token    header   string true        "token"
```

支持的参数类型有: `query`,`path`,`header`,`body`,`formData`

数据类型有: `string (string)`,`number (float64,float32)`,`integer (int, uint, int32, uint32, int64, uint64)`,`boolean (bool)`,`array (array, slice)`,`user defined struct - object (struct)`,`file (multipart.FileHeader)`

## 返回值

定义接口响应的基本格式如下: `@Success <状态码> <数据类型> <备注信息>`

如：
```go
// @Success      200  {array}   model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError

```

比如一般我们会定义一个统一的响应体
```go

type JSONResult struct {
    Code    int          `json:"code" `
    Message string       `json:"message"`
    Data    interface{}  `json:"data"`
}
```
Data字段的类型是不确定的，在描述响应用例时，可以将其组合，如下
```go
// @Success      200  {object}  model.JSONResult{data=model.Account} "{data=model.Account} 对象"
// @Success      200  {object}  model.JSONResult{data=[]model.Account} "[]model.Account 切片对象"
// @Failure      400  {object}  model.JSONResult{data=httputil.HTTPError} "错误对象"
// @Failure      404  {object}  model.JSONResult{data=httputil.HTTPError} "错误对象"

```

## 模型

给结构体字段加注释会被被swagger扫描为模型字段注释

```go
package model

type Account struct {
	// account id
    ID   int    `json:"id" example:"1"`
    // username
    Name string `json:"name" example:"account name"`
}
```
其中example标签的值会被作为示例值在页面中展示，当然它还支持字段限制
```go
type Foo struct {
    Bar string `minLength:"4" maxLength:"16"`
    Baz int `minimum:"10" maximum:"20" default:"15"`
    Qux []string `enums:"foo,bar,baz"`
}

```

## 认证

在认证这块支持
- Basic Auth
- API Key
- OAuth2 app auth
- OAuth2 implicit auth
- OAuth2 password auth
- OAuth2 access code auth

假如接口认证用的是JWT，存放在header中的Authorization字段中，我们可以如下定义
```go
// @SecurityDefinitions.apikey ApiKeyAuth
// in: header
// name: Authorization
func main() {
    ...
}

// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /edit
func Edit(c *gin.Context) {
    ...
}

Basic Auth
// @securityDefinitions.basic BasicAuth
// @in header
// @name Auth
func main() {
    ...
}

/* 账号密码
// @securityDefinitions.basic	BasicAuth
// $ echo -n weiyigeek:123456 | base64
// d2VpeWlnZWVrOjEyMzQ1Ng==
*/

// @Security BasicAuth
// @Success 200 {string} string "ok"
// @Router /edit
func Edit(c *gin.Context) {
    ...
}