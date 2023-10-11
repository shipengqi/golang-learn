---
title: bufio
---

# bufio

bufio 包实现了缓存 IO。提供了数据缓冲功能，能够一定程度减少大块数据读写带来的开销。封装了 `io.Reader` 和 `io.Writer` 对象。


## bufio包中的数据类型
bufio包中的数据类型主要有：
- `Reader`；
- `Scanner`；
- `Writer` 和 `ReadWriter`。

## bufio.Reader
两个用于初始化 `bufio.Reader` 的函数：

- `NewReader` 函数初始化的 `Reader` 值会拥有一个默认尺寸的缓冲区。这个默认尺寸是 `4096` 个字节，即：`4 KB`。
- `NewReaderSize` 函数则将缓冲区尺寸的决定权抛给了使用方。

```go
func NewReader(rd io.Reader) *Reader

func NewReaderSize(rd io.Reader, size int) *Reader // 可以配置缓冲区的大小
```
### bufio.Reader 类型值中的缓冲区的作用
缓冲区其实就是一个**数据存储中介，它介于底层读取器与读取方法及其调用方之间**。所谓的底层读取器是指 `io.Reader`。

`Reader` 值的读取方法一般都会先从其所属值的缓冲区中读取数据。同时，在必要的时候，它们还会预先从底层读取器那里读出一部分数据，并暂
存于缓冲区之中以备后用。

缓冲区的好处是，可以在大多数的时候降低读取方法的执行时间。

```go
type Reader struct {
    buf          []byte
    rd           io.Reader
    r, w         int
    err          error
    lastByte     int
    lastRuneSize int
}
```

`bufio.Reader` 字段：
- `buf`：`[]byte` 类型的字段，即字节切片，代表缓冲区。虽然它是切片类型的，但是其长度却会在初始化的时候指定，并在之后保持不变。
- `rd`：`io.Reader` 类型的字段，代表底层读取器。缓冲区中的数据就是从这里拷贝来的。
- `r`：`int` 类型的字段，代表对缓冲区进行下一次读取时的开始索引。我们可以称它为已读计数。
- `w`：`int` 类型的字段，代表对缓冲区进行下一次写入时的开始索引。我们可以称之为已写计数。
- `err`：`error` 类型的字段。它的值用于表示在从底层读取器获得数据时发生的错误。这里的值在被读取或忽略之后，该字段会被置为 `nil`。
- `lastByte`：`int` 类型的字段，用于记录缓冲区中最后一个被读取的字节。读回退时会用到它的值。
- `lastRuneSize`：`int` 类型的字段，用于记录缓冲区中最后一个被读取的 Unicode 字符所占用的字节数。读回退的时候会用到它的值。这个字
段只会在其所属值的 `ReadRune` 方法中才会被赋予有意义的值。在其他情况下，它都会被置为 `-1`。

### bufio.Reader 类型读取方法

#### ReadSlice、ReadBytes、ReadString 和 ReadLine

后三个方法最终都是调用 `ReadSlice` 来实现的。所以，我们先来看看 `ReadSlice` 方法。

**ReadSlice方法**：
```go
func (b *Reader) ReadSlice(delim byte) (line []byte, err error)
```
`ReadSlice` 从输入中读取，直到遇到第一个界定符（delim）为止，返回一个指向缓存中字节的 `slice`，在下次调用读操作（`read`）时，这些字节会
无效：
```go
reader := bufio.NewReader(strings.NewReader("Hello \nworld"))
line, _ := reader.ReadSlice('\n')
fmt.Printf("the line:%s\n", line) // the line:Hello
n, _ := reader.ReadSlice('\n')
fmt.Printf("the line:%s\n", line) // the line:world
fmt.Println(string(n)) // world
```

从结果可以看出，第一次 `ReadSlice` 的结果 **line**，在第二次调用读操作后，内容发生了变化。也就是说，`ReadSlice` 返回的 `[]byte` 是指
向 `Reader` 中的 `buffer` ，而不是 `copy` 一份返回。正因为 `ReadSlice` 返回的数据会被下次的 I/O 操作重写，因此许多的客户端会选择
使用 `ReadBytes` 或者 `ReadString` 来代替。

