---
title: Go 数据竞争检测器
weight: 2
---

# Go 数据竞争检测器

数据竞争是并发系统中最常见，同时也最难处理的 Bug 类型之一。数据竞争会在两个 Go 程并发访问同一个变量， 且至少有一个访问为写入时产生。

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
    }
    m["2"] = "b" // 第二个冲突的访问
	<-c
    for k, v := range m {
        fmt.Println(k, v)
    }
}

```

## 数据竞争检测器

Go内建了数据竞争检测器。要使用它，请将 `-race` 标记添加到 go 命令之后：

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
// +build !race

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

## 使用

竞争检测器只会寻找在运行时发生的竞争，因此它不能在未执行的代码路径中寻找竞争。若你的测试并未完全覆盖，你可以在实际的工作负载下运行通过 `-race` 编译的二进制程序，以此寻找更多的竞争。

### 典型的数据竞争

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println(i)  // 你要找的不是 i。
			wg.Done()
		}()
	}
	wg.Wait()
}
```

此函数字面中的变量 i 与该循环中使用的是同一个变量， 因此该Go程中对它的读取与该递增循环产生了竞争。（此程序通常会打印55555，而非01234。） 此程序可通过创建该变量的副本来修复。

```go
func main() {
 var wg sync.WaitGroup
 wg.Add(5)
 for i := 0; i < 5; i++ {
  go func(j int) {
   fmt.Println(j)  // 很好。现在读取的是该循环计数器的局部副本。
   wg.Done()
  }(i)
 }
 wg.Wait()
}
```

#### 偶然被共享的变量

```go
// ParallelWrite 将数据写入 file1 和 file2 中，并返回一个错误。
func ParallelWrite(data []byte) chan error {
 res := make(chan error, 2)
 f1, err := os.Create("file1")
 if err != nil {
  res <- err
 } else {
  go func() {
   // 此处的 err 是与主Go程共享的，
   // 因此该写入操作就会与下面的写入操作产生竞争。
   _, err = f1.Write(data)
   res <- err
   f1.Close()
  }()
 }
 f2, err := os.Create("file2")  // 第二个冲突的对 err 的写入。
 if err != nil {
  res <- err
 } else {
  go func() {
   _, err = f2.Write(data)
   res <- err
   f2.Close()
  }()
 }
 return res
}
```

#### 不受保护的全局变量

若以下代码在多个Go程中调用，就会导致 service 映射产生竞争。 对映射的并发读写是不安全的：

```go
var service map[string]net.Addr

func RegisterService(name string, addr net.Addr) {
 service[name] = addr
}

func LookupService(name string) net.Addr {
 return service[name]
}
```

要保证此代码的安全，需通过互斥锁来保护对它的访问：

```go
var (
 service   map[string]net.Addr
 serviceMu sync.Mutex
)

func RegisterService(name string, addr net.Addr) {
 serviceMu.Lock()
 defer serviceMu.Unlock()
 service[name] = addr
}

func LookupService(name string) net.Addr {
 serviceMu.Lock()
 defer serviceMu.Unlock()
 return service[name]
}
```

#### 不受保护的基原类型变量

数据竞争同样会发生在基原类型的变量上（如 bool、int、 int64 等），就像下面这样：

```go
type Watchdog struct { last int64 }

func (w *Watchdog) KeepAlive() {
 w.last = time.Now().UnixNano()  // 第一个冲突的访问。
}

func (w *Watchdog) Start() {
 go func() {
  for {
   time.Sleep(time.Second)
   // 第二个冲突的访问。
   if w.last < time.Now().Add(-10*time.Second).UnixNano() {
    fmt.Println("No keepalives for 10 seconds. Dying.")
    os.Exit(1)
   }
  }
 }()
}
```

甚至“无辜”的数据竞争也会导致难以调试的问题：(1) 非原子性的内存访问 (2) 编译器优化的干扰以及 (3) 进程内存访问的重排序问题。

对此，典型的解决方案就是使用信道或互斥锁。要保护无锁的行为，一种方法就是使用 sync/atomic 包。

```go
type Watchdog struct{ last int64 }

func (w *Watchdog) KeepAlive() {
 atomic.StoreInt64(&w.last, time.Now().UnixNano())
}

func (w *Watchdog) Start() {
 go func() {
  for {
   time.Sleep(time.Second)
   if atomic.LoadInt64(&w.last) < time.Now().Add(-10*time.Second).UnixNano() {
    fmt.Println("No keepalives for 10 seconds. Dying.")
    os.Exit(1)
   }
  }
 }()
}
```

## 运行时开销

竞争检测的代价因程序而异，但对于典型的程序，内存的使用会增加5到10倍， 而执行时间会增加2到20倍。
