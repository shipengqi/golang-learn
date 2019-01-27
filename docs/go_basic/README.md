# 介绍
Go语言非常简单，只有25个关键字：

- `var`和`const`声明变量和常量
- `package`和`import`声明所属包名和导入包。
- `func` 用于定义函数和方法
- `return` 用于从函数返回
- `defer` 用于类似析构函数
- `go` 用于并发
- `select` 用于选择不同类型的通讯
- `interface` 用于定义接口
- `struct` 用于定义抽象数据类型
- `break`、`case`、`continue`、`for`、`fallthrough`、`else`、`if`、`switch`、`goto`、`default`流程控制语句
- `chan`用于`channel`通讯
- `type`用于声明自定义类型
- `map`用于声明`map`类型数据
- `range`用于读取`slice`、`map`、`channel`数据

## 三种文件
- 命令源码文件，如果一个源码文件声明属于`main`包，并且包含一个无参数声明且无结果声明的`main`函数，那么它就是命令源码文件。
- 库源码文件，库源码文件是不能被直接运行的源码文件，它仅用于存放程序实体，这些程序实体可以被其他代码使用
- 测试源码文件

## 命名
所有命名只能以字母或者`_`开头，可以包含字母，数字或者`_`。区分大小写。
关键字不能定义变量名，如`func`，`default`。

**注意：函数内部定义的，只能在函数内部使用（函数级），在函数外部定义的（包级），可以在当前包的所有文件中是使用。
并且，是在函数外定义的名字，如果以大写字母开头，那么会被导出，也就是在包的外部也可以访问，所以定义名字时，要注意大小写。**

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
**类型**或者**表达式**可以省略其中的一个。也就是如果没有类型，可以通过表达式推断出类型，**没有表达式，将会根据类型初始化为对应的零值**。
对应关系：
- 数值类型：`0`
- 布尔类型：`false`
- 字符串：`""`
- 接口或引用类型（包括slice、指针、map、chan和函数）：`nil`

### 声明一组变量
```go
var 变量名字, 变量名字, 变量名字 ... 类型 = 表达式, 表达式, 表达式, ...
```
比如：
```go
// 声明一组`int`类型
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
**`:=`只能在函数内使用，不能提供数据类型**，Go 会自动推断类型：
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
上面的代码中`x := "abc"`相当于重新定义并初始化了同名的局部变量`x`，所以打印出来的结果完全不同。

如何避免重新定义，首先要在同一个作用域中，至少有一个新的变量被定义：
```go
func main() {
	x := 100 
	fmt.Println(&x, x)
	x, y := 200, 300   // 一个新的变量 y，这里的简短声明就是赋值操作
	fmt.Println(&x, x)
}
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
当一个指针被定义后没有分配到任何变量时，它的值为`nil`。`nil`指针也称为空指针。

#### 指向指针的指针
```go
var a int
var ptr *int
var pptr **int

a = 3000

/* 指针 ptr 地址 */
ptr = &a

/* 指向指针 ptr 地址 */
pptr = &ptr

/* 获取 pptr 的值 */
fmt.Printf("变量 a = %d\n", a )
fmt.Printf("指针变量 *ptr = %d\n", *ptr )
fmt.Printf("指向指针的指针变量 **pptr = %d\n", **pptr)
```

#### 为什么需要指针
相比 Java，Python，Javascript 等引用类型的语言，Golang 拥有类似C语言的指针这个相对古老的特性。但不同于 C 语言，Golang 的指针是单独的类型，而不是 C 语言中的~int·类型，
而且也不能对指针做整数运算。从这一点看，Golang 的指针基本就是一种引用。

在学习引用类型语言的时候，总是要先搞清楚，当给一个`函数/方法`传参的时候，传进去的是值还是引用。实际上，在大部分引用型语言里，参数为基本类型时，传进去的大都是值，
也就是另外复制了一份参数到当前的函数调用栈。参数为高级类型时，传进去的基本都是引用。

内存管理中的内存区域一般包括`heap`和`stack`，`stack`主要用来存储当前调用栈用到的简单类型数据：`string`，`boolean`，`int`，`float`等。这些类型的内存占用小，容易回收，
基本上它们的值和指针占用的空间差不多，因此可以直接复制，`GC`也比较容易做针对性的优化。复杂的高级类型占用的内存往往相对较大，存储在`heap`中，`GC`回收频率相对较低，代价也较大，
因此传`引用/指针`可以避免进行成本较高的复制操作，并且节省内存，提高程序运行效率。

