---
title: GPM 模型和一些重要的数据结构
weight: 1
---

**goroutine 建立在操作系统线程基础之上，它与操作系统线程之间实现了一个多对多 (`M:N`) 的两级线程模型**。

这里的 `M:N` 是指 **M 个 goroutine 运行在 N 个操作系统线程之上**。**操作系统内核负责对这 N 个操作系统线程进行调度**，而这 **N 个系统线程又负责对这 M 个 goroutine 进行调度**和运行。

goroutine 调度器的大概工作流程：

```go
// 程序启动时的初始化代码
......
for i = 0; i < N; i++ { // 创建 N 个操作系统线程执行 schedule 函数
    create_os_thread(schedule) // 创建一个操作系统线程执行 schedule 函数
}

// schedule 函数实现调度逻辑
schedule() {
   for { // 调度循环
        // 根据某种算法从 M 个 goroutine 中找出一个需要运行的 goroutine
        g = find_a_runnable_goroutine_from_M_goroutines()
        run_g(g) // CPU 运行该 goroutine，直到需要调度其它 goroutine 才返回
        save_status_of_g(g) // 保存 goroutine 的状态，主要是寄存器的值
    }
}
```

1. 程序运行起来之后创建了 N 个由内核调度的操作系统线程去执行 `shedule` 函数。
2. `schedule` 函数在一个调度循环中反复从 M 个 goroutine 中挑选出一个需要运行的 goroutine 并跳转到该 goroutine 去运行。
3. 直到需要调度其它 goroutine 时才返回到 `schedule` 函数中通过 `save_status_of_g` 保存刚刚正在运行的 goroutine 的状态然后再次去寻找下一个 goroutine。

{{< callout type="info" >}}
系统线程对 goroutine 的调度与内核对系统线程的调度原理是一样的，都是**通过保存和修改 CPU 寄存器的值来达到切换线程 或 goroutine 的目的**。
{{< /callout >}}

## 调度器相关数据结构

### g 结构体（goroutine）

为了实现对 goroutine 的调度，需要引入一个数据结构来保存 CPU 寄存器的值以及 goroutine 的其它一些状态信息，在 Go 调度器源代码中，这个数据结构是一个名叫 **`g` 的结构体**。该结构体的每一个实例对象都代表了一个 goroutine。

调度器代码可以通过 `g` 对象来对 goroutine 进行调度:

- **当 goroutine 被调离 CPU 时，调度器代码负责把 CPU 寄存器的值保存在 `g` 对象的成员变量之中**；
- **当 goroutine 被调度起来运行时，调度器代码又负责把 `g` 对象的成员变量所保存的寄存器的值恢复到 CPU 的寄存器**。

### schedt 结构体（调度器）

只有 `g` 结构体对象是不够的，还需要一个**存放所有（可运行）goroutine 的容器**，便于工作线程寻找需要被调度起来运行的 goroutine，于是 Go 调度器又引入了 **`schedt` 结构体**：

- 用来保存调度器自身的状态信息；
- 保存 goroutine 的运行队列。

#### 全局运行队列

每个 Go 程序只有一个调度器，所以在每个 Go 程序中 **`schedt` 结构体只有一个实例对象**，该实例对象在源代码中被定义成了一个共享的全局变量，这样每个工作线程都可以访问它以及它所拥有的 goroutine 运行队列，我们称这个运行队列为**全局运行队列**。

#### 线程运行队列

这个**全局运行队列是每个工作线程都可以访问的，那就涉及到并发的问题**，因此需要加锁。但是在高并发的场景下，加锁是会导致性能问题的。于是调度器又为**每个工作线程引入了一个私有的局部 goroutine 运行队列**，工作线程优先使用自己的局部运行队列，只有必要时才会去访问全局运行队列，这大大减少了锁冲突，提高了工作线程的并发性。**局部运行队列被包含在 `p` 结构体**的实例对象之中，每一个运行着 Go 代码的**工作线程都会与一个 `p` 结构体的实例对象关联在一起**。

### m 结构体（工作线程）

**每个工作线程都有唯一的一个 `m` 结构体的实例对象与之对应**，`m` 结构体对象除了记录着工作线程的诸如栈的起止位置、当前正在执行的 goroutine 以及是否空闲等等状态信息之外，还通过指针维持着与 `p` 结构体的实例对象之间的绑定关系。于是，通过 `m` 既可以找到与之对应的工作线程正在运行的 goroutine，又可以找到工作线程的局部运行队列等资源。

### GPM 模型

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/gpm.png" alt="gpm" style="width:60%;" />

灰色的 `g` 表示处于运行队列之中正在等待被调度起来运行的 goroutine。

**每个 `m` 都绑定了一个 `p`，每个 `p` 都有一个私有的本地 goroutine 队列，`m` 对应的线程从本地和全局 goroutine 队列中获取 goroutine 并运行**。

#### 工作线程如何绑定到 m 结构体实例对象

