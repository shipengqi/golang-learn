---
title: 系统监控
weight: 5
---

Go 的系统监控 (sysmon) 是一个独立的特殊线程，它在后台持续运行，负责处理各种运行时系统的维护任务和监控工作。它在内部启动了一个不会中止的循环，在循环的内部会轮询网络、抢占长期运行或者处于系统调用的 goroutine 以及触发垃圾回收，通过这些行为，它能够让系统的运行状态变得更健康。

`runtime.main` 函数：

```go
if haveSysmon {
    // 现在执行的是 main goroutine，所以使用的是 main goroutine 的栈，需要切换到 g0 栈去执行 newm()
    systemstack(func() {
        // 创建专用监控线程，该线程独立于调度器，不需要跟 p 关联即可运行
        newm(sysmon, nil, -1)
    })
}
```

## 监控循环

```go
// runtime/proc.go
func sysmon() {
	lock(&sched.lock)
	sched.nmsys++
	checkdead() // 检查是否存在死锁
	unlock(&sched.lock)

	lasttrace := int64(0)
	idle := 0 // how many cycles in succession we had not wokeup somebody
	delay := uint32(0)
	for {
        // 动态调整休眠时间
		if idle == 0 {
			delay = 20 // 20μs
		} else if idle > 50 {
			delay *= 2 // 指数退避
		}
		if delay > 10*1000 { // 上限 10ms
			delay = 10 * 1000
		}
		usleep(delay)
		// ...
	}
}
```

系统监控在每次循环开始时都会通过 `usleep` 挂起当前线程，该函数的参数是微秒。

1. 初始的休眠时间是 20μs；
2. 最长的休眠时间是 10ms；
3. 当系统监控在 50 个循环中都没有唤醒 goroutine 时，休眠时间在每个循环都会倍增；

它除了会检查死锁之外，还会在循环中完成以下的工作：

- 运行计时器 — 获取下一个需要被触发的计时器；
- 轮询网络 — 获取需要处理的到期文件描述符；
- 抢占处理器 — 抢占运行时间较长的或者处于系统调用的 goroutine；
- 垃圾回收 — 在满足条件时触发垃圾收集回收内存；

### 运行计时器

在系统监控的循环中，通过 `runtime.nanotime` 和 `runtime.timeSleepUntil` 获取当前时间和计时器下一次需要唤醒的时间。

系统监控再下面的情况下会陷入休眠：

1. **当前调度器需要执行垃圾回收**。
2. 所有处理器都处于闲置状态时。
3. 没有需要触发的计时器。

休眠的时间会依据强制 GC 的周期 `forcegcperiod` 和计时器下次触发的时间确定。

### 轮询网络

如果上一次轮询网络已经过去了 10ms，那么系统监控还会在循环中轮询网络，检查是否有待执行的文件描述符。

```go
// runtime/proc.go
func sysmon() {
	// ...
	for {
		// ...
		lastpoll := int64(atomic.Load64(&sched.lastpoll))
		if netpollinited() && lastpoll != 0 && lastpoll+10*1000*1000 < now {
			atomic.Cas64(&sched.lastpoll, uint64(lastpoll), uint64(now))
			list := netpoll(0)
			if !list.empty() {
				incidlelocked(-1)
				injectglist(&list)
				incidlelocked(1)
			}
		}
		// ...
	}
}
```

非阻塞地调用 `runtime.netpoll` 检查待执行的文件描述符并通过 `runtime.injectglist` 将所有处于就绪状态的 goroutine 加入全局运行队列中：

```go
func injectglist(glist *gList) {
	if glist.empty() {
		return
	}
	lock(&sched.lock)
	var n int
	for n = 0; !glist.empty(); n++ {
		gp := glist.pop()
		casgstatus(gp, _Gwaiting, _Grunnable)
		globrunqput(gp)
	}
	unlock(&sched.lock)
	for ; n != 0 && sched.npidle != 0; n-- {
		startm(nil, false)
	}
	*glist = gList{}
}
```

该函数会将所有 goroutine 的状态从 `_Gwaiting` 切换至 `_Grunnable` 并加入全局运行队列等待运行，如果当前程序中存在空闲的 `p`，会通过 `runtime.startm` 启动线程来执行这些任务。

### 抢占处理

系统监控会在循环中调用 `runtime.retake` 抢占处于运行或者系统调用中的 `p`：

```go
// runtime/proc.go
func sysmon() {
    // ...
    for {
        // ...

        // retake P's blocked in syscalls
        // and preempt long running G's
        if retake(now) != 0 {
            idle = 0
        } else {
            idle++
        }
        // ...
    }
}
```

`reatke` 函数：

