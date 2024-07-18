

# 微服务组件

微服务架构组件:

* 服务网关：确保服务提供者对客户端的透明，这一层可以进行反向路由、安全认证、灰度发布、日志监控等前置动作
* 服务发现：注册并维护远程服务及服务提供者的地址，供服务消费者发现和调用，为保证可用性，比如etcd，nacos，consul等
* 服务框架：用于实现微服务的 RPC 框架，包含服务接口描述及实现方案、向注册中心发布服务等功能，比如grpc，Thrift等
* 服务监控：对服务消费者与提供者之间的调用情况进行监控和数据展示，比如prometheus等
* 服务追踪：记录对每个请求的微服务调用完整链路，以便进行问题定位和故障分析，比如jeager，zipkin等
* 服务治理：服务治理就是通过一系列的手段来保证在各种意外情况下，服务调用仍然能够正常进行，这些手段包括熔断、隔离、限流、降级、负载均衡等。比如Sentinel，Istio等
* 基础设施：用以提供服务底层的基础数据服务，比如分布式消息队列、日志存储、数据库、缓存、文件服务器、搜索集群等。比如Kafka，Mysql，PostgreSQL，MongoDB，Redis，Minio，ElasticSearch等
* 分布式配置中心：统一配置，比如nacos，consul，apollo等
* 分布式事务：seata，dtm等
* 容器以及容器编排：docker，k8s等
* 定时任务



# 入门

## 环境准备

### go配置

```bash
go版本:
  go version // 1.19.3
  go env  // 
  go env -w GO111MODULE=on
  go env -w GOPROXY=https://goproxy.cn,direct
  # 注意要将%GOPATH%\bin设置到path环境变量中

IDE:
  vscode插件: goctl, vscode-proto3
```

### goctl安装

goctl（官方建议读go control）是go-zero微服务框架下的代码生成工具。使用 goctl 可显著提升开发效率，让开发人员将时间重点放在业务开发上，其功能有：

- api服务生成
- rpc服务生成
- model代码生成
- 模板管理

~~~go
方式一:
# github地址: https://github.com/zeromicro/go-zero/tree/master/tools/goctl

GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@v1.5.0

# 此时会在gopath的bin目录下生成goctl的执行进程（注意要将%GOPATH%\bin设置到path环境变量中）
goctl -v

方式二:
# Go 1.16 及以后版本
GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@latest
~~~

### protoc & protoc-gen-go安装

go-zero提供了便捷的安装方式：

~~~powershell
goctl env check -i -f --verbose
~~~



## API服务生成

### 命令行直接创建项目

```bash
项目创建:
  初始化API项目:
    mkdir -p zero-demo && cd  zero-demo
    // 初始化go.mod文件
    go mod init zero-demo
    // 快捷创建api服务
    goctl api new greet
    // 安装依赖
    go mod tidy
    // 复制依赖到vender目录
    go mod vendor
    // 运行 
    go run greet/greet.go -f greet/etc/greet-api.yaml
    // 访问
    http://localhost:8888/from/{name} 
    
    // 查看参数结构体
    # internal/types/types.go
    type Request struct {
      Name string `path:"name,options=you|me"` // tag是path绑定uri请求,且值为you/me; 可以自定义
    }

    type Response struct {
      Message string `json:"message"` // 返回的参数,及数据类型
    }

    // 查看路由
    # internal/handler/routes.go
    func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
      server.AddRoutes(
        []rest.Route{
          {
            Method:  http.MethodGet,
            Path:    "/from/:name",
            Handler: GreetHandler(serverCtx),
          },
        },
      )
    }

    // 修改返回信息
    # internal/logic/greetlogic.go
    func (l *GreetLogic) Greet(req *types.Request) (resp *types.Response, err error) {
      // todo: add your logic here and delete this line
      resp = &types.Response{
        Message: "hello go-zero!", // 对应的type.go响应结构体
      }
      return
    }

    // 修改服务配置
    # etc/greet-api.go
    Name: greet-api
    Host: localhost # 0.0.0.0 
    Port: 8888

# 服务访问
http://localhost:8888/from/me


```

### 目录结构

```bash
目录结构:
greet
├── etc                                 // 配置
│   └── greet-api.yaml                  // 配置文件
├── greet.api                           // 描述文件用于快速生成代码,通过定义好的内容生成对应的api系列go文件,也可作为文档查阅
├── greet.go                            // 入口文件
└── internal                            // 核心文件,主要操作文件夹，包括路由、业务等
    ├── config                          // 配置模块
    │   └── config.go                   // 配置解析映射结构体
    ├── handler                         // 控制层(路由)
    │   ├── greethandler.go             // 路由对应方法,API的控制类.
    │   └── routes.go                   // 路由文件,方法定义映射类,建议使用api文件生成
    ├── logic                           // 逻辑层
    │   └── greetlogic.go               // 业务逻辑类
    ├── svc                             // 上下文传输层
    │   └── servicecontext.go           // 类似于IOC容器,绑定主要操作依赖;上下文信息的传递内容定义，在主函数中有引用
    └── types                           // 类型字段定义
        └── types.go                    // 请求及响应结构体规则定义,建议使用api文件生成


```

### 使用api文件创建项目

```bash
创建项目
  mkdir -p user-api && cd user-api
  go mod init user-api
快速生成api文件:
  goctl api -o user.api
修改api文件:
syntax = "v1"

info (
	title: // TODO: add title
	desc: // TODO: add description
	author: "hu417"
	email: "hu729919300@163.com"
)

// 方式一
type request {
	// TODO: add members here and delete this comment
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
	Age      int    `json:"age"`
	Sex      string `json:"sex"`
}

type response {
	// TODO: add members here and delete this comment
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 方式二
type (
	// 请求的结构体
	LoginReq {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 响应的结构体
	LoginReply {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

// user-api服务里的metadata，比如url的前缀和group群组
@server(
	prefix: v1       // url的前缀
	group: user            // 群组,主要是在logic逻辑层再加个目录
)

// user-api服务的信息
service user-api {
	// 处理函数
	@handler GetUser // TODO: set handler name and delete this comment
	// 方法与URL路径,绑定type的相关结构体
	get /user(request) returns(response)
	
	// 处理函数
	@handler LoginUser // TODO: set handler name and delete this comment
	// 方法与URL路径
	post /login(LoginReq) returns(LoginReply)
}

根据user.api生成user模块:
  goctl api go -api user.api -dir user

下载依赖
  cd user
  go mod tidy
修改配置
  1、设置返回值
  # internal/logic/user/loginuserlogic.go
    func (l *LoginUserLogic) LoginUser(req *types.LoginReq) (resp *types.LoginReply, err error) {
        // todo: add your logic here and delete this line
        // 设置返回信息
        resp = &types.LoginReply{
            Code:    "200",
            Message: fmt.Sprintf("name=%s passwd=%s", req.Username, req.Password),
        }
        return
  }
  # internal/logic/user/loginuserlogic.go
    func (l *GetUserLogic) GetUser(req *types.Request) (resp *types.Response, err error) {
        // todo: add your logic here and delete this line
        // 设置返回信息
        resp = &types.Response{
            Code:    "200",
            Message: fmt.Sprintf("name=%s passwd=%s", req.Username, req.Passwd),
        }
        return
    }  
  
  2、设置服务配置参数
  # etc/user-api.yaml
    Name: user-api
    Host: 0.0.0.0
    Port: 8080 # 8888

测试请求

```

测试请求

```bash
1、get 127.0.0.1:8080/v1/user
  body:
      json格式
      {
        "username": "laowang",
        "passwd": "123456",
        "age": 19,
        "sex":"man"
      }
  
2、post 127.0.0.1:8080/v1/login
  body
	josn格式:
	{
        "username": "laowang",
        "password": "123456"
	}
```



## goctl常用命令

```bash
生成api模板
  goctl api -o user.api 

根据api文件生成go文件
  goctl api go -api search.api -dir=.

生成rpc模板
  goctl rpc template -o=user.proto

生成model
  goctl model mysql datasource -url="root:123456@tcp(10.0.0.91:3306)/test" -table="user" -dir . -c -style goZero

生成 Dockerfile
  goctl docker -go user.go

生成k8s yaml
  goctl kube deploy -name user-api -namespace blog -image user:v1 -o user.yaml -port 8888

```






## logx日志设置

### logx相关配置

```bash
type LogConf struct {
    ServiceName         string `json:",optional"`
    Mode                string `json:",default=console,options=[console,file,volume]"`
    Encoding            string `json:",default=json,options=[json,plain]"`
    TimeFormat          string `json:",optional"`
    Path                string `json:",default=logs"`
    Level               string `json:",default=info,options=[info,error,severe]"`
    Compress            bool   `json:",optional"`
    KeepDays            int    `json:",optional"`
    StackCooldownMillis int    `json:",default=100"`
}

ServiceName：设置服务名称，可选。在 volume 模式下，该名称用于生成日志文件。在 rest/zrpc 服务中，名称将被自动设置为 rest或zrpc 的名称。
Mode：输出日志的模式，默认是 console
    console 模式将日志写到 stdout/stderr
    file 模式将日志写到 Path 指定目录的文件中
    volume 模式在 docker 中使用，将日志写入挂载的卷中
Encoding: 指示如何对日志进行编码，默认是 json
    json模式以 json 格式写日志
    plain模式用纯文本写日志，并带有终端颜色显示
TimeFormat：自定义时间格式，可选。默认是 2006-01-02T15:04:05.000Z07:00
Path：设置日志路径，默认为 logs
Level: 用于过滤日志的日志级别。默认为 info
    info，所有日志都被写入
    error, info 的日志被丢弃
    severe, info 和 error 日志被丢弃，只有 severe 日志被写入
Compress: 是否压缩日志文件，只在 file 模式下工作
KeepDays：日志文件被保留多少天，在给定的天数之后，过期的文件将被自动删除。对 console 模式没有影响
StackCooldownMillis：多少毫秒后再次写入堆栈跟踪。用来避免堆栈跟踪日志过多

```

基本配置

```bash
# etc/user.yaml
Log:
  ServiceName: app
  # Mode: console       # 日志模式，[console,file,volume]
  Mode: file
  Path: logs
  Encoding: plain   # 输出格式，plain换行，json是一整行
  TimeFormat: "2006-01-02T 15:04:05.000Z07:00"  # 时间格式
  Level: info   # [debug,info,error,severe]
  Compress: true  # 启用压缩
  KeepDays: 7     # 保留天数            int    `json:",optional"`
  StackCooldownMillis: 100  # 多少毫秒后再次写入堆栈跟踪


```

日志输出

```bash
# internal/logic/user/getuserlogic.go

func (l *GetUserLogic) GetUser(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
    
    // 输出日志
	logx.Error("---------- GetUser ---------")

	return &types.Response{
		Code:    "1",
		Message: fmt.Sprintf("username=%v, password= %v", req.Username, req.Passwd),
	}, nil
}
```

禁用stat日志

```bash
# user.go

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 禁止stat日志
	logx.DisableStat()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

```



## mysql配置

### 创建MySQL

```bash
# 创建MySQL
docker run -p 3306:3306 --name mysql \
    -v /app/mysql/log:/var/log/mysql \
    -v /app/mysql/data:/var/lib/mysql \
    -v /app/mysql/conf:/etc/mysql \
    -e MYSQL_ROOT_PASSWORD=123456 \
    -d mysql:5.7


# 创建库
CREATE DATABASE `test` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

# 创建表
CREATE TABLE `user`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `gender` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

```

### 生成model

```bash
# 与internal
mkdir -p model && cd model

goctl model mysql datasource -url="root:123456@tcp(10.0.0.91:3306)/test" -table="user" -dir . -c -style goZero

go mod tidy

```


### 配置

#### 方式一

```bash
# etc/user.yaml
Mysql:
  Datasource: root:123456@tcp(10.0.0.91:3306)/test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

# 需要增加redis配置
CacheRedis:
  - Host: 10.0.0.91:6379
    Pass: 
    Type: node


# internal/config/config.go
import (
	"github.com/zeromicro/go-zero/core/stores/cache"  # 添加缓存依赖
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql struct {   # 定义mysql结构体
		Datasource string
	}
	CacheRedis cache.CacheConf  # 定义redis结构体
}



# internel/svc/serviceContext.go
import (
	"user-api/user/internal/config"
	"user-api/user/model"  // 添加model依赖

    "github.com/tal-tech/go-zero/core/stores/sqlx"  // 使用sqlx依赖

)

type ServiceContext struct {
	Config    config.Config
	Usermodel model.UserModel   // usermodel接口
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Usermodel: model.NewUserModel(sqlx.NewMysql(c.Mysql.Datasource), c.CacheRedis),  // 配置引用
	}
}


#  internal/logic/user/getuserlogic.go
import (
	"context"
	"fmt"

	"user-api/user/internal/svc"
	"user-api/user/internal/types"
	"user-api/user/model"   # 引入model包

	"github.com/zeromicro/go-zero/core/logx"
)

func (l *GetUserLogic) GetUser(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	logx.Error("---------- GetUser ---------")

	user := &model.User{
		Name:   req.Username,
		Gender: "man",
	}

	insert, err := l.svcCtx.Usermodel.Insert(l.ctx, user)
	if err != nil {
		panic(err)
	}

	id, _ := insert.LastInsertId()

	return &types.Response{
		Code:    "1",
		Message: fmt.Sprintf("id=%v, username=%v, password= %v", id, req.Username, req.Passwd),
	}, nil
}


```




#### 方式二

参考: https://blog.csdn.net/kemosisongge/article/details/128444916

