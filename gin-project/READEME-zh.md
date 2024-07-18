
# 笔记

## ssh

```bash
# 生成 SSH Key
ssh-keygen -t github -C "Github SSH Key"
ssh-keygen -t rsa -b 4096 -C 'Github SSH Key' -f ~/.ssh/id_rsa_4096_github -P ''

ls ~/.ssh/
私钥文件 id_rsa_4096_github
公钥文件 id_rsa_4096_github.pub

# 读取公钥文件
cat ~/.ssh/id_rsa_4096_github.pub

# 测试
ssh -T git@gitee.com
```

## README

可以使用 readme.so 工具来协助生成 README 文件，通常 README 文件需要包含以下部分：Features、Installation、Usage/Examples、Documentation、Feedback、Contributing、Authors、License、Related。

## gitignore

在当前目录添加一个 .gitignore 文件，里面包含不期望 Git 跟踪的文件，例如：临时文件等

可以使用生成工具 gitignore.io 来生成 .gitignore

## Git

安装/配置

```bash
$ cd /tmp
$ wget --no-check-certificate https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.36.1.tar.gz
$ tar -xvzf git-2.36.1.tar.gz
$ cd git-2.36.1/
$ ./configure
$ make
$ sudo make install
$ git --version # 输出 git 版本号，说明安装成功
git version 2.36.1

$ tee -a $HOME/.bashrc <<'EOF'
# Configure for git
export PATH=/usr/local/libexec/git-core:$PATH
EOF


# 配置 Git
$ git config --global user.name "Lingfei Kong"    # 用户名改成自己的
$ git config --global user.email "colin404@qq.com"    # 邮箱改成自己的
$ git config --global credential.helper store    # 设置 git，保存用户名和密码
$ git config --global core.longpaths true # 解决 Git 中 'Filename too long' 的错误


# 初始化仓库
$ git init # 初始化当前目录为 Git 仓库
$ git add . # 添加所有被 Git 追踪的文件到暂存区
$ git status -s # 查看仓库当前文件提交状态
$ git commit -m "feat: 第一次提交" # 将暂存区内容添加到本地仓库中
$ git push https://gitee.com/***/test.git  # 将本地的Git仓库信息推送上传到服务器
$ git log  # 查看git提交的日志

# 添加远程仓库
git remote add origin  仓库地址
# 查看当前仓库对应的远程仓库地址
git remote -v
# 修改远程仓库地址
git remote set-url origin 仓库地址

```

## go

安装/配置
```bash

$ wget -P /tmp/ https://go.dev/dl/go1.19.4.linux-amd64.tar.gz
$ mkdir -p $HOME/go
$ tar -xvzf /tmp/go1.19.4.linux-amd64.tar.gz -C $HOME/go
$ mv $HOME/go/go $HOME/go/go1.19.4

$ tee -a $HOME/.bashrc <<'EOF'
# Go envs
export GOVERSION=go1.19.4 # Go 版本设置
export GO_INSTALL_DIR=$HOME/go # Go 安装目录
export GOROOT=$GO_INSTALL_DIR/$GOVERSION # GOROOT 设置
export GOPATH=$WORKSPACE/golang # GOPATH 设置
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH # 将 Go 语言自带的和通过 go install 安装的二进制文件加入到 PATH 路径中
export GO111MODULE="on" # 开启 Go moudles 特性
export GOPROXY=https://goproxy.cn,direct # 安装 Go 模块时，代理服务器设置
export GOPRIVATE=
export GOSUMDB=off # 关闭校验 Go 依赖包的哈希值
EOF


$ bash # 配置 `$HOME/.bashrc` 后，需要执行 `bash` 命令将配置加载到当前 Shell
$ go version
go version go1.19.4 linux/amd64


# 初始化工作区
$ mkdir -p $GOPATH && cd $GOPATH
$ go work init
$ go env GOWORK # 执行此命令，查看 go.work 工作区文件路径
/home/goer/workspace/golang/go.work

# 格式化代码
$ gofmt -s -w ./
```

