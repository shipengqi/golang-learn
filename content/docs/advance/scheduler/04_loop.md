---
title: 调度循环
weight: 4
---

## 非 main goroutine 的退出

`newproc1` 创建 goroutine 的时候已经在其栈上放好了一个返回地址，伪造成 `goexit` 函数调用了 goroutine 的入口函数。**非 main goroutine 执行完成后就会返回到 `goexit` 继续执行**。

`goexit` 函数在 `runtime/asm_amd64.s` 文件中：

```asm
TEXT runtime·goexit(SB),NOSPLIT|TOPFRAME|NOFRAME,$0-0
	BYTE	$0x90	// NOP
	CALL	runtime·goexit1(SB)	// does not return
	// traceback from goexit1 must hit code range of goexit
	BYTE	$0x90	// NOP
```

`CALL runtime·goexit1(SB)` 继续调用 `goexit1` 函数，`goexit1` 函数又调用 `mcall(goexit0)`。

**`mcall` 做的事情跟 `gogo` 函数完全相反。`gogo` 函数实现了从 `g0` 切换到某个 goroutine 去运行，而 `mcall` 实现了从某个 goroutine 切换到 `g0` 来运行**。

切换到 `g0` 栈之后，下面开始在 `g0` 栈执行 `goexit0` 函数，该函数完成最后的清理工作：

1. 把 `g` 的状态从 `_Grunning` 变更为 `_Gdead`；
2. 然后把 `g` 的一些字段清空成 0 值；
3. 调用 `dropg` 函数解除 `g` 和 `m` 之间的关系，其实就是设置 `g->m = nil, m->currg = nil`；
4. 把 `g` 放入 `p` 的 `freeg` 队列缓存起来供下次创建 `g` 时快速获取而不用从内存分配。`freeg` 就是 `g` 的一个对象池；
5. **调用 `schedule` 函数再次进行调度**；

工作线程再次调用了 `schedule` 函数进入新一轮的调度循环。

```go
func goexit0(gp *g) {
	gdestroy(gp)
	schedule()
}
```

调用链：

```
schedule() -> execute() -> gogo() -> g2() -> goexit() -> goexit1() -> mcall() -> goexit0() -> schedule()
```

## 调度策略

1. 从全局运行队列中寻找 goroutine。为了保证调度的公平性，每个工作线程每经过 61 次调度就需要优先尝试从全局运行队列中找出一个 goroutine 来运行，这样才能保证位于全局运行队列中的 goroutine 得到调度的机会。全局运行队列是所有工作线程都可以访问的，所以在访问它之前需要加锁。
2. 从工作线程本地运行队列中寻找 goroutine。如果不需要或不能从全局运行队列中获取到 goroutine 则从本地运行队列中获取。
3. 尝试通过 netpoll 快速获取 I/O 就绪任务
4. 从其它工作线程的运行队列中偷取 goroutine。如果上一步也没有找到需要运行的 goroutine，则从其他工作线程的运行队列中偷取 goroutine，在偷取之前会再次尝试从全局运行队列和当前线程的本地运行队列中查找需要运行的 goroutine。

