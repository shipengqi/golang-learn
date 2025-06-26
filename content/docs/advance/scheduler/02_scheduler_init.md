---
title: 初始化调度器
weight: 2
---

```go
package main
 
import "fmt"
 
func main() {
    fmt.Println("Hello World!")
}
```

程序的启动过程：

1. 从磁盘上把可执行程序读入内存；
2. 创建进程和主线程；
3. 为主线程分配栈空间；
4. 把由用户在命令行输入的参数拷贝到主线程的栈；
5. 把主线程放入操作系统的运行队列等待被调度执起来运行。

主线程第一次被调度起来执行第一条指令之前，函数栈如下：

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/main-thread-init-stack.png" alt="main-thread-init-stack" style="width:50%;" />

## 程序入口

```bash
[root@shcCDFrh75vm7 ~]# dlv exec ./hello
Type 'help' for list of commands.
(dlv) disass
TEXT _rt0_amd64_linux(SB) /usr/local/go/src/runtime/rt0_linux_amd64.s
=>      rt0_linux_amd64.s:8     0x463940        e9fbc8ffff      jmp $_rt0_amd64
(dlv) si
> _rt0_amd64() /usr/local/go/src/runtime/asm_amd64.s:16 (PC: 0x460240)
Warning: debugging optimized function
TEXT _rt0_amd64(SB) /usr/local/go/src/runtime/asm_amd64.s
=>      asm_amd64.s:16  0x460240        488b3c24        mov rdi, qword ptr [rsp]
        asm_amd64.s:17  0x460244        488d742408      lea rsi, ptr [rsp+0x8]
        asm_amd64.s:18  0x460249        e912000000      jmp $runtime.rt0_go
```

使用 dlv 调试程序可以看到程度的入口早 `runtime/rt0_linux_amd64.s` 文件的第 8 行，执行 `jmp $_rt0_amd64` 跳转到 `runtime/asm_amd64.s` 中的 `_rt0_amd64`。

`runtime/asm_amd64.s`：

```asm
TEXT _rt0_amd64(SB),NOSPLIT,$-8
	MOVQ	0(SP), DI	// argc
	LEAQ	8(SP), SI	// argv
	JMP	runtime·rt0_go(SB)
```

前两行指令把操作系统内核传递过来的参数 `argc` 和 `argv` 数组的地址分别放在 `DI` 和 `SI` 寄存器中。第三行指令跳转到 `rt0_go` 去执行。

```asm
TEXT runtime·rt0_go(SB),NOSPLIT|NOFRAME|TOPFRAME,$0
    // copy arguments forward on an even stack
	MOVQ	DI, AX		// argc
	MOVQ	SI, BX		// argv
    SUBQ $(5*8), SP // 3args 2auto
    ANDQ $~15, SP   // 调整栈顶寄存器使其按 16 字节对齐
    MOVQ AX, 16(SP) // argc 放在 SP + 16字节处
    MOVQ BX, 24(SP) // argv 放在 SP + 24字节处
```

第 4 条指令用于调整栈顶寄存器的值使其按 16 字节对齐，也就是让栈顶寄存器 SP 指向的内存的地址为 16 的倍数，**之所以要按 16 字节对齐，是因为 CPU 有一组 SSE 指令，这些指令中出现的内存地址必须是 16 的倍数**，最后两条指令把 `argc` 和 `argv` 搬到新的位置。

## 初始化 g0

后面的代码，开始初始化全局变量 `g0`，**`g0` 的主要作用是提供一个栈供 runtime 代码执行**：

```asm
// create istack out of the given (operating system) stack.
// _cgo_init may update stackguard.
MOVQ	$runtime·g0(SB), DI // g0 的地址放入 DI 寄存器
LEAQ	(-64*1024)(SP), BX
MOVQ	BX, g_stackguard0(DI)
MOVQ	BX, g_stackguard1(DI)
MOVQ	BX, (g_stack+stack_lo)(DI)
MOVQ	SP, (g_stack+stack_hi)(DI)
```

上面的代码主要是**从系统线程的栈空分出一部分当作 `g0` 的栈**，然后初始化 `g0` 的栈信息和 `stackgard`。

![g0-stack]()

## 主线程与 m0 绑定

设置好 `g0` 栈之后，跳过 CPU 型号检查以及 cgo 初始化相关的代码，直接从 258 行继续分析。