### Protobuf 编译环境安装

安装 protoc 命令
```bash
$ cd /tmp/
$ wget https://github.com/protocolbuffers/protobuf/releases/download/v21.9/protobuf-cpp-3.21.9.tar.gz
$ tar -xvzf protobuf-cpp-3.21.9.tar.gz
$ cd protobuf-3.21.9/
# $ libtoolize --automake --copy --debug –force
$ ./autogen.sh

$ ./autogen.sh
$ ./configure
$ make
$ sudo make install
$ protoc --version # 查看 protoc 版本，成功输出版本号，说明安装成功
libprotoc 3.21.9

```

安装 protoc-gen-go 命令
```bash
# 添加 -x 参数打印具体的安装细节
$ go install -x github.com/golang/protobuf/protoc-gen-go@latest

```

### 热加载

安装 air 工具
```bash
$ go install github.com/cosmtrek/air@latest

```

配置 air 工具
```bash
# Config file for [Air](https://github.com/cosmtrek/air) in TOML format

# Working directory
# . or absolute path, please note that the directories following must be under root.
root = "."
tmp_dir = "tmp"

[build]
# Array of commands to run before each build
pre_cmd = []
# Just plain old shell command. You could use `make` as well.

#####  只需要写你平常编译使用的 shell 命令,你也可以使用 `make`  #####
# cmd = "go build -o ./tmp/app ./cmd/main.go"
cmd = "make build"
# Array of commands to run after ^C
post_cmd = []
# Binary file yields from `cmd`.

##### 由 `cmd` 命令得到的二进制文件名. ######
bin = "_output/miniblog"
# Customize binary, can setup environment variables when run your app.
full_bin = "APP_ENV=dev APP_USER=air ./tmp/main"
# Add additional arguments when running binary (bin/full_bin). Will run './tmp/main hello world'.
args_bin = ["hello", "world"]
# Watch these filename extensions.
include_ext = ["go", "tpl", "tmpl", "html"]
# Ignore these filename extensions or directories.
exclude_dir = ["assets", "tmp", "vendor", "frontend/node_modules"]
# Watch these directories if you specified.
include_dir = []
# Watch these files.
include_file = []
# Exclude files.
exclude_file = []
# Exclude specific regular expressions.
exclude_regex = ["_test\\.go"]
# Exclude unchanged files.
exclude_unchanged = true
# Follow symlink for directories
follow_symlink = true
# This log file places in your tmp_dir.
log = "air.log"
# Poll files for changes instead of using fsnotify.
poll = false
# Poll interval (defaults to the minimum interval of 500ms).
poll_interval = 500 # ms
# It's not necessary to trigger build each time file changes if it's too frequent.
delay = 0 # ms
# Stop running old binary when build errors occur.
stop_on_error = true
# Send Interrupt signal before killing process (windows does not support this feature)
send_interrupt = false
# Delay after sending Interrupt signal
kill_delay = 500 # nanosecond
# Rerun binary or not
rerun = false
# Delay after each execution
rerun_delay = 500

[log]
# Show log time
time = false
# Only show main log (silences watcher, build, runner)
main_only = false

[color]
# Customize each part's color. If no color found, use the raw app log.
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"

[misc]
# Delete tmp directory on exit
clean_on_exit = true

[screen]
clear_on_rebuild = true
keep_scroll = true

# Enable live-reloading on the browser.
[proxy]
  enabled = true
  proxy_port = 8090
  app_port = 8080
```

### swagger

在线编辑器: https://editor-next.swagger.io/ 

安装
```bash
$ go install github.com/go-swagger/go-swagger/cmd/swagger@latest

$ swagger serve -F=swagger --no-open --port 65534 ./api/openapi/openapi.yaml
2022/11/22 21:19:49 serving docs at http://localhost:65534/docs

```

### 版权声明

如果你的项目是一个开源项目或者未来准备开源，那么还需要给项目添加一些版权声明，主要包括：

1. 存放在项目根目录下的 LICENSE 文件，用来声明项目遵循的开源协议；

