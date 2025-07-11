---
title: 初始化 main goroutine
weight: 3
---

## 创建 main goroutine

`schedinit` 完成调度系统初始化后，返回到 `rt0_go` 函数中开始调用 `newproc()` 创建一个新的 goroutine 用于执行 mainPC 所对应的 `runtime·main` 函数。

```asm
// ...
CALL	runtime·schedinit(SB)

// create a new goroutine to start program
MOVQ	$runtime·mainPC(SB), AX		// entry
PUSHQ	AX
CALL	runtime·newproc(SB)
POPQ	AX

// start this M
CALL	runtime·mstart(SB)

CALL	runtime·abort(SB)	// mstart should never return
RET
```

另外 `go` 关键字启动一个 goroutine 时，最终也会被编译器转换成 `newproc` 函数。

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
        // new 一个 g 结构体对象，然后从堆上为其分配栈，并设置 g 的 stack 成员和两个 stackgard 成员
		newg = malg(_StackMin)
		casgstatus(newg, _Gidle, _Gdead) // 初始化 g 的状态为 _Gdead
		allgadd(newg) // 放入全局变量 allgs 切片中
	}
	// ...
    // 把 newg.sched 结构体成员的所有成员设置为 0
    memclrNoHeapPointers(unsafe.Pointer(&newg.sched), unsafe.Sizeof(newg.sched))
    
    // 设置 newg 的 sched 成员，调度器需要依靠这些字段才能把 goroutine 调度到 CPU 上运行。
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

- 当 `p` 的 `gfree` 数量充足时，会从列表头部返回一个 goroutine；
- 当 `p` 的 `gfree` 列表为空时，会将调度器持有的空闲 goroutine 转移到当前 `p` 上，直到 `gfree` 列表中的 goroutine 数量达到 32；

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

**当 `p` 的 `gfree` 和调度器的 `gFree` 列表都不存在结构体时，调用 `runtime.malg` 初始化新的 `g`**。

拿到 `g` 之后，调用 `runtime.runqput` 会将 goroutine 放到运行队列 `runq` 上，这既可能是全局的运行队列，也可能是 `p` 本地的运行队列：

1. 当 `next` 为 `true` 时，将 goroutine 设置到处理器的 `runnext `作为下一个处理器执行的任务；
2. 当 `next` 为 `false` 并且本地运行队列还有剩余空间时，将 goroutine 加入处理器持有的本地运行队列；
3. 当 `p` 的**本地运行队列已经没有剩余空间时就会把本地队列中的一部分 goroutine 和待加入的 goroutine 通过 `runtime.runqputslow` 添加到调度器持有的全局运行队列上**；

## 从 g0 切换到 main goroutine

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

	if gp.m != &m0 { // m0 已经绑定了 allp[0]，如果不是 m0 的话，这时还没有 p，所以需要获取一个 p
		acquirep(gp.m.nextp.ptr())
		gp.m.nextp = 0
	}
	// schedule 函数永远不会返回
	schedule()
}
```

1. `mstart1` 函数先保存 `g0` 的调度信息。
2. `GetCallerPC()` 返回的是 `mstart0` 调用 `mstart1` 时被 `call` 指令压栈的返回地址。
3. `GetCallerSP()` 函数返回的是调用 `mstart1` 函数之前 `mstart0` 函数的栈顶地址。

所以 **`mstart1` 最主要做的就是保存当前正在运行的 `g` 的下一条指令的地址和栈顶地址**。

不管是对 `g0` 还是其它 `goroutine` 来说这些信息在调度过程中都是必不可少的。

{{< callout type="info" >}}
上面的 `mstart1` 函数中：

- `g0.sched.pc` 指向的是 `mstart0` 函数中调用 `mstart1` 函数之后下一个指令（也就是 `if mStackIsSystemAllocated()` 语句）的地址。

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
2. 然后把 `gp` 和 `m` 关联起来，这样通过 `m` 就可以找到当前工作线程正在执行哪个 goroutine，反之亦然。
3. 调用 `gogo` 函数完成从 `g0` 到 `gp` 的的切换。

`gogo` 函数是通过汇编语言编写的：

```asm
TEXT gogo<>(SB), NOSPLIT, $0
	get_tls(CX)
	// 把要运行的 g 的指针放入线程本地存储，这样后面的代码就可以通过线程本地存储
	// 获取到当前正在执行的 goroutine 的 g 结构体对象，从而找到与之关联的 m 和 p
	MOVQ	DX, g(CX)
	MOVQ	DX, R14		// set the g register
	// 把 CPU 的 SP 寄存器设置为 sched.sp，完成了栈的切换
	MOVQ	gobuf_sp(BX), SP	// restore SP
	// 恢复调度上下文到 CPU 相关寄存器
	MOVQ	gobuf_ret(BX), AX
	MOVQ	gobuf_ctxt(BX), DX
	MOVQ	gobuf_bp(BX), BP
	// 清空 sched 的值，因为我们已把相关值放入 CPU 对应的寄存器了，不再需要，这样做可以少 gc 的工作量
	MOVQ	$0, gobuf_sp(BX)	// clear to help garbage collector
	MOVQ	$0, gobuf_ret(BX)
	MOVQ	$0, gobuf_ctxt(BX)
	MOVQ	$0, gobuf_bp(BX)
	// 把 sched.pc 值放入 BX 寄存器
	MOVQ	gobuf_pc(BX), BX
	// JMP 把 BX 寄存器的包含的地址值放入 CPU 的 IP 寄存器，于是，CPU 跳转到该地址继续执行指令
	JMP	BX
