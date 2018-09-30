# 语法基础

- 可以直接运行（`go run *.go`，其实也是`run`命令进行编译运行），也可以编译后运行。
- 函数可以返回多个值，函数是第一类型，可以作为参数或返回值。
- 控制语句只有三种`if`，`for`，`switch`。

## 注释
使用`//`添加注释。一般我们会在包声明前添加注释，来对整个包挥着程序做整体的描述。

## package
我们知道，我们在写 Go 语言的代码时，每个文件的头部都有一行`package`声明语句。比如`package main`。这个声明表示这个源文件属于哪个包（类似其他语言的`modules`或者`libraries`）。 Go 语言的代码就是通过这个`package`来组织。

## 行分隔符
Go 中，一行代表一个语句结束，不需要以分号`;`结尾。多个语句写在同一行，则必须使用`;`（不推荐使用）。 

## import
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
`_`代表空标识符，Go 不允许有无用的变量，空标识符可以作为忽略占位符，比如：
```go
var s, sep string
for _, arg := range os.Args[1:] {
	s += sep + arg
	sep = " "
}
```

## 常量
`const`声明常量，运行时不可改变（只读），注意常量的**底层数据类型只能是基础类型（布尔型、数值型和字符串型）**：
```go
const 常量名字 类型 = 表达式
```

类型 可以省略。也就是如果没有类型，可以通过表达式推导出类型。

比如：
```go
// 声明一个`string`类型
const b string = "abc"
const a = "abc"

// 声明一组不同类型
const c, f, s = true, 2.3, "four" // bool, float64, string

// 批量声明多个常量
const (
  Unknown = 0
  Female = 1
  Male = 2
)

const strSize = len("hello, world")
```
常量表达式的值在**编译期计算**。因此常量表达式中，函数必须是内置函数。如`unsafe.Sizeof()`，`len()`, `cap()`**。

常量组中，如果不指定类型和初始值，那么就和上一行非空常量右值相同：
例如：
```go
const (
	a = 1
	b
	c = 2
	d
)

fmt.Println(a, b, c, d) // "1 1 2 2"
```

### iota
Go 中没有枚举的定义，但是可以使用`iota`，`iota`标识符可以认为是一个可以被编译器修改的常量。
在`const`声明中，被重置为0，在第一个声明的常量所在的行，`iota`将会被置为`0`，然后在每一个有常量声明的行加`1`。
```go
const (
	a = iota   //0
	b          //1
	c          //2
	d = "ha"   //"ha", iota += 1
	e          //"ha"   iota += 1
	f = 100    //100, iota +=1
	g          //100  iota +=1
	h = iota   //7, 中断的 iota 计数必须显示恢复
	i          //8
)

const (
	i = 1 << iota //1, 1 << 0
	j = 3 << iota //6, 3 << 1
	k             //12, 3 << 2
	l             //24, 3 << 3
)
```
## 运算符
### 优先级

1. `*`，`/`，`%`，`<<`，`>>`，`&`，`&^`
2. `+`，`-`，`|`，`^`
3. `==`，`!=`，`<`，`<=`，`>`，`>=`
4. `&&`
5. `||`

上面的运算符得优先级，从上到下，从左到右。也就是`*`的优先级最高，`||`的优先级最低。

### 算术运算符
`+`、`-`、`*`和`/`可以适用于整数、浮点数和复数。

在 Go 中，`%`取模运算符的符号和被取模数的符号总是一致的，因此`-5 % 3`和`-5 % -3`结果都是`-2`。`%`仅用于整数间的运算。
除法运算符`/`的行为则依赖于操作数是否为全为整数，比如`5.0/4.0`的结果是`1.2`5，但是`5/4`的结果是`1`，因为整数除法会向着`0`方向截断余数。

`++`自增，`--`自减

### 关系运算符
`==`，`!=`，`<`，`<=`，`>`，`>=`。

布尔型、数字类型和字符串等基本类型都是可比较的，也就是说两个相同类型的值可以用`==`和`!=`进行比较。
### 逻辑运算符
`&&`，`||`，`!`（逻辑 NOT 运算符）。

### 位运算符
`&`，`|`，`^`，`<<`，`>>`，`&^`（位清空 AND NOT）

