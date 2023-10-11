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

以 `fmt.Fprintln` 为例：

```go
func Println(a ...interface{}) (n int, err error) {
	return Fprintln(os.Stdout, a...)
}
```

可以看出 `fmt.Println` 会将内容输出到标准输出中。


## 实现了 io.Reader 接口或 io.Writer 接口的类型

标准库中有哪些类型实现了 `io.Reader` 或 `io.Writer` 接口？

例如 `os.Stdin/Stdout`，它们分别实现了 `io.Reader/io.Writer` 接口：

```go
var (
    Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
    Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
    Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)
```

上面的代码可以看出，`Stdin/Stdout/Stderr` 只是三个特殊的文件类型的标识（都是 `os.File` 的实例），`os.File` 实现
了 `io.Reader` 和 `io.Writer`。

实现了 `io.Reader` 或 `io.Writer` 接口的类型：

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

**io 包本身实现这两个接口的类型**：

- 实现了 `Reader` 的类型：`LimitedReader`、`PipeReader`、`SectionReader`
- 实现了 `Writer` 的类型：`PipeWriter`


## ReaderAt 和 WriterAt

**`ReaderAt` 接口**：

```go
type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
}
```

> `ReadAt` 从基本输入源的偏移量 `off` 处开始，将 `len(p)` 个字节读到 `p` 中。它返回读取的字节数 `n`（`0 <= n <= len(p)`）以及任
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

将文件中的数据全部读取（显示在标准输出）：

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

> `Seek` 设置下一次 `Read` 或 `Write` 的偏移量为 `offset`，它的解释取决于 `whence`：  **0 表示相对于文件的起始处，1 表示相对
于当前的偏移，而 2 表示相对于其结尾处**。 `Seek` 返回新的偏移量和一个错误，如果有的话。

也就是说，`Seek` 方法是用于设置偏移量的，这样可以从某个特定位置开始操作数据流。听起来和 `ReaderAt/WriteAt` 接口有些类似，
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

## 其他接口

### ByteReader 和 ByteWriter

读或写一个字节：

```go
type ByteReader interface {
    ReadByte() (c byte, err error)
}

type ByteWriter interface {
    WriteByte(c byte) error
}
```

下面类型都实现了这两个接口:

- `bufio.Reader/Writer` 分别实现了 `io.ByteReader` 和 `io.ByteWriter`
- `bytes.Buffer` 同时实现了 `io.ByteReader` 和 `io.ByteWriter`
- `bytes.Reader` 实现了 `io.ByteReader`
- `strings.Reader` 实现了 `io.ByteReader`

通过 `bytes.Buffer` 来一次读取或写入一个字节：

```go
var ch byte
fmt.Scanf("%c\n", &ch)

buffer := new(bytes.Buffer)
err := buffer.WriteByte(ch)
if err == nil {
	fmt.Println("写入一个字节成功！准备读取该字节……")
	newCh, _ := buffer.ReadByte()
	fmt.Printf("读取的字节：%c\n", newCh)
} else {
	fmt.Println("写入错误")
}
```

程序从标准输入接收一个字节（ASCII 字符），调用 `buffer` 的 `WriteByte` 将该字节写入 `buffer` 中，之后通过 `ReadByte` 读取该字节。

### ByteScanner、RuneReader 和 RuneScanner

`ByteScanner` 接口：

```go
type ByteScanner interface {
    ByteReader
    UnreadByte() error
}
```

内嵌了 `ByteReader` 接口，`UnreadByte` 方法的意思是：将上一次 `ReadByte` 的字节还原，使得再次调用 `ReadByte` 返回的结果和上一次调
用相同，也就是说，`UnreadByte` 是重置上一次的 `ReadByte`。注意，**`UnreadByte` 调用之前必须调用了 `ReadByte`，且不能连续调
用 `UnreadByte`**。即：

```go
buffer := bytes.NewBuffer([]byte{'a', 'b'})
err := buffer.UnreadByte()
```

和

```go
buffer := bytes.NewBuffer([]byte{'a', 'b'})
buffer.ReadByte()
err := buffer.UnreadByte()
err = buffer.UnreadByte()
```

