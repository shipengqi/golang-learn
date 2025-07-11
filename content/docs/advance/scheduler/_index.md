---
title: 调度器
weight: 2
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
2. `local`：每个 `p` 都有一个本地队列
3. `global`：全局队列

**先看 runnext，再看 local queue，再看 global queue。当然，如果实在找不到，就去其他 `p` 去偷**。

**goroutine 放到哪个可运行队列？**

1. **如果 `runnext` 为空，那么 goroutine 就会顺利地放入 `runnext`，`runnext` 优先级最高，最先被消费**。
2. **`runnext` 不为空，那就先负责把 `runnext` 上的 old goroutine 踢走，再把 new goroutine 放上来**。
3. `runnext` 中被踢走的 goroutine，**在 local queue 不满时，则将它放入 local queue**；否则意味着 **local queue 已满**，需要减负，会将**它和当前 `p` 的 local queue 中的一半 goroutine 一起放到 global queue 中**。

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

输出的顺序：`9, 0, 1, 2, 3, 4, 5, 6, 7, 8`。这就是因为只有一个 `p`，每次生产出来的 goroutine 都会第一时间塞到 `runnext`，而 `i` 从 `1` 开始，`runnext` 已经有 goroutine 在了，所以这时会把 old goroutine 移到 `p` 的本队队列中去，再把 new goroutine 放到 `runnext`。之后会重复这个过程。

因此这后当一次 `i` 为 `9` 时，新 goroutine 被塞到 `runnext`，其余 goroutine 都在本地队列。

之后，main goroutine 执行了一个读 channel 的语句，这是一个好的调度时机：main goroutine 挂起，运行 `p` 的 `runnext` 和本地可运行队列里的 gorotuine。