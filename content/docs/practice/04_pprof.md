---
title: Go 性能分析（上）
weight: 4
---

Go 提供的 pprof 工具可以用来做性能分析。pprof 可以读取分析样本的集合，并生成报告以可视化并帮助分析数据。

pprof 可以用于：

- CPU 分析（CPU Profiling）：按照一定的频率采集所监听的应用程序 CPU（含寄存器）的使用情况，可确定应用程序在主动消耗 CPU 周期
  时花费时间的位置。
- 内存分析（Memory Profiling）：在应用程序进行堆分配时记录堆栈跟踪，用于监视当前和历史内存使用情况，以及检查内存泄漏。
- 阻塞分析（Block Profiling）：记录 goroutine 阻塞等待同步（包括定时器通道）的位置。
- 互斥锁分析（Mutex Profiling）：报告互斥锁的竞争情况。

## 如何生成分析样本

生成分析样本的三种方式：

1. `runtime/pprof`：采集程序（**非 Server**）的运行数据。通过调用如 `runtime.StartCPUProfile`, `runtime.StopCPUProfile` 方法生成分析样本。主要用于**本地测试**。
   - [pkg/profile](https://github.com/pkg/profile) 封装了 `runtime/pprof`，使用起来更加简便。
2. `net/http/pprof`：采集 HTTP Server 的运行时数据，通过 HTTP 服务获取 Profile 分析样本，底层还是调用的 `runtime/pprof`。主要用于**服务器端测试**。
3. `go test -bench`：使用 `go test -bench=. -cpuprofile cpuprofile.out ...` 运行基准测试来生成分析样本，可以指定所需标识来进行数据采集。

以 `net/http/pprof` 为例：

```go
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // net/http/pprof 注册的是默认的 mux
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

`_ "net/http/pprof"` 这行代码会自动添加 `/debug/pprof` 的路由。程序运行后，访问 http://localhost:8080/debug/pprof 就可以查看分析样本。


## 如何查看分析报告

打开 http://localhost:8080/debug/pprof 后会看到下面页面：

![pprof-home](https://raw.gitcode.com/shipengqi/illustrations/blobs/affd23f900f0b1fe9de04d91f201dd0f377bca8c/pprof-home.png)

pprof 包括了几个子页面：

- alloc: 查看所有内存分配的情况
- block（Block Profiling）：`<ip:port>/debug/pprof/block`，查看导致阻塞同步的堆栈跟踪
- cmdline : 当前程序的命令行调用
- goroutine：`<ip:port>/debug/pprof/goroutine`，查看当前所有运行的 goroutines 堆栈跟踪。
- heap（Memory Profiling）: `<ip:port>/debug/pprof/heap`，查看活动对象的内存分配情况，在获取堆样本之前，可以指定 gc GET 参数来运行 gc。
- mutex（Mutex Profiling）: `<ip:port>/debug/pprof/mutex`，查看导致互斥锁竞争的持有者的堆栈跟踪。
- profile（CPU Profiling）: `<ip:port>/debug/pprof/profile`， 默认进行 `30s` 的 CPU Profiling，可以设置 GET 参数 `seconds` 来指定持续时间。获取跟踪文件之后，使用 `go tool pprof` 命令来分析。
- threadcreate：`<ip:port>/debug/pprof/threadcreate`，查看创建新 OS 线程的堆栈跟踪。
- trace: 当前程序的执行轨迹。可以设置 GET 参数 `seconds` 来指定持续时间。获取跟踪文件之后，使用 `go tool trace` 命令来分析。

### 在 Web 查看分析报告

上面有三种方式生成分析样本，这里以 `net/http/pprof` 为例。

#### 下载分析样本

点击 profile，等待 30s 后会下载 CPU profile 文件，或者执行命令 `go tool pprof http://localhost:8080/debug/pprof/profile` ，得到的输出中有一行

```shell
Saved profile in C:\Users\shipeng.CORPDOM\pprof\pprof.samples.cpu.002.pb.gz
```

表示生成的 profile 文件路径。

#### 查看分析报告

执行 `go tool pprof -http=<port> <profile 文件>` 启动 web server，然后就可以访问 `http://localhost:8081` 来查看：

```sh
$ go tool pprof -http=:8081 profile
Serving web UI on http://localhost:8081
```

或者输入 `web`，会在浏览器打开一个 svg 图片：
```shell
$ go tool pprof profile
$ (pprof) web
```

> 如果出现 `Could not execute dot; may need to install graphviz.`，那么需要安裝 Graphviz。

![profile-graph](https://raw.gitcode.com/shipengqi/illustrations/blobs/d9de0b23bc2c8ab7b03fa2b63764108edb03d559/profile-graph.png)

图中框越大，线越粗代表它占用 CPU 的时间越长。

点击 `View -> Flame Graph` 可以查看火焰图：

![profile-flame-graph](https://raw.gitcode.com/shipengqi/illustrations/blobs/cdeb3c01593b09e5c58cdae12703d662efbd19e3/profile-flame-graph.png)

图中调用顺序由上到下，每一块代表一个函数，越大代表占用 CPU 的时间越长。

还可以查看 Top，Peek，Source 等。能够更方便、更直观的看到 Go 应用程序的调用链、使用情况等。

### 在终端查看分析报告

使用 `go tool pprof` 命令可以在交互式终端查看分析报告。`go tool pprof` 可以直接从 HTTP 服务获取分析样本，也可以指定本地样本文件，例如 `go tool pprof cpu.pprof`。

#### CPU Profiling

执行 60s 的 CPU Profiling：

```sh
$ go tool pprof http://localhost:8080/debug/pprof/profile?seconds=60
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=10
Saved profile in C:\Users\shipeng.CORPDOM\pprof\pprof.samples.cpu.001.pb.gz
Type: cpu
Time: Nov 18, 2019 at 11:08am (CST)
Duration: 10.20s, Total samples = 10.03s (98.38%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
```

输入 `top 10`：

```sh
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

- `flat`：当前函数上的运行耗时
- `flat%`：当前函数上的 CPU 运行耗时总比例
- `sum%`：当前函数上累积使用 CPU 总比例
- `cum`：当前函数加上它之上的调用运行总耗时
- `cum%`：当前函数加上它之上的调用的 CPU 运行耗时总比例
- 最后一列为函数名称

#### Heap Profiling

Heap Profiling 支持四种内存概况的分析：

- `inuse_space`：分析程序常驻内存的占用
- `alloc_objects`：分析程序临时分配的内存
- `inuse_objects`：查看函数对应的对象的数量
- `alloc_space`：查看函数分配的内存空间大小

默认就是 `inuse_space`：

```shell
# 默认就是 inuse_space，-inuse_space 可以忽略
$ go tool pprof -inuse_space http://localhost:8080/debug/pprof/heap
Saved profile in C:\Users\shipeng\pprof\pprof.___go_build_github_com_shipengqi_example_v1_advance_go_pprof.exe.alloc_objects.alloc_space.inuse_objects.inuse_space.002.pb.gz
Type: inuse_space
Time: Dec 6, 2023 at 2:05pm (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
```

输入 `top`：

```sh
(pprof) top
Showing nodes accounting for 42464.20kB, 100% of 42464.20kB total
      flat  flat%   sum%        cum   cum%
41952.16kB 98.79% 98.79% 41952.16kB 98.79%  main.Add (inline)
  512.04kB  1.21%   100%   512.04kB  1.21%  unicode/utf16.Encode
         0     0%   100%   512.04kB  1.21%  internal/poll.(*FD).Write
         0     0%   100%   512.04kB  1.21%  internal/poll.(*FD).writeConsole
         0     0%   100%   512.04kB  1.21%  log.(*Logger).output
         0     0%   100%   512.04kB  1.21%  log.Println (inline)
         0     0%   100% 42464.20kB   100%  main.main.func1
         0     0%   100%   512.04kB  1.21%  os.(*File).Write
         0     0%   100%   512.04kB  1.21%  os.(*File).write (inline)
```

输入 `traces` 查看 goroutines 占用内存的大小：

```sh
(pprof) traces
...
Type: inuse_space
Time: Dec 6, 2023 at 2:45pm (CST)
-----------+-------------------------------------------------------
         0   main.Add (inline)
             main.main.func1
-----------+-------------------------------------------------------
     bytes:  2.25MB
    2.28MB   main.Add (inline)
             main.main.func1
-----------+-------------------------------------------------------
     bytes:  1.80MB
    1.85MB   main.Add (inline)
             main.main.func1
-----------+-------------------------------------------------------
     bytes:  1.43MB
         0   main.Add (inline)
             main.main.func1
-----------+-------------------------------------------------------
     bytes:  1.14MB
         0   main.Add (inline)
             main.main.func1
-----------+-------------------------------------------------------

```

#### goroutine Profiling

```sh
$ go tool pprof http://localhost:8080/debug/pprof/goroutine
Saved profile in C:\Users\shipeng\pprof\pprof.___go_build_github_com_shipengqi_example_v1_advance_go_pprof.exe.goroutine.001.pb.gz
...
(pprof)
```

输入 `traces` 会输出所有 goroutines 的调用栈信息，可以很方便的查看整个调用链。

```sh
(pprof) traces
...
-----------+-------------------------------------------------------
         1   runtime.cgocall
             syscall.SyscallN
             syscall.Syscall6
             syscall.WriteConsole
             internal/poll.(*FD).writeConsole
             internal/poll.(*FD).Write
             os.(*File).write (inline)
             os.(*File).Write
             log.(*Logger).output
             log.Println (inline)
             main.main.func1
-----------+-------------------------------------------------------
         1   runtime.goroutineProfileWithLabels
             runtime/pprof.runtime_goroutineProfileWithLabels
             runtime/pprof.writeRuntimeProfile
             runtime/pprof.writeGoroutine
             runtime/pprof.(*Profile).WriteTo
             net/http/pprof.handler.ServeHTTP
             net/http/pprof.Index
             net/http.HandlerFunc.ServeHTTP
             net/http.(*ServeMux).ServeHTTP
             net/http.serverHandler.ServeHTTP
             net/http.(*conn).serve
```

调用栈的顺序是**自下而上**的。

#### Block 和 Mutex Profiling

Block 和 Mutex Profiling 都需要在代码中调用 `runtime` 包的方法进行设置：

```go
package main

import "runtime"

func main() {
	// Rate 小于 0，则不采集
    runtime.SetBlockProfileRate(1)
	// Fraction 小于 0，则不采集
    runtime.SetMutexProfileFraction(1)
    // ...
}
```

然后使用 `go tool pprof` 分析，输入 `top` 查看排名，`list <func>` 可以查看具体的信息。

#### 对比

当需要查看不同时间段的差异时，可以使用 `-base` 参数来对比两个 profile 文件。

```shell
$ go tool pprof -base <profile1> <profile2>
```


