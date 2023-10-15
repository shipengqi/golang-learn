---
title: Go 测试
weight: 3
---

# Go 测试

`go test` 命令测试代码，包目录内，所有以 `_test.go` 为后缀名的源文件在执行 `go build` 时不会被构建成包的一部分，
它们是 `go test` 测试的一部分。

在 `*_test.go` 文件中，有三种类型的函数：

- 测试函数，测试程序的一些逻辑行为是否正确。`go test` 命令会调用这些测试函数并报告测试结果是 `PASS` 或 `FAIL`。
- 基准测试函数，衡量一些函数的性能。`go test` 命令会多次运行基准函数以计算一个平均的执行时间。
- 示例函数，提供一个由编译器保证正确性的示例文档。

`go test` 会生成一个临时 `main` 包调用测试函数。
**参数**

- `-v`，打印每个测试函数的名字和运行时间。
- `-run`，指定一个正则表达式，只有匹配到的测试函数名才会被 `go test` 运行，如 `go test -v -run="French|Canal"`。
- `-cover`，测试覆盖率。
- `-bench`，运行基准测试。例如 `go test -bench=.`（如果在 Windows Powershell 环境下使用 `go test -bench="."`）
- `-c`，生成用于运行测试的可执行文件，但不执行它。这个可执行文件会被命名为 `pkg.test`，其中的 `pkg` 即为被测试代码包的
导入路径的最后一个元素的名称。
- `-i`，安装/重新安装运行测试所需的依赖包，但不编译和运行测试代码。
- `-o`，指定用于运行测试的可执行文件的名称。追加该标记不会影响测试代码的运行，除非同时追加了标记 `-c` 或 `-i`。

## 测试函数

**测试函数必须导入 `testing` 包，并以 `Test` 为函数名前缀，后缀名必须以大写字母开头，并且参数列表中只应有一个 `*testing.T`
类型的参数声明**：

```go
func TestName(t *testing.T) {
  ...
}
```

`t` 参数用于报告测试失败和附加的日志信息。`t.Error` 和 `t.Errorf` 打印错误日志。`t.Fatal` 或 `t.Fatalf` 停止当前测试函数
`go test` 命令如果没有参数指定包那么将默认采用当前目录对应的包。

表格驱动测试在我们要创建一系列相同测试方式的测试用例时很有用。例如:

```go
func TestIsPalindrome(t *testing.T) {
    var tests = []struct {
        input string
        want  bool
    }{
        {"", true},
        {"a", true},
        {"aa", true},
        {"ab", false},
        {"kayak", true},
        {"detartrated", true},
        {"A man, a plan, a canal: Panama", true},
        {"Evil I did dwell; lewd did I live.", true},
        {"Able was I ere I saw Elba", true},
        {"été", true},
        {"Et se resservir, ivresse reste.", true},
        {"palindrome", false}, // non-palindrome
        {"desserts", false},   // semi-palindrome
    }
    for _, test := range tests {
        if got := IsPalindrome(test.input); got != test.want {
            t.Errorf("IsPalindrome(%q) = %v", test.input, got)
        }
    }
}
```

## 覆盖率

`go test` 命令中集成了测试覆盖率工具。
运行 `go tool cover`：

```bash
$ go tool cover
Usage of 'go tool cover':
Given a coverage profile produced by 'go test':
    go test -coverprofile=c.out

Open a web browser displaying annotated source code:
    go tool cover -html=c.out
```

添加 `-coverprofile` 参数，统计覆盖率数据，并将统计日志数据写入指定文件，如 `go test -run=Coverage -coverprofile=c.out`。
`-covermode=count` 参数将在每个代码块插入一个计数器而不是布尔标志量。在统计结果中记录了每个块的执行次数，
这可以用于衡量哪些是被频繁执行的热点代码。

## 基准测试

**测试函数必须导入 `testing` 包，并以 `Benchmark` 为函数名前缀，后缀名必须以大写字母开头，并且唯一参数的类型必须
是 `*testing.B` 类型的**：

```go
func BenchmarkName(b *testing.B) {
  ...
}
```

`*testing.B` 参数除了提供和 `*testing.T` 类似的方法，还有额外一些和性能测量相关的方法。

### 运行基准测试

运行基准测试需要使用 `-bench` 参数，指定要运行的基准测试函数。该参数是一个正则表达式，用于匹配要执行的基准测试函数的名字，
默认值是空的。

`.` 会匹配所有基准测试函数。

### 剖析

基准测试对于衡量特定操作的性能是有帮助的，Go 语言支持多种类型的剖析性能分析：

1. CPU 剖析数据标识了最耗 CPU 时间的函数。
2. 堆剖析则标识了最耗内存的语句。
3. 阻塞剖析则记录阻塞 goroutine 最久的操作，例如系统调用、管道发送和接收，还有获取锁等。

```bash
go test -cpuprofile=cpu.out
go test -blockprofile=block.out
go test -memprofile=mem.out
```

#### go tool pprof