因此，在下列情况可以考虑使用指针：
1. **需要改变参数的值**
2. **避免复制操作**
3. **节省内存**

而在 Golang 中，具体到高级类型`struct`，`slice`，`map`也各有不同。实际上，只有`struct`的使用有点复杂，`slice`，`map`，`chan`都可以直接使用，不用考虑是值还是指针。

##### `struct`

对于函数（`function`），由函数的参数类型指定，传入的参数的类型不对会报错，例如：
```go
func passValue(s struct){}

func passPointer(s *struct){}
```

对于方法（`method`），接收者（`receiver`）可以是指针，也可以是值，Golang 会在传递参数前自动适配以符合参数的类型。也就是：如果方法的参数是值，
那么按照传值的方式 ，方法内部对`struct`的改动无法作用在外部的变量上，例如：
```go
package main

import "fmt"

type MyPoint struct {
    X int
    Y int
}

func printFuncValue(p MyPoint){
    p.X = 1
    p.Y = 1
    fmt.Printf(" -> %v", p)
}

func printFuncPointer(pp *MyPoint){
    pp.X = 1 // 实际上应该写做 (*pp).X，Golang 给了语法糖，减少了麻烦，但是也导致了 * 的不一致
    pp.Y = 1
    fmt.Printf(" -> %v", pp)
}

func (p MyPoint) printMethodValue(){
    p.X += 1
    p.Y += 1
    fmt.Printf(" -> %v", p)
}

// 建议使用指针作为方法（method：printMethodPointer）的接收者（receiver：*MyPoint），一是可以修改接收者的值，二是可以避免大对象的复制
func (pp *MyPoint) printMethodPointer(){
    pp.X += 1
    pp.Y += 1
    fmt.Printf(" -> %v", pp)
}

func main(){
    p := MyPoint{0, 0}
    pp := &MyPoint{0, 0}

    fmt.Printf("\n value to func(value): %v", p)
    printFuncValue(p)
    fmt.Printf(" --> %v", p)
    // Output: value to func(value): {0 0} -> {1 1} --> {0 0}

    //printFuncValue(pp) // cannot use pp (type *MyPoint) as type MyPoint in argument to printFuncValue

    //printFuncPointer(p) // cannot use p (type MyPoint) as type *MyPoint in argument to printFuncPointer

    fmt.Printf("\n pointer to func(pointer): %v", pp)
    printFuncPointer(pp)
    fmt.Printf(" --> %v", pp)
    // Output: pointer to func(pointer): &{0 0} -> &{1 1} --> &{1 1}

    fmt.Printf("\n value to method(value): %v", p)
    p.printMethodValue()
    fmt.Printf(" --> %v", p)
    // Output: value to method(value): {0 0} -> {1 1} --> {0 0}

    fmt.Printf("\n value to method(pointer): %v", p)
    p.printMethodPointer()
    fmt.Printf(" --> %v", p)
    // Output: value to method(pointer): {0 0} -> &{1 1} --> {1 1}

    fmt.Printf("\n pointer to method(value): %v", pp)
    pp.printMethodValue()
    fmt.Printf(" --> %v", pp)
    // Output: pointer to method(value): &{1 1} -> {2 2} --> &{1 1}

    fmt.Printf("\n pointer to method(pointer): %v", pp)
    pp.printMethodPointer()
    fmt.Printf(" --> %v", pp)
    // Output: pointer to method(pointer): &{1 1} -> &{2 2} --> &{2 2}
}
```

##### `slice`
**`slice`实际上相当于对其依附的`array`的引用，它不存储数据，只是对`array`进行描述**。因此，修改`slice`中的元素，改变会体现在`array`上，当然也会体现在该`array`的所有`slice`上。

##### map

使用`make(map[string]string)`返回的本身是个引用，可以直接用来操作：
```go
map["name"]="Jason"；
```

而如果使用`map`的指针，反而会产生错误：
```go
*map["name"]="Jason"  //  invalid indirect of m["title"] (type string)
(*map)["name"]="Jason"  // invalid indirect of m (type map[string]string)
```