```bash
# mysql连接参数
// Config 数据库配置信息
type Config struct {
	Username                  string // 用户名
	Password                  string // 密码
	Host                      string // 链接地址
	Port                      string // 端口
	Database                  string // 数据库名
	Charset                   string // 字符集
	ParseTime                 string // 是否解析时间
	Loc                       string // 时区默认 Local
	TablePrefix               string // 表前缀
	SingularTable             bool   // 是否使用单数表名 true 是 false 否
	PrepareStmt               bool   //  在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
	ConnMaxLifetime           int    // 设置了连接可复用的最大时间
	MaxIdleConn               int    // 连接池里面的连接最大存活时长(秒)
	ConnMaxIdleTime           int    //连接池里面的连接最大空闲时长(秒)
	MaxOpenConn               int    // 设置打开数据库连接的最大数量
	IgnoreRecordNotFoundError bool   // 忽略ErrRecordNotFound（记录未找到）错误
}


Mysql:
  Username: root
  Password: mysql
  Host: 127.0.0.1
  Port: "3308"
  Database: "yoolib"
  Charset: "utf8mb4"
  ParseTime: "false"
  Loc: "Local"
  TablePrefix: ""
  SingularTable: true
  PrepareStmt: true
  ConnMaxLifetime: 300
  MaxIdleConn: 30
  ConnMaxIdleTime: 200
  MaxOpenConn: 300
  IgnoreRecordNotFoundError: false



# etc/user.yaml
MySql:
  Type: mysql
  Host: 10.0.0.91
  Port: 3306
  DBname: test
  Username: root
  Password: 123456
  MaxOpenConn: 1000
  SSLMode: disable
  CacheTime: 5

```


## redis配置

```bash



```













# 图书查阅系统

## 预期实现目标
    用户登录 依靠现有学生系统数据进行登录
    图书检索 根据图书关键字搜索图书，查询图书剩余数量。

## 系统分析
### 服务拆分
- user
    - api 提供用户登录协议
    - rpc 供search服务访问用户数据
- search
    - api 提供图书查询协议


## 环境准备
### golang
```bash
1、go version
go version go1.19.3 windows/amd64

2、go module
查看GO111MODULE开启情况
$ go env GO111MODULE
on

开启GO111MODULE，如果已开启（即执行go env GO111MODULE结果为on）请跳过。
$ go env -w GO111MODULE="on"

设置GOPROXY
$ go env -w GOPROXY=https://goproxy.cn

设置GOMODCACHE
查看GOMODCACHE
$ go env GOMODCACHE

如果目录不为空或者/dev/null，请跳过。
go env -w GOMODCACHE=$GOPATH/pkg/mod

3、goctl
# github地址: https://github.com/zeromicro/go-zero/tree/master/tools/goctl

GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@v1.5.0

# 此时会在gopath的bin目录下生成goctl的执行进程（注意要将%GOPATH%\bin设置到path环境变量中）
goctl -v


4、protoc & protoc-gen-go
方式一
goctl env check -i -f --verbose



```

### etcd
```bash
# 部署
docker run -d --name etcd3 \
    -p 2379:2379 -p 2380:2380     \
    --env ALLOW_NONE_AUTHENTICATION=yes \
    --env ETCD_ADVERTISE_CLIENT_URLS=http://10.0.0.91:2379 \
    bitnami/etcd:3.5.6

# 测试
$ etcdctl --endpoints=0.0.0.0:2379 member list -w table
$ etcdctl --endpoints=0.0.0.0:2379 endpoint health -w table
```

### mysql
```bash
# 安装MySQL
mkdir -p /app/mysql/{conf,data,log}
chmod -R 777 /app/mysql/{conf,data,log}
vim /app/mysql/conf/my.cnf
...
[mysqld]
user=mysql
character-set-server=utf8
collation-server=utf8_general_ci
default_authentication_plugin=mysql_native_password
secure_file_priv=/var/lib/mysql
expire_logs_days=7
sql_mode=STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
max_connections=1000
 
[client]
default-character-set=utf8
 
[mysql]
default-character-set=utf8



docker run --restart=always \
    --privileged=true      \
    -v /app/mysql/data/:/var/lib/mysql \
    -v /app/mysql/log/:/var/log/mysql \
    -v /app/mysql/conf/my.cnf:/etc/my.cnf \
    -p 3306:3306 --name mysql \
    -e LANG="C.utf8" \
    -e MYSQL_ROOT_PASSWORD=123456 \
    -d mysql:8.0.32

# 容器字符集
locale -a 
echo "$LANG"

# 查看字符集
show variables like 'character_%';
show variables like 'collation_%';

# 创建库
CREATE DATABASE `book` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;


# 创建表
CREATE TABLE `test`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `gender` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

# 预设数据
INSERT INTO `test` (number,name,password,gender)values ('666','小明','123456','男');
```

### redis
```bash
mkdir redis/{conf,data,log}
touch redis/log/redis.log
chmod 777 -R redis/{conf,data,log}

# redis/conf/redis.conf
# 任意ip可访问
bind 0.0.0.0
# 自定义启动端口
port 6379
protected-mode no
# daemonize yes
loglevel notice
logfile "/var/log/redis.log"
# rdb或aof文件存储位置
dir /opt/data
save 900 1
save 300 10
save 60 10000
appendonly yes
appendfilename "appendonly.aof"
requirepass 123

docker run -d --name redis \
    --restart=always -p 6379:6379 \
    -v /app/redis/conf/redis.conf:/opt/redis.conf \
    -v /app/redis/data:/opt/data \
    -v /app/redis/log/redis.log:/var/log/redis.log \
    redis:6.2.7 redis-server /opt/redis.conf

```


## 目录拆分
目录拆分是指配合go-zero的最佳实践的目录拆分，这和微服务拆分有着关联，在团队内部最佳实践中， 我们按照业务横向拆分，将一个系统拆分成多个子系统，每个子系统应拥有独立的持久化存储，缓存系统。 如一个商城系统需要有用户系统(user)，商品管理系统(product)，订单系统(order)，购物车系统(cart)，结算中心系统(pay)，售后系统(afterSale)等组成。

### 系统结构分析
在上文提到的商城系统中，每个系统在对外（http）提供服务的同时，也会提供数据给其他子系统进行数据访问的接口（rpc），因此每个子系统可以拆分成一个服务，而且对外提供了两种访问该系统的方式：api和rpc，因此， 以上系统按照目录结构来拆分有如下结构:
```bash
.
├── afterSale
│   ├── api
│   └── rpc
├── cart
│   ├── api
│   └── rpc
├── order
│   ├── api
│   └── rpc
├── pay
│   ├── api
│   └── rpc
├── product
│   ├── api
│   └── rpc
└── user
    ├── api
    └── rpc

```


### rpc调用链建议
在设计系统时，尽量做到服务之间调用链是单向的，而非循环调用，例如：order服务调用了user服务，而user服务反过来也会调用order的服务， 当其中一个服务启动故障，就会相互影响，进入死循环，你order认为是user服务故障导致的，而user认为是order服务导致的，如果有大量服务存在相互调用链， 则需要考虑服务拆分是否合理。


### 常见服务类型的目录结构
在上述服务中，仅列举了api/rpc服务，除此之外，一个服务下还可能有其他更多服务类型，如rmq（消息处理系统），cron（定时任务系统），script（脚本）等， 因此一个服务下可能包含以下目录结构：
```bash
user
    ├── api //  http访问服务，业务需求实现
    ├── cronjob // 定时任务，定时数据更新业务
    ├── rmq // 消息处理系统：mq和dq，处理一些高并发和延时消息业务
    ├── rpc // rpc服务，给其他子系统提供基础数据访问
    └── script // 脚本，处理一些临时运营需求，临时数据修复

```


### 完整工程目录结构示例
```bash
mall // 工程名称
├── common // 通用库
│   ├── randx
│   └── stringx
├── go.mod
├── go.sum
└── service // 服务存放目录
    ├── afterSale
    │   ├── api
    │   └── model
    │   └── rpc
    ├── cart
    │   ├── api
    │   └── model
    │   └── rpc
    ├── order
    │   ├── api
    │   └── model
    │   └── rpc
    ├── pay
    │   ├── api
    │   └── model
    │   └── rpc
    ├── product
    │   ├── api
    │   └── model
    │   └── rpc
    └── user
        ├── api
        ├── cronjob
        ├── model
        ├── rmq
        ├── rpc
        └── script

```


## 项目初始化
```bash
mkdir -p book && cd book 
go mod book

# 创建相关服务目录
mkdir -p {common,service/{search/api,user/{api,model,rpc}}}

# 目录结构如下
book/
|-- README.md
|-- common
|-- go.mod
`-- service
    |-- search
    |   `-- api
    `-- user
        |-- api
        |-- model
        `-- rpc

```

### model生成
model是服务访问持久化数据层的桥梁，业务的持久化数据常存在于mysql，mongo等数据库中，我们都知道，对于一个数据库的操作莫过于CURD， 而这些工作也会占用一部分时间来进行开发，我曾经在编写一个业务时写了40个model文件，根据不同业务需求的复杂性，平均每个model文件差不多需要 10分钟，对于40个文件来说，400分钟的工作时间，差不多一天的工作量，而goctl工具可以在10秒钟来完成这400分钟的工作。
```bash
# 创建表
> user.sql
CREATE TABLE `user`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名称',
  `gender` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '男｜女｜未公开',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户密码',
  `number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '学号',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `number_unique` (`number`),
  UNIQUE KEY `name_unique` (`number`), # 此字段会生成对应的FindOneByName()接口方法
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

# 预设数据
INSERT INTO `user` (number,name,password,gender)values ('666','小明','123456','男');

# 生成model
cd service/user/model/

【生成带缓存的model, -c】
方式一: datasource
goctl model mysql datasource -url="root:123456@tcp(10.0.0.91:3306)/boot" -table="user" -dir . -c -style goZero
# model层缓存默认过期时间为7天，如果没有查到数据会设置一个空缓存，空缓存的过期时间为1分钟

方式二: ddl
goctl model mysql ddl -src="*.sql" -dir="." -c

说明: 如果不想带有缓存，则去掉-c

# 下载依赖
go mod tidy

model/
|-- user.sql
|-- userModel.go         // 扩展代码
|-- userModel_gen.go     // 扩展代码
`-- vars.go              // 定义常量和变量
```


### api服务生成
#### api文件编写
```bash
cd service/user/api/
# 生成api文件模板
goctl api -o user.api

# user.api  
...
type (
    LoginReq {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    LoginReply {
        Id           int64 `json:"id"`
        Name         string `json:"name"`
        Gender       string `json:"gender"`
        AccessToken  string `json:"accessToken"`
        AccessExpire int64 `json:"accessExpire"`
        RefreshAfter int64 `json:"refreshAfter"`
    }
)

service user-api {
    ...
    @handler login
    post /user/login (LoginReq) returns (LoginReply)
}

```

#### 生成api服务
```bash
goctl api go -api user.api -dir .

go mod tidy
```

### 业务编码
#### 添加Mysql配置
```bash
# service/user/api/internal/config/config.go

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"  # 引入redis依赖
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Mysql struct {       
		DataSource string  # 添加mysql配置信息
	}
	CacheRedis cache.CacheConf  # 添加redis配置信息
}

```


#### 添加连接信息
```bash
# service/user/api/etc/user-api.yaml
Name: user-api
Host: 0.0.0.0
Port: 8888

Mysql:
  # 字段与config保持一致
  DataSource: root:123456@tcp(10.0.0.91:3306)/book?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

# 字段与config保持一致
CacheRedis:
  - Host: 10.0.0.91:6379
    Pass: "123"
    Type: node

```

#### 添加服务依赖
```bash
# service/user/api/internal/svc/servicecontext.go

package svc

import (
	"book/service/user/api/internal/config"

	"book/service/user/model"   # 引入model依赖

	"github.com/zeromicro/go-zero/core/stores/sqlx"  # 引入sqlx依赖
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel  # 添加UserModel接口
}

func NewServiceContext(c config.Config) *ServiceContext {

	conn := sqlx.NewMysql(c.Mysql.DataSource) # 获取mysql连接
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, c.CacheRedis), # 建立连接
	}
}


```


#### 登录逻辑
```bash
# service/user/api/internal/logic/loginlogic.go
...
import (
	"context"
	"errors"
	"strings"  // 添加依赖包


	"book/service/user/api/internal/svc"
	"book/service/user/api/internal/types"
	"book/service/user/model"  // 引入model依赖包

	"github.com/zeromicro/go-zero/core/logx"
)

...
...
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginReply, error) {
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, errors.New("参数错误")
	}

                              // 调用model中FindOneByUsername接口，需要实现FindOneByUsername接口
	userInfo, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username) 
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("用户名不存在")
	default:
		return nil, err
	}

	if userInfo.Password != req.Password {
		return nil, errors.New("用户密码不正确")
	}

	return &types.LoginReply{
		Id:     userInfo.Id,
		Name:   userInfo.Name,
		Gender: userInfo.Gender,
		// AccessToken:  jwtToken,
		// AccessExpire: now + accessExpire,
		// RefreshAfter: now + accessExpire/2,
	}, nil
}
```

