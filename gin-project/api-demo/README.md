

参考：https://juejin.cn/post/7395218310399295523

## db

```bash
create database `gin-base` charset utf8mb4;
show processlist;

INSERT INTO `admin` (`id`, `name`, `password`, `created_at`, `updated_at`) VALUES
(1, '管理员', 'NOEFNIE', '2024-07-15 15:12:26', '2024-07-15 16:01:08');



```

## cobra

Cobra 是一个 Go 语言开发的命令行（CLI）框架，它提供了简洁、灵活且强大的方式来创建命令行程序

Cobra分为cobra和cobra-cli两部分 cobra-cli用来初始化项目文件以及生成新的命令文件 cobra用来执行项目中生成的命令



```go
go install github.com/spf13/cobra-cli@v1.3.0 

// 初始化项目
GOWORK=off cobra-cli init --license=MIT .
```

## air

Air 是一个 Go 语言开发的热重载工具，它允许你无需重启服务器即可在代码更改后自动重新编译和重启应用程序。
```go

$ go install github.com/air-verse/air
$ air -v
```