```go
func findRunnable() (gp *g, inheritTime, tryWakeP bool) {
	mp := getg().m

    // ...
    
    // Check the global runnable queue once in a while to ensure fairness.
	// Otherwise two goroutines can completely occupy the local runqueue
	// by constantly respawning each other.
    // 为了保证调度的公平性，每进行 61 次调度就需要优先从全局运行队列中获取 goroutine，
    // 因为如果只调度本地队列中的 g，那么全局运行队列中的 goroutine 将得不到运行
	if pp.schedtick%61 == 0 && sched.runqsize > 0 {
		lock(&sched.lock)  // 所有工作线程都能访问全局运行队列，所以需要加锁
		gp := globrunqget(pp, 1)
		unlock(&sched.lock)
		if gp != nil {
			return gp, false, false
		}
	}

    // ...

    // local runq
    // 从与 m 关联的 p 的本地运行队列中获取 goroutine
	if gp, inheritTime := runqget(pp); gp != nil {
		return gp, inheritTime, false
	}

	// global runq
    // 从全局运行队列中获取 goroutine
	if sched.runqsize != 0 {
		lock(&sched.lock)
		gp := globrunqget(pp, 0)
		unlock(&sched.lock)
		if gp != nil {
			return gp, false, false
		}
	}

	// Poll network.
	// This netpoll is only an optimization before we resort to stealing.
	// We can safely skip it if there are no waiters or a thread is blocked
	// in netpoll already. If there is any kind of logical race with that
	// blocked thread (e.g. it has already returned from netpoll, but does
	// not set lastpoll yet), this thread will do blocking netpoll below
	// anyway.
    // 这里是在偷取 goroutine 之前的额一个优化。尝试通过 netpoll 快速获取 I/O 就绪任务
    // 如果系统中已经有线程在处理 netpoll，就可以跳过这一步
    if netpollinited() && netpollAnyWaiters() && sched.lastpoll.Load() != 0 {
        // ...
	}

    // Steal work from other P's.
    // If number of spinning M's >= number of busy P's, block.
    // This is necessary to prevent excessive CPU consumption
    // when GOMAXPROCS>>1 but the program parallelism is low.
    // 这个判断主要是为了防止因为寻找可运行的 goroutine 而消耗太多的 CPU。
    // 因为已经有足够多的工作线程正在寻找可运行的 goroutine，让他们去找就好了，自己偷个懒去睡觉
    if mp.spinning || 2*sched.nmspinning.Load() < gomaxprocs-sched.npidle.Load() {
        if !mp.spinning {
            mp.becomeSpinning() // 设置 m 的状态为 spinning
        }

        gp, inheritTime, tnow, w, newWork := stealWork(now) // 从其它 p 的本地运行队列盗取 goroutine
        // ...
    }
}
```

对于多个线程同时窃取同一个 P 的本地队列的情况，只有一个线程能窃取成功，其他线程只能继续从全局队列或者当前线程的本地队列中查找。

这里使用的 `for` 循环加原子操作 CAS （`atomic.CasRel`）来保证只有一个线程能窃取成功。`atomic.CasRel(&pp.runqhead, h, h+n)` 中 `runqhead` 是本地丢列的头指针。

## 调度时机

触发调度的几个路径：

- 主动挂起 — `runtime.gopark` -> `runtime.park_m`。
- 系统调用 — `runtime.exitsyscall` -> `runtime.exitsyscall0`。
- 协作式调度 — `runtime.Gosched` -> `runtime.gosched_m` -> `runtime.goschedImpl`。
- 系统监控 — `runtime.sysmon` -> `runtime.retake` -> `runtime.preemptone`。

### 主动挂起

`runtime.gopark` 会通过 `runtime.mcall` 切换到 `g0` 的栈上调用 `runtime.park_m`：

```go
func park_m(gp *g) {
	// ...
	casgstatus(gp, _Grunning, _Gwaiting)
	dropg()

	schedule()
}
```

1. 将当前 goroutine 的状态从 `_Grunning` 切换至 `_Gwaiting`。
2. 调用 `runtime.dropg` 移除 `m` 和 `g` 之间的关联。
3. 调用 `runtime.schedule` 触发新一轮的调度。

当 goroutine 等待的特定条件满足后，运行时会调用 `runtime.goready` 将因为调用 `runtime.gopark` 而陷入休眠的 goroutine 唤醒。

```go
func goready(gp *g, traceskip int) {
	systemstack(func() {
		ready(gp, traceskip, true)
	})
}

func ready(gp *g, traceskip int, next bool) {
	// ...
	mp := acquirem() // disable preemption because it can be holding p in a local var
    // ...
	casgstatus(gp, _Gwaiting, _Grunnable)
	runqput(mp.p.ptr(), gp, next)
	if atomic.Load(&sched.npidle) != 0 && atomic.Load(&sched.nmspinning) == 0 {
		wakep()
	}
	// ...
}
```

1. 将 goroutine 的 `_Gwaiting` 状态切换至 `_Grunnable`。
2. 将其加入处理器的运行队列中，等待调度器的调度。


{{< callout type="info" >}}
`gopark` 需要调用 `schedule` 而 `goready` 不需要，原因：

- `gopark` 将 `g` 从 `Grunning` 变为 `Gwaiting`，必须让出 `m`，找新 `g` 来运行。
- `goready` 将 `g` 从 `Gwaiting` 变为 `Grunnable`，只需将 `g` 放入 `runq` 队列即可。

正常结束的非 main goroutine 会返回到 `goexit` 函数，切换到 `g0` 继续执行 `shcedule`。

**`gopark`（`mcall`）和 `goready`（`systemstack`）都会切换到 `g0` 栈去执行**。
{{< /callout >}}

