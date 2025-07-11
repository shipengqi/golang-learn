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

![mutex-state](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/mutex-state.png)

在默认状态下，互斥锁的所有状态位都是 0，`int32` 中的不同位分别表示了不同的状态：

- `locked`：表示互斥锁的锁定状态 (1: 锁定, 0: 未锁定)；
- `woken`：表示是否有被唤醒的 goroutine；
- `starving`：表示当前锁是否进入饥饿状态；
- `waitersCount`：表示等待当前锁的 goroutine 的数量（这个在 Fast path 的判断中会用到）；

### 正常模式和饥饿模式

`sync.Mutex` 有两种模式：**正常模式**和**饥饿模式**。

- 正常模式下，锁的等待者会按照**先进先出**的顺序获取锁。
- Go 1.9 中为 mutex 增加了**饥饿模式**。饥饿模式是指，刚被唤起的 goroutine 与新创建的 goroutine 竞争时，大概率会获取不到锁。为了减少这种情况，一旦 **goroutine 超过 1ms 没有获取到锁，它就会将当前互斥锁切换到饥饿模式**，保证锁的公平性。

**在饥饿模式中，互斥锁会直接交给等待队列最前面的 goroutine**。新来的 goroutine 在该状态下不能获取锁、也不会进入自旋状态，只会在队列的末尾等待。

下面两种情况，mutex 会切换为正常模式:

- 一个 goroutine 获得了锁并且它在队列的末尾。
- 一个 goroutine 获得了锁并且等待的时间少于 1ms。

### Lock

`Lock` 的实现：

