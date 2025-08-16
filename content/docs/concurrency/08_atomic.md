---
title: 原子操作
weight: 8
---

原子操作就是执行过程中不能被中断的操作。

Go 的标准库 `sync/atomic` 提供了一些实现原子操作的方法：

- Add
- CompareAndSwap（简称 CAS）
- Load
- Swap
- Store

这些函数针对的数据类型有：

- `int32`
- `int64`
- `uint32`
- `uint64`
- `uintptr`
- `unsafe` 包中的 `Pointer`

以 `Add` 为例，上面类型对应的原子操作函数为：

- `func AddInt32(addr *int32, delta int32) (new int32)`
- `func AddInt64(addr *int64, delta int64) (new int64)`
- `func AddUint32(addr *uint32, delta uint32) (new uint32)`
- `func AddUint64(addr *uint64, delta uint64) (new uint64)`
- `func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)`

> `unsafe.Pointer` 类型，并未提供进行原子加法操作的函数。

`sync/atomic` 包还提供了一个名为 `Value` 的类型，它可以被用来存储（Store）和加载（Load）任意类型的值。

它只有两个指针方法：

- `Store` 
- `Load`。

**尽量不要向原子值中存储引用类型的值**。

```go
var box6 atomic.Value
v6 := []int{1, 2, 3}
box6.Store(v6)
v6[1] = 4 // 此处的操作不是并发安全的
```

上面的代码 `v6[1] = 4` 绕过了原子值而进行了非并发安全的操作。可以改为：

```go
store := func(v []int) {
    replica := make([]int, len(v))
    copy(replica, v)
    box6.Store(replica)
}
store(v6)
v6[2] = 5
```

## 使用

### 互斥锁与原子操作

区别：

- **互斥锁是用来保护临界区，原子操作用于对一个变量的更新保护**。
- 互斥锁由操作系统的调度器实现，原子操作由底层硬件指令直接提供支持

对于一个变量更新的保护，原子操作通常会更有效率，并且更能利用计算机多核的优势。而互斥锁保护的共享资源每次只给一个线程使用，其它线程阻塞，用完后再把资源转让给其它线程。

使用互斥锁实现并发计数：

```go
func MutexAdd() {
	var a int32 =  0
	var wg sync.WaitGroup
	var mu sync.Mutex
	start := time.Now()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			a += 1
			mu.Unlock()
		}()
	}
	wg.Wait()
	timeSpends := time.Now().Sub(start).Nanoseconds()
    fmt.Printf("mutex value %d, spend time: %v\n", a, timeSpends)
}
```

使用原子操作替换互斥锁：

```go
func AtomicAdd() {
	var a int32 =  0
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt32(&a, 1)
		}()
	}
	wg.Wait()
	timeSpends := time.Now().Sub(start).Nanoseconds()
    fmt.Printf("atomic value %d, spend time: %v\n", atomic.LoadInt32(&a), timeSpends)
}
```

运行后得到的结果：

```
mutex value 10000, spend time: 5160800 
atomic value 10000, spend time: 2577300
```

原子操作节省了大概一半的时间。

### 利用 CAS 实现自旋锁

```go
func addValue(v int32)  {
	for {
		// 在进行读取 value 的操作的过程中,其他对此值的读写操作是可以被同时进行的,那么这个读操作很可能会读取到一个只被修改了一半的数据.
		// 因此要使用原子读取
		old := atomic.LoadInt32(&value)
		if atomic.CompareAndSwapInt32(&value, old, old + v) {
			break
		}
	}
}
```

在高并发的情况下，单次 CAS 的执行成功率会降低，因此需要配合循环语句 `for`，形成一个 `for+atomic` 的类似自旋乐观锁。

**自旋锁的使用场景**：

1. 读多写少的场景，线程能够很快地获得锁，则自旋锁非常有效，因为它避免了线程调度和上下文切换的开销。如果有大量的写操作，CAS 操作无法获取到锁，线程会在不断的自旋中消耗大量的 CPU 时间。线程需要反复尝试获取锁，而不释放 CPU，这可能导致性能下降。
2. 自旋锁适合用在加锁粒度很小的场景，锁的持有时间非常短。如果锁定时间较长，使用自旋锁可能导致线程长时间占用 CPU。


### ABA 问题

使用 CAS，会有 ABA 问题，ABA 问题是什么？

例如，一个 goroutine a 从内存位置 V 中取出 1，这时候另一个 goroutine b 也从内存位置 V 中取出 1，并且 goroutine b 将 V 位置的值更新为 0，接着又将 V 位置的值改为 1，这时候 goroutine a 进行 CAS 操作发现位置 V 的值仍然是 1，然后 goroutine a 操作成功。虽然 goroutine a 的 CAS 操作成功，但是这个值其实已经被修改过。

可以给变量附加时间戳、版本号等信息来解决。


