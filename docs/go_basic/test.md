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
测试函数必须导入`testing`包，并以`Test`为函数名前缀，后缀名必须以大写字母开头：
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
测试函数必须导入`testing`包，并以`Benchmark`为函数名前缀，后缀名必须以大写字母开头：
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