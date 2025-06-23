---
title: Go 性能优化
weight: 6
---

## 预分配容量

对于切片和 map，尽量预分配容量来避免触发扩容机制，扩容是一个比较耗时的操作。

## 字符串拼接

使用 `strings.Builder` 或 `bytes.Buffer` 操作字符串，参考 [字符串拼接](../../basic/01_basic_type/#%e5%ad%97%e7%ac%a6%e4%b8%b2%e6%8b%bc%e6%8e%a5)。

## JSON 优化

Go 的标准库 `encoding/json` 是通过反射来实现的。性能相对有些慢。 可以使用第三方库来替代标准库：

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

## 使用原子操作保护变量

在并发编程中，对于一个变量更新的保护，原子操作通常会更有效率。参考 [互斥锁与原子操作](../../concurrency/08_atomic/#互斥锁与原子操作)。

## 减少循环中的内存读写操作

```go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	input, _ := strconv.Atoi(os.Args[1]) // Get an input number from the command line
	u := int32(input)
	r := int32(rand.Uint32() % 10000)   // Use Uint32 for faster random number generation
	var a [10000]int32                  // Array of 10k elements initialized to 0
	for i := int32(0); i < 10000; i++ { // 10k outer loop iterations
		for j := int32(0); j < 100000; j++ { // 100k inner loop iterations, per outer loop iteration
			a[i] = a[i] + j%u // Simple sum
		}
		a[i] += r // Add a random value to each element in array
	}
	z := a[r]
	fmt.Println(z) // Print out a single element from the array
}

```

编译测试：

```bash
$go build -o code code.go
$time ./code 10
456953

real 0m3.766s
user 0m3.767s
sys 0m0.007s
```

修改代码，将数组元素累积到一个临时变量中，并在外层循环结束后写回数组，这样做可以减少内层循环中的内存读写操作，充分利用 CPU 缓存和寄存器，加速数据处理。

{{< callout type="info" >}}
数组从内存或缓存读，而一个临时变量很大可能是从寄存器读，那读取速度相差还是很大的。
{{< /callout >}}

```go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	input, e := strconv.Atoi(os.Args[1]) // Get an input number from the command line
	if e != nil {
		panic(e)
	}
	u := int32(input)
	r := int32(rand.Intn(10000))        // Get a random number 0 <= r < 10k
	var a [10000]int32                  // Array of 10k elements initialized to 0
	for i := int32(0); i < 10000; i++ { // 10k outer loop iterations
		temp := a[i]
		for j := int32(0); j < 100000; j++ { // 100k inner loop iterations, per outer loop iteration
			temp += j % u // Simple sum
		}
		temp += r // Add a random value to each element in array
		a[i] = temp
	}
	fmt.Println(a[r]) // Print out a single element from the array
}
```

编译测试：

```bash
$go build -o code code.go
$time ./code 10
459169

real 0m3.017s
user 0m3.017s
sys 0m0.007s
```

参考文章：[惊！Go 在十亿次循环和百万任务中表现不如 Java，究竟为何？](https://mp.weixin.qq.com/s/hTQiEmf3ztRS-77fBET91A)

## 内存对齐

### 为什么需要内存对齐？

CPU 访问内存时，并不是逐个字节访问，而是以**字长**（word size）为单位访问。比如：

- 64 位系统 1 个字长等于 8 个字节
- 32 位系统 1 个字长等于 4 个字节

因此 **CPU 在读取内存时是一块一块进行读取的**。这么设计的目的，是**减少 CPU 访问内存的次数，加大 CPU 访问内存的吞吐量**。比如同样读取 8 个字节的数据，一
次读取 4 个字节那么只需要读取 2 次。

进行**内存对齐，就是为了减少 CPU 访问内存的次数**。

![mem-align](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/mem-align.png)

上图中，假如 CPU 字长为 4 个字节。变量 a 和 b 的大小为 3 个字节，没有内存对齐之前，CPU 读取 b 时，需要访问两次内存：

1. 第一次读取 0-3 字节，移除不需要的 0-2 字节，拿到 b 的第一个字节，
2. 第二次读取 4-7 字节，读取到 b 的后面两个字节，并移除不需要的 6，7 字节。
3. 合并 4 个字节的数据
4. 放入寄存器

内存对齐后，a 和 b 都占据了 4 个字节空间，CPU 读取 b 就只需要访问一次内存，读取到 4-7 字节。

### 各类型的对齐系数

`unsafe` 标准库提供了 **`Alignof` 方法，可以返回一个类型的对齐系数**。例如：

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

### 对齐计算公式

对于任意类型 T，其对齐系数 align 为：

- 基本类型：参考上面的类型对齐系数。
- 数组类型：与其元素类型相同。
- 结构体类型：等于其字段中最大的对齐系数。
- 其他类型(如指针、接口等)：与平台相关。

### Go 结构体内存对齐

`struct` 中的字段的顺序会对 `struct` 的大小产生影响吗？

```go
type Part1 struct {
    a int8
	b int16
    c int32
}

type Part2 struct {
	a int8
	c int32
	b int16
}

func main()  {
	part1 := Part1{}
	fmt.Printf("part1 size: %d, align: %d\n", unsafe.Sizeof(part1), unsafe.Alignof(part1))
	// part1 size: 8, align: 4
	part2 := Part2{}
	fmt.Printf("part2 size: %d, align: %d\n", unsafe.Sizeof(part2), unsafe.Alignof(part2))
	// part2 size: 12, align: 4
}
```

`Part1` 只是对成员变量的字段顺序进行了调整，就减少了结构体占用大小。

![mem-align](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/struct-mem-align.png)

`part1`：

- a 从第 0 个位置开始占据 1 字节。
- b 对齐系数为 2，因此，必须空出 1 个字节，偏移量才是 2 的倍数，从第 2 个位置开始占据 2 字节。
- c 对齐系数为 4，此时，内存已经是对齐的，从第 4 个位置开始占据 4 字节即可。

`part2`：

- a 从第 0 个位置开始占据 1 字节。
- c 对齐系数为 4，因此，必须空出 3 个字节，偏移量才是 4 的倍数，从第 4 个位置开始占据 4 字节。
- b 对齐系数为 2，从第 8 个位置开始占据 2 字节。

### 空 `struct{}` 的对齐

**空 `struct{}` 大小为 0**，作为其他 `struct` 的字段时，**一般不需要内存对齐**。

但是**当 `struct{}` 作为结构体最后一个字段时，需要内存对齐**。

因为虽然空结构体本身不占用内存，且**如果存在指向该字段的指针，可能会返回超出结构体范围的地址。这个返回的地址可能指向另一个被分配的内存块**。如果此指针一直存活不释放对应的内存，就会有内存泄露的问题（该内存不因结构体释放而释放）。

因此，**当空 `struct{}` 作为一个结构体的最后一个字段时，需要填充额外的内存保证安全**。

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

可以看到 `Part1{}` 额外填充了 4 字节的空间。

## 逃逸分析

编译器决定内存分配位置的方式，就称之为**逃逸分析**(escape analysis)。逃逸分析由编译器完成，作用于编译阶段。

变量逃逸是指编译器将一个变量从栈上分配到堆上的情况。

在 Go 中，栈是跟函数绑定的，函数结束时栈被回收。如果一个变量分配在栈中，则函数执行结束可自动将内存回收。如果分配在堆中，则函数执行结束可交给 GC（垃圾回收）处理。

变量逃逸常见的情况：

1. 指针逃逸：**返回指针，当一个函数返回一个局部变量的指针**时，编译器就不得不吧该变量分配到堆上，以便函数返回后还可以访问它。
2. **发送指针或带有指针的值到 channel 中**，编译时，是没有办法知道哪个 goroutine 会在 channel 上接收数据。所以编译器没法知道变量什么时候才会被释放。该值就会被分配到堆上。
3. **在一个切片上存储指针或带指针的值**。例如 `[]*string` 。这会导致切片的内容逃逸。尽管其后面的数组可能是在栈上分配的，但其引用的值一定是在堆上。
4. **切片的底层数组被重新分配了**，因为 `append` 时可能会超出其容量。切片初始化的地方在编译时是可以知道的，它最开始会在栈上分配。如果切片背后的存储要基于运行时的数据进行扩充，就会在堆上分配。
5. 在 `interface` 类型上调用方法都是**动态调度**的，方法的实现只能在运行时才知道。比如 `io.Reader` 类型的变量 `r`，调用 `r.Read(b)` 会使 `r` 的值和切片 `b` 的底层数组都逃逸掉，在堆上分配。
6. **数据类型不确定**，如调用 `fmt.Sprintf`，`json.Marshal` 等接受变量为 `...interface{}` 的函数，会导致传入的变量逃逸到堆上。
7. **闭包引用**：如果一个局部变量被一个闭包函数引用，那么编译器也可能把它分配到堆上，确保闭包可以继续访问它。
   ```go
   func isaclosure() func() {
       v := 1
       return func() {
           println(v)
       }
   }
   ```
8. **栈空间不足**

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

总结，**类型不确定的，大小不确定的，被函数外引用的变量都会逃逸到堆上**。

### 传值还是传指针？

传值会拷贝整个对象，而传指针只会拷贝指针地址，指向的对象是同一个。传指针可以减少值的拷贝，但是会导致内存分配逃逸到堆中，增加垃圾回收(GC)的负担。在**对象频繁创建和删除的场景下，传递指针导致的 GC 开销可能会严重影响性能**。

一般情况下，对于**需要修改原对象值，或占用内存比较大的结构体，选择传指针**。对于**只读的占用内存较小的结构体，直接传值**能够获得更好的性能。

## 死码消除

死码消除(dead code elimination, DCE)是一种编译器优化技术，用处是**在编译阶段去掉对程序运行结果没有任何影响的代码**。

**死码消除可以减小程序体积，程序运行过程中避免执行无用的指令，缩短运行时间**。

### 使用常量提升性能

有些场景下，使用常量不仅可以减少程序的体积，性能也会有很大的提升。

`usevar.go`：

```go
func Max(num1, num2 int) int {
   if num1 > num2 {
      return num1
   }
   return num2
}

var a, b = 10, 20

func main() {
   if Max(a, b) == a {
      fmt.Println(a)
   }
}
```

`useconst.go`：

```go
func Max(num1, num2 int) int {
	if num1 > num2 {
		return num1
	}
	return num2
}

const a, b = 10, 20

func main() {
	if Max(a, b) == a {
		fmt.Println(a)
	}
}
```

上面两个文件编译后的文件大小：

```
$ ls -lh
-rwxr-xr-x 1 pshi2 1049089 1.9M Oct 24 13:45 usevar.exe
-rwxr-xr-x 1 pshi2 1049089 1.5M Oct 24 13:44 useconst.exe
```

只是使用了常量代替变量，两个文件的大小就相差 0.3 M，为什么？

使用 `-gcflags=-m` 参数可以查看编译器做了哪些优化：

```
$ go build -gcflags=-m ./useconst.go
# command-line-arguments
./main.go:5:6: can inline Max
./main.go:15:8: inlining call to Max
./main.go:16:14: inlining call to fmt.Println
./main.go:16:14: ... argument does not escape
./main.go:16:15: a escapes to heap
```

`Max` 函数被内联了，内联后的代码是这样的：

```go
func main() {
	var result int
	if a > b {
		result = a
	} else {
		result = b
    }
	if result == a {
		fmt.Println(a)
	}
}
```

由于 a 和 b 均为常量，在编译阶段会直接计算：

```go
func main() {
	var result int
	if 10 > 20 {
		result = 10
	} else {
		result = 20
    }
	if result == 10 {
		fmt.Println(a)
	}
}
```

`10 > 20` 永远为假，那么分支消除，`result` 永远等于 20：

```go
func main() {
	if 20 == 10 {
		fmt.Println(a)
	}
}
```

`20 == 10` 也永远为假，再次消除分支：

```go
func main() {
	fmt.Println(10)
}
```

但是对于变量 a 和 b，编译器并不知道运行过程中 a、b 会不会发生改变，因此不能够进行死码消除，这部分代码被编译到最终的二进制程序中。因此编译后的二进制程序体积大了 0.3 M。

因此，**在声明全局变量时，如果能够确定为常量，尽量使用 `const` 而非 `var`**。这样很多运算在编译器即可执行。死码消除后，既减小了二进制的体积，又可以提高运行时的效率。

### 可推断的局部变量

Go 编译器只对函数的局部变量做了优化，当可以推断出函数的局部变量的值时，死码消除仍然会生效，例如：

```go
func main() {
	var a, b = 10, 20
	if max(a, b) == a {
		fmt.Println(a)
	}
}
```

上面的代码与 `useconst.go` 的编译结果是一样的，因为编译器可以推断出 a、b 变量的值。

如果增加了并发操作：

```go
func main() {
	var a, b = 10, 20
	go func() {
		b, a = a, b
	}()
	if max(a, b) == a {
		fmt.Println(a)
	}
}
```

上面的代码，a、b 的值不能有效推断，死码消除失效。

包级别的变量推断难度是非常大的。函数内部的局部变量的修改只会发生在该函数中。但是如果是包级别的变量，对该变量的修改可能出现在：

- 包初始化函数 `init()` 中，`init()` 函数可能有多个，且可能位于不同的 `.go` 源文件。
- 包内的其他函数。
- 如果是 public 变量（首字母大写），其他包引用时可修改。

因此，Go 编译器只对局部变量作了优化。

## 利用 sync.Pool 减少堆分配

[sync.Pool 使用](/golang-learn/docs/concurrency/06_pool/)。

## 控制 goroutine 的并发数量

基于 GPM 的 Go 调度器，可以大规模的创建 goroutine 来执行任务，可能 1k，1w 个 goroutine 没有问题，但是当 goroutine 非常大时，比如 10w，100w 甚至更多
就会出现问题。

1. 即使每个 goroutine 只分配 2KB 的内存，但是数量太多会导致内存占用暴涨，对 GC 造成极大的压力，GC 是有 STW 机制的，运行时会挂起用户程序直到垃圾回收完。虽然 Go 1.8 去掉了 STW 以及改成了并行 GC，性能上有了不
   小的提升但是，如果太过于频繁地进行 GC，依然会有性能瓶颈。
2. runtime 和 GC 也都是 goroutine，如果 goroutine 规模太大，内存吃紧，Go 调度器就会阻塞 goroutine，进而导致内存溢出，甚至 crash。

### 利用 channel 的缓存区控制并发数量

```go
func main() {
	var wg sync.WaitGroup
	// 创建缓冲区大小为 3 的 channel
	ch := make(chan struct{}, 3)
	for i := 0; i < 10; i++ {
		// 如果缓存区满了，则会阻塞在这里
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Println(i)
			time.Sleep(time.Second)
			// 释放缓冲区
			<-ch
		}(i)
	}
	wg.Wait()
}
```

### 使用第三方 goroutine pool

常用的第三方 goroutine pool：

- [ants](https://github.com/panjf2000/ants)
- [conc](https://github.com/sourcegraph/conc)

## 零拷贝优化

### 优化字符串与 []byte 转换，减少内存分配

在开发中，字符串与 `[]byte` 相互转换是经常用到的。直接通过类型转换 `string(bytes)` 或者 `[]byte(str)` 会带来数据的复制，性能不佳。

在 Go 1.20 之前的版本可以采用下面的方式来优化：

```go
// B2S convert []byte to string.
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B convert string to []byte.
func S2B(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))

	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len

	return b
}
```


Go 1.20 提供了新的方式：

```go
// B2S convert []byte to string.
func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// S2B convert string to []byte.
func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
```

## 设置 GOMAXPROCS

`GOMAXPROCS` 是 Go 提供的一个非常重要的环境变量。设置它的值可以调整调度器 Processor 的数量，每个 Processor 都会绑定一个系统线程。所以 Processor 的数量，会影响 Go 的并发性能。

Go 1.5 版本以后，`GOMAXPROCS` 的默认值是机器的 CPU 核数（`runtime.NumCPU()` 的返回值）。

但是 `runtime.NumCPU()` 在容器中是无法获取正确的 CPU 核数的，因为容器是使用 `cgroup` 技术对 CPU 资源进行隔离限制的，但 `runtime.NumCPU()` 获取的却是**宿主机的 CPU 核数**。
例如一个 Kubernetes 集群中 Node 核数是 36，然后创建一个 Pod，并且限制 Pod 的 CPU 核数是 1。Pod 中的进程在设置 `GOMAXPROCS` 后，线程数量是 36。导致线程过多，线程频繁切换，增加上线文切换的负担。

Uber 提供了一个库 `go.uber.org/automaxprocs` 可以解决这个问题：

```go
package main

import (
	_ "go.uber.org/automaxprocs"
)

func main() {
	// ...
}
```

{{< callout type="info" >}}
`go.uber.org/automaxprocs` **只会在程序启动时执行一次**。如果容器在运行过程中，CPU 的 limit 被调整了（比如 k8s 调整了 Pod 的 CPU limit），`go.uber.org/automaxprocs` 是感知不到的。

**Go 1.25 的 runtime，可以周期性的检查 cgroup 的限制。如果限制改变了，会自动调整 `GOMAXPROCS` 的值**。
{{< /callout >}}