`go tool pprof` 命令可以用来分析上面的命令生成的数据。

## 示例函数

并以 `Benchmark` 为函数名前缀，示例函数没有函数参数和返回值：

```go
func ExampleName() {
  ...
}
```

三个用处:

1. 作为文档，如 `ExampleIsPalindrome` 示例函数将是 `IsPalindrome` 函数文档的一部分。
2. `go test` 会运行示例函数测试。
3. 提供 Go Playground，可以在浏览器中在线编辑和运行每个示例函数。

## go test 命令执行的主要测试流程

`go test` 命令在开始运行时，会先做一些准备工作，比如，确定内部需要用到的命令，检查我们指定的代码包或源码文件的有效性，
以及判断我们给予的标记是否合法，等等。

在准备工作顺利完成之后，go test 命令就会针对每个被测代码包，依次地进行构建、执行包中符合要求的测试函数，清理临时文件，
打印测试结果。这就是通常情况下的主要测试流程。

对于每个被测代码包，`go test` 命令会**串行地执行测试流程中的每个步骤**。

但是，为了加快测试速度，它通常会并发地对多个被测代码包进行功能测试，只不过，在最后打印测试结果的时候，它会依照我们给定的
顺序逐个进行，这会让我们感觉到它是在完全串行地执行测试流程。

由于**并发的测试会让性能测试的结果存在偏差，所以性能测试一般都是串行进行的**。

## 功能测试的测试结果

```bash
$ go test puzzlers/article20/q2
ok   puzzlers/article20/q2 (cached)
```

`(cached)` 表明，由于测试代码与被测代码都没有任何变动，所以 `go test` 命令直接把之前缓存测试成功的结果打印出来了。

go 命令通常会缓存程序构建的结果，以便在将来的构建中重用。我们可以通过运行 `go env GOCACHE` 命令来查看缓存目录的路径。

运行 `go clean -testcache` 将会删除所有的测试结果缓存。不过，这样做肯定不会删除任何构建结果缓存。

设置环境变量 `GODEBUG` 的值也可以稍稍地改变 go 命令的缓存行为。比如，设置值为 `gocacheverify=1` 将会导致 go 命令绕
过任何的缓存数据，而真正地执行操作并重新生成所有结果，然后再去检查新的结果与现有的缓存数据是否一致。

## 性能测试的测试结果

```bash
$ go test -bench=. -run=^$ puzzlers/article20/q3
goos: darwin
goarch: amd64
pkg: puzzlers/article20/q3
BenchmarkGetPrimes-8      500000       2314 ns/op
PASS
ok   puzzlers/article20/q3 1.192s
```

**第一个标记及其值为 `-bench=.`，只有有了这个标记，命令才会进行性能测试**。该标记的值 `.` 表明需要执行任意名称的性能测试函数。

第二个标记及其值是 `-run=^$`，这个标记用于表明需要执行哪些功能测试函数，这同样也是以函数名称为依据的。该标记的值 `^$` 意味着：
只执行名称为空的功能测试函数，换句话说，不执行任何功能测试函数。

这两个标记的值都是正则表达式。实际上，它们只能以正则表达式为值。此外，如果运行 `go test` 命令的时候不加 `-run` 标记，
那么就会使它执行被测代码包中的所有功能测试函数。

测试结果，重点在倒数第三行的内容。`BenchmarkGetPrimes-8` 被称为单个性能测试的名称，它表示命令执行了性能测试
函数 `BenchmarkGetPrimes`，并且当时所用的最大 P 数量为 8。

最大 P 数量相当于可以同时运行 goroutine 的逻辑 CPU 的最大个数。这里的逻辑 CPU，也可以被称为 CPU 核心，但它并不等同
于计算机中真正的 CPU 核心，只是 Go 语言运行时系统内部的一个概念，代表着它同时运行 goroutine 的能力。

可以通过调用 `runtime.GOMAXPROCS` 函数改变最大 P 数量，也可以在运行 `go test` 命令时，加入标记 `-cpu` 来设置一个最大 P 数量
的列表，以供命令在多次测试时使用。

测试名称右边的是执行次数。**它指的是被测函数的执行次数，而不是性能测试函数的执行次数**。

## `-parallel` 标记

该标记的作用是：设置同一个被测代码包中的功能测试函数的最大并发执行数。
该标记的默认值是测试运行时的最大 P 数量（这可以通过调用表达 式`runtime.GOMAXPROCS(0)` 获得）。

对于功能测试，为了加快测试速度，命令通常会并发地测试多个被测代码包。但是，在默认情况下，**对于同一个被测代码包中的多个功
能测试函数，命令会串行地执行它们**。除非我们在一些功能测试函数中显式地调用 `t.Parallel`方 法。

这个时候，这些包含了 `t.Parallel` 方法调用的功能测试函数就会被 `go test` 命令并发地执行，而并发执行的最大数量正是
由 `-parallel` 标记值决定的。要注意，同一个功能测试函数的多次执行之间一定是串行的。

## 性能测试函数中的计时器

