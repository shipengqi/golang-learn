---
title: ErrGroup
weight: 12
---

# ErrGroup

Go 的扩展库 `golang.org/x/sync` 提供了 `errgroup` 包，它是基于 `WaitGroup` 实现的，功能上和 `WaitGroup` 类似，不过可以通过上下文取消，控制并发数量，还能返回错误。

## 使用

最简单的使用方式：

```go
package main

import (
    "errors"
    "fmt"
    "time"

    "golang.org/x/sync/errgroup"
)

func main() {
    var g errgroup.Group
    // g, ctx := errgroup.WithContext(context.Background())
	
    g.Go(func() error {
        time.Sleep(5 * time.Second)
        fmt.Println("exec 1")
        return nil
    })

    g.Go(func() error {
        time.Sleep(10 * time.Second)
        fmt.Println("exec 2")
        return errors.New("failed to exec 2")
    })

    if err := g.Wait(); err == nil {
        fmt.Println("exec done")
    } else {
        fmt.Println("failed: ", err)
    }
}
```

- `errgroup.WithContext` 返回一个 `Group` 实例，同时还会返回一个使用 `context.WithCancel(ctx)` 生成的新 `Context`。
- `Group.Go` 方法能够创建一个 goroutine 并在其中执行传入的函数
- `Group.Wait` 会等待所有 goroutine 全部返回，该方法的不同返回结果也有不同的含义： 
  - 如果返回 `error`，那么这组 goroutine 至少有一个返回了 `error`。
  - 如果返回 `nil`，表示所有 goroutine 都成功执行。


### 限制 goroutine 的并发数量

```go
package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group
	g.SetLimit(2)
	g.TryGo(func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("exec 1")
		return nil
	})

	g.TryGo(func() error {
		time.Sleep(10 * time.Second)
		fmt.Println("exec 2")
		return errors.New("failed to exec 2")
	})

	if err := g.Wait(); err == nil {
		fmt.Println("exec done")
	} else {
		fmt.Println("failed: ", err)
	}
}
```

- `Group.SetLimit` 设置并发数量。
- `Group.TryGo` 替换 `Group.Go` 方法。

  
## 原理

`errgroup.Group` 的结构体：

```go
type Group struct {
	cancel func(error) // 创建 context.Context 时返回的取消函数，用于在多个 goroutine 之间同步取消信号

	wg sync.WaitGroup // 用于等待一组 goroutine 的完成

	sem chan token // 利用这个 channel 的缓冲区大小，来控制并发的数量

	errOnce sync.Once // 保证只接收一个 goroutine 返回的错误
	err     error
}
```

`errgroup` 的实现很简单：

```go
func (g *Group) done() {
	if g.sem != nil {
		// 从 channel 获取一个值，释放资源
		<-g.sem
	}
	//  WaitGroup 并发数量 -1
	g.wg.Done()
}

// golang/sync/errgroup/errgroup.go
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := withCancelCause(ctx)
	return &Group{cancel: cancel}, ctx
}

func (g *Group) Go(f func() error) {
	// g.sem 的值不为 nil，说明调用了 SetLimit 设置并发数量
	if g.sem != nil {
		// 尝试从 channel 发送一个值
		// - 发送成功，缓冲区还没有满，意味着并发数还没有达到 SetLimit 设置的数量
		// - 发送不成功，缓冲区已满，阻塞在这里，等待其他 goroutine 释放一个资源
		g.sem <- token{}
	}

    // 调用 WaitGroup.Add 并发数量 +1
	g.wg.Add(1)
	// 创建新的 goroutine 运行传入的函数
	go func() {
		defer g.done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				// 返回错误时，调用 context 的 cancel 并对 err 赋值
				g.err = err
				if g.cancel != nil {
					g.cancel(g.err)
				}
			})
		}
	}()
}

func (g *Group) Wait() error {
	// 只是调用了 WaitGroup.Wait
	g.wg.Wait()
	// 在所有 goroutine 完成时，取消 context
	if g.cancel != nil {
		g.cancel(g.err)
	}
	return g.err
}
```

限制 goroutine 并发数量的实现：

```go
func (g *Group) SetLimit(n int) {
	// 小于 0 时，直接给 g.sem 赋值为 nil，表示不限制并发数量
	if n < 0 {
		g.sem = nil
		return
	}
	// 已有 goroutine 运行时，不能在设置并发数量
	if len(g.sem) != 0 {
		panic(fmt.Errorf("errgroup: modify limit while %v goroutines in the group are still active", len(g.sem)))
	}
	// 创建一个大小为 n 的有缓冲 channel
	g.sem = make(chan token, n)

}
func (g *Group) TryGo(f func() error) bool {
	// 与 Go 方法的主要区别，就在对 sem 的处理上
	// 尝试获取资源，当无法拿到资源时，直接返回 false，表示执行失败
	if g.sem != nil {
		select {
		case g.sem <- token{}:
			// Note: this allows barging iff channels in general allow barging.
		default:
			return false
		}
	}

    // 调用 WaitGroup.Add 并发任务 +1
	g.wg.Add(1)
	go func() {
		defer g.done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel(g.err)
				}
			})
		}
	}()
	return true
}
```