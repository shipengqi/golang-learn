---
title: 调度器
weight: 3
---

goroutine 是 Go 实现的用户态线程，主要用来解决操作系统线程两个方面的问题：

1. 创建和切换太重：操作系统线程的创建和切换都需要进入内核，而进入内核所消耗的性能代价比较高，开销较大；
2. 内存使用太重：一方面，为了尽量避免极端情况下操作系统线程栈的溢出，内核在创建操作系统线程时默认会为其分配一个较大的栈内存（虚拟地址空间，内核并不会一开始就分配这么多的物理内存），然而在绝大多数情况下，系统线程远远用不了这么多内存，这导致了浪费；另一方面，栈内存空间一旦创建和初始化完成之后其大小就不能再有变化，这决定了在某些特殊场景下系统线程栈还是有溢出的风险。

用户态的 goroutine 则轻量得多：

1. goroutine 是用户态线程，其**创建和切换都在用户代码中完成而无需进入操作系统内核**，所以其开销要远远小于系统线程的创建和切换；
2. goroutine 启动时默认栈大小只有 2k，这在多数情况下已经够用了，即使不够用，goroutine 的栈也会自动扩大，同时，如果栈太大了过于浪费它还能自动收缩，这样既没有栈溢出的风险，也不会造成栈内存空间的大量浪费。

## Go 调度的本质

Go 调度的本质是一个**生产-消费流程**。`m` 拿到 goroutine 并运行它的过程就是一个消费过程。

![scheduler-queue](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/scheduler-queue.png)

生产出的 goroutine 就放在可运行队列中。可运行队列是分为三级：

1. `runnext`：实际上只能指向一个 goroutine。
2. `local`：每个 P 都有一个本地队列
3. `global`：全局队列

**先看 runnext，再看 local queue，再看 global queue。当然，如果实在找不到，就去其他 `p` 去偷**。

**goroutine 放到哪个可运行队列？**

1. 如果 `runnext` 为空，那么 goroutine 就会顺利地放入 `runnext`，`runnext` 优先级最高，最先被消费。
2. `runnext` 不为空，那就先负责把 `runnext` 上的 old goroutine 踢走，再把 new goroutine 放上来。
3. `runnext` 中被踢走的 goroutine，在 local queue 不满时，则将它放入 local queue；否则意味着 local queue 已满，需要减负，会将它和当前 p 的 local queue 中的一半 goroutine 一起放到 global queue 中。

```go
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
    runtime.GOMAXPROCS(1)
    for i := 0; i < 10; i++ {
        i := i
        go func() {
            fmt.Println(i)
        }()
    }

    var ch = make(chan int)
    <- ch
}

// 输出
// 9
// 0
// 1
// 2
// 3
// 4
// 5
// 6
// 7
// 8
// fatal error: all goroutines are asleep - deadlock!

// goroutine 1 [chan receive]:
// main.main()
// 	C:/Code/my-repos/example.v1/advance/scheduler/v1/main.go:18 +0x6c
```

 输出的顺序：`9, 0, 1, 2, 3, 4, 5, 6, 7, 8`。这就是因为只有一个 `p`，每次生产出来的 goroutine 都会第一时间塞到 `runnext`，而 `i` 从 `1` 开始，`runnext` 已经有 goroutine 在了，所以这时会把 old goroutine 移到 `p` 的本队队列中去，再把 new goroutine 放到 runnext。之后会重复这个过程。

因此这后当一次 `i` 为 `9` 时，新 goroutine 被塞到 runnext，其余 goroutine 都在本地队列。

之后，main goroutine 执行了一个读 channel 的语句，这是一个好的调度时机：main goroutine 挂起，运行 `p` 的 `runnext` 和本地可运行队列里的 gorotuine。

## Go 调度器

**goroutine 建立在操作系统线程基础之上，它与操作系统线程之间实现了一个多对多 (`M:N`) 的两级线程模型**。

这里的 `M:N` 是指 **M 个goroutine运行在 N 个操作系统线程之上**。**操作系统内核负责对这 N 个操作系统线程进行调度**，而这 **N 个系统线程又负责对这 M 个 goroutine 进行调度**和运行。

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

每个**工作线程在刚刚被创建出来进入调度循环之前就利用线程本地存储机制为该工作线程实现了一个指向 `m` 结构体实例对象的私有全局变量**，这样在之后的代码中就**使用该全局变量来访问自己的 `m` 结构体对象以及与 `m` 相关联的 `p` 和 `g` 对象**。

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

p 是 processor 的意思。

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

## 调度器初始化

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

### 初始化 m0

[`mcommoninit`](https://github.com/golang/go/blob/6796ebb2cb66b316a07998cdcd69b1c486b8579e/src/runtime/proc.go#L942) 初始化 `m0`：

```go
func mcommoninit(mp *m, id int64) {
   gp := getg() //初始化过程中_g_ = g0
 
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

### 初始化 allp

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

## 创建 main goroutine

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


## 调度循环