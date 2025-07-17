---
title: Channel
weight: 9
---

```
Don’t communicate by sharing memory; share memory by communicating.
不要通过共享内存来通信，通过通信来共享内存。
```

这是 Go 语言最重要的编程理念。goroutine 通过 channel 向另一个 goroutine 发送消息，channel 和 goroutine 结合，可以实现用通信代替共享内存的 CSP （Communicating Sequential Process）模型。

## 使用

创建 channel：

```go
// 无缓冲 channel
ch := make(chan int)

// 带缓冲 channel，缓冲区为 3
ch = make(chan int, 3)

// ok 为 false 表示通道已经关闭
val, ok := <- ch
```

> channel 的零值是 `nil`。

### 无缓冲 channel

无缓冲 channel 也叫做同步 channel：

- 一个 goroutine 基于一个无缓冲 channel 发送数据，那么就会阻塞，直到另一个 goroutine 在相同的 channel 上执行接收操作。
- 一个 goroutine 基于一个无缓冲 channel 先执行了接收操作，也会阻塞，直到另一个 goroutine 在相同的 channel 上执行发送操作

### 带缓冲 channel

带缓冲的 channel 有一个缓冲区：

- 若缓冲区未满则不会阻塞，发送者可以不断的发送数据。当缓冲区满了后，发送者就会阻塞。
- 当缓冲区为空时，接受者就会阻塞，直至有新的数据

### 关闭 channel

使用 `close` 函数关闭 channel：

- channel 关闭后不能再发送数据
- channel 关闭后可以接收已经发送成功的数据。
- channel 关闭后如果 channel 中没有数据，那么接收者会收到一个 channel 元素的零值。

`close` 表示这个 channel 不会再继续发送数据，所以要**在发送者所在的 goroutine 去关闭 channel**。

{{< callout type="warning" >}}
- 关闭一个 `nil` 的 channel 会导致 panic。
- 重复关闭 channel 会导致 panic。
- 向已关闭的 channel 发送值会导致 panic。
{{< /callout >}}


### 单向 channel

当一个 channel 作为一个函数参数时，它一般总是被专门用于**只发送或者只接收**。

- `chan<- int` 表示一个只发送 `int` 的 channel。
- `<-chan int` 表示一个只接收 `int` 的 channel。

### cap 和 len

- `cap` 函数可以获取 channel 内部缓冲区的容量。
- `len` 函数可以获取 channel 内部缓冲区有效元素的个数。

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

我们可以使用容量只有 `1` 的 channel 来保证最多只有一个 goroutine 在同一时刻访问一个共享变量：

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

channel 本质上就是一个有锁的环形队列，channel 的结构体 `hchan`：

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

