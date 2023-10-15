---
title: Go 性能分析
weight: 4
---

# Go 性能分析

PProf 是 Go 提供的用于可视化和分析性能分析数据的工具。

- `runtime/pprof`：采集程序（非 Server）的运行数据进行分析
- `net/http/pprof`：采集 HTTP Server 的运行时数据进行分析

主要可以用于：

- CPU Profiling：CPU 分析，按照一定的频率采集所监听的应用程序 CPU（含寄存器）的使用情况，可确定应用程序在主动消耗 CPU 周期
时花费时间的位置。
- Memory Profiling：内存分析，在应用程序进行堆分配时记录堆栈跟踪，用于监视当前和历史内存使用情况，以及检查内存泄漏。
- Block Profiling：阻塞分析，记录 goroutine 阻塞等待同步（包括定时器通道）的位置。
- Mutex Profiling：互斥锁分析，报告互斥锁的竞争情况。

## 性能分析

### 分析 HTTP Server

#### Web

```go
import (
 "log"
 "net/http"
 _ "net/http/pprof"
)

var datas []string

func Add(str string) string {
 data := []byte(str)
 sData := string(data)
 datas = append(datas, sData)

 return sData
}

func main() {
 go func() {
  for {
   log.Println(Add("https://github.com/shipengqi"))
  }
 }()

 _ = http.ListenAndServe("0.0.0.0:8080", nil)
}
```

注意要引入 `_ "net/http/pprof"`，这样程序运行以后，就会自动添加 `/debug/pprof` 的路由，可以
访问 `ttp://127.0.0.1:8080/debug/pprof/`。

![profile](../imgs/profile1.png)

- alloc: 查看所有内存分配的情况
- block（Block Profiling）：`$HOST/debug/pprof/block`，查看导致阻塞同步的堆栈跟踪
- cmdline : 当前程序的命令行调用
- goroutine：`$HOST/debug/pprof/goroutine`，查看当前所有运行的 goroutines 堆栈跟踪
- heap（Memory Profiling）: `$HOST/debug/pprof/heap`，查看活动对象的内存分配情况，在获取堆样本之前，可以指定 gc GET 参数来运行 gc。
- mutex（Mutex Profiling）: `$HOST/debug/pprof/mutex`，查看导致互斥锁的竞争持有者的堆栈跟踪
- profile: `$HOST/debug/pprof/profile`， 默认进行 30s 的 CPU Profiling，可以 GET 参数 `seconds` 中指定持续时间。
获得 profile 文件之后，使用 `go tool pprof` 命令分析 profile 文件。
- threadcreate：`$HOST/debug/pprof/threadcreat`e，查看创建新 OS 线程的堆栈跟踪
- trace: 当前程序的执行轨迹。可以在 GET 参数 `seconds` 中指定持续时间。获取跟踪文件之后，使用 `go tool trace` 命令来分析。

#### 交互式终端

```sh
# seconds 可以调整等待的时间，当前命令设置等待 60 秒后会进行 CPU Profiling
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=60

Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=10
Saved profile in C:\Users\shipeng.CORPDOM\pprof\pprof.samples.cpu.001.pb.gz
Type: cpu
Time: Nov 18, 2019 at 11:08am (CST)
Duration: 10.20s, Total samples = 10.03s (98.38%)
Entering interactive mode (type "help" for commands, "o" for options)

# 进入交互式命令模式
(pprof) top 10
Showing nodes accounting for 9.54s, 95.11% of 10.03s total
Dropped 73 nodes (cum <= 0.05s)
Showing top 10 nodes out of 14
      flat  flat%   sum%        cum   cum%
     9.42s 93.92% 93.92%      9.46s 94.32%  runtime.cgocall
     0.02s   0.2% 94.12%      9.62s 95.91%  internal/poll.(*FD).writeConsole
     0.02s   0.2% 94.32%      9.81s 97.81%  log.(*Logger).Output
     0.02s   0.2% 94.52%      0.10s     1%  log.(*Logger).formatHeader
     0.02s   0.2% 94.72%      0.06s   0.6%  main.Add
     0.02s   0.2% 94.92%      9.50s 94.72%  syscall.Syscall6
     0.01s   0.1% 95.01%      0.07s   0.7%  runtime.systemstack
     0.01s   0.1% 95.11%      9.51s 94.82%  syscall.WriteConsole
         0     0% 95.11%      0.07s   0.7%  fmt.Sprintln
         0     0% 95.11%      9.69s 96.61%  internal/poll.(*FD).Write
```

上面的输出：

- `flat`：给定函数上运行耗时
- `flat%`：同上的 CPU 运行耗时总比例
- `sum%`：给定函数累积使用 CPU 总比例
- `cum`：当前函数加上它之上的调用运行总耗时
- `cum%`：同上的 CPU 运行耗时总比例
- 最后一列为函数名称

```sh
go tool pprof http://localhost:6060/debug/pprof/heap

Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in C:\Users\shipeng.CORPDOM\pprof\pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.008.pb.gz
Type: inuse_space
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 837.48MB, 100% of 837.48MB total
      flat  flat%   sum%        cum   cum%
  837.48MB   100%   100%   837.48MB   100%  main.main.func1

# 其他分析
go tool pprof http://localhost:6060/debug/pprof/block

go tool pprof http://localhost:6060/debug/pprof/mutex
```

- `-inuse_space`：分析应用程序的常驻内存占用情况
- `-alloc_objects`：分析应用程序的内存临时分配情况

## PProf 可视化界面

`data.go`：

```go
package pdata

var datas []string

func Add(str string) string {
 data := []byte(str)
 sData := string(data)
 datas = append(datas, sData)

 return sData
}
```

`data_test.go`：

```go
package pdata

import "testing"

const url = "https://github.com/"

func TestAdd(t *testing.T) {
 s := Add(url)
 if s == "" {
  t.Errorf("Test.Add error!")
 }
}

func BenchmarkAdd(b *testing.B) {
 for i := 0; i < b.N; i++ {
  Add(url)
 }
}
```

运行基准测试：

```sh
# 下面的命令会生成 cprof 文件, 使用 go tool pprof 分析
go test -bench . -cpuprofile=cprof
goos: windows
goarch: amd64
pkg: github.com/shipengqi/golang-learn/demos/pprof/pdata
BenchmarkAdd-8          10084636               143 ns/op
PASS
ok      github.com/shipengqi/golang-learn/demos/pprof/pdata     2.960s
```

启动可视化界面：

```sh
$ go tool pprof -http=:8080 cpu.prof

# 或者
$ go tool pprof cpu.prof
$ (pprof) web
```

如果出现 `Could not execute dot; may need to install graphviz.`，参考 "安裝 Graphviz"

![](../imgs/profile2.png)

上图中的框越大，线越粗代表它消耗的时间越长。

![](../imgs/profile3.png)

![](../imgs/profile4.png)

PProf 的可视化界面能够更方便、更直观的看到 Go 应用程序的调用链、使用情况等。

火焰图：
![](../imgs/profile5.png)

### 安裝 Graphviz

官网 [下载地址](http://www.graphviz.org/download/)

### 配置环境变量

将 bin 目录添加到 Path 环境变量中，如 `C:\Program Files (x86)\Graphviz2.38\bin`。

### 验证

```sh
dot -version
```

**部分内容来自** [Go 大杀器之性能剖析 PProf](https://github.com/EDDYCJY/blog/blob/7b021d0dee/tools/go-tool-pprof.md)