`testing.B` 类型有这么几个指针方法：`StartTimer`、`StopTimer` 和 `ResetTimer`。这些方法都是用于操作当前的性能测试函数
专属的计时器的。

这些字段用于记录：当前测试函数在当次执行过程中耗费的时间、分配的堆内存的字节数以及分配次数。

## 性能分析

Go 语言为程序开发者们提供了丰富的性能分析 API，和非常好用的标准工具。这些 API 主要存在于：

- `runtime/pprof`；
- `net/http/pprof`；
- `runtime/trace`；

至于标准工具，主要有 `go tool pprof` 和 `go tool trace` 这两个。它们可以解析概要文件中的信息，并以人类易读的方式把这些
信息展示出来。

在 Go 语言中，用于分析程序性能的概要文件有三种，分别是：**CPU 概要文件（CPU Profile）、内存概要文件（Mem Profile）和阻塞概
要文件（Block Profile）**。

- CPU 概要文件，其中的每一段独立的概要信息都记录着，在进行某一次采样的那个时刻，CPU 上正在执行的 Go 代码。
- 内存概要文件，其中的每一段概要信息都记载着，在某个采样时刻，正在执行的 Go 代码以及堆内存的使用情况，这里包含已分配和已释放的
字节数量和对象数量。
- 阻塞概要文件，其中的每一段概要信息，都代表着 Go 程序中的一个 goroutine 阻塞事件。

### 程序对 CPU 概要信息进行采样

这需要用到 `runtime/pprof` 包中的 API。想让程序开始对 CPU 概要信息进行采样的时候，需要调用这个代码包中
的 `StartCPUProfile` 函数，而在停止采样的时候则需要调用该包中的`StopCPUProfile`函数。

### 设定内存概要信息的采样频率

针对内存概要信息的采样会按照一定比例收集 Go 程序在运行期间的堆内存使用情况。设定内存概要信息采样频率的方法很简单，
只要为 `runtime.MemProfileRate` 变量赋值即可。

这个变量的含义是，平均每分配多少个字节，就对堆内存的使用情况进行一次采样。如果把该变量的值设为0，那么，Go 语言运行时系统就
会完全停止对内存概要信息的采样。该变量的缺省值是 512 KB，也就是 512 千字节。

**如果你要设定这个采样频率，那么越早设定越好，并且只应该设定一次，否则就可能会对 Go 语言运行时系统的采样工作，造成不良影响**。
比如，只在 `main` 函数的开始处设定一次。

当我们想获取内存概要信息的时候，还需要调用 `runtime/pprof` 包中的 `WriteHeapProfile` 函数。该函数会把收集好的内存概要信息，
写到我们指定的写入器中。

注意，我们通过 **`WriteHeapProfile` 函数得到的内存概要信息并不是实时的，它是一个快照，是在最近一次的内存垃圾收集工作完成时产
生的**。如果你想要实时的信息，那么可以调用 `runtime.ReadMemStats` 函数。不过要特别注意，该函数会引起 Go 语言调度器的短暂停顿。

### 获取到阻塞概要信息

调用 `runtime` 包中的 `SetBlockProfileRate` 函数，即可对阻塞概要信息的采样频率进行设定。该函数有一个名叫 `rate` 的参数，
它是 `int` 类型的。

这个参数的含义是，只要发现一个阻塞事件的持续时间达到了多少个纳秒，就可以对其进行采样。如果这个参数的值小于或等于0，那么就意
味着 Go 语言运行时系统将会完全停止对阻塞概要信息的采样。

当我们需要获取阻塞概要信息的时候，需要先调用 `runtime/pprof` 包中的 `Lookup` 函数并传入参数值 "block"，从而得到一
个 `*runtime/pprof.Profile` 类型的值（以下简称Profile值）。在这之后，我们还需要调用这个 `Profile` 值的 `WriteTo` 方法，
以驱使它把概要信息写进我们指定的写入器中。

`WriteTo` 方法有两个参数，一个参数就是我们刚刚提到的写入器，它是 `io.Writer` 类型的。而另一个参数则是代表了概要信息
详细程度的 `int` 类型参数 `debug`。

`debug` 参数主要的可选值有两个，即：0 和 1。当 `debug` 的值为 0 时，通过 `WriteTo` 方法写进写入器的概要信息仅会包含
 `go tool pprof` 工具所需的内存地址，这些内存地址会以十六进制的形式展现出来。

当该值为 1 时，相应的包名、函数名、源码文件路径、代码行号等信息就都会作为注释被加入进去。另外，`debug` 为 0 时的概要信息，
会经由 protocol buffers 转换为字节流。而在 `debug` 为 1 的时候，`WriteTo` 方法输出的这些概要信息就是我们可以读懂
的普通文本了。

除此之外，`debug` 的值也可以是 2。这时，被输出的概要信息也会是普通的文本，并且通常会包含更多的细节。至于这些细节都包含了哪些
内容，那就要看们调用 `runtime/pprof.Lookup` 函数的时候传入的是什么样的参数值了。