注意，这里的界定符可以是任意的字符。同时，返回的结果是包含界定符本身的。

如果 `ReadSlice` 在找到界定符之前遇到了 `error`，它就会返回缓存中所有的数据和错误本身（经常是 `io.EOF`）。如果在找到界定符之前缓存已经
满了，`ReadSlice` 会返回 `bufio.ErrBufferFull` 错误。当且仅当返回的结果（`line`）没有以界定符结束的时候，`ReadSlice` 返
回 `err != nil`，也就是说，如果 `ReadSlice` 返回的结果 `line` 不是以界定符 `delim` 结尾，那么返回的 `err` 也一定不等于 `nil`。

**ReadBytes 方法**：
```go
func (b *Reader) ReadBytes(delim byte) (line []byte, err error)
```
该方法的参数和返回值类型与 `ReadSlice` 都一样。 `ReadBytes` 从输入中读取直到遇到界定符（delim）为止，返回的 `slice` 包含了从当前到
界定符的内容 **（包括界定符）**。

`ReadBytes` 源码：
```go
func (b *Reader) ReadBytes(delim byte) ([]byte, error) {
	// Use ReadSlice to look for array,
	// accumulating full buffers.
	var frag []byte
	var full [][]byte
	var err error
	for {
		var e error
		frag, e = b.ReadSlice(delim)
		if e == nil { // got final fragment
			break
		}
		if e != ErrBufferFull { // unexpected error
			err = e
			break
		}

		// Make a copy of the buffer.
		buf := make([]byte, len(frag)) // 这里把 ReadSlice 的返回值 copy 了一份，不再是指向 Reader 中的 buffer
		copy(buf, frag)
		full = append(full, buf)
	}

	// Allocate new buffer to hold the full pieces and the fragment.
	n := 0
	for i := range full {
		n += len(full[i])
	}
	n += len(frag)

	// Copy full pieces and fragment in.
	buf := make([]byte, n)
	n = 0
	for i := range full {
		n += copy(buf[n:], full[i])
	}
	copy(buf[n:], frag)
	return buf, err
}
```
**ReadString 方法**

`ReadString` 源码：
```go
func (b *Reader) ReadString(delim byte) (line string, err error) {
    bytes, err := b.ReadBytes(delim)
    return string(bytes), err
}
```
调用了 `ReadBytes` 方法，并将结果的 `[]byte` 转为 `string` 类型。

**ReadLine 方法**
```go
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
```
`ReadLine` 是一个底层的原始行读取命令。可以使用 `ReadBytes('\n')` 或者 `ReadString('\n')` 来代替这个方法。

**`ReadLine` 尝试返回单独的行，不包括行尾的换行符**。如果一行大于缓存，`isPrefix` 会被设置为 `true`，同时返回该行的开始部分
（等于缓存大小的部分）。该行剩余的部分就会在下次调用的时候返回。当下次调用返回该行剩余部分时，`isPrefix` 将会是 `false` 。
跟 `ReadSlice` 一样，**返回的 `line` 是 `buffer` 的引用**，在下次执行 IO 操作时，`line` 会无效。

建议读取一行使用下面的方式：
```go
line, err := reader.ReadBytes('\n')
line = bytes.TrimRight(line, "\r\n")
```

### Peek 方法

`Peek` 是 "窥视" 的意思，`Peek` 一个鲜明的特点，就是：即使它读取了缓冲区中的数据，也不会更改已读计数的值。

```go
func (b *Reader) Peek(n int) ([]byte, error)
```
**返回的 `[]byte` 是 `buffer` 中的引用**，该切片引用缓存中前 `n` 字节数据。

**`Peek` 方法、`ReadSlice` 方法和 `ReadLine` 方法都有可能会造成内容泄露。这主要是因为它们在正常的情况下都会返回直接基于缓冲区的字节切片，
也因为为这个原因对多 goroutine 是不安全的，也就是在多并发环境下，不能依赖其结果。**。

另外，Reader 的 Peek 方法如果返回的 []byte 长度小于 n，这时返回的 `err != nil` ，用于解释为啥会小于 n。如果 n 大于 reader 的 buffer 长度，err 会是 ErrBufferFull。

