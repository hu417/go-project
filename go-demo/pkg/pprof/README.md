
# pprof性能分析

pprof可以分析以下9中数据

| **Profile项** | **说明**    | **详情**                                                                        |
| ------------ | --------- | ----------------------------------------------------------------------------- |
| allocs       | 内存分配      | 从程序启动开始，分配的全部内存                                                               |
| block        | 阻塞        | 导致同步原语阻塞的堆栈跟踪                                                                 |
| cmdline      | 命令行调用     | 当前程序的命令行调用                                                                    |
| goroutine    | gorouting | 所有当前 goroutine 的堆栈跟踪                                                          |
| heap         | 堆         | 活动对象的内存分配抽样。您可以指定 gc 参数以在获取堆样本之前运行 GC                                         |
| mutex        | 互斥锁       | 争用互斥锁持有者的堆栈跟踪                                                                 |
| profile      | CPU分析     | CPU 使用率分析。可以在url中，通过seconds指定持续时间（默认30s）。获取配置文件后，使用 go tool pprof 命令分析CPU使用情况 |
| threadcreate | 线程创建      | 导致创建新操作系统线程的堆栈跟踪                                                              |
| trace        | 追踪        | 当前程序的执行轨迹。可以在url中，通过seconds指定持续时间（默认30s）。获取跟踪文件后，使用 go tool trace 命令调查跟踪      |


## 数据采集

### web采集

使用net/http/pprof进行采集
```go

package main

import (
    "fmt"
    "net/http"
    _ "net/http/pprof"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello world")
}

func main() {
    http.HandleFunc("/", HelloWorld)

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println(err)
    }
}


// 注意：net/http/pprof 当引入该包后，会自动触发该包的init()函数，进行注册路由
func init() {
    http.HandleFunc("/debug/pprof/", Index)
    http.HandleFunc("/debug/pprof/cmdline", Cmdline)
    http.HandleFunc("/debug/pprof/profile", Profile)
    http.HandleFunc("/debug/pprof/symbol", Symbol)
    http.HandleFunc("/debug/pprof/trace", Trace)
}


```

可以通过浏览器访问：http://127.0.0.1:8080/debug/pprof/
> allocs： 查看过去所有内存分配的样本，访问路径为 `$HOST/debug/pprof/allocs`。
> block： 查看导致阻塞同步的堆栈跟踪，访问路径为 `$HOST/debug/pprof/block`。
> cmdline： 当前程序的命令行的完整调用路径。
> goroutine：查看当前所有运行的 goroutines 堆栈跟踪，访问路径为 `$HOST/debug/pprof/goroutine`。
> heap: 查看活动对象的内存分配情况， 访问路径为 `$HOST/debug/pprof/heap`。
> mutex：查看导致互斥锁的竞争持有者的堆栈跟踪，访问路径为 `$HOST/debug/pprof/mutex`。
> profile： 默认进行 30s 的 CPU Profiling，得到一个分析用的 profile 文件，访问路径为 `$HOST/debug/pprof/profile`。
> threadcreate： 查看创建新 OS 线程的堆栈跟踪，访问路径为 `$HOST/debug/pprof/threadcreate`。

### 基准测试采集

testing 支持生成 CPU、memory 和 block 的 profile 文件
- cpuprofile=$FILE
- memprofile=$FILE
- blockprofile=$FILE

```go

// add.go
var datas []string
​
func add(str string) int {
    data := []byte(str)
    datas = append(datas, string(data))
    return len(datas)
}
​
// add_test.go
import "testing"
​
func TestAdd(t *testing.T) {
    _ = add("go pprof add text")
}
​
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        add("go pprof add text")
    }
}

```

以下命令分别对应CPU分析和内存分析
```go

go test -bench=. -cpuprofile=cpu.profile
go test -bench=. -memprofile=mem.profile

```



### 代码采集

一些命令行工具，执行完毕后直接退出的应用，这种不会常驻进程的，需要手动些代码进行采集
使用runtime/pprof 包进行程序运行时分析

```go

package main
​
import (
    "fmt"
    "os"
    "runtime/pprof"
)
​
func main() {
    cpuProfile, err := os.Create("./pprof/cpu_profile")
    if err != nil {
        fmt.Printf("创建文件失败:%s", err.Error())
        return
    }
    defer cpuProfile.Close()
​
    memProfile, err := os.Create("./pprof/mem_profile")
    if err != nil {
        fmt.Printf("创建文件失败:%s", err.Error())
        return
    }
    defer memProfile.Close()
    
    // 采集CPU信息
    pprof.StartCPUProfile(cpuProfile)
    defer pprof.StopCPUProfile()
​
    // 采集内存信息
    pprof.WriteHeapProfile(memProfile)
​
    // 逻辑代码
    for i := 0; i < 100; i++ {
        fmt.Println("pprof 工具型测试")
    }
}



```




## 数据分析

### web页面分析 
需要下载 graphviz,并配置到环境变量中

```go

// 检测是否安装好了graphviz
dot --version

// 直接分析profile文件 (-http是开启web端口浏览器访问)
go tool pprof -http=:8081 profile
​
// CPU耗时文件 （前提是使用的3.1的web采集的方式，这个方式会自动下载并分析文件，后面的这段url为当前服务的路由）
go tool pprof -http=:8081 http://127.0.0.1:8080/debug/pprof/profile
​
// 内存分配分析 (这个方式会自动下载并分析文件)
go tool pprof -http=:8081 http://127.0.0.1:8080/debug/pprof/allocs

```

访问访问web界面：http://localhost:8081/ui/
- localhost:8081/ui/top?si=cpu
- localhost:8081/ui/top?si=mem

### 终端分析

可以通过终端交互进行分析
```go

// 执行分析文件
go tool pprof cpu.pprof
​
go tool pprof  http://127.0.0.1:8080/debug/pprof/profile
​
go tool pprof  http://127.0.0.1:8080/debug/pprof/allocs


```