#### sql接口
```bash
# service/user/model/userModel_gen.go
# 默认只有Insert，FindOne，FindOneByNumber，Update，Delete等接口方法
type (
	userModel interface {
		Insert(ctx context.Context, data *User) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*User, error)
		FindOneByNumber(ctx context.Context, number string) (*User, error)
		FindOneByUsername(ctx context.Context, name string) (*User, error)  // 新增
		Update(ctx context.Context, data *User) error
		Delete(ctx context.Context, id int64) error
	}
    ...


...
# 实现接口方法
func (m *defaultUserModel) FindOneByUsername(ctx context.Context, name string) (*User, error) {
	bookUserNumberKey := fmt.Sprintf("%s%v", cacheBookUserNumberPrefix, name)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, bookUserNumberKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
        
        // 注意这原生sql语句
		query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, name); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

# findall
func (m *defaultUserInfoModel) FindAll(ctx context.Context) ([]*UserInfo, error) {
	query := fmt.Sprintf("select %s from %s ", userRows, m.table)

	var resp []*UserInfo
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
/*
# 业务代码
userInfoList, err := l.svcCtx.UserInfoModel.FindAll(l.ctx)

	if err != nil {
		return nil, err
	}

	var userInfoListTmp []types.UserInfo //
	_ = copier.Copy(&userInfoListTmp, userInfoList)

	return &types.ListResponse{
		Status:  200,
		Message: "success",
		Data:    userInfoListTmp,
	}, nil

*/

# 分页
func (m *defaultUserModel) FindPaginations(ctx context.Context, where string, Page, PageSize int64) ([]User, error) {
	var resp []User
	// 声明错误
	var err error
	// 分页使用的偏移量
	offset := (Page - 1) * PageSize
	// 没有检索条件：只有 `Page` 、`PageSize`
	if len(where) == 0 {
		query := fmt.Sprintf("select %s from %s order by create_time desc limit ?,?", productRows, m.table)
		err = m.conn.QueryRowsCtx(ctx, &resp, query, offset, PageSize)
	} else {
		// 有检索条件
		query := fmt.Sprintf("select %s from %s where %s order by create_time desc limit ?,?", productRows, m.table, where)
		err = m.conn.QueryRowsCtx(ctx, &resp, query, offset, PageSize)
	}
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

```


#### 测试
```bash
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Connection: keep-alive' \
--data-raw '{
    "username":"小明",   // 注意中文可能查不到
    "password":"123456"
}'

```





### JWT鉴权
#### 概述
JSON Web令牌（JWT）是一个开放标准（RFC 7519），它定义了一种紧凑而独立的方法，用于在各方之间安全地将信息作为JSON对象传输。由于此信息是经过数字签名的，因此可以被验证和信任。可以使用秘钥（使用HMAC算法）或使用RSA或ECDSA的公钥/私钥对对JWT进行签名。

#### 使用JWT场景
- 授权：这是使用JWT的最常见方案。一旦用户登录，每个后续请求将包括JWT，从而允许用户访问该令牌允许的路由，服务和资源。单一登录是当今广泛使用JWT的一项功能，因为它的开销很小并且可以在不同的域中轻松使用。
- 信息交换：JSON Web令牌是在各方之间安全地传输信息的一种好方法。因为可以对JWT进行签名（例如，使用公钥/私钥对），所以您可以确保发件人是他们所说的人。此外，由于签名是使用标头和有效负载计算的，因此您还可以验证内容是否未被篡改。


#### go-zero中使用jwt
##### user api生成jwt token
添加配置定义和yaml配置项
```bash
# service/user/api/internal/config/config.go

type Config struct {
    rest.RestConf
    Mysql struct{
        DataSource string
    }
    CacheRedis cache.CacheConf

    // 添加认证参数
    Auth      struct {
        AccessSecret string
        AccessExpire int64
    }
}



# service/user/api/etc/user-api.yaml

Auth:
  AccessSecret: "cfdvfwvwre324safceds"  # 生成jwt token的密钥
  AccessExpire: 10  # jwt token有效期，单位：秒


```

登录逻辑
```bash
# service/user/api/internal/logic/loginlogic.go

import (
    ...
    "time"

	"github.com/golang-jwt/jwt"  // 添加依赖
	
)


// 新增 getJwtToken接口方法
func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
  claims := make(jwt.MapClaims)
  claims["exp"] = iat + seconds
  claims["iat"] = iat
  claims["userId"] = userId
  token := jwt.New(jwt.SigningMethodHS256)
  token.Claims = claims
  return token.SignedString([]byte(secretKey))
}


func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginReply, error) {
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, errors.New("参数错误")
	}

    // userinfo 是sql查询的结果
	userInfo, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("用户名不存在")
	default:
		return nil, err
	}

	if userInfo.Password != req.Password {
		return nil, errors.New("用户密码不正确")
	}

    // jwt token获取
	// ---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, userInfo.Id)
	if err != nil {
		return nil, err
	}
	// ---end---

	return &types.LoginReply{
		Id:           userInfo.Id,
		Name:         userInfo.Name,
		Gender:       userInfo.Gender,
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

```

下载依赖
```bash
git mod tidy

```

##### search api使用jwt token鉴权
编写search.api文件
```bash
# service/search/api/

# 生成模板
goctl api -o search.api


# 添加请求/响应参数
type (
    SearchReq {
        // 图书名称
        Name string `form:"name"`
    }

    SearchReply {
        Name string `json:"name"`
        Count int `json:"count"`
    }
)

# 添加jwt
@server(
    jwt: Auth
)

# 添加请求方法
service search-api {
    @handler search
    get /search/do (SearchReq) returns (SearchReply)
}

service search-api {
    @handler ping
    get /search/ping
}


# 根据api文件生成go文件
goctl api go -api search.api -dir=.

```

添加yaml配置项
```bash
# service/search/api/etc/search-api.yaml

Name: search-api
Host: 0.0.0.0
Port: 8889   # 修改端口

Auth:
  AccessSecret: "cfdvfwvwre324safceds"  # 生成jwt token的密钥，必须要和user api中声明的一致
  AccessExpire: 10  # jwt token有效期，单位：秒；和user api中声明的一致

```

业务逻辑
```bash
# service/search/api/internal/logic/searchlogic.go

...
func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchReply, err error) {
	// todo: add your logic here and delete this line

	return &types.SearchReply{  // 返回参数
		Name:  req.Name,
		Count: 1,
	}, nil
}

```


##### 验证 jwt token
启动user api服务，登录
```bash
cd service/user/api/
go run user.go -f etc/user-api.yaml

# 发送请求，获取token
curl -i -X POST \
  http://127.0.0.1:8888/user/login \
  -H 'Content-Type: application/json' \
  -d '{
    "username":"666",   // 数据库信息要匹配
    "password":"123456"
}'
```

启动search api服务，调用/search/do验证jwt鉴权是否通过
```bash
cd service/search/api/
go run search.go -f etc/search-api.yaml

```

验证
- 不带token的请求
```bash
$ curl -i -X GET \
   'http://127.0.0.1:8889/search/do?name=%E8%A5%BF%E6%B8%B8%E8%AE%B0'

HTTP/1.1 401 Unauthorized
Traceparent: 00-3ed2ef862322d1f9ab3ece3a5e6358f3-c35ab610b9341214-00
Date: Tue, 14 Mar 2023 10:18:32 GMT
Content-Length: 0

```
- 使用token的请求
```bash
$ curl -i -X GET \
  'http://127.0.0.1:8889/search/do?name=laoshe' \
  -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzg4MDU5NDQsImlhdCI6MTY3ODgwMjM0NCwidXNlcklkIjoxfQ.j44lF-lacEEUta8M4SLvalJh0Nhg4MIOG1Gl_z4e2lw'

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Traceparent: 00-e2f4d34a41bf400554b50fa185747397-bf08aa0329189a8a-00
Date: Tue, 14 Mar 2023 10:22:19 GMT
Content-Length: 4

{"name":"laoshe","count":1}
```

#### 获取jwt token中携带的信息
go-zero从jwt token解析后会将用户生成token时传入的kv原封不动的放在http.Request的Context中，因此我们可以通过Context就可以拿到你想要的值
```bash
# /service/search/api/internal/logic/searchlogic.go

func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchReply, err error) {
	// todo: add your logic here and delete this line

	// 添加一个log来输出从jwt解析出来的userId
	logx.Infof("userId: %v", l.ctx.Value("userId")) // 这里的key和生成jwt token时传入的key一致

	return &types.SearchReply{
		Name:  req.Name,
		Count: 1,
	}, nil
}

```
测试请求
```bash
curl --location --request GET 'http://localhost:8889/search/do?name=%22%E9%AA%86%E9%A9%BC%E7%A5%A5%E5%AD%90%22' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzg4MDcyNDMsImlhdCI6MTY3ODgwMzY0MywidXNlcklkIjoxfQ.kJd4tnd_M4QvZFq9MGGmqEI8HX6hmg4OvvGwtc5-8mM' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Accept: */*' \
--header 'Host: localhost:8889' \
--header 'Connection: keep-alive'

# log
2023-03-14T 22:29:50.389+08:00   info   userId: 1       caller=logic/searchlogic.go:29

```



### 中间件使用

#### 中间件分类
在go-zero中，中间件可以分为路由中间件和全局中间件，
- 路由中间件是指某一些特定路由需要实现中间件逻辑，其和jwt类似，没有放在jwt:xxx下的路由不会使用中间件功能， 
- 而全局中间件的服务范围则是整个服务。

#### 中间件使用
这里以search服务为例来演示中间件的使用


##### 路由中间件
重新编写search.api文件，添加middleware声明
```bash
# service/search/api
...
@server(
	jwt: Auth
	middleware: Example // 路由中间件声明,只加这一个
)
service search-api {
	@handler search
	get /search/do (SearchReq) returns (SearchReply)
}

```

重新生成api代码
```bash

goctl api go -api search.api -dir .
```


完善资源依赖ServiceContext
```bash

# service/search/api/internal/svc/servicecontext.go
...

import (
	"book/service/search/api/internal/config"
	"book/service/search/api/internal/middleware"  // 引入middleware

	"github.com/zeromicro/go-zero/rest"  // 引入第三方依赖
)

type ServiceContext struct {
	Config config.Config

	Example rest.Middleware  // 添加依赖
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Example: middleware.NewExampleMiddleware().Handle,  // 结构体方法
	}
}


```

编写中间件逻辑
```bash
# service/search/api/internal/middleware/examplemiddleware.go

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx" // 引入日志
)

...
func (m *ExampleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {


	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

        logx.Info("---- 这是路由中间件 ----") // 打印日志

		// Passthrough to next handler if need
		next(w, r)
	}
}

```

启动服务验证
```bash
go run search.go -f etc/search-api.yaml


# 请求
curl --location --request GET 'http://localhost:8889/search/do?name=%22%E9%AA%86%E9%A9%BC%E7%A5%A5%E5%AD%90%22' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzg4MDcyNDMsImlhdCI6MTY3ODgwMzY0MywidXNlcklkIjoxfQ.kJd4tnd_M4QvZFq9MGGmqEI8HX6hmg4OvvGwtc5-8mM' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Accept: */*' \
--header 'Host: localhost:8889' \
--header 'Connection: keep-alive'

# log
2023-03-14T 22:52:29.544+08:00   info   ---- 这是路由中间件 ----  caller=middleware/examplemiddleware.go:18
```


##### 全局中间件
通过rest.Server提供的Use方法即可
```bash
# service/search/api/search.go
import (

	"net/http"  // 添加依赖

)

...
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	// 全局中间件
    // -- start -- //
    server.Use(func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            logx.Info("---- 这是全局中间件 ----")
            next(w, r)
        }
    })
    // -- end -- //

	handler.RegisterHandlers(server, ctx)

	// 禁用stat日志
	logx.DisableStat()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}


```

测试请求
```bash

# log
023-03-14T 23:23:21.777+08:00   info   ---- 这是全局中间件 ----        caller=api/search.go:33
2023-03-14T 23:23:21.778+08:00   info   ---- 这是路由中间件 ----        caller=middleware/examplemiddleware.go:20
2023-03-14T 23:23:21.779+08:00   info   userId: 1       caller=logic/searchlogic.go:29
```


##### 在中间件里调用其它服务

通过闭包的方式把其它服务传递给中间件
```bash
// 模拟的其它服务
type AnotherService struct{}

func (s *AnotherService) GetToken() string {
    return stringx.Rand()
}

// 常规中间件
func middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("X-Middleware", "static-middleware")
        next(w, r)
    }
}

// 调用其它服务的中间件
func middlewareWithAnotherService(s *AnotherService) rest.Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            w.Header().Add("X-Middleware", s.GetToken())
            next(w, r)
        }
    }
}

```
完整代码参考：https://github.com/zeromicro/zero-examples/tree/main/http/middleware



### rpc编写与调用
在一个大的系统中，多个子系统（服务）间必然存在数据传递，有数据传递就需要通信方式，你可以选择最简单的http进行通信，也可以选择rpc服务进行通信， 在go-zero，我们使用zrpc来进行服务间的通信，zrpc是基于grpc。

#### 场景
在前面我们完善了对用户进行登录，用户查询图书等接口协议，但是用户在查询图书时没有做任何用户校验，如果当前用户是一个不存在的用户则我们不允许其查阅图书信息， 从上文信息我们可以得知，需要user服务提供一个方法来获取用户信息供search服务使用，因此我们就需要创建一个user rpc服务，并提供一个getUser方法。

#### rpc服务编写
编译proto文件
```bash
# service/user/rpc/
生成proto模板
goctl rpc -o user.proto

# user.proto
syntax = "proto3";

package user;
option go_package="./user";

message Request {
  string ping = 1;
}

message Response {
  string pong = 1;
}

// 新增请求/响应参数
message IdReq{
  int64 id = 1;
}

message UserInfoReply{
  int64 id = 1;
  string name = 2;
  string number = 3;
  string gender = 4;
}

// 绑定路由
service User {
  rpc getUser(IdReq) returns(UserInfoReply);
  rpc Ping(Request) returns(Response);
}



```

生成rpc服务代码
```bash
goctl env check -i -f

goctl rpc protoc user.proto --go_out=./types --go-grpc_out=./types --zrpc_out=.

注意： 每一个 *.proto文件只允许有一个service
```

