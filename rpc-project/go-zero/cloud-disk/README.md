# 轻量级云盘系统
> 轻量级云盘系统，基于 go-zero，xxorm 实现
功能模块:
- 用户模块
  - 密码登录
  - 邮箱注册
  - 个人资料
- 存储池模块
  - 中心存储池资源管理(公共资源)
    - 文件存储
  - 个人存储池资源管理
    - 文件关联存储
    - 文件夹层级管理
- 文件分享模块
  - 文件分享

涉及的技术栈:
- go
```bash
$ go version
go version go1.19.3 windows/amd64
```
- go-zero/goctl
```bash
# Go 1.16 及以后版本
$ GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@latest
$ goctl -v
goctl version 1.5.0 windows/amd64
```
- mysql
```bash
mkdir -p /app/mysql/{conf,data,log}
vim /app/mysql/my.cnf 
[mysqld]
user=mysql
character-set-server=utf8
default_authentication_plugin=mysql_native_password
secure_file_priv=/var/lib/mysql
expire_logs_days=7
sql_mode=STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
max_connections=1000
 
[client]
default-character-set=utf8
 
[mysql]
default-character-set=utf8

chmod -R 777 /app/mysql/
docker run --restart=always --privileged=true \
    -v /app/mysql/data/:/var/lib/mysql \
    -v /app/mysql/log/:/var/log/mysql \
    -v /app/mysql/conf/my.cnf:/etc/my.cnf \
    -p 3306:3306 --name mysql \
    -e MYSQL_ROOT_PASSWORD=123456 \
    -e LANG=C.utf8 \
    -d mysql:8.0.32
```
- redis
```bash
mkdir -p /app/redis/{conf,data,log}
vim /app/redis/conf/redis.conf
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

chmod -R 777 /app/redis
docker run -d --name redis \
    --restart=always \
    -p 6379:6379 \
    -v /app/redis/conf/redis.conf:/opt/redis.conf \
    -v /app/redis/log/redis.log:/var/log/redis.log \
    -v /app/redis/data:/opt/data
    redis:6.2.7 \

```

## 项目
### 初始化
```bash
# 项目目录
mkdir -p cloud-disk && cd cloud-disk
go mod init cloud-disk

```

### 数据库设计
```bash
mkdir -p sql && cd sql

# 创建数据库
$ mysql -uroot -p123456
mysql> create database `cloud-disk`;
mysql> use cloud-di
```
- 用户信息
存储用户基本信息，用于登录
```bash
# user_basic.sql
CREATE TABLE `user_basic` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`name` varchar(60) DEFAULT NULL,
	`password` varchar(32) DEFAULT NULL,
	`email` varchar(100) DEFAULT NULL,

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

```
- 中心存储池
公共存储文件信息
```bash
# repository_pool.sql
CREATE TABLE `repository_pool` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`hash` varchar(32) DEFAULT NULL COMMENT '文件的唯一标识',
	`name` varchar(255) DEFAULT NULL COMMENT '文件名称',
	`ext` varchar(30) DEFAULT NULL COMMENT '文件扩展名',
	`size` int(11) DEFAULT NULL COMMENT '文件大小',
	`path` varchar(255) DEFAULT NULL COMMENT '文件路径',

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

```
- 用户存储池
对公共文件存储池中文件信息的引用
```bash
# user_repository.sql
CREATE TABLE `user_repository` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`parent_id` int(11) DEFAULT NULL COMMENT '父级文件层级, 0-【文件夹】',
	`user_identity` varchar(36) DEFAULT NULL COMMENT '对应用户的唯一标识',
	`repository_identity` varchar(36) DEFAULT NULL COMMENT '公共池中文件的唯一标识',
	`ext` varchar(255) DEFAULT NULL COMMENT '文件或文件夹类型',
	`name` varchar(255) DEFAULT NULL COMMENT '用户定义的文件名',

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
```
- 文件分享
```bash
# share_basic.sql
CREATE TABLE `share_basic` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`user_identity` varchar(36) DEFAULT NULL COMMENT '对应用户的唯一标识',
	`repository_identity` varchar(36) DEFAULT NULL COMMENT '公共池中文件的唯一标识',
	`user_repository_identity` varchar(36) DEFAULT NULL COMMENT '用户池子中的唯一标识',
	`expired_time` int(11) DEFAULT NULL COMMENT '失效时间，单位秒,【0-永不失效】',
	`click_num` int(11) DEFAULT '0' COMMENT '点击次数',

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

```

### xorm
官方地址: https://xorm.io/zh/
- 安装
```bash
go get xorm.io/xorm # 安装xorm
go get xorm.io/builder # 安装builder
go get xorm.io/reverse # 安装反转工具


```
- 添加测试数据
```bash
insert into `user_basic`(identity,name,password,email) value("ident","name","password","123@admin.com");
```

- 测试xorm

```bash
mkdir -p test models

# models/user_basic.go
package models

// 字段映射
type User_basic struct {
	Id       int
	Identity string
	Name     string
	Password string
	Email    string
}

// 表明初始化
func (table *User_basic) TableName() string {
	return "user_basic"
}

```

```bash
# test/xorm_test.go
package test

import (
	"bytes"
	"cloud-disk/models"
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func TestXorm(t *testing.T) {
	// 1、建立连接
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(10.0.0.91:3306)/cloud-disk?charset=utf8")
	// 开启显示sql语句
	engine.ShowSQL(true)

	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	// 2、执行curd
	data := make([]*models.User_basic, 0)
	err = engine.Find(&data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("data => ", data) // [0xc0001942d0]
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	dst := new(bytes.Buffer)
	err = json.Indent(dst, b, "", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("dst => ", dst.String())
	fmt.Printf("data 类型: %T", dst)

}

```
- 执行单元测试
```bash
# 日志
Running tool: D:\study\go\bin\go.exe test -timeout 30s -run ^TestXorm$ cloud-disk/test

=== RUN   TestXorm
data =>  [0xc0001922d0]
dst =>  [
{
"Id": 5,
"Identity": "ident",
"Name": "name",
"Password": "password",
"Email": "123@admin.com"
}
]
data 类型: *bytes.Buffer--- PASS: TestXorm (0.01s)
PASS
ok      cloud-disk/test 1.203s


> 测试运行完成时间: 2023/3/19 23:53:54 <
```

### go-zero
#### 初始化项目
```bash
# cloud-disk
goctl api new core
go mod tidy

// 业务逻辑
# internal/logic/corelogic.go
...
func (l *CoreLogic) Core(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return &types.Response{
		Message: "hello world",
	}, nil
}

// 启动服务
cd core
go run core.go -f etc/core-api.yaml


// 测试，根据api文件里接口信息
get http://localhost:8888/from/you

```

#### 集成xorm
- 配置xorm
```bash
# core/internal/models/
# 1、init_db.go
package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"xorm.io/xorm"
)

var Engine = Init()

func Init() *xorm.Engine {
	// 建立连接
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(10.0.0.91:3306)/cloud-disk?charset=utf8")
	engine.ShowSQL(true)

	if err != nil {
		logx.Error("Xorm Engine Error: ", err)
		return nil
	}
	return engine
}


# 2、user_basic.go
package models

// 字段映射
type User_basic struct {
	Id       int
	Identity string
	Name     string
	Password string
	Email    string
}

// 表明初始化
func (table *User_basic) TableName() string {
	return "user_basic"
}

```

- 业务逻辑
```bash
# core/internal/logic/corelogic.go
import (
	"bytes"
	"context"
	"encoding/json"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *CoreLogic) Core(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	data := make([]*models.User_basic, 0)
	err = models.Engine.Find(&data)
	if err != nil {
		logx.Infof("engine find error: %v", err)
	}

	logx.Info("data => ", data) // [0xc0001942d0]
	b, err := json.Marshal(data)
	if err != nil {
		logx.Infof("json Marshal error: %v", err)
	}
	dst := new(bytes.Buffer)
	err = json.Indent(dst, b, "", "")
	if err != nil {
		logx.Infof("json Indent error: %v", err)
	}
	logx.Info("dst => ", dst.String())

	resp = new(types.Response)
	resp.Message = dst.String()

	return
	// return &types.Response{
	// 	Message: dst.String(),
	// }, nil
}
```


- 日志配置
```bash
# core/core.go
import (

	"github.com/zeromicro/go-zero/core/logx"

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

	// 关闭stat日志
	logx.DisableStat()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

```

- 启动服务
```bash
# core
go run core.go -f etc/core-api.yaml


# 发送请求
get http://localhost:8888/from/you

```

#### 用户模块

##### 密码登录

请求逻辑: 前端发送请求(参数:用户名/密码) -> 根据请求参数中password进行加密 -> 将加密后的password,name从数据库中查询 -> 将数据库返回的数据进行jwt token生成并加盐 -> 将生成的token返回前端
- api文件
```bash
# core.api
// 用户登录
service core-api {
	@doc(
		summary: "用户登录"
	)
	@handler UserLogin
	post /user/login(LoginRequest) returns (LoginResponse)
	
}

// 用户登录
type LoginRequest {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse {
	Token string `json:"token"`
}


# 生成go文件
goctl api go --api core.api --dir=. --style go_zero


// 删除相关文件: core_logic.go, core_handler.go, servicecontext.go
```

- logic
1. 插入数据
```bash
echo -n "123456" | md5sum | cut -d " " -f1
// e10adc3949ba59abbe56e057f20f883e

mysql> insert into user_basic(`identity`,`name`,`password`,`email`,`created_at`,`updated_at`) values("qazwsx","李四","e10adc3949ba59abbe56e057f20f883e","123456@admin.com","2023-03-20 21:56","2023-03-20 21:57");

```
2. md5加密
helper方法
```bash
# cloud-disk/core/helper/helper.go
package helper

import (
	"crypto/md5"
	"fmt"

	"cloud-disk/core/define"

	"github.com/dgrijalva/jwt-go"
)

func Md5(s string) string {

	// md5加密
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

```
3. 查询数据库
user_login_logic.go
```bash
# core\internal\logic\user_login_logic.go
import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserLogic) User(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line

	// 1、从数据库中查询用户
	user := new(models.User_basic)
	has, err := models.Engine.Where("name = ? AND password = ?", req.Name, helper.Md5(req.Passwprd)).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户名/密码错误")
	}
	
	return
}

```
4. jwt token生成
define.go
```bash
# core\define\define.go

package define

// jwt
import "github.com/dgrijalva/jwt-go"

// 定义用户 jwt token相关参数
type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.StandardClaims
}

// 加盐
var Shar = "cloud-disk"

```

