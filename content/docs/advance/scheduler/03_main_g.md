---
title: 初始化 main goroutine
weight: 3
---

`go` 关键字启动一个 goroutine 时，最终会被编译器转换成 `newproc` 函数：

```go
func newproc(fn *funcval) {
    // 函数调用参数入栈顺序是从右向左，而且栈是从高地址向低地址增长的
	gp := getg()
	pc := sys.GetCallerPC()
	systemstack(func() {
		newg := newproc1(fn, gp, pc, false, waitReasonZero)

		pp := getg().m.p.ptr()
		runqput(pp, newg, true)

		if mainStarted {
			wakep()
		}
	})
}
```

1. `newproc1` 函数的第一个参数 `fn` 是新创建的 goroutine 需要执行的函数；
2. `newproc1` 根据传入参数初始化一个 `g` 结构体。
  
```go
func newproc1(fn *funcval, callergp *g, callerpc uintptr, parked bool, waitreason waitReason) *g {
    // ...
	mp := acquirem() // disable preemption because we hold M and P in local vars.
	pp := mp.p.ptr()
	newg := gfget(pp) // 从 p 的本地缓冲里获取一个没有使用的 g，初始化时没有，返回nil

	if newg == nil {
        // new 一个 g 结构体对象，然后从堆上为其分配栈，并设置 g 的 stack 成员和两个stackgard 成员
		newg = malg(_StackMin)
		casgstatus(newg, _Gidle, _Gdead) // 初始化 g 的状态为 _Gdead
		allgadd(newg) // 放入全局变量 allgs 切片中
	}
	// ...
    // 把 newg.sched 结构体成员的所有成员设置为 0
    memclrNoHeapPointers(unsafe.Pointer(&newg.sched), unsafe.Sizeof(newg.sched))
    
    //设置 newg 的 sched 成员，调度器需要依靠这些字段才能把 goroutine 调度到 CPU 上运行。
    newg.sched.sp = sp // newg 的栈顶
    newg.stktopsp = sp
    // newg.sched.pc 表示当 newg 被调度起来运行时从这个地址开始执行指令
    // 把 pc 设置成了 goexit 这个函数偏移 1（sys.PCQuantum 等于 1）的位置
    newg.sched.pc = funcPC(goexit) + sys.PCQuantum // +PCQuantum so that previous instruction is in same function
    newg.sched.g = guintptr(unsafe.Pointer(newg))

    gostartcallfn(&newg.sched, fn) // 调整 sched 成员和 newg 的栈
}
```

`runtime.gfget` 中包含两部分逻辑，它会根据处理器中 `gFree` 列表中 goroutine 的数量做出不同的决策：

- 当 `p` 的 goroutine 列表为空时，会将 `sched` 调度器持有的空闲 goroutine 转移到当前 `p` 上，直到 `gFree` 列表中的 goroutine 数量达到 32；
- 当 `p` 的 goroutine 数量充足时，会从列表头部返回一个 goroutine；

```go
func gfget(pp *p) *g {
retry:
	if pp.gFree.empty() && (!sched.gFree.stack.empty() || !sched.gFree.noStack.empty()) {
		lock(&sched.gFree.lock)
		// Move a batch of free Gs to the P.
		for pp.gFree.n < 32 {
			// Prefer Gs with stacks.
			gp := sched.gFree.stack.pop()
			if gp == nil {
				gp = sched.gFree.noStack.pop()
				if gp == nil {
					break
				}
			}
			sched.gFree.n--
			pp.gFree.push(gp)
			pp.gFree.n++
		}
		unlock(&sched.gFree.lock)
		goto retry
	}
    gp := pp.gFree.pop()
	if gp == nil {
		return nil
	}
    // ...
}
```

当 `p` 的 `gFree` 和调度器的 `gFree` 列表都不存在结构体时，调用 `runtime.malg` 初始化新的 `g`。

拿到 `g` 之后，调用 `runtime.runqput` 会将 goroutine 放到运行队列上，这既可能是全局的运行队列，也可能是 `p` 本地的运行队列：

1. 当 `next` 为 `true` 时，将 goroutine 设置到处理器的 `runnext `作为下一个处理器执行的任务；
2. 当 `next` 为 `false` 并且本地运行队列还有剩余空间时，将 goroutine 加入处理器持有的本地运行队列；
3. 当处理器的本地运行队列已经没有剩余空间时就会把本地队列中的一部分 goroutine 和待加入的 goroutine 通过 `runtime.runqputslow` 添加到调度器持有的全局运行队列上；


从 `newproc` 继续往下执行 `mstart0`，继续调用 `mstart1` 函数：

