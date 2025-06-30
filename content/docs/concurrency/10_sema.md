---
title: 信号量
weight: 10
---

信号量（Semaphore）是一种用于实现多进程或多线程之间同步和互斥的机制。

信号量可以简单理解为一个整型数，包含两种操作：`P`（Proberen，测试）操作和 `V`（Verhogen，增加）操作。其中，`P` 操作会尝试获取一个信号量，如果信号量的值大于 0，则将信号量的值减 1 并
继续执行。否则，当前进程或线程就会被阻塞，直到有其他进程或线程释放这个信号量为止。V 操作则是释放一个信号量，将信号量的值加 1。

`P` 操作和 `V` 操作可以看做是对资源的获取和释放。

Go 的 `WaitGroup` 和 `Metux` 都是通过信号量来控制 goroutine 的阻塞和唤醒，例如 `Mutex` 结构体中的 `sema`：

```go
type Mutex struct {
    state int32
	sema  uint32
}
```

`Metux` 本质上就是基于信号量（sema）+ 原子操作来实现并发控制的。

Go 操作信号量的方法：

```go
// src/sync/runtime.go
// 阻塞等待直到 s 大于 0，然后立刻将 s 减去 1
func runtime_Semacquire(s *uint32)

// 类似于 runtime_Semacquire
// 如果 lifo 为 true，waiter 将会被插入到队列的头部，否则插入到队列尾部
// skipframes 是跟踪过程中要省略的帧数，从这里开始计算
func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)

// 将 s 增加 1，然后通知阻塞在 runtime_Semacquire 的 goroutine
// 如果 handoff 为 true，传递信号到队列头部的 waiter
// skipframes 是跟踪过程中要省略的帧数，从这里开始计算
func runtime_Semrelease(s *uint32, handoff bool, skipframes int)
```

Acquire 和 Release 分别对应了 `P` 操作和 `V` 操作。

## Acquire 信号量

```go
// src/runtime/sema.go
//go:linkname sync_runtime_Semacquire sync.runtime_Semacquire
func sync_runtime_Semacquire(addr *uint32) {
	semacquire1(addr, false, semaBlockProfile, 0, waitReasonSemacquire)
}

//go:linkname sync_runtime_SemacquireMutex sync.runtime_SemacquireMutex
func sync_runtime_SemacquireMutex(addr *uint32, lifo bool, skipframes int) {
	semacquire1(addr, lifo, semaBlockProfile|semaMutexProfile, skipframes, waitReasonSyncMutexLock)
}
```

`runtime_Semacquire` 和 `runtime_SemacquireMutex` 最终都是调用了 `semacquire1` 函数：

```go
func semacquire1(addr *uint32, lifo bool, profile semaProfileFlags, skipframes int, reason waitReason) {
	// 检查当前 goroutine 是否在 G 栈上
	gp := getg()
	if gp != gp.m.curg {
		throw("semacquire not on the G stack")
	}

	// Easy case.
	// 快速路径：信号量大于 0，直接返回，信号量 -1
	if cansemacquire(addr) {
		return
	}

	// Harder case:
	// 慢路径：从池中获取 sudog 结构（避免频繁内存分配）
	// sudog 表示一个等待中的 goroutine
	s := acquireSudog()
	// 将信号量的地址放到到 semtable 中
	// 返回一个 semaRoot 类型
	root := semtable.rootFor(addr)
	// ...
	for {
		lockWithRank(&root.lock, lockRankRoot)
		// 等待计数 +1
		root.nwait.Add(1)
		// 再次检查信号量是否大于 0，避免错误唤醒
		if cansemacquire(addr) {
			root.nwait.Add(-1)
			unlock(&root.lock)
			break
		}
		// 将 sudog 放入到 semaRoot 的等待者队列
		// queue 会将 sudog 和 g 关联起来
		root.queue(addr, s, lifo)
		// 挂起当前 goroutine
		goparkunlock(&root.lock, reason, traceBlockSync, 4+skipframes)
		// 被唤醒后重新检查
		if s.ticket != 0 || cansemacquire(addr) {
			break
		}
	}
	if s.releasetime > 0 {
		blockevent(s.releasetime-t0, 3+skipframes)
	}
	// 释放 sudog 放回池内
	releaseSudog(s)
}
```

