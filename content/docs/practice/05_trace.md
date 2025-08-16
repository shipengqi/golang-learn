---
title: Go 性能分析（下）
weight: 5
---

Go 提供了完善的性能分析工具：`pprof` 和 `trace`。

- `pprof` 主要适用于 CPU 占用、内存分配等资源的分析。
- `trace` 记录了程序运行中的行为，更适合于找出程序在一段时间内正在做什么。例如指定的 goroutine 在何时执行、执行了多长时间、什么时候陷入了堵塞、什么时候解除了堵塞、GC 如何影响了 goroutine 的执行。

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

![trace-home](https://raw.gitcode.com/shipengqi/illustrations/blobs/534089cf545c92c024a21a63d968b05ad23292b4/trace-home.png)

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

![view-trace](https://raw.gitcode.com/shipengqi/illustrations/blobs/0a5e63c1461ed7c1dc0250237f16dac8fc92a7df/view-trace.png)

1. 时间线：显示执行的时间。
2. Goroutines/Heap/Threads 的详细信息。
    - Goroutines：显示在执行期间的有多少个 goroutine 在运行，包含 GC 等待（GCWaiting）、可运行（Runnable）、 运行中（Running）这三种状态。
    - Heap：显示执行期间的内存分配和释放情况，包含当前堆使用量（Allocated）和下一次 GC 的阈值（NextGC）统计。
    - Threads：显示执行期间有多少个系统线程在运行，包含正在调用 SysCall（InSysCall）和运行中（Running）两种状态。
3. `PROCS`：每个 Processor 显示一行。默认显示系统内核数量，可以使用 `runtime.GOMAXPROCS(n)` 来控制数量。
    - `GC`：显示执行期间垃圾回收执行的次数和时间。**每次执行 GC，堆内存都会被释放一部分**。
    - 协程和事件：显示在每个虚拟处理器上有什么 goroutine 正在运行，而连线行为代表事件关联。

> 快捷键：w（放大），s（缩小），a（左移），d（右移）。

#### 查看某个时间点 goroutines 情况

![view-trace-g-counter](https://raw.gitcode.com/shipengqi/illustrations/blobs/66c698d9ef1f36f6671be54f76adf11f7c76972d/view-trace-g-counter.png)

图中正在运行的 goroutine 数量为 3，其他状态的 goroutine 数量都是 0。

#### 查看某个时间点堆的使用情况

![view-trace-heap](https://raw.gitcode.com/shipengqi/illustrations/blobs/a3eecdc08ebe49517599e4eeaf159b6610cb4369/view-trace-heap.png)

1. 红色部分表示已经占用的内存
2. 绿色部分的上边沿表示下次 GC 的目标内存，也就是绿色部分用完之后，就会触发 GC。

#### 查看某个时间点的系统线程

![view-trace-threads](https://raw.gitcode.com/shipengqi/illustrations/blobs/6d8792caf43219fec29ed801b6b1f38c2598b7fc/view-trace-threads.png)

图中正在运行的线程数量为 3，正在调用 SysCall 的线程数量为 0。

#### 查看 GC

![view-trace-gc](https://raw.gitcode.com/shipengqi/illustrations/blobs/81d595d195276ff92b76d757c0655fd68d0d514f/view-trace-gc.png)

可以选择多个查看统计信息。

#### 查看某个时间点的 goroutine 运行情况

![view-trace-folw-events](https://raw.gitcode.com/shipengqi/illustrations/blobs/7bc9ec23661e28a97187a01507cf301b499285d5/view-trace-folw-events.png)

点击具体的 goroutine 可以查看详细信息：

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

### goroutine analysis

进入 `Goroutine analysis` 可查看整个运行过程中，每个函数块有多少个 goroutine 在跑，并且观察每个的 goroutine 的运行开销都花费在哪个阶段。

![goroutines-analysis](https://raw.gitcode.com/shipengqi/illustrations/blobs/6987ab04e99c07e81ce40f765d8226b4cf520671/goroutines-analysis.png)

点击一个 goroutine 查看详细信息，例如 `main.main.func1`：

![goroutines-analysis-n](https://raw.gitcode.com/shipengqi/illustrations/blobs/36e800f4584a3c1b13be3e87291717339b0f310b/goroutines-analysis-n.png)

| 名称                    | 含义     | 耗时    |
|-----------------------|--------|-------|
| Execution Time        | 执行时间   | 983ms |
| Network Wait Time     | 网络等待时间 | 0ns   |
| Sync Block Time       | 同步阻塞时间 | 0ns   |
| Blocking Syscall Time | 调用阻塞时间 | 2ns   |
| Scheduler Wait Time   | 调度等待时间 | 194µs |
| GC Sweeping           | GC 清扫  | 0ns   | 
| GC Pause              | GC 暂停  | 14ms  |

## 查看 GC 的另一种方式

设置 `GODEBUG=gctrace=1`：

```bash

$ GODEBUG=gctrace=1 go run main.go 

gc 1 @0.059s 0%: 0.11+2.9+0.071 ms clock, 0.95+0.46/2.2/0+0.57 ms cpu, 4->4->1 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 2 @0.071s 1%: 0.17+1.9+0.031 ms clock, 1.3+0.46/1.5/1.1+0.24 ms cpu, 3->4->1 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 3 @0.086s 1%: 0.043+0.91+0.14 ms clock, 0.35+0/1.1/0.72+1.1 ms cpu, 3->3->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 4 @0.116s 1%: 0.047+2.3+0.009 ms clock, 0.38+0/2.0/0.40+0.075 ms cpu, 3->3->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
# command-line-arguments
gc 1 @0.005s 4%: 0.041+2.6+0.076 ms clock, 0.32+0.10/2.2/1.4+0.61 ms cpu, 5->5->4 MB, 5 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 2 @0.057s 1%: 0.036+2.4+0.097 ms clock, 0.29+0.23/2.7/0.31+0.78 ms cpu, 9->9->6 MB, 10 MB goal, 0 MB stacks, 0 MB globals, 8 P
Hello World!
```

### 格式

```
gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
```

- `gc #`：GC 执行次数的编号，每次叠加。例如 `gc 1`。
- `@#s`：自程序启动后到当前的具体秒数。
- `#%`：自程序启动以来在GC中花费的时间百分比。
- `#+...+#`：GC 的标记工作共使用的 CPU 时间占总 CPU 时间的百分比。
- `#->#-># MB`：分别表示 GC 启动时, GC 结束时, GC 活动时的堆大小.
- `#MB goal`：下一次触发 GC 的内存占用阈值。
- `#P`：当前使用的处理器 P 的数量。

示例：

```
gc 2 @0.057s 1%: 0.036+2.4+0.097 ms clock, 0.29+0.23/2.7/0.31+0.78 ms cpu, 9->9->6 MB, 10 MB goal, 0 MB stacks, 0 MB globals, 8 P
# gc 2：第 2 次 GC
# @0.057s 1%：当前是程序启动后的 0.057s。
# 1%：程序启动后到现在共花费 1% 的时间在 GC 上
# 0.036+2.4+0.097 ms clock：
#   0.036：表示单个 P 在 mark 阶段的 STW 时间。
#   2.4：表示所有 P 的 mark concurrent（并发标记）所使用的时间。
#   0.097：表示单个 P 的 markTermination 阶段的 STW 时间。
# 0.29+0.23/2.7/0.31+0.78 ms cpu：
#   0.29：表示整个进程在 mark 阶段 STW 停顿的时间。
#   0.23/2.7/0.31：0.23 表示 mutator assist 占用的时间，2.7 表示 dedicated + fractional 占用的时间，0.31 表示 idle 占用的时间。
#   0.78：0.78 表示整个进程在 markTermination 阶段 STW 时间。
# 9->9->6 MB：
#   9：表示开始 mark 阶段前的 heap_live 大小。
#   9：表示开始 markTermination 阶段前的 heap_live 大小。
#   6：表示被标记对象的大小。
# 10 MB goal：表示下一次触发 GC 回收的阈值是 10 MB。
# 0 MB stacks：表示本次 GC 没有任何栈上对象被移动。
# 0 MB globals：表示本次 GC 没有任何全局对象被移动。
# 8 P：本次 GC 一共涉及多少个 P。
```