添加配置及完善yaml配置项
```bash
# service/user/rpc/internal/config/config.go
...
import (
	"github.com/zeromicro/go-zero/core/stores/cache"  // 添加redis依赖
	"github.com/zeromicro/go-zero/zrpc"
)


type Config struct {
	zrpc.RpcServerConf

	Mysql struct {          # 添加mysql配置
        DataSource string
    }
    CacheRedis cache.CacheConf  # 添加redis配置
}


# service/user/rpc/etc/user.yaml
Name: user.rpc
ListenOn: 0.0.0.0:8090   # 修改端口
Etcd:
  Hosts:
  - 10.0.0.91:2379    # etcd连接地址
  Key: user.rpc       # 当前服务注册的key

Mysql:   # MySQL连接信息
  DataSource: root:123456@tcp(10.0.0.91:3306)/book?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

CacheRedis:  # redis连接信息
  - Host: 10.0.0.91:6379
    Pass: "123"
    Type: node   # redis类型，单机

```

添加资源依赖
```bash
# service/user/rpc/internal/svc/servicecontext.go
package svc

import (
	"book/service/user/model"  // 引入model依赖
	"book/service/user/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"  // 引入sqlx依赖
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel  // 引入usermodel接口
}

func NewServiceContext(c config.Config) *ServiceContext {

	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, c.CacheRedis), // 建立关联
	}
}


```

添加rpc逻辑
```bash
# service/user/rpc/internal/logic/getuserlogic.go
...
func (l *GetUserLogic) GetUser(in *user.IdReq) (*user.UserInfoReply, error) {
	// todo: add your logic here and delete this line

	// 调用usermodel的findone接口
	one, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 返回值
	return &user.UserInfoReply{
		Id:     one.Id,
		Name:   one.Name,
		Number: one.Number,
		Gender: one.Gender,
	}, nil

	// return &user.UserInfoReply{}, nil
}


```

#### 使用rpc
> 在search服务中调用user rpc

添加UserRpc配置及yaml配置项
```bash
# service/search/api/internal/config/config.go
...
import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"  // 添加依赖
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	UserRpc zrpc.RpcClientConf  // 添加rpc配置
}


# service/search/api/etc/search-api.yaml
...
UserRpc:
  Etcd:
    Hosts:
    - 10.0.0.91:2379   # etcd连接信息
    Key: user.rpc  # user rpc的keys
  NonBlock: true  // 弱依赖，不会因为其他服务无法启动而不能正常运行

```

添加依赖
```bash
# service/search/api/internal/svc/servicecontext.go
package svc

import (
	"book/service/search/api/internal/config"
	"book/service/search/api/internal/middleware"
	"book/service/user/rpc/userclient" // 引入userclient 依赖

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	Example rest.Middleware

	UserRpc userclient.User // 使用user接口
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Example: middleware.NewExampleMiddleware().Handle,

		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)), // 绑定连接
	}
}


```

补充逻辑
```bash
# service/search/api/internal/logic/searchlogic.go
...
import (
    ...
	"encoding/json"
	"fmt"

	"book/service/user/rpc/types/user"  // 引入user

)
...
func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchReply, err error) {
	// todo: add your logic here and delete this line

	// 添加一个log来输出从jwt解析出来的userId
	logx.Infof("userId: %v", l.ctx.Value("userId")) // 这里的key和生成jwt token时传入的key一致

    // 类型转换
	userIdNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	logx.Infof("userId: %s", userIdNumber)
	userId, err := userIdNumber.Int64()
	if err != nil {
		return nil, err
	}
	// 使用user rpc的响应结构体
	var people *user.UserInfoReply
    // 向user rpc发送请求
	people, err = l.svcCtx.UserRpc.GetUser(l.ctx, &user.IdReq{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("user rpc返回的值是: ---> ", people, " <---")


	return &types.SearchReply{
		Name:  people.Name,  // 获取属性值
		Count: 100,
	}, nil
}



```
启动并验证服务
> 注意,先启动运行 user rpc服务,然后再启动user api服务(提供jwt token),然后启动search api服务
>

启动user rpc服务
```bash
#  service/user/rpc
go run user.go -f etc/user.yaml

```

启动user api服务
```bash
# service/user/api
go run user.go -f etc/user-api.yaml
```
启动search api服务
```bash
# service/search/api
go run search.go -f etc/search-api.yaml

```

验证服务
```bash
# 获取token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "username":"小明",
    "password":"123456"
}'

响应:
{
    "id": 1,
    "name": "小明",
    "gender": "男",
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzg4NTE1MDUsImlhdCI6MTY3ODg0NzkwNSwidXNlcklkIjoxfQ.ZFEta8tcDd-3ByTm3XjuN2JvtVhHSX-iANB1G6Ozjug",
    "accessExpire": 1678851505,
    "refreshAfter": 1678849705
}


# 请求search api
curl --location --request GET 'http://localhost:8889/search/do?name=%22%E9%AA%86%E9%A9%BC%E7%A5%A5%E5%AD%90%22' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzg4NTE1MDUsImlhdCI6MTY3ODg0NzkwNSwidXNlcklkIjoxfQ.ZFEta8tcDd-3ByTm3XjuN2JvtVhHSX-iANB1G6Ozjug' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Accept: */*' \
--header 'Host: localhost:8889' \
--header 'Connection: keep-alive'


# search api日志
2023-03-15T 10:56:08.922+08:00   info   userId: 1       caller=logic/searchlogic.go:36
user rpc返回的值是: --->  id:1 name:"小明" number:"666" gender:"男"  <---
2023-03-15T 10:56:08.926+08:00   info   [HTTP]  200  -  GET  /search/do?name=%22%E9%AA%86%E9%A9%BC%E7%A5%A5%E5%AD%90%22 - [::1]:1870 - Apifox/1.0.0 (https://www.apifox.cn)      duration=5.9ms  caller=handler/loghandler.go:160 trace=fb23092705e4cdb742a39372b9aba35b  span=89f304a1c6191f08


# search api响应
{
    "name": "小明",
    "count": 100
}

```

#### 调试gRPC

```bash

# 用go install安装grpcurl工具
go install github.com/fullstorydev/grpcurl/cmd/grpcurl

# 使用
grpcurl -plaintext 127.0.0.1:8081 list  // 查看接口
grpcurl -plaintext 127.0.0.1:8081 list product.Product

grpcurl -plaintext -d '{"product_id": 1}' 127.0.0.1:8081 product.Product.Product  // 调用Product接口查询id为1的商品数据

```


### 短链服务
1. 什么是短链服务？
短链服务就是将长的URL网址，通过程序计算等方式，转换为简短的网址字符串。

写此短链服务是为了从整体上演示go-zero构建完整微服务的过程，算法和实现细节尽可能简化了，所以这不是一个高阶的短链服务。
2. 短链微服务架构
                      |<-> Shorten/rpc <-> Redis <-> Mysql
Reqs -> Api Gateway ->|
                      |<-> Expand/rpc  <-> Redis <-> Mysql

- 这里把shorten和expand分开为两个微服务，并不是说一个远程调用就需要拆分为一个微服务，只是为了最简演示多个微服务而已
- 后面的redis和mysql也是共用的，但是在真正项目里要尽可能每个微服务使用自己的数据库，数据边界要清晰

3. 准备工作
数据库: etcd,redis,mysql
go配置/依赖: go,goctl
初始化目录:
```bash
mkdir -p service/shorturl/{api,rpc}
```

#### shorturl rpc
此服务是实现
- proto文件
```bash
# 生成proto模板文件, service/shorturl/rpc
goctl rpc template -o transform.proto

# 编辑proto文件
syntax = "proto3";

package transform;
option go_package="./transform";

message expandReq {
  string shorten = 1;
}

message expandResp {
  string url = 1;
}

message shortenReq {
  string url = 1;
}

message shortenResp {
  string shorten = 1;
}

service transformer {
  rpc expand(expandReq) returns(expandResp);
  rpc shorten(shortenReq) returns(shortenResp);
}


# 生成go文件
goctl rpc protoc transform.proto --go_out=./types --go-grpc_out=./types --zrpc_out=.

```
目录结构:
```bash
expand/rpc/
|-- etc
|   `-- transform.yaml             // 配置文件
|-- internal
|   |-- config
|   |   `-- config.go              // 配置定义一些依赖，对应于yaml文件
|   |-- logic
|   |   |-- expandlogic.go         // expand 业务逻辑编写的地方
|   |   `-- shortenlogic.go        // shorten 业务逻辑编写的地方
|   |-- server
|   |   `-- transformerserver.go   // 调用入口，不需要修改
|   `-- svc
|       `-- servicecontext.go      // 定义 ServiceContext，传递依赖
|-- transform.go
|-- transform.proto
|-- transformer
|   `-- transformer.go             // 提供了外部调用方法，无需修改
`-- types
    `-- transform
        |-- transform.pb.go
        `-- transform_grpc.pb.go   // request/response结构体定义

```

- 修改配置文件
```bash
# service\shorturl\rpc\etc\transform.yaml
Name: transform.rpc
ListenOn: 0.0.0.0:8090  # 修改端口
Etcd:
  Hosts:
  - 10.0.0.91:2379   # 修改etcd服务地址
  Key: transform.rpc


```
- 启动运行
```bash
# service/shorturl/rpc/
go run transform.go -f etc/transform.yaml

# etcd查看
$ etcdctl get transform.rpc --prefix
$ etcdctl get transform.rpc --prefix
transform.rpc/7587869393952790860
192.168.206.1:8090
```

#### api getway
- api文件
```bash
# service/shorturl/api
goctl api --o=shorturl.api

# 编辑shorturl.api
syntax = "v1"

info (
	title: // TODO: add title
	desc: // TODO: add description
	author: "hu417"
	email: "hu729919300@163.com"
)

type (
    expandReq {
        shorten string `form:"shorten"`
    }
    expandResp {
        url string `json:"url"`
    }
)

type (
    shortenReq {
        url string `form:"url"`
    }
    shortenResp {
        shorten string `json:"shorten"`
    }
)

service shorturl-api { 
	@doc "短链请求接口"   
	@server(
        handler: ShortenHandler
    )
    get /shorten(shortenReq) returns(shortenResp)
    
	@doc "短链生成接口"
	@server(
        handler: ExpandHandler
    )
    get /expand(expandReq) returns(expandResp)
}


# 生成go文件
goctl api go --api shorturl.api --dir=.
``
目录结构如下:
​```bash
shorturl/api/
|-- etc
|   `-- shorturl-api.yaml        // 配置文件
|-- internal
|   |-- config
|   |   `-- config.go            // 配置文件
|   |-- handler
|   |   |-- expandhandler.go     // 实现expandHandler
|   |   |-- routes.go            // 定义路由处理
|   |   `-- shortenhandler.go    // 实现shortenHandler
|   |-- logic 
|   |   |-- expandlogic.go       // 实现ExpandLogic业务逻辑
|   |   `-- shortenlogic.go      // 实现ShortenLogic业务逻辑
|   |-- svc
|   |   `-- servicecontext.go    // 定义ServiceContext
|   `-- types 
|       `-- types.go             // 定义请求、返回结构体
|-- shorturl.api                 // api文件
`-- shorturl.go                  // main入口定义

```
- 业务逻辑
```bash
# service\shorturl\api\internal\logic\shortenlogic.go
...
func (l *ShortenLogic) Shorten(req *types.ShortenReq) (resp *types.ShortenResp, err error) {
	// todo: add your logic here and delete this line

	// 响应参数设置
	resp = &types.ShortenResp{
		ShortUrl: req.Url,  // 返回值是请求参数值
	}
	return
}

```
- 启动服务
```bash
# 下载依赖
go mod tidy
# service\shorturl\api
go run shorturl.go -f etc/shorturl-api.yaml

# 测试
Get http://localhost:8888/shorten?url=http://www.xiaoheiban.cn
```


#### api gateway调用expand rpc服务
- api yaml文件
```bash
# service/shorturl/api/shorturl-api.yaml
Name: shorturl-api
Host: 0.0.0.0
Port: 8888

# Transform rpc
Transform:
  Etcd:
    Hosts:
    - 10.0.0.91:2379  # 通过 etcd 自动去发现可用的 transform 服务
    Key: transform.rpc

```
- api config文件
```bash

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Transform zrpc.RpcClientConf
}
```
- api servicecontext文件
```bash

import (
	"book/service/shorturl/api/internal/config"
	"book/service/shorturl/rpc/transformer"   // 添加rpc客户端

	"github.com/zeromicro/go-zero/zrpc"   // zrpc依赖
)

type ServiceContext struct {
	Config      config.Config
	Transformer transformer.Transformer  // 接口
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		Transformer: transformer.NewTransformer(zrpc.MustNewClient(c.Transform)),
	}
}

```
- 业务逻辑处理
  - api expandlogic文件
```bash




```





### 错误处理

错误的处理是一个服务必不可缺的环节。在平时的业务开发中，我们可以认为http状态码不为2xx系列的，都可以认为是http请求错误， 并伴随响应的错误信息，但这些错误信息都是以plain text形式返回的。除此之外，我在业务中还会定义一些业务性错误，常用做法都是通过 code、msg 两个字段来进行业务处理结果描述，并且希望能够以json响应体来进行响应。

#### 业务错误响应格式
- 业务处理正常
```bash
{
  "code": 0,
  "msg": "successful",
  "data": {
    ....
  }
}

# 参考格式
{"timestamp":"2023-03-22 09:23:03","status":999,"error":"None","message":"No message available"}
{"code":"NoSuchKey","message":"The specified key does not exist.","requestId":"736eb45b-4371-4873-a366-ab3307aa9cf4"}
```
- 业务处理异常
```bash
{
  "code": 10001,
  "msg": "参数错误"
}
```

#### user api之login
在之前，我们在登录逻辑中处理用户名不存在时，直接返回来一个error。我们来登录并传递一个不存在的用户名看看效果。
```bash
curl -X POST \
  http://127.0.0.1:8888/user/login \
  -H 'content-type: application/json' \
  -d '{
    "username":"111",
    "password":"123456"
}'