`err` 都 **非 nil**，错误为：`bytes.Buffer: UnreadByte: previous operation was not a read`

`RuneReader` 接口和 `ByteReader` 类似，只是 `ReadRune` 方法读取单个 UTF-8 字符，返回其 `rune` 和该字符占用的字节数。

`RuneScanner` 接口和 `ByteScanner` 类似。

### ReadCloser、ReadSeeker、ReadWriteCloser、ReadWriteSeeker、ReadWriter、WriteCloser 和 WriteSeeker

这些接口是上面介绍的接口的两个或三个组合而成的新接口。`ReadWriter` 接口：

```go
type ReadWriter interface {
    Reader
    Writer
}
```

`Reader` 和 `Writer` 接口的组合。

这些接口的作用是：有些时候同时需要某两个接口的所有功能，即必须同时实现了某两个接口的类型才能够被传入使用。

## SectionReader 类型

`SectionReader` 是一个 `struct`，实现了 `Read`, `Seek` 和 `ReadAt`，同时，内嵌了 `ReaderAt` 接口。结构定义如下：

```go
type SectionReader struct {
	r     ReaderAt	// 该类型最终的 Read/ReadAt 最终都是通过 r 的 ReadAt 实现
	base  int64		// NewSectionReader 会将 base 设置为 off
	off   int64		// 从 r 中的 off 偏移处开始读取数据
	limit int64		// limit - off = SectionReader 流的长度
}
```

该类型读取数据流中部分数据。

```go
func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader
```


> `NewSectionReader` 返回一个 `SectionReader`，它从 `r` 中的偏移量 `off` 处读取 `n` 个字节后以 `EOF` 停止。

也就是说，SectionReader 只是内部 ReaderAt 表示的数据流的一部分：从 off 开始后的 n 个字节。

这个类型的作用是：方便重复操作某一段 (section) 数据流；或者同时需要 ReadAt 和 Seek 的功能。

由于该类型所支持的操作，前面都有介绍，因此不提供示例代码了。

## LimitedReader 类型

```go
type LimitedReader struct {
    R Reader // underlying reader，最终的读取操作通过 R.Read 完成
    N int64  // max bytes remaining
}
```


> 从 `R` 读取但将返回的数据量限制为 `N` 字节。每调用一次 `Read` 都将更新 `N` 来反应新的剩余数量。

也就是说，最多只能返回 `N` 字节数据。

`LimitedReader` 只实现了 `Read` 方法。

示例：

```go
content := "This Is LimitReader Example"
reader := strings.NewReader(content)
limitReader := &io.LimitedReader{R: reader, N: 8}
for limitReader.N > 0 {
	tmp := make([]byte, 2)
	limitReader.Read(tmp)
	fmt.Printf("%s", tmp) // This Is
}
```

通过该类型可以达到 **只允许读取一定长度数据** 的目的。

在 `io` 包中，`LimitReader` 函数的实现其实就是调用 `LimitedReader`：

```go
func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }
```

## PipeReader 和 PipeWriter 类型

`PipeReader` 是管道的读取端。它实现了 `io.Reader` 和 `io.Closer` 接口：
```go
type PipeReader struct {
	p *pipe
}
```

> 从管道中读取数据。该方法会堵塞，直到管道写入端开始写入数据或写入端被关闭。如果写入端关闭时带有 `error`（即调用 `CloseWithError` 关闭），
该 `Read` 返回的 `err` 就是写入端传递的 `error`；否则 `err` 为 `EOF`。

`PipeWriter` 是管道的写入端。它实现了 `io.Writer` 和 `io.Closer` 接口：
```go
type PipeWriter struct {
	p *pipe
}
```

> 写数据到管道中。该方法会堵塞，直到管道读取端读完所有数据或读取端被关闭。如果读取端关闭时
带有 `error`（即调用 `CloseWithError` 关闭），该 `Write` 返回的 `err` 就是读取端传递的 `error`；否则 `err` 为 `ErrClosedPipe`。

示例：

