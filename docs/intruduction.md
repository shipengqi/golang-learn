# 介绍
## 命名
所有命名只能以字母或者`_`开头，可以包含字母，数字或者`_`。区分大小写。
关键字不能定义变量名，如`func`，`default`。

注意：函数内部定义的，只能在函数内部使用（函数级），在函数外部定义的（包级），可以在当前包的所有文件中是使用。
并且，是在函数外定义的名字，如果以大写字母开头，那么会被导出，也就是在包的外部也可以访问，所以定义名字时，要注意大小写。

## 声明
- `var`声明变量
- `const`声明常量
- `type`声明类型
- `func`声明函数

每个文件以`package`声明语句。比如`package main`。

## 变量
`var`声明变量，必须使用空格隔开：
```go
var 变量名字 类型 = 表达式
```
类型 或者 表达式可以省略其中的一个。也就是如果没有类型，可以通过表达式推导出类型，**没有表达式，将会根据类型初始化为对应的零值**。
对应关系：
- 数值类型：`0`
- 布尔类型：`false`
- 字符串：`""`
- 接口或引用类型（包括slice、指针、map、chan和函数）：`nil`

### 声明一组变量
```go
var 变量名字 变量名字 变量名字 ... 类型 = 表达式, 表达式, 表达式, ...
```
比如：
```go
// 声明一组`int`类型
var i, j, k int                 // int, int, int

// 声明一组不同类型
var b, f, s = true, 2.3, "four" // bool, float64, string
```

### 简短声明
`:=`只能在函数内使用，也就是只能声明局部变量，Go 会自动推导出类型：
```go
变量名字 := 表达式
```

### 指针
```go
x := 1
p := &x         // p, of type *int, points to x
fmt.Println(*p) // "1"
*p = 2          // equivalent to x = 2
fmt.Println(x)  // "2"
```

上面的代码，初始化一个变量`x`，`&`是取地址操作，`&x`就是取变量`x`的内存地址，那么`p`就是一个指针，
类型是`*int`，`p`这个指针保存了变量`x`的内存地址。接下来`*p`表示读取指针指向的变量的值，也就是变量`x`的值`1`。
`*p`也可以被赋值。

任何类型的指针的零值都是`nil`。当指针指向同一个变量或者`nil`时是相等的。

### new
`new`函数可以创建变量，返回变量的地址。`new`函数很少使用，直接使用字面量语法创建更灵活。
```go
new(T)
```
返回指针类型为`*T`。
如：
```go
p := new(int)   // p, *int 类型, 指向匿名的 int 变量
fmt.Println(*p) // "0"
*p = 2          // 设置 int 匿名变量的值为 2
fmt.Println(*p) // "2"
```
### 生命周期
对于在包一级声明的变量，它们的生命周期和程序的运行周期是一致的。
局部变量（包括函数的参数和返回值也是局部变量）的生命周期则是动态的：每次从创建一个新变量的声明语句开始，
直到该变量不再被引用为止，然后变量的存储空间可能被回收。

## 赋值
常见的赋值的方式：
```go
x = 1                       // 命名变量的赋值
*p = true                   // 通过指针间接赋值
person.name = "bob"         // 结构体字段赋值
count[x] = count[x] * scale // 数组、slice或map的元素赋值
count[x] *= scale           // 等价于 count[x] = count[x] * scale，但是省去了对变量表达式的重复计算
x, y = y, x                 // 交换值
f, err = os.Open("foo.txt") // 左边变量的数目必须和右边一致，函数一般会返回一个`error`类型
v, ok = m[key]              // map查找，返回布尔值类表示操作是否成功
v = m[key]                  // map查找，也可以返回一个值，失败时返回零值
```

不管是隐式还是显式地赋值，在赋值语句左边的变量和右边最终的求到的值必须有相同的数据类型。这就是**可赋值性**。

## 类型
创建了一个新的类型名称，一般出现在包一级，与变量一样，如果类型名字的首字母是大写，则在包外部也可以使用：
```go
type 类型名字 底层类型
```

如不同温度单位分别定义为不同的类型：
```go
type Celsius float64    // 摄氏温度
type Fahrenheit float64 // 华氏温度

const (
    AbsoluteZeroC Celsius = -273.15 // 绝对零度
    FreezingC     Celsius = 0       // 结冰点温度
    BoilingC      Celsius = 100     // 沸水温度
)

func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }
```

