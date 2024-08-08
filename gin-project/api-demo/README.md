

## cobra

Cobra 是一个 Go 语言开发的命令行（CLI）框架，它提供了简洁、灵活且强大的方式来创建命令行程序

Cobra分为cobra和cobra-cli两部分 cobra-cli用来初始化项目文件以及生成新的命令文件 cobra用来执行项目中生成的命令



```go
go install github.com/spf13/cobra-cli@v1.3.0 

// 初始化项目
GOWORK=off cobra-cli init --license=MIT .
```