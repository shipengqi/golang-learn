# 介绍
## 语法

- 可以直接运行（`go run *.go`，其实也是`run`命令进行编译运行），也可以编译后运行。
- 使用`var`声明变量，函数内部可以省略`var`。
- 函数可以返回多个值，函数是第一类型，可以作为参数或返回值。
- 控制语句只有三种`if`，`for`，`switch`。
- `os.Args`获取命令行参数。
## package
我们知道，我们在写 Go 语言的代码时，每个文件的头部都有一行`package`声明语句。比如`package main`。这个声明表示这个源文件属于哪个包（类似其他语言的`modules`或者`libraries`）。 Go 语言的代码就是通过这个`package`来组织。

## import
在`package`声明下面，我们需要导入一系列需要使用的包。比如`import "fmt"`。注意如果导入了不需要的包，或者缺少了必要的包，编译会失败。

## main
`main`是一个特殊的包，`main`包代表一个独立运行的程序，而不是一个`modules`或者`libraries`。`main`包里必须有`main`函数，这个是程序的入口函数，并且`mian`函数没有参数。比如：
```go
func main() {
	fmt.Println("Hello, 世界")
}
```

## hello world
```go
package main

import "fmt"

func main() {
	fmt.Println(x)
}
```
**函数声明使用`func`关键字。Go 不需要在语句或者声明的末尾添加分号。除非一行代码上有多条语句。**

## os.Args
程序的命令行参数可使用`os.Args`访问。`os.Args`是一个字符串的切片。我们打印看一下：
```go
package main

import (
  "fmt"
  "os"
)

func main() {
  for i := 1; i < len(os.Args); i ++ {
  	fmt.Println(os.Args[i])
  }
}
```
然后运行：
```bash
go run .\args1.go arg1 arg2 arg3
```



