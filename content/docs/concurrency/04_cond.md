---
title: 条件变量
weight: 4
---

# 条件变量

Go 标准库提供了条件变量 `sync.Cond` 它可以让一组的 goroutine 都在满足特定条件时被唤醒。

`sync.Cond` 不是一个常用的同步机制，但是在条件长时间无法满足时，**与使用 `for {}` 进行忙碌等待相比，`sync.Cond` 能够让出处理器的使用权，提高 CPU 的利用率**。

`sync.Cond` 基于互斥锁/读写锁，它和互斥锁的区别是什么？

互斥锁 `sync.Mutex` 通常用来保护临界区和共享资源，条件变量 `sync.Cond` 用来协调想要访问共享资源的 goroutine。

`sync.Cond` 经常用在多个 goroutine 等待，一个 goroutine 通知的场景。

比如有一个 goroutine 在异步地接收数据，剩下的多个 goroutine 必须等待这个协程接收完数据，才能读取到正确的数据。这个时候，就需要有个全局的变量来标志第一
个 goroutine 数据是否接受完毕，剩下的 goroutine，反复检查该变量的值，直到满足要求。

当然也可以创建多个 channel，每个 goroutine 阻塞在一个 channel 上，由接收数据的 goroutine 在数据接收完毕后，逐个通知。但是这种方式更复杂一点。

## 使用

`NewCond` 用来创建 `sync.Cond` 实例，`sync.Cond` 暴露了几个方法：

- `Broadcast` 用来唤醒所有等待条件变量的 goroutine，无需锁保护。
- `Signal` 唤醒一个 goroutine。
- `Wait` 调用 `Wait` 会自动释放锁，并挂起调用者所在的 goroutine，也就是当前 goroutine 会阻塞在 `Wait` 方法调用的地方。如果其他 goroutine 调用了 `Signal` 或 `Broadcast` 唤醒
了该 goroutine，那么 `Wait` 方法在结束阻塞时，会重新加锁，并且继续执行 `Wait` 后面的代码。

```go
var status int64

func main() {
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		go listen(c)
	}
	time.Sleep(1 * time.Second)
	go broadcast(c)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func broadcast(c *sync.Cond) {
	c.L.Lock()
	atomic.StoreInt64(&status, 1)
	c.Broadcast()
	c.L.Unlock()
}

func listen(c *sync.Cond) {
	c.L.Lock()
	// 使用了 for !condition() 而非 if，是因为当前 goroutine 被唤醒时，条件不一定符合要求，需要再次 Wait 等待下次被唤醒
	// 例如，如果 broadcast 没有调用 atomic.StoreInt64(&status, 1) 将 status 设置为 1，这里判断条件后会再次阻塞
	for atomic.LoadInt64(&status) != 1 { 
		c.Wait()
	}
	fmt.Println("listen")
	c.L.Unlock()
}
```

- `status`：互斥锁需要保护的条件变量。
- `listen()` 调用 `Wait()` 等待通知，直到 `status` 为 1。
- `broadcast()` 将 `status` 置为 1，调用 `Broadcast()` 通知所有等待的 goroutine。

运行：
```
$ go run main.go
listen
...
listen
```

打印出 10 次 “listen” 并结束调用。


## 原理

`sync.Cond` 结构体：

```go
// src/sync/cond.go
type Cond struct {
    noCopy  noCopy
    L       Locker
    notify  notifyList
    checker copyChecker
}

type notifyList struct {
	// wait 和 notify 分别表示当前正在等待的和已经通知到的 goroutine 的索引
    wait uint32
    notify uint32
    
    lock mutex
	// head 和 tail 分别指向的链表的头和尾
    head *sudog
    tail *sudog
}
```

- `noCopy`：用于保证结构体不会在编译期间拷贝
- `copyChecker`：用于禁止运行期间发生的拷贝
- `L`：用于保护 `notify` 字段
- `notify`：一个 goroutine 链表，它是实现同步机制的核心结构

`Wait` 方法会将当前 goroutine 陷入休眠状态，它的执行过程分成以下两个步骤：

- 调用 `runtime.notifyListAdd` 将等待计数器加 1 并解锁；
- 调用 `runtime.notifyListWait` 等待其他 goroutine 的唤醒并加锁：

```go
func (c *Cond) Wait() {
	c.checker.check()
	t := runtime_notifyListAdd(&c.notify)
	c.L.Unlock()
	// 休眠直到被唤醒
	runtime_notifyListWait(&c.notify, t)
	c.L.Lock()
}

func notifyListAdd(l *notifyList) uint32 {
	return atomic.Xadd(&l.wait, 1) - 1
}

// notifyListWait 获取当前 goroutine 并将它追加到 goroutine 通知链表的最末端
func notifyListWait(l *notifyList, t uint32) {
    s := acquireSudog()
    s.g = getg()
    s.ticket = t
    if l.tail == nil {
        l.head = s
    } else {
        l.tail.next = s
    }
    l.tail = s
	// 调用 runtime.goparkunlock 使当前 goroutine 陷入休眠
	// 该函数会直接让出当前处理器的使用权并等待调度器的唤醒
    goparkunlock(&l.lock, waitReasonSyncCondWait, traceEvGoBlockCond, 3)
    releaseSudog(s)
}
```

`Signal` 方法会唤醒队列最前面的 goroutine，`Broadcast` 方法会唤醒队列中全部的 goroutine：

```go
func (c *Cond) Signal() {
	c.checker.check()
	runtime_notifyListNotifyOne(&c.notify)
}

func (c *Cond) Broadcast() {
	c.checker.check()
	runtime_notifyListNotifyAll(&c.notify)
}
```

`notifyListNotifyOne` 从 `notifyList` 链表中找到满足 `sudog.ticket == l.notify` 条件的 goroutine 并通过 `runtime.readyWithTime` 唤醒：

```go
// src/runtime/sema.go#L554
func notifyListNotifyOne(l *notifyList) {
    t := l.notify
    atomic.Store(&l.notify, t+1)

    for p, s := (*sudog)(nil), l.head; s != nil; p, s = s, s.next {
        if s.ticket == t {
            n := s.next
            if p != nil {
                p.next = n
            } else {
                l.head = n
            }
			if n == nil {
                l.tail = p
            }
            s.next = nil
            readyWithTime(s, 4)
            return
        }
    }
}
```

`notifyListNotifyAll` 会依次通过 `runtime.readyWithTime` 唤醒链表中所有 goroutine：

```go
func notifyListNotifyAll(l *notifyList) {
	s := l.head
	l.head = nil
	l.tail = nil

	atomic.Store(&l.notify, atomic.Load(&l.wait))

	for s != nil {
		next := s.next
		s.next = nil
		readyWithTime(s, 4)
		s = next
	}
}
```

**goroutine 的唤醒顺序也是按照加入队列的先后顺序，先加入的会先被唤醒**。