```go
type sysmontick struct {
	schedtick   uint32 // p 的调度次数
	schedwhen   int64  // p 的上次调度时间
	syscalltick uint32 // 系统调用次数
	syscallwhen int64  // 上次系统调用时间
}

func retake(now int64) uint32 {
	n := 0
	// Prevent allp slice changes. This lock will be completely
	// uncontended unless we're already stopping the world.
	lock(&allpLock)
    // 遍历所有的 p
	for i := 0; i < len(allp); i++ {
		pp := allp[i]
		if pp == nil {
			// This can happen if procresize has grown
			// allp but not yet created new Ps.
			continue
		}
		pd := &pp.sysmontick
		s := pp.status
		sysretake := false
        // 当 p 处于 _Prunning 或者 _Psyscall 状态时，如果上一次触发调度的时间已经过去了 10ms，通过 runtime.preemptone 抢占当前 p
		if s == _Prunning || s == _Psyscall {
			// Preempt G if it's running on the same schedtick for
			// too long. This could be from a single long-running
			// goroutine or a sequence of goroutines run via
			// runnext, which share a single schedtick time slice.
			t := int64(pp.schedtick) // schedtick 表示 p 的调度次数
			if int64(pd.schedtick) != t {
				pd.schedtick = uint32(t)
				pd.schedwhen = now // schedwhen 表示 p 上次调度时间
			} else if pd.schedwhen+forcePreemptNS <= now { // 上一次触发调度的时间已经超过了 10ms
				preemptone(pp)
				// In case of syscall, preemptone() doesn't
				// work, because there is no M wired to P.
				sysretake = true
			}
		}
        // 当 p 处于 _Psyscall 状态时
		if s == _Psyscall {
            // ...
			// On the one hand we don't want to retake Ps if there is no other work to do,
			// but on the other hand we want to retake them eventually
			// because they can prevent the sysmon thread from deep sleep.
            // 判断当亲 p 的运行队列是否为空
            // 是否存在空闲的 p
            // 系统调用时间是否超过了 10ms
			if runqempty(pp) && sched.nmspinning.Load()+sched.npidle.Load() > 0 && pd.syscallwhen+10*1000*1000 > now {
				continue
			}
			// Drop allpLock so we can take sched.lock.
			unlock(&allpLock)
			// Need to decrement number of idle locked M's
			// (pretending that one more is running) before the CAS.
			// Otherwise the M from which we retake can exit the syscall,
			// increment nmidle and report deadlock.
			incidlelocked(-1)
			trace := traceAcquire()
			if atomic.Cas(&pp.status, s, _Pidle) {
				if trace.ok() {
					trace.ProcSteal(pp, false)
					traceRelease(trace)
				}
				n++
				pp.syscalltick++
				handoffp(pp)
			} else if trace.ok() {
				traceRelease(trace)
			}
			incidlelocked(1)
			lock(&allpLock)
		}
	}
	unlock(&allpLock)
	return uint32(n)
}
```
1. 当处 `p` 处于 `_Prunning` 或者 `_Psyscall` 状态时，如果上一次触发调度的时间已经过去了 10ms，通过 `runtime.preemptone` 抢占当前 `p`；
2. 当 `p` 处于 `_Psyscall` 状态时，在满足以下两种情况下会调用 `runtime.handoffp` 让出 `p` 的使用权：
   - **当 `p` 的运行队列不为空或者不存在空闲的 `p` 时**（运行队列不为空说明还有 goroutine 等待运行；没有空闲的 `p` 说明有点忙）；
   - **当系统调用时间超过了 10ms 时**；
3. 被抢占的 goroutine，状态从 `Grunning` 变为 `Grunnable`，然后被放回全局运行队列。

`preemptone` 函数的主要流程：

1. 设置当前 goroutine 的抢占标志 `gp.preempt` 为 `true`；
2. 异步抢占,发送 `SIGURG` 信号;
3. 信号处理程序保存上下文;
4. 执行调度;

`handoffp` 函数的主要流程：

1. 解除 `p` 和 `m` 的绑定关系，设置 `p` 状态为` _Pidle`；
2. 启动一个新的 `m`，绑定到 `p` 上；
3. 切换到新的 `m` 运行，执行调度；

#### 信号抢占的流程

```
sysmon(检测超时) → preemptone() → preemptM() → 发送 SIGURG → 
信号处理程序 → asyncPreempt → asyncPreempt2 → mcall(gopreempt_m) → 抢占 goroutine → 执行调度
```

### 垃圾回收

在最后，系统监控还会决定是否需要触发强制垃圾回收，`runtime.sysmon` 会构建 `runtime.gcTrigger` 并调用 `runtime.gcTrigger.test` 方法判断是否需要触发垃圾回收：

```go
func sysmon() {
	// ...
	for {
		// ...
		if t := (gcTrigger{kind: gcTriggerTime, now: now}); t.test() && atomic.Load(&forcegc.idle) != 0 {
			lock(&forcegc.lock)
			forcegc.idle = 0
			var list gList
			list.push(forcegc.g)
			injectglist(&list)
			unlock(&forcegc.lock)
		}
		// ...
	}
}
```

如果需要触发垃圾回收，会**将用于垃圾回收的 goroutine 加入全局队列，让调度器选择合适的 `p` 去执行**。