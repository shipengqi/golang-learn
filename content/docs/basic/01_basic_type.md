---
title: 基础数据类型
weight: 1
---

## 数值类型

### 整型

- `uint`，无符号 32 或 64 位整型
- `uint8`，无符号 8 位整型 (0 到 255)
- `uint16`，无符号 16 位整型 (0 到 65535)
- `uint32`，无符号 32 位整型 (0 到 4294967295)
- `uint64`，无符号 64 位整型 (0 到 18446744073709551615)
- `int`，有符号 32 或 64 位整型
- `int8`，有符号 8 位整型 (-128 到 127)
- `int16`，有符号 16 位整型 (-32768 到 32767)
- `int32`，有符号 32 位整型 (-2147483648 到 2147483647)
- `int64`，有符号 64 位整型 (-9223372036854775808 到 9223372036854775807)

`int` 和 `uint` 对应的是 CPU 平台机器的字大小。

### 浮点数

`float32` 和 `float64` 的算术规范由 IEEE-754 浮点数国际标准定义。

- `float32`，32 位浮点型数，`math.MaxFloat32` 表示 `float32` 能表示的最大数值，大约是 `3.4e38`。
- `float64`，64 位浮点型数，`math.MaxFloat64` 表示 `float64` 能表示的最大数值，大约是 `1.8e308`。

### 复数

- `complex64`，对应 `float32` 浮点数精度。
- `complex128`，对应 `float64` 浮点数精度。

内置 `complex` 函数创建复数。标准库 `math/cmplx` 提供了处理复数的函数。

### 其他数值类型

- **`byte`，`uint8`的别名**，一般用于强调数值是一个原始的数据而不是一个小的整数。
- **`rune`，`int32`的别名**，通常用于表示一个 `Unicode` 码点。
- `uintptr`，无符号整型，没有指定具体的 `bit` 大小，用于存放一个指针。

## 布尔类型

布尔类型的值只有两种：`true` 和 `false`。

## 字符串

字符串实际上是由**字符组成的数组**，C 语言中的字符串使用字符数组 `char[]` 表示。数组会占用一片连续的内存空间，而内存空间存储的字节共同组成了字符串，Go 中的字符串只是一个**只读的字节数组**。

字符串的结构体：

```go
// src/reflect/value.go#L1983
type StringHeader struct {
	Data uintptr
	Len  int
}
```

与切片的结构体很像，只不过少了一个容量 `Cap`。

因为字符串是一个只读的类型，不可以直接向字符串直接追加元素改变其本身的内存空间，所有在**字符串上的写入操作都是通过拷贝实现的**。

### 字符串拼接

拼接字符串的几种方式：

#### `+` 拼接字符串

例如 `fmt.Println("hello" + s[5:])` 输出 `"hello, world"`。使用 `+` 来拼接两个字符串时，它会申请一块新的内存空间，大小是两个字符串的大小之和。拼接第三个字符串时，再申请一块新的内存空间，大小是三个字符串大小之和。这种方式每次运算都需要重新分配内存，会给内存分配和 GC 带来额外的负担，所以性能比较差。

#### fmt.Sprintf

`fmt.Sprintf()` 拼接字符串，内部使用 `[]byte` 实现，不像直接运算符这种会产生很多临时的字符串，但是内部的逻辑比较复杂，有很多额外的判断，还用到了 `interface`，所以性能一般。

#### bytes.Buffer

利用 `bytes.Buffer` 拼接字符串，是比较理想的一种方式。对内存的增长有优化，如果能预估字符串的长度，还可以用 `buffer.Grow` 接口来设置 `capacity`。

```go
var buffer bytes.Buffer
buffer.WriteString("hello")
buffer.WriteString(", ")
buffer.WriteString("world")

fmt.Print(buffer.String())
```

#### strings.Builder

`strings.Builder` 内部通过 `slice` 来保存和管理内容。`strings.Builder` 是非线程安全，性能上和 `bytes.Buffer` 相差无几。

```go
var b1 strings.Builder
b1.WriteString("ABC")
b1.WriteString("DEF")

fmt.Print(b1.String())
```

`Builder.Grow` 方法可以预分配内存。

推荐使用 `strings.Builder` 来拼接字符串。

`strings.Builder` 性能上比 `bytes.Buffer` 略快，一个比较重要的区别在于，`bytes.Buffer` 转化为字符串时重新申请了一块空间，存放生成的字符串变量，而 `strings.Builder` 直接将底层的 `[]byte` 转换成了字符串类型并返回。

`bytes.Buffer`：

```go
func (b *Buffer) String() string {
	if b == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(b.buf[b.off:])
}
```

`strings.Builder`：

```go
func (b *Builder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}
```

### 类型转换

在日常开发中，`string` 和 `[]byte` 之间的转换是很常见的，不管是 `string` 转 `[]byte` 还是 `[]byte` 转 `string` 都需要拷贝数据，而内存拷贝带来的性能损耗会随着字符串和 `[]byte` 长度的增长而增长。

## interface{} 和 any

`interface{}` 和 `any` 都是 Go 语言中表示**任意类型**的类型，但是它们的含义和使用方式有一些不同的背景和语境。

- `interface{}` 是 Go 语言的一个**空接口类型，表示没有方法集合的接口**。任何类型都实现了空接口，因为空接口没有要求具体实现任何方法。换句话说，**`interface{}` 可以持有任何类型的值**。
- `any` 是 Go 1.18 引入的一个新的别名，它是 `interface{}` 的类型**别名**。从语义上讲，`any` 和 `interface{}` 是等价的，但 `any` 是为了增强代码的可读性和清晰度。在新的 Go 代码中，使用 `any` 可以更明确地表达类型含义，避免误解。

```go
var x interface{}
x = 42       // 可以存储 int 类型
x = "hello"  // 也可以存储 string 类型
x = true     // 甚至可以存储 bool 类型

var x2 any
x2 = 42       // 可以存储 int 类型
x2 = "hello"  // 也可以存储 string 类型
x2 = true     // 甚至可以存储 bool 类型
```
