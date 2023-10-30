---
title: Go 数据竞争检测器
weight: 2
---

# Go 数据竞争检测器

数据竞争是并发系统中最常见，同时也最难处理的 Bug 类型之一。数据竞争会在两个 goroutine 并发访问同一个变量，且至少有一个访问为写入时产生。

这个数据竞争的例子可导致程序崩溃和内存数据损坏（memory corruption）。

```go
package main

import "fmt"

func main() {
	c := make(chan bool)
	m := make(map[string]string)
	go func() {
        m["1"] = "a"  // 第一个冲突的访问
		c <- true
    }()
    m["2"] = "b" // 第二个冲突的访问
	<-c
    for k, v := range m {
        fmt.Println(k, v)
    }
}
```

运行 `go run -race ./main.go` 或者 `go build -race ./main.go` 编译后再运行会抛出类似的错误：

```
==================
WARNING: DATA RACE
Write at 0x00c00010a090 by goroutine 7:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  main.main.func1()
      /root/workspace/main.go:9 +0x4a

Previous write at 0x00c00010a090 by main goroutine:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  main.main()
      /root/workspace/main.go:12 +0x108

Goroutine 7 (running) created at:
  main.main()
      /root/workspace/main.go:8 +0xeb
==================
2 b
1 a
Found 1 data race(s)
```

## 数据竞争检测器

Go 内建了数据竞争检测器。要使用它，请将 `-race` 标记添加到 go 命令之后：

```bash
go test -race mypkg    // 测试该包
go run -race mysrc.go  // 运行其源文件
go build -race mycmd   // 构建该命令
go install -race mypkg // 安装该包
```

### 选项

`GORACE` 环境变量可以设置竞争检测的选项：

```bash
GORACE="option1=val1 option2=val2"
```

选项：

- `log_path`（默认为 `stderr`）：竞争检测器会将其报告写入名为 `log_path.pid` 的文件中。特殊的名字 `stdout` 和 `stderr` 会将报告分别写入到标准输出和标准错误中。
- `exitcode`（默认为 66）：当检测到竞争后使用的退出状态。
- `strip_path_prefix`（默认为 ""）：从所有报告文件的路径中去除此前缀， 让报告更加简洁。
- `history_size`（默认为 1）：每个 Go 程的内存访问历史为 `32K * 2**history_size` 个元素。增加该值可避免在报告中避免 "failed to restore the stack"（栈恢复失败）的提示，但代价是会增加内存的使用。
- `halt_on_error`（默认为 0）：控制程序在报告第一次数据竞争后是否退出。

例如：

```bash
GORACE="log_path=/tmp/race/report strip_path_prefix=/my/go/sources/" go test -race
```

### 编译标签

可以通过编译标签来排除某些竞争检测器下的代码/测试：

```go
//go:build !race

package foo

// 此测试包含了数据竞争。见123号问题。
func TestFoo(t *testing.T)  {
 // ...
}

// 此测试会因为竞争检测器的超时而失败。
func TestBar(t *testing.T)  {
 // ...
}

// 此测试会在竞争检测器下花费太长时间。
func TestBaz(t *testing.T)  {
 // ...
}
```

## 运行时开销

竞争检测器只会寻找在运行时发生的竞争，因此它不能在未执行的代码路径中寻找竞争。若你的测试并未完全覆盖，你可以运行通过 `-race` 编译的二进制程序，以此寻找更多的竞争。

竞争检测的代价因程序而异，但对于典型的程序，内存的使用会增加 5 到 10 倍， 而执行时间会增加 2 到 20 倍。