用户名不存在
```
接下来我们将其以json格式进行返回

#### 自定义错误
- 先在common中添加一个baseerror.go文件，并填入代码
```bash
# 回到项目根目录
cd common
mkdir errorx && cd errorx
touch baseerror.go

# baseerror.go
...
package errorx

const defaultCode = 1001

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CodeErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func NewDefaultError(msg string) error {
	return NewCodeError(defaultCode, msg)
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}


```

- 将登录逻辑中错误用CodeError自定义错误替换
```bash
# service/user/api/internal/logic/loginlogic.go
...
import (

	"book/common/errorx" // 引入依赖

)
...
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginReply, error) {
	// if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
	// 	return nil, errors.New("参数错误")
	// }
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, errorx.NewDefaultError("参数错误")   // 使用自定义error
	}

	userInfo, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	switch err {
	case nil:
	case model.ErrNotFound:
		// return nil, errors.New("用户名不存在")
		return nil, errorx.NewDefaultError("用户名不存在")   // 使用自定义error
	default:
		return nil, err
	}

	if userInfo.Password != req.Password {
		// return nil, errors.New("用户密码不正确")
		return nil, errorx.NewDefaultError("用户密码不正确")  // 使用自定义error
	}

	// ---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, userInfo.Id)
	if err != nil {
		return nil, err
	}
	// ---end---

	return &types.LoginReply{
		Id:           userInfo.Id,
		Name:         userInfo.Name,
		Gender:       userInfo.Gender,
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

```

- 开启自定义错误
```bash
# service/user/api/user.go
...
import (
	"context"
	"net/http"

	"book/common/errorx"

	"github.com/zeromicro/go-zero/rest/httpx"
)
...
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.DisableStat()

	// 自定义错误
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.CodeError:
			return http.StatusOK, e.Data()
		default:
			return http.StatusInternalServerError, nil
		}
	})
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

```

- 重启服务验证
```bash
# service/user/api
go run user.go -f etc/user-api.yaml

# 请求
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "username":"小",
    "password":"123456"
}'

# 响应
{
    "code": 1001,
    "msg": "用户名不存在"
}
```


### 响应模板修改
#### 场景
实现统一格式的body响应，格式如下
```bash
{
  "code": 0,
  "msg": "OK",
  "data": {} // 实际响应数据
}

```

#### 准备工作
我们提前在module为greet的工程下的response包中写一个Response方法，目录树类似如下
以user api示例
```bash
# service/user
mkdir -p greet/response && cd greet/response/
touch response.go

# response.go
...
package response

import (
    "net/http"

    "github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
    var body Body
    if err != nil {
        body.Code = -1
        body.Msg = err.Error()
    } else {
        body.Msg = "OK"
        body.Data = resp
    }
    httpx.OkJson(w, body)
}

```

#### 修改handler模板
- 全局模板设置
初始化handler模板
```bash
goctl template init
```

修改handler模板
```bash
# ~/.goctl/${goctl版本号}/api/handler.tpl
import (

	"/{mod}/service/user/greet/response" // 引入依赖

)
...
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		{{end}}
		
		l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})

        # 使用respnse方法,替换注释的内容
        response.Response(w, resp, err)

		# if err != nil {
		# 	httpx.ErrorCtx(r.Context(), w, err)
		# } else {
		# 	{{if .HasResp}}httpx.OkJsonCtx(r.Context(), w, resp){{else}}httpx.Ok(w){{end}}
		# }
	}
}

```
生成go文件
```bash
goctl api go -api xxx.api -dir .

```

- 当前服务使用
```bash
# service/user/api/internal/handler/loginhandler.go
...
import (

	"book/service/user/greet/response" // 添加依赖

)
...
func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)

		response.Response(w, resp, err)  // 添加response方法

		// if err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// } else {
		// 	httpx.OkJsonCtx(r.Context(), w, resp)
		// }
	}
}


```

#### 测试响应
重新启动运行user api服务
```bash
# service/user/api
go run user.go -f etc/user-api.yaml

```

发送请求
```bash
1. 发送正确的参数
# 发送请求
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "username":"小明",
    "password":"123456"
}'

# 响应
{
    "code": 0,
    "msg": "OK",
    "data": {
        "id": 1,
        "name": "小明",
        "gender": "男",
        "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzg4NjM3NTksImlhdCI6MTY3ODg2MDE1OSwidXNlcklkIjoxfQ.7SK87B2uQaQzqQHbNTRTV5PdTPRqglANIPrwW2t6pak",
        "accessExpire": 1678863759,
        "refreshAfter": 1678861959
    }
}


2. 发送错误的参数
# 发送请求
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "username":"111",
    "password":"123456"
}'

# 响应
{
    "code": -1,
    "msg": "用户名不存在"
}

```

#### 总结
本文档仅对http相应为例讲述了自定义模板的流程，除此之外，自定义模板的场景还有：
- model 层添加kmq
- model 层生成待有效期option的model实例
- http自定义相应格式

### 监控网络连接

当我在做后端开发和写 go-zero 的时候，经常会需要监控网络连接，分析请求内容。比如：

- 分析 gRPC 连接何时连接、何时重连，并据此调整各种参数，比如：MaxConnectionIdle
- 分析 MySQL 连接池，当前多少连接，连接的生命周期是什么策略
- 也可以用来观察和分析任何 TCP 连接，看服务端主动断，还是客户端主动断等等

#### tproxy 的安装

```bash
# go 
$ GOPROXY=https://goproxy.cn/,direct go install github.com/kevwan/tproxy@latest


使用 docker 镜像：
$ docker run --rm -it -p <listen-port>:<listen-port> -p <remote-port>:<remote-port> kevinwan/tproxy:v1 tproxy -l 0.0.0.0 -p <listen-port> -r host.docker.internal:<remote-port>

参数
$ tproxy --help
Usage of tproxy:
  -d duration
            the delay to relay packets
  -l string
            Local address to listen on (default "localhost")
  -p int
            Local port to listen on
  -q        Quiet mode, only prints connection open/close and stats, default false
  -r string
            Remote address (host:port) to connect
  -t string
            The type of protocol, currently support grpc
```

#### 分析 gRPC 连接

```bash
tproxy -p 8088 -r localhost:8081 -t grpc -d 100ms
侦听在 localhost 和 8088 端口
重定向请求到 localhost:8081
识别数据包格式为 gRPC
数据包延迟100毫秒


```

#### 分析 MySQL 连接

```bash

hey -c 10 -z 10s "http://localhost:8888/lookup?url=go-zero.dev"
发为10QPS且持续10秒钟的压测


mysql连接池参数:
maxIdleConns = 3
maxOpenConns = 8
maxLifetime  = time.Minute
# ConnMaxLifetime 一定要设置的小于 wait_timeout
```







## 集成

### Swagger

1. 安装swagger插件

```bash
GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/goctl-swagger@latest

```

2. 配置环境
   将$GOPATH/bin中的goctl-swagger添加到环境变量

3. 使用

- 创建api文件

```bash
info(
 title: "type title here"
 desc: "type desc here"
 author: "type author here"
 email: "type email here"
 version: "type version here"
)


type (
 RegisterReq {
  Username string `json:"username"`
  Password string `json:"password"`
  Mobile string `json:"mobile"`
 }
 
 LoginReq {
  Username string `json:"username"`
  Password string `json:"password"`
 }
 
 UserInfoReq {
  Id string `path:"id"`
 }
 
 UserInfoReply {
  Name string `json:"name"`
  Age int `json:"age"`
  Birthday string `json:"birthday"`
  Description string `json:"description"`
  Tag []string `json:"tag"`
 }
 
 UserSearchReq {
  KeyWord string `form:"keyWord"`
 }
)

service user-api {
    @doc(
    summary: "注册"
    )
    @handler register
    post /api/user/register (RegisterReq)
    
    @doc(
    summary: "登录"
    )
    @handler login
    post /api/user/login (LoginReq)
    
    @doc(
    summary: "获取用户信息"
    )
    @handler getUserInfo
    get /api/user/:id (UserInfoReq) returns (UserInfoReply)
    
    @doc(
    summary: "用户搜索"
    )
    @handler searchUser
    get /api/user/search (UserSearchReq) returns (UserInfoReply)
}
```

- 生成swagger.json 文件

```bash
goctl api plugin -plugin goctl-swagger="swagger -filename user.json" -api user.api -dir .

```

- 指定Host，basePath api-host-and-base-path

```bash
goctl api plugin -plugin goctl-swagger="swagger -filename user.json -host 0.0.0.0 -basepath /api" -api user.api -dir .

```

- swagger ui 查看生成的文档

```bash
docker run --rm -p 8083:8080 -e SWAGGER_JSON=/foo/user.json -v $PWD:/foo swaggerapi/swagger-ui
```

- Swagger Codegen 生成客户端调用代码(go,javascript,php)

```bash
for l in go javascript php; do
  docker run --rm -v "$(pwd):/go-work" swaggerapi/swagger-codegen-cli generate \
    -i "/go-work/rest.swagger.json" \
    -l "$l" \
    -o "/go-work/clients/$l"
done
```



### 日志

#### logx

go-zero的 ***logx*** 包提供了日志功能，默认不需要做任何配置就可以在stdout中输出日志。当我们请求/v1/order/list接口的时候输出日志如下，默认是json格式输出，包括时间戳，http请求的基本信息，接口耗时，以及链路追踪的span和trace信息。

~~~go
type Logger interface {
    // Error logs a message at error level.
    Error(...interface{})
    // Errorf logs a message at error level.
    Errorf(string, ...interface{})
    // Errorv logs a message at error level.
    Errorv(interface{})
    // Errorw logs a message at error level.
    Errorw(string, ...LogField)
    // Info logs a message at info level.
    Info(...interface{})
    // Infof logs a message at info level.
    Infof(string, ...interface{})
    // Infov logs a message at info level.
    Infov(interface{})
    // Infow logs a message at info level.
    Infow(string, ...LogField)
    // Slow logs a message at slow level.
    Slow(...interface{})
    // Slowf logs a message at slow level.
    Slowf(string, ...interface{})
    // Slowv logs a message at slow level.
    Slowv(interface{})
    // Sloww logs a message at slow level.
    Sloww(string, ...LogField)
    // WithContext returns a new logger with the given context.
    WithContext(context.Context) Logger
    // WithDuration returns a new logger with the given duration.
    WithDuration(time.Duration) Logger
}
~~~

logx 配置

```bash
type LogConf struct {
    ServiceName         string `json:",optional"`
    Mode                string `json:",default=console,options=[console,file,volume]"`
    Encoding            string `json:",default=json,options=[json,plain]"`
    TimeFormat          string `json:",optional"`
    Path                string `json:",default=logs"`
    Level               string `json:",default=info,options=[info,error,severe]"`
    Compress            bool   `json:",optional"`
    KeepDays            int    `json:",optional"`
    StackCooldownMillis int    `json:",default=100"`
}
```

- `ServiceName`：设置服务名称，可选。在 `volume` 模式下，该名称用于生成日志文件。在 `rest/zrpc` 服务中，名称将被自动设置为 `rest`或`zrpc` 的名称。
- `Mode`：输出日志的模式，默认console
  - `console` 模式将日志写到 `stdout/stderr`
  - `file` 模式将日志写到 `Path` 指定目录的文件中
  - `volume` 模式在 docker 中使用，将日志写入挂载的卷中
- `Encoding`: 指示如何对日志进行编码，默认是json
  - `json`模式以 json 格式写日志
  - `plain`模式用纯文本写日志，并带有终端颜色显示
- `TimeFormat`：自定义时间格式，可选。默认是 `2006-01-02T15:04:05.000Z07:00`
- `Path`：设置日志路径，默认为 `logs`
- `Level`: 用于过滤日志的日志级别。默认为`info`
  - `info`，所有日志都被写入
  - `error`, `info` 的日志被丢弃
  - `severe`, `info` 和 `error` 日志被丢弃，只有 `severe` 日志被写入
- `Compress`: 是否压缩日志文件，只在 `file` 模式下工作
- `KeepDays`：日志文件被保留多少天，在给定的天数之后，过期的文件将被自动删除。对 `console` 模式没有影响
- `StackCooldownMillis`：多少毫秒后再次写入堆栈跟踪。用来避免堆栈跟踪日志过多

配置

```bash
# user/etc/user-api.yaml
// 1、控制台输出
Log:
  ServiceName: app
  Mode: console       # 日志模式，[console,file,volume]
  Encoding: plain   # 输出格式，plain换行，json是一整行
  TimeFormat: "2006-01-02T 15:04:05.000Z07:00"  # 时间格式
  Level: info   # [debug,info,error,severe]

// 2、写入到文件
Log:
  ServiceName: app
  Mode: console       # 日志模式，[console,file,volume]
  # Mode: file
  # Path: logs
  Encoding: plain   # 输出格式，plain换行，json是一整行
  TimeFormat: "2006-01-02T 15:04:05.000Z07:00"  # 时间格式
  Level: info   # [debug,info,error,severe]
  # Compress: true  # 启用压缩
  # KeepDays: 7     # 保留天数            int    `json:",optional"`
  # StackCooldownMillis: 100  # 多少毫秒后再次写入堆栈跟踪