2. 项目源文件中的版权头信息，用来说明文件所遵循的开源协议。

添加版权声明的第一步就是选择一个开源协议

#### 添加 LICENSE 文件

一般项目的根目录下会存放一个 LICENSE 文件用来声明本项目所遵循的协议，所以我们这里也要为 miniblog 初始化一个 LICENSE 文件。LICENSE 文件怎么写？不用慌，作为一个懒惰的程序员，我们可以使用 license 工具来生成，具体操作命令如下：

```bash
$ go install github.com/nishanths/license/v5@latest
$ license -list # 查看支持的代码协议
$ license -n 'colin404(孔令飞) <colin404@foxmail.com>' -o LICENSE mit # 在 miniblog 项目根目录下执行
$ ls LICENSE 
LICENSE
```

#### 给源文件添加版本声明

给项目中的源文件添加版权头信息，用来声明文件所遵循的开源协议。miniblog 的版权头信息保存在 boilerplate.txt 文件中。
> 提示：版权头信息保存的文件名，通常命名为：boilerplate。

安装 addlicense 工具
```bash
$ go install github.com/marmotedu/addlicense@latest

```

运行 addlicense 工具添加版权头信息
```bash
$ addlicense -v -f ./scripts/boilerplate.txt --skip-dirs=third_party,vendor,_output .
cmd/miniblog/main.go

```

可以看到 main.go 文件已经被追加上了版权头信息，内容如下：
```bash
// Copyright 2022 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package main

import "fmt"

// Go 程序的默认入口函数(主函数).
func main() {
    fmt.Println("Hello MiniBlog!")
}
```

### 构建工具

选择 make 作为构建工具

#### 编写简单的 Makefile

通过以下方式来学习 Makefile 编程：

- 学习 Makefile 基本语法：可参考文档 [Makefile基础知识.md](https://github.com/marmotedu/geekbang-go/blob/master/makefile/Makefile%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86.md)；

- 学习 Makefile 高级语法（如果有时间/感兴趣）：陈皓老师编写的 跟我一起写 Makefile (PDF 重制版) 。

编写后的 Makefile 文件位于项目根目录下，内容为：
```bash
# ==============================================================================
# 定义全局 Makefile 变量方便后面引用

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 获取项目根目录绝对路径
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(ROOT_DIR)/_output

# ==============================================================================
# 定义 Makefile all 伪目标，执行 `make` 时，会默认会执行 all 伪目标
.PHONY: all
all: add-copyright format build

# ==============================================================================
# 定义其他需要的伪目标

.PHONY: build
build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
        @go build -v -o $(OUTPUT_DIR)/miniblog $(ROOT_DIR)/cmd/miniblog/main.go

.PHONY: format
format: # 格式化 Go 源码.
        @gofmt -s -w ./

.PHONY: add-copyright
add-copyright: # 添加版权头信息.
        @addlicense -v -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=third_party,vendor,$(OUTPUT_DIR)

.PHONY: swagger
swagger: # 启动 swagger 在线文档.
        @swagger serve -F=swagger --no-open --port 65534 $(ROOT_DIR)/api/openapi/openapi.yaml

.PHONY: tidy
tidy: # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
        @go mod tidy

.PHONY: clean
clean: # 清理构建产物、临时文件等. 实现幂等删除
        @-rm -vrf $(OUTPUT_DIR)

```


## DB

### mysql



### MariaDB

登录数据库并创建 miniblog 用户
```bash
$ mysql -h127.0.0.1 -P3306 -uroot -p'miniblog1234' # 连接 MariaDB，-h 指定主机，-P 指定监听端口，-u 指定登录用户，-p 指定登录密码
MariaDB [(none)]> grant all on miniblog.* TO miniblog@127.0.0.1 identified by  'miniblog1234'  ; Query OK, 0 rows affected (0.000 sec)
MariaDB [(none)]> flush privileges; Query OK, 0 rows affected (0.000 sec)

```







## 1


## 2