多个工作线程和多个 `m` 需要一一对应，如何实现？**线程本地存储**。线程本地存储其实就是线程私有的全局变量，这正是我们需要的。只要**每个工作线程拥有了各自私有的 `m` 结构体全局变量，就能在不同的工作线程中使用相同的全局变量名来访问不同的 `m` 结构体对象**。

每个**工作线程在刚刚被创建出来进入调度循环之前就利用线程本地存储机制为该工作线程实现了一个指向 `m` 结构体实例对象的私有全局变量**，这样在之后的代码中就**使用该全局变量来访问自己的 `m` 结构体对象以及与 `m` 相关联的 `p` 和 `g` 对象**（工作线程可以直接从本地线程存储取出来 `m`）。

调度伪代码：

```go
// 程序启动时的初始化代码
......
for i = 0; i < N; i++ { // 创建 N 个操作系统线程执行 schedule 函数
    create_os_thread(schedule) // 创建一个操作系统线程执行 schedule 函数
}

// 定义一个线程私有全局变量，注意它是一个指向m结构体对象的指针
// ThreadLocal 用来定义线程私有全局变量
ThreadLocal self *m

// schedule 函数实现调度逻辑
schedule() {
    // 创建和初始化 m 结构体对象，并赋值给私有全局变量 self
    self = initm()   
    for { // 调度循环
        if(self.p.runqueue is empty) { // 本地运行队列为空
            // 从全局运行队列中找出一个需要运行的 goroutine
            g = find_a_runnable_goroutine_from_global_runqueue()
        } else {
            // 从私有的本地运行队列中找出一个需要运行的 goroutine
            g = find_a_runnable_goroutine_from_local_runqueue()
        }
        run_g(g) // CPU 运行该 goroutine，直到需要调度其它 goroutine 才返回
        save_status_of_g(g) // 保存 goroutine 的状态，主要是寄存器的值
     }
}
```

### 重要的结构体

这些结构体的定义全部在 `runtime/runtime2.go` 源码文件中：

#### stack 结构体

记录 goroutine 所使用的栈的信息，包括栈顶和栈底位置：

```go
// Stack describes a Go execution stack.
// The bounds of the stack are exactly [lo, hi),
// with no implicit data structures on either side.
// 用于记录 goroutine 使用的栈的起始和结束位置
type stack struct{ 
    lo uintptr   // 栈顶，低地址
    hi uintptr   // 栈底，高地址
}
```

#### gobuf 结构体

用于保存 goroutine 的调度信息，主要包括 CPU 的几个寄存器的值：

```go
type gobuf struct {
    sp  uintptr  // 保存 CPU 的 rsp 寄存器的值
    pc  uintptr  // 保存 CPU 的 rip 寄存器的值
    g   guintptr // 记录当前这个 gobuf 对象属于哪个 goroutine
    ctxt unsafe.Pointer
  
   // 保存系统调用的返回值，因为从系统调用返回之后如果 p 被其它工作线程抢占，
   // 则这个 goroutine 会被放入全局运行队列被其它工作线程调度，其它线程需要知道系统调用的返回值。
    ret sys.Uintreg
    lr  uintptr
    bp  uintptr// for GOEXPERIMENT=framepointer
}
```

#### g 结构体

代表一个 goroutine，该结构体保存了 goroutine 的所有信息，包括栈，`gobuf` 结构体和其它的一些状态信息：

```go
type g struct {
	// goroutine 使用的栈
	stack       stack   // offset known to runtime/cgo
	// 下面两个成员用于栈溢出检查，实现栈的自动伸缩，抢占调度也会用到 stackguard0
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

	_panic         *_panic // innermost panic - offset known to liblink
	_defer         *_defer // innermost defer
    // 当前与 g 绑定的 m
	m              *m      // current m; offset known to arm liblink
	// 保存调度信息，主要是几个寄存器的值
	sched          gobuf
	syscallsp      uintptr        // if status==Gsyscall, syscallsp = sched.sp to use during gc
	syscallpc      uintptr        // if status==Gsyscall, syscallpc = sched.pc to use during gc
	stktopsp       uintptr        // expected sp at top of stack, to check in traceback
	param          unsafe.Pointer // passed parameter on wakeup
	atomicstatus   uint32
	stackLock      uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
	goid           int64
	// schedlink 字段指向全局运行队列中的下一个 g，所有位于全局运行队列中的 g 形成一个链表
	schedlink      guintptr
	// g 被阻塞之后的近似时间
	waitsince      int64      // approx time when the g become blocked
	// g 被阻塞的原因
	waitreason     waitReason // if status==Gwaiting
	// 抢占调度标志，如果需要抢占调度，设置 preempt 为 true
	preempt        bool       // preemption signal, duplicates 
	// ...
}
```

#### m 结构体

代表工作线程，保存了 `m` 自身使用的栈信息，当前正在运行的 goroutine 以及与 `m` 绑定的 `p` 等信息：

`g` 需要调度到 `m` 上才能运行，`m` 是真正工作的人。

**当 `m` 没有工作可做的时候，在它休眠前，会“自旋”地来找工作：检查全局队列，查看 network poller，试图执行 gc 任务，或者“偷”工作**。