# user.go
func main(){
	...
    // 关闭stat日志
    logx.DisableStat()
    // logx.SetLevel(logx.ErrorLevel) // 只打印error日志
	...
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
```

使用

```bash
# user/api/internal/logic/loginlogic.go
...
	logx.info("xxxx")
	// 可以扩展日志输出的字段，添加了uid字段记录请求的用户的uid
	logx.Infow("order list", logx.Field("uid",req.UID))
```



#### zap

在实际的使用中，记录日志可能更倾向于使用一些更成熟的库，比如zap

我们可以将zap做为go-zero的实现

~~~go
// 以user api为例
# user/zap/zap.go
package zapx

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

const callerSkipOffset = 3

type ZapWriter struct {
	logger *zap.Logger
}

func NewZapWriter(opts ...zap.Option) (logx.Writer, error) {
	opts = append(opts, zap.AddCallerSkip(callerSkipOffset))
	logger, err := zap.NewProduction(opts...)
	if err != nil {
		return nil, err
	}

	return &ZapWriter{
		logger: logger,
	}, nil
}

func (w *ZapWriter) Alert(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

func (w *ZapWriter) Close() error {
	return w.logger.Sync()
}

func (w *ZapWriter) Debug(v interface{}, fields ...logx.LogField) {
	w.logger.Debug(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.Error(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Severe(v interface{}) {
	w.logger.Fatal(fmt.Sprint(v))
}

func (w *ZapWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.Warn(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Stack(v interface{}) {
	w.logger.Error(fmt.Sprint(v), zap.Stack("stack"))
}

func (w *ZapWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func toZapFields(fields ...logx.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
~~~

应用：

~~~go
# user.go
...
func main() {
    ...
    writer, err := zapx.NewZapWriter()
	logx.Must(err)
	logx.SetWriter(writer)
	...	
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

~~~

业务逻辑

```bash
# user/api/internal/logic/loginlogic.go
...
	logx.info("xxxx")

```





### Prometheus

在微服务的开发当中，监控也是一件非常重要的事情，很多线上问题都需要通过监控来触发告警，从而进行及时处理。

Prometheus是目前应用最广，使用最多的监控中间件。

同样，我们先部署prometheus

~~~yaml
  prometheus:
    container_name: prometheus
    image: bitnami/prometheus:2.40.7
    environment:
      - TZ=Asia/Shanghai
    privileged: true
    volumes:
      - ${PRO_DIR}/prometheus.yml:/opt/bitnami/prometheus/conf/prometheus.yml  # 将 prometheus 配置文件挂载到容器里
      - ${PRO_DIR}/target.json:/opt/bitnami/prometheus/conf/targets.json  # 将 prometheus 配置文件挂载到容器里
    ports:
      - "9090:9090"                     # 设置容器9090端口映射指定宿主机端口，用于宿主机访问可视化web
    restart: always
~~~

prometheus.yml

~~~yaml
# my global config
global:
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
# - "first_rules.yml"
# - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'file_ds'
    file_sd_configs:
      - files:
          - targets.json
~~~



targets.json

~~~json
[
  {
    "targets": ["host.docker.internal:9081"],
    "labels": {
      "job": "user-api",
      "app": "user-api",
      "env": "test",
      "instance": "host.docker.internal:8888"
    }
  },
  {
    "targets": ["host.docker.internal:9091"],
    "labels": {
      "job": "user-rpc",
      "app": "user-rpc",
      "env": "test",
      "instance": "host.docker.internal:8080"
    }
  }
]
~~~

在userapi添加配置

~~~yaml
# user/api/etc/user-api.yaml
Prometheus:
  Host: 127.0.0.1
  Port: 9081
  Path: /metrics
~~~

在user模块添加配置

~~~yaml
# user/rpc/etc/user-rpc.yaml
Prometheus:
  Host: 127.0.0.1
  Port: 9091
  Path: /metrics
~~~

访问：http://localhost:9090/targets?search=



### jaeger

jaeger是一个用于链路追踪的中间件。

同样，先安装jaeger

docker-compose.yaml

~~~yaml
  jaeger:
    container_name: jaeger
    image: rancher/jaegertracing-all-in-one:1.20.0
    environment:
      - TZ=Asia/Shanghai
      - SPAN_STORAGE_TYPE=elasticsearch
      - ES_SERVER_URLS=http://elasticsearch:9200
      - LOG_LEVEL=debug
    privileged: true
    ports:
      - "6831:6831/udp"
      - "6832:6832/udp"
      - "5778:5778"
      - "16686:16686"
      - "4317:4317"
      - "4318:4318"
      - "14250:14250"
      - "14268:14268"
      - "14269:14269"
      - "9411:9411"
    restart: always
  elasticsearch:
    container_name: elasticsearch
    image: elasticsearch:7.13.1
    environment:
      - TZ=Asia/Shanghai
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    privileged: true
    ports:
      - "9200:9200"
    restart: always
~~~

在之前的userapi和user模块yaml文件中加入配置：

~~~yaml
# user/api/etc/user-api.yaml
Telemetry:
  Name: user-api
  Endpoint: http://localhost:14268/api/traces
  Sampler: 1.0
  Batcher: jaeger
~~~

~~~yaml
# user/rpc/etc/user-rpc.yaml
Telemetry:
  Name: user-rpc
  Endpoint: http://localhost:14268/api/traces
  Sampler: 1.0
  Batcher: jaeger
~~~

启动访问`http://localhost:16686/`进行测试。



### 分布式事务

分布式事务也是微服务架构中必不可少的。

go-zero使用了dtm的方案来解决分布式事务问题，dtm也是国人开发。

dtm地址：https://www.dtm.pub/

集成地址：https://dtm.pub/ref/gozero.html

同样我们先安装dtm

下载dtm源码，可以挑选一个版本

~~~shell
git clone https://github.com/dtm-labs/dtm.git
~~~

然后将conf.sample.yml改名为conf.yml

放入以下内容：

~~~yaml
MicroService:
 Driver: 'dtm-driver-gozero' # name of the driver to handle register/discover
 Target: 'etcd://localhost:2379/dtmservice' # register dtm server to this url
 EndPoint: 'localhost:36790'
~~~

然后运行

~~~shell
go run main.go -c conf.yml
~~~

创建表：

~~~sql
create database if not exists dtm_barrier
/*!40100 DEFAULT CHARACTER SET utf8mb4 */
;
drop table if exists dtm_barrier.barrier;
create table if not exists dtm_barrier.barrier(
  id bigint(22) PRIMARY KEY AUTO_INCREMENT,
  trans_type varchar(45) default '',
  gid varchar(128) default '',
  branch_id varchar(128) default '',
  op varchar(45) default '',
  barrier_id varchar(45) default '',
  reason varchar(45) default '' comment 'the branch type who insert this record',
  create_time datetime DEFAULT now(),
  update_time datetime DEFAULT now(),
  key(create_time),
  key(update_time),
  UNIQUE key(gid, branch_id, op, barrier_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
~~~



#### 创建积分服务

假设我们这里有个业务，在注册的时候，需要给用户增加积分，这个积分可以后续在系统的商城进行兑换商品。

我们创建一个独立的积分服务，注册成功后，进行积分服务调用，增加积分，如果积分增加失败，回滚。

sql

~~~sql
CREATE TABLE `user_score`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(0) NOT NULL,
  `score` int(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;
~~~

~~~go
go get github.com/dtm-labs/dtm
go get github.com/dtm-labs/driver-gozero
~~~

#### 实现

api中初始逻辑：

~~~go
func (l *UserLogic) Register(req *types.Request) (resp *types.Response, err error) {
	userResponse, err := l.svcCtx.UserRpc.SaveUser(context.Background(), &user.UserRequest{
		Name:   req.Name,
		Gender: req.Gender,
	})
	if err != nil {
		return nil, err
	}
	userId, _ := strconv.ParseInt(userResponse.Id, 10, 64)
	scoreRequest := &userscore.UserScoreRequest{
		UserId: userId,
		Score:  10,
	}
	score, err := l.svcCtx.UserScoreRpc.SaveUserScore(context.Background(), scoreRequest)
	if err != nil {
		return nil, err
	}
	logx.Infof("register add score %d", score.Score)
	return &types.Response{
		Message: "success",
		Data:    userResponse,
	}, nil
}

~~~

当添加积分失败时，无法回滚，所以我们加入dtm事务。

~~~go
// 下面这行导入gozero的dtm驱动
import _ "github.com/dtm-labs/driver-gozero"
// dtm已经通过前面的配置，注册到下面这个地址，因此在dtmgrpc中使用该地址
var dtmServer = "etcd://localhost:2379/dtmservice"
~~~



~~~go
func (l *UserLogic) Register(req *types.Request) (resp *types.Response, err error) {
	gid := dtmgrpc.MustGenGid(dtmServer)
	//消息型
	msgGrpc := dtmgrpc.NewSagaGrpc(dtmServer, gid)
	userRequest := &user.UserRequest{
		Name:   req.Name,
		Gender: req.Gender,
	}
	userServer, err := l.svcCtx.Config.UserRpc.BuildTarget()
	if err != nil {
		return nil, err
	}
	userScoreServer, err := l.svcCtx.Config.UserScoreRpc.BuildTarget()
	if err != nil {
		return nil, err
	}
	msgGrpc.Add(userServer+"/user.User/saveUser", userServer+"/user.User/saveUserCallback", userRequest)
	//userResponse, err := l.svcCtx.UserRpc.SaveUser(context.Background(), userRequest)
	//if err != nil {
	//	return nil, err
	//}
	//userId, _ := strconv.ParseInt(userResponse.Id, 10, 64)
	scoreRequest := &userscore.UserScoreRequest{
		UserId: 100,
		Score:  10,
	}
	msgGrpc.Add(userScoreServer+"/userscore.UserScore/saveUserScore", "", scoreRequest)
	//score, err := l.svcCtx.UserScoreRpc.SaveUserScore(context.Background(), scoreRequest)
	//if err != nil {
	//	return nil, err
	//}
	msgGrpc.WaitResult = true
	err = msgGrpc.Submit()
	if err != nil {
		fmt.Println("-----------------------")
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}
	//logx.Infof("register add score %d", score.Score)
	return &types.Response{
		Message: "success",
		Data:    "",
	}, nil
}
~~~

~~~go
func (l *GetUserLogic) SaveUser(in *user.UserRequest) (*user.UserResponse, error) {
	// todo: add your logic here and delete this line
	fmt.Println("------------------start-----------------")
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, err
	}
	data := &model.User{
		Name:   in.GetName(),
		Gender: in.GetGender(),
	}
	err = barrier.CallWithDB(l.svcCtx.Db, func(tx *sql.Tx) error {
		err := l.svcCtx.UserRepo.Save(context.Background(), tx, data)
		return err
	})
	if err != nil {
		fmt.Println(err)
		//Internal重试，Aborted 回滚
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	return &user.UserResponse{
		Id:     strconv.FormatInt(data.Id, 10),
		Name:   data.Name,
		Gender: data.Gender,
	}, nil
}
~~~

~~~go
func (l *GetUserLogic) SaveUserCallback(in *user.UserRequest) (*user.UserResponse, error) {
	fmt.Println("call back..........")
	return &user.UserResponse{}, nil
}
~~~





# 项目部署

## CI/CD
> 在软件工程中，CI/CD或CICD通常指的是持续集成和持续交付或持续部署的组合实践
从概念上来看，CI/CD包含部署过程，我们这里将部署(CD)单独放在一节服务部署， 本节就以gitlab来做简单的CI（Run Unit Test）演示。
## gitlab CI
Gitlab CI/CD是Gitlab内置的软件开发工具，提供
- 持续集成(CI)
- 持续交付(CD)
- 持续部署(CD)

### 准备工作
- gitlab安装
- git安装
- gitlab runner安装

### 开启gitlab CI
- 上传代码
  - 在gitlab新建一个仓库go-zero-demo
  - 将本地代码上传到go-zero-demo仓库
- 在项目根目录下创建.gitlab-ci.yaml文件，通过此文件可以创建一个pipeline，其会在代码仓库中有内容变更时运行，pipeline由一个或多个按照顺序运行， 每个阶段可以包含一个或者多个并行运行的job。
- 添加CI内容(仅供参考)
```bash
stages:
- analysis

analysis:
stage: analysis
image: golang
script:
- go version && go env
- go test -short $(go list ./...) | grep -v "no test"

```

## 日志收集
为了保证业务稳定运行，预测服务不健康风险，日志的收集可以帮助我们很好的观察当前服务的健康状况， 在传统业务开发中，机器部署还不是很多时，我们一般都是直接登录服务器进行日志查看、调试，但随着业务的增大，服务的不断拆分， 服务的维护成本也会随之变得越来越复杂，在分布式系统中，服务器机子增多，服务分布在不同的服务器上，当遇到问题时， 我们不能使用传统做法，登录到服务器进行日志排查和调试，这个复杂度可想而知。

### 准备工作
- kafka
- elasticsearch
- kibana
- filebeat、Log-Pilot（k8s）
- go-stash

### filebeat配置
前提说明: go服务日志是写入指定文件中
```bash
# filebeat.yaml
filebeat.inputs:
- type: log
  enabled: true
  # 开启json解析
  json.keys_under_root: true
  json.add_error_key: true
  # 日志文件路径
  paths:
    - /var/log/order/*.log

setup.template.settings:
  index.number_of_shards: 1

# 定义kafka topic field
fields:
  log_topic: log-collection

# 输出到kafka
output.kafka:
  hosts: ["127.0.0.1:9092"]
  topic: '%{[fields.log_topic]}'
  partition.round_robin:
    reachable_only: false
  required_acks: 1
  keep_alive: 10s

# ================================= Processors =================================
processors:
  - decode_json_fields:
      fields: ['@timestamp','level','content','trace','span','duration']
      target: ""

```

### go-stash配置
- 新建config.yaml文件
- 添加配置内容
```bash
# config.yaml
Clusters:
- Input:
    Kafka:
      Name: go-stash
      Log:
        Mode: file
      Brokers:
      - "127.0.0.1:9092"
      Topics: 
      - log-collection
      Group: stash
      Conns: 3
      Consumers: 10
      Processors: 60
      MinBytes: 1048576
      MaxBytes: 10485760
      Offset: first
  Filters:
  - Action: drop
    Conditions:
      - Key: status
        Value: "503"
        Type: contains
      - Key: type
        Value: "app"
        Type: match
        Op: and
  - Action: remove_field
    Fields:
    - source
    - _score
    - "@metadata"
    - agent
    - ecs
    - input
    - log
    - fields
  Output:
    ElasticSearch:
      Hosts:
      - "http://127.0.0.1:9200"
      Index: "go-stash-{{yyyy.MM.dd}}"
      MaxChunkBytes: 5242880
      GracePeriod: 10s
      Compress: false
      TimeZone: UTC

```

### 启动服务(按顺序启动)
- 启动kafka
- 启动elasticsearch
- 启动kibana
- 启动go-stash
- 启动filebeat
- 启动order-api服务及其依赖服务（go-zero-demo工程中的order-api服务）

### 访问kibana
进入127.0.0.1:5601

## 服务部署
本节通过jenkins来进行简单的服务部署到k8s演示
### 准备工作
- k8s集群安装
- gitlab环境安装
- jenkins环境安装
- redis&mysql&nginx&etcd安装
- goctl安装
> goctl确保k8s每个node节点上都有

### 服务部署
#### 1、gitlab代码仓库相关准备#
##### 1.1、添加SSH Key
进入gitlab，点击用户中心，找到Settings，在左侧找到SSH Keystab
```bash
1、在jenkins所在机器上查看公钥
$ cat ~/.ssh/id_rsa.pub

2、如果没有，则需要生成，如果存在，请跳转到第3步
$ ssh-keygen -t rsa -b 2048 -C "email@example.com"
"email@example.com" 可以替换为自己的邮箱
完成生成后，重复第一步操作

3、将公钥添加到gitlab中
```

##### 1.2、上传代码到gitlab仓库#
新建工程go-zero-demo并上传代码，这里不做细节描述。

#### 2、jenkins#
##### 2.1、添加凭据
查看jenkins所在机器的私钥，与前面gitlab公钥对应
```bash
cat id_rsa
```
- 进入jenkins，依次点击Manage Jenkins-> Manage Credentials
- 进入全局凭据页面，添加凭据，Username是一个标识，后面添加pipeline你知道这个标识是代表gitlab的凭据就行，Private Key`即上面获取的私钥

##### 2.2、 添加全局变量
进入Manage Jenkins->Configure System，滑动到全局属性条目，添加docker私有仓库相关信息，如图为docker用户名、docker用户密码、docker私有仓库地址

##### 2.3、配置git
进入Manage Jenkins->Global Tool Configureation，找到Git条目，填写jenkins所在机器git可执行文件所在path，如果没有的话，需要在jenkins插件管理中下载Git插件。

##### 2.4、 添加一个Pipeline
> pipeline用于构建项目，从gitlab拉取代码->生成Dockerfile->部署到k8s均在这个步骤去做，这里是演示环境，为了保证部署流程顺利， 需要将jenkins安装在和k8s集群的其中过一个节点所在机器上，我这里安装在master上的。
>

1. 获取凭据id 进入凭据页面，找到Username为gitlab的凭据id
2. 进入jenkins首页，点击新建Item，名称为user
3. 查看项目git地址
4. 添加服务类型Choice Parameter
   在General中勾选This project is parameterized ,点击添加参数选择Choice Parameter，按照图中添加选择的值常量(api、rpc)及接收值的变量(type)，后续在Pipeline script中会用到。
5. 配置user
   在user配置页面，向下滑动找到Pipeline script,填写脚本内容
```bash
pipeline {
  agent any
  parameters {
      gitParameter name: 'branch', 
      type: 'PT_BRANCH',
      branchFilter: 'origin/(.*)',
      defaultValue: 'master',
      selectedValue: 'DEFAULT',
      sortMode: 'ASCENDING_SMART',
      description: '选择需要构建的分支'
  }

  stages {
      stage('服务信息')    {
          steps {
              sh 'echo 分支：$branch'
              sh 'echo 构建服务类型：${JOB_NAME}-$type'
          }
      }


      stage('check out') {
          steps {
              checkout([$class: 'GitSCM', 
              branches: [[name: '$branch']],
              doGenerateSubmoduleConfigurations: false, 
              extensions: [], 
              submoduleCfg: [],
              # ${credentialsId} 要替换为你的具体凭据值, 即 [添加凭据] 模块中的一串字符串, ${gitUrl} 需要替换为你代码的git仓库地址，其他的 ${..} 形式的变量无需修改, 保持原样即可.
              userRemoteConfigs: [[credentialsId: '${credentialsId}', url: '${gitUrl}']]])
          }   
      }

      stage('获取commit_id') {
          steps {
              echo '获取commit_id'
              git credentialsId: '${credentialsId}', url: '${gitUrl}'
              script {
                  env.commit_id = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
              }
          }
      }


      stage('goctl版本检测') {
          steps{
              sh '/usr/local/bin/goctl -v'
          }
      }

      stage('Dockerfile Build') {
          steps{
                 sh '/usr/local/bin/goctl docker -go service/${JOB_NAME}/${type}/${JOB_NAME}.go'
                 script{
                     env.image = sh(returnStdout: true, script: 'echo ${JOB_NAME}-${type}:${commit_id}').trim()
                 }
                 sh 'echo 镜像名称：${image}'
                 sh 'docker build -t ${image} .'
          }
      }

      stage('上传到镜像仓库') {
          steps{
              sh '/root/dockerlogin.sh'
              sh 'docker tag  ${image} ${dockerServer}/${image}'
              sh 'docker push ${dockerServer}/${image}'
          }
      }

      stage('部署到k8s') {
          steps{
              script{
                  env.deployYaml = sh(returnStdout: true, script: 'echo ${JOB_NAME}-${type}-deploy.yaml').trim()
                  env.port=sh(returnStdout: true, script: '/root/port.sh ${JOB_NAME}-${type}').trim()
              }

              sh 'echo ${port}'

              sh 'rm -f ${deployYaml}'
              sh '/usr/local/bin/goctl kube deploy -secret dockersecret -replicas 2 -nodePort 3${port} -requestCpu 200 -requestMem 50 -limitCpu 300 -limitMem 100 -name ${JOB_NAME}-${type} -namespace hey-go-zero -image ${dockerServer}/${image} -o ${deployYaml} -port ${port}'
              sh '/usr/bin/kubectl apply -f ${deployYaml}'
          }
      }
      
      stage('Clean') {
          steps{
              sh 'docker rmi -f ${image}'
              sh 'docker rmi -f ${dockerServer}/${image}'
              cleanWs notFailBuild: true
          }
      }
  }
}

```
注意: 
- port.sh参考内容如下
```bash
case $1 in
    "user-api") 
        echo 1000
        ;;
    "user-rpc") 
        echo 1001
        ;;
    "course-api") 
        echo 1002
        ;;
    "course-rpc") 
        echo 1003
        ;;
    "selection-api") 
        echo 1004
