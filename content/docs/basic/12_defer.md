---
title: defer
weight: 12
---

`defer` 语句一般被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。因为 `defer` 可以保证让你更任何情况下，资源都会被释放。所在的 goroutine 发生 `panic` 时依然可以执行。

**注意：** 调用 `os.Exit` 时 `defer` 不会被执行。

`defer` 延迟函数为什么是按照**代码中的出现的顺序逆序执行**的？

因为 Go 的底层实现中每个 goroutine 的 `g` 对象上都有一个 `defer` 链表，当有新的 `_defer` 结构体就会挂在到这个链表的头部。因此执行顺序是**先进后出**。

```go
type _defer struct {
    started   bool    // defer 语句是否已经执行
    heap      bool    // 区分对象是在堆上分配还是栈上分配
    sp        uintptr // 调用方的 sp (栈底) 寄存器
    pc        uintptr // 调用方的 pc (程序计数器) 寄存器，下一条汇编指令的地址
    fn        func()  // 传入 defer 的函数，包括函数地址及参数
    _panic    *_panic // 正在执行 defer 的 panic 对象
    link      *_defer // next defer, 链表指针，可以指向栈或者堆
}
```

由于每个 goroutine 有自己的 `defer` 链表，因此 `defer` 是无法跨 goroutine 的。

## defer 初始化


编译器不仅将 `defer` 关键字都转换成 `runtime.deferproc` 函数，还会为所有调用 `defer` 的函数末尾插入 `runtime.deferreturn` 的函数调用。

- `runtime.deferproc` 负责创建新的延迟调用；
- `runtime.deferreturn` 负责在函数调用结束时执行所有的延迟调用；

`runtime.deferproc` 会为 `defer` 创建一个新的 `runtime._defer` 结构体、设置它的函数指针 `fn`、程序计数器 `pc` 和栈指针 `sp` 并将相关的参数拷贝到相邻的内存空间中：

```go
func deferproc(siz int32, fn *funcval) {
	sp := getcallersp()
	argp := uintptr(unsafe.Pointer(&fn)) + unsafe.Sizeof(fn)
	callerpc := getcallerpc() // 调用 deferproc 的函数的程序计数器

	d := newdefer(siz)
	if d._panic != nil {
		throw("deferproc: d.panic != nil after newdefer")
	}
	d.fn = fn // 传入 defer 的函数，包括函数地址及参数
	d.pc = callerpc // 调用 deferproc 的函数的程序计数器
	d.sp = sp // 调用 deferproc 的函数的栈指针

	switch siz {
	case 0:
	case sys.PtrSize:
		*(*uintptr)(deferArgs(d)) = *(*uintptr)(unsafe.Pointer(argp))
	default:
		memmove(deferArgs(d), unsafe.Pointer(argp), uintptr(siz))
	}

	return0()
}
```

最后调用的 `runtime.return0` 是唯一一个不会触发延迟调用的函数，它可以避免递归调用 `runtime.deferreturn` 函数。

`runtime.newdefer` 的作用是获得一个 `runtime._defer` 结构体，有三种方式：

- 从调度器的延迟调用缓存池 `sched.deferpool` 中取出结构体并将该结构体追加到当前 goroutine 的缓存池中；
- 从 goroutine 的延迟调用缓存池 `pp.deferpool` 中取出结构体；
- 通过 `runtime.mallocgc` 在堆上创建一个新的结构体；

无论使用哪种方式，只要获取到 `runtime._defer` 结构体，它都会被追加到所在 `g._defer` 链表的最前面。这也是为什么后调用的 `defer` 会优先执行。

`runtime.deferreturn` 会从 goroutine 的 `_defer` 链表中取出最前面的 `runtime._defer` 结构体并调用 `runtime.jmpdefer` 函数传入需要执行的函数和参数：

```go
func deferreturn(arg0 uintptr) {
	gp := getg()
	d := gp._defer
	if d == nil {
		return
	}
	sp := getcallersp()
	// ...

	switch d.siz {
	case 0:
	case sys.PtrSize:
		*(*uintptr)(unsafe.Pointer(&arg0)) = *(*uintptr)(deferArgs(d))
	default:
		memmove(unsafe.Pointer(&arg0), deferArgs(d), uintptr(d.siz))
	}
	fn := d.fn
	gp._defer = d.link
	freedefer(d)
	jmpdefer(fn, uintptr(unsafe.Pointer(&arg0)))
}
```