#### 哪些值是不可寻址的
1. **不可变的值**不可寻址。常量、基本类型的值字面量、字符串变量的值、函数以及方法的字面量都是如此。其实这样规定也有安全性方面的考虑。
2. 绝大多数被视为**临时结果的值**都是不可寻址的。算术操作的结果值属于临时结果，针对值字面量的表达式结果值也属于临时结果。但有一个例外，
对切片字面量的索引结果值虽然也属于临时结果，但却是可寻址的。函数的返回值也是临时结果。`++`和`--`并不属于操作符。
3. **不安全的值**不可寻址，若拿到某值的指针可能会破坏程序的一致性，那么就是不安全的。由于字典的内部机制，对字典的索引结果值的取址
操作都是不安全的。另外，获取由字面量或标识符代表的函数或方法的地址显然也是不安全的。


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

## 自定义类型
使用`type`自定义类型，一般出现在包一级，与变量一样，如果类型名字的首字母是大写，则在包外部也可以使用：
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

自定义类型虽然置顶了底层类型，但是只是底层数据结构相同，不会继承底层类型的其他信息，比如（方法）。
不能隐式转换，不能直接用于比较表达式。

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

### 点操作
```go
 import(
    . "fmt"
 )
```
这个点操作的含义就是这个包导入之后在你调用这个包的函数时，你可以省略前缀的包名，也就是前面你调用的`fmt.Println("hello world")`
可以省略的写成`Println("hello world")`。

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

#### 声明所属的代码包与其所在目录的名称不同时
源码文件所在的目录相对于 src 目录的相对路径就是它的代码包导入路径，而实际使用其程序实体时给定的限定符要与它声明所属的代码包名称对应。
为了不让该代码包的使用者产生困惑，我们总是应该让声明的包名与其父目录的名称一致。

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
`init`函数不能有任何的参数和返回值，虽然一个`package`里面可以写任意多个`init`函数，但这无论是对于可读性还是以后的可维护性来说，
我们都强烈建议用户在一个`package`中每个文件只写一个`init`函数。

程序的初始化和执行都起始于`main`包。如果`main`包还导入了其它的包，那么就会在编译时将它们依次导入。有时一个包会被多个包同时导入，
那么它只会被导入一次（例如很多包可能都会用到`fmt`包，但它只会被导入一次，因为没有必要导入多次）。当一个包被导入时，如果该包还导入了其它的包，
那么会先将其它包导入进来，然后再对这些包中的包级常量和变量进行初始化，接着执行`init`函数（如果有的话），依次类推。等所有被导入的包都加载完毕了，就
会开始对`main`包中的包级常量和变量进行初始化，然后执行`main`包中的`init`函数（如果存在的话），最后执行`main`函数。

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

在这几种编码格式的名称中，**“-”右边的整数的含义是，以多少个比特位作为一个编码单元**。以`UTF-8`为例，它会以 8 个比特，也就是一个字节，作为一个编码单元。
它与标准的 ASCII 编码是完全兼容的。也就是说，在`[0x00, 0x7F]`的范围内，这两种编码表示的字符都是相同的。这也是 UTF-8 编码格式的一个巨大优势。

**UTF-8 是一种可变宽的编码方案**。换句话说，**它会用一个或多个字节的二进制数来表示某个字符，最多使用四个字节**。比如，对于一个英文字符，它仅用一个字
节的二进制数就可以表示，而对于一个中文字符，它需要使用三个字节才能够表示。不论怎样，一个受支持的字符总是可以由 UTF-8 编码为一个字节序列。以下会简称后
者为 UTF-8 编码值。

### 一个`string`类型的值在底层怎样被表达
在底层，一个string类型的值是由一系列相对应的 Unicode 代码点的 UTF-8 编码值来表达的。

一个string类型的值既可以被拆分为一个包含多个字符的序列，也可以被拆分为一个包含多个字节的序列。
前者可以由一个以`rune`（`int32`的别名）为元素类型的切片来表示，而后者则可以由一个以`byte`为元素类型的切片代表。

