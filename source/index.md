---
title: Golang 介绍
---

# Golang 介绍
Go 语言非常简单，只有 25 个关键字：

- `var` 和 `const` 声明变量和常量
- `package` 和 `import` 声明所属包名和导入包。
- `func` 用于定义函数和方法
- `return` 用于从函数返回
- `defer` 用于类似析构函数
- `go` 用于并发
- `select` 用于选择不同类型的通讯
- `interface` 用于定义接口
- `struct` 用于定义抽象数据类型
- `break`、`case`、`continue`、`for`、`fallthrough`、`else`、`if`、`switch`、`goto`、`default` 流程控制语句
- `chan` 用于 `channel` 通讯
- `type` 用于声明自定义类型
- `map` 用于声明 `map` 类型数据
- `range` 用于读取 `slice`、`map`、`channel` 数据

## 数据类型
Go 语言的四类数据类型
- 基础类型，数值、字符串和布尔型
- 复合类型，数组和结构体
- 引用类型，指针、切片、字典、函数、通道
- 接口类型

## 三种文件
- 命令源码文件，如果一个源码文件声明属于 `main` 包，并且包含一个无参数声明且无结果声明的 `main` 函数，那么它就是命令源码文件。
- 库源码文件，库源码文件是不能被直接运行的源码文件，它仅用于存放程序实体，这些程序实体可以被其他代码使用
- 测试源码文件


## package
在写 Go 语言的代码时，每个文件的头部都有一行 `package` 声明语句。比如 `package main`。这个声明表示这个源
文件属于哪个包（类似其他语言的 `modules` 或者 `libraries`）。 Go 语言的代码就是通过这个 `package` 来组织。

## 注释
使用 `//` 添加注释。一般我们会在包声明前添加注释，来对整个包挥着程序做整体的描述。

## 行分隔符
Go 中，一行代表一个语句结束，不需要以分号 `;` 结尾。多个语句写在同一行，则必须使用 `;`（不推荐使用）。 

## os.Args
程序的命令行参数可使用 `os.Args` 访问。`os.Args` 是一个字符串的切片。我们打印看一下：
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

## 空标识符
`_` 代表空标识符，**Go 不允许有无用的变量，空标识符可以作为忽略占位符**，比如：
```go
var s, sep string
for _, arg := range os.Args[1:] {
	s += sep + arg
	sep = " "
}
```


## 命名
所有命名只能以字母或者 `_` 开头，可以包含字母，数字或者 `_`。区分大小写。
关键字不能定义变量名，如 `func`，`default`。

**注意：函数内部定义的，只能在函数内部使用（函数级），在函数外部定义的（包级），可以在当前包的所有文件中是使用。
并且，是在函数外定义的名字，如果以大写字母开头，那么会被导出，也就是在包的外部也可以访问，所以定义名字时，要注意大小写。**

## 声明
- `var` 声明变量
- `const` 声明常量
- `type` 声明类型
- `func` 声明函数

每个文件以 `package` 声明语句。比如 `package main`。

## 变量
`var` 声明变量，必须使用空格隔开：
```go
var 变量名字 类型 = 表达式
```
**类型**或者**表达式**可以省略其中的一个。也就是如果没有类型，可以通过表达式推断出类型，**没有表达式，将会根据类型初始化为对应的零值**。
对应关系：
- 数值类型：`0`
- 布尔类型：`false`
- 字符串: `""`
- 接口或引用类型（包括 `slice`、指针、`map`、`chan` 和函数）：`nil`

### 声明一组变量
```go
var 变量名字, 变量名字, 变量名字 ... 类型 = 表达式, 表达式, 表达式, ...
```
比如：
```go
// 声明一组 `int` 类型
var i, j, k int                 // int, int, int

// 声明一组不同类型
var b, f, s = true, 2.3, "four" // bool, float64, string

var (
  i int
  pi float32
  prefix string
)
```

### 简短声明
** `:=` 只能在函数内使用，不能提供数据类型**，Go 会自动推断类型：
```go
变量名字 := 表达式
```

```go
var x = 100

func main() {
	fmt.Println(&x, x)
	x := "abc"
	fmt.Println(&x, x)
}
```
上面的代码中 `x := "abc"` 相当于重新定义并初始化了同名的局部变量 `x`，所以打印出来的结果完全不同。

如何避免重新定义，首先要在同一个作用域中，至少有一个新的变量被定义：
```go
func main() {
	x := 100
	fmt.Println(&x, x)
	x, y := 200, 300   // 一个新的变量 y，这里的简短声明就是赋值操作
	fmt.Println(&x, x)
}
```

## make 和 new
1. **`make` 只能用于内建类型（`map`、`slice` 和`channel`）的内存分配。`new` 用于各种类型的内存分配**。
2. **`make` 返回初始化后的（非零）值**。
3. **`new` 返回指针**。

## 类型转换
**Go 强制使用显示类型转换**。这样可以确定语句和表达式的明确含义。**类型转换在编译期完成，包括强制转换和隐式转换**。

类型转换用于将一种数据类型的变量转换为另外一种类型的变量：
```go
类型名(表达式)
```

实例：
```go
var sum int = 17
var count int = 5
var mean float32

mean = float32(sum)/float32(count)
fmt.Printf("mean 的值为: %f\n",mean)
```