`runtime.jmpdefer` 是一个用汇编语言实现的运行时函数，它的主要工作是跳转到 `defer` 所在的代码段并在执行结束之后跳转回 `runtime.deferreturn`。

### 总结

整个过程可以简单概括为：

1. 函数遇到 `defer` 关键字时，调用 `runtime.deferproc` 函数创建一个新的 `runtime._defer` 结构体并将该结构体追加到所在 goroutine 的 `_defer` 链表的最前面；
2. 然后在函数的末尾插入 `runtime.deferreturn` 函数；
3. 调用 `runtime.deferreturn` 函数时，会从所在 goroutine 的 `_defer` 链表中取出最前面的 `runtime._defer` 结构体并调用 `runtime.jmpdefer` 函数传入需要执行的函数和参数；
4. `runtime.jmpdefer` 函数会跳转到 `defer` 所在的代码段并在执行结束之后跳转回 `runtime.deferreturn`；


## 开放编码

在 1.14 中通过开发编码（Open Coded）实现 `defer` 关键字，该设计使用代码内联优化 `defer` 关键字的额外开销，该优化可以将 `defer` 的调用开销从 1.13 版本的 ~35ns 降低至 ~6ns 左右：

通过静态分析和代码转换，将部分 defer 调用从动态链表管理转为静态化处理。

开发编码只会在满足以下的条件时启用：

- 函数的 `defer` 数量少于或者等于 8 个；
- 函数的 `defer` 关键字不能在循环中执行；
- 未使用 `recover`（需精确控制 `defer` 执行顺序时禁用优化）。

开发编码在函数返回前直接插入 `defer` 函数调用，省去链表遍历。

## defer 的性能

相比直接用 CALL 汇编指令调用函数，`defer` 要花费更大代价，包括注册，调用操作，额为的缓存开销。

```go
func call () {
  m.Lock()
  m.Unlock()
}

func deferCall()  {
  m.Lock()
  defer m.Unlock()
}

func BenchmarkCall(b *testing.B)  {
  for i := 0; i < b.N; i ++ {
    call()
  }
}


func BenchmarkDeferCall(b *testing.B)  {
  for i := 0; i < b.N; i ++ {
    deferCall()
  }
}
```

```sh
$ go test -bench=.
goos: windows
goarch: amd64
pkg: github.com/shipengqi/golang-learn/demos/defers
BenchmarkCall-8         92349604                12.9 ns/op
BenchmarkDeferCall-8    34305316                36.3 ns/op
PASS
ok      github.com/shipengqi/golang-learn/demos/defers  2.571s

```

性能相差三倍，尽量避免使用 `defer`。

## 什么时候不应该使用 defer

比如处理日志文件，不恰当的 `defer` 会导致关闭文件延时。

```go
func main() {
  for i := 0; i < 100; i ++ {
    f, err := os.Open(fmt.Sprintf("%d.log", i))
    if err != nil {
      continue
    }
    defer f.Close()
    // something
  }
}
```

上面的 `defer` 导致所有的 `f` 都是在 `main` 函数退出时才调用，白白消耗了资源。所以应该直接调用 `Close` 函数，将文件操作封装到一个函数中，在该函数中调用 `Close` 函数。


## defer 函数的参数会立即求值

`defer` 声明时，会立即对参数进行求值，而不是等延迟调用时才求值。

```go
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
        defer func(val int) {
            fmt.Println(val) // 打印传递的值
        }(i) // 在这里立即捕获当前的 i 值
    }
    fmt.Println("Loop ended")
}
```

输出：

```bash
Loop ended
2
1
0
```

但是如果不是函数传参，而是直接使用变量，会导致延迟调用时使用的是最新的值。

```go
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
        defer func() {
            fmt.Println(i) // 直接引用 i
        }()
    }
    fmt.Println("Loop ended")
}
```

输出：

```bash
Loop ended
3
3
3
```

`defer` 遇到链式调用时，会先通过计算得到最后一个要执行的函数，然后保留这个函数的指针、参数（值复制）。所以下面的代码的输出是 `1 3 2`。

```go
type T struct{}

func (t T) f(n int) T {
    fmt.Print(n)
    return t
}

func main() {
    var t T
    defer t.f(1).f(2)
    fmt.Print(3)
}

// defer t.f(1).f(2)
// 类似于
// var t T
// tmpT := t.f(1)
// defer tmpT.f(2)
// fmt.Print(3)
```


