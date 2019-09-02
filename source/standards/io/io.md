---
title: io
---

# io
`io` 是对输入输出设备的抽象。`io` 库对这些功能进行了抽象，通过统一的接口对输入输出设备进行操作。
最重要的是两个接口：`Reader` 和 `Writer`。

## Reader

Reader 接口：

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

> `Read` 将 `len(p)` 个字节读取到 `p` 中。它返回读取的字节数 `n`（`0 <= n <= len(p)`） 以及任何遇到的错误。
即使 `Read` 返回的 `n < len(p)`，它也会在调用过程中占用 `len(p)` 个字节作为暂存空间。若可读取的数据不到 `len(p)` 个
字节，`Read` 会返回可用数据，而不是等待更多数据。

> 当 `Read` 在成功读取 `n > 0` 个字节后遇到一个错误或 `EOF` (`end-of-file`)，它会返回读取的字节数。它可能会同时在本次的调
用中返回一个 `non-nil` 错误,或在下一次的调用中返回这个错误（且 `n` 为 0）。 一般情况下, `Reader` 会返回一个 非 0 字节数 `n`, 
若 `n = len(p)` 个字节从输入源的结尾处由 `Read` 返回，`Read` 可能返回 `err == EOF` 或者 `err == nil`。并且之后的 `Read` 
都应该返回 (`n:0, err:EOF`)。

> 调用者在考虑错误之前应当首先处理返回的数据。这样做可以正确地处理在读取一些字节后产生的 I/O 错误，同时允许 `EOF` 的出现。

```go
func ReadFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}
```

`ReadFrom` 函数将 `io.Reader` 作为参数，也就是说，`ReadFrom` 可以从任意的地方读取数据，只要来源实现了 `io.Reader` 接口。
比如，我们可以从标准输入、文件、字符串等读取数据，示例代码如下：

```go
// 从标准输入读取
data, err = ReadFrom(os.Stdin, 11)

// 从普通文件读取，其中 file 是 os.File 的实例
data, err = ReadFrom(file, 9)

// 从字符串读取
data, err = ReadFrom(strings.NewReader("from string"), 12)
```

`io.EOF` 变量的定义：`var EOF = errors.New("EOF")`，是 `error` 类型。根据 `reader` 接口的说明，在 `n > 0` 且数据被读完了
的情况下，当次返回的 `error` 有可能是 `EOF` 也有可能是 `nil`。

## Writer

Writer 接口：

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

> `Write` 将 `len(p)` 个字节从 `p` 中写入到基本数据流中。它返回从 `p` 中被写入的字节数 `n`（`0 <= n <= len(p)`）以及任何遇到的引
起写入提前停止的错误。若 `Write` 返回的 `n < len(p)`，它就必须返回一个 **非 nil** 的错误。

所有实现了 `Write` 方法的类型都实现了 `io.Writer` 接口。

以 `fmt.Fprintln` 为例，`fmt.Println` 函数的源码。

```go
func Println(a ...interface{}) (n int, err error) {
	return Fprintln(os.Stdout, a...)
}
```

`fmt.Println` 会将内容输出到标准输出中。


## 实现了 io.Reader 接口或 io.Writer 接口的类型

标准库中有哪些类型实现了 `io.Reader` 或 `io.Writer` 接口？

`os.File` 同时实现了这两个接口。我们还看到 `os.Stdin/Stdout` 这样的代码，它们分别实现了 `io.Reader/io.Writer` 接口：

```go
var (
    Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
    Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
    Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)
```

也就是说，`Stdin/Stdout/Stderr` 只是三个特殊的文件类型的标识（即都是 `os.File` 的实例），自然也实现了 `io.Reader` 和 `io.Writer`。

列出实现了 io.Reader 或 io.Writer 接口的类型（导出的类型）：