#### 使用场景


- channel 阻塞（`hchan.sendq` 向 channel 发送数据而被阻塞的 goroutine 队列，`hchan.recvq` 读取 channel 的数据而被阻塞的 goroutine 队列）-> `gopark/goready`。
- `sync.Metux` -> 信号量（`semaRoot.treap` 等待者队列）-> `gopark/goready`。
- `sync.WaitGroup` -> 信号量（`semaRoot.treap` 等待者队列）-> `gopark/goready`。
- `sync.Cond` -> `gopark/goready`。
- `golang.org/x/sync/semaphore` -> channel 阻塞、通知。
- `golang.org/x/sync/singleflight` -> `sync.Metux` -> 信号量（`semaRoot.treap` 等待者队列）-> `gopark/goready`。
- `golang.org/x/sync/errgroup` -> `sync.WaitGroup` -> 信号量（`semaRoot.treap` 等待者队列）-> `gopark/goready`。


上面的几种方式，都有一个被阻塞的 goroutine 队列， `goready` 唤醒时，可以直接使用阻塞队列中的 `g` 对象。

### 系统调用

系统调用也会触发运行时调度器的调度，goroutine 有一个 `_Gsyscall` 状态用来表示系统调用。

Go 通过汇编语言封装了系统调用：

```asm
#define INVOKE_SYSCALL	INT	$0x80

TEXT ·Syscall(SB),NOSPLIT,$0-28
	CALL	runtime·entersyscall(SB)
	...
	INVOKE_SYSCALL
	...
	CALL	runtime·exitsyscall(SB)
	RET
ok:
	...
	CALL	runtime·exitsyscall(SB)
	RET
```

1. `runtime.entersyscall` 完成 goroutine 进入系统调用前的准备工作。
2. `INVOKE_SYSCALL` 系统调用指令。
3. `runtime.exitsyscall` 为当前 goroutine 重新分配资源。
4. 释放当前 `m` 上的锁，**锁被释放后，当前线程会陷入系统调用等待返回**，在锁被释放后，**会有其他 goroutine 抢占 `p`**（这是后面 `exitsyscall` 会有两种路径的原因）。

#### runtime.entersyscall

`runtime.entersyscall` 主要做以下几件事：

1. 保存当前 goroutine 的上下文信息，程序计数器 PC 和栈指针 SP 中的内容。
2. 切换当前 goroutine 为 `_Gsyscall` 状态。
3. 将 goroutine 的 `p` 和 `m` 暂时分离并更新 `p` 的状态到 `_Psyscall`；

{{< callout type="info" >}}
这里的当前 goroutine 并没有和 `m` 解绑，只是 `p` 和 `m` 解绑。当前 goroutine 的保存上下文信息是执行系统调用前的 PC 和 SP 等。

然后 `m` 陷入了阻塞，等待系统调用返回。

返回之后才会将当前 goroutine 切换至 `_Grunnable` 状态，并移除 `m` 和当前 goroutine 的关联，放入运行队列，触发 `runtime.schedule` 调度。
{{< /callout >}}

#### runtime.exitsyscall

系统调用结束后，会调用退出系统调用的函数 `runtime.exitsyscall` 为当前 goroutine 重新分配资源，该函数有两个不同的执行路径：

1. 调用 `runtime.exitsyscallfast`；
2. 切换至 `g0` 并调用 `runtime.exitsyscall0`，将当前 goroutine 切换至 `_Grunnable` 状态；

对于当前 goroutine 放入哪个运行队列有两种策略：

1. 如果当前 goroutine 的执行系统调用前就绑定的 `p` 仍处于 `_Psyscall` 状态，会直接调用 `wirep` 将 goroutine 与处理器进行关联；
2. 如果调度器中存在闲置的 `p`，会调用 `runtime.acquirep` 使用闲置的 `p` 处理当前 goroutine；

**最后都会调用 `runtime.schedule` 触发调度器的调度**。

### 协作式调度

**`runtime.Gosched` 函数会主动让出处理器**，允许其他 goroutine 运行。**该函数无法挂起 goroutine，调度器可能会将当前 goroutine 调度到其他线程上**。

```go
func Gosched() {
	checkTimeouts()
	mcall(gosched_m)
}

func gosched_m(gp *g) {
	goschedImpl(gp)
}

func goschedImpl(gp *g) {
	casgstatus(gp, _Grunning, _Grunnable)
	dropg()
	lock(&sched.lock)
	globrunqput(gp)
	unlock(&sched.lock)

	schedule()
}
```

