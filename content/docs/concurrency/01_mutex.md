---
title: 互斥锁
weight: 1
---

Go 的标准库 `sync` 提供了两种锁类型：`sync.Mutex` 和 `sync.RWMutex`，前者是互斥锁（排他锁），后者是读写锁。

互斥锁是并发控制的一个基本手段，是为了避免竞争而建立的一种并发控制机制。

Go 定义的锁接口只有两个方法：

```go
type Locker interface {
    Lock() // 请求锁
    Unlock() // 释放锁
}
```

## 使用

```go
import "sync"

var (
	mu sync.Mutex // guards balance 
	balance int
)

func Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()
	balance = balance + amount
}

func Balance() int {
	mu.Lock()
	defer mu.Unlock()
	b := balance
	return b
}
```

当已经有 goroutine 调用 `Lock` 方法获得了这个锁，再有 goroutine 请求这个锁就会阻塞在 `Lock` 方法的调用上，
直到持有这个锁的 goroutine 调用 `UnLock` 释放这个锁。

**使用 `defer` 来 `UnLock` 锁，确保在函数返回之后或者发生错误返回时一定会执行 `UnLock`**。

### 为什么一定要加锁？

```go
import (
    "fmt"
    "sync"
)
    
func main() {
    var count = 0
    // 使用 WaitGroup 等待 10 个 goroutine 完成
    var wg sync.WaitGroup
    wg.Add(10)
    for i := 0; i < 10; i++ {
        go func() {
            defer wg.Done()
            // 对变量 count 执行 10 次加 1
            for j := 0; j < 1000; j++ {
                count++
            }
        }()
    }
    // 等待 10 个 goroutine 完成
    wg.Wait()
    fmt.Println(count)
}
```

上面的例子中期望的最后计数的结果是 `10 * 1000 = 10000`。但是每次运行都可能得到不同的结果，基本上不会得到的一万的结果。

这是因为，`count++` 不是一个原子操作，它至少包含 3 个步骤

1. 读取变量 count 的当前值，
2. 对这个值加 1，
3. 把结果保存到 count 中。

因为不是原子操作，就会有数据竞争的问题。例如，两个 goroutine 同时读取到 count 的值为 8888，接着各自按照自己的逻辑加 1，值变成了 8889，把这个结果再写回到 count 变量。
此时总数只增加了 1，但是应该是增加 2 才对。这是并发访问共享数据的常见问题。

数据竞争的问题可以再编译时通过数据竞争检测器（race detector）工具发现计数器程序的问题以及修复方法。

## 原理

`sync.Mutex` 的结构体：

```go
// src/sync/mutex.go#L34
type Mutex struct {
    state int32
	sema  uint32
}
```

`state` 和 `sema` 加起来占用 8 个字节。

`state` 是一个复合型的字段，包含多个意义：