## 包
Go 语言的包与其他语言的`modules`或者`libraries`类似。Go语言有超过100个的标准包，可以使用`go list std | wc -l`查看包的数量。
更多 Go 语言开源包，可以在[这里](http://godoc.org)搜索。
我们知道 Go 语言编译速度很快，主要依赖下面三点：
1. 导入的包必须在文件的头部显式声明，这样的话编译器就没有必要读取和分析整个源文件来判断包的依赖关系。
2. 禁止包的循环依赖，每个包可以被独立编译，而且很可能是被并发编译。
3. 编译后包的目标文件不仅仅记录包本身的导出信息，同时还记录了包的依赖关系。
因此，在编译一个包的时候，编译器只需读取每个导入包的目标文件，而不需要遍历所有依赖的的文件。

### 导入包

每个包都有一个全局唯一的导入路径，一个导入路径如`import "packages/test"`，那么这个包的完整路径应该是`GOPTAH/src/packages/test`，
`test`应该是一个目录，包含了一个或多个`test`（包的名字和包的导入路径的最后一个字段相同）包的源文件。
```go
// 导入一个包
import "fmt"

// 导入多个包
import (
    "fmt"
    "os"
)
```

导入的包之间可以通过添加空行来分组；通常将来自不同组织的包独自分组。
```go
import (
    "fmt"
    "os"

    "golang.org/x/net/ipv4"
)
```


如果你的包会发布出去，那么导入路径最好是全球唯一的。为了避免冲突，所有非标准库包的导入路径建议以所在组织的互联网域名为前缀；而且这样也有利于包的检索。例如`import "github.com/go-sql-driver/mysql"`。

#### 导入包重命名
如果导入两个相同名字的包，如`math/rand`包和`crypto/rand`包，可以为一个包重命名来解决名字冲突：
```go
import (
    "crypto/rand"
    mrand "math/rand" // alternative name mrand avoids conflict
)
```
注意，重命名的包名只在当前源文件有效。

有些情况下也可以使用包重命名：
1. 包名很长。重命名一个简短的包名。
2. 与变量名冲突。

选择用简短名称重命名导入包时候最好统一，以避免包名混乱。

#### 匿名导入
比如`import _ "image/png"`，`_`是空白标识符，不能被访问。
匿名导入有什么用？我们知道如果导入一个包而不使用会导致编译错误`unused import`。当我们想要导入包，
仅仅只是想计算导入包的包级变量的初始化表达式和执行导入包的`init`初始化函数，就可以使用匿名导入。

### 包声明
包声明语句（包名）必须在每个源文件的开头。被其它包导入时默认的标识符。每个包都对应一个独立的名字空间，
如：`image`包和`unicode/utf16`包中都包含了`Decode`。要在外部引用该函数，必须显式使用`image.Decode`或`utf16.Decode`形式访问。

**包内以大写字母开头定义的名字（包括变量，类型，函数等等），会被导出，可以在包的外部访问。**

默认包名一般采用导入路径名的最后一段，比如`GOPTAH/src/packages/test`的`test`就是包名。三种情况例外：
1. `main`包，`go build`命令编译完之后生成一个可执行程序。
2. 以`_test`为后缀包名的测试外部扩展包都由`go test`命令独立编译。(以`_`或`.`开头的源文件会被构建工具忽略)
3. 如"gopkg.in/yaml.v2"。包的名字包含版本号后缀`.v2`，这种情况下包名是`yaml`。

### 包命名
包命名尽量有描述性且无歧义，简短，避免冲突。

### 初始化包
包的初始化首先是解决包级变量的依赖顺序，然后按照包级变量声明出现的顺序依次初始化：
```go
var a = b + c // a 第三个初始化, 为 3
var b = f()   // b 第二个初始化, 为 2, 通过调用 f (依赖c)
var c = 1     // c 第一个初始化, 为 1

func f() int { return c + 1 }
```
如果包中含有多个源文件，构建工具首先会将`.go`文件根据文件名排序，然后依次调用编译器编译。

每个包在解决依赖的前提下，以导入声明的顺序初始化，每个包只会被初始化一次。因此，如果一个 p 包导入了 q 包，
那么在 p 包初始化的时候可以认为 q 包必然已经初始化过了。初始化工作是自下而上进行的，`main`包最后被初始化。以这种方式，
可以确保在`main`函数执行之前，所有依赖的包都已经完成初始化工作了。

#### 使用`init`函数
使用`init`函数来简化初始化工作，`init`函数和普通函数类似，但是不能被调用或引用。
程序开始执行时按照它们声明的顺序自动调用。

## 作用域
声明语句的作用域是指源代码中可以有效使用这个名字的范围。我觉得 Go 语言的作用域和`Javascript`很相似。
句法块：由花括弧所包含的语句，比如函数体或者循环体花括弧包裹的内容。句法块内部声明的名字是无法被外部块访问的。
词法块：未显式地使用花括号包裹起来，比如对于全局的源代码，存在一个整体的词法块，称为全局词法块。对于每个包；
每个`for`、`if`和`switch`语句，也都有对应词法块；每个`switch`或`select`的分支也有独立的词法块；当然也包括显式书写的词法块（花括弧包含的语句）。

声明语句对应的词法域决定了作用域范围的大小。
- 内置的类型、函数和常量，比如`int`、`len`和`true`是全局作用域
- 在函数外部（也就是包级语法域）声明的名字可以在同一个包的任何源文件中访问
- 导入的包，如`import "packages/test"`，是对应源文件级的作用域，只能在当前的源文件中访问
- 在函数内部声明的名字，只能在函数内部访问

**一个程序可能包含多个同名的声明，只要它们在不同的词法域就可以。内层的词法域会屏蔽外部的声明。** 编译器遇到一个名字引用是，
会从最内层的词法域向全局查找（和`Javascript`相似）。

## 工具
Go 语言工具箱的命令：
```bash
$ go
Go is a tool for managing Go source code.

Usage:

        go command [arguments]

The commands are:

        build       compile packages and dependencies
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        bug         start a bug report
        fix         update packages to use new APIs
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         download and install packages and dependencies
        install     compile and install packages and dependencies
        list        list packages
        run         compile and run Go program
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help [command]" for more information about a command.

Additional help topics:

        c           calling between Go and C
        buildmode   build modes
        cache       build and test caching
        filetype    file types
        gopath      GOPATH environment variable
        environment environment variables
        importpath  import path syntax
        packages    package lists
        testflag    testing flags
        testfunc    testing functions

Use "go help [topic]" for more information about that topic.
```

### 工作区 GOPATH
### 下载包
### 构建包
### 包文档
### 内部包
### 查询包