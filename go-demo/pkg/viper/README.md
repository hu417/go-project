

# Viper

安装依赖: `go get github.com/spf13/viper`

## 基础

Viper是适用于Go应用程序（包括Twelve-Factor App）的完整配置解决方案
支持以下特性：
- 设置默认值
- 从JSON、TOML、YAML、HCL、envfile和Java properties格式的配置文件读取配置信息
- 实时监控和重新读取配置文件（可选）
- 从环境变量中读取
- 从远程配置系统（etcd或Consul）读取并监控配置变化
- 从命令行参数读取配置
- 从buffer读取配置
- 显式配置值

读取优先级:
- 显式调用Set设置值
- 命令行参数（flag）
- 环境变量
- 配置文件
- key/value存储
- 默认值
重要： 目前Viper配置的键（Key）是大小写不敏感的。

## 把值存入Viper

### 设置默认值
```go

viper.SetDefault("ContentDir", "content")
viper.SetDefault("LayoutDir", "layouts")
viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})

```

### 读取配置文件

Viper支持JSON、TOML、YAML、HCL、envfile和Java properties格式的配置文件
```go
/* 方式一 */
viper.SetConfigName("config") // 配置文件名称(无扩展名);如果找到多个，则报错
// viper.SetConfigType("yaml") // 如果配置文件的名称中没有扩展名，则需要配置此项,针对远程配置
viper.AddConfigPath("/etc/appname/")   // 查找配置文件所在的路径
viper.AddConfigPath("$HOME/.appname")  // 多次调用以添加多个搜索路径
viper.AddConfigPath(".")               // 还可以在工作目录中查找配置

/* 方式二 */
// 读取配置文件
// viper.SetConfigFile("config.json")

// 查找并读取配置文件
if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
        // 配置文件未找到错误；如果需要可以忽略
    } else {
        // 配置文件被找到，但产生了另外的错误
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }
}

// 配置文件找到并成功解析

```

### 写入配置文件
存储在运行时所做的所有修改
相关命令
- WriteConfig - 将当前的viper配置写入预定义的路径并覆盖（如果存在的话）。如果没有预定义的路径，则报错。
- SafeWriteConfig - 将当前的viper配置写入预定义的路径。如果没有预定义的路径，则报错。如果存在，将不会覆盖当前的配置文件。
- WriteConfigAs - 将当前的viper配置写入给定的文件路径。将覆盖给定的文件(如果它存在的话)。
- SafeWriteConfigAs - 将当前的viper配置写入给定的文件路径。不会覆盖给定的文件(如果它存在的话)。
  - 根据经验，标记为safe的所有方法都不会覆盖任何文件，而是直接创建（如果不存在），而默认行为是创建或截断

```go

viper.WriteConfig() // 将当前配置写入“viper.AddConfigPath()”和“viper.SetConfigName”设置的预定义路径
viper.SafeWriteConfig()
viper.WriteConfigAs("/path/to/my/.config")
viper.SafeWriteConfigAs("/path/to/my/.config") // 因为该配置文件写入过，所以会报错
viper.SafeWriteConfigAs("/path/to/my/.other_config")

```

### 监控并重新读取配置文件
viper可以在运行时读取配置文件的更新
```go

// 需要先配置读取文件的操作

viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
  // 配置文件发生变更之后会调用的回调函数
	fmt.Println("Config file changed:", e.Name)
})

```

### 从io.Reader读取配置

```go

viper.SetConfigType("yaml") // 或者 viper.SetConfigType("YAML")

// 任何需要将此配置添加到程序中的方法。
var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)

viper.ReadConfig(bytes.NewBuffer(yamlExample))

viper.Get("name") // 这里会得到 "steve"

```

### 覆盖设置
```go

viper.Set("Verbose", true)
viper.Set("LogFile", LogFile)

```

### 注册和使用别名

别名允许多个键引用单个值
```go

viper.RegisterAlias("loud", "Verbose")  // 注册别名（此处loud和Verbose建立了别名）

viper.Set("verbose", true) // 结果与下一行相同
viper.Set("loud", true)   // 结果与前一行相同

viper.GetBool("loud") // true
viper.GetBool("verbose") // true

```

### 使用环境变量

使用ENV变量时，务必要意识到Viper将ENV变量视为区分大小写

### 使用Flags

ackages/morphling/star_user_v2/src/pages/index
ackages/morphling/star_user_v2/src/pages/inde
### eee

## fff