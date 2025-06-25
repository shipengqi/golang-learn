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

## 初始化 m0

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