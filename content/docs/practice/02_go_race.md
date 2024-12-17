---
title: Go 数据竞争检测器
weight: 2
---

数据竞争是并发系统中最常见，同时也最难处理的 Bug 类型之一。数据竞争会在两个 goroutine 并发访问同一个变量，且至少有一个访问为写入时产生。

下面是一个会导致程序崩溃的例子：

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

运行 `go run -race ./main.go` 程序会马上崩溃：

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

Go 内置了数据竞争检测器。使用时将 `-race` 标记添加到 go 命令之后：

```bash
go test -race mypkg    // 测试该包
go run -race mysrc.go  // 运行其源文件
go build -race mycmd   // 构建该命令
go install -race mypkg // 安装该包
```

### 选项

Go 提供的 `GORACE` 环境变量可以用来设置竞争检测器的选项，格式为 `GORACE="option1=val1 option2=val2"`

支持的选项：

- `log_path`（默认为 `stderr`）：竞争检测器会将报告写入名为 `<log_path>.pid` 的文件中。如果值为 `stdout` 或 `stderr` 时会将报告分别写入到标准输出和标准错误中。
- `exitcode`（默认为 66）：检测到竞争后使用的退出状态码。
- `strip_path_prefix`（默认为 ""）：从所有报告文件的路径中去除此前缀，使报告更加简洁。
- `history_size`（默认为 1）：每个 Go 程序的内存访问历史为 `32K * 2**history_size` 个元素。增加该值可以在报告中避免 "failed to restore the stack" 的提示，但代价是会增加内存的使用。
- `halt_on_error`（默认为 0）：控制程序在报告第一次数据竞争后是否退出。

例如：

```bash
GORACE="log_path=/tmp/race/report strip_path_prefix=/my/go/sources/" go test -race
```

### 编译标签

如果某些代码不需要被竞争检测器检查，可以通过编译标签来排除：

```go
//go:build !race

package foo
```

## 运行时开销

竞争检测器只会寻找在运行时发生的竞争，因此它不能在未执行的代码路径中寻找竞争。如果你的测试覆盖率比较低，可以通过 `go build -race` 来编译，来寻找更多的竞争。

竞争检测的代价因程序而异，但对于典型的程序，内存的使用会增加 5 到 10 倍， 而执行时间会增加 2 到 20 倍。
