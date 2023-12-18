---
title: Go 性能分析（下）
weight: 5
---

# Go 性能分析（下）

Go 提供了完善的性能分析工具：`pprof` 和 `trace`。

- `pprof` 主要适用于 CPU 占用、内存分配等资源的分析。
- `trace` 记录了程序运行中的行为，更适合于找出程序在一段时间内正在做什么。例如指定的 goroutine
  在何时执行、执行了多长时间、什么时候陷入了堵塞、什么时候解除了堵塞、GC 如何影响了 goroutine 的执行。

## 如何生成分析样本

生成 Trace 分析样本的方式主要有三种：

**1. 使用 `runtime/trace` 标准库来生成**：

```go
package main

import (
	"os"
	"runtime/trace"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	ch := make(chan string)

	go func() {
		ch <- "hello"
	}()
	// read from channel
	<-ch
}

```

执行程序就可以生成跟踪文件 `trace.out`：

```
go run main.go
```

**2. 使用 `net/http/pprof` 来生成**，查看 [Go 性能分析](../04_pprof)。

**3. 使用 `go test -trace` 来生成**，例如 `go test -trace trace.out demo_test.go`。

## 如何查看分析报告

使用上面例子生成的 `trace.out` 文件，运行下面的命令：

```
$ go tool trace trace.out
2019/11/18 15:17:28 Parsing trace...
2019/11/18 15:17:28 Splitting trace...
2019/11/18 15:17:28 Opening browser. Trace viewer is listening on http://127.0.0.1:59181
```

访问 http://127.0.0.1:59181，可以看到类似的界面：

![trace-home](https://raw.githubusercontent.com/shipengqi/illustrations/f700fde7450ccdb21c410611ff949769e301c364/go/trace-home.png)

- `View trace`：查看所有 goroutines 的执行过程。
- `Goroutine analysis`：goroutines 分析，查看具体的 goroutine 的信息。
- `Network blocking profile`：网络阻塞概况。
- `Synchronization blocking profile`：同步阻塞概况。
- `Syscall blocking profile`：系统调用阻塞概况。
- `Scheduler latency profile`：调度延迟的概况，可以调度在哪里最耗费时间。
- `User defined tasks`：用户自定义任务。
- `User defined regions`：用户自定义区域。
- `Minimum mutator utilization`：最低 Mutator 利用率。

> Network/Sync/Syscall blocking profile 是分析锁竞争的最佳选择。

### View trace

进入 `View trace` 页面：

![view-trace](https://raw.githubusercontent.com/shipengqi/illustrations/f700fde7450ccdb21c410611ff949769e301c364/go/view-trace.png)

1. 时间线：显示执行的时间。
2. Goroutines/Heap/Threads 的详细信息。
    - Goroutines：显示在执行期间的有多少个 goroutine 在运行，包含 GC 等待（GCWaiting）、可运行（Runnable）、 运行中（Running）这三种状态。
    - Heap：显示执行期间的内存分配和释放情况，包含当前堆使用量（Allocated）和下一次 GC 的阈值（NextGC）统计。
    - Threads：显示执行期间有多少个系统线程在运行，包含正在调用 SysCall （InSysCall）和运行中（Running）两种状态。
3. `PROCS`：每个 Processor 显示一行。默认显示系统内核数量，可以使用 `runtime.GOMAXPROCS(n)` 来控制数量。
    - `GC`：显示执行期间垃圾回收执行的次数和时间。**每次执行 GC，堆内存都会被释放一部分**。
    - 协程和事件：显示在每个虚拟处理器上有什么 Goroutine 正在运行，而连线行为代表事件关联。

> 快捷键：w（放大），s（缩小），a（左移），d（右移）。

#### 查看某个时间点 goroutines 情况

![view-trace-g-counter](https://raw.githubusercontent.com/shipengqi/illustrations/f700fde7450ccdb21c410611ff949769e301c364/go/view-trace-g-counter.png)

图中正在运行的 goroutine 数量为 3，其他状态的 goroutine 数量都是 0。

#### 查看某个时间点堆的使用情况

![view-trace-heap](https://raw.githubusercontent.com/shipengqi/illustrations/aa2cb61b817e520379277b90937ac21909b4abd5/go/view-trace-heap.png)

1. 红色部分表示已经占用的内存
2. 绿色部分的上边沿表示下次 GC 的目标内存，也就是绿色部分用完之后，就会触发 GC。

#### 查看某个时间点的系统线程

![view-trace-threads](https://raw.githubusercontent.com/shipengqi/illustrations/aa2cb61b817e520379277b90937ac21909b4abd5/go/view-trace-threads.png)

图中正在运行的线程数量为 3，正在调用 SysCall 的线程数量为 0。

#### 查看 GC

![view-trace-gc](https://raw.githubusercontent.com/shipengqi/illustrations/aa2cb61b817e520379277b90937ac21909b4abd5/go/view-trace-gc.png)

可以选择多个查看统计信息。

#### 查看某个时间点的 goroutine 运行情况

![view-trace-folw-events](https://raw.githubusercontent.com/shipengqi/illustrations/aa2cb61b817e520379277b90937ac21909b4abd5/go/view-trace-folw-events.png)

点击具体的 Goroutine 可以查看详细信息：

- `Start`：开始时间
- `Wall Duration`：持续时间
- `Self Time`：执行时间
- `Start Stack Trace`：开始时的堆栈信息
- `End Stack Trace`：结束时的堆栈信息
- `Incoming flow`：输入流
- `Outgoing flow`：输出流
- `Preceding events`：之前的事件
- `Following events`：之后的事件
- `All connected`：所有连接的事件

点击 `Flow events` 选择 `All`，可以查看程序运行中的事件流情况。

### Goroutine analysis

进入 `Goroutine analysis` 可查看整个运行过程中，每个函数块有多少个 goroutine 在跑，并且观察每个的 goroutine 的运行开销都花费在哪个阶段。

![goroutines-analysis](https://raw.githubusercontent.com/shipengqi/illustrations/aa2cb61b817e520379277b90937ac21909b4abd5/go/goroutines-analysis.png)

点击一个 goroutine 查看详细信息，例如 `main.main.func1`：

![goroutines-analysis-n](https://raw.githubusercontent.com/shipengqi/illustrations/aa2cb61b817e520379277b90937ac21909b4abd5/go/goroutines-analysis-n.png)

| 名称                    | 含义     | 耗时    |
|-----------------------|--------|-------|
| Execution Time        | 执行时间   | 983ms |
| Network Wait Time     | 网络等待时间 | 0ns   |
| Sync Block Time       | 同步阻塞时间 | 0ns   |
| Blocking Syscall Time | 调用阻塞时间 | 2ns   |
| Scheduler Wait Time   | 调度等待时间 | 194µs |
| GC Sweeping           | GC 清扫  | 0ns   | 
| GC Pause              | GC 暂停  | 14ms  |
