# 语法基础

- 可以直接运行（`go run *.go`，其实也是`run`命令进行编译运行），也可以编译后运行。
- 函数可以返回多个值，函数是第一类型，可以作为参数或返回值。
- 控制语句只有三种`if`，`for`，`switch`。

### 注释
使用`//`添加注释。一般我们会在包声明前添加注释，来对整个包挥着程序做整体的描述。

### package
我们知道，我们在写 Go 语言的代码时，每个文件的头部都有一行`package`声明语句。比如`package main`。这个声明表示这个源文件属于哪个包（类似其他语言的`modules`或者`libraries`）。 Go 语言的代码就是通过这个`package`来组织。

### 行分隔符
Go 中，一行代表一个语句结束，不需要以分号`;`结尾。多个语句写在同一行，则必须使用`;`（不推荐使用）。 

### import
在`package`声明下面，我们需要导入一系列需要使用的包。比如`import "fmt"`。注意如果导入了不需要的包，或者缺少了必要的包，编译会失败。
```go
// 导入一个包
import "fmt"

// 导入多个
import (
  "fmt"
  "os"
)
```

### main
`main`是一个特殊的包，`main`包代表一个独立运行的程序，而不是一个`modules`或者`libraries`。`main`包里必须有`main`函数，这个是程序的入口函数，并且`mian`函数没有参数。比如：
```go
func main() {
	fmt.Println("Hello, 世界")
}
```

### hello world
```go
package main

import "fmt"

func main() {
	fmt.Println(x)
}
```
**函数声明使用`func`关键字。Go 不需要在语句或者声明的末尾添加分号。除非一行代码上有多条语句。**

### os.Args
程序的命令行参数可使用`os.Args`访问。`os.Args`是一个字符串的切片。我们打印看一下：
```go
package main

import (
  "fmt"
  "os"
)

// ++ 和 -- 都只能放在变量名,如 i ++
func main() {
  for i := 1; i < len(os.Args); i ++ {
  	fmt.Println(os.Args[i])
  }
}
```
然后运行：
```bash
go run args1.go arg1 arg2 arg3
```

### 空标识符
`_`代表空标识符，Go 不允许有无用的变量，空标识符可用于任何语法需要变量名但程序逻辑不需要的时候，比如：
```go
var s, sep string
for _, arg := range os.Args[1:] {
	s += sep + arg
	sep = " "
}
```