esac

```

- dockerlogin.sh内容
```bash
#!/bin/bash
docker login --username=$docker-user --password=$docker-pass $docker-server

```

##### 2.5 测试运行
查看pipeline
查看k8s


## 服务监控
在微服务治理中，服务监控也是非常重要的一个环节，监控一个服务是否正常工作，需要从多维度进行，如：
- mysql指标
- mongo指标
- redis指标
- 请求日志
- 服务指标统计
- 服务健康检测 ...
监控的工作非常大，本节仅以其中的服务指标监控作为例子进行说明

### 基于prometheus的微服务指标监控
服务上线后我们往往需要对服务进行监控，以便能及早发现问题并做针对性的优化，监控又可分为多种形式，比如日志监控，调用链监控，指标监控等等。而通过指标监控能清晰的观察出服务指标的变化趋势，了解服务的运行状态，对于保证服务稳定起着非常重要的作用 prometheus是一个开源的系统监控和告警工具，支持强大的查询语言PromQL允许用户实时选择和汇聚时间序列数据，时间序列数据是服务端通过HTTP协议主动拉取获得，也可以通过中间网关来推送时间序列数据，可以通过静态配置文件或服务发现来获取监控目标

### Prometheus 的架构
Prometheus Server直接从监控目标中或者间接通过推送网关来拉取监控指标，它在本地存储所有抓取到样本数据，并对此数据执行一系列规则，以汇总和记录现有数据的新时间序列或生成告警。可以通过 Grafana 或者其他工具来实现监控数据的可视化

### go-zero基于prometheus的服务指标监控
go-zero 框架中集成了基于prometheus的服务指标监控，下面我们通过go-zero官方的示例shorturl来演示是如何对服务指标进行收集监控的：

- 第一步需要先安装Prometheus，安装步骤请参考官方文档
- go-zero默认不开启prometheus监控，开启方式很简单，只需要在shorturl-api.yaml文件中增加配置如下，其中Host为Prometheus Server地址为必填配置，Port端口不填默认9091，Path为用来拉取指标的路径默认为/metrics
```bash
Prometheus:
  Host: 127.0.0.1
  Port: 9091
  Path: /metrics

```
- 编辑prometheus的配置文件prometheus.yml，添加如下配置，并创建targets.json
```bash
- job_name: 'file_ds'
  file_sd_configs:
  - files:
    - targets.json

```
- 编辑targets.json文件，其中targets为shorturl配置的目标地址，并添加了几个默认的标签
```bash
[
    {
        "targets": ["127.0.0.1:9091"],
        "labels": {
            "job": "shorturl-api",
            "app": "shorturl-api",
            "env": "test",
            "instance": "127.0.0.1:8888"
        }
    }
]

```
- 启动prometheus服务，默认侦听在9090端口
```bash
prometheus --config.file=prometheus.yml
```
- 在浏览器输入http://127.0.0.1:9090/，然后点击Status -> Targets即可看到状态为Up的Job，并且Lables栏可以看到我们配置的默认的标签
通过以上几个步骤我们完成了prometheus对shorturl服务的指标监控收集的配置工作，为了演示简单我们进行了手动的配置，在实际的生产环境中一般采用定时更新配置文件或者服务发现的方式来配置监控目标，篇幅有限这里不展开讲解，感兴趣的同学请自行查看相关文档


### go-zero监控的指标类型
go-zero目前在http的中间件和rpc的拦截器中添加了对请求指标的监控。

主要从请求耗时和请求错误两个维度，请求耗时采用了Histogram指标类型定义了多个Buckets方便进行分位统计，请求错误采用了Counter类型，并在http metric中添加了path标签rpc metric中添加了method标签以便进行细分监控。 接下来演示如何查看监控指标： 首先在命令行多次执行如下命令
```bash
curl -i "http://localhost:8888/shorten?url=http://www.xiaoheiban.cn"

```
打开Prometheus切换到Graph界面，在输入框中输入{path="/shorten"}指令，即可查看监控指标

我们通过PromQL语法查询过滤path为/shorten的指标，结果中显示了指标名以及指标数值，其中http_server_requests_code_total指标中code值为http的状态码，200表明请求成功，http_server_requests_duration_ms_bucket中对不同bucket结果分别进行了统计，还可以看到所有的指标中都添加了我们配置的默认指标 Console界面主要展示了查询的指标结果，Graph界面为我们提供了简单的图形化的展示界面，在实际的生产环境中我们一般使用Grafana做图形化的展示

### grafana可视化界面
grafana是一款可视化工具，功能强大，支持多种数据来源Prometheus、Elasticsearch、Graphite等，安装比较简单请参考官方文档，grafana默认端口3000，安装好后再浏览器输入http://localhost:3000/，默认账号和密码都为admin 下面演示如何基于以上指标进行可视化界面的绘制： 点击左侧边栏Configuration->Data Source->Add data source进行数据源添加，其中HTTP的URL为数据源的地址

### 总结
以上演示了go-zero中基于prometheus+grafana服务指标监控的简单流程，生产环境中可以根据实际的场景做不同维度的监控分析。现在go-zero的监控指标主要还是针对http和rpc，这对于服务的整体监控显然还是不足的，比如容器资源的监控，依赖的mysql、redis等资源的监控，以及自定义的指标监控等等，go-zero在这方面后续还会持续优化。希望这篇文章能够给您带来帮助


## go-zero链路追踪
### 序言
微服务架构中，调用链可能很漫长，从 http 到 rpc ，又从 rpc 到 http 。而开发者想了解每个环节的调用情况及性能，最佳方案就是 全链路跟踪。

追踪的方法就是在一个请求开始时生成一个自己的 spanID ，随着整个请求链路传下去。我们则通过这个 spanID 查看整个链路的情况和性能问题。

下面来看看 go-zero 的链路实现。

### 代码结构
- spancontext ：保存链路的上下文信息「traceid，spanid，或者是其他想要传递的内容」
- span ：链路中的一个操作，存储时间和某些信息
- propagator ： trace 传播下游的操作「抽取，注入」
- noop ：实现了空的 tracer 实现

### 概念
#### SpanContext
在介绍 span 之前，先引入 context 。SpanContext 保存了分布式追踪的上下文信息，包括 Trace id，Span id 以及其它需要传递到下游的内容。OpenTracing 的实现需要将 SpanContext 通过某种协议 进行传递，以将不同进程中的 Span 关联到同一个 Trace 上。对于 HTTP 请求来说，SpanContext 一般是采用 HTTP header 进行传递的。

下面是 go-zero 默认实现的 spanContext
```bash
type spanContext struct {
    traceId string      // TraceID 表示tracer的全局唯一ID
    spanId  string      // SpanId 标示单个trace中某一个span的唯一ID，在trace中唯一
}

```
同时开发者也可以实现 SpanContext 提供的接口方法，实现自己的上下文信息传递：
```bash
type SpanContext interface {
    TraceId() string                        // get TraceId
    SpanId() string                         // get SpanId
    Visit(fn func(key, val string) bool)    // 自定义操作TraceId，SpanId
}

```
#### Span
一个 REST 调用或者数据库操作等，都可以作为一个 span 。 span 是分布式追踪的最小跟踪单位，一个 Trace 由多段 Span 组成。追踪信息包含如下信息：
```bash