`rune`是 Go 语言特有的一个基本数据类型，它的一个值就代表一个字符，即：一个 Unicode 字符。比如，'G'、'o'、'爱'、'好'、'者'代表的就都是一个 Unicode 字符。
一个`rune`类型的值会由四个字节宽度的空间来存储。它的存储空间总是能够存下一个 UTF-8 编码值。

**一个`rune`类型的值在底层其实就是一个 UTF-8 编码值**。前者是（便于我们人类理解的）外部展现，后者是（便于计算机系统理解的）内在表达。

```go
str := "Go 爱好者 "
fmt.Printf("The string: %q\n", str)
fmt.Printf("  => runes(char): %q\n", []rune(str))
fmt.Printf("  => runes(hex): %x\n", []rune(str))
fmt.Printf("  => bytes(hex): [% x]\n", []byte(str))
```
字符串值"Go 爱好者"如果被转换为`[]rune`类型的值的话，其中的每一个字符（不论是英文字符还是中文字符）就都会独立成为一个`rune`类型的元素值。因
此，这段代码打印出的第二行内容就会如下所示：
```bash
=> runes(char): ['G' 'o' '爱' '好' '者']
```
又由于，每个`rune`类型的值在底层都是由一个 UTF-8 编码值来表达的，所以我们可以换一种方式来展现这个字符序列：
```bash
=> runes(hex): [47 6f 7231 597d 8005]
```
我们还可以进一步地拆分，把每个字符的 UTF-8 编码值都拆成相应的字节序列。上述代码中的第五行就是这么做的。它会得到如下的输出：
```bash
=> bytes(hex): [47 6f e7 88 b1 e5 a5 bd e8 80 85]
```

### 使用带有`range`子句的`for`语句遍历字符串值的时候应该注意
带有`range`子句的`for`语句会先把被遍历的字符串值拆成一个**字节序列**（注意是字节序列），然后再试图找出这个字节序列中包含的每一个 UTF-8 编码值，
或者说每一个 Unicode 字符。

这样的`for`语句可以为两个迭代变量赋值。如果存在两个迭代变量，那么赋给第一个变量的值就将会是当前字节序列中的某个 UTF-8 编码值的第一个字节所
对应的那个索引值。而赋给第二个变量的值则是这个 UTF-8 编码值代表的那个 Unicode 字符，其类型会是`rune`。

```go
str := "Go 爱好者 "
for i, c := range str {
    fmt.Printf("%d: %q [% x]\n", i, c, []byte(string(c)))
}
```
完整的打印内容如下：
```bash
0: 'G' [47]
1: 'o' [6f]
2: '爱' [e7 88 b1]
5: '好' [e5 a5 bd]
8: '者' [e8 80 85]
```

**注意了，'爱'是由三个字节共同表达的，所以第四个 Unicode 字符'好'对应的索引值并不是3，而是2加3后得到的5**。

## strings包与字符串操作
```go
/*字符串基本操作--strings*/
str := "wangdy"
//是否包含
fmt.Println(strings.Contains(str, "wang"), strings.Contains(str, "123")) //true false
//获取字符串长度
fmt.Println(len(str)) //6
//获取字符在字符串的位置 从0开始,如果不存在，返回-1
fmt.Println(strings.Index(str, "g")) //3
fmt.Println(strings.Index(str, "x")) //-1
//判断字符串是否以 xx 开头
fmt.Println(strings.HasPrefix(str, "wa")) //true
//判断字符串是否以 xx 结尾
fmt.Println(strings.HasSuffix(str, "dy")) //true
//判断2个字符串大小，相等0，左边大于右边-1，其他1
str2 := "hahaha"
fmt.Println(strings.Compare(str, str2)) //1
//分割字符串
strSplit := strings.Split("1-2-3-4-a", "-")
fmt.Println(strSplit) //[1 2 3 4 a]
//组装字符串
fmt.Println(strings.Join(strSplit, "#")) //1#2#3#4#a
//去除字符串2端空格
fmt.Printf("%s,%s\n", strings.Trim("  我的2边有空格   1  ", " "), "/////") //我的2边有空格   1,/////
//大小写转换
fmt.Println(strings.ToUpper("abDCaE")) //ABDCAE
fmt.Println(strings.ToLower("abDCaE")) //abdcae
//字符串替换:意思是：在sourceStr中，把oldStr的前n个替换成newStr，返回一个新字符串，如果n<0则全部替换
sourceStr := "123123123"
oldStr := "12"
newStr := "ab"
n := 2
fmt.Println(strings.Replace(sourceStr, oldStr, newStr, n))
```

