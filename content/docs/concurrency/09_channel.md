---
title: Channel
weight: 9
---

# Channel

```
Don’t communicate by sharing memory; share memory by communicating.
不要通过共享内存来通信，通过通信来共享内存。
```

这是 Go 语言最重要的编程理念。goroutine 通过 channel 向另一个 goroutine 发送消息 channel 和 goroutine 结合，可以实现用通信代替共享内存的 CSP （Communicating Sequential Process）模型。

## 使用

创建 channel：

```go
// 无缓冲 channel
ch := make(chan int)

// 带缓冲 channel，缓冲区为 3
ch = make(chan int, 3)
```

### 无缓冲 channel

无缓冲 `channel` 也叫做同步 `channel`：

- 一个 goroutine 基于一个无缓冲 `channel` 发送数据，那么就会阻塞，直到另一个 goroutine 在相同的 `channel` 上执行接收操作。
- 一个 goroutine 基于一个无缓冲 `channel` 先执行了接收操作，也会阻塞，直到另一个 goroutine 在相同的 `channel` 上执行发送操作

### 带缓冲 channel

带缓冲的 `channel` 有一个缓冲区：

- 若缓冲区未满则不会阻塞，发送者可以不断的发送数据。当缓冲区满了后，发送者就会阻塞。
- 当缓冲区为空时，接受者就会阻塞，直至有新的数据

### close

使用 `close` 函数关闭 `channel`，`channel` 关闭后不能再发送数据，但是可以接收已经发送成功的数据。
`close` 以后如果 `channel` 中没有数据，那么接收者会收到一个零值。

`close` 表示这个 `channel` 不会再继续发送数据，所以要在**发送者**所在的 goroutine 调用。

> channel 的零值是 `nil`。关闭一个 `nil` 的 channel 会导致程序 panic。

### 单向 channel

当一个 `channel` 作为一个函数参数时，它一般总是被专门用于**只发送或者只接收**。

类型 `chan<- int` 表示一个只发送 `int` 的 `channel`。相反，类型 `<-chan int` 表示一个只接收 `int` 的 `channel`。

### cap 和 len

`cap` 函数可以获取 `channel` 内部缓冲区的容量。
`len` 函数可以获取 `channel` 内部缓冲区有效元素的个数。

### 使用 range 遍历 channel

使用 `range` 循环可以遍历 channel，它依次从 channel 中接收数据，当 channel 被关闭并且没有值可接收时跳出循环：

```go
ch := make(chan int, 3)
ch <- 1
ch <- 2
ch <- 3

// 关闭 channel
// 如果不关闭 channel，range 就会阻塞当前 goroutine, 直到 channel 关闭
close(ch)

for v := range ch {
    fmt.Println(v) 
}
```

### 使用 channel 实现互斥锁

我们可以使用容量只有 `1` 的 `channel` 来保证最多只有一个 goroutine 在同一时刻访问一个共享变量：

```go
var (
	sema = make(chan struct{}, 1) // a binary semaphore guarding balance 
	balance int
)

func Deposit(amount int) {
	sema <- struct{}{} // acquire lock 
	balance = balance + amount
	<-sema // release lock
}

func Balance() int {
	sema <- struct{}{} // acquire lock 
	b := balance
	<-sema // release lock 
	// return b
}
```

## 原理

channel 的结构体 `hchan`：

```go
// src/runtime/chan.go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}
```

- `qcount`：channel 中的元素个数
- `dataqsiz`：channel 中的循环队列的长度
- `buf`：channel 的缓冲区数据指针，指向底层的循环数组，只针对有缓冲的 channel。
- `elemsize`：channel 中元素大小
- `elemtype`：channel 中元素类型
- `closed`：channel 是否被关闭的标志位
- `sendx`：表示当前可以发送的元素在底层循环数组中位置索引
- `recvx`：表示当前可以发送的元素在底层循环数组中位置索引
- `sendq`：向 channel 发送数据而被阻塞的 goroutine 队列
- `recvq`：读取 channel 的数据而被阻塞的 goroutine 队列
- `lock`：保护 `hchan` 中所有字段

`waitq` 是一个双向链表，链表中所有的元素都是 `sudog`：

```go
type waitq struct {
	first *sudog
	last  *sudog
}
```

### 向 channel 发送数据

发送操作对应底层的 `chansend` 函数：

```go
// src/runtime/chan.go
```

### 从 channel 接收数据