```go
func main() {
    pipeReader, pipeWriter := io.Pipe()
    go PipeWrite(pipeWriter)
    go PipeRead(pipeReader)
    time.Sleep(30 * time.Second)
}

func PipeWrite(writer *io.PipeWriter){
	data := []byte("Go语言中文网")
	for i := 0; i < 3; i++{
		n, err := writer.Write(data)
		if err != nil{
			fmt.Println(err)
			return
		}
		fmt.Printf("写入字节 %d\n",n)
	}
	writer.CloseWithError(errors.New("写入段已关闭"))
}

func PipeRead(reader *io.PipeReader){
	buf := make([]byte, 128)
	for{
		fmt.Println("接口端开始阻塞5秒钟...")
		time.Sleep(5 * time.Second)
		fmt.Println("接收端开始接受")
		n, err := reader.Read(buf)
		if err != nil{
			fmt.Println(err)
			return
		}
		fmt.Printf("收到字节: %d\n buf内容: %s\n",n,buf)
	}
}
```

`io.Pipe()` 用于创建一个同步的内存管道：

```go
func Pipe() (*PipeReader, *PipeWriter)
```

它将 `io.Reader` 连接到 `io.Writer`。一端的读取匹配另一端的写入，直接在这两端之间复制数据；它没有内部缓存。它对于并行调用 `Read`
 和 `Write` 以及其它函数或 `Close` 来说都是安全的。一旦等待的 I/O 结束，`Close` 就会完成。并行调用 `Read` 或并行调用 `Write` 也
同样安全：同种类的调用将按顺序进行控制。

正因为是**同步**的，因此不能在一个 goroutine 中进行读和写。

另外，对于管道的 `close` 方法（非 `CloseWithError` 时），`err` 会被置为 `EOF`。

## Copy 和 CopyN 函数

**Copy 函数**：

```go
func Copy(dst Writer, src Reader) (written int64, err error)
```

函数文档：

> `Copy` 将 `src` 复制到 `dst`，直到在 `src` 上到达 `EOF` 或发生错误。它返回复制的字节数，如果有错误的话，还会返回在复制时遇到的第一
个错误。

> 成功的 `Copy` 返回 `err == nil`，而非 `err == EOF`。由于 `Copy` 被定义为从 `src` 读取直到 `EOF` 为止，因此它不会将来
自 `Read` 的 `EOF` 当做错误来报告。

> 若 `dst` 实现了 `ReaderFrom` 接口，其复制操作可通过调用 `dst.ReadFrom(src)` 实现。此外，若 `src` 实现了 `WriterTo` 接口，其复
制操作可通过调用 `src.WriteTo(dst)` 实现。

代码：

```go
io.Copy(os.Stdout, strings.NewReader("Go语言中文网"))
```