type Span struct {
    ctx           spanContext       // 传递的上下文
    serviceName   string            // 服务名 
    operationName string            // 操作
    startTime     time.Time         // 开始时间戳
    flag          string            // 标记开启trace是 server 还是 client
    children      int               // 本 span fork出来的 childsnums
}

```
从 span 的定义结构来看：在微服务中， 这就是一个完整的子调用过程，有调用开始 startTime ，有标记自己唯一属性的上下文结构 spanContext 以及 fork 的子节点数。

### 实例应用
在 go-zero 中http，rpc中已经作为内置中间件集成。我们以 http ，rpc 中，看看 tracing 是怎么使用的
#### HTTP
```bash
func TracingHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // **1**
        carrier, err := trace.Extract(trace.HttpFormat, r.Header)
        // ErrInvalidCarrier means no trace id was set in http header
        if err != nil && err != trace.ErrInvalidCarrier {
            logx.Error(err)
        }

        // **2**
        ctx, span := trace.StartServerSpan(r.Context(), carrier, sysx.Hostname(), r.RequestURI)
        defer span.Finish()
        // **5**
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}

func StartServerSpan(ctx context.Context, carrier Carrier, serviceName, operationName string) (
    context.Context, tracespec.Trace) {
    span := newServerSpan(carrier, serviceName, operationName)
    // **4**
    return context.WithValue(ctx, tracespec.TracingKey, span), span
}

func newServerSpan(carrier Carrier, serviceName, operationName string) tracespec.Trace {
    // **3**
    traceId := stringx.TakeWithPriority(func() string {
        if carrier != nil {
            return carrier.Get(traceIdKey)
        }
        return ""
    }, func() string {
        return stringx.RandId()
    })
    spanId := stringx.TakeWithPriority(func() string {
        if carrier != nil {
            return carrier.Get(spanIdKey)
        }
        return ""
    }, func() string {
        return initSpanId
    })

    return &Span{
        ctx: spanContext{
            traceId: traceId,
            spanId:  spanId,
        },
        serviceName:   serviceName,
        operationName: operationName,
        startTime:     timex.Time(),
        // 标记为server
        flag:          serverFlag,
    }
}

```
1. 将 header -> carrier，获取 header 中的traceId等信息
2. 开启一个新的 span，并把「traceId，spanId」封装在context中
3. 从上述的 carrier「也就是header」获取traceId，spanId
  1. 看header中是否设置
  2. 如果没有设置，则随机生成返回
4. 从 request 中产生新的ctx，并将相应的信息封装在 ctx 中，返回
5. 从上述的 context，拷贝一份到当前的 request

这样就实现了 span 的信息随着 request 传递到下游服务。


#### RPC
在 rpc 中存在 client, server ，所以从 tracing 上也有 clientTracing, serverTracing 。 serveTracing 的逻辑基本与 http 的一致，来看看 clientTracing 是怎么使用的？
```bash
func TracingInterceptor(ctx context.Context, method string, req, reply interface{},
    cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
    // open clientSpan
    ctx, span := trace.StartClientSpan(ctx, cc.Target(), method)
    defer span.Finish()

    var pairs []string
    span.Visit(func(key, val string) bool {
        pairs = append(pairs, key, val)
        return true
    })
    // **3** 将 pair 中的data以map的形式加入 ctx
    ctx = metadata.AppendToOutgoingContext(ctx, pairs...)

    return invoker(ctx, method, req, reply, cc, opts...)
}

func StartClientSpan(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
    // **1**
    if span, ok := ctx.Value(tracespec.TracingKey).(*Span); ok {
        // **2**
        return span.Fork(ctx, serviceName, operationName)
    }

    return ctx, emptyNoopSpan
}
```
1. 获取上游带下来的 span 上下文信息
2. 从获取的 span 中创建新的 ctx，span「继承父span的traceId」
3. 将生成 span 的data加入ctx，传递到下一个中间件，流至下游

### 总结
go-zero 通过拦截请求获取链路traceID，然后在中间件函数入口会分配一个根Span，然后在后续操作中会分裂出子Span，每个span都有自己的具体的标识，Finsh之后就会汇集在链路追踪系统中。开发者可以通过 ELK 工具追踪 traceID ，看到整个调用链。

同时 go-zero 并没有提供整套 trace 链路方案，开发者可以封装 go-zero 已有的 span 结构，做自己的上报系统，接入 jaeger, zipkin 等链路追踪工具。



# go-zero goctl

## 基本介绍

goctl是go-zero微服务框架下的代码生成工具。使用 goctl 可显著提升开发效率，让开发人员将时间重点放在业务开发上，其功能有：

- api服务生成
- rpc服务生成
- model代码生成
- 模板管理

## 安装 goctl

```bash

# Go 1.16 及以后版本
GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@latest
```

自动补全设置(仅支持linux系统)

```bash

# 生成自动补全文件
$ goctl completion
generation auto completion success!
executes the following script to setting shell:
echo PROG=goctl source /Users/keson/.goctl/.auto_complete/zsh/goctl_autocomplete >> ~/.zshrc && source ~/.zshrc
or
echo PROG=goctl source /Users/keson/.goctl/.auto_complete/bash/goctl_autocomplete >> ~/.bashrc && source ~/.bashrc

# shell 配置
1、zsh
$ echo PROG=goctl source /Users/keson/.goctl/.auto_complete/zsh/goctl_autocomplete >> ~/.zshrc && source ~/.zshrc

2、bash
$ echo PROG=goctl source /Users/keson/.goctl/.auto_complete/bash/goctl_autocomplete >> ~/.bashrc && source ~/.bashrc

```

## 常见命令
### help命令
```bash
goctl --help
```

### env
```bash
goctl env

# 修改变量
goctl env -w 
```

### 快速启动服务
- 单体服务
```bash
mkdir -p demo && cd demo
go mod init demo
goctl quickstart
```
- 微服务
```bash
goctl quickstart -t micro
# 包含api和rpc服务
```

### upgrade
```bash
# goctl 升级
goctl upgrade
```

### api命令

goctl api是goctl中的核心模块之一，其可以通过.api文件一键快速生成一个api服务，如果仅仅是启动一个go-zero的api演示项目， 你甚至都不用编码，就可以完成一个api服务开发及正常运行。在传统的api项目中，我们要创建各级目录，编写结构体， 定义路由，添加logic文件，这一系列操作，如果按照一条协议的业务需求计算，整个编码下来大概需要5～6分钟才能真正进入业务逻辑的编写， 这还不考虑编写过程中可能产生的各种错误，而随着服务的增多，随着协议的增多，这部分准备工作的时间将成正比上升， 而goctl api则可以完全替代你去做这一部分工作，不管你的协议要定多少个，最终来说，只需要花费10秒不到即可完成。

```bash

方式一: 快速生成api服务文件
goctl api new user

方式二: 根据api文件生成go文件
# 生成api文件模板
goctl api -o user.api

# 根据api文件生成go文件
goctl api go -api user.api -dir=. --style gozero


# 文件校验
goctl api validate --api user.api

# 文件格式化
goctl api format --dir=. 

# 输出mk文档, 需要@doc标识
service test-api {
	@doc "get 路由"  
	@handler GetUser // TODO: set handler name and delete this comment
	get /users/id/:userId(request) returns(response)

	@doc(
		summary: "post 路由"
	)
	
	@handler CreateUser // TODO: set handler name and delete this comment
	post /users/create(request)
}


# 返回api文件上一级目录
goctl api doc --dir . --o .
```


### rpc命令

Goctl Rpc是goctl脚手架下的一个rpc服务代码生成模块，支持proto模板生成和rpc服务代码生成，通过此工具生成代码你只需要关注业务逻辑编写而不用去编写一些重复性的代码。这使得我们把精力重心放在业务上，从而加快了开发效率且降低了代码出错率。

特性

- 简单易用
- 快速提升开发效率
- 出错率低
- 贴近protoc

```bash

方式一: 快速生成greet服务
goctl rpc new greet

方式二: 通过指定proto生成rpc服务
# 生成proto模板
goctl rpc template -o=user.proto
# 通过指定proto生成rpc服务
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

### model命令

goctl model 为go-zero下的工具模块中的组件之一，目前支持识别mysql ddl进行model层代码生成，通过命令行或者idea插件（即将支持）可以有选择地生成带redis cache或者不带redis cache的代码逻辑。

```bash
方式一: 通过ddl生成
goctl model mysql ddl -src="./*.sql" -dir="./sql/model" -c

方式二: 通过datasource生成
goctl model mysql datasource -url="user:password@tcp(127.0.0.1:3306)/database" -table="*"  -dir="./model"

```

### plugin命令

goctl支持针对api自定义插件，那我怎么来自定义一个插件了？来看看下面最终怎么使用的一个例子。

```bash

goctl api plugin -p goctl-android="android -package com.tal" -api user.api -dir .

上面这个命令可以分解成如下几步：
    goctl 解析api文件
    goctl 将解析后的结构 ApiSpec 和参数传递给goctl-android可执行文件
    goctl-android 根据 ApiSpec 结构体自定义生成逻辑。
此命令前面部分 goctl api plugin -p 是固定参数，goctl-android="android -package com.tal" 是plugin参数，其中goctl-android是插件二进制文件，android -package com.tal是插件的自定义参数，-api user.api -dir .是goctl通用自定义参数。

```

### template
```bash
# 模板相关命令
$ goctl template -h
Template operation

Usage:
  goctl template [command]

Available Commands:
  clean       Clean the all cache templates
  init        Initialize the all templates(force update)
  revert      Revert the target template to the latest
  update      Update template of the target category to the latest

Flags:
  -h, --help   help for template

Use "goctl template [command] --help" for more information about a command.

# 备份模板
mkdir -p template
goctl template init --home template/
# 生成api服务指定模板文件
goctl api -o user.api --home template/
```

### 路由前缀/分组/tag
```bash
# api文件
syntax = "v1"

info (
	title: // TODO: add title
	desc: // TODO: add description
	author: "hu417"
	email: "hu729919300@163.com"
)

type request {
	// TODO: add members here and delete this comment
    ping string `json:"ping"` // 
    /*
    json:"name,optional" // json是指tag绑定body中json格式参数,optional可选参数
    path:"id"           // path是指tag绑定uri中的参数
    form:"name"         // form是指tag绑定form表单格式参数

    */
}

type response {
	// TODO: add members here and delete this comment
}


@server (
    prefix: /v1
    group: ping
)

service greet {
    @handler ping
    get /ping
}

@server(
    prefix: /v1
    group: login
)

service greet {
    @handler login
    get /user/login (LoginReq) returns (LoginResp)
}
```

### 其他命令

#### goctl docker

goctl docker 可以极速生成一个 Dockerfile，帮助开发/运维人员加快部署节奏，降低部署复杂度。

```bash
# 查看帮助
goctl docker -h

# 生成默认Dockerfile
goctl docker -go hello.go
```

Dockerfile 内容如下：

```bash

  FROM golang:alpine AS builder
  LABEL stage=gobuilder
  ENV CGO_ENABLED 0
  ENV GOOS linux
  ENV GOPROXY https://goproxy.cn,direct
  WORKDIR /build/zero
  ADD go.mod .
  ADD go.sum .
  RUN go mod download
  COPY . .
  COPY service/hello/etc /app/etc
  RUN go build -ldflags="-s -w" -o /app/hello service/hello/hello.go
  FROM alpine
  RUN apk update --no-cache
  RUN apk add --no-cache ca-certificates
  RUN apk add --no-cache tzdata
  ENV TZ Asia/Shanghai
  WORKDIR /app
  COPY --from=builder /app/hello /app/hello
  COPY --from=builder /app/etc /app/etc
  CMD ["./hello", "-f", "etc/hello-api.yaml"]
```

镜像构建

```bash
$ docker build -t hello:v1 -f service/hello/Dockerfile .

```

#### goctl kube

goctl kube提供了快速生成一个 k8s 部署文件的功能，可以加快开发/运维人员的部署进度，减少部署复杂度。

```bash

goctl kube deploy -name redis -namespace adhoc -image redis:6-alpine -o redis.yaml -port 6379

# 其他参数
goctl kube deploy -h
```

服务部署

```bash

 kubectl create namespace adhoc 
 kubectl apply -f redis.yaml
```


# 其他
- api文件
```bash
# api文件

@server(
	jwt: Auth
	middleware: Example // 路由中间件声明

    group: user  // 业务/中间件文件分组
    prefix: v1   // 路由前缀
)

service search-api {
    @doc "获取jwt token"
	@handler search
	get /search/do (SearchReq) returns (SearchReply)
}

说明: 
    @server()的内容会对最近的service{}生效
    service xxx {} 可以写多个，但是xxx要一样


```
- model控制缓存时间
在 sqlc.NewNodeConn 的时候可以通过可选参数 cache.WithExpiry 传递，如缓存时间控制为1天，代码如下
```bash
sqlc.NewNodeConn(conn,redis,cache.WithExpiry(24*time.Hour))  
```

- rpc直连与服务发现连接模式写法
```bash
// mode1: 集群直连
// conf:=zrpc.NewDirectClientConf([]string{"ip:port"},"app","token")

// mode2: etcd 服务发现
// conf:=zrpc.NewEtcdClientConf([]string{"ip:port"},"key","app","token")
// client, _ := zrpc.NewClient(conf)

// mode3: ip直连mode
// client, _ := zrpc.NewClientWithTarget("127.0.0.1:8888")
```
- 跨域
```bash
srv := rest.MustNewServer(c, rest.WithCors())

# 单个域名的情况
srv := rest.MustNewServer(c, rest.WithCors("http://example.com"))

```

- yaml配置文件中使用环境变量
```bash
    conf.MustLoad(*configFile, &c, conf.UseEnv())

```



