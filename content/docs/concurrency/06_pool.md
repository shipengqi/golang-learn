---
title: Pool
weight: 6
---

Go 从 1.3 版本开始提供了对象重用的机制，即 `sync.Pool`。`sync.Pool` 用来保存可以被重复使用的**临时**对象，避免了重复创建和销毁临时对象带来的消耗，降低 GC 压力，提高性能。

**`sync.Pool` 是可伸缩的，也是并发安全的**。可以在多个 goroutine 中并发调用 `sync.Pool` 存取对象。

## 使用

```go
var buffers = sync.Pool{
	New: func() interface{} { 
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	buffers.Put(buf)
}
```

`New`：类型是 `func() interface{}`，用来创建新的元素。
`Get`：从 Pool 中取出一个元素，如果没有更多的空闲元素，就调用 `New` 创建新的元素。如果没有设置 `New` 那么可能返回 `nil`。
`Put`：将一个元素放回 Pool 中，使该元素可以重复使用，如果 `Put` 的值是 nil，会被忽略。

## 原理

Go 1.13 之前的 `sync.Pool` 的问题：

1. 每次 GC 都会回收创建的对象。
   - 缓存元素数量太多，就会导致 STW 耗时变长；
   - 缓存元素都被回收后，会导致 `Get` 命中率下降，`Get` 方法不得不新创建很多对象。
2. 底层使用了 `Mutex`，并发请求竞争锁激烈的时候，会导致性能的下降。

Go 1.13 进行了优化，移除了 `Mutex`，增加了 `victim` 缓存。

Pool 的结构体：

```go
type Pool struct {
	noCopy noCopy

    // 每个 P 的本地队列，实际类型为 [P]poolLocal
	local     unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal
	// [P]poolLocal的大小
	localSize uintptr        // size of the local array

	victim     unsafe.Pointer // local from previous cycle
	victimSize uintptr        // size of victims array

	// 自定义的对象创建回调函数，当 pool 中无可用对象时会调用此函数
	New func() interface{}
}
```

重要的两个字段是 `local` 和 `victim`，都是要用来存储空闲的元素。

`local` 字段存储指向 `[P]poolLocal` 数组（严格来说，它是一个切片）的指针。访问时，P 的 id 对应 `[P]poolLocal` 下标索引。通过这样的设计，多个 goroutine 使用同一个 Pool 时，减少了竞争，提升了性能。

在 `src/sync/pool.go` 文件的 `init` 函数里，注册了 GC 发生时，如何清理 Pool 的函数：

```go
func init() {
	runtime_registerPoolCleanup(poolCleanup)
}
```

GC 时 `sync.Pool` 的处理逻辑：

```go
func poolCleanup() {
    // 丢弃当前 victim, STW 所以不用加锁
    for _, p := range oldPools {
        p.victim = nil
        p.victimSize = 0
    }

    // 将 local 复制给 victim, 并将原 local 置为 nil
    for _, p := range allPools {
        p.victim = p.local
        p.victimSize = p.localSize
        p.local = nil
        p.localSize = 0
    }

    oldPools, allPools = allPools, nil
}
```

`poolCleanup` 会在 STW 阶段被调用。主要是将 `local` 和 `victim` 作交换，这样也就不致于让 GC 把所有的 Pool 都清空了。

> 如果 `sync.Pool` 的获取、释放速度稳定，那么就不会有新的池对象进行分配。如果获取的速度下降了，那么对象可能会在两个 GC 周期内被释放，而不是 Go 1.13 以前的一个 GC 周期。

调用 `Get` 时，会先从 `victim` 中获取，如果没有找到，则就会从 `local` 中获取，如果 `local` 中也没有，就会执行 `New` 创建新的元素。

**Pool 中的缓存对象虽然还被清除，但是在两次 GC 之间的窗口期内，对象可被重复复用多次**。例如一个 HTTP 服务在 1 秒内处理 1 万次请求，期间可能触发 0 次 GC，所有请求共享 Pool 中的对象。

### 内存泄露

前面的示例代码中实现了一个 `buffer` 池，这个实现可能会有内存泄漏的风险。为什么？

因为在取出 `bytes.Buffer` 之后，可以给这个 `buffer` 中增加大量的 `byte` 数据，这会导致底层的 byte slice 的容量可能会变得很大。这个时候，即使 Reset 再放回到池子中，这些 byte slice 的容量不会改变，所占的空间依然很大。

`Reset` 的实现：

```go
// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (b *Buffer) Reset() {
	// 基于已有 slice 创建新 slice 对象，不会拷贝原数组或者原切片中的数据，新 slice 和老 slice 共用底层数组
	// 它只会创建一个 指向原数组的 切片结构体，新老 slice 对底层数组的更改都会影响到彼此。
	b.buf = b.buf[:0]
	b.off = 0
	b.lastRead = opInvalid
}
```

切片结构体：

```go
// runtime/slice.go
type slice struct {
    array unsafe.Pointer // 元素指针，指向底层数组
    len   int // 长度 
    cap   int // 容量
}
```

因为 Pool 回收的机制，这些大的 Buffer 可能不会被立即回收，而是会占用很大的空间，这属于内存泄漏的问题。

Go 的标准库 `encoding/json` 和 `fmt` 修复这个问题的方法是增加了检查逻辑：如果放回的 `buffer` 超过一定大小，就直接丢弃掉，不再放到池子中。

```go
// 超过一定大小，直接丢弃掉
if cap(p.buf) > 64<<0 {
	return
}

// 放回 pool
```

所以在使用 `sync.Pool` 时，回收 `buffer` 的时候，**一定要检查回收的对象的大小**。如果 `buffer` 太大，就直接丢弃掉。

### 优化内存使用

使用 `buffer` 池的时候，可以根据实际元素的大小来分为几个 `buffer` 池。比如：

- 小于 512 byte 的元素的 `buffer` 占一个池子；
- 其次，小于 1K byte 大小的元素占一个池子；
- 再次，小于 4K byte 大小的元素占一个池子。

这样分成几个池子以后，就可以根据需要，到所需大小的池子中获取 `buffer` 了。

例如标准库 [net/http/server.go](https://github.com/golang/go/blob/master/src/net/http/server.go) 的实现：

```go
var (
	bufioReaderPool   sync.Pool
	bufioWriter2kPool sync.Pool
	bufioWriter4kPool sync.Pool
)

var copyBufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 32*1024)
		return &b
	},
}

func bufioWriterPool(size int) *sync.Pool {
	switch size {
	case 2 << 10:
		return &bufioWriter2kPool
	case 4 << 10:
		return &bufioWriter4kPool
	}
	return nil
}
```

还有第三方的实现：

- [bytebufferpool](https://github.com/valyala/bytebufferpool)
