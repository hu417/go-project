

# miniblog



参考: https://juejin.cn/book/7176608782871429175
代码: https://github.com/marmotedu/miniblog/blob/master/README.md

## 介绍
本博客系统实现了以下 2 类功能：

- 用户管理： 支持 用户注册、用户登录、获取用户列表、获取用户详情、更新用户信息、修改用户密码、注销用户 7 种用户操作；
- 博客管理： 支持 创建博客、获取博客列表、获取博客详情、更新博客内容、删除博客、批量删除博客 6 种博客操作。

博客系统功能及使用流程如下图所示:

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/bd2a760daa324d21a7d485f75125a120~tplv-k3u1fbpfcp-jj-mark:1512:0:0:0:q75.awebp" alt="miniblog" style="zoom:50%;" />

操作流程如下：

- 普通用户：注册用户后，登录系统，之后可以执行博客的 CURD 操作。还可以获取注册时的详细信息，再次更新用户信息，包括密码；
- root 用户：root 用户为博客系统的管理员用户，整个系统只有一个 root 用户。root 用户除了可以执行普通用户的所有操作之外，还可以执行 1 类管理员操作：获取用户列表。


## 规范设计

- 规范设计内容也比较多，通常包含代码规范、Commit 规范、版本规范、接口规范、日志规范、错误码规范等。miniblog 项目也有自己的规范设计，具体如下：

- 代码规范： 编码规范内容比较多，本课程就不做详细介绍。minibolog 所遵循的代码规范见：Go 代码开发规范。

- Commit 规范： 社区有多种 Commit 规范，例如 jQuery、Angular 等。在这些规范中，Angular 规范在功能上能够满足开发者 commit 需求，在格式上清晰易读，目前也是用得最多的。Commit 规范对于开发者来说是一个建议遵守的规范，但不是必须的规范。如果你感兴趣可以参考 Conventional Commits （中文翻译版本：Conventional Commits：一份让代码提交记录人机友好的规范）。

- 版本规范： 当前绝大部分优秀的 Go 项目都有自己的版本规范，其中语义化版本规范（Semantic Versioning，SemVer）是用的最多的规范。

- 接口规范： 接口文档又称为 API 文档，一般由后台开发人员编写，用来描述组件提供的 API 接口，以及如何调用这些 API 接口。一个好的 API 接口文档，应该是规范的。当前用的最多的接口规范是 OpenAPI 3.0 规范。因为内容比较多，本课程也不详细介绍，你可以参考官方文档进行学习：OpenAPI Specification（一个不错的中文翻译文档：开放 API 规范中文翻译）。

- 日志规范： miniblog 为了规范打印日志，也制定了日志规范

- 错误码规范： 现代的软件架构，很多都是对外暴露 RESTful API 接口，内部系统通信采用 RPC 协议。因为 RESTful API 接口有一些天生的优势，比如规范、调试友好、易懂，所以通常作为直接面向用户的通信规范。既然是直接面向用户，那么首先就要求消息返回格式是规范的；其次，如果接口报错，还要能给用户提供一些有用的报错信息，通常需要包含错误 Code（用来唯一定位一次错误）和 Message（用来展示出错的信息）。这就需要我们设计一套规范的、科学的错误码。


## 目录结构设计

目录结构是一个项目的门面。很多时候，根据目录结构就能看出开发者对这门语言的掌握程度。所以，在我看来，遵循一个好的目录规范，把代码目录设计的可维护、可扩展，甚至比文档规范、Commit 规范都要重要。通常，根据功能，可以将目录结构分为以下两种。

- 平铺式目录结构： 主要用在 Go 包中，相对简单。

- 结构化目录结构： 主要用在 Go 应用中，相对复杂。

### 平铺式目录结构


一个 Go 项目可以是一个应用，也可以是一个代码库，当项目是代码库时，比较适合采用平铺式目录结构。

平铺方式就是在项目的根目录下存放项目的代码，整个目录结构看起来更像是一层的，这种方式在很多库中存在，使用这种方式的好处是引用路径长度明显减少。例如 glog 包就是平铺式的，目录内容如下：

```bash
$ ls glog/
glog_file.go  glog.go  glog_test.go  go.mod  LICENSE  README.md

```

### 结构化目录结构

当前 Go 社区比较推荐的结构化目录结构是 project-layout。虽然它并不是官方和社区的规范，但因为组织方式比较合理，被很多 Go 开发人员接受。所以，我们可以把它当作是一个事实上的规范。本实战项目 miniblog 目录也遵循了 project-layout 项目的目录规范。

### miniblog 目录结构设计

miniblog 目录结构遵循了 project-layout 的目录结构设计。项目目录及功能介绍如下：