对于整数类型值、整数常量之间的类型转换，原则上只要源值在目标类型的可表示范围内就是合法的。比如，`uint8(255)` 
可以把无类型的常量 255 转换为 `uint8` 类型的值，是因为 255 在 `[0, 255]` 的范围内。

这种类型转换主要在**切换同一基础类型不同精度范围**时使用，比如我们要将 `int` 型转为 `int64` 类型时。

```go
a := 100
b := byte(a)
c := a + int(b) 混合类型表达式，类型必须保持一致
```
在 Go 中，非布尔值不能当做 `true/false` 使用，这点和我常用的js不同：
```go
x := 100

if x { // 错误 x 不是布尔值

}
```

如果**要转换为指针类型，或者单向 `channel`，或者函数，要给类型加上 `()`，避免编译器分析错误**，如：
```go
x := 100
(*int)(&x) // *int 加括号，否则会被解析为*(int(&x))

(<- channel int)(c)
(func())(f)
(func()int)(f) // 有返回值的函数其实可以不加括号，但是加括号的话，语义清晰
```

## 类型断言

断言，顾名思义就是果断的去猜测一个未知的事物。在 go 语言中，`interface{}` 就是这个神秘的未知类型，其**断言操作就是用来
判断 `interface{}` 的类型**。因为 `interface{}` 是个未知类型，在编译时无法确定，所以类型断言在运行时确定。

Go 语言里面有一个语法，可以直接**判断是否是该类型的变量：`value, ok = x.(T)`**，这里 `value` 就是变量的值，`ok` 是
一个 `bool` 类型，`x` 是 `interface{}` 变量，`T` 是断言的类型。

该语法返回两个参数，第一个参数是 `x` 转化为 `T` 类型后的变量，第二个值是一个布尔值，若为 `true` 则表示断言成功，
`false` 则表示断言失败。

```go
// comma-ok
for index, element := range list {
	if value, ok := element.(int); ok {
		fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
	} else if value, ok := element.(string); ok {
		fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
	} else if value, ok := element.(Person); ok {
		fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
	} else {
		fmt.Printf("list[%d] is of a different type\n", index)
	}
}


// 或者 使用 switch
for index, element := range list{
	switch value := element.(type) {
		case int:
			fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
		case string:
			fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
		case Person:
			fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
		default:
			fmt.Println("list[%d] is of a different type", index)
	}
}
```

**注意，`x.(type)` 语法不能在 `switch` 外的任何逻辑里面使用，如果你要在 `switch` 外面判断一个类型就使用 `comma-ok`**。

## 零值
“零值”，所指并非是空值，而是一种“变量未填充前”的默认值，通常为 `0`：
```
int     0
int8    0
int32   0
int64   0
uint    0x0
rune    0 //rune的实际类型是 int32
byte    0x0 // byte的实际类型是 uint8
float32 0 //长度为 4 byte
float64 0 //长度为 8 byte
bool    false
string  ""
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

## 编码
Go 语言的源码文件必须使用 UTF-8 编码格式进行存储。如果源码文件中出现了非 UTF-8 编码的字符，那么在构建、安装以及运行的时候，
go 命令就会报告错误“illegal UTF-8 encoding”。

### ASCII 编码
ASCII 是英文“American Standard Code for Information Interchange”的缩写，中文译为美国信息交换标准代码。

ASCII 编码方案使用单个字节（byte）的二进制数来编码一个字符。标准的 ASCII 编码用一个字节的最高比特（bit）位作为奇偶校验位，
而扩展的 ASCII 编码则将此位也用于表示字符。ASCII 编码支持的可打印字符和控制字符的集合也被叫做 ASCII 编码集。

### unicode 编码
**unicode 编码规范，实际上是另一个更加通用的、针对书面字符和文本的字符编码标准。它为世界上现存的所有自然语言中的每一个字符，
都设定了一个唯一的二进制编码**。它定义了不同自然语言的文本数据在国际间交换的统一方式，并为全球化软件创建了一个重要的基础。

Unicode 编码规范通常使用十六进制表示法来表示 Unicode 代码点的整数值，并使用“U+”作为前缀。比如，英文字母字符“a”的 Unicode
代码点是`U+0061`。在 Unicode 编码规范中，一个字符能且只能由与它对应的那个代码点表示。

Unicode 编码规范提供了三种不同的编码格式，即：`UTF-8`、`UTF-16`和`UTF-32`。其中的 UTF 是 UCS Transformation Format 的缩写。
而 UCS 又是 Universal Character Set 的缩写，但也可以代表 Unicode Character Set。所以，UTF 也可以被翻译为 Unicode 转换格式。
它代表的是字符与字节序列之间的转换方式。

在这几种编码格式的名称中，**“-”右边的整数的含义是，以多少个比特位作为一个编码单元**。以`UTF-8`为例，它会以 8 个比特，也就是一个字节，
作为一个编码单元。它与标准的 ASCII 编码是完全兼容的。也就是说，在`[0x00, 0x7F]`的范围内，这两种编码表示的字符都是相同的。
这也是 UTF-8 编码格式的一个巨大优势。

**UTF-8 是一种可变宽的编码方案**。换句话说，**它会用一个或多个字节的二进制数来表示某个字符，最多使用四个字节**。比如，对于一个英文字符，
它仅用一个字节的二进制数就可以表示，而对于一个中文字符，它需要使用三个字节才能够表示。不论怎样，一个受支持的字符总是可以由 UTF-8 编码
为一个字节序列。以下会简称后者为 UTF-8 编码值。