```asm
// 初始化 tls (thread local storage, 线程本地存储)
LEAQ	runtime·m0+m_tls(SB), DI // DI=&m0.tls，取 m0 的 tls 成员的地址到 DI 寄存器
CALL	runtime·settls(SB)       // 调用 settls 设置线程本地存储，settls 函数的参数在 DI 寄存器中

// store through it, to make sure it works
// 验证 settls 是否可以正常工作，如果有问题则 abort 退出程序
get_tls(BX)
MOVQ	$0x123, g(BX)
MOVQ	runtime·m0+m_tls(SB), AX
CMPQ	AX, $0x123
JEQ 2(PC)
CALL	runtime·abort(SB)
```

1. 先调用 `settls` 函数初始化主线程的线程本地存储 (TLS)，目的是把 `m0` 与主线程关联在一起。
2. 证 TLS 功能是否正常，如果不正常则直接 `abort` 退出程序。

`settls` 函数在 `runtime/sys_linx_amd64.s` 文件中：

```asm
// set tls base to DI
TEXT runtime·settls(SB),NOSPLIT,$32
#ifdef GOOS_android
	// Android stores the TLS offset in runtime·tls_g.
	SUBQ	runtime·tls_g(SB), DI
#else
    // DI 寄存器中存放的是 m.tls[0] 的地址，m 的 tls 成员是一个数组
	// 把 DI 寄存器中的地址加 8，存放的就是 m.tls[1] 的地址了
	ADDQ	$8, DI	// ELF wants to use -8(FS)
#endif
	MOVQ	DI, SI
	MOVQ	$0x1002, DI	// ARCH_SET_FS
	MOVQ	$SYS_arch_prctl, AX
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET
```

上面的 `arch_prctl` 系统调用把 `m0.tls[1]` 的地址设置成了 `fs` 段的段基址。CPU 中有个叫 `fs` 的段寄存器。这样通过 `m0.tls[1]` 就可以访问到线程的 TLS 区域了。工作线程代码也可以通过 fs 寄存器来找到 `m.tls`。

CPU 的 FS 寄存器主要用于线程本地存储（TLS），用于在每个线程中快速访问“当前线程的本地数据”。

`rt0_go` 下面的代码会**把 g0 的地址放入主线程的线程本地存储中**，然后通过：

```go
m0.g0 = &g0
g0.m = &m0
```

把 `m0` 和 `g0` 绑定在一起，这样，之后在主线程中通过 `get_tls` 可以获取到 `g0`，通过 `g0` 的 `m` 成员又可以找到 `m0`，于是这里就实现了 `m0` 和 `g0` 与主线程之间的关联。

## 初始化 m0

