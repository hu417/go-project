
# GRPC

## 安装

- 安装protoc编译器
```go
// mac
brew install protobuf

// linux
sudo apt-get install protobuf-compiler

// win
https://github.com/protocolbuffers/protobuf/releases

```

- 安装go插件编译器
```go
// 从Proto文件(gRPC接口描述文件)生成go文件的编译器插件
go get -a google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

go env | grep "GOPATH"
go env | grep "GOROOT"
注:
① 此时会在你的GOPATH 的bin目录下生成可执行文件protoc-gen-go(protobuf的编译器插件).
② 执行protoc命令时,就会自动调用这个插件. // 需要将插件的路径添加到环境变量PATH中,或者GOROOT目录下
③ 可以使用不同的语言插件不同语言的文件.

```
- 安装grpc框架
```go
go mod init grpc.com
// gRPC运行时接口编解码支持库
go get -u google.golang.org/grpc

```

- 示例demo.proto
```go
syntax = "proto3";

package services;

// option go_package = "services/";

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```
生成Go代码
```go
protoc --go_out=./services/ --go-grpc_out=./services  ./pd/helloworld.proto

go mod tidy
```

## xx