helper.go
```bash
# core\helper\helper.go

package helper

import (
	"crypto/md5"
	"fmt"

	"cloud-disk/core/define"

	"github.com/dgrijalva/jwt-go"
)

func Md5(s string) string {

	// md5加密
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int, identity, name string) (string, error) {

	// jwt token生成
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
	}
	// jwt token加密
	token := jwt.NewWithClaims(jwt.SigningMethodES256, uc)
	// token加盐
	tikenstring, err := token.SignedString([]byte(define.Shar))
	if err != nil {
		return "", err
	}
	return tikenstring, nil

}

```
user_login_logic.go
```bash
# core\internal\logic\user_login_logic.go
import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserLogic) User(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line

	// 1、从数据库中查询用户
	user := new(models.User_basic)
	has, err := models.Engine.Where("name = ? AND password = ?", req.Name, helper.Md5(req.Passwprd)).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户名/密码错误")
	}
	// 2、生产token
	token, err := helper.GenerateToken(user.Id, user.Identity, user.Name)

	if err != nil {
		return nil, err
	}

	resp = new(types.LoginResponse)
	resp.Token = token
	return
}



```



5. 启动服务
```bash
go mod tidy
go run core.go -f etc/core-api.yaml

# 测试请求
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'

# 响应
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8"
}

```
##### 刷新Token
两个 Token：
	Token 有效期较短，用于对用户做鉴权操作
	RefreshToken 有效期比较长，用于刷新上面那个 Token
- api
```bash
# core.api
service core-api {
	@doc(
		summary: "用户登录"
	)
	@handler UserLogin
	post /user/login(LoginRequest) returns (LoginResponse)
}

@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	@doc(
		summary: "刷新token过期时间"
	)
	@handler UserRefreshToken
	post /user/refresh/token(UserRefreshTokenReq) returns ( UserRefreshTokenResp)
	
}

// 用户登录
type LoginRequest {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// 用户-刷新token
type (
	UserRefreshTokenReq {}
	UserRefreshTokenResp {
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		Message string `json:"message"`
	}
)
```
- define
```bash
# core\define\define.go

// 定义token有效期
var (
	TokenExpires        = 3600
	RefreshTokenExpires = 24 * 60 * 60
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- helper
```bash
# core\helper\helper.go
// jwt token生成
func GenerateToken(id int, identity, name string, second int) (string, error) {

	// jwt token生成
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
		// 添加过期时间，当前时间+有效期时间=最终过期时间
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(second)).Unix(),
		},
	}
	// jwt token加密
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	// token加盐
	tikenstring, err := token.SignedString([]byte(define.Shar))
	if err != nil {
		return "", err
	}
	return tikenstring, nil

}

```
- handler
```bash
# core\internal\handler\user_refresh_token_handler.go
func UserRefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRefreshTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserRefreshTokenLogic(r.Context(), svcCtx)

		// 传递用户token
		resp, err := l.UserRefreshToken(&req, r.Header.Get("Authorization"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

```
- logic
```bash
# core\internal\logic\user_login_logic.go
import (
	"context"
	"strconv"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserRefreshTokenLogic) UserRefreshToken(req *types.UserRefreshTokenReq, authorization string) (resp *types.UserRefreshTokenResp, err error) {
	// todo: add your logic here and delete this line

	// 根据用户token进行解析，并重新生成
	uc, err := helper.AnalyzeJwtTkoen(authorization)
	if err != nil {
		return
	}

	token, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.TokenExpires)
	if err != nil {
		return
	}
	refreshToken, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.RefreshTokenExpires)

	resp = new(types.UserRefreshTokenResp)
	resp.Token = token
	resp.RefreshToken = refreshToken
	resp.Message = "token 过期时间: " + strconv.Itoa(define.TokenExpires) + "s"

	return
}

```
- test
```bash
1、修改用户登录时默认token过期时间
func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line
...
	// 2、生成token
	//token, err := helper.GenerateToken(user.Id, user.Identity, user.Name, define.TokenExpires)
	token, err := helper.GenerateToken(user.Id, user.Identity, user.Name, 20)
...

}

2、启动服务
go run core.go -f etc/core-api.yaml

3、请求登录接口获取token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'

4、使用token请求refreshtoken接口
// 等待20S
curl --location --request POST 'http://localhost:8888/user/refresh/token' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIiwiZXhwIjoxNjc5ODExOTY5fQ.swWM3Iet2gOCuXKhB6r_o0XvXL8hnOHeMwT1Dm0MoYk' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw ''
// 返回信息
认证失败  # 此提示是: middleware\auth_middleware.go

5、使用refreshToken请求refreshtoken接口
curl --location --request POST 'http://localhost:8888/user/refresh/token' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIiwiZXhwIjoxNjc5ODk4MzQ5fQ.YLpF2qWQzYQbHHykfCIHGaEB-LJTO_NSsnsCdaw-dPw' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw ''

// 返回值
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIiwiZXhwIjoxNjc5ODE2ODAxfQ.DH_wsKRmkiKP3_YMQpdBlkhVUDlBiQx4bu9rh3S466c",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIiwiZXhwIjoxNjc5ODk5NjAxfQ.mLAUpiL_cOAInLr1qkChIIr0-FlV2DpebOk4wGTID9Y",
    "message": "token 过期时间: 3600s"
}


```

##### 用户详情
- api文件
```bash
# core\core.api
// 用户登录
service core-api {
	@doc(
		summary: "用户登录"
	)
	@handler UserLogin
	post /user/login(LoginRequest) returns (LoginResponse)
	
	@doc(
		summary: "用户详情"
	)
	@handler UserDetails
	get /user/details(UserDetailsReq) returns (UserDetailsResp)
	
}

// 用户登录
type LoginRequest {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse {
	Token string `json:"token"`
}

// 用户详情
type (
	UserDetailsReq {
		Identity string `json:"username"`
	}

	UserDetailsResp {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style go_zero

```
- 业务逻辑
```bash
# core\internal\logic\user_details_logic.go

import (
	"context"
	"errors"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

...
func (l *UserDetailsLogic) UserDetails(req *types.UserDetailsReq) (resp *types.UserDetailsResp, err error) {
	// todo: add your logic here and delete this line

	resp = &types.UserDetailsResp{}
	ub := new(models.User_basic)
	has, err := models.Engine.Where("identity = ?", req.Identity).Get(ub)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("identity 错误")
	}
	resp.Name = ub.Name
	resp.Email = ub.Email
	return
}


```



- 运行服务
```bash
go run core.go -f etc/core-api.yaml



# 测试请求
curl --location --request GET 'http://localhost:8888/user/details' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "identity": "qazwsx"
}'

```

##### 邮箱注册
- 安装email
```bash
go get github.com/jordan-wright/email

```

- api文件
```bash
# core/core.api
// 用户登录
service core-api {
	@doc(
		summary: "用户登录"
	)
	@handler UserLogin
	post /user/login(LoginRequest) returns (LoginResponse)
	
	@doc(
		summary: "用户详情"
	)
	@handler UserDetails
	get /user/details(UserDetailsReq) returns (UserDetailsResp)
	
	@doc(
		summary: "邮箱验证码"
	)
	@handler MailCodeSend
	post /mail/send(MailCodeReq) returns (MailCodeResp)
	
}

// 用户登录
type LoginRequest {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse {
	Token string `json:"token"`
}

// 用户详情
type (
	UserDetailsReq {
		Identity string `json:"identity"`
	}

	UserDetailsResp {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

// 邮箱验证码

type (
	MailCodeReq {
		Email string `json:"email"`
	}

	MailCodeResp {
        Msg string `json:"msg"`
		Code string `json:"code"`
	}
)


# 生成go文件
goctl api go --api core.api --dir=. --style go_zero
```

- email发送
```bash
# core\helper\helper.go
package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"cloud-disk/core/define"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
)

func Md5(s string) string {

	// md5加密
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int, identity, name string) (string, error) {

	// jwt token生成
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
	}
	// jwt token加密
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	// token加盐
	tikenstring, err := token.SignedString([]byte(define.Shar))
	if err != nil {
		return "", err
	}
	return tikenstring, nil

}

// 邮件发送
func SendEmail(emails, code string) error {
	e := email.NewEmail()
	e.From = "Get <hu729919300@163.com>" // 来自
	e.To = []string{emails}              // 发送
	// e.Cc = []string{"xxx@126.com"} // 抄送
	e.Subject = "验证码测试发送" // 主题

	e.HTML = []byte("你的验证码是: <h3>" + code + "</h3>") // 内容

	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "hu729919300@163.com", define.Pwd, "smtp.163.com"), &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.163.com",
	})

	if err != nil {
		return err
	}
	return nil
}


```
- 业务逻辑
```bash
# core\internal\logic\mail_code_send_logic.go
import (
	"context"
	"math/rand"
	"strconv"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...

func (l *MailCodeSendLogic) MailCodeSend(req *types.MailCodeReq) (resp *types.MailCodeResp, err error) {
	// todo: add your logic here and delete this line

	code := strconv.Itoa(rand.Intn(100000) + 100000)
	err = helper.SendEmail(req.Email, code)
	if err != nil {
		return nil, err
	}
	resp = &types.MailCodeResp{
		Msg:  "验证码发送成功!",
		Code: code,
	}

	return
}

```


##### 用户注册
1、判断邮箱是否存在 ---> 不存在则发送验证码
- 安装redis
```bash

# redis依赖
go get github.com/go-redis/redis/v8

```
- redis配置
```bash
# core/internal/modules/init_db.go
// redis初始化
import (
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"xorm.io/xorm"
)
...
var RDB = InitRDB()

func InitRDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "10.0.0.91:6379",
		Password: "123",
		DB:       0,
	})
}

