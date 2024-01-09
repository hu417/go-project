### 此仓库包含go web框架学习的项目

Go 学习路线

  基础学习:
  
    参考学习: https://github.com/hu417/go-study/
    
  gin框架:
  
    官方文档: https://gin-gonic.com/zh-cn/docs/
    参考博客: https://www.mszlu.com/go/gin/01/01.html
    示例项目: https://www.bilibili.com/video/BV1f24y1G72C

  gorm框架(mysql)
  
    参考博客: https://www.mszlu.com/go/gorm/01/01.html

  xorm框架(mysql):
  
    官方文档: http://xorm.topgoer.com/
    参考博客: 

  go-redis框架:
  
    参考博客: https://www.mszlu.com/go/go-redis/01/01.html

  go-zero框架:
  
    入门视频: 
    官方文档: https://go-zero.dev/cn/docs/introduction
    示例项目: 
    - https://www.bilibili.com/video/BV1cr4y1s7H4?p=1
    - https://github.com/zhoushuguang/beyond
    微服务项目: https://juejin.cn/post/7036011047391592485
    进阶学习: https://www.bilibili.com/video/BV1LS4y1U72n
    进阶项目: https://github.com/Mikaelemmmm/go-zero-looklook/tree/main/doc/chinese


Gin框架模版

```bash
.
├── conf (配置文件)
│   └── config.yaml
├── controller（控制类）
│   ├── admin.go
│   ├── auth.go
│   ├── base.go
│   ├── blog.go
│   └── controller.go
├── dao （数据库连接）
│   └── mysql.go
├── log（日志文件）
├── logger（zap logger工具类）
│   └── logger.go
├── main.go （入口文件）
├── models（GORM 访问数据库）
│   ├── base.go
│   ├── category.go
│   ├── comment.go
│   ├── config.go
│   ├── post.go
│   ├── response.go
│   └── user.go
├── routers （路由）
│   └── routers.go
├── settings （文件配置管理）
│   └── settings.go
├── static （静态文件）
├── templates（页面文件）
└── util（工具类）
    ├── RediStore.go
    ├── SessionStore.go
    ├── functions.go
    ├── localtime.go
    └── pager.go

```