```

`gogo` 函数就只做了两件事：

1. 把 `gp.sched` 的成员恢复到 CPU 的寄存器完成状态以及栈的切换；
2. 跳转到 `gp.sched.pc` 所指的指令地址（`runtime.main`）处执行。

现在已经从 `g0` 切换到了 `gp` 这个 goroutine（main goroutine），它的入口函数是 `runtime.main`：

```go
// The main goroutine.
func main() {
	mp := getg().m // g = main goroutine，不再是 g0 了

	// Racectx of m0->g0 is used only as the parent of the main goroutine.
	// It must not be used for anything else.
	mp.g0.racectx = 0

	// Max stack size is 1 GB on 64-bit, 250 MB on 32-bit.
	// Using decimal instead of binary GB and MB because
	// they look nicer in the stack overflow failure message.
	if goarch.PtrSize == 8 { // 64 位系统上每个 goroutine 的栈最大可达 1G，也就是说 gorputine 的栈虽然可以自动扩展，但它并不是无限扩展的
		maxstacksize = 1000000000
	} else {
		maxstacksize = 250000000
	}

	// An upper limit for max stack size. Used to avoid random crashes
	// after calling SetMaxStack and trying to allocate a stack that is too big,
	// since stackalloc works with 32-bit sizes.
	maxstackceiling = 2 * maxstacksize

	// Allow newproc to start new Ms.
	mainStarted = true

	if haveSysmon {
		// 现在执行的是 main goroutine，所以使用的是 main goroutine 的栈，需要切换到 g0 栈去执行 newm()
		systemstack(func() {
			// 创建监控线程，该线程独立于调度器，不需要跟 p 关联即可运行
			newm(sysmon, nil, -1)
		})
	}

	// ...

	gcenable() // 开启垃圾回收器
    
	// ...

    // main 包的初始化函数，也是由编译器实现，会递归的调用 import 进来的包的初始化函数
	for m := &firstmoduledata; m != nil; m = m.next {
		doInit(m.inittasks)
	}

    // 调用 main.main 函数
	fn := main_main // make an indirect call, as the linker doesn't know the address of the main package when laying down the runtime
	fn()
	
	// ...

    // 进入系统调用，退出进程，可以看出 main goroutine 并未返回，而是直接进入系统调用退出进程了
	exit(0)
	// 保护性代码，如果 exit 意外返回，下面的代码也会让该进程 crash 死掉
	for {
		var x *int32
		*x = 0
	}
}
```

1. 启动一个 sysmon 系统监控线程，该线程负责整个程序的 gc、抢占调度以及 netpoll 等功能的监控。
2. 执行 runtime 包的初始化；
3. 执行 main 包以及 main 包 import 的所有包的初始化；
4. 执行 `main.main` 函数；
5. 从 `main.main` 函数返回后调用 `exit` 系统调用退出进程；

### goexit 函数

main goroutine 调用 exit 直接退出进程了！！

`runtime.main` 是 main goroutine 的入口函数，是在 `schedule()-> execute()-> gogo()` 这个调用链的 `gogo` 函数中用汇编代码直接跳转过来的，而且运行完后会直接退出。

goexit 函数为什么没有调用？

但是在 **`newproc1` 创建 goroutine 的时候已经在其栈上放好了一个返回地址，伪造成 `goexit` 函数调用了 goroutine 的入口函数，这里怎么没有用到这个返回地址啊？**

`newproc1` 函数部分插入 goexit：

```go
`newg.sched.pc = funcPC(goexit) + sys.PCQuantum`
```

因为那是为非 main goroutine 准备的，**非 main goroutine 执行完成后就会返回到 `goexit` 继续执行**，而 main goroutine 执行完成后整个进程就结束了。

## 流程图

<div class="img-zoom lg">
  <img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/go-main-run-flow.png" alt=go-main-run-flow">
</div>