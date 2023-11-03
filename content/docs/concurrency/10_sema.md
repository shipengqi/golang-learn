---
title: 信号量
weight: 10
---

mutex 和 waitgroup 不保存 goroutine 的信息，所以通过信号量（底层依然是 gopark goready 来实现让出cpu，和唤醒）来通知所有阻塞的 goroutine。
但是 channel 需要在 goroutine 之间传递数据的，需要拷贝内存，仅通过信号量无法实现，信号量只能做到通知，所以 channel 保存了 goroutine 的信息。

# 信号量

Go 通过信号量来控制 goroutine 的阻塞和唤醒，例如 `Mutex` 结构体重的 `sema`：

```go
type Mutex struct {
    state int32
	sema  uint32
}
```