```
- 验证码保存
```bash
# core/define/define.go
// 定义验证码长度
var CodeLength = 6

// 验证码过期时间
var CodeExoire = 300

# core\internal\logic\mail_code_send_logic.go
import (
	"context"
	"errors"

	// "math/rand"
	// "strconv"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

...
func (l *MailCodeSendLogic) MailCodeSend(req *types.MailCodeReq) (resp *types.MailCodeResp, err error) {
	// todo: add your logic here and delete this line

	// 1、邮箱不存在
	cnt, err := models.Engine.Where("email = ?", req.Email).Count(new(models.User_basic))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("该邮箱已存在!")
	}

	// 2、获取随机验证码
	// code := strconv.Itoa(rand.Intn(100000) + 100000)
	code := helper.RandCode()

	// 3、存储验证码
	models.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExoire))
	if err != nil {
		return nil, err
	}

	// 4、发送验证码
	err = helper.SendEmail(req.Email, code)
	if err != nil {
		return nil, err
	}
	resp = &types.MailCodeResp{
		Msg:  "验证码发送成功!",
		Code: code,
	}

	return
}

```


- 测试发送
```bash
go run core.go -f etc/core-api.yaml

# 发送请求
curl --location --request POST 'http://localhost:8888/mail/send' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "email": "1234@admin.com "
}'

// 响应
{
    "msg": "验证码发送成功!",
    "code": "029602"
}
```

2、用户注册
- api文件
```bash
service core-api {
	@doc(
		summary: "用户登录"
	)
	@handler UserLogin
	post /user/login(LoginRequest) returns (LoginResponse)
	
	@doc(
		summary: "用户详情"
	)
	@handler UserDetails
	get /user/details(UserDetailsReq) returns (UserDetailsResp)
	
	@doc(
		summary: "邮箱验证码"
	)
	@handler MailCodeSend
	post /mail/send(MailCodeReq) returns (MailCodeResp)
	
	@doc(
		summary: "用户注册"
	)
	@handler UserRegister
	post /user/register(UserRegisterReq) returns (UserRegisterResp)
	
}

