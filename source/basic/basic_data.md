---
title: 基础据类型
---

# 基础据类型

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
- `float32`，IEEE-754 32 位浮点型数，`math.MaxFloat32` 表示 `float32` 能表示的最大数值，大约是 `3.4e38`。
- `float64`，IEEE-754 64 位浮点型数，`math.MaxFloat64` 表示 `float64` 能表示的最大数值，大约是 `1.8e308`。

### 复数
- `complex64`，对应 `float32` 浮点数精度。
- `complex128`，对应 `float64` 浮点数精度。

内置 `complex` 函数创建复数。`math/cmplx` 包提供了复数处理的许多函数。

### 其他数值类型
- `byte`，`uint8`的别名，通常用于表示一个`Unicode`码点。
- `rune`，`int32`的别名，一般用于强调数值是一个原始的数据而不是一个小的整数。
- `uintptr`，无符号整型，用于存放一个指针，没有指定具体的`bit`大小。

## 布尔类型
布尔类型的值只有两种：`true` 和 `false`。

## 字符串
字符串就是一串固定长度的字符连接起来的字符序列，不可改变。Go 的字符串是由单个字节连接起来的。Go 的字符串的字节使
用 `UTF-8` 编码标识 `Unicode` 文本。

### 字符串操作
- 内置函数 `len` 可以获取字符串的长度。
- **可以通过 `string[index]` 获取某个索引位置的字节值，字符串是不可修改的，不能使用 `string[index] = "string2"`
这种方式改变字符串**。
- `string[i, l]` 获取 `string` 从第 `i` 个字节位置开始的 `l` 个字节，返回一个新的字符串。如：
  ```go
    s := "hello, world"
    fmt.Println(s[0:5]) // "hello"
    
    fmt.Println(s[:5]) // "hello"
    fmt.Println(s[7:]) // "world"
    fmt.Println(s[:])  // "hello, world"
  ```
- `+` 拼接字符串，如 `fmt.Println("goodbye" + s[5:])` 输出 `"goodbye, world"`。这种方式每次运算都会产生一个新的字
符串，所以会产生很多临时的无用的字符串，不仅没有用，还会给内存分配和 GC 带来额外的负担，所以性能比较差。
- `fmt.Sprintf()` 拼接字符串，内部使用 `[]byte` 实现，不像直接运算符这种会产生很多临时的字符串，但是内部的逻辑比较复杂，有很
多额外的判断，还用到了 `interface`，所以性能一般。
- `strings.Join()` 拼接字符串，Join会先根据字符串数组的内容，计算出一个拼接之后的长度，然后申请对应大小的内存，一个一个字符
串填入，在已有一个数组的情况下，这种效率会很高，但是本来没有，去构造这个数据的代价也不小。
- `bytes.Buffer` 拼接字符串，比较理想，可以当成可变字符使用，对内存的增长也有优化，如果能预估字符串的长度，还可
以用 `buffer.Grow()` 接口来设置 `capacity`。
```go
var buffer bytes.Buffer
buffer.WriteString("hello")
buffer.WriteString(", ")
buffer.WriteString("world")

fmt.Print(buffer.String())
```
- `strings.Builder` 内部通过 `slice` 来保存和管理内容。`slice` 内部则是通过一个指针指向实际保存内容的数组。`strings.Builder` 
是非线程安全，性能上和 `bytes.Buffer` 相差无几。
```go
var b1 strings.Builder
b1.WriteString("ABC")
b1.WriteString("DEF")

fmt.Print(b1.String())
```
- 使用 `==` 和 `<` 进行字符串比较。