```go
type m struct {
	// g0 主要用来记录工作线程使用的栈信息，在执行调度代码时需要使用这个栈
    // 执行用户 goroutine 代码时，使用用户 goroutine 自己的栈，调度时会发生栈的切换
	g0      *g     // goroutine with scheduling stack
    // ...

	// 通过 TLS 实现 m 结构体对象与工作线程之间的绑定
	tls           [6]uintptr   // thread-local storage (for x86 extern register)
	mstartfn      func()
	// 指向正在运行的 gorutine 对象
	curg          *g       // current running goroutine
	caughtsig     guintptr // goroutine running during fatal signal
	
    // 当前工作线程绑定的 p
	p             puintptr // attached p for executing go code (nil if not executing go code)
	nextp         puintptr
	oldp          puintptr // the p that was attached before executing a 
	// ...
	// spinning 状态：表示当前工作线程正在试图从其它工作线程的本地运行队列偷取goroutine
	// 
	spinning      bool // m is out of work and is actively looking for work
	// m 正阻塞在 note 上
	blocked       bool // m is blocked on a note
	// ...
	// 正在执行 cgo 调用
	incgo         bool   // m is executing a cgo call
	// ...
	// 没有 goroutine 需要运行时，工作线程睡眠在这个 park 成员上，其它线程通过这个 park 唤醒该工作线程
	park          note
	// 记录所有工作线程的一个链表
	alllink       *m // on allm
	// ...
	// Linux 平台 thread 的值就是操作系统线程 ID
	thread        uintptr // thread handle
	freelink      *m      // on sched.freem
	// ...
}
```

#### p 结构体

`p` 是 processor 的意思。

保存工作线程执行 Go 代码时所必需的资源，比如 goroutine 的运行队列，内存分配用到的缓存等等。

```go
type p struct {
	// ...
	     
    // 在 allp 中的索引
    id          int32
    // 每次调用 schedule 时会加一
    schedtick   uint32
	// 每次系统调用时加一
    syscalltick uint32
    // 用于 sysmon 线程记录被监控 p 的系统调用时间和运行时间
    sysmontick  sysmontick // last tick observed by sysmon
	// 指向绑定的 m，如果 p 是 idle 的话，那这个指针是 nil
	m           muintptr   // back-link to associated m (nil if idle)
    // ...

	// Queue of runnable goroutines. Accessed without lock.
	// 本地 goroutine 运行队列
	runqhead uint32 // 队列头
	runqtail uint32 // 队列尾
	runq     [256]guintptr // 使用数组实现的循环队列
  
    // runnext 非空时，代表的是一个 runnable 状态的 G，
    // 这个 G 被当前 G 修改为 ready 状态，相比 runq 中的 G 有更高的优先级。
    // 如果当前 G 还有剩余的可用时间，那么就应该运行这个 G
    // 运行之后，该 G 会继承当前 G 的剩余时间
	runnext guintptr

    // 空闲的 g
    gfree    *g
	// ...
}
```

#### schedt 结构体

保存调度器的状态信息和 goroutine 的全局运行队列：

```go
type schedt struct {
    // ...
 
    // 由空闲的工作线程组成链表
    midle       muintptr // idle m's waiting for work
    // 空闲的工作线程的数量
    nmidle        int32 // number of idle m's waiting for work
    nmidlelocked  int32 // number of locked m's waiting for work
    mnext         int64 // number of m's that have been created and next M ID
    // 最多只能创建 maxmcount 个工作线程
    maxmcount    int32 // maximum number of m's allowed (or die)
    nmsys        int32 // number of system m's not counted for deadlock
    nmfreed      int64 // cumulative number of freed m's
 
    ngsys        uint32 // number of system goroutines; updated atomically
 
    // 由空闲的 p 结构体对象组成的链表
    pidle     puintptr // idle p's
    // 空闲的 p 结构体对象的数量
    npidle     uint32
    nmspinning uint32 // See "Worker thread parking/unparking" comment in proc.go.
 
    // Global runnable queue.
    // goroutine 全局运行队列
    runq       gQueue
    runqsize   int32
 
    ......
 
    // Global cache of dead G's.
    // gFree 是所有已经退出的 goroutine 对应的 g 结构体对象组成的链表
    // 用于缓存 g 结构体对象，避免每次创建 goroutine 时都重新分配内存
    gFree struct{
        lock        mutex
        stack       gList // Gs with stacks
        noStack     gList // Gs without stacks
        n           int32
    }
  
    ......
}
```

### 重要的全局变量

```go

allgs    []*g   // 保存所有的 g
allm     *m     // 所有的 m 构成的一个链表，包括下面的 m0
allp     []*p  // 保存所有的 p，len(allp) == gomaxprocs
 
ncpu         int32  // 系统中 cpu 核的数量，程序启动时由 runtime 代码初始化
gomaxprocs   int32  // p 的最大值，默认等于 ncpu，但可以通过 GOMAXPROCS 修改
 
sched     schedt    // 调度器结构体对象，记录了调度器的工作状态
 
m0 m        // 代表进程的主线程
g0 g        // m0 的 g0，也就是 m0.g0 = &g0
```