// 用户登录
type LoginRequest {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse {
	Token string `json:"token"`
}

// 用户详情
type (
	UserDetailsReq {
		Identity string `json:"identity"`
	}

	UserDetailsResp {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

// 邮箱验证码
type (
	MailCodeReq {
		Email string `json:"email"`
	}

	MailCodeResp {
		Msg  string `json:"msg"`
		Code string `json:"code"`
	}
)

// 用户注册
type (
	UserRegisterReq {
		// 用户名
		Name string `json:"name"`
		// 密码
		Password string `json:"password"`
		// 邮箱
		Email string `json:"email"`
		// 验证码
		Code string `json:"code"`
	}

	UserRegisterResp {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Code     string `json:"code"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```

- 更新models表字段
```bash
# core\internal\models\user_basic.go
...
// 字段映射
type User_basic struct {
	Id        int
	Identity  string
	Name      string
	Password  string
	Email     string
	CreatedAt time.Time `xorm:"created" description:"创建时间"`
	UpdatedAt time.Time `xorm:"updated" description:"更新时间"`
	DeletedAt time.Time `xorm:"deleted" description:"删除时间"`
}

```
- uuid方法
```bash
# 安装依赖
go get github.com/satori/go.uuid

# core\helper\helper.go
import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/smtp"

	"time"

	"cloud-disk/core/define"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
)
...
// 生成uuid
func GetUuid() string {

	return uuid.NewV4().String()
}

```
- 业务逻辑
```bash
# core\internal\logic\user_register_logic.go

import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (resp *types.UserRegisterResp, err error) {
	// todo: add your logic here and delete this line

	// 1、判断code是否一致
	code, err := models.RDB.Get(l.ctx, req.Email).Result()
	if err != nil {
		return nil, err
	}
	if code != req.Code {
		return nil, errors.New("验证码不一致")
	}

	// 2、判断用户是否存在
	cnt, err := models.Engine.Where("name = ?", req.Name).Count(new(models.User_basic))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("用户已存在")
	}

	// 3、新建用户
	user := &models.User_basic{
		Identity: helper.GetUuid(),
		Name:     req.Name,
		Password: helper.Md5(req.Password),
		Email:    req.Email,
	}
	n, err := models.Engine.Insert(user)
	if err != nil {
		return nil, err
	}
	logx.Infof("sql受影响的行数: %v", n)

	// 4、返回参数
	resp = &types.UserRegisterResp{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		Code:     "******",
	}

	return
}


```
- 测试
```bash
# 启动服务
go run core.go -f etc/core-api.yaml

# 1、获取验证码
curl --location --request POST 'http://localhost:8888/mail/send' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "email": "12345@admin.com"
}'

# 2、用户注册
curl --location --request POST 'http://localhost:8888/user/register' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"laoli",
    "password":"123456",
    "email":"12345@admin.com",
    "code":"474937"
}'

```

##### 配置提取
将MySQL，Redis的连接配置写在config文件里，使用go-zero来解析配置信息
- api yaml
```bash
# core\etc\core-api.yaml
Name: core-api
Host: 0.0.0.0
Port: 8888


# MySQL
Mysql:
  DataSource: root:123456@tcp(10.0.0.91:3306)/cloud-disk?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai


# redis
Redis:
  Addr: 10.0.0.91:6379
```
- config 
```bash
# core\internal\config\config.go
...
type Config struct {
	rest.RestConf

	// 添加配置信息
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr string
	}
}

```
- service
```bash
# core\internal\svc\service_context.go
package svc

import (
	"cloud-disk/core/internal/config"

	"cloud-disk/core/internal/models"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config config.Config

	// 定义客户端
	Engine *xorm.Engine
	RDB    *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// 配置引用，调用方法
		Engine: models.Init(c.Mysql.DataSource),
		RDB:    models.InitRDB(c.Redis.Addr),
	}
}

```
- 配置models
```bash
# core\internal\models\init_db.go
package models

import (
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"xorm.io/xorm"
)

// mysql初始化
func Init(datasource string) *xorm.Engine {
	// 建立连接
	engine, err := xorm.NewEngine("mysql", datasource)
	engine.ShowSQL(true)

	if err != nil {
		logx.Error("Xorm Engine Error: ", err)
		return nil
	}
	return engine
}

// redis初始化
func InitRDB(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "123",
		DB:       0,
	})
}

```

- 修改业务代码
```bash
core\internal\logic\mail_code_send_logic.go
core\internal\logic\user_details_logic.go
core\internal\logic\user_login_logic.go
core\internal\logic\user_register_logic.go

# 替换Engine，RDB的引用
models.Engine -> l.svcCtx.Engine
models.RDB -> l.svcCtx.RDB
```
- 验证服务
```bash
go run core.go -f etc/core-api.yaml

```

#### 文件存储

##### 存储池资源管理

###### 对象存储
使用腾讯云cos对象存储
参考文档: https://cloud.tencent.com/document/product/436/31215
> 域名: https://cloud-1304907914.cos.ap-guangzhou.myqcloud.com
cos相关设置
- 子用户
	子账号ID: 100030324927 
	用户名: cloud-test 
	登录密码: tCKgQ03` 
	SecretId: AKIDykGajdIjuMaAdttkOo33pOJ1iIcjhT0O
	SecretKey: SijW5PWO0NXR7ZpreCfAevxtd7LL8lAO
- 权限
	QcloudCOSDataFullControl 对象存储（COS）数据读、写、删除、列出的访问权限
	// QcloudCOSFullAccess		 对象存储（COS）全读写访问权限
- API 密钥
	SecretId: AKIDykGajdIjuMaAdttkOo33pOJ1iIcjhT0O
	SecretKey: SijW5PWO0NXR7ZpreCfAevxtd7LL8lAO
- 存储桶访问权限
  - 私有读写 
  - 子账号	100030324927(cloud-test)	数据读取、数据写入
	
```bash
# 安装依赖； go get -u github.com/tencentyun/cos-go-sdk-v5
// linux
export COS_SECRETID=AKIDykGajdIjuMaAdttkOo33pOJ1iIcjhT0O
export COS_SECRETKEY=SijW5PWO0NXR7ZpreCfAevxtd7LL8lAO

// win
set COS_SECRETID = "AKIDykGajdIjuMaAdttkOo33pOJ1iIcjhT0O"
set COS_SECRETKEY = "SijW5PWO0NXR7ZpreCfAevxtd7LL8lAO"

// cos_test.go
package test

import (
	"context"
	"fmt"
	"os"

	"net/url"

	"net/http"

	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func TestCosInit(t *testing.T) {
	fmt.Println(os.Getenv("COS_SECRETID"))
	fmt.Println(os.Getenv("COS_SECRETKEY"))
	fmt.Println("--------------------------------")

	// 存储桶名称，由 bucketname-appid + 地域组成
	u, _ := url.Parse("https://cloud-1304907914.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId
			SecretID: os.Getenv("COS_SECRETID"),
			// 环境变量 SECRETKEY 表示用户的 SecretKey
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	// 查询
	// opt := &cos.BucketGetOptions{
	// 	Prefix:    "cloud-disk", // prefix 表示要查询的文件夹
	// 	Delimiter: "/",          // deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
	// 	MaxKeys:   3,            // 设置最大遍历出多少个对象, 一次 listobject 最大支持1000
	// }
	// v, _, err := c.Bucket.Get(context.Background(), opt)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, c := range v.Contents {
	// 	fmt.Printf("%s, %d\n", c.Key, c.Size)
	// }

	// 上传,需要指定路径及文件名
	key := "cloud-disk/vue1.png"

	if _, _, err := c.Object.Upload(
		context.Background(), key, "../img/vue.png", nil,
	); err != nil {
		panic(err)
	}

	// 下载
	// key = "cloud-disk/vue.png"
	// file := "../img/vue.png"

	// opt := &cos.MultiDownloadOptions{
	// 	ThreadPoolSize: 5,
	// }
	// _, err := c.Object.Download(
	// 	context.Background(), key, file, opt,
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// 列出
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:    "cloud-disk/", // prefix 表示要查询的文件夹
		Delimiter: "/",           // deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
		MaxKeys:   1000,          // 设置最大遍历出多少个对象, 一次 listobject 最大支持1000
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := c.Bucket.Get(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			break
		}
		// common prefix 表示表示所有以 Prefix 开头，被 delimiter的值 截断的路径, 如 delimter 设置为/, common prefix 则表示所有子目录的路径
		for _, commonPrefix := range v.CommonPrefixes {
			fmt.Printf("当前目录: %v\n", commonPrefix)
		}

		for _, content := range v.Contents {
			fmt.Printf("cos对象文件夹/文件: %v\n", content.Key)
		}

		isTruncated = v.IsTruncated // 是否还有数据
		marker = v.NextMarker       // 设置下次请求的起始 key
	}

	// 删除
	// 1、删除文件
	key = "cloud-disk/vue1.png"
	_, err := c.Object.Delete(context.Background(), key)
	if err != nil {
		panic(err)
	}
	// 2、删除文件夹
	keys := "test/"
	if _, err := c.Object.Delete(context.Background(), keys); err != nil {
		panic(err)
	}
}


//  测试运行
go test -v cos_test.go

```

###### 文件上传

**未认证**

- models
```bash
package models

import "time"

// 字段映射
type repository_pool struct {
	Id        int
	Identity  string
	Hash      string
	Name      string
	Ext       string
	Size      int
	Path      string
	CreatedAt time.Time `xorm:"created" description:"创建时间"`
	UpdatedAt time.Time `xorm:"updated" description:"更新时间"`
	DeletedAt time.Time `xorm:"deleted" description:"删除时间"`
}

// 表明初始化
func (table *repository_pool) TableName() string {
	return "repository_pool"
}


```
- api文件
```bash
# core/core.api
service core-api {
	...
	@doc(
		summary: "文件上传"
	)
	@handler FileUpload
	post /file/upload(FileUploadReq) returns (FileUploadResp)
	
}

// 文件上传
type (
	FileUploadReq {
		Hash string `json:"hash,optional"`
		Name string `json:"name,optional"`
		Ext  string `json:"ext,optional"`
		Size int    `json:"size,optional"`
		Path string `json:"path,optional"`
	}

	FileUploadResp {
		Identity string `json:"identity"`
		Ext      string `json:"ext"`
		Name     string `json:"name"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero
```
- define
```bash
# core\define\define.go

// cos存储相关
var COS_SECRETID = os.Getenv("COS_SECRETID")
var COS_SECRETKEY = os.Getenv("COS_SECRETKEY")
var BucketUrl = "https://cloud-1304907914.cos.ap-guangzhou.myqcloud.com/"

```

- helper
```bash

# core\helper\helper.go
// cos上传文件
func CosUpload(r *http.Request) (string, error) {
	// 存储桶名称，由 bucketname-appid + 地域组成
	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId
			SecretID: define.COS_SECRETID,
			// 环境变量 SECRETKEY 表示用户的 SecretKey
			SecretKey: define.COS_SECRETKEY,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	// 上传,需要指定路径及文件名
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}

	key := "cloud-disk/" + GetUuid() + path.Ext(fileHeader.Filename)

	_, err = c.Object.Put(
		context.Background(), key, file, nil,
	)
	if err != nil {
		panic(err)
	}

	return define.BucketUrl + key, nil

}

```
- handler
```bash
# core\internal\handler\file_upload_handler.go

import (
	"cloud-disk/core/helper"
	"path"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 1、根据文件hash判断文件是否存在
		// 获取请求文件信息
		file, fileheader, err := r.FormFile("file")
		if err != nil {
			return
		}
		// 根据文件size的hash判断
		b := make([]byte, fileheader.Size)
		_, err = file.Read(b)
		if err != nil {
			return
		}
		hash := fmt.Sprintf("%x", md5.Sum(b))
		rp := new(models.Repository_pool)
		has, err := svcCtx.Engine.Where("hash = ?", hash).Get(rp)
		if err != nil {
			return
		}
		if has {
			httpx.OkJson(w, &types.FileUploadResp{
				Identity: rp.Identity,
				Name:     rp.Name,
				Ext:      rp.Ext,
			})
			return
		}

		// 2、req参数赋值
		cospath, err := helper.CosUpload(r)
		if err != nil {
			return
		}
		req.Name = fileheader.Filename
		req.Ext = path.Ext(fileheader.Filename)
		req.Size = int(fileheader.Size)
		req.Hash = hash
		req.Path = cospath

		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

```
- logic
```bash
# core\internal\logic\file_upload_logic.go
import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *FileUploadLogic) FileUpload(req *types.FileUploadReq) (resp *types.FileUploadResp, err error) {
	// todo: add your logic here and delete this line

	// 根据handeler传递过来的req参数，保存db中
	rp := &models.Repository_pool{
		Identity: helper.GetUuid(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}

	_, err = l.svcCtx.Engine.Insert(rp)
	if err != nil {
		return nil, err
	}
	resp = &types.FileUploadResp{
		Identity: rp.Identity,
		Name:     rp.Name,
		Ext:      rp.Ext,
	}

	return
}


```
- 发送请求
```bash
# 启动服务
go run core.go -f etc/core-api.yaml

# 发送请求
curl --location --request POST 'http://localhost:8888/file/upload' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--header 'Content-Type: multipart/form-data; boundary=--------------------------406765020503899960111884' \
--form 'file=@"C:\\Users\\hu\\Pictures\\xx.jpg"'

// 响应
{
    "identity": "7147a874-afab-4f78-b956-67f8f2740d82",
    "ext": ".jpg",
    "name": "xx.jpg"
}

// 检测数据库/cos存储桶
mysql> select * from repository_pool;
```

**已认证**

- api 
```bash
# core.api
// 用户登录
service core-api {
	@doc(
		summary: "用户登录"
	)
	@handler UserLogin
	post /user/login(LoginRequest) returns (LoginResponse)
	
	@doc(
		summary: "用户详情"
	)
	@handler UserDetails
	get /user/details(UserDetailsReq) returns (UserDetailsResp)
	
	@doc(
		summary: "邮箱验证码"
	)
	@handler MailCodeSend
	post /mail/send(MailCodeReq) returns (MailCodeResp)
	
	@doc(
		summary: "用户注册"
	)
	@handler UserRegister
	post /user/register(UserRegisterReq) returns (UserRegisterResp)
	
}

@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	@doc(
		summary: "文件上传"
	)
	@handler FileUpload
	post /file/upload(FileUploadReq) returns (FileUploadResp)
}

// 用户登录
type LoginRequest {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse {
	Token string `json:"token"`
}

// 用户详情
type (
	UserDetailsReq {
		Identity string `json:"identity"`
	}

	UserDetailsResp {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

// 邮箱验证码
type (
	MailCodeReq {
		Email string `json:"email"`
	}

	MailCodeResp {
		Msg  string `json:"msg"`
		Code string `json:"code"`
	}
)

// 用户注册
type (
	UserRegisterReq {
		// 用户名
		Name string `json:"name"`
		// 密码
		Password string `json:"password"`
		// 邮箱
		Email string `json:"email"`
		// 验证码
		Code string `json:"code"`
	}

	UserRegisterResp {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Code     string `json:"code"`
	}
)

// 文件上传
type (
	FileUploadReq {
		Hash string `json:"hash,optional"`
		Name string `json:"name,optional"`
		Ext  string `json:"ext,optional"`
		Size int    `json:"size,optional"`
		Path string `json:"path,optional"`
	}

	FileUploadResp {
		Identity string `json:"identity"`
		Ext      string `json:"ext"`
		Name     string `json:"name"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero
```

- helper
```bash
# core\helper\helper.go
// jwt token解析
func AnalyzeJwtTkoen(token string) (*define.UserClaim, error) {
	uc := new(define.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		// jwt shar盐
		return []byte(define.Shar), nil
	})

	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token 已过期")
	}

	return uc, err
}

```

- auth_middleware
```bash
# core\internal\middleware\auth_middleware.go
...
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// 认证中间件，请求时进行
		// 1、判断auth是否为空
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization不能为空"))
			return
		}

		// 2、验证auth的token信息
		uc, err := helper.AnalyzeJwtTkoen(auth)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("认证失败"))
			return
		}
		r.Header.Set("UserId", string(rune(uc.Id)))
		r.Header.Set("UserIdentity", uc.Identity)
		r.Header.Set("UserName", uc.Name)

		// Passthrough to next handler if need
		next(w, r)
	}
}

```

- service
```bash
# core\internal\svc\service_context.go
package svc

import (
	"cloud-disk/core/internal/config"

	"cloud-disk/core/internal/middleware" // 引入中间件
	"cloud-disk/core/internal/models"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/rest"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config config.Config

	// 定义客户端
	Engine *xorm.Engine
	RDB    *redis.Client

	// auth中间件定义
	Auth rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// 配置引用
		Engine: models.Init(c.Mysql.DataSource),
		RDB:    models.InitRDB(c.Redis.Addr),

		// 引入认证中间件
		Auth: middleware.NewAuthMiddleware().Handle,
	}
}

```
- logic
```bash
# core\internal\logic\file_upload_logic.go
import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *FileUploadLogic) FileUpload(req *types.FileUploadReq) (resp *types.FileUploadResp, err error) {
	// todo: add your logic here and delete this line

	// 根据handeler传递过来的req参数，保存db中
	rp := &models.Repository_pool{
		Identity: helper.GetUuid(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}

	_, err = l.svcCtx.Engine.Insert(rp)
	if err != nil {
		return nil, err
	}
	resp = &types.FileUploadResp{
		Identity: rp.Identity,
		Name:     rp.Name,
		Ext:      rp.Ext,
	}

	return
}

```


- 发送请求
```bash
# 启动服务
goctl api go --api core.api --dir=. --style=go_zero

# 发送请求
1、用户登录获取token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'

2、根据token进行文件上传
curl --location --request POST 'http://localhost:8888/file/upload' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--header 'Content-Type: multipart/form-data; boundary=--------------------------749655764149154342307149' \
--form 'file=@"C:\\Users\\hu\\Pictures\\xx.jpg"'

```
###### 文件秒传
- api
```bash
# core.api
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	@doc(
		summary: "文件秒传"
	)
	@handler UserFileSecondUpload
	post /user/file/secondupload(UserFileSecondUploadReq) returns ( UserFileSecondUploadResp)
	
}

