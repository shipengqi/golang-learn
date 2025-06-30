---
title: WaitGroup
weight: 3
---

`sync.WaitGroup` 可以等待一组 goroutine 的返回，常用于处理批量的并发任务。它是并发安全的。

## 使用

并发发送 HTTP 请求的示例：

```go
requests := []*Request{...}
wg := &sync.WaitGroup{}
wg.Add(len(requests))

for _, request := range requests {
    go func(r *Request) {
        defer wg.Done()
        // res, err := service.call(r)
    }(request)
}
wg.Wait()
```

`WaitGroup` 提供了三个方法：

- `Add`：用来设置 `WaitGroup` 的计数值。
- `Done`：用来将 `WaitGroup` 的计数值减 1，其实就是调用了 `Add(-1)`。
- `Wait`：调用这个方法的 `goroutine` 会一直阻塞，直到 `WaitGroup` 的计数值变为 0。

**不要把 `Add` 和 `Wait` 方法的调用放在不同的 goroutine 中执行**，以免 `Add` 还未执行，`Wait` 已经退出：

```go
var wg sync.WaitGroup
go func(){
	wg.Add(1)
	fmt.Println("test")
}()

wg.Wait()
fmt.Println("exit.")
```

### 1.25 Go 方法

新增 Go 方法：

```go
func (wg *WaitGroup) Go(f func()) {
    wg.Add(1)
    go func() {
        defer wg.Done()
        f()
    }()
}
```

简单的封装了 `wg.Add(1)` 和 `wg.Done()`。

使用：

```go
for _, request := range requests {
    wg.Go(func() {
        // res, err := service.call(request)
    }
}
wg.Wait()
```

### sync.WaitGroup 类型值中计数器的值可以小于 0 么？

不可以。小于 0，会引发 panic。所以尽量不要传递负数给 `Add` 方法，只通过 `Done` 来给计数值减 1。

### sync.WaitGroup 可以复用么？

可以。但是必须在 `Wait` 方法返回之后才能被重新使用。否则会引发 panic。所以尽量不要重用 `WaitGroup`。新建一个 `WaitGroup` 不会带来多大的资源
开销，重用反而更容易出错。

### Wait 可以在多个 goroutine 调用多次么？

可以。当前 `sync.WaitGroup` 计数器的归零时，这些 goroutine 会被同时唤醒。

## 原理

`sync.WaitGroup` 结构体：

```go
// src/sync/waitgroup.go#L20
type WaitGroup struct {
	noCopy noCopy
	state1 [3]uint32
}
```

`noCopy` 是 go 1.7 开始引入的一个静态检查机制，它只是一个辅助类型：

```go
// src/sync/cond.go#L117
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
```

`tools/go/analysis/passes/copylock` 包中的分析器会在编译期间检查被拷贝的变量中是否包含 `noCopy` 或者实现了 `Lock` 和 `Unlock` 方法，如果包含该结构体或者实现了对应的方法就会报错：

```
$ go vet proc.go
./prog.go:10:10: assignment copies lock value to yawg: sync.WaitGroup
./prog.go:11:14: call of fmt.Println copies lock value: sync.WaitGroup
./prog.go:11:18: call of fmt.Println copies lock value: sync.WaitGroup
```

`state1` 包含一个总共占用 12 字节的数组，这个数组会存储当前结构体的状态，在 64 位与 32 位的机器上表现也非常不同。

![waitgroup-state1](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/waitgroup-state1.png)

`state` 方法用来从 `state1` 字段中取出它的状态和信号量。

```go
// 得到 state 的地址和信号量的地址
func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
    if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
        // 如果地址是 64bit 对齐的，数组前两个元素做 state，后一个元素做信号量
        return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
    } else {
        // 如果地址是 32bit 对齐的，数组后两个元素用来做 state，它可以用来做 64bit 的原子操作，第一个元素 32bit 用来做信号量
        return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
    }
}
```

`Add` 的实现：

```go
func (wg *WaitGroup) Add(delta int) {
    statep, semap := wg.state()
    // 高 32bit 是计数值 v，所以把 delta 左移 32，更新计数器 counter
    state := atomic.AddUint64(statep, uint64(delta)<<32)
    v := int32(state >> 32) // 当前计数值
    w := uint32(state) // waiter count

    if v < 0 {
        panic("sync: negative WaitGroup counter")
    }
	// 并发的 Add 会导致 panic
    if w != 0 && delta > 0 && v == int32(delta) {
        panic("sync: WaitGroup misuse: Add called concurrently with Wait")
    }
    if v > 0 || w == 0 {
        return
    }
	
    // 将 waiter 调用计数器归零，也就是 *statep 直接设置为 0 即可。
	// 通过 sync.runtime_Semrelease 唤醒处于等待状态的 goroutine。
    *statep = 0
    for ; w != 0; w-- {
        runtime_Semrelease(semap, false, 0)
    }
}


// Done 方法实际就是计数器减 1
func (wg *WaitGroup) Done() {
    wg.Add(-1)
}
```

`Wait` 方法的实现逻辑：不断检查 state 的值。如果其中的计数值变为了 0，那么说明所有的任务已完成，调用者不必再等待，直接返回。如果计数值大于 0，说明此时还有任
务没完成，那么调用者就变成了等待者，需要加入 waiter 队列，并且阻塞住自己。

```go
func (wg *WaitGroup) Wait() {
    statep, semap := wg.state()
    
    for {
        state := atomic.LoadUint64(statep)
        v := int32(state >> 32) // 当前计数值
        w := uint32(state) // waiter 的数量
        if v == 0 {
            // 如果计数值为 0, 调用这个方法的 goroutine 不必再等待，继续执行它后面的逻辑即可
            return
        }
        // 否则把 waiter 数量加 1。期间可能有并发调用 Wait 的情况，所以最外层使用了一个 for 循环
        if atomic.CompareAndSwapUint64(statep, state, state+1) {
            // 阻塞休眠等待
            runtime_Semacquire(semap)
            // 被唤醒，不再阻塞，返回
            return
        }
    }
}
```