**一个原生的字符串面值形式是 \`...\`，使用反引号代替双引号。在原生的字符串面值中，没有转义操作；全部的内容都是字面的意思，
包含退格和换行**。

### strings 包与字符串操作
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

在 Go 语言中，**`string` 类型的值是不可变的。如果我们想获得一个不一样的字符串，那么就只能基于原字符串进行裁剪、拼接等操作，
从而生成一个新的字符串**。裁剪操作可以使用切片表达式，而拼接操作可以用操作符`+`实现。

在底层，一个 `string` 值的内容会被存储到一块连续的内存空间中。同时，这块内存容纳的字节数量也会被记录下来，并用于表示
该 `string` 值的长度。

你可以把这块内存的内容看成一个字节数组，而相应的 `string` 值则包含了指向字节数组头部的指针值。如此一来，**我们在
一个 `string` 值上应用切片表达式，就相当于在对其底层的字节数组做切片**。

另一方面，我们在**进行字符串拼接的时候，Go 语言会把所有被拼接的字符串依次拷贝到一个崭新且足够大的连续内存空间中，
并把持有相应指针值的 `string` 值作为结果返回**。

显然，当**程序中存在过多的字符串拼接操作的时候，会对内存的分配产生非常大的压力**。

#### 与 `string` 值相比，`strings.Builder` 类型的值有哪些优势
- 已存在的内容不可变，但可以拼接更多的内容；
- 减少了内存分配和内容拷贝的次数；
- 可将内容重置，可重用值。

`Builder` 值中有一个用于承载内容的容器（以下简称内容容器）。它是一个以 `byte` 为元素类型的切片（以下简称字节切片）。

**由于这样的字节切片的底层数组就是一个字节数组，所以我们可以说它与 `string` 值存储内容的方式是一样的**。实际上，它们都是通过
一个 `unsafe.Pointer` 类型的字段来持有那个指向了底层字节数组的指针值的。

因为这样的内部构造，`Builder` 值同样拥有高效利用内存的前提条件。

已存在于 `Builder` 值中的内容是不可变的。因此，我们可以利用 `Builder` 值提供的方法拼接更多的内容，而丝毫不用担心这些方法
会影响到已存在的内容。

`Builder` 值拥有的一系列指针方法，包括：`Write`、`WriteByte`、`WriteRune` 和 `WriteString`。我们可以把它们统称
为**拼接方法**。

调用上述方法把新的内容拼接到已存在的内容的尾部（也就是右边）。这时，如有必要，`Builder` 值会自动地对自身的内容容器进行扩容。
这里的自动扩容策略与切片的扩容策略一致。

除了 `Builder` 值的自动扩容，我们还可以选择手动扩容，这通过调用 `Builder` 值的 `Grow` 方法就可以做到。`Grow` 方法也可以被称
为**扩容方法**，它接受一个 `int` 类型的参数 `n`，该参数用于代表将要扩充的字节数量。

`Grow` 方法会把其所属值中内容容器的容量增加 `n` 个字节。更具体地讲，它会生成一个字节切片作为新的内容容器，该切片的容量会是原
容器容量的二倍再加上 `n`。之后，它会把原容器中的所有字节全部拷贝到新容器中。

#### 使用 `strings.Builder` 类型的约束
**只要调用了 `Builder` 值的拼接方法或扩容方法，就不能再以任何的方式对其所属值进行复制了**。否则，只要在任何副本上调用上述方
法就都会引发 panic。这里所说的复制方式，包括但不限于在函数间传递值、通过通道传递值、把值赋予变量等等。

正是由于已使用的 `Builder` 值不能再被复制，所以肯定不会出现多个 `Builder` 值中的内容容器（也就是那个字节切片）共用一个底层字
节数组的情况。这样也就避免了多个同源的` Builder` 值在拼接内容时可能产生的冲突问题。

**不过，虽然已使用的 `Builder` 值不能再被复制，但是它的指针值却可以。无论什么时候，我们都可以通过任何方式复制这样的指针值**。
注意，这样的指针值指向的都会是同一个 `Builder` 值。

### `strings.Reader` 类型
`strings.Reader` 类型是为了高效读取字符串而存在的。可以让我们很方便地读取一个字符串中的内容。在读取的过程中，`Reader` 值会
保存已读取的字节的计数（以下简称已读计数）。

**已读计数也代表着下一次读取的起始索引位置。`Reader` 值正是依靠这样一个计数，以及针对字符串值的切片表达式，从而实现快速读取**。
### bytes 包与字节串操作
`strings` 包和 `bytes` 包可以说是一对孪生兄弟，它们在 API 方面非常的相似。单从它们提供的函数的数量和功能上讲，差别微乎其微。
只不过，`strings`包主要面向的是 `Unicode` 字符和经过 `UTF-8` 编码的字符串，而 `bytes` 包面对的则主要是字节和字节切片。


#### `bytes.Buffer`
`bytes.Buffer` 类型的用途主要是作为字节序列的缓冲区。`bytes.Buffer` 是开箱即用的。`bytes.Buffer` 不但可以拼接、截断其中
的字节序列，以各种形式导出其中的内容，还可以顺序地读取其中的子序列。

在内部，`bytes.Buffer` 类型同样是使用字节切片作为内容容器的。并且，与 `strings.Reader` 类型类似，`bytes.Buffer` 有一个 `int` 
类型的字段，用于代表已读字节的计数，可以简称为**已读计数**。

**注意，与 `strings.Reader` 类型的 `Len` 方法一样，`bytes.Buffer` 的` Len` 方法返回的也是内容容器中未被读取部分的长度，
而不是其中已存内容的总长度（以下简称内容长度）。**

```go
// 示例1。
var buffer1 bytes.Buffer
contents := "Simple byte buffer for marshaling data."
fmt.Printf("Write contents %q ...\n", contents)
buffer1.WriteString(contents)
fmt.Printf("The length of buffer: %d\n", buffer1.Len()) // => 39
fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // => 64
fmt.Println()

