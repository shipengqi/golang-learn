# 测试
`go test`命令测试代码，包目录内，所有以`_test.go`为后缀名的源文件在执行`go build`时不会被构建成包的一部分，它们是`go test`测试的一部分。

在`*_test.go`文件中，有三种类型的函数：
- 测试函数，测试程序的一些逻辑行为是否正确。`go test`命令会调用这些测试函数并报告测试结果是`PASS`或`FAIL`。
- 基准测试函数，衡量一些函数的性能。`go test`命令会多次运行基准函数以计算一个平均的执行时间。
- 示例函数，提供一个由编译器保证正确性的示例文档。

`go test`会生成一个临时`main`包调用测试函数。
**参数**
- `-v`，打印每个测试函数的名字和运行时间。
- `-run`，指定一个正则表达式，只有匹配到的测试函数名才会被`go test`运行，如`go test -v -run="French|Canal"`。

## 测试函数
**测试函数必须导入`testing`包，并以`Test`为函数名前缀，后缀名必须以大写字母开头，并且参数列表中只应有一个`*testing.T`类型的参数声明**：
```go
func TestName(t *testing.T) {
  ...
}
```
`t`参数用于报告测试失败和附加的日志信息。`t.Error`和`t.Errorf`打印错误日志。`t.Fatal`或`t.Fatalf`停止当前测试函数
`go test`命令如果没有参数指定包那么将默认采用当前目录对应的包。示例查看`src/unittest/`。

将所有测试数据合并到一个测试中的表格中:
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
`go test`命令中集成了测试覆盖率工具。
运行`go tool cover`：
```bash
$ go tool cover
Usage of 'go tool cover':
Given a coverage profile produced by 'go test':
    go test -coverprofile=c.out

Open a web browser displaying annotated source code:
    go tool cover -html=c.out
```

添加`-coverprofile`参数，统计覆盖率数据，并将统计日志数据写入指定文件，如`go test -run=Coverage -coverprofile=c.out`。
`-covermode=count`参数将在每个代码块插入一个计数器而不是布尔标志量。在统计结果中记录了每个块的执行次数，
这可以用于衡量哪些是被频繁执行的热点代码。

## 基准测试
**测试函数必须导入`testing`包，并以`Benchmark`为函数名前缀，后缀名必须以大写字母开头，并且唯一参数的类型必须是`*testing.B`类型的**：
```go
func BenchmarkName(b *testing.B) {
  ...
}
```
`*testing.B`参数除了提供和`*testing.T`类似的方法，还有额外一些和性能测量相关的方法。

### 运行基准测试
运行基准测试需要使用`-bench`参数，指定要运行的基准测试函数。该参数是一个正则表达式，用于匹配要执行的基准测试函数的名字，默认值是空的。

`.`会匹配所有基准测试函数。

### 剖析
基准测试对于衡量特定操作的性能是有帮助的，Go语言支持多种类型的剖析性能分析：
1. CPU剖析数据标识了最耗CPU时间的函数。
2. 堆剖析则标识了最耗内存的语句。
3. 阻塞剖析则记录阻塞goroutine最久的操作，例如系统调用、管道发送和接收，还有获取锁等。

```bash
$ go test -cpuprofile=cpu.out
$ go test -blockprofile=block.out
$ go test -memprofile=mem.out
```

#### go tool pprof
`go tool pprof`命令可以用来分析上面的命令生成的数据。

## 示例函数
并以`Benchmark`为函数名前缀，示例函数没有函数参数和返回值：
```go
func ExampleName() {
  ...
}
```

三个用处:
1. 作为文档，如`ExampleIsPalindrome`示例函数将是`IsPalindrome`函数文档的一部分。
2. `go test`会运行示例函数测试。
3. 提供 Go Playground，可以在浏览器中在线编辑和运行每个示例函数。

## go test命令执行的主要测试流程
go test命令在开始运行时，会先做一些准备工作，比如，确定内部需要用到的命令，检查我们指定的代码包或源码文件的有效性，以及判断我们给予的标记是否合法，等等。