- `os.File` 同时实现了 `io.Reader` 和 `io.Writer`
- `strings.Reader` 实现了 `io.Reader`
- `bufio.Reader/Writer` 分别实现了 `io.Reader` 和 `io.Writer`
- `bytes.Buffer` 同时实现了 `io.Reader` 和 `io.Writer`
- `bytes.Reader` 实现了 `io.Reader`
- `compress/gzip.Reader/Writer` 分别实现了 `io.Reader` 和 `io.Writer`
- `crypto/cipher.StreamReader/StreamWriter` 分别实现了 `io.Reader` 和 `io.Writer`
- `crypto/tls.Conn` 同时实现了 `io.Reader` 和 `io.Writer`
- `encoding/csv.Reader/Writer` 分别实现了 `io.Reader` 和 `io.Writer`
- `mime/multipart.Part` 实现了 `io.Reader`
- `net/conn` 分别实现了 `io.Reader` 和 `io.Writer`(Conn接口定义了Read/Write)

除此之外，io 包本身也有这两个接口的实现类型。如：

	实现了 Reader 的类型：LimitedReader、PipeReader、SectionReader
	实现了 Writer 的类型：PipeWriter

以上类型中，常用的类型有：os.File、strings.Reader、bufio.Reader/Writer、bytes.Buffer、bytes.Reader

## ReaderAt 和 WriterAt

**`ReaderAt` 接口**：

```go
type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
}
```

> `ReadAt` 从基本输入源的偏移量 `off` 处开始，将 `len(p)` 个字节读取到 `p` 中。它返回读取的字节数 `n`（`0 <= n <= len(p)`）以及任
何遇到的错误。

> 当 `ReadAt` 返回的 `n < len(p)` 时，它就会返回一个 **非 nil** 的错误来解释为什么没有返回更多的字节。

> 即使 `ReadAt` 返回的 `n < len(p)`，它也会在调用过程中使用 `p` 的全部作为暂存空间。若可读取的数据不到 `len(p)` 字节，`ReadAt` 就会
阻塞,直到所有数据都可用或一个错误发生。

> 若 `n = len(p)` 个字节从输入源的结尾处由 `ReadAt` 返回，`Read` 可能返回 `err == EOF` 或者 `err == nil`

> 若 `ReadAt` 携带一个偏移量从输入源读取，`ReadAt` 应当既不影响偏移量也不被它所影响。

> 可对相同的输入源并行执行 `ReadAt` 调用。

可见，`ReadAt` 接口使得可以从指定偏移量处开始读取数据。

简单示例代码如下：

```go
reader := strings.NewReader("Hello world")
p := make([]byte, 6)
n, err := reader.ReadAt(p, 2)
if err != nil {
    panic(err)
}
fmt.Printf("%s, %d\n", p, n) // llo wo, 6
```

**`WriterAt` 接口**：

```go
type WriterAt interface {
    WriteAt(p []byte, off int64) (n int, err error)
}
```


> `WriteAt` 从 `p` 中将 `len(p)` 个字节写入到偏移量 `off` 处的基本数据流中。它返回从 `p` 中被写入的字节数 `n`（`0 <= n <= len(p)`）
以及任何遇到的引起写入提前停止的错误。若 `WriteAt` 返回的 `n < len(p)`，它就必须返回一个 **非 nil** 的错误。

> 若 `WriteAt` 携带一个偏移量写入到目标中，`WriteAt` 应当既不影响偏移量也不被它所影响。

> 若被写区域没有重叠，可对相同的目标并行执行 `WriteAt` 调用。

我们可以通过该接口将数据写入到数据流的特定偏移量之后。

```go
file, err := os.Create("writeAt.txt")
if err != nil {
    panic(err)
}
defer file.Close()
_, _ = file.WriteString("Hello world----ignore")
n, err := file.WriteAt([]byte("Golang"), 15)
if err != nil {
    panic(err)
}
fmt.Println(n)
```

打开文件 `WriteAt.txt`，内容是：`Hello world----Golang`。

分析：

`file.WriteString("Hello world----ignore")` 往文件中写入 `Hello world----ignore`，之后 
`file.WriteAt([]byte("Golang"), 15)` 在文件流的 `offset=15` 处写入 `Golang`（会覆盖该位置的内容）。

## ReaderFrom 和 WriterTo
这两个接口实现了**一次性从某个地方读或写到某个地方去**。
**ReaderFrom**：

```go
type ReaderFrom interface {
    ReadFrom(r Reader) (n int64, err error)
}
```