经过连续几次跳转，最终在 `g0` 的栈上调用 `runtime.goschedImpl`：

1. 运行时会更新 goroutine 的状态到 `_Grunnable`。
2. 让出当前的处理器并将 goroutine 重新放回全局队列。
3. 在最后，该函数会调用 `runtime.schedule` 触发调度。

### 总结

goroutine 的调度，总体就是一个循环，伪代码：

```go
// --------------------------------
// 线程部分

// 定义一个线程私有全局变量，注意它是一个指向 m 结构体对象的指针
// ThreadLocal 用来定义线程私有全局变量
ThreadLocal self *m

// schedule 函数实现调度逻辑
schedule() {
    // 创建和初始化 m 结构体对象，并赋值给私有全局变量 self
    self = initm()   
    for { // 调度循环
        g = find_a_runnable_goroutine_from_local_runqueue()
        run_g(g) // CPU 运行该 goroutine，直到需要调度其它 goroutine 才返回
        save_status_of_g(g) // 保存 goroutine 的状态，主要是寄存器的值
     }
}
```

- 正常执行结束的 goroutine，会返回到 `goexit` 函数，然后切换到 `g0` 栈继续执行 `schedule` 函数。
- 调用 `gopark` 的 goroutine，将状态设置为 `_Gwaiting`，然后切换到 `g0` 栈继续执行 `schedule` 函数。当前 goroutine 会放到某个队列中，方便 `goready` 时唤醒。唤醒时将状态设置为 `_Grunnable`，并放入可运行队列。
- 调用 `Gosched` 函数会让出处理器并将 goroutine 重新放回全局队列。状态仍然是 `_Grunnable`。可能会被调度到其他的 `p`。然后切换到 `g0` 栈继续执行 `schedule` 函数。
- 执行系统调用的 goroutine，将状态设置为 `_Gsyscall`。当前 goroutine 仍然和 `m` 绑定，`m` 被阻塞，系统调用返回时，将状态设置为 `_Grunnable`，并将 goroutine 放到可运行队列。然后切换到 `g0` 栈继续执行 `schedule` 函数。
- 被抢占调度的 goroutine，将状态设置为 `_Grunnable`，并放入可运行队列。然后切换到 `g0` 栈继续执行 `schedule` 函数。

## 线程管理

`runtime.LockOSThread` 和 `runtime.UnlockOSThread` 可以绑定 goroutine 和线程完成一些比较特殊的操作。

```go
func LockOSThread() {
	if atomic.Load(&newmHandoff.haveTemplateThread) == 0 && GOOS != "plan9" {
		startTemplateThread()
	}
	_g_ := getg()
	_g_.m.lockedExt++
	dolockOSThread()
}

func dolockOSThread() {
	_g_ := getg()
	_g_.m.lockedg.set(_g_)
	_g_.lockedm.set(_g_.m)
}
```

`runtime.dolockOSThread` 会分别设置线程的 `lockedg` 字段和 goroutine 的 `lockedm` 字段，这两行代码会绑定线程和 goroutine。

`runtime.UnlockOSThread` 用户解绑 goroutine 和线程。

### 查看 goroutine 数量

可以使用 `runtime.NumGoroutine` 函数查看当前 goroutine 的数量。

### 线程生命周期

Go 语言的运行时会通过 `runtime.startm` 启动线程来执行处理器 `p`，如果在该函数中没能从闲置列表中获取到线程 `m` 就会调用 `runtime.newm` 创建新的线程：

```go
func newm(fn func(), _p_ *p, id int64) {
	mp := allocm(_p_, fn, id)
	mp.nextp.set(_p_)
	mp.sigmask = initSigmask
	...
	newm1(mp)
}

func newm1(mp *m) {
	if iscgo {
		...
	}
	newosproc(mp)
}
```

创建新的线程需要使用如下所示的 `runtime.newosproc`，该函数在 Linux 平台上会通过系统调用 clone 创建新的操作系统线程：

```go
func newosproc(mp *m) {
	stk := unsafe.Pointer(mp.g0.stack.hi)
	...
	ret := clone(cloneFlags, stk, unsafe.Pointer(mp), unsafe.Pointer(mp.g0), unsafe.Pointer(funcPC(mstart)))
	...
}
```