`cansemacquire` 就是判断信号量的值，若等于 0，则直接返回 `false`，否则，CAS 操作信号量 -1，成功则返回 `true`：

```go
func cansemacquire(addr *uint32) bool {
    for {
        v := atomic.Load(addr)
		// 等于 0，表示没有资源
        if v == 0 {
            return false
        }
        if atomic.Cas(addr, v, v-1) {
            return true
        }
    }
}
```

`semtable` 是一个 `semTable` 类型，`semTable.rootFor` 返回的是一个 `semaRoot` 类型：

```go
// src/runtime/sema.go
type semaRoot struct {
	// 保护本结构的自旋锁（非 Go 级别的 mutex，是更底层的锁定机制）
	lock  mutex
	treap *sudog        // 等待者队列（平衡树）的根节点
    nwait atomic.Uint32 // 等待者的数量
}

var semtable semTable

type semTable [semTabSize]struct {
	root semaRoot
	pad  [cpu.CacheLinePadSize - unsafe.Sizeof(semaRoot{})]byte
}

// rootFor 本质上就是将 semaRoot 与信号量绑定
func (t *semTable) rootFor(addr *uint32) *semaRoot {
    return &t[(uintptr(unsafe.Pointer(addr))>>3)%semTabSize].root
}


func (root *semaRoot) queue(addr *uint32, s *sudog, lifo bool) {
	// 释放信号量时，唤醒 g 需要用到
	s.g = getg()
	// ...
}
```

## Release 信号量

```go
// src/runtime/sema.go
//go:linkname sync_runtime_Semrelease sync.runtime_Semrelease
func sync_runtime_Semrelease(addr *uint32, handoff bool, skipframes int) {
	semrelease1(addr, handoff, skipframes)
}
```

`runtime_Semrelease` 最终是调用了 `semrelease1`：

```go
func semrelease1(addr *uint32, handoff bool, skipframes int) {
	// 取出信号量对应的 semaRoot
	root := semtable.rootFor(addr)
	// 信号量 +1
	atomic.Xadd(addr, 1)

	// Easy case
	// 没有等待者，直接返回
	if root.nwait.Load() == 0 {
		return
	}

	// Harder case
	lockWithRank(&root.lock, lockRankRoot)
	// 再次检查等待者计数
	if root.nwait.Load() == 0 {
		// 计数已经被其他 goroutine 消费，不需要唤醒其他 goroutine
		unlock(&root.lock)
		return
	}
	// 出队当前信号量上的 sudog
	s, t0, tailtime := root.dequeue(addr)
	if s != nil {
		// 等待者计数 -1
		root.nwait.Add(-1)
	}
	unlock(&root.lock)
	if s != nil { // May be slow or even yield, so unlock first
		// ...
		// 唤醒 goroutine
		readyWithTime(s, 5+skipframes)
		if s.ticket == 1 && getg().m.locks == 0 {
			goyield()
		}
	}
}
```
`goparkunlock` 的实现：

```go
func goparkunlock(lock *mutex, reason waitReason, traceReason traceBlockReason, traceskip int) {
	// 调用 gopark 函数，将 goroutine 阻塞
	gopark(parkunlock_c, unsafe.Pointer(lock), reason, traceReason, traceskip)
}
```

`readyWithTime` 的实现：

```go
func readyWithTime(s *sudog, traceskip int) {
	if s.releasetime != 0 {
		s.releasetime = cputicks()
	}
	// 设置 goroutine 的状态为 runnable 等待被重新调度
	goready(s.g, traceskip)
}
```

## semaphore 扩展库

前面 Go 对信号量的实现都是隐藏在 runtime 中的，并没有标准库来供外部使用。不过 Go 的扩展库 `golang.org/x/sync` 提供了 `semaphore` 包实现的信号量操作。

使用 `func NewWeighted(n int64) *Weighted` 来创建信号量。

`Weighted` 有三个方法：

- `Acquire(ctx contex.Context, n int64) error`：对应 `P` 操作，可以一次获取 n 个资源，如果没有足够多的资源，调用者就会被阻塞。
- `Release(n int64)`：对应 `V` 操作，可以释放 n 个资源。
- `TryAcquire(n int64) bool`：尝试获取 n 个资源，但是它不会阻塞，成功获取 n 个资源则返回 `true`。否则一个也不获取，返回 `false`。


