---
title: Once
weight: 5
---

# Once

Go 标准库中 `sync.Once` 可以保证 Go 程序运行期间的某段代码只会执行一次。常常用于单例对象的初始化场景。

`sync.Once` 只有一个对外唯一暴露的方法 `Do`，可以多次调用，但是只第一次调用时会执行一次。

```go
func main() {
    o := &sync.Once{}
    for i := 0; i < 10; i++ {
        o.Do(func() {
            fmt.Println("only once")
        })
    }
}
```

运行：
```
$ go run main.go
only once
```

## 原理

`sync.Once` 的实现：

```go
// src/sync/once.go
type Once struct {
	done uint32
	m    Mutex
}

func (o *Once) Do(f func()) {
	// 如果传入的参数 f 已经执行过，直接返回
    if atomic.LoadUint32(&o.done) == 0 {
        o.doSlow(f)
    }
}

func (o *Once) doSlow(f func()) {
	// 为当前 goroutine 加锁
    o.m.Lock()
    defer o.m.Unlock()
    if o.done == 0 {
		// 将 done 设置为 1
        defer atomic.StoreUint32(&o.done, 1)
		// 执行参数 f
        f()
    }
}
```

`sync.Once` 使用互斥锁和原子操作实现了某个函数在程序运行期间只能执行一次的语义。

使用互斥锁，同时利用双检查的机制（double-checking），再次判断 `o.done` 是否为 0，如果为 0，则是第一次执行，执行完毕后，就将 `o.done` 设置为 1，然后释放锁。

即使有多个 goroutine 同时进入了 `doSlow` 方法，因为双检查的机制，后续的 goroutine 会看到 `o.done` 的值为 1，也不会再次执行 `f`。