---
title: Go 性能优化
weight: 5
---

# Go 性能优化

## JSON 优化

Go 官方的 `encoding/json` 是通过反射来实现的。性能相对有些慢。 可以使用第三方库来替代标准库：

- [json-iterator/go](https://github.com/json-iterator/go)，完全兼容标准库，性能有很大提升。
- [go-json](https://github.com/goccy/go-json)，完全兼容标准库，性能强于 `json-iterator/go`。
- [sonic](https://github.com/bytedance/sonic)，字节开发的的 JSON 序列化/反序列化库，速度快，但是对硬件有一些要求。

实际开发中可以根据编译标签来选择 `JSON` 库，参考 [component-base/json](https://github.com/shipengqi/component-base/tree/main/json)。

## 使用空结构体

在 Go 中空结构体 `struct{}` 不占据内存空间：

```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println(unsafe.Sizeof(struct{}{})) // 0
}
```

空结构体不占据内存空间，因此被广泛作为各种场景下的占位符使用，可以节省资源。

### 集合 Set

要实现一个 `Set`，通常会使用 `map` 来实现，比如 `map[string]bool`。 但是对于集合来说， 只需要 `map` 的键，而不需要值。将值设置为 `bool` 
类型，就会多占据 1 个字节。这个时候就可以使用空结构体 `map[string]struct{}`。

### channel 通知

有时候使用 `channel` 不需要发送任何的数据，只用来通知 goroutine 执行任务，或结束等。这个时候就可以使用空结构体。

## 内存对齐

### 为什么需要内存对齐？

CPU 访问内存时，并不是逐个字节访问，而是以**字长**（word size）为单位访问。比如：
- 64 位系统 1 个字长等于 8 个字节
- 32 位系统 1 个字长等于 4 个字节

因此 CPU 在读取内存时是一块一块进行读取的。这么设计的目的，是减少 CPU 访问内存的次数，加大 CPU 访问内存的吞吐量。比如同样读取 8 个字节的数据，一
次读取 4 个字节那么只需要读取 2 次。

进行内存对齐，就是为了减少 CPU 访问内存的次数。

![mem-align](https://raw.githubusercontent.com/shipengqi/illustrations/0a6f953f0387c30638ae8b0e03dda230194d10ab/go/mem-align.png)

上图中，假如 CPU 字长为 4 个字节。变量 a 和 b 的大小为 3 个字节，没有内存对齐之前，CPU 读取 b 时，需要访问两次内存：

1. 第一次读取 0-3 字节，移除不需要的 0-2 字节，拿到 b 的第一个字节，
2. 第二次读取 4-7 字节，读取到 b 的后面两个字节，并移除不需要的 6，7 字节。
3. 合并 4 个字节的数据
4. 放入寄存器

内存对齐后，a 和 b 都占据了 4 个字节空间，CPU 读取 b 就只需要访问一次内存，读取到 4-7 字节。

### 对齐系数

不同平台上的编译器都有自己默认的 “对齐系数”，常用的平台的系数如下：
- 64 位系统：8
- 32 位系统：4

`unsafe` 标准库提供了 `Alignof` 方法，可以返回一个类型的对齐系数。例如：

```go
func main() {
    fmt.Printf("bool align: %d\n", unsafe.Alignof(bool(true))) // bool align: 1
    fmt.Printf("int8 align: %d\n", unsafe.Alignof(int8(0))) // int8 align: 1
    fmt.Printf("int16 align: %d\n", unsafe.Alignof(int16(0))) // int16 align: 2
    fmt.Printf("int32 align: %d\n", unsafe.Alignof(int32(0))) // int32 align: 4
    fmt.Printf("int64 align: %d\n", unsafe.Alignof(int64(0))) // int64 align: 8
    fmt.Printf("byte align: %d\n", unsafe.Alignof(byte(0))) // byte align: 1
    fmt.Printf("string align: %d\n", unsafe.Alignof("EDDYCJY")) // string align: 8
    fmt.Printf("map align: %d\n", unsafe.Alignof(map[string]string{})) // map align: 8
}
```

### 对齐规则

1. 对于任意类型的变量 `x`，`unsafe.Alignof(x)` 至少为 1。
2. 对于 `struct` 结构体类型的变量 `x`，计算 `x` 每一个字段 `f` 的 `unsafe.Alignof(x.f)`，`unsafe.Alignof(x)` 等于其中的最大值。
3. 对于 `array` 数组类型的变量 `x`，`unsafe.Alignof(x)` 等于构成数组的元素类型的对齐倍数。

### Go 结构体内存对齐

`struct` 中的字段的顺序会对 `struct` 的大小产生影响吗？

```go
type Part1 struct {
    a int8
    c int32
    b int16
}

type Part2 struct {
	a int8
	c int32
	b int16
}

func main()  {
	part1 := Part1{}
	fmt.Printf("part1 size: %d, align: %d\n", unsafe.Sizeof(part1), unsafe.Alignof(part1))
	part2 := Part2{}
	fmt.Printf("part2 size: %d, align: %d\n", unsafe.Sizeof(part2), unsafe.Alignof(part2))
}
```

输出：

```
// Output:
// part1 size: 8, align: 4
// part2 size: 12, align: 4
```

`Part1` 只是对成员变量的字段顺序进行了调整，就减少了结构体占用大小。

![mem-align](https://raw.githubusercontent.com/shipengqi/illustrations/ca9f00e7c3f54f02935d6615da69123d09ee8c7c/go/struct-mem-align.png)

`part1`：

- a 从第 0 个位置开始占据 1 字节。
- b 对齐系数为 2，因此，必须空出 1 个字节，偏移量才是 2 的倍数，从第 2 个位置开始占据 2 字节。
- c 对齐系数为 4，此时，内存已经是对齐的，从第 4 个位置开始占据 4 字节即可。

`part2`：

- a 从第 0 个位置开始占据 1 字节。
- c 对齐系数为 4，因此，必须空出 3 个字节，偏移量才是 4 的倍数，从第 4 个位置开始占据 4 字节。
- b 对齐系数为 2，从第 8 个位置开始占据 2 字节。

### 空 `struct{}` 的对齐

空 `struct{}` 大小为 0，作为其他 struct 的字段时，一般不需要内存对齐。但是当 `struct{}` 作为结构体最后一个字段时，需要内存对齐。
因为如果有指针指向该字段, 返回的地址将在结构体之外，如果此指针一直存活不释放对应的内存，就会有内存泄露的问题（该内存不因结构体释放而释放）。

因此，当 `struct{}` 作为其他 `struct` 最后一个字段时，需要填充额外的内存保证安全。

```go
type Part1 struct {
	c int32
	a struct{}
}

type Part2 struct {
	a struct{}
	c int32
}

func main() {
	fmt.Println(unsafe.Sizeof(Part1{})) // 8
	fmt.Println(unsafe.Sizeof(Part2{})) // 4
}
```

可以看到 `Part1{`} 额外填充了 4 字节的空间。

## 逃逸分析

编译器决定内存分配位置的方式，就称之为**逃逸分析**(escape analysis)。逃逸分析由编译器完成，作用于编译阶段。

变量逃逸是指编译器将一个变量从栈上分配到对上的情况。

在 Go 中，栈是跟函数绑定的，函数结束时栈被回收。如果一个变量分配在栈中，则函数执行结束可自动将内存回收。如果分配在堆中，则函数执行结束可交给 GC（垃圾回收）处理。

变量逃逸常见的情况：

1. 指针逃逸：返回指针，当一个函数返回一个局部变量的指针时，编译器就不得不吧该变量分配到堆上，以便函数返回后还可以访问它。
2. 发送指针或带有指针的值到 channel 中，编译时，是没有办法知道哪个 goroutine 会在 channel 上接收数据。所以编译器没法知道变量什么时候才会被释放。该值就会被分配到堆上。
3. 在一个切片上存储指针或带指针的值。例如 `[]*string` 。这会导致切片的内容逃逸。尽管其后面的数组可能是在栈上分配的，但其引用的值一定是在堆上。
4. 切片的底层数组被重新分配了，因为 append 时可能会超出其容量。切片初始化的地方在编译时是可以知道的，它最开始会在栈上分配。如果切片背后的存储要基于运行时的数据进行扩充，就会在堆上分配。
5. 在 `interface` 类型上调用方法都是**动态调度**的，方法的实现只能在运行时才知道。比如 `io.Reader` 类型的变量 `r`，调用 `r.Read(b)` 会使 `r` 的值和切片 `b` 的底层数组都逃逸掉，在堆上分配。
6. 数据类型不确定，如调用 `fmt.Sprintf`，`json.Marshal` 等接受变量为 `...interface{}` 的函数，会导致传入的变量逃逸到堆上。
7. 闭包引用：如果一个局部变量被一个闭包函数引用，那么编译器也可能把它分配到堆上，确保闭包可以继续访问它。
   ```go
   func isaclosure() func() {
       v := 1
       return func() {
           println(v)
       }
   }
   ```
8. 栈空间不足

变量逃逸就意味着增加了堆中的对象个数，影响 GC 耗时，影响性能。所以编写代码时，避免返回指针，限制闭包的作用范围等来要尽量避免逃逸。

可以使用编译器的 `gcflags="-m"` 来查看变量逃逸的情况：

```go
package main

import "fmt"

type A struct {
	s string
}

// 在方法内返回局部变量的指针
func foo(s string) *A {
	a := new(A)
	a.s = s
	return a // a 会逃逸到堆上
}

func main() {
	a := foo("hello")
	b := a.s + " world"
	c := b + "!"
	fmt.Println(c) // c 数据类型不确定，所以 escapes to heap
}
```

运行 `go run -gcflags=-m ./main.go` 会得到下面类似的输出：

```
# command-line-arguments
./main.go:10:6: can inline foo
./main.go:17:10: inlining call to foo
./main.go:20:13: inlining call to fmt.Println
./main.go:10:10: leaking param: s
./main.go:11:10: new(A) escapes to heap
./main.go:17:10: new(A) does not escape
./main.go:18:11: a.s + " world" does not escape
./main.go:19:9: b + "!" escapes to heap
./main.go:20:13: c escapes to heap
./main.go:20:13: []interface {} literal does not escape
<autogenerated>:1: .this does not escape
<autogenerated>:1: .this does not escape
hello world!
```

### 传值还是传指针？

传值会拷贝整个对象，而传指针只会拷贝指针地址，指向的对象是同一个。传指针可以减少值的拷贝，但是会导致内存分配逃逸到堆中，增加垃圾回收(GC)的负担。在对
象频繁创建和删除的场景下，传递指针导致的 GC 开销可能会严重影响性能。

一般情况下，对于需要修改原对象值，或占用内存比较大的结构体，选择传指针。对于只读的占用内存较小的结构体，直接传值能够获得更好的性能。

## 死码消除

死码消除(dead code elimination, DCE)是一种编译器优化技术，用处是在编译阶段去掉对程序运行结果没有任何影响的代码。

死码消除可以减小程序体积，程序运行过程中避免执行无用的指令，缩短运行时间。

### 使用常量提升性能

有些场景下，使用常量不仅可以减少程序的体积，性能也会有很大的提升。


## 利用 sync.Pool 减少堆分配

## 字符串与字节转换优化，减少内存分配

## 垃圾回收优化

## 设置 GOMAXPROCS
