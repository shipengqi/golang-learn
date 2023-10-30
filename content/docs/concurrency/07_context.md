---
title: Context
weight: 7
---

# Context

使用 `WaitGroup` 值的时候，我们最好用**先统一 `Add`，再并发 `Done`，最后 `Wait`** 的标准模式来构建协作流程。如果在调用
该值的 `Wait` 方法的同时，为了增大其计数器的值，而并发地调用该值的 `Add` 方法，那么就很可能会引发 panic。

但是**如果，我们不能在一开始就确定执行子任务的 goroutine 的数量，那么使用 `WaitGroup` 值来协调它们和分发子任务的 goroutine，就是有一定风险的**。一个解决方案是：**分批地启用执行子任务的 goroutine**。

只要我们在严格遵循上述规则的前提下，分批地启用执行子任务的 goroutine，就肯定不会有问题。

```go
func coordinateWithWaitGroup() {
    total := 12
    stride := 3
    var num int32
    fmt.Printf("The number: %d [with sync.WaitGroup]\n", num)
    var wg sync.WaitGroup
    for i := 1; i <= total; i = i + stride {
        wg.Add(stride)
        for j := 0; j < stride; j++ {
            go addNum(&num, i+j, wg.Done)
        }
        wg.Wait()
    }
    fmt.Println("End.")
}
```

### 使用 `context` 包中的程序实体，实现一对多的 goroutine 协作流程

用 `context` 包中的函数和 `Context` 类型作为实现工具，实现 `coordinateWithContext` 的函数。这个函数应该具有上
面 `coordinateWithWaitGroup` 函数相同的功能。

```go
func coordinateWithContext() {
 total := 12
 var num int32
 fmt.Printf("The number: %d [with context.Context]\n", num)
 cxt, cancelFunc := context.WithCancel(context.Background())
 for i := 1; i <= total; i++ {
  go addNum(&num, i, func() {
   if atomic.LoadInt32(&num) == int32(total) {
    cancelFunc()
   }
  })
 }
 <-cxt.Done()
 fmt.Println("End.")
}
```

先后调用了 `context.Background` 函数和 `context.WithCancel` 函数，并得到了一个可撤销的 `context.Context` 类型的值
（由变量 `cxt` 代表），以及一个 `context.CancelFunc`类型的撤销函数（由变量 `cancelFunc` 代表）。

注意我给予 `addNum` 函数的最后一个参数值。它是一个匿名函数，其中只包含了一条 `if` 语句。这条 `if` 语句会**原子地**加载
`num` 变量的值，并判断它是否等于 `total` 变量的值。

如果两个值相等，那么就调用 `cancelFunc` 函数。其含义是，如果所有的 `addNum` 函数都执行完毕，那么就立即通知分发子任务
的 goroutine。

**这里分发子任务的 goroutine，即为执行 `coordinateWithContext` 函数的 goroutine**。它在执行完 `for` 语句后，会
立即调用 `cxt` 变量的 `Done` 函数，并试图针对该函数返回的通道，进行接收操作。

一旦 `cancelFunc` 函数被调用，针对该通道的接收操作就会马上结束，所以，这样做就可以实现“等待所有的 `addNum` 函数都执
行完毕”的功能。

### context.Context 类型

`Context` 类型的值（以下简称 `Context` 值）是可以繁衍的，这意味着我们可以通过一个 `Context` 值产生出任意个子值。这些子值
可以携带其父值的属性和数据，也可以响应通过其父值传达的信号。

正因为如此，所有的 `Context` 值共同构成了一颗代表了上下文全貌的树形结构。这棵树的**树根（或者称上下文根节点）是一个已经
在 `context` 包中预定义好的 `Context` 值**，它是**全局唯一**的。通过调用 `context.Background` 函数，我们就可以获取到
它（在 `coordinateWithContext` 函数中就是这么做的）。

注意一下，这个**上下文根节点仅仅是一个最基本的支点，它不提供任何额外的功能**。也就是说，它既不可以被撤销（`cancel`），
也不能携带任何数据。

`context` 包中还包含了**四个用于繁衍 `Context` 值的函数，即：`WithCancel`、`WithDeadline`、`WithTimeout` 和 `WithValue`**。

这些函数的第一个参数的类型都是 `context.Context`，而名称都为 `parent`。顾名思义，**这个位置上的参数对应的都是它们将会产生
的 `Context` 值的父值**。

**`WithCancel` 函数用于产生一个可撤销的 parent 的子值**。

在 `coordinateWithContext` 函数中，通过调用该函数，获得了一个衍生自上下文根节点的 `Context` 值，和一个用于触发撤销信号的函数。

`WithDeadline` 函数和 `WithTimeout` 函数则都可以被用来产生一个会**定时撤销**的 `parent` 的子值。至于 `WithValue` 函数，
我们可以通过调用它，产生一个会携带额外数据的 `parent` 的子值。

### “可撤销的”在 context 包中代表着什么？“撤销”一个 Context 值又意味着什么？

这需要从 `Context` 类型的声明讲起。这个接口中有两个方法与“撤销”息息相关。`Done` 方法会返回一个元素类型为 `struct{}` 的接
收通道。不过，这个接收通道的用途并不是传递元素值，而是**让调用方去感知“撤销”当前Context值的那个信号**。

一旦当前的 `Context` 值被撤销，这里的接收通道就会被立即关闭。我们都知道，对于一个未包含任何元素值的通道来说，它的关闭会
使任何针对它的接收操作立即结束。

正因为如此，在 `coordinateWithContext` 函数中，基于调用表达式 `cxt.Done()` 的接收操作，才能够起到感知撤销信号的作用。

### 撤销信号是如何在上下文树中传播的

`context`包的 `WithCancel` 函数在被调用后会产生两个结果值。第一个结果值就是那个可撤销的 `Context` 值，而第二个结果值则是
用于触发撤销信号的函数。

在撤销函数被调用之后，对应的 `Context` 值会先关闭它内部的接收通道，也就是它的 `Done` 方法会返回的那个通道。

然后，它会向它的所有子值（或者说子节点）传达撤销信号。这些子值会如法炮制，把撤销信号继续传播下去。最后，这个 `Context` 值会
断开它与其父值之间的关联。

**通过调用 `context.WithValue` 函数得到的 `Context` 值是不可撤销的**。

### 怎样通过 Context 值携带数据

**`WithValue` 函数在产生新的 `Context` 值（以下简称含数据的 `Context` 值）的时候需要三个参数，即：父值、键和值**。
“字典对于键的约束”类似，这里**键的类型必须是可判等**的。

原因很简单，当我们从中获取数据的时候，它需要根据给定的键来查找对应的值。不过，这种 `Context` 值并不是用字典来存储键和值的，
后两者只是被简单地存储在前者的相应字段中而已。