type sudog struct {
	// 指向当前的 goroutine
	g *g
	
	// 指向下一个 goroutine
	next *sudog
	// 指向上一个 goroutine
	prev *sudog
	// 指向元素数据
	elem unsafe.Pointer
    // ...
}
```

### 创建 channel

创建 channel 要使用 `make`，编译器会将 `make` 转换成 `makechan` 或者 `makechan64` 函数：

```go
// src/runtime/chan.go#L72
func makechan(t *chantype, size int) *hchan {
	elem := t.Elem

	// compiler checks this but be safe.
	// ...
	
	var c *hchan
	switch {
	case mem == 0:
		// 无缓冲 channel
		// 调用 mallocgc 方法分配一段连续的内存空间
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		c.buf = c.raceaddr()
	case elem.PtrBytes == 0:
		// channel 存储的元素类型不是指针
		// 分配一块连续的内存给 hchan 和底层数组
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
        // 默认情况下，进行两次内存分配操作，分别为 hchan 和缓冲区分配内存
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}

	// 设置元素大小，元素类型，循环数组的长度
	c.elemsize = uint16(elem.Size_)
	c.elemtype = elem
	c.dataqsiz = uint(size)
	lockInit(&c.lock, lockRankHchan)
    // ...
	return c
}
```

使用 mallocgc 函数创建 channel，就意味着 channel 都是分配在堆上的。所以**当一个 channel 没有被任何 goroutine 引用时，是会被 GC 回收的**。

### 向 channel 发送数据

发送操作，也就是 `ch <- i` 语句，编译器最终会将该语句转换成 `chansend` 函数：

```go
// src/runtime/chan.go
// block 为 true 时，表示当前操作是阻塞的 
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	if c == nil {
		// 不可以阻塞，直接返回 false，表示未发送成功
		if !block {
			return false
		}
        // 挂起当前 goroutine
		gopark(nil, nil, waitReasonChanSendNilChan, traceBlockForever, 2)
		throw("unreachable")
	}
	// ...
	if !block && c.closed == 0 && full(c) {
		return false
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	// 执行发送数据的逻辑之前，先为当前 channel 加锁，防止多个线程并发修改数据
	lock(&c.lock)

	// 如果 channel 已经关闭，那么向该 channel 发送数据会导致 panic：send on closed channel
	if c.closed != 0 {
		// 解锁
		unlock(&c.lock)
		// panic
		panic(plainError("send on closed channel"))
	}

	// 当前接收队列里存在 goroutine，通过 runtime.send 直接将数据发送给阻塞的接收者
	if sg := c.recvq.dequeue(); sg != nil {
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	// 走到这里，说明没有等待数据的接收者
	
	// 对于有缓冲的 channel，并且还有缓冲空间
	if c.qcount < c.dataqsiz {
        // 计算出下一个可以存储数据的位置
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			racenotify(c, c.sendx, nil)
		}
		// 将发送的数据拷贝到缓冲区中并增加 sendx 索引和 qcount 计数器
		typedmemmove(c.elemtype, qp, ep)
		// sendx 索引 +1
		c.sendx++
		// 由于 buf 是一个循环数组，所以当 sendx 等于 dataqsiz 时会重新回到数组开始的位置。
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		// 释放锁
		unlock(&c.lock)
		return true
	}

	// 走到这里，说明缓冲空间已满，或者是无缓冲 channel
	
	// 如果不可以阻塞，直接返回 false，表示未发送成功
	if !block {
		unlock(&c.lock)
		return false
	}

	// 缓冲空间已满或者是无缓冲 channel，发送方会被阻塞
	
	// 获取当前发送数据的 goroutine 的指针
	gp := getg()
	// 构造一个 sudog
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// 设置这一次阻塞发送的相关信息
	mysg.elem = ep // 待发送数据的内存地址
	mysg.waitlink = nil
	mysg.g = gp // 当前发送数据的 goroutine 的指针
	mysg.isSelect = false // 是否在 select 中
	mysg.c = c // 发送的 channel
	gp.waiting = mysg
	gp.param = nil
	// 将 sudog 放入到发送等待队列
	c.sendq.enqueue(mysg)
	// 挂起当前 goroutine，等待唤醒
	gp.parkingOnChan.Store(true)
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceBlockChanSend, 2)
	KeepAlive(ep)

	// goroutine 开始被唤醒了
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	closed := !mysg.success
	gp.param = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	// 移除 mysg 上绑定的 channel
	mysg.c = nil
	releaseSudog(mysg)
	if closed {
		if c.closed == 0 {
			throw("chansend: spurious wakeup")
		}
		// 被唤醒了，但是 channel 已经关闭了，panic
		panic(plainError("send on closed channel"))
	}
	// 返回 true 表示已经成功向 channel 发送了数据
	return true
}
```

`send` 发送数据：

```go
func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
	// ...
	// sg 是接收者的 sudog 结构
	// sg.elem 指向接收到的值存放的位置，如 val <- ch，指的就是 &val
	if sg.elem != nil {
		// 直接拷贝内存到 val <- ch 表达式中变量 val 所在的内存地址（&val）上
		sendDirect(c.elemtype, sg, ep)
		sg.elem = nil
	}
	// 获取 sudog 上绑定的等待接收的 goroutine 的指针
	gp := sg.g
	unlockf()
	gp.param = unsafe.Pointer(sg)
	// 唤醒等待接收的 goroutine
	goready(gp, skip+1)
}
```

> `goready` 是将 goroutine 的状态改成 `runnable`，然后需要等待调度器的调度。

```go
func sendDirect(t *_type, sg *sudog, src unsafe.Pointer) {
	// src 是当前 goroutine 发送的数据的内存地址
	// dst 是接收者的值的存放位置
	dst := sg.elem
	// 写屏障
	typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.size)
	// 拷贝内存数据
	memmove(dst, src, t.size)
}
```

### 从 channel 接收数据

Go 中可以使用两种不同的方式去接收 channel 中的数据：

```go
i <- ch
i, ok <- ch
```

编译器的处理后分别会转换成 `chanrecv1`，`chanrecv2`：

```go
// src/runtime/chan.go
func chanrecv1(c *hchan, elem unsafe.Pointer) {
	chanrecv(c, elem, true)
}