`&^`：如果对应`y`中`bit`位为`1`的话, 表达式`z = x &^ y`结果`z`的对应的`bit`位为`0`，否则`z`对应的`bit`位等于`x`相应的`bit`位的值。如：
```go
var x uint8 = 00100010
var y uint8 = 00000110
fmt.Printf("%08b\n", x&^y) // "00100000"
```
### 赋值运算符
除了`=`外，还有`+=`（相加后再赋值），`-=`（相减后再赋值），`*=`（相乘后再赋值）等等，其他的赋值运算符也都是一个套路。

### 其他运算符
`&`（取地址操作），`*`（指针变量）。

## 类型转换
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

## 条件语句
### if
```go
if 布尔表达式 {
   
}
```
### if...else

```go
if 布尔表达式 {
   
} else {
  
}
```

### switch
```go
switch var1 {
    case val1:
        ...  // 不需要显示的break，case 执行完会自动中断
    case val2:
				...
		case val3,val4,...:		
    default:
        ...
}
```
`val1`,`val2` ... 类型不被局限于常量或整数，但必须是相同的类型。

### select
`select`类似于用于通信的`switch`语句。每个`case`必须是一个通信操作，要么是发送要么是接收。

当条件满足时，`select`会去通信并执行`case`之后的语句，这时候其它通信是不会执行的。
如果多个`case`同时满足条件，`select`会随机地选择一个执行。如果没有`case`可运行，它将阻塞，直到有`case`可运行。

一个默认的子句应该总是可运行的。

```go
select {
  case communication clause:
      ...     
  case communication clause:
      ... 
  default: /* 可选 */
			... 
}			
```

## 循环语句
### for
```go
for init; condition; post { }

// 相当于  while (x < 5) { ... }
for x < 5 {
  ...
}

// 相当于 while (true) { ... }
for {
	...
}

for key, value := range oldMap { // 第二个循环变量可以忽略，但是第一个变量要忽略可以使用空标识符 _ 代替
    newMap[key] = value
}
```
`for range`支持遍历数组，切片，字符串，字典，通道，并返回索引和键值。**`for range`会复制目标数据。可改用数组指针或者切片**。

循环控制语句：
- `break`，用于中断当前`for`循环或跳出`switch`语句
- `continue`，跳过当前循，继续进行下一轮循环。
- `goto`，将控制转移到被标记的语句。通常与条件语句配合使用。可用来实现条件转移， 构成循环，跳出循环体等功能。不推荐
使用，以免造成流程混乱。

`goto`实例：
```go
LOOP: for a < 20 {
	if a == 15 {
			/* 跳过迭代 */
			a = a + 1
			goto LOOP
	}
	fmt.Printf("a的值为 : %d\n", a)
	a ++  
}  
```

## JSON
`JavaScript`对象表示法（JSON）是一种用于发送和接收结构化信息的标准协议。Go 对于其他序列化协议如`XML`，`Protocol Buffers`，都有良好的支持，
由标准库中的`encoding/json`、`encoding/xml`、`encoding/asn1`等包提供支持，`Protocol Buffers`的由`github.com/golang/protobuf`包提供支持，
并且这类包都有着相似的API接口。

GO 中结构体转为`JSON`使用`json.Marshal`，也就是编码操作：
```go
type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

var movies = []Movie{
	{
		Title: "Casablanca", 
		Year: 1942, 
		Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{
		Title: "Cool Hand Luke",
		Year: 1967, 
		Color: true,
		Actors: []string{"Paul Newman"}},
	{
		Title: "Bullitt", 
		Year: 1968, 
		Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}}}	

data, err := json.Marshal(movies)
if err != nil {
    log.Fatalf("JSON marshaling failed: %s", err)
}
fmt.Printf("%s\n", data)
```

`json.MarshalIndent`格式化输出`JSON`。

**注意，只有导出的结构体成员才会被编码**

解码操作，使用`json.Unmarshal`：
```go
var titles []struct{ Title string }
if err := json.Unmarshal(data, &titles); err != nil {
    log.Fatalf("JSON unmarshaling failed: %s", err)
}
fmt.Println(titles) // "[{Casablanca} {Cool Hand Luke} {Bullitt}]"
```
通过定义合适的Go语言数据结构，我们可以选择性地解码JSON中感兴趣的成员。

基于流式的解码器`json.Decoder`。针对输出流的`json.Encoder`编码对象