// 用户-文件秒传
type (
	UserFileSecondUploadReq {
		Md5  string `json:"md5"`
		Name string `json:"name"`
	}
	UserFileSecondUploadResp {
		Identity string `json:"identity"`
		UploadId string `json:"uploadId"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```

- logic
```bash
# core\internal\logic\user_file_second_upload_logic.go
import (
	"context"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserFileSecondUploadLogic) UserFileSecondUpload(req *types.UserFileSecondUploadReq) (resp *types.UserFileSecondUploadResp, err error) {
	// todo: add your logic here and delete this line

	// 1、根据md5值进行判断
	rp := new(models.Repository_pool)
	hs, err := l.svcCtx.Engine.Where("hash = ? ", req.Md5).Get(rp)
	if err != nil {
		return
	}
	if hs {
		resp = &types.UserFileSecondUploadResp{
			// 秒传成功
			Identity: rp.Identity,
		}

	} else {
		// 进行分片上传
		logx.Info("进行分片上传")
	}

	return
}
```
- test
```bash
1、启动服务
go run core.go -f etc/core-api.yaml

2、获取token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'

3、获取md5
select * from repository_pool;

4、测试秒传
curl --location --request POST 'http://localhost:8888/user/file/secondupload' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIiwiZXhwIjoxNjc5ODk4MzQ5fQ.YLpF2qWQzYQbHHykfCIHGaEB-LJTO_NSsnsCdaw-dPw' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "md5": "929f45aa8b6e1fb650218c44a63de101",
    "name": ""
}'

```


###### 文件分片/合并/一致性验证
- test
```bash
# file_test.go

package test

import (
	"crypto/md5"
	"fmt"
	"io" // io/ioutil在go 1.19版本已经合并到io包中

	"math"
	"os"
	"strconv"
	"testing"
)

// 分片大小
const chunkSize = 10 * 1024 * 1024 // 10MB

// 1、文件分片
func TestChunkNumFile(t *testing.T) {
	// 读取文件
	fileInfo, err := os.Stat("../mv/yequ.mp4")
	if err != nil {
		t.Fatal(err)
	}

	// 分片个数 = 文件大小 / 分片大小
	// 390 / 100 ==> 4, 向上取整
	chunkNum := math.Ceil(float64(fileInfo.Size()) / chunkSize)
	// 只读方式打开文件
	myFile, err := os.OpenFile("../mv/yequ.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	// 存放每一次的分片数据
	b := make([]byte, chunkSize)
	// 遍历所有分片
	for i := 0; i < int(chunkNum); i++ {
		// 指定读取文件的起始位置
		myFile.Seek(int64(i*chunkSize), 0)
		// 最后一次的分片数据不一定是整除下来的数据
		// 例如: 文件 120M, 第一次读了 100M, 剩下只有 20M
		if chunkSize > fileInfo.Size()-int64(i*chunkSize) {
			b = make([]byte, fileInfo.Size()-int64(i*chunkSize))
		}
		myFile.Read(b)

		f, err := os.OpenFile("../mv/"+strconv.Itoa(i)+".chunk", os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		f.Write(b)
		f.Close()
	}
	defer myFile.Close()

}

// 2、分片文件的合并
func TestMergeChunkFile(t *testing.T) {
	myFile, err := os.OpenFile("../mv/yequ1.mp4", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// 计算分片个数, 正常应该由前端传来, 这里测试时自行计算
	fileInfo, err := os.Stat("../mv/yequ.mp4")
	if err != nil {
		t.Fatal(err)
	}
	// 分片个数 = 文件大小 / 分片大小
	chunkNum := math.Ceil(float64(fileInfo.Size()) / chunkSize)

	// 合并分片
	for i := 0; i < int(chunkNum); i++ {
		f, err := os.OpenFile("../mv/"+strconv.Itoa(i)+".chunk", os.O_RDONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		b, err := io.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		myFile.Write(b)
		f.Close()
	}

	defer myFile.Close()
}

// 3、文件的一致性
func TestCheckFile(t *testing.T) {
	// 获取第一个文件的信息
	f1, err := os.OpenFile("../mv/yequ.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b1, err := io.ReadAll(f1)
	if err != nil {
		t.Fatal(err)
	}

	// 获取第二个文件的信息
	f2, err := os.OpenFile("../mv/yequ1.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := io.ReadAll(f2)
	if err != nil {
		t.Fatal(err)
	}

	s1 := fmt.Sprintf("%x", md5.Sum(b1))
	s2 := fmt.Sprintf("%x", md5.Sum(b2))

	fmt.Println("s1 = ", s1)
	fmt.Println("s2 = ", s2)
	fmt.Println(s1 == s2)
}


# 测试运行
go test -v file_test.go
```

###### cos测试分片文件上传
```bash
# test\cos_file_chunk_test.go

package test

import (
	"bytes"
	"cloud-disk/core/define"
	"context"

	"fmt"

	"math"
	"os"
	"strconv"

	"net/http"
	"net/url"

	// "strconv"

	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// 文件分片
func TestCosChunkNumFile(t *testing.T) {
	// 读取文件
	fileInfo, err := os.Stat("../img/vue.png")
	if err != nil {
		t.Fatal(err)
	}

	// 分片个数 = 文件大小 / 分片大小
	// 390 / 100 ==> 4, 向上取整
	chunkNum := math.Ceil(float64(fileInfo.Size()) / chunkSize)
	// 只读方式打开文件
	myFile, err := os.OpenFile("../img/vue.png", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	// 存放每一次的分片数据
	b := make([]byte, chunkSize)
	// 遍历所有分片
	for i := 0; i < int(chunkNum); i++ {
		// 指定读取文件的起始位置
		myFile.Seek(int64(i*chunkSize), 0)
		// 最后一次的分片数据不一定是整除下来的数据
		// 例如: 文件 120M, 第一次读了 100M, 剩下只有 20M
		if chunkSize > fileInfo.Size()-int64(i*chunkSize) {
			b = make([]byte, fileInfo.Size()-int64(i*chunkSize))
		}
		myFile.Read(b)

		f, err := os.OpenFile("./"+strconv.Itoa(i)+".chunk", os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		f.Write(b)
		f.Close()
	}
	defer myFile.Close()

}

// 分片上传初始化
func TestInitPartUpload(t *testing.T) {
	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.COS_SECRETID,
			SecretKey: define.COS_SECRETKEY,
		},
	})

	key := "cloud-disk/111.png"
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		t.Fatal(err)
	}
	UploadID := v.UploadID
	fmt.Println(UploadID) // 16798358340b5b37b4bc47a02a4dca772975088fd07c8e27493e89c5cd7ad6c06eec2ee107
}

// 分片上传
func TestPartUpload(t *testing.T) {
	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.COS_SECRETID,
			SecretKey: define.COS_SECRETKEY,
		},
	})

	key := "cloud-disk/111.png"
	UploadID := "16798358340b5b37b4bc47a02a4dca772975088fd07c8e27493e89c5cd7ad6c06eec2ee107"
	f, err := os.ReadFile("0.chunk")
	if err != nil {
		t.Fatal(err)
	}
	// opt可选
	resp, err := client.Object.UploadPart(
		context.Background(), key, UploadID, 1, bytes.NewReader(f), nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	PartETag := resp.Header.Get("ETag") // chunk的Md5值: 82b9c7a5a3f405032b1db71a25f67021
	fmt.Println(PartETag)
}

// 分片上传完成
func TestPartUploadComplete(t *testing.T) {
	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.COS_SECRETID,
			SecretKey: define.COS_SECRETKEY,
		},
	})

	key := "cloud-disk/111.png"
	UploadID := "16798358340b5b37b4bc47a02a4dca772975088fd07c8e27493e89c5cd7ad6c06eec2ee107"

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, cos.Object{
		PartNumber: 1, ETag: "82b9c7a5a3f405032b1db71a25f67021"},
	)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, UploadID, opt,
	)
	if err != nil {
		t.Fatal(err)
	}
}

# 测试
go test -v file_test.go
```

###### COS切片文件上传 - 初始化

- api
```bash
# core.api
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	@doc(
		summary: "文件上传前预处理"
	)
	@handler UserFileUploadPrepare
	post /user/file/upload/prepare(UserFileUploadPrepareReq) returns ( UserFileUploadPrepareResp)
	
}

type (
	UserFileUploadPrepareReq {
		Md5  string `json:"md5"`
		Name string `json:"name"`
		Ext  string `json:"ext"`
	}
	UserFileUploadPrepareResp {
		Identity string `json:"identity"`
		UploadId string `json:"upload_id"`
		Key      string `json:"key"`
	}
)


# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- helper
```bash
# core\helper\helper.go
// COS 文件分片上传前
// 1.初始化
func CosInitPart(ext string) (string, string, error) {

	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.COS_SECRETID,
			SecretKey: define.COS_SECRETKEY,
		},
	})

	// cos存储的文件路径
	key := "cloud-disk/" + GetUuid() + ext
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		return "", "", err
	}

	return key, v.UploadID, nil
}

```
- logic
```bash
# core\internal\logic\user_file_upload_prepare_logic.go
import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

...
func (l *UserFileUploadPrepareLogic) UserFileUploadPrepare(req *types.UserFileUploadPrepareReq) (resp *types.UserFileUploadPrepareResp, err error) {
	// todo: add your logic here and delete this line

	// 1、根据md5值进行判断
	rp := new(models.Repository_pool)
	hs, err := l.svcCtx.Engine.Where("hash = ? ", req.Md5).Get(rp)
	if err != nil {
		return
	}

	resp = new(types.UserFileUploadPrepareResp)
	if hs {
		// 秒传成功
		resp.Identity = rp.Identity

	} else {
		// 获取该文件的UploadID、Key,用来进行文件的分片上传
		key, uploadId, err := helper.CosInitPart(req.Ext)
		if err != nil {
			return nil, err
		}
		resp.Key = key
		resp.UploadId = uploadId
	}

	return
}

```

- test
```bash
1、启动服务
go run core.go -f etc/core-api.yaml