func chanrecv2(c *hchan, elem unsafe.Pointer) (received bool) {
	_, received = chanrecv(c, elem, true)
	return
}
```

两个方法最终还是调用了 `chanrecv` 函数：

```go
// src/runtime/chan.go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
    // ...
	// channel 是 nil
	if c == nil {
		// 不可以阻塞，直接返回
		if !block {
			return
		}
		// 挂起当前 goroutine
		gopark(nil, nil, waitReasonChanReceiveNilChan, traceBlockForever, 2)
		throw("unreachable")
	}
	
	if !block && empty(c) {
		if atomic.Load(&c.closed) == 0 {
			return
		}
		if empty(c) {
			if raceenabled {
				raceacquire(c.raceaddr())
			}
			if ep != nil {
				typedmemclr(c.elemtype, ep)
			}
			return true, false
		}
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	// 执行接收数据的逻辑之前，先为当前 channel 加锁
	lock(&c.lock)

	// channel 已关闭
	if c.closed != 0 {
		// 底层的循环数组 buf 中没有元素
		if c.qcount == 0 {
			if raceenabled {
				raceacquire(c.raceaddr())
			}
			// 释放锁
			unlock(&c.lock)
			if ep != nil {
				// typedmemclr 根据类型清理相应地址的内存
				typedmemclr(c.elemtype, ep)
			}
			return true, false
		}
	} else {
		// channel 未关闭，并且等待发送队列里存在 goroutine
		// 发送的 goroutine 被阻塞，那有两种情况：
		// 1. 这是一个非缓冲型的 channel
		// 2. 缓冲型的 channel，但是 buf 满了
		// recv 直接进行内存拷贝
		if sg := c.sendq.dequeue(); sg != nil {
			recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
			return true, true
		}
	}
	
	// channel 未关闭
	// 缓冲型 channel 并且 buf 里有元素，可以正常接收
	if c.qcount > 0 {
        // 直接从循环数组里取出要接收的元素
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			racenotify(c, c.recvx, nil)
		}
		// 这里表示，代码中没有忽略要接收的值，不是 "<- ch"，而是 "val <- ch"，ep 指向 val
		if ep != nil {
			// 拷贝数据
			typedmemmove(c.elemtype, ep, qp)
		}
		// 清理掉循环数组里相应位置的值
		typedmemclr(c.elemtype, qp)
		// recvx 索引 +1
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		// 元素个数 -1
		c.qcount--
		unlock(&c.lock)
		return true, true
	}

	// 非阻塞接收，释放锁
	// selected 返回 false，因为没有接收到值
	if !block {
		unlock(&c.lock)
		return false, false
	}
	
	// 走到这里说明 buf 是空的
	
	// 没有数据可接收，阻塞当前接收的 goroutine
	
	// 获取当前接收的 goroutine
	gp := getg()
	// 构造一个 sudog
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
    // 设置这一次阻塞接收的相关信息
	mysg.elem = ep // 待接收数据的地址
	mysg.waitlink = nil
	gp.waiting = mysg
	mysg.g = gp // 当前接收的 goroutine 指针
	mysg.isSelect = false // 是否在 select 中
	mysg.c = c // 接收的 channel
	gp.param = nil
	// 将 sudog 放入到接收等待队列
	c.recvq.enqueue(mysg)
	gp.parkingOnChan.Store(true)
	// 挂起当前接收 goroutine
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanReceive, traceBlockChanRecv, 2)

	// 被唤醒了
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	success := mysg.success
	gp.param = nil
	mysg.c = nil
	releaseSudog(mysg)
	return true, success
}
```

`recv` 接收数据：

```go
func recv(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
	// 无缓冲的 channel
	if c.dataqsiz == 0 {
		if raceenabled {
			racesync(c, sg)
		}
        // 这里表示，代码中没有忽略要接收的值，不是 "<- ch"，而是 "val <- ch"，ep 指向 val
		if ep != nil {
			// 直接拷贝数据
			recvDirect(c.elemtype, sg, ep)
		}
	} else {
		// 缓冲型的 channel，但是 buf 已满
		// 将底层的循环数组 buf 队首的元素拷贝到接收数据的地址
		// 将发送者的数据放入 buf
		qp := chanbuf(c, c.recvx)
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}

		// 将发送者数据拷贝到 buf
		typedmemmove(c.elemtype, qp, sg.elem)
		// 增加 recvx 索引
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.sendx = c.recvx
	}
	sg.elem = nil
	gp := sg.g

	// 释放锁
	unlockf()
	gp.param = unsafe.Pointer(sg)
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}

	// 唤醒发送的 goroutine
	goready(gp, skip+1)
}
```

### 关闭 channel

`close` 关闭 channel 会被编译器转换成 `closechan` 函数：

```go
// src/runtime/chan.go#L357
func closechan(c *hchan) {
	// 关闭一个 nil 的 channel，panic
	if c == nil {
		panic(plainError("close of nil channel"))
	}

	// 先加锁
	lock(&c.lock)
	// 重复关闭，panic
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("close of closed channel"))
	}
	// ...
	// 设置 channel 关闭的标志位
	c.closed = 1

	var glist gList

    // 将 channel 等待接收队列的里 sudog 释放
	for {
		// 从接收队列里取出一个 sudog
		sg := c.recvq.dequeue()
		// 接收队列空了，跳出循环
		if sg == nil {
			break
		}
		// 
		if sg.elem != nil {
			typedmemclr(c.elemtype, sg.elem)
			sg.elem = nil
		}
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		// 获取接收 goroutine 的指针
		gp := sg.g
		gp.param = unsafe.Pointer(sg)
		sg.success = false
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		// 放入链表
		glist.push(gp)
	}

	// 将 channel 等待发送队列的里 sudog 释放
	// 如果存在，这些 goroutine 将会 panic
	// 可以查看 chansend 函数中的逻辑：
	// 对于发送者，如果被唤醒后 channel 已关闭，则会 panic
	for {
        // 从发送队列里取出一个 sudog
		sg := c.sendq.dequeue()
        // 发送队列空了，跳出循环
		if sg == nil {
			break
		}
		sg.elem = nil
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
        // 获取发送 goroutine 的指针
		gp := sg.g
		gp.param = unsafe.Pointer(sg)
		sg.success = false
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
        // 放入链表
		glist.push(gp)
	}
	// 释放锁
	unlock(&c.lock)

	// 遍历链表，唤醒所有 goroutine
	for !glist.empty() {
		gp := glist.pop()
		gp.schedlink = 0
		goready(gp, 3)
	}
}
```

`recvq` 和 `sendq` 中的所有 goroutine 被唤醒后，会分别去执行 `chanrecv` 和 `chansend` 中 `gopark` 后面的代码。