在 Go 语言中，**string类型的值是不可变的。如果我们想获得一个不一样的字符串，那么就只能基于原字符串进行裁剪、拼接等操作，
从而生成一个新的字符串**。裁剪操作可以使用切片表达式，而拼接操作可以用操作符`+`实现。

在底层，一个string值的内容会被存储到一块连续的内存空间中。同时，这块内存容纳的字节数量也会被记录下来，并用于表示该string值的长度。

你可以把这块内存的内容看成一个字节数组，而相应的`string`值则包含了指向字节数组头部的指针值。如此一来，**我们在一个`string`值上应用切片表达式，
就相当于在对其底层的字节数组做切片**。

另一方面，我们在**进行字符串拼接的时候，Go 语言会把所有被拼接的字符串依次拷贝到一个崭新且足够大的连续内存空间中，并把持有相应指针值的`string`值作为结果返回**。

显然，当**程序中存在过多的字符串拼接操作的时候，会对内存的分配产生非常大的压力**。

### 与`string`值相比，`strings.Builder`类型的值有哪些优势
- 已存在的内容不可变，但可以拼接更多的内容；
- 减少了内存分配和内容拷贝的次数；
- 可将内容重置，可重用值。

`Builder`值中有一个用于承载内容的容器（以下简称内容容器）。它是一个以`byte`为元素类型的切片（以下简称字节切片）。

**由于这样的字节切片的底层数组就是一个字节数组，所以我们可以说它与string值存储内容的方式是一样的**。实际上，它们都是通过一个`unsafe.Pointer`类型的字段
来持有那个指向了底层字节数组的指针值的。

因为这样的内部构造，`Builder`值同样拥有高效利用内存的前提条件。

已存在于`Builder`值中的内容是不可变的。因此，我们可以利用`Builder`值提供的方法拼接更多的内容，而丝毫不用担心这些方法会影响到已存在的内容。

`Builder`值拥有的一系列指针方法，包括：`Write`、`WriteByte`、`WriteRune`和`WriteString`。我们可以把它们统称为**拼接方法**。

调用上述方法把新的内容拼接到已存在的内容的尾部（也就是右边）。这时，如有必要，`Builder`值会自动地对自身的内容容器进行扩容。这里的自动扩容策略与切片的扩容策略一致。

除了Builder值的自动扩容，我们还可以选择手动扩容，这通过调用`Builder`值的`Grow`方法就可以做到。`Grow`方法也可以被称为**扩容方法**，它接受
一个`int`类型的参数`n`，该参数用于代表将要扩充的字节数量。

`Grow`方法会把其所属值中内容容器的容量增加`n`个字节。更具体地讲，它会生成一个字节切片作为新的内容容器，该切片的容量会是原容器容量的二倍再加上`n`。之
后，它会把原容器中的所有字节全部拷贝到新容器中。

### 使用`strings.Builder`类型的约束
**只要调用了`Builder`值的拼接方法或扩容方法，就不能再以任何的方式对其所属值进行复制了**。否则，只要在任何副本上调用上述方法就都会引发 panic。
这里所说的复制方式，包括但不限于在函数间传递值、通过通道传递值、把值赋予变量等等。

正是由于已使用的`Builder`值不能再被复制，所以肯定不会出现多个`Builder`值中的内容容器（也就是那个字节切片）共用一个底层字节数组的情况。这样也就避免
了多个同源的`Builder`值在拼接内容时可能产生的冲突问题。

**不过，虽然已使用的`Builder`值不能再被复制，但是它的指针值却可以。无论什么时候，我们都可以通过任何方式复制这样的指针值**。注意，这样的指针值指向的都会
是同一个`Builder`值。

### `strings.Reader`类型
`strings.Reader`类型是为了高效读取字符串而存在的。可以让我们很方便地读取一个字符串中的内容。在读取的过程中，`Reader`值会保存已读取的字节的计数（以下简称已读计数）。

**已读计数也代表着下一次读取的起始索引位置。Reader值正是依靠这样一个计数，以及针对字符串值的切片表达式，从而实现快速读取**。