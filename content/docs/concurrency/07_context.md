---
title: Context
weight: 7
---

Go 1.7 版本中正式引入新标准库 `context`。主要的作用是在在一组 goroutine 之间传递共享的值、取消信号、deadline 等。

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}
```

- `Deadline` — 返回当前 context 的截止时间。
- `Done` — 返回一个只读的 channel，可用于识别当前 channel 是否已经被关闭，其原因可能是到期，也可能是被取消了。多次调用 `Done` 方法会返回同一个 channel。
- `Err` — 返回当前 context 被关闭的原因。
  - 如果 context 被取消，会返回 `Canceled` 错误。
  - 如果 context 超时，会返回 `DeadlineExceeded` 错误。
- `Value` — 返回当前 context 对应所存储的 context信息，可以用来传递请求特定的数据。

创建 context：

- `Background`：创建一个空的 context，一般用在主函数、初始化、测试以及创建 root context 的时候。
- `TODO`：创建一个空的 context，不知道要传递一些什么上下文信息的时候，就用这个。
- `WithCancel`：基于 parent context 创建一个可以取消的新 context。
- `WithTimeout`：基于 parent context 创建一个具有**超时时间**的新 context。
- `WithDeadline`：和 `WithTimeout` 一样，只不过参数是**截止时间**（超时时间加上当前时间）。
- `WithValue`：基于某个 context 创建并存储对应的上下文信息。

最常用的场景，使用 context 来取消一个 goroutine 的运行：

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())

    go func() {
        defer func() {
            fmt.Println("goroutine exit")
        }()

        for {
            select {
            case <-ctx.Done():
                return
            default:
                time.Sleep(time.Second)
            }
        }
    }()

    time.Sleep(time.Second)
    cancel()
    time.Sleep(2 * time.Second)
}
```

可以多个 goroutine 同时订阅 `ctx.Done()` 管道中的消息，一旦接收到取消信号就立刻停止当前正在执行的工作。

## 原理

context 的最大作用就是在一组 goroutine 构成的树形结构中对信号进行同步，以减少计算资源的浪费。

例如，Go 的 HTTP server，处理每一个请求，都是启动一个单独的 goroutine，处理过程中还会启动新的 goroutine 来访问数据库和其他服务。而 context 在不同 Goroutine 之间可以同步请求特定数据、取消信号以及处理
请求的截止日期。

