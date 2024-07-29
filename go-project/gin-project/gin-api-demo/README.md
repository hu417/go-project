

# gin-api后端框架

## 项目结构
编码规范: https://juejin.cn/post/7157594175846744071
RESTful API: https://juejin.cn/post/7227444929948319805
ddd: github.com/dddplayer/

```bash

mkdir -p gin-api-demo
go mod init github.com/go/gin/gin-api-demo
```


## 依赖注入
1、参考: https://github.com/jassue/gin-wire

```bash
Godotenv: github.com/joho/godotenv
Gin: https://github.com/gin-gonic/gin
Gorm: https://github.com/go-gorm/gorm
Wire: https://github.com/google/wire
Viper: https://github.com/spf13/viper
Zap: https://github.com/uber-go/zap
Jwt: github.com/golang-jwt/jwt/v5
Golang-jwt: https://github.com/golang-jwt/jwt
Go-redis: https://github.com/go-redis/redis
Testify: https://github.com/stretchr/testify
Sonyflake: https://github.com/sony/sonyflake
Gocron: https://github.com/go-co-op/gocron
Go-sqlmock: https://github.com/DATA-DOG/go-sqlmock
Gomock: https://github.com/golang/mock
Swaggo: https://github.com/swaggo/swag
Prometheus: github.com/prometheus/client_golang
TimeOut: github.com/gin-contrib/timeout


```


## 常见分层下的 error 处理
// 参考: https://mp.weixin.qq.com/s/SnaurQfXDVidrl_ihBQtDA
github.com/pkg/errors 提供了很多实用的函数，例如：
- Wrap(err error, message string) error：该函数基于原始错误 err，返回一个带有堆栈跟踪信息和附加信息 message 的新 error
- Wrapf(err error, format string, args ...interface{}) error: 和上面的函数功能是一样的，只不过可以对附加信息进行格式化封装
- WithMessage(err error, message string) error：该函数基于原始错误 err，返回一个附加信息 message 的新 error
- WithMessagef(err error, format string, args ...interface{}) error: 和上面的函数功能是一样的，只不过可以对附加信息进行格式化封装
- Cause(err error) error：该函数用于提取 err 中的原始 error，它会递归地检查 error，直到找到最底层的原始 error，如果存在的话


以典型的 MVC ( dao → service → controller/middleware) 分层结构举例，常见的错误处理大致如下：
```bash
// controller / middleware  
res, err := service.GetById(ctx, id)  
if err != nil {  
  // 获取根因，打印堆栈等信息
  log.Errorf(ctx, "service.GetById failed, original error: %T %v", errors.Cause(err), errors.Cause(err))  
  log.Errorf(ctx, "stack trace: \n%+v\n", err)  
······  
}  
······  
  
// service  
article, err := dao.GetById(ctx, id)  
if err != nil {  
  // 2、附带额外信息
  return errors.WithMessage(err, "dao.GetById failed")  
}  
······  
  
// dao  
······  
if err != nil {  
  // 1、初始化堆栈，附带额外信息
  return errors.Wrapf(err, "GetById failed, id=%s, error=%v", id, err)  
}  
······


当在 Dao 层遇到原始错误 Original Error 后，使用 errors.Wrap() 对错误进行封装。这个封装操作可以在保留根因（Origin error）的同时，提供堆栈信息，并添加额外的上下文信息，然后将封装后的错误传递给上一层处理。

当 service 层接收到 error 之后，使用 errors.WithMessage() 函数，将额外的信息附加到错误上，并继续将错误向上层传递，直至到达 controller 层。在 controller 层，我们可以打印出根因的类型、信息以及堆栈信息，以便更好地进行问题排查。
```