### 使用

```go
var (
    maxWorkers = runtime.GOMAXPROCS(0)                    // worker 数量和 CPU 核数一样
    sema       = semaphore.NewWeighted(int64(maxWorkers)) // 信号量
    task       = make([]int, maxWorkers*4)                // 任务数，是 worker 的四倍
)

func main() {
    ctx := context.Background()

    for i := range task {
        // 如果没有 worker 可用，会阻塞在这里，直到某个 worker 被释放
        if err := sema.Acquire(ctx, 1); err != nil {
            break
        }

        // 启动 worker goroutine
        go func(i int) {
            defer sema.Release(1)
            time.Sleep(100 * time.Millisecond) // 模拟一个耗时操作
            task[i] = i + 1
        }(i)
    }

    // 获取最大计数值的信号量，这样能确保前面的 worker 都执行完
    if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
        log.Printf("获取所有的 worker 失败: %v", err)
    }

    fmt.Println(task)
}
```

### 原理

`Weighted` 是使用互斥锁和 List 实现的，信号量 `semaphore.Weighted` 的结构体：

```go
type Weighted struct {
    size    int64         // 最大资源数
    cur     int64         // 当前已被使用的资源
    mu      sync.Mutex    // 互斥锁，保证并发安全 
    waiters list.List     // 等待者队列
}
```

List 实现了一个等待队列，等待者的通知是通过 channel 实现的。

`Acquire` 实现：

```go
func (s *Weighted) Acquire(ctx context.Context, n int64) error {
    s.mu.Lock()
    // 剩余的资源大于 n，直接返回
    if s.size-s.cur >= n && s.waiters.Len() == 0 {
		// 已被使用的资源 +n
        s.cur += n
        s.mu.Unlock()
        return nil
    }

    // 请求的资源数 n 大于最大的资源数 size
    if n > s.size {
        s.mu.Unlock()
        // 依赖 ctx 的状态返回，否则会一直阻塞
        <-ctx.Done()
        return ctx.Err()
    }
	
	// 走到这里，说明资源不足

    // 把调用者加入到等待队列中
    // 创建一个 ready chan,以便被通知唤醒
    ready := make(chan struct{})
    w := waiter{n: n, ready: ready}
	// 插入到队列尾部，elem 是新插入的元素
    elem := s.waiters.PushBack(w)
    s.mu.Unlock()


    // 阻塞等待，直到 ctx 被取消或者超时，或者被唤醒
    select {
    case <-ctx.Done(): // ctx 被取消或者超时
        err := ctx.Err()
        s.mu.Lock()
        select {
        case <-ready: // 被唤醒了，那么就忽略 ctx 的状态
            err = nil
        default: 
			// s.waiters.Front() 取出队列的第一个 等待者
            isFront := s.waiters.Front() == elem
			// 直接移除当前 等待者
            s.waiters.Remove(elem)
            // 还有资源，通知其它的 等待者
            if isFront && s.size > s.cur {
                s.notifyWaiters()
            }
        }
        s.mu.Unlock()
        return err
    case <-ready: // 被唤醒了
        return nil
    }
}
```

`Release` 的实现：

```go
func (s *Weighted) Release(n int64) {
    s.mu.Lock()
	// 已被使用的资源 -n
    s.cur -= n
    if s.cur < 0 {
        s.mu.Unlock()
        panic("semaphore: released more than held")
    }
	// 唤醒等待队列中等待者
    s.notifyWaiters()
    s.mu.Unlock()
}
```

`notifyWaiters` 就是遍历等待队列中的等待者，如果资源不够，或者等待队列是空的，就返回：

```go
func (s *Weighted) notifyWaiters() {
	for {
		next := s.waiters.Front()
		// 没有等待者了
		if next == nil {
			break // No more waiters blocked.
		}

		w := next.Value.(waiter)
		// 资源不足，退出
		// s.waiters.Front() 是以先入先出的方式取出等待者，如果第一个等待者没有足够的资源，那么队列中的所有等待者都会继续等待
		if s.size-s.cur < w.n {
			break
		}

		// 资源足够
		// 已被使用的资源 +n
		s.cur += w.n
		// 将等待者移出队列
		s.waiters.Remove(next)
		// 关闭 channel，唤醒等待者
		close(w.ready)
	}
}
```