2、查询hash值，也就是请求的md5值
select * from repository_pool;

3、发送请求(先获取token)
// 不存在的md5值
curl --location --request POST 'http://localhost:8888/user/file/upload/prepare' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIiwiZXhwIjoxNjc5OTMxOTg1fQ.qBXgy_7fKkY_N8QytocvHrndSJyZd30ZQLBdO9aMnRM' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "md5": "929f45aa8b6e1fb650218c44a63de132",
    "name": "111",
    "ext":".png"
}'

// 响应
{
    "identity": "",
    "upload_id": "16798465646ee3526446247c4aaefdff614567f2e20b131f57580350fc7a62a46efcbfb66f",
    "key": "cloud-disk/ea34dc36-7c3b-42c2-92ab-dd2ead7cb332.png"
}

```
###### COS切片文件上传 - 切片文件上传
- api
```bash
# core.api
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个	
	@doc(
		summary: "COS文件上传-分片上传"
	)
	@handler UserFileChunkUpload
	post /user/file/chunk/upload(UserFileChunkUploadReq) returns ( UserFileChunkUploadResp)
	
}
...
// COS 文件上传 - 分片上传
type (
	UserFileChunkUploadReq { // form_data
		// key, upload_id, chunk文件，part_number
	}
	UserFileChunkUploadResp {
		Etag string `json:"etag" description:"获取分片文件chunk的md5值"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- helper
```bash
# core\helper\helper.go
import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/smtp"
	"path"
	"strconv"
	"strings"

	"cloud-disk/core/define"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"

	"context"

	"net/url"

	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)
...
// 2.分片上传
func CosPartUpload(r *http.Request) (string, error) {
	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.COS_SECRETID,
			SecretKey: define.COS_SECRETKEY,
		},
	})

	key := r.PostForm.Get("key")
	UploadID := r.PostForm.Get("upload_id")
	part_number, err := strconv.Atoi(r.PostForm.Get("part_number"))
	if err != nil {
		return "", err
	}
	f, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}

	// 此处f可能会报错: 411 MissingContentLength
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	// opt可选
	resp, err := client.Object.UploadPart(
		//context.Background(), key, UploadID, part_number, f, nil,
		context.Background(), key, UploadID, part_number, bytes.NewBuffer(buf.Bytes()), nil,
	)
	if err != nil {
		return "", err
	}
	// PartETag := resp.Header.Get("ETag") // chunk的Md5值: 82b9c7a5a3f405032b1db71a25f67021
	// fmt.Println(PartETag)
	return strings.Trim(resp.Header.Get("ETag"), "\""), nil
}