// 示例2。
p1 := make([]byte, 7)
n, _ := buffer1.Read(p1)
fmt.Printf("%d bytes were read. (call Read)\n", n)
fmt.Printf("The length of buffer: %d\n", buffer1.Len()) // => 32
fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // => 64
```
上面的代码，示例一输出 39 和 64，但是示例二，从` buffer1` 中读取一部分内容，并用它们填满长度为7的字节切片 `p1` 之后，
`buffer1` 的 `Len` 方法返回的结果值变为了 32。因为我们并没有再向该缓冲区中写入任何内容，所以它的容量会保持不变，仍是 64。

> 对于处在零值状态的 `Buffer` 值来说，如果第一次扩容时的另需字节数不大于 64，那么该值就会基于一个预先定义好的、长度为 64 
的字节数组来创建内容容器。

由于 `strings.Reader` 还有一个 `Size` 方法可以给出内容长度的值，所以我们用内容长度减去未读部分的长度，就可以很方便地得
到它的已读计数。

然而，`bytes.Buffer` 类型却没有这样一个方法，它只有 `Cap` 方法。可是 `Cap` 方法提供的是内容容器的容量，也不是内容长度。

#### bytes.Buffer 的扩容策略
`Buffer` 值既可以被手动扩容，也可以进行自动扩容。并且，这两种扩容方式的策略是基本一致的。所以，除非我们完全确定后续内容所需
的字节数，否则让 `Buffer` 值自动去扩容就好了。

在扩容的时候，`Buffer` 值中相应的代码（以下简称扩容代码）会先判断内容容器的剩余容量，是否可以满足调用方的要求，或者是否足
够容纳新的内容。

如果可以，那么扩容代码会在当前的内容容器之上，进行长度扩充。更具体地说，如果内容容器的容量与其长度的差，大于或等于另需的字
节数，那么扩容代码就会通过切片操作对原有的内容容器的长度进行扩充，就像下面这样：
```go
b.buf = b.buf[:length+need]
```
反之，如果内容容器的剩余容量不够了，那么扩容代码可能就会用新的内容容器去替代原有的内容容器，从而实现扩容。不过，这里还一步优化。

如果当前内容容器的容量的一半仍然大于或等于其现有长度再加上另需的字节数的和，即：
```go
cap(b.buf)/2 >= len(b.buf)+need
```
那么，扩容代码就会复用现有的内容容器，并把容器中的未读内容拷贝到它的头部位置。这也意味着其中的已读内容，将会全部被未读内容和
之后的新内容覆盖掉。

这样的复用预计可以至少节省掉一次后续的扩容所带来的内存分配，以及若干字节的拷贝。

若这一步优化未能达成，也就是说，当前内容容器的容量小于新长度的二倍，那么扩容代码就只能再创建一个新的内容容器，并把原有容器
中的未读内容拷贝进去，最后再用新的容器替换掉原有的容器。这个新容器的容量将会等于原有容量的二倍再加上另需字节数的和。
```
新容器的容量 =2* 原有容量 + 所需字节数
```

#### bytes.Buffer 中的哪些方法可能会造成内容的泄露
什么叫内容泄露？这里所说的内容泄露是指，使用 `Buffer` 值的一方通过某种非标准的（或者说不正式的）方式得到了本不该得到的内容。

在` bytes.Buffer` 中，**`Bytes` 方法和` Next`方法都可能会造成内容的泄露**。原因在于，它们都把基于内容容器的切片直接返
回给了方法的调用方。

我们都知道，**通过切片，我们可以直接访问和操纵它的底层数组。不论这个切片是基于某个数组得来的，还是通过对另一个切片做切片操作
获得的**，都是如此。
```go
contents := "ab"
buffer1 := bytes.NewBufferString(contents)
fmt.Printf("The capacity of new buffer with contents %q: %d\n",
    contents, buffer1.Cap()) // 内容容器的容量为：8。