```go
const (
    mutexLocked = 1 << iota // 1 (二进制: 0001)
    mutexWoken // 2 (二进制: 0010)
    mutexStarving // 4 (二进制: 0100)
    mutexWaiterShift = iota
	starvationThresholdNs = 1e6 // 1000000 (进入饥饿模式的阈值)
)

func (m *Mutex) Lock() {
    // Fast path: grab unlocked mutex.
	// 锁的状态是 0，没有 goroutine 持有锁，也没有等待的 goroutine，当前 goroutine 可以直接获得锁，设置为 mutexLocked
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        if race.Enabled {
            race.Acquire(unsafe.Pointer(m))
        }
	    return
	}
    // Slow path (outlined so that the fast path can be inlined)
	// 互斥锁的状态不是 0，尝试通过自旋（Spinnig）或信号量阻塞等方式等待锁的释放
    m.lockSlow()
}

func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false // 当前 goroutine 的饥饿标记
	awoke := false // 唤醒标记
	iter := 0 // 自旋次数
	old := m.state // 当前锁的状态
	for {
		// 情况1: 锁是非饥饿模式并且还没被释放，且可以自旋
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// 尝试设置唤醒标志
            // 条件:
            // 1. 当前 goroutine 还未被唤醒 (!awoke)
            // 2. 锁没有设置唤醒标志 (old&mutexWoken == 0)
            // 3. 存在等待者 (old>>mutexWaiterShift != 0)
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 && 
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
			// 执行自旋，执行 30 次 PAUSE 指令
			runtime_doSpin()
			iter++
			old = m.state // 重新加载锁的状态
			continue
		}
		// 情况2: 准备新状态（没有进入自旋）
        new := old
		// 如果不是饥饿模式，表示当前 goroutine 要获取锁了，这里还没有真正修改锁的状态
        if old&mutexStarving == 0 {
			new |= mutexLocked 
		}
		// 如果锁已被持有或处于饥饿模式，增加等待计数，waiter 加 1
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}
		// 如果当前 goroutine 已处于饥饿状态且锁仍被持有
        // 则设置饥饿模式标志
        if starving && old&mutexLocked != 0 {
			new |= mutexStarving
		}
		// 如果当前 goroutine 是被唤醒的
		if awoke {
			// The goroutine has been woken from sleep, 
			// so we need to reset the flag in either case. 
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken // 清除唤醒标记 (因为当前 goroutine 要么获取锁，要么再次休眠)
		}
        // 尝试去更新锁的新状态
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			// CAS 成功并不代表获取了锁，而是代表成功更新了锁的状态，CAS 成功说明当前 goroutine 是获取锁的竞争者
			// 更新了锁的状态之后，具体分为两种情况

            // 情况1：
			// 如果锁原先未被持有，并且不是饥饿模式，成功获取锁，直接返回
			if old&(mutexLocked|mutexStarving) == 0 {
				break // 通过 CAS 函数获取了锁
			}
			
		    // ...

            // 情况2：需要排队等待
			// 真正的锁获取需要等待前一个持有者释放
            // runtime_SemacquireMutex 通过信号量保证资源不会被两个 goroutine 获取
			// runtime_SemacquireMutex 会在方法中不断尝试获取锁并陷入休眠等待信号量的释放
			// 也就是在这里会阻塞等待
			// 一旦当前 goroutine 可以获取信号量，它就会立刻返回，剩余代码也会继续执行，尝试获取锁
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)

            // 走到这里说明锁被释放了，要开始重新尝试获取锁了：
			// 在正常模式下，这段代码会设置唤醒和饥饿标记、重置迭代次数并重新执行获取锁的循环，超过 1ms 会将当前 goroutine 设置为饥饿状态
            // 在饥饿模式下，当前 goroutine 快要饿死了，直接会获得锁，如果等待队列中只存在当前 goroutine，锁还会从饥饿模式中退出

			// 判断是否应该进入饥饿状态
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			// 如果当前锁已经处于饥饿模式
			if old&mutexStarving != 0 {
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				// 设置 locked 位(1)
                // 减少一个等待者(-1<<3)
				delta := int32(mutexLocked - 1<<mutexWaiterShift)

				// 如果当前 goroutine 不是饥饿状态
                // 或者它是最后一个等待者
                // 则退出饥饿模式
				if !starving || old>>mutexWaiterShift == 1 {
					delta -= mutexStarving // 清除饥饿标志
				}
				atomic.AddInt32(&m.state, delta)
				break
			}
			
			// 不是饥饿模式则继续尝试获取锁
			// 下一轮循环大概率会获取到锁，因为现在处于非饥饿模式，下一轮循环也不会饥饿（饥饿模式是在获取信号量后面执行的），锁也被释放了
			// 而且信号量唤醒的是等待队列头部的 g
			awoke = true
			iter = 0
		} else {
			// CAS 失败，重新加载状态，进入下一轮循环
			old = m.state
		}
	}
	if race.Enabled {
		race.Acquire(unsafe.Pointer(m))
	}	
}
```

#### 自旋

自旋是一种多线程同步机制，**当前的进程在进入自旋的过程中会一直保持 CPU 的占用**，持续检查某个条件是否为真。在多核的 CPU 上，自旋可以避免 goroutine 的切换，使用恰当会对性能带来很大的增益，但是使用的不恰当就会拖慢整个程序，所以 goroutine 进入自旋的条件非常苛刻：

1. `old&(mutexLocked|mutexStarving) == mutexLocked` 只有在普通模式下才能进入自旋；
2. `runtime_canSpin(iter)` 为真：
   - 运行在多 CPU 的机器上
   - 自旋的次数小于四次
   - 当前机器上至少存在一个正在运行的处理器 P 并且处理的运行队列为空

进入自旋会调用 `runtime_doSpin()`，并执行 30 次的 `PAUSE` 指令，该指令只会占用 CPU 并消耗 CPU 时间：

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
	// 快速解锁，new == 0 成功释放锁
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 { // 意味着当前锁有等待者，需要唤醒等待者
		// Outlined slow path to allow inlining the fast path.
		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
		// 慢速解锁
		m.unlockSlow(new)
	}
}

func (m *Mutex) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 { // 如果当前互斥锁已经被解锁过了会直接抛出异常
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