在准备工作顺利完成之后，go test命令就会针对每个被测代码包，依次地进行构建、执行包中符合要求的测试函数，清理临时文件，打印测试结果。这就
是通常情况下的主要测试流程。

对于每个被测代码包，go test命令会**串行地执行测试流程中的每个步骤**。

但是，为了加快测试速度，它通常会并发地对多个被测代码包进行功能测试，只不过，在最后打印测试结果的时候，它会依照我们给定的顺序逐个进行，这会让我们感觉到它是
在完全串行地执行测试流程。

由于**并发的测试会让性能测试的结果存在偏差，所以性能测试一般都是串行进行的**。

## 功能测试的测试结果
```bash
$ go test puzzlers/article20/q2
ok   puzzlers/article20/q2 (cached)
```
``(cached)`表明，由于测试代码与被测代码都没有任何变动，所以`go test`命令直接把之前缓存测试成功的结果打印出来了。

go 命令通常会缓存程序构建的结果，以便在将来的构建中重用。我们可以通过运行`go env GOCACHE`命令来查看缓存目录的路径。

运行`go clean -testcache`将会删除所有的测试结果缓存。不过，这样做肯定不会删除任何构建结果缓存。

设置环境变量`GODEBUG`的值也可以稍稍地改变 go 命令的缓存行为。比如，设置值为`gocacheverify=1`将会导致 go 命令绕过任何的缓存数据，而真正地执行
操作并重新生成所有结果，然后再去检查新的结果与现有的缓存数据是否一致。

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

**第一个标记及其值为`-bench=.`，只有有了这个标记，命令才会进行性能测试**。该标记的值`.`表明需要执行任意名称的性能测试函数。

第二个标记及其值是`-run=^$`，这个标记用于表明需要执行哪些功能测试函数，这同样也是以函数名称为依据的。该标记的值`^$`意味着：只执行名称为空的功能测试函数，
换句话说，不执行任何功能测试函数。

这两个标记的值都是正则表达式。实际上，它们只能以正则表达式为值。此外，如果运行`go test`命令的时候不加`-run`标记，那么就会使它执行被测代码包中的所有功能测试函数。

测试结果，重点在倒数第三行的内容。`BenchmarkGetPrimes-8`被称为单个性能测试的名称，它表示命令执行了性能测试函数`BenchmarkGetPrimes`，并且当时所用的最大 P 数量为8。

最大 P 数量相当于可以同时运行 goroutine 的逻辑 CPU 的最大个数。这里的逻辑 CPU，也可以被称为 CPU 核心，但它并不等同于计算机中真正的 CPU 核心，只是 Go 语言运行时系统
内部的一个概念，代表着它同时运行 goroutine 的能力。

可以通过调用`runtime.GOMAXPROCS`函数改变最大 P 数量，也可以在运行`go test`命令时，加入标记`-cpu`来设置一个最大 P 数量的列表，以供命令在多次测试时使用。

测试名称右边的是执行次数。**它指的是被测函数的执行次数，而不是性能测试函数的执行次数**。

## `-parallel`标记
该标记的作用是：设置同一个被测代码包中的功能测试函数的最大并发执行数。
该标记的默认值是测试运行时的最大 P 数量（这可以通过调用表达式`runtime.GOMAXPROCS(0)`获得）。

对于功能测试，为了加快测试速度，命令通常会并发地测试多个被测代码包。但是，在默认情况下，**对于同一个被测代码包中的多个功能测试函数，命令会串行地执行它们**。
除非我们在一些功能测试函数中显式地调用`t.Parallel`方法。

这个时候，这些包含了`t.Parallel`方法调用的功能测试函数就会被`go test`命令并发地执行，而并发执行的最大数量正是由`-parallel`标记值决定的。要注意，同一个功能测试函数
的多次执行之间一定是串行的。

## 性能测试函数中的计时器
`testing.B`类型有这么几个指针方法：`StartTimer`、`StopTimer`和`ResetTimer`。这些方法都是用于操作当前的性能测试函数专属的计时器的。

这些字段用于记录：当前测试函数在当次执行过程中耗费的时间、分配的堆内存的字节数以及分配次数。