## io包中的接口和工具
`strings.Reader`类型主要用于读取字符串，它的指针类型实现的接口比较多，包括：
- io.Reader；
- io.ReaderAt；
- io.ByteReader；
- io.RuneReader；
- io.Seeker；
- io.ByteScanner；
- io.RuneScanner；
- io.WriterTo；

`io.ByteScanner`是`io.ByteReader`的扩展接口，而`io.RuneScanner`又是`io.RuneReader`的扩展接口。

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

为了提高不同程序实体之间的互操作性。以io包中的一些函数为例。

io包中，有这样几个用于拷贝数据的函数，它们是：`io.Copy`、`io.CopyBuffer`和`io.CopyN`。这几个函数在功能上都略有差别，但是它们都首先会接受两个参数，即：
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

## bufio包中的数据类型
bufio包中的数据类型主要有：
- `Reader`；
- `Scanner`；
- `Writer`和`ReadWriter`。

### `bufio.Reader`类型值中的缓冲区的作用
缓冲区其实就是一个**数据存储中介，它介于底层读取器与读取方法及其调用方之间**。所谓的底层读取器，就是在初始化此类值的时候传入的`io.Reader`类型的参数值。

Reader值的读取方法一般都会先从其所属值的缓冲区中读取数据。同时，在必要的时候，它们还会预先从底层读取器那里读出一部分数据，并暂存于缓冲区之中以备后用。

缓冲区的好处是，可以在大多数的时候降低读取方法的执行时间。

`bufio.Reader`类型并不是开箱即用的，因为它包含了一些需要显式初始化的字段。一些字段：
- `buf`：`[]byte`类型的字段，即字节切片，代表缓冲区。虽然它是切片类型的，但是其长度却会在初始化的时候指定，并在之后保持不变。
- `rd`：`io.Reader`类型的字段，代表底层读取器。缓冲区中的数据就是从这里拷贝来的。
- `r`：`int`类型的字段，代表对缓冲区进行下一次读取时的开始索引。我们可以称它为已读计数。
- `w`：`int`类型的字段，代表对缓冲区进行下一次写入时的开始索引。我们可以称之为已写计数。
- `err`：`error`类型的字段。它的值用于表示在从底层读取器获得数据时发生的错误。这里的值在被读取或忽略之后，该字段会被置为`nil`。
- `lastByte`：`int`类型的字段，用于记录缓冲区中最后一个被读取的字节。读回退时会用到它的值。
- `lastRuneSize`：`int`类型的字段，用于记录缓冲区中最后一个被读取的 Unicode 字符所占用的字节数。读回退的时候会用到它的值。这个字段只会在其所
属值的`ReadRune`方法中才会被赋予有意义的值。在其他情况下，它都会被置为`-1`。

两个用于初始化`Reader`值的函数，分别叫`NewReader`和`NewReaderSize`，它们都会返回一个`*bufio.Reader`类型的值。

- `NewReader`函数初始化的`Reade`r值会拥有一个默认尺寸的缓冲区。这个默认尺寸是 4096 个字节，即：4 KB。
- `NewReaderSize`函数则将缓冲区尺寸的决定权抛给了使用方。

### bufio.Writer类型值中缓冲的数据什么时候会被写到它的底层写入器
`bufio.Writer`类型的字段:
- `err`：`error`类型的字段。它的值用于表示在向底层写入器写数据时发生的错误。
- `buf`：`[]byte`类型的字段，代表缓冲区。在初始化之后，它的长度会保持不变。
- `n`：`int`类型的字段，代表对缓冲区进行下一次写入时的开始索引。我们可以称之为已写计数。
- `wr`：`io.Writer`类型的字段，代表底层写入器。

`bufio.Writer`类型有一个名为`Flush`的方法，它的主要功能是把相应缓冲区中暂存的所有数据，都写到底层写入器中。数据一旦被写进底层写入器，该方法就会把它们
从缓冲区中删除掉。

`bufio.Writer`类型值（以下简称Writer值）拥有的所有数据写入方法都会在必要的时候调用它的`Flush`方法。

比如，`Write`方法有时候会在把数据写进缓冲区之后，调用`Flush`方法，以便为后续的新数据腾出空间。`WriteString`方法的行为与之类似。

`WriteByte`方法和`WriteRune`方法，都会在发现缓冲区中的可写空间不足以容纳新的字节，或 Unicode 字符的时候，调用`Flush`方法。

在**通常情况下，只要缓冲区中的可写空间无法容纳需要写入的新数据，`Flush`方法就一定会被调用**。


### bufio.Reader类型读取方法
`bufio.Reader`类型拥有很多用于读取数据的指针方法，这里面有 4 个方法可以作为不同读取流程的代表，它们是：`Peek`、`Read`、`ReadSlice`和`ReadBytes`。

- `Peek`方法的特点是即使读取了缓冲区中的数据，也不会更改已读计数的值。而`Read`方法会在参数值的长度过大，且缓冲区中已无未读字节时，跨过缓冲区并直接向底层读取器索要数据。
`Peek`方法有一个鲜明的特点，那就是：即使它读取了缓冲区中的数据，也不会更改已读计数的值。
- `ReadSlice`方法会在缓冲区的未读部分中寻找给定的分隔符，并在必要时对缓冲区进行填充。如果在填满缓冲区之后仍然未能找到分隔符，那么该方法就会把整个缓冲区作为第一个结果值返回，
同时返回缓冲区已满的错误。
- `ReadBytes`方法会通过调用`ReadSlice`方法，一次又一次地填充缓冲区，并在其中寻找分隔符。除非发生了未预料到的错误或者找到了分隔符，否则这一过程将会一直进行下去。
- Reader值的`ReadLine`方法会依赖于它的`ReadSlice`方法，而其`ReadString`方法则完全依赖于`ReadBytes`方法。

**`Peek`方法、`ReadSlice`方法和`ReadLine`方法都有可能会造成内容泄露。这主要是因为它们在正常的情况下都会返回直接基于缓冲区的字节切片**。
