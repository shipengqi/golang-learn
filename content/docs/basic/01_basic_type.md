---
title: 基础据类型
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

字符串是只读的，虽然可以通过 `string[index]` 获取某个索引位置的字节值，但是不能使用 `string[index] = "string2"` 这种方式来修改字符串。可以将其转为可变类型（`[]rune` 或 `[]byte`），完成后再转回来。

拼接字符串的几种方式：

#### `+` 拼接字符串

例如 `fmt.Println("hello" + s[5:])` 输出 `"hello, world"`。这种方式每次运算都会产生一个新的字符串，需要重新分配内存，会给内存分配和 GC 带来额外的负担，所以性能比较差。

#### fmt.Sprintf

`fmt.Sprintf()` 拼接字符串，内部使用 `[]byte` 实现，不像直接运算符这种会产生很多临时的字符串，但是内部的逻辑比较复杂，有很多额外的判断，还用到了 `interface`，所以性能一般。

#### strings.Join

`strings.Join()` 拼接字符串，`Join` 会先根据字符串数组的内容，计算出一个拼接之后的长度，然后申请对应大小的内存，一个一个字符串填入，在已有一个数组的情况下，这种效率会很高，但是本来没有，去构造
这个数据的代价也不小。

#### bytes.Buffer

利用 `bytes.Buffer` 拼接字符串，是比较理想的一种方式。对内存的增长有优化，如果能预估字符串的长度，还可以用 `buffer.Grow()` 接口来设置 `capacity`。

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