```bash
├── api # Swagger / OpenAPI 文档存放目录
│   └── openapi
│       └── openapi.yaml # OpenAPI 3.0 API 接口文档
├── cmd # main 文件存放目录
│   └── miniblog
│       └── miniblog.go
│   └── cron # 定时任务
│       └── email.go
├── configs # 配置文件存放目录
│   ├── miniblog.sql # 数据库初始化 SQL
│   ├── miniblog.yaml # miniblog 配置文件
│   └── nginx.conf # Nginx 配置
├── docs # 项目文档
│   ├── devel # 开发文档
│   │   ├── en-US # 英文文档
│   │   └── zh-CN # 中文文档
│   │       ├── architecture.md # miniblog 架构介绍
│   │       ├── conversions # 规范文档存放目录
│   │       │   ├── api.md # 接口规范
│   │       │   ├── commit.md # Commit 规范
│   │       │   ├── directory.md # 目录结构规范
│   │       │   ├── error_code.md # 错误码规范
│   │       │   ├── go_code.md # 代码规范
│   │       │   ├── log.md # 日志规范
│   │       │   └── version.md # 版本规范
│   │       └── README.md
│   ├── guide # 用户文档
│   │   ├── en-US # 英文文档
│   │   └── zh-CN # 中文文档
│   │       ├── announcements.md # 动态与公告
│   │       ├── best-practice # 最佳实践
│   │       ├── faq # 常见问题
│   │       ├── installation # 安装指南
│   │       ├── introduction # 产品介绍
│   │       ├── operation-guide # 操作指南
│   │       ├── quickstart # 快速入门
│   │       └── README.md
│   └── images # 项目图片存放目录
├── examples # 示例源码
├── go.mod
├── go.sum
├── init # Systemd Unit 文件保存目录
│   ├── miniblog.service # miniblog systemd unit
├── internal # 内部代码保存目录，这里面的代码不能被外部程序引用
│   ├── miniblog # miniblog 代码实现目录
│   │   ├── biz # biz 层代码
│   │   ├── controller # controller 层代码
│   │   │   └── v1 # API 接口版本
│   │   │       ├── post # 博客相关代码实现
│   │   │       │   ├── create.go # 创建博客
│   │   │       │   ├── delete_collection.go #批量删除博客
│   │   │       │   ├── delete.go # 删除博客
│   │   │       │   ├── get.go # 获取博客详情
│   │   │       │   ├── list.go # 获取博客列表
│   │   │       │   ├── post.go # 博客 Controller 结构定义、创建
│   │   │       │   └── update.go # 更新博客
│   │   │       └── user
│   │   │           ├── change_password.go # 修改用户密码
│   │   │           ├── create.go #创建用户
│   │   │           ├── delete.go # 删除用户
│   │   │           ├── get.go # 获取用户详情
│   │   │           ├── list.go # 获取用户列表
│   │   │           ├── login.go # 用户登录
│   │   │           ├── update.go  # 更新用户
│   │   │           └── user.go # 用户 Controller 结构定义、创建
│   │   ├── helper.go # 工具类代码存放文件
│   │   ├── miniblog.go # miniblog 主业务逻辑实现代码
│   │   ├── router.go # Gin 路由加载代码
│   │   └── store # store 层代码
│   ├── cron  # 定时任务
│   │   ├── biz # biz 层代码
│   │   └── store # store 层代码
│   └── pkg # 内部包保存目录
│       ├── core # core 包，用来保存一些核心的函数
│       ├── errno # errno 包，实现了 miniblog 的错误码功能
│       │   ├── code.go # 错误码定义文件
│       │   └── errno.go # errno 包功能函数文件
│       ├── known # 存放项目级的常量定义
│       ├── log # miniblog 自定义 log 包
│       ├── middleware # Gin 中间件包
│       │   ├── authn.go # 认证中间件
│       │   ├── authz.go # 授权中间件
│       │   ├── header.go # 指定 HTTP Response Header
│       │   └── requestid.go # 请求 / 返回头中添加 X-Request-ID
│       └── model # GORM Model
├── LICENSE # 声明代码所遵循的开源协议
├── Makefile # Makefile 文件，一般大型软件系统都是采用 make 来作为编译工具
├── _output # 临时文件存放目录
├── pkg # 可供外部程序直接使用的 Go 包存放目录
│   ├── api # REST API 接口定义存放目录
│   ├── proto # Protobuf 接口定义存放目录
│   ├── auth # auth 包，用来完成认证、授权功能
│   │   ├── authn.go # 认证功能
│   │   └── authz.go # 授权功能
│   ├── db # db 包，用来完成 MySQL 数据库连接
│   ├── token # JWT Token 的签发和解析
│   ├── util # 工具类包存放目录
│   │   └── id # id 包，用来生成唯一短 ID
│   └── version # version 包，用来保存 / 输出版本信息
├── README-en.md # 英文 README
├── README.md # 中文 README
├── scripts # 脚本文件
│   ├── boilerplate.txt # 指定版权头信息
│   ├── coverage.awk # awk 脚本，用来计算覆盖率
│   ├── make-rules # 子 Makefile 保存目录
│   │   ├── common.mk # 存放通用的 Makefile 变量
│   │   ├── golang.mk # 用来编译源码
│   │   └── tools.mk # 用来完成工具的安装
│   └── wrktest.sh # wrk 性能测试脚本
└── third_party # 第三方 Go 包存放目录
```