```
- handler
```bash
# core\internal\handler\user_file_chunk_upload_handler.go
import (
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

...

func UserFileChunkUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFileChunkUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 对请求参数的判断
		if r.PostForm.Get("key") == "" || r.PostForm.Get("part_number") == "" || r.PostForm.Get("upload_id") == "" {
			httpx.Error(w, errors.New("key or upload_id or post_id is empty"))
			return
		}

		// 获取etag
		etag, err := helper.CosPartUpload(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserFileChunkUploadLogic(r.Context(), svcCtx)
		resp, err := l.UserFileChunkUpload(&req)

		resp = new(types.UserFileChunkUploadResp)
		resp.Etag = etag
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

```
- test
```bash
1、启动服务
go run core.go -f etc/core-api.yaml

2、发送请求(先获取token)
curl --location --request POST 'http://localhost:8888/user/file/chunk/upload' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIiwiZXhwIjoxNjc5OTMxOTg1fQ.qBXgy_7fKkY_N8QytocvHrndSJyZd30ZQLBdO9aMnRM' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--header 'Content-Type: multipart/form-data; boundary=--------------------------345030900949343756079767' \
--form 'key="cloud-disk/ea34dc36-7c3b-42c2-92ab-dd2ead7cb332.png"' \
--form 'part_number="1"' \
--form 'upload_id="16798465646ee3526446247c4aaefdff614567f2e20b131f57580350fc7a62a46efcbfb66f"' \
--form 'file=@"D:\\study\\project\\golang\\7-go-zero\\cloud-disk\\test\\0.chunk"'

// 响应
{
    "etag": "82b9c7a5a3f405032b1db71a25f67021"
}
```

###### COS切片文件上传 - 切片文件上传完成
- api
```bash
# core.api
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个

	@doc(
		summary: "COS切片文件上传完成"
	)
	@handler UserFileUploadChunkComplete
	post /user/file/upload/chunk/complete(UserFileUploadChunkCompleteReq) returns ( UserFileUploadChunkCompleteResp)
	
}
...
// COS文件上传：切片文件上传完成
type (
	UserFileUploadChunkCompleteReq {
		Key string `json:"key"`
		UploadId string `json:"upload_id"`
		CosObjects []CosObjects `json:"cos_objects"`

	}
	CosObjects {
		PartNumber string `json:"part_number"`
		Etag string `json:"etag"`
	}
	UserFileUploadChunkCompleteResp {

	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- helper
```bash
# core\helper\helper.go

// 3.分片上传完成
func CosPartUploadComplete(key, uploadid string, co []cos.Object) error {
	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.COS_SECRETID,
			SecretKey: define.COS_SECRETKEY,
		},
	})

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, co...)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, uploadid, opt,
	)
	if err != nil {
		return err
	}
	return nil
}


```
- logic
```bash
# core\internal\logic\user_file_upload_chunk_complete_logic.go
import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserFileUploadChunkCompleteLogic) UserFileUploadChunkComplete(req *types.UserFileUploadChunkCompleteReq) (resp *types.UserFileUploadChunkCompleteResp, err error) {
	// todo: add your logic here and delete this line

	// 分片文件上传完成
	co := make([]cos.Object, 0)
	for _, v := range req.CosObjects {
		co = append(co, cos.Object{
			ETag:       v.Etag,
			PartNumber: v.PartNumber,
		})
	}

	err = helper.CosPartUploadComplete(req.Key, req.UploadId, co)

	return
}


```
- test
```bash
1、启动服务
go run core.go -f etc/core-api.yaml

2、发送请求(先获取token)
curl --location --request POST 'http://localhost:8888/user/file/upload/chunk/complete' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIiwiZXhwIjoxNjc5OTMxOTg1fQ.qBXgy_7fKkY_N8QytocvHrndSJyZd30ZQLBdO9aMnRM' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "key":"cloud-disk/ea34dc36-7c3b-42c2-92ab-dd2ead7cb332.png",
    "upload_id": "16798465646ee3526446247c4aaefdff614567f2e20b131f57580350fc7a62a46efcbfb66f",
    "cos_objects": [{
            "etag": "82b9c7a5a3f405032b1db71a25f67021",
            "part_number": 1
    }]
}'

// 登录腾讯云，进入COS对象存储查看验证
```

##### 个人存储资源管理

###### 用户 - 文件存储
- api
```bash
...
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	@doc(
		summary: "文件上传"
	)
	@handler FileUpload
	post /file/upload(FileUploadReq) returns (FileUploadResp)
	
	@doc(
		summary: "用户文件的关联存储"
	)
	@handler UserRepositorySave
	post /user/repository/save(UserRepositorySaveReq) returns (UserRepositorySaveResp)
	
}

// 用户关联存储
type (
	UserRepositorySaveReq {
		ParentId           int    `json:"parentid" description:"父级id"`
		RepositoryIdentity string `json:"repositoryIdentity"`
		Ext                string `json:"ext"`
		Name               string `json:"name"`
	}
	UserRepositorySaveResp {
		Identity string `json:"identity"`
	}
)


# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- models
```bash
# core\internal\models\user_repository.go
package models

import "time"

type UserRepository struct {
	Id                 int
	Identity           string
	UserIdentity       string
	ParentId           int
	RepositoryIdentity string
	Ext                string
	Name               string
	CreatedAt          time.Time `xorm:"created"`
	UpdatedAt          time.Time `xorm:"updated"`
	DeletedAt          time.Time `xorm:"deleted"`
}

func (table UserRepository) TableName() string {
	return "user_repository"
}

```

- logic
```bash
# core\internal\logic\user_repository_save_logic.go
import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...

func (l *UserRepositorySaveLogic) UserRepositorySave(req *types.UserRepositorySaveReq, userIdentity string) (resp *types.UserRepositorySaveResp, err error) {
	// todo: add your logic here and delete this line

	// 将文件信息跟用户进行关联
	ur := &models.UserRepository{
		Identity:           helper.GetUuid(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                req.Ext,
		Name:               req.Name,
	}
	_, err = l.svcCtx.Engine.Insert(ur)
	if err != nil {
		return nil, err
	}
	resp = &types.UserRepositorySaveResp{
		Identity: ur.Identity,
	}

	return
}
```
- handler
```bash
# core\internal\handler\user_repository_save_handler.go
...
func UserRepositorySaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRepositorySaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserRepositorySaveLogic(r.Context(), svcCtx)
		// 获取请求Header的UserIdentity,即jwt token的identity，即user_basic中的identity
		resp, err := l.UserRepositorySave(&req, r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
```
- 发送请求
```bash
# 启动服务
go run core.go -f etc/core-api.yaml

# 发送请求
1、用户登录获取token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'

2、文件上传获取文件信息
curl --location --request POST 'http://localhost:8888/file/upload' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--header 'Content-Type: multipart/form-data; boundary=--------------------------349240731018021951541591' \
--form 'file=@"C:\\Users\\hu\\Pictures\\xx.jpg"'

3、根据token，文件信息绑定用户
curl --location --request POST 'http://localhost:8888/user/repository/save' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "parentid": 0,
    "repositoryIdentity": "f0d830cb-1285-4d27-8654-5e5d8ddc737f",
    "ext": ".jpg",
    "name": "xx.jpg"
}'

```

###### 用户 - 文件列表

- api
```bash
# core.api
...
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
...
	@doc(
		summary: "用户文件列表"
	)
	@handler UserFileList
	get /user/file/list(UserFileListReq) returns (UserFileListResp)
	
}


// 用户文件列表
type (
	UserFileListReq {
		Id   int `json:"id,optional"`
		Page int `json:"page,optional" description:"页数"`
		Size int `json:"size,optional" description:"每页条数"`
	}
	UserFileListResp {
		List  []*UserFile `json:"list,optional"`
		Count int64       `json:"count,optional"`
	}

	UserFile {
		Id                 int    `json:"id"`
		Identity           string `json:"identity"`
		RepositoryIdentity string `json:"repository_identity"`
		Name               string `json:"name"`
		Ext                string `json:"ext"`
		Path               string `json:"path"`
		Size               string `json:"size"`
	}
)


# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- define
```bash
# core\define\define.go
// 定义时间格式
var Datetime = "2006-01-02 15:04:05"
```
- hendeler
```bash
# core\internal\handler\user_file_list_handler.go
...
func UserFileListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFileListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserFileListLogic(r.Context(), svcCtx)

		// 传递用户信息
		resp, err := l.UserFileList(&req, r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

```
- logic
```bash
# core\internal\logic\user_file_list_logic.go
import (
	"context"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
// 传递UserIdentity
func (l *UserFileListLogic) UserFileList(req *types.UserFileListReq, UserIdentity string) (resp *types.UserFileListResp, err error) {
	// todo: add your logic here and delete this line

	// 定义list,count
	uf := make([]*types.UserFile, 0)
	resp = new(types.UserFileListResp)
	//
	// type UserFile struct {
	// 	Id                 int    `json:"id"`
	// 	Identity           string `json:"identity"`
	// 	RepositoryIdentity string `json:"repository_identity"`
	// 	Name               string `json:"name"`
	// 	Ext                string `json:"ext"`
	// 	Path               string `json:"path"`
	// 	Size               string `json:"size"`
	// }

	// type UserRepository struct {
	// 	Id                 int
	// 	Identity           string
	// 	UserIdentity       string
	// 	ParentId           int
	// 	RepositoryIdentity string
	// 	Ext                string
	// 	Name               string
	// 	CreatedAt          time.Time `xorm:"created"`
	// 	UpdatedAt          time.Time `xorm:"updated"`
	// 	DeletedAt          time.Time `xorm:"deleted"`
	// }

	// 分页参数
	size := req.Size
	if size == 0 {
		size = define.PageSize
	}
	page := req.Page
	if page == 0 {
		page = define.Page
	}
	offsize := (page - 1) * size

	// 查询用户文件列表
	err = l.svcCtx.Engine.Table("user_repository").Where("parent_id = ? AND user_identity = ? ", req.Id, UserIdentity).Select("user_repository.id, user_repository.identity, user_repository.repository_identity, user_repository.ext, user_repository.name, "+" repository_pool.path, repository_pool.size").Join("LEFT", "repository_pool", "user_repository.repository_identity = repository_pool.identity").Where("user_repository.deleted_at = ? OR user_repository.deleted_at IS NULL", time.Time{}.Format(define.Datetime)).Limit(size, offsize).Find(&uf)

	if err != nil {
		return
	}

	// 查询用户文件总数
	cnt, err := l.svcCtx.Engine.Where("parent_id = ? AND user_identity = ?", req.Id, UserIdentity).Count(new(models.UserRepository))
	if err != nil {
		return
	}
	resp.List = uf
	resp.Count = cnt

	return
}


```

- 发送请求
```bash
# 启动服务
go run core.go -f etc/core-api.yaml

# 发送请求
1、获取token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'

2、请求文件列表
curl --location --request GET 'http://localhost:8888/user/file/list' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{}'
```

###### 用户 - 文件名修改
- api
```bash
# core.api
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	...	
	@doc(
		summary: "用户文件名修改"
	)
	@handler UserFileUpdate
	post /user/file/update(UserFileUpdateReq) returns (UserFileUpdateResp)
	
}

// 用户文件名修改
type (
	UserFileUpdateReq {
		Identity string `json:"identity"`
		Name     string `json:"name"`
	}

	UserFileUpdateResp {
		Identity string `json:"identity"`
		Name     string `json:"name"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- hendler
```bash
# core\internal\handler\user_file_update_handler.go
...
func UserFileUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFileUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserFileUpdateLogic(r.Context(), svcCtx)

		// 绑定用户信息
		resp, err := l.UserFileUpdate(&req,r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

```
- logic
```bash
# core\internal\logic\user_file_update_logic.go
import (
	"context"
	"errors"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...

func (l *UserFileUpdateLogic) UserFileUpdate(req *types.UserFileUpdateReq, UserIdentity string) (resp *types.UserFileUpdateResp, err error) {
	// todo: add your logic here and delete this line

	// 判断文件名是否已存在
	cnt, err := l.svcCtx.Engine.Where("name = ? AND parent_id = (SELECT parent_id FROM user_repository ur WHERE ur.identity = ?)", req.Name, req.Identity).Count(new(models.UserRepository))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("文件名已存在")
	}

	// 修改文件名称
	data := &models.UserRepository{
		Name: req.Name,
	}

	_, err = l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.Identity, UserIdentity).Update(data)
	if err != nil {
		return nil, err
	}
	resp = &types.UserFileUpdateResp{
		Identity: req.Identity,
		Name:     req.Name,
	}
	return
}


```
- 发送请求
```bash
# 启动服务
go run core.go -f etc/core-api.yaml

# 发送请求
curl --location --request POST 'http://localhost:8888/user/file/update' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "identity":"94d73053-eeec-450d-bc29-bceae316e603",
    "name":"hhhhh"
}'
// identity的值是user_repository表中identity字段的值;
```

###### 用户 - 文件夹创建
- api
```bash
# core.api
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	@doc(
		summary: "用户文件夹创建"
	)
	@handler UserFolderCreate
	post /user/folder/create(UserFolderCreateReq) returns (UserFolderCreateResp)
}


// 用户-文件夹创建
type (
	UserFolderCreateReq {
		ParentId int    `json:"parent_id"`
		Name     string `json:"name"`
	}

	UserFolderCreateResp {
		Identity string `json:"identity"`
		Name     string `json:"name"`
	}
)


# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- handler
```bash
# core\internal\handler\user_folder_create_handler.go
...
func UserFolderCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFolderCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserFolderCreateLogic(r.Context(), svcCtx)

		// 传递用户信息
		resp, err := l.UserFolderCreate(&req, r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

```
- logic
```bash
# core\internal\logic\user_folder_create_logic.go
import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserFolderCreateLogic) UserFolderCreate(req *types.UserFolderCreateReq, UserIdentity string) (resp *types.UserFolderCreateResp, err error) {
	// todo: add your logic here and delete this line

	// 判断文件夹名是否存在
	cnt, err := l.svcCtx.Engine.Where("name = ? AND parent_id = ?", req.Name, req.ParentId).Count(new(models.UserRepository))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("文件夹已存在")
	}

	// 创建文件夹
	folder := &models.UserRepository{
		Identity:     helper.GetUuid(),
		UserIdentity: UserIdentity,
		ParentId:     req.ParentId,
		Name:         req.Name,
	}
	_, err = l.svcCtx.Engine.Insert(folder)
	if err != nil {
		return nil, err
	}

	resp = &types.UserFolderCreateResp{
		Identity: folder.Identity,
		Name:     folder.Name,
	}

	return
}

```
- 测试请求
```bash
# 启动服务
go run core.go -f etc/core-api.yaml

# 发送请求
1、用户登录获取token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'

2、使用token修改指定文件夹名
curl --location --request POST 'http://localhost:8888/user/folder/create' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "parent_id":0,
    "name":"hhh"
}'
```

###### 用户 - 文件删除
- api
```bash
# core.api


	@doc(
		summary: "用户-文件删除"
	)
	@handler UserFileDelete
	delete /user/file/delete(UserFileDeleteReq) returns (UserFileDeleteResp)
	
}
...
// 用户-文件删除
type (
	UserFileDeleteReq {
		Identity string `json:"identity"`
	}
	UserFileDeleteResp {
		Message string `json:"message"`
	}
)


# 生产go文件
goctl api go --api core.api --dir=. --style=go_zero

```

- handler
```bash
# core\internal\handler\user_file_delete_handler.go
func UserFileDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFileDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserFileDeleteLogic(r.Context(), svcCtx)

		// 绑定用户信息
		resp, err := l.UserFileDelete(&req, r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

```

- logic
```bash
# core\internal\logic\user_file_delete_logic.go
import (
	"context"
	"fmt"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserFileDeleteLogic) UserFileDelete(req *types.UserFileDeleteReq, UserIdentity string) (resp *types.UserFileDeleteResp, err error) {
	// todo: add your logic here and delete this line

	// 软删除，只是删除数据库文件信息，添加deletetime字段值
	_, err = l.svcCtx.Engine.Where("user_identity = ? AND identity = ?", UserIdentity, req.Identity).Delete(new(models.UserRepository))
	if err != nil {
		return
	}
	resp = &types.UserFileDeleteResp{
		Message: fmt.Sprintf("identity = %v 删除成功", req.Identity),
	}
	return
}

```

- test
```bash
# 启动服务
go run core.go -f etc/core-api.yaml

# 发送请求
1、用户登录
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'

2、文件删除
curl --location --request DELETE 'http://localhost:8888/user/file/delete' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "identity":"94d73053-eeec-450d-bc29-bceae316e603"
}'

// identity是user_repository表中的字段

```

###### 用户 - 文件移动
- api
```bash
# core.api
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	@doc(
		summary: "用户-文件移动"
	)
	@handler UserFileMove
	put /user/file/move(UserFileMoveReq) returns (UserFileMoveResp)
	
}