直接将内容输出（写入 Stdout 中）：

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	_, _ = io.Copy(os.Stdout, os.Stdin)
	fmt.Println("Got EOF -- bye")
}
```

执行：`echo "Hello, World" | go run main.go`


**CopyN 函数**：

```go
func CopyN(dst Writer, src Reader, n int64) (written int64, err error)
```

> `CopyN` 将 `n` 个字节(或到一个 `error`)从 `src` 复制到 `dst`。 它返回复制的字节数以及在复制时遇到的最早的错误。
当且仅当 `err == nil` 时,`written == n` 。

> 若 `dst` 实现了 `ReaderFrom` 接口，复制操作也就会使用它来实现。

代码：

```go
io.CopyN(os.Stdout, strings.NewReader("Go语言中文网"), 8) // Go语言
```

## ReadAtLeast 和 ReadFull 函数

**ReadAtLeast 函数**：

```go
func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error)
```

函数文档：

> `ReadAtLeast` 将 `r` 读取到 `buf` 中，直到读了最少 `min` 个字节为止。它返回复制的字节数，如果读取的字节较少，还会返回一个错误。
若没有读取到字节，错误就只是 `EOF`。如果一个 `EOF` 发生在读取了少于 `min` 个字节之后，`ReadAtLeast` 就会返回 `ErrUnexpectedEOF`。
若 `min` 大于 `buf` 的长度，`ReadAtLeast` 就会返回 `ErrShortBuffer`。对于返回值，当且仅当 `err == nil` 时，才有 `n >= min`。

**ReadFull 函数**：

```go
func ReadFull(r Reader, buf []byte) (n int, err error)
```

函数文档：

> `ReadFull` 精确地从 `r` 中将 `len(buf)` 个字节读取到 `buf` 中。它返回复制的字节数，如果读取的字节较少，还会返回一个错误。若没有读
取到字节，错误就只是 `EOF`。如果一个 `EOF` 发生在读取了一些但不是所有的字节后，`ReadFull` 就会返回 `ErrUnexpectedEOF`。对于返回值，
当且仅当 `err == nil` 时，才有 `n == len(buf)`。

注意该函数和 `ReadAtLeast` 的区别：
- `ReadFull` 将 `buf` 读满
- `ReadAtLeast` 是最少读取 `min` 个字节。


## WriteString 函数

这是为了方便写入 `string` 类型提供的函数，函数签名：

```go
func WriteString(w Writer, s string) (n int, err error)
```

函数文档：

> `WriteString` 将 ``s 的内容写入 `w` 中，当 `w` 实现了 `WriteString` 方法时，会直接调用该方法，否则执行 `w.Write([]byte(s))`。

## MultiReader 和 MultiWriter 函数

```go
func MultiReader(readers ...Reader) Reader
func MultiWriter(writers ...Writer) Writer
```

它们接收多个 `Reader` 或 `Writer`，返回一个 `Reader` 或 `Writer`。我们可以猜想到这两个函数就是操作多个 `Reader` 或 `Writer` 就像
操作一个。

事实上，在 `io` 包中定义了两个非导出类型：`mutilReader` 和 `multiWriter`，它们分别实现了 `io.Reader` 和 `io.Writer` 接口：

```go
type multiReader struct {
	readers []Reader
}

type multiWriter struct {
	writers []Writer
}
```

对于这两种类型对应的实现方法（`Read` 和 `Write` 方法）的使用，示例：

**MultiReader 的使用**：

```go
readers := []io.Reader{
	strings.NewReader("from strings reader"),
	bytes.NewBufferString("from bytes buffer"),
}
reader := io.MultiReader(readers...)
data := make([]byte, 0, 128)
buf := make([]byte, 10)
	
for n, err := reader.Read(buf); err != io.EOF ; n, err = reader.Read(buf){
	if err != nil{
		panic(err)
	}
	data = append(data,buf[:n]...)
}
fmt.Printf("%s\n", data) // from strings readerfrom bytes buffer
```

代码中首先构造了一个 `io.Reader` 的 `slice`，然后通过 `MultiReader` 得到新的 `Reader`，循环读取新 `Reader` 中的内容。从输出结果
可以看到，第一次调用 `Reader` 的 `Read` 方法获取到的是 `slice` 中第一个元素的内容……也就是说，`MultiReader` 只是逻辑上将多
个 `Reader` 组合起来，并不能通过调用一次 `Read` 方法获取所有 `Reader` 的内容。在所有的 `Reader` 内容都被读
完后，`Reader` 会返回 `EOF`。

**MultiWriter 的使用**：

```go
file, err := os.Create("tmp.txt")
if err != nil {
    panic(err)
}
defer file.Close()
writers := []io.Writer{
	file,
	os.Stdout,
}
writer := io.MultiWriter(writers...)
writer.Write([]byte("Go语言中文网"))
```

这段程序执行后在生成 `tmp.txt` 文件，同时在文件和屏幕中都输出：`Go语言中文网`。这和 Unix 中的 tee 命令类似。



## TeeReader 函数


```go
func TeeReader(r Reader, w Writer) Reader
```

> `TeeReader` 返回一个 `Reader`，它将从 `r` 中读到的数据写入 `w` 中。所有经由它处理的从 `r` 的读取都匹配于对应的对 `w` 的写入。它没
有内部缓存，即写入必须在读取完成前完成。任何在写入时遇到的错误都将作为读取错误返回。

也就是说，我们通过 `Reader` 读取内容后，会自动写入到 `Writer` 中去。例子代码如下：

```go
reader := io.TeeReader(strings.NewReader("Go语言中文网"), os.Stdout)
reader.Read(make([]byte, 20)) // Go语言中文网
```


这种功能的实现其实挺简单，无非是在 `Read` 完后执行 `Write`。