```go
func mstart0() {
	gp := getg() // gp = g0

    // 对于启动过程来说，g0 的 stack.lo 早已完成初始化，所以 onStack = false
	osStack := gp.stack.lo == 0
	if osStack {
		// Initialize stack bounds from system stack.
		// Cgo may have left stack size in stack.hi.
		// minit may update the stack bounds.
		//
		// Note: these bounds may not be very accurate.
		// We set hi to &size, but there are things above
		// it. The 1024 is supposed to compensate this,
		// but is somewhat arbitrary.
		size := gp.stack.hi
		if size == 0 {
			size = 16384 * sys.StackGuardMultiplier
		}
		gp.stack.hi = uintptr(noescape(unsafe.Pointer(&size)))
		gp.stack.lo = gp.stack.hi - size + 1024
	}
	// Initialize stack guard so that we can start calling regular
	// Go code.
	gp.stackguard0 = gp.stack.lo + stackGuard
	// This is the g0, so we can also call go:systemstack
	// functions, which check stackguard1.
	gp.stackguard1 = gp.stackguard0


	mstart1()

	// Exit this thread.
	if mStackIsSystemAllocated() {
		// Windows, Solaris, illumos, Darwin, AIX and Plan 9 always system-allocate
		// the stack, but put it in gp.stack before mstart,
		// so the logic above hasn't set osStack yet.
		osStack = true
	}
	mexit(osStack)
}

func mstart1() {
	gp := getg() // gp = g0

	if gp != gp.m.g0 {
		throw("bad runtime·mstart")
	}

	// Set up m.g0.sched as a label returning to just
	// after the mstart1 call in mstart0 above, for use by goexit0 and mcall.
	// We're never coming back to mstart1 after we call schedule,
	// so other calls can reuse the current frame.
	// And goexit0 does a gogo that needs to return from mstart1
	// and let mstart0 exit the thread.
	gp.sched.g = guintptr(unsafe.Pointer(gp))
	gp.sched.pc = sys.GetCallerPC() // 获取 mstart1 执行完的返回地址
	gp.sched.sp = sys.GetCallerSP() // 获取调用 mstart1 时的栈顶地址

	asminit() // 在 AMD64 Linux 平台中，这个函数什么也没做，是个空函数
	minit()

	// Install signal handlers; after minit so that minit can
	// prepare the thread to be able to handle the signals.
	if gp.m == &m0 { //启动时 gp.m 是 m0，所以会执行下面的 mstartm0 函数
		mstartm0()
	}

	if debug.dataindependenttiming == 1 {
		sys.EnableDIT()
	}

	if fn := gp.m.mstartfn; fn != nil { // 初始化过程中 fn == nil
		fn()
	}

	if gp.m != &m0 { // m0 已经绑定了 allp[0]，不是 m0 的话还没有 p，所以需要获取一个 p
		acquirep(gp.m.nextp.ptr())
		gp.m.nextp = 0
	}
	// schedule 函数永远不会返回
	schedule()
}
```

1. `mstart1` 函数先保存 `g0` 的调度信息
2. `GetCallerPC()` 返回的是 `mstart0` 调用 `mstart1` 时被 `call` 指令压栈的返回地址。
3. `GetCallerSP()` 函数返回的是调用 `mstart1` 函数之前 `mstart0` 函数的栈顶地址

所以 **`mstart1` 最主要做的就是保存当前正在运行的 `g` 的下一条指令的地址和栈顶地址**。

不管是对 `g0` 还是其它 `goroutine` 来说这些信息在调度过程中都是必不可少的。

{{< callout type="info" >}}
上面的 `mstart1` 函数中：
- `g0.sched.sp` 指向了 `mstart1` 函数执行完成后的返回地址，该地址保存在了 `mstart0` 函数的栈帧之中；
- `g0.sched.pc` 指向的是 `mstart0` 函数中调用 `mstart1` 函数之后的 `if mStackIsSystemAllocated()` 语句。

从 `mstart0` 函数可以看到，`if mStackIsSystemAllocated()` 语句之后就要退出线程了。为什么要这么做？

原因就在核心函数 `schedule`。
{{< /callout >}}

```go
func schedule() {
	mp := getg().m // getg().m = g0.m, 初始化时 g0.m = m0

    // ...

    // 从本地运行队列和全局运行队列寻找需要运行的 goroutine，
	
	// 为了保证调度的公平性，每进行 61 次调度就需要优先从全局运行队列中获取 goroutine，
    // 因为如果只调度本地队列中的 g，那么全局运行队列中的 goroutine 将得不到运行

    // 如果本地运行队列和全局运行队没有则从其它工作线程的运行队列中偷取，如果偷取不到，则当前工作线程进入睡眠，
    // 直到获取到需要运行的 goroutine 之后 findrunnable 函数才会返回。 
	gp, inheritTime, tryWakeP := findRunnable() // blocks until work is available

	// ...

    // 当前运行的是 runtime 的代码，函数调用栈使用的是 g0 的栈空间
	// 调用 execte 切换到 gp 的代码和栈空间去运行
	execute(gp, inheritTime)
}
```

```go
func execute(gp *g, inheritTime bool) {
	mp := getg().m // getg().m = g0.m, 初始化时 g0.m = m0

	if goroutineProfile.active {
		// Make sure that gp has had its stack written out to the goroutine
		// profile, exactly as it was when the goroutine profiler first stopped
		// the world.
		tryRecordGoroutineProfile(gp, nil, osyield)
	}

	// Assign gp.m before entering _Grunning so running Gs have an
	// M.
	// 把待运行 g 和 m 关联起来
	mp.curg = gp
	gp.m = mp
	// 先设置待运行 g 的状态为 _Grunning
	casgstatus(gp, _Grunnable, _Grunning)
	// ...

    // gogo 完成从 g0 到 gp 真正的切换
	gogo(&gp.sched)
}
```

1. `execute` 函数的第一个参数 `gp` 即是需要调度起来运行的 goroutine，这里首先把 `gp` 的状态从 `_Grunnable` 修改为 `_Grunning`
2. 然后把 `gp` 和 `m` 关联起来，这样通过 `m` 就可以找到当前工作线程正在执行哪个goroutine，反之亦然。
3. 调用 `gogo` 函数完成从 `g0` 到 `gp` 的的切换。

`gogo` 函数是通过汇编语言编写的：

```asm
```

1. `execute` 函数在调用 `gogo` 时把 `gp` 的 `sched` 成员的地址作为实参传递了过来。
2. 