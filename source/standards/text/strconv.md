---
title: strconv
---

# strconv

`strconv` 包包含了一系列字符串与相关的类型转换的函数。

## 转换错误处理

`strconv` 中的错误处理。

由于将字符串转为其他数据类型可能会出错，*strconv* 包定义了两个 *error* 类型的变量：*ErrRange* 和 *ErrSyntax*。其中，*ErrRange* 表示
值超过了类型能表示的最大范围，比如将 "128" 转为 `int8` 就会返回这个错误；*ErrSyntax* 表示语法错误，比如将 `""` 转为 `int` 类型会返
回这个错误。

然而，在返回错误的时候，通过构造一个 *NumError* 类型的 *error* 对象返回。*NumError* 结构的定义如下：
```go
// A NumError records a failed conversion.
type NumError struct {
    Func string // the failing function (ParseBool, ParseInt, ParseUint, ParseFloat)
    Num  string // the input
    Err  error  // the reason the conversion failed (ErrRange, ErrSyntax)
}
```

该结构记录了转换过程中发生的错误信息。该结构不仅包含了一个 *error* 类型的成员，记录具体的错误信息，包的实现中，定义了两个便捷函数，
用于构造 *NumError* 对象：
```go
func syntaxError(fn, str string) *NumError {
    return &NumError{fn, str, ErrSyntax}
}

func rangeError(fn, str string) *NumError {
    return &NumError{fn, str, ErrRange}
}
```


## 字符串转为整型

包括三个函数：`ParseInt`、`ParseUint` 和 `Atoi`，函数原型如下：
```go
// 转为有符号整型
func ParseInt(s string, base int, bitSize int) (i int64, err error)
// 转为无符号整型
func ParseUint(s string, base int, bitSize int) (n uint64, err error)
func Atoi(s string) (i int, err error)
```

**`Atoi` 内部通过调用 `ParseInt(s, 10, 0)` 来实现的**。

### ParseInt
```go
func ParseInt(s string, base int, bitSize int) (i int64, err error)
```
参数：
-  `base` 进制，取值为 `2~36`，如果 `base` 的值为 0，则会根据字符串的前缀来确定 `base` 的值：
"0x" 表示 16 进制； "0" 表示 8 进制；否则就是 10 进制。
- `bitSize` 表示的是整数取值范围，或者说整数的具体类型。取值 0、8、16、32 和 64 分别代表 `int`、`int8`、`int16`、`int32` 和 `int64`。
当 `bitSize == 0` 时。

Go 中，`int/uint` 类型，不同系统能表示的范围是不一样的，目前的实现是，32 位系统占 4 个字节；64 位系统占 8 个字节。当 `bitSize==0` 时，
应该表示 32 位还是 64 位呢？`strconv.IntSize` 变量用于获取程序运行的操作系统平台下 int 类型所占的位数。

下面的代码 n 和 err 的值分别是什么？
```go
n, err := strconv.ParseInt("128", 10, 8)
```
	
在 `ParseInt/ParseUint` 的实现中，如果字符串表示的整数超过了 `bitSize` 参数能够表示的范围，则会返回 `ErrRange`，同时会
返回 `bitSize` 能够表示的最大或最小值。

`int8` 占 8 位，最高位代表符号位 （1-负号；0-正号）。所以这里 n 是 127。

另外，`ParseInt` 返回的是 `int64`，这是为了能够容纳所有的整型，在实际使用中，可以根据传递的 `bitSize`，然后将结果转为实际需要的类型。


## 整型转为字符串

遇到需要将字符串和整型连接起来，在 Go 语言中，你需要将整型转为字符串类型，然后才能进行连接。
```go
// 无符号整型转字符串
func FormatUint(i uint64, base int) string	
// 有符号整型转字符串
func FormatInt(i int64, base int) string	
func Itoa(i int) string
```

**`Itoa` 内部直接调用 *FormatInt(i, 10)* 实现的**。


除了使用上述方法将整数转为字符串外，经常见到有人使用 *fmt* 包来做这件事。如：
```go
fmt.Sprintf("%d", 127)
```

那么，这两种方式我们该怎么选择呢？我们主要来考察一下性能。
```go
startTime := time.Now()
for i := 0; i < 10000; i++ {
    fmt.Sprintf("%d", i)
}   
fmt.Println(time.Now().Sub(startTime))

startTime := time.Now()
for i := 0; i < 10000; i++ {
    strconv.Itoa(i)
}   
fmt.Println(time.Now().Sub(startTime))
```
`Sprintf` 的时间是 3.549761ms，而 `Itoa` 的时间是 848.208us，相差 4 倍多。

`Sprintf` 性能差些可以预见，因为它接收的是 `interface`，需要进行反射等操作。建议使用 `strconv` 包中的方法进行转换。


## 字符串和布尔值之间的转换

Go 中字符串和布尔值之间的转换比较简单，主要有三个函数：
```go
// 接受 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False 等字符串；
// 其他形式的字符串会返回错误
// 返回转换后的布尔值
func ParseBool(str string) (value bool, err error)
// 直接返回 "true" 或 "false"
func FormatBool(b bool) string
```


## 字符串和浮点数之间的转换

类似的，包含三个函数：
```go
// 无论 bitSize 取值如何，函数返回值类型都是 float64。
func ParseFloat(s string, bitSize int) (f float64, err error)
// 第一个参数是输入浮点数
// 第二个是浮点数的显示格式（可以是 b, e, E, f, g, G）
func FormatFloat(f float64, fmt byte, prec, bitSize int) string
```

`prec` 表示有效数字（对 `fmt='b'` 无效），对于 'e', 'E' 和 'f'，有效数字用于小数点之后的位数；
对于 'g' 和 'G'，则是所有的有效数字。例如：
```go
strconv.FormatFloat(1223.13252, 'e', 3, 32)	// 结果：1.223e+03
strconv.FormatFloat(1223.13252, 'g', 3, 32)	// 结果：1.22e+03
```


由于浮点数有精度的问题，精度不一样，`ParseFloat` 和 `FormatFloat` 可能达不到互逆的效果。如：
```go
s := strconv.FormatFloat(1234.5678, 'g', 6, 64)
strconv.ParseFloat(s, 64)
```

另外，`fmt='b'` 时，得到的字符串是无法通过 `ParseFloat` 还原的。

同样的，基于性能的考虑，应该使用 `FormatFloat` 而不是 `fmt.Sprintf`。
