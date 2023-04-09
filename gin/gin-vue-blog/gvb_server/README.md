
gin-vue-blog
    文檔:
    https://www.wolai.com/fengfeng/brbaVLeTzLxegLKpuqWy6w
    wxid_6516/wxid_6516
    視頻:
    https://www.bilibili.com/video/BV1f24y1G72C?p=1&vd_source=e8ffec33a64e60ef7e166f307c32e001

项目结构
```bash
api         接口目录
bin         可执行文件
config      服务配置的结构体目录
core        初始化操作
docs        swag生成的api文件目录
flag        命令行相关的初始化
global      全局变量的包
middleware  中间件
models      表结构
routers     gin路由目录
service     项目与服务有关的目录,如settings.yaml  配置文件
test        测试文件目录
utils/valid  常用的工具目录
.gitignore   git忽略文件
main.go     程序入口


```

mkdir -p api config core docs flag global middleware models routers service testdata utils
touch settings.yaml
go mod init gvb_server

## 配置读取
1、编写配置文件信息
   setting.yaml
2、定义配置参数结构体
   config/enter.go,conf_mysql.go
3、定义全局变量进行存储
   global/global.go
4、读取配置文件
   core/conf.go
5、测试
   main.go

## zap日志初始化
1、安装依赖
    go get -u go.uber.org/zap
    go get -u gopkg.in/natefinch/lumberjack.v2
2、添加配置
    config/enter.go
    config/conf_logger.go
    service/settings.yaml

    core/logger.go
    global/global.go

3、初始化
    bin/appStart.sh

## gorm配置
1、安装依赖
    go get gorm.io/gorm
2、初始化配置
    config/enter.go
    config/conf_mysql.go
    core/gorm.go
    service/setting.yaml
3、定义全局变量
    global/global.go
4、测试
    main.go


## router路由初始化
1、安装依赖
    go get -u github.com/gin-gonic/gin
2、初始化配置
    system
        config/enter.go
        config/conf_system.go
        service/setting.yaml
    routers
        enter.go
3、定义全局变量
    global/global.go
4、初始化路由
    service/appStart.go
4、测试
    main.go

## router路由配置
// 示例: 系统设置settings
1、定义响应
    api/settings_api/enter.go settings_info.go
    api/enter.go
2、router封装
    router/enter.go
    router/settings_router.go

## 响应封装
1、对响应数据的封装
    models/res/response.go
    models/res/err_code.go // 状态码的定义
2、重新修改api
    api/settings_api/settings_info.go

## 错误状态码的封装测试
1、定义错误状态码
    models/res/err_code.json
2、测试
    test/err_code_json_test.go


## elasticsearch配置
1、安装依赖库
    go get -u github.com/olivere/elastic/v7
2、配置es连接参数
    gvb_server/config/enter.go conf_elastic.go
    gvb_server/core/es.go
    gvb_server/service/settings.yaml
3、全局配置
    gvb_server/global/global.go
4、启动
    gvb_server/bin/appStart.go