运行时通过 [runtime.schedinit](https://github.com/golang/go/blob/6796ebb2cb66b316a07998cdcd69b1c486b8579e/src/runtime/proc.go#L798) 初始化调度器：

```go
func schedinit() {
    // ...
    // getg 函数在源代码中没有对应的定义，由编译器插入类似下面两行代码
    // get_tls(CX) 
    // MOVQ g(CX), BX; 
    // BX 存器里面现在放的是当前 g 结构体对象的地址
	gp := getg()
	// ...

	sched.maxmcount = 10000

    mcommoninit(gp.m, -1)

	// ...
	sched.lastpoll = uint64(nanotime())
	procs := ncpu
	if n, ok := atoi32(gogetenv("GOMAXPROCS")); ok && n > 0 {
		procs = n
	}
	if procresize(procs) != nil {
		throw("unknown runnable goroutine during bootstrap")
	}
    // ...
}
```

1. `g0` 的地址已经被设置到了线程本地存储之中，通过 `getg` 函数（`getg` 函数是编译器实现的，源代码中是找不到其定义）从线程本地存储中获取当前正在运行的 `g`。
2. `mcommoninit` 对 `m0` 进行必要的初始化。
3. 调用 `procresize` 初始化系统需要用到的 `p` 结构体对象。它的数量决定了最多可以有都少个 goroutine 同时并行运行。
4. `sched.maxmcount = 10000` 一个 Go 程序最多可以创建 10000 个线程。
5. 线程数可以通过 `GOMAXPROCS` 变量控制。

[`mcommoninit`](https://github.com/golang/go/blob/6796ebb2cb66b316a07998cdcd69b1c486b8579e/src/runtime/proc.go#L942) 初始化 `m0`：

```go
func mcommoninit(mp *m, id int64) {
   gp := getg() // 初始化过程中 gp = g0
 
    // g0 stack won't make sense for user (and is not necessary unwindable).
    if gp != gp.m.g0 { // 函数调用栈 traceback，不需要关心
		callers(1, mp.createstack[:])
	}

    lock(&sched.lock)
	if id >= 0 {
		mp.id = id
	} else {
		mp.id = mReserveID()
	}
    // random 初始化
	mrandinit(mp)
    // 创建用于信号处理的 gsignal，只是简单的从堆上分配一个 g 结构体对象,然后把栈设置好就返回了
	mpreinit(mp)
	if mp.gsignal != nil {
		mp.gsignal.stackguard1 = mp.gsignal.stack.lo + stackGuard
	}

	// Add to allm so garbage collector doesn't free g->m
	// when it is just in a register or thread-local storage.
    // 把 m0 加入到 allm 全局链表中
	mp.alllink = allm

	// NumCgoCall() and others iterate over allm w/o schedlock,
	// so we need to publish it safely.
	atomicstorep(unsafe.Pointer(&allm), unsafe.Pointer(mp))
	unlock(&sched.lock)

	// Allocate memory to hold a cgo traceback if the cgo call crashes.
	if iscgo || GOOS == "solaris" || GOOS == "illumos" || GOOS == "windows" {
		mp.cgoCallers = new(cgoCallers)
	}
	mProfStackInit(mp)
}
```

这里并未对 `m0` 做什么关于调度相关的初始化，可以简单的认为这个函数只是把 `m0` 放入全局链表 `allm` 之中就返回了。

## 初始化 allp

```go
func procresize(nprocs int32) *p {
    // ...

    // 系统初始化时 gomaxprocs = 0
	old := gomaxprocs
    // ...

	// Grow allp if necessary.
	if nprocs > int32(len(allp)) { // 初始化时 len(allp) == 0
		// Synchronize with retake, which could be running
		// concurrently since it doesn't run on a P.
		lock(&allpLock)
		if nprocs <= int32(cap(allp)) {
			allp = allp[:nprocs]
		} else {  // 初始化时进入此分支，创建 allp 切片
			nallp := make([]*p, nprocs)
			// Copy everything up to allp's cap so we
			// never lose old allocated Ps.
			copy(nallp, allp[:cap(allp)])
			allp = nallp
		}
        // ...
		unlock(&allpLock)
	}

	// initialize new P's
    // 循环创建 nprocs 个 p 并完成基本初始化
	for i := old; i < nprocs; i++ {
		pp := allp[i]
		if pp == nil {
			pp = new(p) // 调用内存分配器从堆上分配一个 struct p
		}
		pp.init(i)
		atomicstorep(unsafe.Pointer(&allp[i]), unsafe.Pointer(pp))
	}

	gp := getg()
	if gp.m.p != 0 && gp.m.p.ptr().id < nprocs { // 初始化时 m0->p 还未初始化，所以不会执行这个分支
		// continue to use the current P
		gp.m.p.ptr().status = _Prunning
		gp.m.p.ptr().mcache.prepareForSweep()
	} else { // 初始化时执行这个分支
		// release the current P and acquire allp[0].
		//
		// We must do this before destroying our current P
		// because p.destroy itself has write barriers, so we
		// need to do that from a valid P.
		if gp.m.p != 0 { // 初始化时这里不执行
			trace := traceAcquire()
			if trace.ok() {
				// Pretend that we were descheduled
				// and then scheduled again to keep
				// the trace consistent.
				trace.GoSched()
				trace.ProcStop(gp.m.p.ptr())
				traceRelease(trace)
			}
			gp.m.p.ptr().m = 0
		}
		gp.m.p = 0
		pp := allp[0]
		pp.m = 0
		pp.status = _Pidle
		acquirep(pp) // 把 p 和 m0 关联起来，其实是这两个 strct 的成员相互赋值
		trace := traceAcquire()
		if trace.ok() {
			trace.GoStart()
			traceRelease(trace)
		}
	}

	// g.m.p is now set, so we no longer need mcache0 for bootstrapping.
	mcache0 = nil

	// ...

    // 循环把所有空闲的 p 放入空闲链表
	var runnablePs *p
	for i := nprocs - 1; i >= 0; i-- {
		pp := allp[i]
		if gp.m.p.ptr() == pp { // allp[0] 跟 m0 关联了，所以是不能放到空闲链表
			continue
		}
		pp.status = _Pidle
		if runqempty(pp) { // 初始化时除了 allp[0] 其它 p 全部执行这个分支，放入空闲链表
			pidleput(pp, now)
		} else {
            // ...
		}
	}
    // ...
}
```

1. 使用 `make([]*p, nprocs)` 初始化全局变量 `allp`，即 `allp = make([]*p, nprocs)`；
2. 循环创建并初始化 `nprocs` 个 `p` 结构体对象并依次保存在 `allp` 切片之中；
3. 把 `m0` 和 `allp[0]` 绑定在一起，即 `m0.p = allp[0]`, `allp[0].m = m0`；
4. 把除了 `allp[0]` 之外的所有 `p` 放入到全局变量 `sched` 的 `pidle` 空闲队列之中。

`procresize` 函数执行完后，调度器相关的初始化工作就基本结束了。