## 代码结构设计

miniblog 项目代码设计遵循简洁架构设计，一个简洁架构具有以下 5 个特性：

- 独立于框架： 该架构不会依赖于某些功能强大的软件库存在。这可以让你使用这样的框架作为工具，而不是让你的系统陷入到框架的约束中。

- 可测试性： 业务规则可以在没有 UI、数据库、Web 服务或其他外部元素的情况下进行测试，在实际的开发中，我们通过 Mock 来解耦这些依赖。

- 独立于UI ： 在无需改变系统其他部分的情况下，UI 可以轻松地改变。例如，在没有改变业务规则的情况下，Web UI 可以替换为控制台 UI。

- 独立于数据库： 你可以用 Mongo、Oracle、Etcd 或者其他数据库来替换 MariaDB，你的业务规则不要绑定到数据库。

- 独立于外部媒介： 实际上，你的业务规则可以简单到根本不去了解外部世界。

所以，基于这些约束，每一层都必须是独立的和可测试的。miniblog 代码架构分为 4 层：模型层（Model）、控制层（Controller）、业务层 （Biz）、仓库层（Store）。从控制层、业务层到仓库层，从左到右层级依次加深。模型层独立于其他层，可供其他层引用。代码架构如下图所示：

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/2da85aa1405945519dbf5f645587fb5f~tplv-k3u1fbpfcp-jj-mark:1512:0:0:0:q75.awebp" width="100%" />

层与层之间导入包时，都有严格的导入关系，这可以防止包的循环导入问题。导入关系如下：

- 模型层的包可以被仓库层、业务层和控制层导入。

- 控制层能够导入业务层和仓库层的包。这里需要注意，如果没有特殊需求，控制层要避免导入仓库层的包，控制层需要完成的业务功能都通过业务层来完成。这样可以使代码逻辑更加清晰、规范。

- 业务层能够导入仓库层的包。

### miniblog 四层架构

1. 模型层（Model）：模型层在有些软件架构中也叫做实体层（Entities），模型会在每一层中使用，在这一层中存储对象的结构和它的方法。

2. 控制层（Controller）：控制层接收 HTTP 请求，并进行参数解析、参数校验、逻辑分发处理、请求返回这些操作。控制层会将逻辑分发给业务层，业务层处理后返回，返回数据在控制层中被整合再加工，最终返回给请求方。控制层相当于实现了业务路由的功能。具体流程如下图所示：

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/cf1b54fea8bd4e2e9a059b4b10599787~tplv-k3u1fbpfcp-jj-mark:1512:0:0:0:q75.awebp" width="100%" /> 

3. 业务层 (Biz)：业务层主要用来完成业务逻辑处理，我们可以把所有的业务逻辑处理代码放在业务层。业务层会处理来自控制层的请求，并根据需要请求仓库层完成数据的 CURD 操作。业务层功能如下图所示：

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/89451e92ff3f4f888e56f16c638bbc55~tplv-k3u1fbpfcp-jj-mark:1512:0:0:0:q75.awebp" width="100%" />

1. 仓库层（Store)：仓库层用来跟数据库/第三方服务进行 CURD 交互，作为应用程序的数据引擎进行应用数据的输入和输出。这里需要注意，仓库层仅对数据库/第三方服务执行 CRUD 操作，不封装任何业务逻辑。这一层也会起到数据转换的作用：将从数据库/微服务中获取的数据转换为控制层、业务层能识别的数据结构，将控制层、业务层的数据格式转换为数据库或微服务能识别的数据格式。

### 层之间的通信

上面，我介绍了 miniblog 采用的 4 层结构，接下来我们再看看每一层之间是如何通信的。

除了模型层，控制层、业务层、仓库层之间都是通过接口进行通信的。通过接口通信，一方面可以使相同的功能支持不同的实现（也就是说具有插件化能力），另一方面也使得每一层的代码变得可测试。关于层通信我会在 第 12 节 详细介绍。



### miniblog 代码测试

控制层、业务层和仓库层之间是通过接口来通信的。通过接口通信有一个好处，就是可以让各层变得可测。那接下来，我们就来看下如何测试各层的代码。因为 第 18 节 ****会详细介绍如何测试 Go 代码，所以这里只介绍下测试思路。

- 模型层
  - 因为模型层不依赖其他任何层，我们只需要测试其中定义的结构及其函数和方法即可。

- 控制层
  - 控制层依赖于业务层，意味着该层需要业务层来支持测试。你可以通过 golang/mock 来 mock 业务层，测试用例可参考 TestPostController_Create。

- 业务层
  - 因为该层依赖于仓库层，意味着该层需要仓库层来支持测试。你可以通过 golang/mock 来 mock 仓库层，测试用例可以参考 Test_userBiz_List。

- 仓库层
  - 仓库层依赖于数据库，如果调用了其他微服务，那还会依赖第三方服务。我们可以通过 sqlmock 来模拟数据库连接，通过 httpmock 来模拟 HTTP 请求。