fmt.Println()

unreadBytes := buffer1.Bytes()
fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes)
```

前面通过调用 `buffer1` 的` Bytes` 方法得到的结果值 `unreadBytes`，包含了在那时其中的所有未读内容。

但是，由于这个结果值与 `buffer1` 的内容容器在此时还共用着同一个底层数组，所以，我只需通过简单的再切片操作，就可以利用这个
结果值拿到 `buffer1` 在此时的所有未读内容。如此一来，`buffer1` 的新内容就被泄露出来了。
### 一个 `string` 类型的值在底层怎样被表达
在底层，一个 `string` 类型的值是由一系列相对应的 Unicode 代码点的 UTF-8 编码值来表达的。

一个 `string` 类型的值既可以被拆分为一个包含多个字符的序列，也可以被拆分为一个包含多个字节的序列。
前者可以由一个以 `rune`（`int32` 的别名）为元素类型的切片来表示，而后者则可以由一个以 `byte` 为元素类型的切片代表。

`rune` 是 Go 语言特有的一个基本数据类型，它的一个值就代表一个字符，即：一个 Unicode 字符。比如，'G'、'o'、'爱'、'好'、
'者'代表的就都是一个 Unicode 字符。一个 `rune` 类型的值会由四个字节宽度的空间来存储。它的存储空间总是能够存下一
个 UTF-8 编码值。

**一个 `rune` 类型的值在底层其实就是一个 UTF-8 编码值**。前者是（便于我们人类理解的）外部展现，后者是（便于计算机系统理解的）
内在表达。

```go
str := "Go 爱好者 "
fmt.Printf("The string: %q\n", str)
fmt.Printf("  => runes(char): %q\n", []rune(str))
fmt.Printf("  => runes(hex): %x\n", []rune(str))
fmt.Printf("  => bytes(hex): [% x]\n", []byte(str))
```
字符串值 "Go 爱好者" 如果被转换为 `[]rune` 类型的值的话，其中的每一个字符（不论是英文字符还是中文字符）就都会独立成为一
个 `rune` 类型的元素值。因
此，这段代码打印出的第二行内容就会如下所示：
```bash
=> runes(char): ['G' 'o' '爱' '好' '者']
```
又由于，每个 `rune` 类型的值在底层都是由一个 UTF-8 编码值来表达的，所以我们可以换一种方式来展现这个字符序列：
```bash
=> runes(hex): [47 6f 7231 597d 8005]
```
我们还可以进一步地拆分，把每个字符的 UTF-8 编码值都拆成相应的字节序列。上述代码中的第五行就是这么做的。它会得到如下的输出：
```bash
=> bytes(hex): [47 6f e7 88 b1 e5 a5 bd e8 80 85]
```