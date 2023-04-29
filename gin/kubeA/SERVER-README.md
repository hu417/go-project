
【项目】

1、项目初始化

```bash
mkdir -p kubea-demo && cd kubea-demo 
go mod init kubea-demo

```

2、项目骨架文件目录

```bash
mkdir -p {config,controller,dao,db,middle,model,service,utils}

kubea-demo/
|-- README.md
|-- config      # 定义全局配置，如监听地址、管理员账号等
|-- controller  # controller层，定义路由规则，及接口入参和响应
|-- dao         # 数据库操作，包含数据库的增删改查
|-- db          # 用于初始化数据库连接以及配置
|-- go.mod      # 定义项目的依赖包以及版本
|-- middle      # 中间件层，添加全局的逻辑处理，如跨域、jwt验证等
|-- model       # 定义数据库的表的字段
|-- service     # 服务层，处理接口的业务逻辑
|-- utils       # 工具目录，定义常用工具，如token解析，文件操作等
|-- main.go     # 项目的主入口 main函数

```

3、部分操作记录

```bash
# 下载依赖
go mod tidy
# 拷贝依赖到当前目录下
go mod vendor

# http各种状态码对应的函数
https://pkg.go.dev/net/http


# client-go是kubernetes官方提供的go语言的客户端库
# client-go地址: https://github.com/kubernetes/client-go
go get k8s.io/client-go@v0.23.15


# 日志依赖包-logger
go get github.com/wonderivan/logger@latest


```
