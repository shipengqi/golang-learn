---
title: 读写锁
weight: 2
---

读写互斥锁 `sync.RWMutex` 是细粒度的互斥锁，一般来说有几种情况：

- 读锁之间不互斥
- 写锁之间是互斥的
- 写锁与读锁是互斥的

`sync.RWMutex` 类型中的 `Lock` 方法和 `Unlock` 方法用于对写锁进行锁定和解锁，`RLock` 方法和 `RUnlock` 方法则分别用于对读锁进行锁定和解锁。

## 原理

```go
type RWMutex struct {
	w           Mutex  // 复用互斥锁提供的能力，解决多个 writer 的竞争
	writerSem   uint32 // writer 的信号量
	readerSem   uint32 // reader 的信号量
	readerCount atomic.Int32 // 正在执行的 reader 的数量
	readerWait  atomic.Int32 // 当写操作被阻塞时需要等待 read 完成的 reader 的数量
}

const rwmutexMaxReaders = 1 << 30
```

`rwmutexMaxReaders`：定义了最大的 reader 数量。

### RLock 和 RUnlock

移除了 race 等无关紧要的代码：
```go
func (rw *RWMutex) RLock() {
	if rw.readerCount.Add(1) < 0 {
		// rw.readerCount 是负值，意味着此时有其他 goroutine 获得了写锁
		// 当前 goroutine 就会调用 runtime_SemacquireRWMutexR 陷入休眠等待锁的释放
		runtime_SemacquireRWMutexR(&rw.readerSem, false, 0)
	}
}

func (rw *RWMutex) RUnlock() {
	// 先减少正在读资源的 readerCount 整数
	// 如果返回值大于等于零，读锁直接解锁成功
	if r := rw.readerCount.Add(-1); r < 0 {
		// 如果返回值小于零，有一个正在执行的写操作
		rw.rUnlockSlow(r)
	}
}

func (rw *RWMutex) rUnlockSlow(r int32) {
	// 减少 readerWait
	if rw.readerWait.Add(-1) == 0 {
        // 在所有读操作都被释放之后触发写操作的信号量 writerSem，
		// 该信号量被触发时，调度器就会唤醒尝试获取写锁的 goroutine。
		runtime_Semrelease(&rw.writerSem, false, 1)
	}
}
```

### Lock 和 Unlock

移除了 race 等无关紧要的代码：

```go
func (rw *RWMutex) Lock() {

	// 写锁加锁，其他 goroutine 在获取写锁时会进入自旋或者休眠
	rw.w.Lock()
	// 将 readerCount 变为负数，阻塞后续的读操作
	r := rw.readerCount.Add(-rwmutexMaxReaders) + rwmutexMaxReaders
	// 如果仍然有其他 goroutine 持有互斥锁的读锁，当前 goroutine 会调用 runtime_SemacquireRWMutex 进入休眠状态等待所有读锁所有者执
	// 行结束后释放 writerSem 信号量将当前协程唤醒
	if r != 0 && rw.readerWait.Add(r) != 0 {
		runtime_SemacquireRWMutex(&rw.writerSem, false, 0)
	}
}

func (rw *RWMutex) Unlock() {
	// 将 readerCount 变回正数，释放读锁
	r := rw.readerCount.Add(rwmutexMaxReaders)
	if r >= rwmutexMaxReaders {
		race.Enable()
		fatal("sync: Unlock of unlocked RWMutex")
	}
	// 通过 for 循环释放所有因为获取读锁而陷入等待的 goroutine
	for i := 0; i < int(r); i++ {
		runtime_Semrelease(&rw.readerSem, false, 0)
	}
	// 释放写锁
	rw.w.Unlock()
}
```

获取写锁时会**先阻塞写锁的获取，后阻塞读锁的获取**，这种策略能够保证读操作不会被连续的写操作**饿死**。
