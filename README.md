IDE:

  VSCode:
  
    mac通过终端code 命令打开vscode

    ```bash
    $ cat ~/.bash_profile 
    alias code="/Users/hellokitty/Downloads/Visual\ Studio\ Code.app/Contents/Resources/app/bin/code"
    $ source ~/.bash_profile 
    ```
前端:

  项目实战: https://gitee.com/zhengqingya/java-developer-document/tree/master/%E7%9F%A5%E8%AF%86%E5%BA%93/%E5%89%8D%E7%AB%AF/03-%E2%98%86%E5%AE%9E%E6%88%98%E9%A1%B9%E7%9B%AE%E2%98%86/web%E9%A1%B9%E7%9B%AE%E5%AE%9E%E6%88%98/02-%E3%80%90%E7%AC%AC%E4%BA%8C%E7%89%88%E3%80%91vue3+vite4/small-web/doc

Go 学习路线

  基础学习:
  
    参考学习: https://github.com/hu417/go-study/

  项目结构:
    ddd: https://jishuzhan.net/article/1721803979219800066
    
  gin框架:
  
    官方文档: https://gin-gonic.com/zh-cn/docs/
    参考博客: https://www.mszlu.com/go/gin/01/01.html
    示例项目: 
      博客: https://www.bilibili.com/video/BV1f24y1G72C
      Todo List: https://github.com/borntodie-new/todo-list-backup
      从零开发企业级 Go 应用: https://www.golangblogs.com/books/go-zero-cloud
      casbin: https://github.com/kubesre/go-easy-admin
      k8s dashboard: https://github.com/kubesre/genbu
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