// 用户-文件移动
type (
	UserFileMoveReq {
		Identity string `json:"identity"`
		ParentIdentity string `json:"parent_identity"`
	}
	UserFileMoveResp {
		Message string `json:"message"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```

- handler
```bash
# core\internal\handler\user_file_move_handler.go
func UserFileMoveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFileMoveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserFileMoveLogic(r.Context(), svcCtx)

		// 用户身份绑定
		resp, err := l.UserFileMove(&req,r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
```
- logic
```bash
# core\internal\logic\user_file_move_logic.go
import (
	"context"
	"errors"
	"fmt"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)


...
func (l *UserFileMoveLogic) UserFileMove(req *types.UserFileMoveReq, UserIdentity string) (resp *types.UserFileMoveResp, err error) {
	// todo: add your logic here and delete this line

	// 1、判断文件夹存不存在
	has, err := l.svcCtx.Engine.Where("identity = ? AND user_identity = ? ", req.Identity, UserIdentity).Get(new(models.UserRepository))
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New("文件夹不存在")
	}

	// 2、更新 ParentID,注意，文件信息deletetime未被记录
	_, err = l.svcCtx.Engine.Where("identity = ? ", req.Identity).Update(&models.UserRepository{
		ParentId: req.ParentIdentity,
	})
	if err != nil {
		return
	}

	resp = &types.UserFileMoveResp{
		Message: fmt.Sprintf("文件: %v,文件夹: %v", req.Identity, req.ParentIdentity),
	}
	return
}

```
- test
```bash
1、登录获取token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'
2、文件存储
curl --location --request POST 'http://localhost:8888/user/repository/save' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "parentid": 2,
    "repositoryIdentity": "f0d830cb-1285-4d27-8654-5e5d8ddc737f",
    "ext": ".jpg",
    "name": "xx.jpg"
}'
// 新建parentid

3、文件移动
curl --location --request PUT 'http://localhost:8888/user/file/move' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "parent_identity": 2,
    "identity":"86eaf75e-0270-49d2-8c79-c68953e1af18"
}'
// 数据库查询: select * from user_repository;

```
#### 文件分享
##### 用户 - 创建分享记录
- models
```bash
# core\internal\models\share_basic.go
package models

import "time"

// 字段映射
type Share_basic struct {
	Id                       int
	Identity                 string    `xorm:"identity"`
	User_identity            string    `xorm:"user_identity"`
	User_Repository_Identity string    `xorm:"user_repository_identity"`
	Repository_identity      string    `xorm:"repository_identity"`
	Expired_time             int       `xorm:"expired_time" description:"失效时间"`
	Click_num                int       `xorm:"click_num" description:"点击次数"`
	CreatedAt                time.Time `xorm:"created" description:"创建时间"`
	UpdatedAt                time.Time `xorm:"updated" description:"更新时间"`
	DeletedAt                time.Time `xorm:"deleted" description:"删除时间"`
}

// 表明初始化
func (table *Share_basic) TableName() string {
	return "share_basic"
}

```
- api
```bash
# core.api 

@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	@doc(
		summary: "用户-创建文件分享"
	)
	@handler UserShareBasicCreate
	post /user/share/create(UserShareBasicCreateReq) returns (UserShareBasicCreateResp)
	
}

// 用户-创建文件分享
type (
	UserShareBasicCreateReq {
		UserRepositoryIdentity string `json:"user_repository_identity"`
		Expiredtime            int    `json:"expired_time"`
	}

	// 返回该信息的Identity,即uuid
	UserShareBasicCreateResp {
		Identity string `json:"identity"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- handler
```bash
# core\internal\handler\user_share_basic_create_handler.go

func UserShareBasicCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserShareBasicCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserShareBasicCreateLogic(r.Context(), svcCtx)

		// 绑定用户身份信息
		resp, err := l.UserShareBasicCreate(&req,r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

```
- logic
```bash
# core\internal\logic\user_share_basic_create_logic.go
import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserShareBasicCreateLogic) UserShareBasicCreate(req *types.UserShareBasicCreateReq, UserIdentity string) (resp *types.UserShareBasicCreateResp, err error) {
	// todo: add your logic here and delete this line

	// 创建文件分享记录
	uuid := helper.GetUuid()
	ur := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Where("identity = ? ", req.UserRepositoryIdentity).Get(ur)
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New(req.UserRepositoryIdentity + " not found or delete")
	}
	data := &models.Share_basic{
		Identity:                 uuid,
		User_identity:            UserIdentity,
		User_Repository_Identity: req.UserRepositoryIdentity,
		Repository_identity:      ur.RepositoryIdentity,
		Expired_time:             req.Expiredtime,
	}

	_, err = l.svcCtx.Engine.Insert(data)
	if err != nil {
		return
	}

	resp = &types.UserShareBasicCreateResp{
		Identity: uuid,
	}
	return
}

```
- test
```bash
1、启动服务
go run core.go -f etc/core-api.yaml
2、获取需要分享的文件repository_identity
select * from user_repository;
3、用户登录获取token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'
4、创建文件分享
curl --location --request POST 'http://localhost:8888/user/share/create' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "repository_identity": "f0d830cb-1285-4d27-8654-5e5d8ddc737f",
    "expired_time": 0
}'
5、数据库查询
select * from share_basic;
```

##### 获取资源详情
根据分享记录信息，获取资源详情
- api
```bash
# core.api

service core-api {   // 注意: service只有一个
	@doc(
		summary: "获取分享文件详情"
	)
	@handler UserShareBasicInfo
	get /user/share/info(UserShareBasicInfoReq) returns (UserShareBasicInfoResp)
	
}

// 用户分享文件详情
type (
	UserShareBasicInfoReq {
		Identity string `json:"identity"`
	}
	UserShareBasicInfoResp {
		RepositoryIdentity string `json:"repository_identity"`
		Name               string `json:"name"`
		Ext                string `json:"ext"`
		Size               int    `json:"size"`
		Path               string `json:"path"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```

- logic
```bash
# core\internal\logic\user_share_basic_info_logic.go
import (
	"context"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)
...
func (l *UserShareBasicInfoLogic) UserShareBasicInfo(req *types.UserShareBasicInfoReq) (resp *types.UserShareBasicInfoResp, err error) {
	// todo: add your logic here and delete this line

	// 1、对分享记录的点击+1

	_, err = l.svcCtx.Engine.Table("share_basic").Exec("Update share_basic Set click_num = click_num + 1 Where identity = ? ", req.Identity)
	if err != nil {
		return nil, err
	}

	// 2、获取分享文件的详情
	resp = new(types.UserShareBasicInfoResp)
	_, err = l.svcCtx.Engine.Table("share_basic").Select("share_basic.repository_identity,user_repository.name,repository_pool.ext,repository_pool.size,repository_pool.path").Join("LEFT", "repository_pool", "share_basic.repository_identity = repository_pool.identity").Join("LEFT", "user_repository", "user_repository.identity = user_repository_identity").Where("share_basic.identity = ? ", req.Identity).Get(resp)
	if err != nil {
		return nil, err
	}

	return
}

```
- test
```bash
1、启动服务
go run core.go -f etc/core-api.yaml

2、获取分享的文件的identity
select * from share_basic;  -- identity

3、获取该文件详情
curl --location --request GET 'http://localhost:8888/user/share/info' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "identity": "71ca49a3-df2a-46ac-83ef-41ad299b5e5c"
}'

4、数据库验证
select * from share_basic;
SELECT share_basic.repository_identity,user_repository.name,repository_pool.ext,repository_pool.size,repository_pool.path FROM `share_basic` LEFT JOIN `repository_pool` ON share_basic.repository_identity = repository_pool.identity LEFT JOIN `user_repository` ON user_repository.identity = user_repository_identity WHERE share_basic.identity = "d228aa82-6193-4498-b81e-2374824ac718" LIMIT 1;
```

##### 用户 - 资源分享保存
- api
```bash
# core.api
@server(
	middleware: Auth  // 添加Auth中间件
)
service core-api {   // 注意: service只有一个
	@doc(
		summary: "用户-资源保存"
	)
	@handler UserShareBasicSave
	post /user/share/save(UserShareBasicSaveReq) returns ( UserShareBasicSaveResp)
}
...
// 用户 - 分享文件保存
type (
	UserShareBasicSaveReq {
		RepositoryIdentity string `json:"repository_identity"`
		ParentId           int    `json:"parent_id"`
	}
	UserShareBasicSaveResp {
		Identity string `json:"identity"`
	}
)

# 生成go文件
goctl api go --api core.api --dir=. --style=go_zero

```
- handler
```bash
# core\internal\handler\user_share_basic_save_handler.go
func UserShareBasicSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserShareBasicSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserShareBasicSaveLogic(r.Context(), svcCtx)

		// 绑定用户Identity
		resp, err := l.UserShareBasicSave(&req, r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
```
- logic
```bash
# core\internal\logic\user_share_basic_save_logic.go
import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

...
func (l *UserShareBasicSaveLogic) UserShareBasicSave(req *types.UserShareBasicSaveReq, UserIdentity string) (resp *types.UserShareBasicSaveResp, err error) {
	// todo: add your logic here and delete this line

	// 1、获取资源详情
	rp := new(models.Repository_pool)
	has, err := l.svcCtx.Engine.Where("identity = ?", req.RepositoryIdentity).Get(rp)
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New(req.RepositoryIdentity + " not found")
	}

	// 2、资源保存
	ur := &models.UserRepository{
		Identity:           helper.GetUuid(),
		UserIdentity:       UserIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                rp.Ext,
		Name:               rp.Name,
	}
	_, err = l.svcCtx.Engine.Insert(ur)
	if err != nil {
		return
	}
	resp = &types.UserShareBasicSaveResp{
		Identity: ur.Identity,
	}
	return
}

```
- test
```bash
1、启动服务
go run core.go -f etc/core-api.yaml

2、获取资源信息
select * from share_basic; -- repository_identity
select * from user_repository; -- parent_id 自定义

4、获取用户token
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name":"李四",
    "password":"123456"
}'

5、保存文件分享资源
curl --location --request POST 'http://localhost:8888/user/share/save' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiSWRlbnRpdHkiOiJxYXp3c3giLCJOYW1lIjoi5p2O5ZubIn0.i0heXn_gRiKAxP8miHtu_BZEWhx8SLeih9bWCdDL3W8' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8888' \
--header 'Connection: keep-alive' \
--data-raw '{
    "repository_identity": "f0d830cb-1285-4d27-8654-5e5d8ddc737f",
    "parent_id": 0
}'

6、查询
select * from user_repository;
```









```bash

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var users = []User{
	{"user1", "password1"},
	{"user2", "password2"},
	{"user3", "password3"},
}

func main() {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, u := range users {
			if u.Username == user.Username && u.Password == user.Password {
				c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

```