> `ReadFrom` 从 `r` 中读取数据，直到 `EOF` 或发生错误。其返回值 `n` 为读取的字节数。除 `io.EOF` 之外，在读取过程中遇到的任何错误也
将被返回。

> 如果 `ReaderFrom` 可用，`Copy` 函数就会使用它。

注意：`ReadFrom` 方法不会返回 `err == EOF`。

下面的例子简单的实现将文件中的数据全部读取（显示在标准输出）：

```go
file, err := os.Open("writeAt.txt")
if err != nil {
    panic(err)
}
defer file.Close()
writer := bufio.NewWriter(os.Stdout)
writer.ReadFrom(file)
writer.Flush()
```

也可以通过 `ioutil` 包的 `ReadFile` 函数获取文件全部内容。其实，跟踪一下 `ioutil.ReadFile` 的源码，会发现其实也是通过 `ReadFrom` 方
法实现（用的是 `bytes.Buffer`，它实现了 `ReaderFrom` 接口）。

**WriterTo**：

```go
type WriterTo interface {
    WriteTo(w Writer) (n int64, err error)
}
```

> `WriteTo` 将数据写入 `w` 中，直到没有数据可写或发生错误。其返回值 `n` 为写入的字节数。 在写入过程中遇到的任何错误也将被返回。

> 如果 `WriterTo` 可用，`Copy` 函数就会使用它。

将一段文本输出到标准输出：

```go
reader := bytes.NewReader([]byte("Hello world"))
reader.WriteTo(os.Stdout)
```

## Seeker

```go
type Seeker interface {
    Seek(offset int64, whence int) (ret int64, err error)
}
```

> `Seek` 设置下一次 `Read` 或 `Write` 的偏移量为 `offset`，它的解释取决于 `whence`：  0 表示相对于文件的起始处，1 表示相对
于当前的偏移，而 2 表示相对于其结尾处。 `Seek` 返回新的偏移量和一个错误，如果有的话。

也就是说，`Seek` 方法是用于设置偏移量的，这样可以从某个特定位置开始操作数据流。听起来和 `ReaderAt/WriteA`t 接口有些类似，
不过 `Seeker` 接口更灵活，可以更好的控制读写数据流的位置。

获取倒数第二个字符（需要考虑 UTF-8 编码，这里的代码只是一个示例）：

```go
reader := strings.NewReader("Hello world")
reader.Seek(-6, io.SeekEnd)
r, _, _ := reader.ReadRune()
fmt.Printf("%c\n", r)
```

`whence` 的值，在 io 包中定义了相应的常量，应该使用这些常量

```go
const (
  SeekStart   = 0 // seek relative to the origin of the file
  SeekCurrent = 1 // seek relative to the current offset
  SeekEnd     = 2 // seek relative to the end
)
```

而原先 `os` 包中的常量已经被标注为 Deprecated

```go
// Deprecated: Use io.SeekStart, io.SeekCurrent, and io.SeekEnd.
const (
  SEEK_SET int = 0 // seek relative to the origin of the file
  SEEK_CUR int = 1 // seek relative to the current offset
  SEEK_END int = 2 // seek relative to the end
)
```

## Closer

```go
type Closer interface {
    Close() error
}
```

该接口比较简单，只有一个 `Close()` 方法，用于关闭数据流。

文件 (`os.File`)、归档（压缩包）、数据库连接、`Socket` 等需要手动关闭的资源都实现了 `Closer` 接口。

实际编程中，经常将 `Close` 方法的调用放在 `defer` 语句中。

```go
file, err := os.Open("studygolang.txt")
defer file.Close()
if err != nil {
	...
}
```

当文件 `studygolang.txt` 不存在或找不到时，`file.Close()` 会返回错误，因为 `file` 是 `nil`。
因此，应该**将 `defer file.Close()` 放在错误检查之后**。

```go
func (f *File) Close() error {
	if f == nil {
		return ErrInvalid
	}
	return f.file.close()
}
```

## io 包中的接口和工具
`strings.Reader` 类型主要用于读取字符串，它的指针类型实现的接口比较多，包括：
- io.Reader；
- io.ReaderAt；
- io.ByteReader；
- io.RuneReader；
- io.Seeker；
- io.ByteScanner；
- io.RuneScanner；
- io.WriterTo；

