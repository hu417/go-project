

参考：https://github.com/hbinr/XS-bbs/blob/main/README.md



目录结构
```bash
├── cmd                # 程序入口
│   ├── main.go
│   ├── wire_gen.go    # 已删除
│   └── wire.go        # 已删除
├── docs               # swagger接口文档
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal           # 私有模块，业务代码和业务严重依赖的库
│   ├── app            # app 项目，按功能模块划分，方便后续扩展微服务
│   └── pkg            # 业务严重依赖的公共库
├── pkg                # 公共模块，和业务无关，可以对外使用的库
│   ├── cache          # 缓存初始化封装
│   ├── conf           # 配置定义及初始化封装
│   ├── database       # 数据库初始化封装
│   ├── logger         # 日志库初始化封装
│   ├── servers        # http 路由初识化、注册相关,后续可以支持 grpc server 
│   └── utils          # 一些工具封装
├── config.yaml        # 配置文件
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── script             # 脚本文件
    └── app.sql

```