![context](https://raw.githubusercontent.com/shipengqi/illustrations/505be84aff62a8b94292bf3468f0d9a1c8c049cf/go/context.png)

每一个 context 都会从 root goroutine 一层层传递到底层。context 可以在上层 goroutine 执行出现错误时，将信号及时同步给下层。

### WithCancel

```go
// src/context/context.go#L235
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := withCancel(parent)
	return c, func() { c.cancel(true, Canceled, nil) }
}

func withCancel(parent Context) *cancelCtx {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := &cancelCtx{}
	// 构建 父子 context 之间的关联，当 父 context 被取消时，子 context 也会被取消
	c.propagateCancel(parent, c)
	return c
}

func (c *cancelCtx) propagateCancel(parent Context, child canceler) {
	c.Context = parent

	done := parent.Done()
	if done == nil { // parent context 是个空 context
		return // parent is never canceled
	}

	select {
	case <-done:
		// parent context 已经被取消，child 也会立刻被取消
		child.cancel(false, parent.Err(), Cause(parent))
		return
	default:
	}

    // 找到可以取消的 parent context
	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil {
            // parent context 已经被取消，child 也会立刻被取消
			child.cancel(false, p.err, p.cause)
		} else {
			// 将 child 加入到 parent 的 children 列表中
			// 等待 parent 释放取消信号
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
		return
	}

	if a, ok := parent.(afterFuncer); ok {
		// parent implements an AfterFunc method.
		c.mu.Lock()
		stop := a.AfterFunc(func() {
			child.cancel(false, parent.Err(), Cause(parent))
		})
		c.Context = stopCtx{
			Context: parent,
			stop:    stop,
		}
		c.mu.Unlock()
		return
	}

	goroutines.Add(1)
	// 没有找到可取消的 parent context
	// 运行一个新的 goroutine 同时监听 parent.Done() 和 child.Done() 两个 channel
	go func() {
		select {
		case <-parent.Done():
			// 在 parent.Done() 关闭时调用 child.cancel 取消 子 context
			child.cancel(false, parent.Err(), Cause(parent))
		case <-child.Done(): // 这个空的 case 表示如果子节点自己取消了，那就退出这个 select，父节点的取消信号就不用管了。
		                     // 如果去掉这个 case，那么很可能父节点一直不取消，这个 goroutine 就泄漏了
		}
	}()
}


func (c *cancelCtx) Done() <-chan struct{} {
	c.mu.Lock()
	// 有调用了 Done() 方法的时候才会被创建
	if c.done == nil {
		c.done = make(chan struct{})
	}
	// 返回的是一个只读的 channel
	// 这个 channel 不会被写入数据，直接调用读这个 channel，协程会被 block 住。
	// 一般通过搭配 select 来使用。一旦关闭，就会立即读出零值。
	d := c.done
	c.mu.Unlock()
	return d
}
```

`propagateCancel` 的作用就是向上寻找可以“挂靠”的“可取消”的 context，并且“挂靠”上去。这样，调用上层 `cancel` 方法的时候，就可以层层传递，
将那些挂靠的子 context 同时“取消”。

`cancelCtx.cancel` 会关闭 context 中的 channel 并向所有的 子 context 同步取消信号：

```go
func (c *cancelCtx) cancel(removeFromParent bool, err, cause error) {
	// ...
	if d == nil {
		c.done.Store(closedchan)
	} else {
		close(d)
	}
	// 遍历所有 子 context，取消所有 子 context
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		child.cancel(false, err, cause)
	}
    // 将子节点置空
    c.children = nil
	// ...
    if removeFromParent {
		// 从父节点中移除自己 
		removeChild(c.Context, c)
	}
}
```

### WithTimeout 和 WithDeadline

`WithTimeout` 和 `WithDeadline` 创建的 context 也都是可以被取消的。

`WithTimeout` 和 `WithDeadline` 创建的是 `timeCtx`，`timerCtx` 基于 `cancelCtx`，多了一个 `time.Timer` 和 `deadline`：

```go
type timerCtx struct {
	cancelCtx
	timer *time.Timer // Under cancelCtx.mu.
	
	deadline time.Time
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
	// 直接调用 cancelCtx 的取消方法
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		// 从父节点中删除子节点
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		// 关掉定时器，这样，在deadline 到来时，不会再次取消
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}
```

`WithTimeout` 实际就时调用了 `WithDeadline`，传入的 deadline 是当前时间加上 timeout 的时间：

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}
```

`WithDeadline` 的实现：

```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	return WithDeadlineCause(parent, d, nil)
}

func WithDeadlineCause(parent Context, d time.Time, cause error) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	// 如果 parent context 的 deadline 早于指定时间。直接构建一个可取消的 context
	// 原因是一旦 parent context 超时，自动调用 cancel 函数，子节点也会随之取消
	// 所以没有必要再处理 子 context 的计时器
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		return WithCancel(parent)
	}
	c := &timerCtx{
		deadline: d,
	}
	// 构建一个 cancelCtx，挂靠到一个可取消的 parent context 上
	// 也就是说一旦 parent context 取消了，这个子 context 随之取消。
	c.cancelCtx.propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
        // 超过了截止日期，直接取消
		c.cancel(true, DeadlineExceeded, cause)
		return c, func() { c.cancel(false, Canceled, nil) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		// 到了截止时间，timer 会自动调用 cancel 函数取消
		c.timer = time.AfterFunc(dur, func() {
			// 传入错误 DeadlineExceeded
			c.cancel(true, DeadlineExceeded, cause)
		})
	}
	return c, func() { c.cancel(true, Canceled, nil) }
}
```

> 如果要创建的这个 子 context 的 deadline 比 parent context 的要晚，parent context 到时间了会自动取消，子 context 也会取消，
> 导致 子 context 的 deadline 时间还没到就会被取消

### WithValue

```go
// src/context/context.go#L713
func WithValue(parent Context, key, val any) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}

type valueCtx struct {
	Context
	key, val interface{}
}

func (c *valueCtx) Value(key any) any {
	if c.key == key {
		return c.val
	}
	// 如果 valueCtx 中存储的键值对与传入的参数不匹配
	// 就会从 parent context 中查找该键对应的值直到某个 parent context 中返回 nil 或者查找到对应的值。
	return value(c.Context, key)
}
```