![mutex-state](https://raw.githubusercontent.com/shipengqi/illustrations/132b0f97ec250e221f725cbeb2d5c323a35f1cfe/go/mutex-state.png)

在默认状态下，互斥锁的所有状态位都是 0，`int32` 中的不同位分别表示了不同的状态：

- `locked`：表示这个锁是否被持有
- `woken`：表示是否从有唤醒的 goroutine
- `starving`：表示此锁是否进入饥饿状态
- `waitersCount`：表示等待此锁的 goroutine 的数量

### 饥饿模式

请求锁的 goroutine 有两类，一类是新来请求锁的 goroutine，另一类是被唤醒的等待请求锁的 goroutine。

由于新来的 goroutine 也参与竞争锁，极端情况下，等待中的 goroutine 可能一直获取不到锁，这就是**饥饿问题**。

为了解决饥饿，Go 1.9 中为 mutex 增加了**饥饿模式**。

在正常模式下，等待中的 goroutine 会按照先进先出的顺序获取锁。但是如果新来的 goroutine 竞争锁，等待中的 goroutine 大概率是获取不到锁的。一旦 goroutine 超
过 1ms 没有获取到锁，它就会将当前互斥锁切换到饥饿模式，保证锁的公平性。

在饥饿模式中，互斥锁会直接交给等待队列最前面的 goroutine。新来的 goroutine 在该状态下不能获取锁、也不会进入自旋状态，只会在队列的末尾等待。

下面两种情况，mutex 会切换为正常模式:

- 一个 goroutine 获得了锁并且它在队列的末尾
- 一个 goroutine 等待的时间少于 1ms

### Lock

`Lock` 的实现：

```go
const (
    mutexLocked = 1 << iota // 1
    mutexWoken // 2
    mutexStarving // 4
    mutexWaiterShift = iota // 3
	starvationThresholdNs = 1e6 // 1000000
)

func (m *Mutex) Lock() {
    // Fast path: grab unlocked mutex.
	// 没有 goroutine 持有锁，也没有等待的 goroutine，当前 goroutine 可以直接获得锁
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        if race.Enabled {
            race.Acquire(unsafe.Pointer(m))
        }
	    return
	}
    // Slow path (outlined so that the fast path can be inlined)
	// 通过自旋等方式竞争锁
    m.lockSlow()
}

func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false // 当前 goroutine 的饥饿标记
	awoke := false // 唤醒标记
	iter := 0 // 自旋次数
	old := m.state // 当前锁的状态
	for {
		// 锁是非饥饿模式并且还没被释放，尝试自旋
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// 尝试设置 mutexWoken 标志来通知解锁，以避免唤醒其他阻塞的 goroutine
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 && 
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state // 再次获取锁的状态，后面会检查锁是否被释放了
			continue
		}
        new := old
        if old&mutexStarving == 0 {
			new |= mutexLocked // 非饥饿状态，加锁
		}
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift // waiter 数量加 1
		}
        if starving && old&mutexLocked != 0 {
			new |= mutexStarving // 设置饥饿状态
		}
		if awoke {
			// The goroutine has been woken from sleep, 
			// so we need to reset the flag in either case. 
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken // 新状态清除唤醒标记
		}
        // 设置新状态
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
            // 再次检查，原来锁的状态已释放，并且不是饥饿状态，正常请求到了锁，返回
			if old&(mutexLocked|mutexStarving) == 0 {
				break // locked the mutex with CAS
			}
			// 处理饥饿状态
			// 如果之前就在该队列里面，就加入到队列头
			queueLifo : waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
            // runtime_SemacquireMutex 通过信号量保证资源不会被两个 goroutine 获取
			// runtime_SemacquireMutex 会在方法中不断尝试获取锁并陷入休眠等待信号量的释放
			// 也就是这里会阻塞等待
			// 一旦当前 goroutine 可以获取信号量，它就会立刻返回，剩余代码也会继续执行
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)
			// 在正常模式下，这段代码会设置唤醒和饥饿标记、重置迭代次数并重新执行获取锁的循环
            // 在饥饿模式下，当前 goroutine 会获得锁，如果等待队列中只存在当前 goroutine，锁还会从饥饿模式中退出
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			if old&mutexStarving != 0 {
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				if !starving || old>>mutexWaiterShift == 1 {
					delta -= mutexStarving
				}
				atomic.AddInt32(&m.state, delta)
				break
			}
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}
	if race.Enabled {
		race.Acquire(unsafe.Pointer(m))
	}	
}
```

#### 自旋

自旋是一种多线程同步机制，**当前的进程在进入自旋的过程中会一直保持 CPU 的占用**，持续检查某个条件是否为真。在多核的 CPU 上，自旋可以避免 goroutine 的切换，使用恰当
会对性能带来很大的增益，但是使用的不恰当就会拖慢整个程序，所以 goroutine 进入自旋的条件非常苛刻：

1. `old&(mutexLocked|mutexStarving) == mutexLocked` 只有在普通模式
2. `runtime_canSpin(iter)` 为真：
   - 运行在多 CPU 的机器上
   - 自旋的次数小于四次
   - 当前机器上至少存在一个正在运行的处理器 P 并且处理的运行队列为空

进入自旋会调用 `runtime_doSpin()`，并执行 30 次的 PAUSE 指令，该指令只会占用 CPU 并消耗 CPU 时间：

```go
//go:linkname sync_runtime_doSpin sync.runtime_doSpin
//go:nosplit
func sync_runtime_doSpin() {
	procyield(active_spin_cnt)
}

TEXT runtime·procyield(SB),NOSPLIT,$0-0
    MOVL	cycles+0(FP), AX
again:
    PAUSE
    SUBL	$1, AX
    JNZ	again
    RET
```
### Unlock

```go
func (m *Mutex) Unlock() {
	if race.Enabled {
		_ = m.state
		race.Release(unsafe.Pointer(m))
	}

	// Fast path: drop lock bit.
	// new == 0 成功释放锁
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 {
		// Outlined slow path to allow inlining the fast path.
		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
		m.unlockSlow(new)
	}
}

func (m *Mutex) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 { // unlock 一个未加锁的锁
		fatal("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 { // 正常模式
		old := new
		for {
			// 不存在等待者 或者 mutexLocked、mutexStarving、mutexWoken 状态不都为 0
			// 则不需要唤醒其他等待者
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}
			// 存在等待者，通过 runtime_Semrelease 唤醒等待者并移交锁的所有权
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		}
	} else { // 饥饿模式
		// 直接调用 runtime_Semrelease 将当前锁交给下一个正在尝试获取锁的等待者，等待者被唤醒后会得到锁，在这时还不会退出饥饿状态
		runtime_Semrelease(&m.sema, true, 1)
	}
}
```