### 其他方法
```go
func (b *Reader) Read(p []byte) (n int, err error)
func (b *Reader) ReadByte() (c byte, err error)
func (b *Reader) ReadRune() (r rune, size int, err error)
func (b *Reader) UnreadByte() error
func (b *Reader) UnreadRune() error
func (b *Reader) WriteTo(w io.Writer) (n int64, err error)
```

## bufio.Writer
`bufio.Writer` 结构封装了一个 `io.Writer` 对象。同时实现了 `io.Writer` 接口。
```go
type Writer struct {
    err error		// 写过程中遇到的错误
    buf []byte		// 缓存
    n   int			// 当前缓存中的字节数
    wr  io.Writer	// 底层的 io.Writer 对象
}
```

`bufio.Writer` 类型的字段:
- `err`：`error` 类型的字段。它的值用于表示在向底层写入器写数据时发生的错误。
- `buf`：`[]byte` 类型的字段，代表缓冲区。在初始化之后，它的长度会保持不变。
- `n`：`int` 类型的字段，代表对缓冲区进行下一次写入时的开始索引。我们可以称之为已写计数。
- `wr`：`io.Writer` 类型的字段，代表底层写入器。

两个用于初始化 `bufio.Writer` 的函数：

- `NewWriter` 函数初始化的 `Writer` 值会拥有一个默认尺寸的缓冲区。这个默认尺寸是 `4096` 个字节，即：`4 KB`。
- `NewWriterSize` 函数则将缓冲区尺寸的决定权抛给了使用方。

```go
func NewWriter(wr io.Writer) *Writer

func NewWriterSize(wr io.Writer, size int) *Writer // 可以配置缓冲区的大小
```

### 方法
- `Available` 方法获取缓存中还未使用的字节数（缓存大小 - 字段 n 的值）
- `Buffered` 方法获取写入当前缓存中的字节数（字段 n 的值）
- `Flush` 方法将缓存中的所有数据写入底层的 io.Writer 对象中。

其他实现了 `io` 包的接口方法：
```go
// 实现了 io.ReaderFrom 接口
func (b *Writer) ReadFrom(r io.Reader) (n int64, err error)

// 实现了 io.Writer 接口
func (b *Writer) Write(p []byte) (nn int, err error)

// 实现了 io.ByteWriter 接口
func (b *Writer) WriteByte(c byte) error

// io 中没有该方法的接口，它用于写入单个 Unicode 码点，返回写入的字节数（码点占用的字节），内部实现会根据当前 rune 的范围调用 WriteByte 或 WriteString
func (b *Writer) WriteRune(r rune) (size int, err error)

// 写入字符串，如果返回写入的字节数比 len(s) 小，返回的error会解释原因
func (b *Writer) WriteString(s string) (int, error)
```

### bufio.Writer 类型值中缓冲的数据什么时候会被写到它的底层写入器

`bufio.Writer` 类型有一个名为 `Flush` 的方法，它的主要功能是把相应缓冲区中暂存的所有数据，都写到底层写入器中。数据一旦被写进底层写入器，
该方法就会把它们从缓冲区中删除掉。

`bufio.Writer` 类型值（以下简称 `Writer` 值）拥有的所有数据写入方法都会在必要的时候调用它的 `Flush` 方法。

比如，`Write` 方法有时候会在把数据写进缓冲区之后，调用 `Flush` 方法，以便为后续的新数据腾出空间。`WriteString` 方法的行为与之类似。

`WriteByte` 方法和 `WriteRune` 方法，都会在发现缓冲区中的可写空间不足以容纳新的字节，或 Unicode 字符的时候，调用 `Flush` 方法。

在**通常情况下，只要缓冲区中的可写空间无法容纳需要写入的新数据，`Flush` 方法就一定会被调用**。

## ReadWriter
```go
type ReadWriter struct {
    *Reader
    *Writer
}
```

通过调用 `bufio.NewReadWriter` 函数来初始化：
```go
func NewReadWriter(r *Reader, w *Writer) *ReadWriter
```