`io.ByteScanner` 是`io.ByteReader`的扩展接口，而`io.RuneScanner`又是`io.RuneReader`的扩展接口。

`bytes.Buffer`该指针类型实现的读取相关的接口有下面几个：
- io.Reader；
- io.ByteReader；
- io.RuneReader；
- io.ByteScanner；
- io.RuneScanner；
- io.WriterTo；

实现的写入相关的接口：
- io.Writer；
- io.ByteWriter；
- io.stringWriter；
- io.ReaderFrom；

这些类型实现了这么多的接口，目的是什么？

为了提高不同程序实体之间的互操作性。以 io 包中的一些函数为例。

io 包中，有这样几个用于拷贝数据的函数，它们是：`io.Copy`、`io.CopyBuffer`和`io.CopyN`。这几个函数在功能上都略有差别，但是它们都首先会接受两个参数，即：
用于代表**数据目的地、`io.Writer`类型的参数`dst`**，以及用于代表**数据来源的、`io.Reader`类型的参数`src`**。大致上都是把数据从`src`拷贝到`dst`。

**不论第一个参数值是什么类型的，只要这个类型实现了`io.Writer`接口即可**。同样的第二个参数值只要该类型实现了`io.Reader`接口就行。

很多数据类型实现了`io.Reader`接口，是因为它们提供了从某处读取数据的功能。类似的，许多能够把数据写入某处的数据类型，也都会去实现`io.Writer`接口。

### io.Reader的扩展接口和实现类型
`io.Reader`的扩展接口：
- `io.ReadWriter`：此接口既是`io.Reader`的扩展接口，也是`io.Writer`的扩展接口。
- `io.ReadCloser`：此接口除了包含基本的字节序列读取方法之外，还拥有一个基本的关闭方法`Close`。后者一般用于关闭数据读写的通路。这个接口其实是`io.Reader`接口和`io.Closer`接口的组合。
- `io.ReadWriteCloser`：`io.Reader`、`io.Writer`和`io.Closer`这三个接口的组合。
- `io.ReadSeeker`：此接口的特点是拥有一个用于寻找读写位置的基本方法`Seek`。更具体地说，该方法可以根据给定的偏移量基于数据的起始位置、末尾位置，或者当前读写
位置去寻找新的读写位置。这个新的读写位置用于表明下一次读或写时的起始索引。`Seek`是`io.Seeker`接口唯一拥有的方法。
- `io.ReadWriteSeeker`：`io.Reader`、`io.Writer`和`io.Seeker`的组合。

`io.Reader`接口的实现类型：
- `*io.LimitedReader`：此类型的基本类型会包装`io.Reader`类型的值，并提供一个额外的受限读取的功能。。
- `*io.SectionReader`：此类型的基本类型可以包装`io.ReaderAt`类型的值，并且会限制它的`Read`方法，只能够读取原始数据中的某一个部分（或者说某一段）。
- `*io.teeReader`：此类型是一个包级私有的数据类型，也是io.TeeReader函数结果值的实际类型。这个函数接受两个参数r和w，类型分别是`io.Reader`和`io.Writer`。
- `io.multiReader`：此类型也是一个包级私有的数据类型。类似的，io包中有一个名为`MultiReader`的函数，它可以接受若干个`io.Reader`类型的参数值，并返回一个实
际类型为`io.multiReader`的结果值。
- `io.pipe`：此类型为一个包级私有的数据类型，它比上述类型都要复杂得多。它不但实现了`io.Reader`接口，而且还实现了`io.Writer`接口。
实际上，`io.PipeReader`类型和`io.PipeWriter`类型拥有的所有指针方法都是以它为基础的。这些方法都只是代理了`io.pipe`类型值所拥有的某一个方法而已。
又因为`io.Pipe`函数会返回这两个类型的指针值并分别把它们作为其生成的同步内存管道的两端，所以可以说，`*io.pipe`类型就是io包提供的同步内存管道的核心实现。
- `io.PipeReader`：此类型可以被视为`io.pipe`类型的代理类型。