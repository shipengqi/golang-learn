---
title: 垃圾回收
weight: 2
---

Go 语言中使用的垃圾回收使用的是**标记清扫算法**。标记清理最典型的做法是三⾊标记。进行垃圾回收时会 STW(stop the world），
就是 **runtime 把所有的线程全部冻结掉，意味着⽤户逻辑都是暂停的，所有的⽤户对象都不会被修改了**，这时候去扫描肯定是安全的，
对象要么活着要么死着，所以会造成中间暂停时间可能会很⻓，⽤户逻辑对于⽤户的反应就中⽌了。

Go GC 的基本特征：非分代，非紧缩，写屏障，并发标记清理。

## 三色标记和写屏障

白色对象 — 潜在的垃圾，其内存可能会被垃圾收集器回收；
黑色对象 — 活跃的对象，包括不存在任何引用外部指针的对象以及从根对象可达的对象；
灰色对象 — 活跃的对象，因为存在指向白色对象的外部指针，垃圾收集器会扫描这些对象的子对象；

三色标记算法原理如下：

1. 起初所有对象都是白色。
2. 从根出发扫描所有可达对象，标记为灰色，放入待处理队列。
3. 从队列取出灰色对象，将其引用对象标记为灰色放入队列，自身标记为黑色。
4. 重复 3，直到灰色对象队列为空。

扫描和标记完成后，只剩下白色（待回收）和黑色（活跃对象）的对象，清理操作将白色对象内存回收。

在垃圾收集器开始工作时，程序中不存在任何的黑色对象，垃圾收集的根对象会被标记成灰色，垃圾收集器只会从灰色对象集合中取出对象开始扫描，当灰色集合中不存在任何对象时，标记阶段就会结束。

![](tri-color-mark-sweep.png)

因为用户程序可能在标记执行的过程中修改对象的指针，所以三色标记清除算法本身是不可以并发或者增量执行的，它仍然需要 **STW**。在如下所示的三色标记过程中，用户程序建立了从 A 对象到 D 对象的引用，但是因为程序中已经不存在灰色对象了，所以 D 对象会被垃圾收集器错误地回收。

![](tri-color-mark-sweep-and-mutator.png)

本来不应该被回收的对象却被回收了，这在内存管理中是非常严重的错误，我们将这种错误称为**悬挂指针**，即指针没有指向特定类型的合法对象，影响了内存的安全性

### 屏障技术

想要在并发或者增量的标记算法中保证正确性，我们需要达成以下两种三色不变性（Tri-color invariant）中的任意一种：

- 强三色不变性 — 黑色对象不会指向白色对象，只会指向灰色对象或者黑色对象；
- 弱三色不变性 — 黑色对象指向的白色对象必须包含一条从灰色对象经由多个白色对象的可达路径

![](strong-weak-tricolor-invariant.png)

遵循上述两个不变性中的任意一个，我们都能保证垃圾收集算法的正确性。而屏障技术就是在并发或者增量标记过程中保证三色不变性的重要技术。

**垃圾收集中的屏障技术更像是一个钩子方法**，它是在用户程序读取对象、创建新对象以及更新对象指针时执行的一段代码，根据操作类型的不同，我们可以将它们分成**读屏障**（Read barrier）和**写屏障**（Write barrier）两种，因为**读屏障需要在读操作中加入代码片段，对用户程序的性能影响很大，所以编程语言往往都会采用写屏障保证三色不变性**。

### 增量和并发

增量式（Incremental）的垃圾收集是减少程序最长暂停时间的一种方案，它可以将原本时间较长的暂停时间切分成多个更小的 GC 时间片，虽然从垃圾收集开始到结束的时间更长了，但是这也减少了应用程序暂停的最大时间

增量式的垃圾收集需要与三色标记法一起使用，为了保证垃圾收集的正确性，我们需要在垃圾收集开始前打开写屏障，这样用户程序对内存的修改都会先经过写屏障的处理，保证了堆内存中对象关系的强三色不变性或者弱三色不变性。虽然增量式的垃圾收集能够减少最大的程序暂停时间，但是增量式收集也会增加一次 GC 循环的总时间，在垃圾收集期间，因为写屏障的影响用户程序也需要承担额外的计算开销，所以增量式的垃圾收集也不是只有优点的。

并发（Concurrent）的垃圾收集不仅能够减少程序的最长暂停时间，还能减少整个垃圾收集阶段的时间，通过开启读写屏障、利用多核优势与用户程序并行执行，并发垃圾收集器确实能够减少垃圾收集对应用程序的影响

虽然并发收集器能够与用户程序一起运行，但是并不是所有阶段都可以与用户程序一起运行，部分阶段还是需要暂停用户程序的，不过与传统的算法相比，并发的垃圾收集可以将能够并发执行的工作尽量并发执行；当然，因为读写屏障的引入，并发的垃圾收集器也一定会带来额外开销，不仅会增加垃圾收集的总时间，还会影响用户程序，这是我们在设计垃圾收集策略时必须要注意的。

## 何时触发 GC

垃圾回收器在初始化时，设置 `gcpercent` 和 `next_gc` 阈值。

### 自动垃圾回收

为对象分配内存以后，`mallocgc` 函数会检查 GC 触发条件。
**在堆上分配大于 `maxSmallSize` （32K byte）的对象时进行检测此时是否满足垃圾回收条件，如果满足则进行垃圾回收**。

```go
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    ...
    shouldhelpgc := false
    // 分配的对象小于 maxSmallSize (32K byte)
    if size <= maxSmallSize {
        ...
    } else {
        shouldhelpgc = true
        ...
    }
    ...
    // gcShouldStart() 函数进行触发条件检测
    if shouldhelpgc && gcShouldStart(false) {
        // gcStart() 函数进行垃圾回收
        gcStart(gcBackgroundMode, false)
    }
}
```

#### GC 触发条件

触发时机
运行时会通过如下所示的 runtime.gcTrigger.test 方法决定是否需要触发垃圾收集，当满足触发垃圾收集的基本条件时 — 允许垃圾收集、程序没有崩溃并且没有处于垃圾收集循环，该方法会根据三种不同的方式触发进行不同的检查：

```go
func (t gcTrigger) test() bool {
 if !memstats.enablegc || panicking != 0 || gcphase != _GCoff {
  return false
 }
 switch t.kind {
 case gcTriggerHeap:
  return memstats.heap_live >= memstats.gc_trigger
 case gcTriggerTime:
  if gcpercent < 0 {
   return false
  }
  lastgc := int64(atomic.Load64(&memstats.last_gc_nanotime))
  return lastgc != 0 && t.now-lastgc > forcegcperiod
 case gcTriggerCycle:
  return int32(t.n-work.cycles) > 0
 }
 return true
}
```

gcTriggerHeap — 堆内存的分配达到达控制器计算的触发堆大小；
gcTriggerTime — 如果一定时间内没有触发，就会触发新的循环，该出发条件由 runtime.forcegcperiod 变量控制，默认为 2 分钟；
gcTriggerCycle — 如果当前没有开启垃圾收集，则触发新的循环；

runtime.sysmon 和 runtime.forcegchelper — 后台运行定时检查和垃圾收集；
runtime.GC — 用户程序手动触发垃圾收集；
runtime.mallocgc — 申请内存时根据堆大小触发垃圾收集；

触发条件主要关注下面代码中的中间部分：`forceTrigger || memstats.heap_live >= memstats.gc_trigger`。
`forceTrigger` 是 `forceGC` 的标志；后面半句的意思是当前堆上的活跃对象大于我们初始化时候设置的 GC 触发阈值。
在 malloc 以及 free 的时候 `heap_live` 会一直进行更新。

```go
// gcShouldStart returns true if the exit condition for the _GCoff
// phase has been met. The exit condition should be tested when
// allocating.
//
// If forceTrigger is true, it ignores the current heap size, but
// checks all other conditions. In general this should be false.
func gcShouldStart(forceTrigger bool) bool {
    return gcphase == _GCoff && (forceTrigger || memstats.heap_live >= memstats.gc_trigger) && memstats.enablegc && panicking == 0 && gcpercent >= 0
}

// 初始化的时候设置 GC 的触发阈值
func gcinit() {
    _ = setGCPercent(readgogc())
    memstats.gc_trigger = heapminimum
    ...
}
// 启动的时候通过 GOGC 传递百分比 x
// 触发阈值等于 x * defaultHeapMinimum (defaultHeapMinimum 默认是 4M)
func readgogc() int32 {
    p := gogetenv("GOGC")
    if p == "off" {
        return -1
    }
    if n, ok := atoi32(p); ok {
        return n
    }
    return 100
}
```

`heap_live` 是活跃对象总量。

### 主动垃圾回收

主动垃圾回收，通过调用 `runtime.GC()`，这是阻塞式的。

```go
// GC runs a garbage collection and blocks the caller until the
// garbage collection is complete. It may also block the entire
// program.
func GC() {
    gcStart(gcForceBlockMode, false)
}
```

## 监控

在一个场景中：服务重启，海量的客户端接入，瞬间分配了大量对象，这会将 GC 的触发条件 `next_gc` 推到一个很大的值。
在服务正常以后，由于活跃对象远远小于改阈值，会导致 GC 无法触发，大量白色对象不能被回收，最终造成内存泄露。

所以 GC 的最后一道保险，就是监控线程 sysmon，sysmon 每隔 2 分钟会检查一次 GC 状态，超过 2 分钟则强制执行。

## 逃逸分析

手动分配内存会导致如下的两个问题：

不需要分配到堆上的对象分配到了堆上 — 浪费内存空间；
需要分配到堆上的对象分配到了栈上 — 悬挂指针、影响内存安全；

逃逸分析（Escape analysis）是用来决定指针动态作用域的方法

Go 语言的逃逸分析遵循以下两个不变性：

- **指向栈对象的指针不能存在于堆中**；
- **指向栈对象的指针不能在栈对象回收后存活**；

![](escape-analysis-and-key-invariants.png)

上图展示两条不变性存在的意义，当我们违反了第一条不变性，堆上的绿色指针指向了栈中的黄色内存，一旦当前函数返回函数栈被回收，该绿色指针指向的值就不再合法；如果我们违反了第二条不变性，因为寄存器 SP 下面的内存由于函数返回已经被释放掉，所以黄色指针指向的内存已经不再合法。