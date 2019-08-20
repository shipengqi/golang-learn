---
title: 并发编程
---

# 并发编程
```
Don’t communicate by sharing memory; share memory by communicating.
（不要通过共享内存来通信，而应该通过通信来共享内存。）
```
这是作为 Go 语言的主要创造者之一的 Rob Pike 的至理名言，这也充分体现了 Go 语言最重要的编程理念。

## Goroutines
`goroutine`可以简单理解为一个线程，但是它比线程更小，十几个`goroutine`可能体现在底层就是五六个线程，Go语言内部帮你实现
了这些`goroutine`之间的内存共享。执行`goroutine`只需极少的栈内存(大概是`4~5KB`)，当然会根据相应的数据伸缩。也正因为如此，
可同时运行成千上万个并发任务。`goroutine`比`thread`更易用、更高效、更轻便。我们程序运行的`main`函数在一个单独的`goroutine`中运行，叫做`main goroutine`。
在代码中可以使用`go`关键字创建`goroutine`。
```go
go f()
```

主函数返回时，所有`goroutine`都会被打断，程序退出。除了从主函数退出或者直接终止程序之外，没有其它的编程方法能够让一个`goroutine`来打断另一个的执行，但是之后可以看到一种方式来实
现这个目的，通过`goroutine`之间的通信来让一个`goroutine`请求其它的`goroutine`，使被请求`goroutine`自行结束执行。


### goroutines泄漏
`goroutines`被永远卡住，就会导致`goroutines`泄漏，例如当使用了无缓存的`channel`，`goroutines`因为`channel`的数据没有被接收而被卡住。
泄漏的`goroutines`不会被自动回收。

### 什么是主 goroutine，它与我们启用的其他 goroutine 有什么不同
```go
package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}
```
上面的代码会打印出什么内容？

回答是：不会有任何内容被打印出来。

Go 语言不但有着独特的并发编程模型，以及用户级线程 goroutine，还拥有强大的用于调度 goroutine、对接系统级线程的调度器。

这个调度器是 Go 语言运行时系统的重要组成部分，它主要负责统筹调配 Go 并发编程模型中的三个主要元素，即：G（goroutine 的缩写）、P（processor 的缩写）和 M（machine 的缩写）
M 指代的就是系统级线程。而 P 指的是一种可以承载若干个 G，且能够使这些 G 适时地与 M 进行对接，并得到真正运行的中介。

与一个进程总会有一个主线程类似，每一个独立的 Go 程序在运行时也总会有一个主 goroutine。这个主 goroutine 会在 Go 程序的运行准备工作完成后被自动地启用，并不需要我们做任何手动的操作。

每条go语句一般都会携带一个函数调用，这个被调用的函数常常被称为go函数。而主 goroutine 的go函数就是那个作为程序入口的`main`函数。

**go函数真正被执行的时间总会与其所属的go语句被执行的时间不同**。

当程序执行到一条go语句的时候，Go 语言的运行时系统，会先试图从某个存放空闲的 G 的队列中获取一个 G（也就是 goroutine），它只有在找不到空闲 G 的情况下才会去创建一个新的 G。
已存在的 goroutine 总是会被优先复用。

在拿到了一个空闲的 G 之后，Go 语言运行时系统会用这个 G 去包装当前的那个go函数（或者说该函数中的那些代码），然后再把这个 G 追加到某个存放可运行的 G 的队列中。
这类队列中的 G 总是会按照先入先出的顺序，很快地由运行时系统内部的调度器安排运行。虽然这会很快，但是由于上面所说的那些准备工作还是不可避免的，所以耗时还是存在的。

因此，go函数的执行时间总是会明显滞后于它所属的go语句的执行时间。当然了，这里所说的“明显滞后”是对于计算机的 CPU 时钟和 Go 程序来说的。我们在大多数时候都不会有明显的感觉。

请记住，**只要go语句本身执行完毕，Go 程序完全不会等待go函数的执行，它会立刻去执行后边的语句。这就是所谓的异步并发地执行**。

上面的代码中那 10 个包装了go函数的 goroutine 往往还没有获得运行的机会。但是如果有机会运行，打印的结果是什么，全是10？

当`for`语句的最后一个迭代运行的时候，其中的那条go语句即是最后一条语句。所以，在执行完这条go语句之后，主 goroutine 中的代码也就执行完了，Go 程序会立即结束运行。
那么，如果这样的话，还会有任何内容被打印出来吗？

Go 语言并不会去保证这些 goroutine 会以怎样的顺序运行。由于主 goroutine 会与我们手动启用的其他 goroutine 一起接受调度，又因为调度器很可能会在 goroutine 中的代码只执行
了一部分的时候暂停，以期所有的 goroutine 有更公平的运行机会。

所以哪个 goroutine 先执行完、哪个 goroutine 后执行完往往是不可预知的，除非我们使用了某种 Go 语言提供的方式进行了人为干预。

### 怎样才能让主 goroutine 等待其他 goroutine
刚才说过，一旦主 goroutine 中的代码执行完毕，当前的 Go 程序就会结束运行，无论其他的 goroutine 是否已经在运行了。那么，怎样才能做到等其他的 goroutine 运行完毕之后，
再让主 goroutine 结束运行呢？

**使用`time`包**

可以简单粗暴的`time.Sleep(time.Millisecond * 500)`让主 goroutine“小睡”一会儿。在这里传入了“500 毫秒”

问题是我们让主 goroutine“睡眠”多长时间才是合适的呢？如果“睡眠”太短，则很可能不足以让其他的 goroutine 运行完毕，而若“睡眠”太长则纯属浪费时间，这个时间就太难把握了。

**使用通道**。

**使用`sync`包的`sync.WaitGroup`类型**

### 怎样让启用的多个 goroutine 按照既定的顺序运行
首先，我们需要稍微改造一下for语句中的那个go函数:
```go
for i := 0; i < 10; i++ {
    go func(i int) {
        fmt.Println(i)
    }(i)
}
```
只有这样，Go 语言才能保证每个 goroutine 都可以拿到一个唯一的整数。这里有点像js。

在go语句被执行时，我们**传给go函数的参数`i`会先被求值**，如此就得到了当次迭代的序号。之后，无论go函数会在什么时候执行，这个参数值都不会变。也就是说，
go函数中调用的`fmt.Println`函数打印的一定会是那个当次迭代的序号。

```go
	var count uint32 = 0
	trigger := func(i uint32, fn func()) { // func()代表的是既无参数声明也无结果声明的函数类型
		for {
			if n := atomic.LoadUint32(&count); n == i {
				fn()
				atomic.AddUint32(&count, 1)
				break
			}
			time.Sleep(time.Nanosecond)
		}
	}
	for i := uint32(0); i < 10; i++ {
		go func(i uint32) {
			fn := func() {
				fmt.Println(i)
			}
			trigger(i, fn)
		}(i)
	}
	trigger(10, func(){})
```
调用了一个名叫`trigger`的函数，并把go函数的参数`i`和刚刚声明的变量`fn`作为参数传给了它。**func()代表的是既无参数声明也无结果声明的函数类型**。

`trigger`函数会不断地获取一个名叫`count`的变量的值，并判断该值是否与参数i的值相同。如果相同，那么就立即调用fn代表的函数，然后把`count`变量的值加`1`，最后显式地退出当前的循环。
否则，我们就先让当前的 goroutine“睡眠”一个纳秒再进入下一个迭代。

操作变量`count`的时候使用的都是原子操作。这是由于`trigger`函数会被多个 goroutine 并发地调用，所以它用到的非本地变量`count`，就被多个用户级线程共用了。因此，对它的操作就产生
了竞态条件（race condition），破坏了程序的并发安全性。在`sync/atomic`包中声明了很多用于原子操作的函数。由于我选用的原子操作函数对被操作的数值的类型有约束，所以我才对`count`以及相
关的变量和参数的类型进行了统一的变更（由`int`变为了`uint32`）。

纵观`count`变量、`trigger`函数以及改造后的`for`语句和go函数，我要做的是，让`count`变量成为一个信号，它的值总是下一个可以调用打印函数的go函数的序号。

这个序号其实就是启用 goroutine 时，那个当次迭代的序号。

依然想让主 goroutine 最后一个运行完毕，所以还需要加一行代码。不过既然有了trigger函数，我就没有再使用通道。
```go
trigger(10, func(){})
```

调用`trigger`函数完全可以达到相同的效果。由于当所有我手动启用的 goroutine 都运行完毕之后，`count`的值一定会是`10`，所以我就把`10`作为了第一个参数值。
又由于我并不想打印这个`10`，所以我把一个什么都不做的函数作为了第二个参数值。
## Channels
**通道类型的值本身就是并发安全的，这也是 Go 语言自带的、唯一一个可以满足并发安全性的类型**。

`channels`是`goroutine`之间的通信机制。`goroutine`通过`channel`向另一个`goroutine`发送消息。`channel`和`goroutine`结合，
可以实现用通信代替共享内存的`CSP`模型。

创建`channel`：
```go
ch := make(chan int)

ch = make(chan int, 3) // buffered channel with capacity 3
```

上面的代码中，`int`代表这个`channel`要发送的数据的类型。第二个参数代表创建一带缓存的`channel`，容量为`3`。
`channel`的零值是`nil`。

发送和接收两个操作使用`<-`运算符，一个左尖括号紧接着一个减号形象地代表了元素值的传输方向：
```go
// 发送一个值
ch <- x // 我的理解就是这里将数据push到channel

// 接受一个值
x = <-ch // 取出channel的值并复制给变量x

<-ch // 接受的值会被丢弃
```

### close

使用`close`函数关闭`channel`，`channel`关闭后不能再发送数据，但是可以接受已经发送成功的数据，如果`channel`中没有
数据，那么返回一个零值。

**注意，`close`函数不是一个清理操作，而是一个控制操作，在确定这个channel不会在发送数据时调用。**

**因为关闭操作只用于断言不再向`channel`发送新的数据，所以只有在发送者所在的`goroutine`才会调用`close`函数**，因此对一个只接收的`channel`调用`close`将是一个编译错误。

使用`range`循环可直接在`channels`上面迭代。它依次从`channel`接收数据，当`channel`被关闭并且没有值可接收时跳出循环。
```go
naturals := make(chan int)
for x := 0; x < 100; x++ {
	naturals <- x
}
for x := range naturals {

}
```
### 无缓存channel
无缓存`channel`也叫做同步`channel`，这是因为如果一个`goroutine`基于一个无缓存`channel`发送数据，那么就会阻塞，直到
另一个`goroutine`在相同的`channel`上执行接收操作。同样的，如果一个`goroutine`基于一个无缓存`channel`先执行了接受操作，
也会阻塞，直到另一个`goroutine`在相同的`channel`上执行发送操作。在`channel`成功传输之后，两个`goroutine`之后的语句才会
继续执行。

### 单向channel

当一个`channel`作为一个函数参数时，它一般总是被专门用于只发送或者只接收。

类型`chan<- int`表示一个只发送`int`的`channel`。相反，类型`<-chan int`表示一个只接收`int`的`channel`。

```go
var uselessChan = make(chan<- int, 1)
```

### 带缓存channel
```go
ch = make(chan int, 3)
```
带缓存的`channel`内部持有一个元素队列。`make`函数创建`channel`时通过第二个参数指定队列的最大容量。

发送操作会向`channel`的缓存队列`push`元素，接收操作则是`pop`元素，如果队列被塞满了，那么发送操作将阻塞直到另一个`goroutine`执行接收操作而释放了新的队列空间。
相反，如果`channel`是空的，接收操作将阻塞直到有另一个`goroutine`执行发送操作而向队列插入元素。

在大多数情况下，缓冲通道会作为收发双方的中间件。正如前文所述，元素值会先从发送方复制到缓冲通道，之后再由缓冲通道复制给接收方。

但是，当发送操作在执行的时候发现空的通道中，正好有等待的接收操作，那么它会直接把元素值复制给接收方。

### cap 和 len
`cap`函数可以获取`channel`内部缓存的容量。
`len`函数可以获取`channel`内部缓存有效元素的个数。

```go
ch = make(chan int, 3)
fmt.Println(cap(ch)) // 3

ch <- "A"
ch <- "B"

fmt.Println(len(ch)) // 2
fmt.Println(<-ch) // A
fmt.Println(len(ch)) // 1
```

### 通道的发送和接收操作的特性
1. 对于同一个通道，发送操作之间是互斥的，接收操作之间也是互斥的。，Go 语言的运行时系统（以下简称运行时系统）只会执行对同一个通道的任意个发送操作中的某一个。直到这个元素值
被完全复制进该通道之后，其他针对该通道的发送操作才可能被执行。
2. 发送操作和接收操作中对元素值的处理都是不可分割的。发送操作要么还没复制元素值，要么已经复制完毕，绝不会出现只复制了一部分的情况。接收操作在准备好元素值的副本
之后，一定会删除掉通道中的原值，绝不会出现通道中仍有残留的情况。
3. 发送操作在完全完成之前会被阻塞。接收操作也是如此。

**元素值从外界进入通道时会被复制。更具体地说，进入通道的并不是在接收操作符右边的那个元素值，而是它的副本**。

**对于通道中的同一个元素值来说，发送操作和接收操作之间也是互斥的。例如，虽然会出现，正在被复制进通道但还未复制完成的元素值，但是这时它绝不会被想接收它的一方看到和取走**。

### 发送操作和接收操作在什么时候可能被长时间的阻塞
- 针对**缓冲通道**的情况。如果通道已满，那么对它的所有发送操作都会被阻塞，直到通道中有元素值被接收走。相对的，如果通道已空，那么对它的所有接收操作都会被阻塞，直到通道中有新的元素值出现。这时，通道会通知最早等待的那个接收操作所在的 goroutine，并使它再次执行接收操作。
- 对于**非缓冲通道**，情况要简单一些。无论是发送操作还是接收操作，一开始执行就会被阻塞，直到配对的操作也开始执行，才会继续传递。
- **对于值为nil的通道，不论它的具体类型是什么，对它的发送操作和接收操作都会永久地处于阻塞状态**。它们所属的 goroutine 中的任何代码，都不再会被执行。注意，由于通道类型是引用类型，所以它的零值就是nil。**当我们只声明该类型的变量但没有用make函数对它进行初始化时，该变量的值就会是nil。我们一定不要忘记初始化通道**！

## select 多路复用
`select`语句是专为通道而设计的，**所以每个`case`表达式中都只能包含操作通道的表达式**，比如接收表达式。

```go
select {
  case communication clause  :
      ...     
  case communication clause  :
      ... 
  default : /* 可选 */
			... 
}			
```

如果有多个`channel`需要接受消息，如果第一个`channel`没有消息发过来，那么程序会被阻塞，第二个`channel`的消息就也无法接收了。
这时候就需要使用`select`多路复用。
```go
select {
  case <-ch1:
      ...     
  case x := <-ch2:
			... 
	case ch3 <- y:
	    ...		
  default:
			... 
}	
```
每一个`case`代表一个通信操作，发送或者接收。**如果没有`case`可运行，它将阻塞，直到有`case`可运行**。
如果多个`case`同时满足条件，`select`会随机地选择一个执行。

**为了避免因为发送或者接收导致的阻塞，尤其是当`channel`没有准备好写或者读时。`default`可以设置当其它的操作都不能够马上被处理时程序需要执行哪些逻辑**。

### 超时
我们可以利用`select`来设置超时，避免`goroutine`阻塞的情况：
```go
func main() {
	c := make(chan int)
	o := make(chan bool)
	go func() {
		for {
			select {
				case v := <- c:
					fmt.println(v)
				case <- time.After(5 * time.Second):
					fmt.println("timeout")
					o <- true
					break
			}
		}
	}()
	<- o
}
```

### 使用select语句的时候，需要注意的事情
1. 如果加入了默认分支，那么无论涉及通道操作的表达式是否有阻塞，`select`语句都不会被阻塞。如果那几个表达式都阻塞了，或者说都没有满足求值的条件，那么默认分支就会被选中并执行。
2. 如果没有加入默认分支，那么一旦所有的`case`表达式都没有满足求值条件，那么`select`语句就会被阻塞。直到至少有一个`case`表达式满足条件为止。
3. 还记得吗？我们可能会因为通道关闭了，而直接从通道接收到一个其元素类型的零值。所以，在很多时候，我们需要通过接收表达式的第二个结果值来判断通道是否已经关闭。一旦发现某个通道关闭了，我们就应该及时地屏蔽掉对应的分支或者采取其他措施。这对于程序逻辑和程序性能都是有好处的。
4. `select`语句只能对其中的每一个`case`表达式各求值一次。所以，如果我们想连续或定时地操作其中的通道的话，就往往需要通过在`for`语句中嵌入`select`语句的方式实现。但这时要注意，简单地在`select`语句的分支中使用`break`语句，只能结束当前的`select`语句的执行，而并不会对外层的`for`语句产生作用。这种错误的用法可能会让这个`for`语句无休止地运行下去。

```go
intChan := make(chan int, 1)
// 一秒后关闭通道。
time.AfterFunc(time.Second, func() {
  close(intChan)
})
select {
  case _, ok := <-intChan:
    if !ok {
      fmt.Println("The candidate case is closed.")
      break
    }
    fmt.Println("The candidate case is selected.")
  }
```
上面的代码`select`语句只有一个候选分支，我在其中利用接收表达式的第二个结果值对`intChan`通道是否已关闭做了判断，并在得到肯定结果后，通过`break`语句立即结束当前`select`语句的执行。

## 共享变量
无论任何时候，只要有两个以上`goroutine`并发访问同一变量，且至少其中的一个是写操作的时候就会发生数据竞争。
避免数据竞争的三种方式：
1. 不去写变量。读取不可能出现数据竞争。
2. 避免从多个`goroutine`访问变量，尽量把变量限定在了一个单独的`goroutine`中。(**不要使用共享数据来通信，使用通信来共享数据**)
3. 互斥锁

### 互斥锁
我们可以使用容量只有`1`的`channel`来保证最多只有一个`goroutine`在同一时刻访问一个共享变量：
```go
var (
  sema = make(chan struct{}, 1) // a binary semaphore guarding balance
  balance int
)

func Deposit(amount int) {
  sema <- struct{}{} // acquire lock
  balance = balance + amount
  <-sema // release lock
}

func Balance() int {
  sema <- struct{}{} // acquire lock
  b := balance
  <-sema // release lock
  return b
}
```
#### sync.Mutex
使用`sync.Mutex`互斥锁：
```go
import "sync"

var (
  mu sync.Mutex // guards balance
  balance int
)

func Deposit(amount int) {
  mu.Lock()
  balance = balance + amount
  mu.Unlock()
}

func Balance() int {
  mu.Lock()
  b := balance
  mu.Unlock()
  return b
}
```
`mutex`会保护共享变量，当已经有`goroutine`获得这个锁，再有`goroutine`访问这个加锁的变量就会被阻塞，直到持有这个锁的`goroutine`
`unlock`这个锁。

我们可以使用`defer`来`unlock`锁，保证在函数返回之后或者发生错误返回时一定会执行`unlock`。

#### 读写锁`sync.RWMutex`
如果有多个`goroutine`读取变量，那么是并发安全的，这个时候使用`sync.Mutex`加锁就没有必要。可以使用`sync.RWMutex`读写锁（多读单写锁）。
**读写锁是把对共享资源的“读操作”和“写操作”区别对待了。它可以对这两种操作施加不同程度的保护**。

一个读写锁中实际上包含了两个锁，即：读锁和写锁。sync.RWMutex类型中的`Lock`方法和`Unlock`方法分别用于对写锁进行锁定和解锁，
而它的`RLock`方法和`RUnlock`方法则分别用于对读锁进行锁定和解锁。

对于同一个读写锁来说有如下规则。

- 在写锁已被锁定的情况下再试图锁定写锁，会阻塞当前的 goroutine。
- 在写锁已被锁定的情况下试图锁定读锁，也会阻塞当前的 goroutine。
- 在读锁已被锁定的情况下试图锁定写锁，同样会阻塞当前的 goroutine。
- 在读锁已被锁定的情况下再试图锁定读锁，并不会阻塞当前的 goroutine。

**对于某个受到读写锁保护的共享资源，多个写操作不能同时进行，写操作和读操作也不能同时进行，但多个读操作却可以同时进行**。

```go
var mu sync.RWMutex
var balance int
func Balance() int {
  mu.RLock() // readers lock
  defer mu.RUnlock()
  return balance
}
```
`RLock`只能在共享变量没有任何写入操作时可用。

为什么只读操作也需要加锁？
```go
var x, y int
go func() {
  x = 1 // A1
  fmt.Print("y:", y, " ") // A2
}()
go func() {
  y = 1                   // B1
  fmt.Print("x:", x, " ") // B2
}()
```

上面的代码打印的结果可能是：
```bash
y:0 x:1
x:0 y:1
x:1 y:1
y:1 x:1

# 还可能是
x:0 y:0
y:0 x:0
```

为什么会有`x:0 y:0`这种结果，在一个`goroutine`中，语句的执行顺序可以保证，在声明的例子，可以保证执行`x = 1`后打印`y:`，但是不能保证
打印`y:`时，另一个`goroutine`中`y = 1`是否已经执行。

所以可能的话，将变量限定在`goroutine`内部；如果是多个`goroutine`都需要访问的变量，使用互斥条件来访问。

#### 注意事项
- 不要重复锁定互斥锁；对一个已经被锁定的互斥锁进行锁定，是会立即阻塞当前的 goroutine 的。这个 goroutine 所执行的流程，会一直停滞在调用该互斥锁的`Lock`方法的那行代码上。
直到该互斥锁的Unlock方法被调用，并且这里的锁定操作成功完成，后续的代码（也就是临界区中的代码）才会开始执行。这也正是互斥锁能够保护临界区的原因所在。
- 不要忘记解锁互斥锁，必要时使用defer语句；避免重复锁定。
- 不要对尚未锁定或者已解锁的互斥锁解锁；解锁“读写锁中未被锁定的写锁”，会立即引发 panic，对于其中的读锁也是如此，并且同样是不可恢复的。
- 不要在多个函数之间直接传递互斥锁。

一旦，你把一个互斥锁同时用在了多个地方，就必然会有更多的 goroutine 争用这把锁。这不但会让你的程序变慢，还会大大增加死锁（deadlock）的可能性。

所谓的**死锁**，指的就是当前程序中的主 goroutine，以及我们启用的那些 goroutine 都已经被阻塞。这些 goroutine 可以被统称为用户级的 goroutine。这就相当于整个程序都已经停滞不前了。

Go 语言运行时系统是不允许这种情况出现的，只要它发现所有的用户级 goroutine 都处于等待状态，就会自行抛出一个带有如下
信息的 panic：`fatal error: all goroutines are asleep - deadlock!`

**注意，这种由 Go 语言运行时系统自行抛出的 panic 都属于致命错误，都是无法被恢复的，调用recover函数对它们起不到任何作用。也就是说，一旦产生死锁，程序必然崩溃**。

**最简单、有效的方式就是让每一个互斥锁都只保护一个临界区或一组相关临界区**。

### 条件变量sync.Cond
条件变量是基于互斥锁的，它必须有互斥锁的支撑才能发挥作用。条件变量并不是被用来保护临界区和共享资源的，它是用于协调想要访问共享资源的那些线程的。
**当共享资源的状态发生变化时，它可以被用来通知被互斥锁阻塞的线程**。

条件变量在这里的最大优势就是在效率方面的提升。当共享资源的状态不满足条件的时候，想操作它的线程再也不用循环往复地做检查了，只要等待通知就好了。

#### 条件变量怎样与互斥锁配合使用
条件变量的初始化离不开互斥锁，并且它的方法有的也是基于互斥锁的。

条件变量提供的方法有三个：等待通知（wait）、单发通知（signal）和广播通知（broadcast）。我们在利用条件变量等待通知的时候，需要在它基于的那个互斥锁
保护下进行。而在进行单发通知或广播通知的时候，却是恰恰相反的，也就是说，需要在对应的互斥锁解锁之后再做这两种操作。

```go
var mailbox uint8
var lock sync.RWMutex
sendCond := sync.NewCond(&lock)
recvCond := sync.NewCond(lock.RLocker())
```
`lock`是一个类型为`sync.RWMutex`的变量，是一个读写锁。基于这把锁，我还创建了两个代表条件变量的变量，名字分别叫`sendCond`和`recvCond`。
`sync.Cond`类型并不是开箱即用的。我们只能利用`sync.NewCond`函数创建它的指针值。

`lock`变量的`Lock`方法和`Unlock`方法分别用于对其中写锁的锁定和解锁，它们与`sendCond`变量的含义是对应的。被视为对共享资源的写操作。

初始化`recvCond`这个条件变量，我们需要的是`lock`变量中的读锁，`sync.RWMutex`类型的`RLocker`方法可以实现这一需求。
`lock.RLocker()`，在其内部会分别调用lock变量的RLock方法和RUnlock方法。

`mailbox`是一个信箱，如果在放置的时候发现信箱里还有未被取走的情报，那就不再放置，而先返回。另一方面，如果你在获取的时候发现信箱里没有情报，那也只能先回去了。

```go
lock.Lock()
for mailbox == 1 {
    sendCond.Wait()
}
mailbox = 1
lock.Unlock()
recvCond.Signal()
```
先调用lock变量的Lock方法。注意，这个Lock方法在这里意味的是：持有信箱上的锁，并且有打开信箱的权利，而不是锁上这个锁。

检查`mailbox`变量的值是否等于1，也就是说，要看看信箱里是不是还存有情报。如果还有情报，那么我就回家去等通知。

如果信箱里没有情报，那么我就把新情报放进去，关上信箱、锁上锁，然后离开。用代码表达出来就是`mailbox = 1`和`lock.Unlock()`。
然后发通知，“信箱里已经有新情报了”，我们调用`recvCond`的`Signal`方法就可以实现这一步骤。

另一方面，你现在是另一个 goroutine，想要适时地从信箱中获取情报，然后通知我。
```go
lock.RLock()
for mailbox == 0 {
    recvCond.Wait()
}
mailbox = 0
lock.RUnlock()
sendCond.Signal()
```
事情在流程上其实基本一致，只不过每一步操作的对象是不同的。

**为什么先要锁定条件变量基于的互斥锁，才能调用它的Wait方法？**

`Wait`方法主要做了四件事。

1. 把调用它的 goroutine（也就是当前的 goroutine）加入到当前条件变量的通知队列中。
2. 解锁当前的条件变量基于的那个互斥锁。
3. 让当前的 goroutine 处于等待状态，等到通知到来时再决定是否唤醒它。此时，这个 goroutine 就会阻塞在调用这个Wait方法的那行代码上。
4. 如果通知到来并且决定唤醒这个 goroutine，那么就在唤醒它之后重新锁定当前条件变量基于的互斥锁。自此之后，当前的 goroutine 就会继续执行后面的代码了。

因为条件变量的`Wait`方法在阻塞当前的 goroutine 之前会解锁它基于的互斥锁，所以在调用该`Wait`方法之前我们必须先锁定那个互斥锁，否则在调用这个`Wait`方法时，就会引发一个不可恢复的 panic。

为什么条件变量的`Wait`方法要这么做呢？你可以想象一下，如果Wait方法在互斥锁已经锁定的情况下，阻塞了当前的 goroutine，那么又由谁来解锁呢？别的 goroutine 吗？

先不说这违背了互斥锁的重要使用原则，即：成对的锁定和解锁，就算别的 goroutine 可以来解锁，那万一解锁重复了怎么办？由此引发的 panic 可是无法恢复的。

如果当前的 goroutine 无法解锁，别的 goroutine 也都不来解锁，那么又由谁来进入临界区，并改变共享资源的状态呢？只要共享资源的状态不变，即使当前的 goroutine 因收到通知而被唤醒，也依然会再次执行这个Wait方法，并再次被阻塞。

所以说，如果条件变量的Wait方法不先解锁互斥锁的话，那么就只会造成两种后果：不是当前的程序因 panic 而崩溃，就是相关的 goroutine 全面阻塞。

**为什么要用for语句来包裹调用其Wait方法的表达式，用if语句不行吗？**

`if`语句只会对共享资源的状态检查一次，而`for`语句却可以做多次检查，直到这个状态改变为止。

那为什么要做多次检查呢？

为了保险起见。如果一个 goroutine 因收到通知而被唤醒，但却发现共享资源的状态，依然不符合它的要求，那么就应该再次调用条件变量的Wait方法，并继续等待下次通知的到来。
这种情况是很有可能发生的。

#### 条件变量的Signal方法和Broadcast方法有哪些异同
条件变量的Signal方法和Broadcast方法都是被用来发送通知的，不同的是，前者的通知只会唤醒一个因此而等待的 goroutine，而后者的通知却会唤醒所有为此等待的 goroutine。

条件变量的Wait方法总会把当前的 goroutine 添加到通知队列的队尾，而它的Signal方法总会从通知队列的队首开始查找可被唤醒的 goroutine。所以，因Signal方法的通知而
被唤醒的 goroutine 一般都是最早等待的那一个。

## 原子操作
Go 语言的原子操作当然是基于 CPU 和操作系统的，所以它也只针对少数数据类型的值提供了原子操作函数。这些函数都存在于标准库代码包`sync/atomic`中。

`sync/atomic`包中的函数可以做的原子操作有：加法（add）、比较并交换（compare and swap，简称 CAS）、加载（load）、存储（store）和交换（swap）。

这些函数针对的数据类型并不多。对这些类型中的每一个，sync/atomic包都会有一套函数给予支持。这些数据类型有：int32、int64、uint32、uint64、uintptr，以及unsafe包中的Pointer。
不过，针对unsafe.Pointer类型，该包并未提供进行原子加法操作的函数。

`sync/atomic`包还提供了一个名为Value的类型，它可以被用来存储任意类型的值。

`atomic.AddInt32`函数的第一个参数，为什么不是`int32`而是`*int32`呢？

因为原子操作函数需要的是被操作值的指针，而不是这个值本身；被传入函数的参数值都会被复制，像这种基本类型的值一旦被传入函数，就已经与函数外的那个值毫无关系了。

所以，传入值本身没有任何意义。unsafe.Pointer类型虽然是指针类型，但是那些原子操作函数要操作的是这个指针值，而不是它指向的那个值，所以需要的仍然是指向这个指针值的指针。

只要原子操作函数拿到了被操作值的指针，就可以定位到存储该值的内存地址。只有这样，它们才能够通过底层的指令，准确地操作这个内存地址上的数据。

### 比较并交换操作与交换操作相比有什么不同
比较并交换操作即 CAS 操作，是有条件的交换操作，**只有在条件满足的情况下才会进行值的交换**。

**所谓的交换指的是，把新值赋给变量，并返回变量的旧值**。

CAS 操作用途要更广泛一些。例如，我们将它与for语句联用就可以实现一种简易的自旋锁（spinlock）。
```go
for {
    if atomic.CompareAndSwapInt32(&num2, 10, 0) {
        fmt.Println("The second number has gone to zero.")
        break
    }
    time.Sleep(time.Millisecond * 500)
}
```
在`for`语句中的 CAS 操作可以不停地检查某个需要满足的条件，一旦条件满足就退出`for`循环。这就相当于，只要条件未被满足，当前的流程就会被一直“阻塞”在这里。

这在效果上与互斥锁有些类似。不过，它们的适用场景是不同的。我们在使用互斥锁的时候，总是假设共享资源的状态会被其他的 goroutine 频繁地改变。

而`for`语句加 CAS 操作的假设往往是：共享资源状态的改变并不频繁，或者，它的状态总会变成期望的那样。这是一种更加乐观，或者说更加宽松的做法。

**假设我已经保证了对一个变量的写操作都是原子操作，比如：加或减、存储、交换等等，那我对它进行读操作的时候，还有必要使用原子操作吗**？

很有必要。其中的道理你可以对照一下读写锁。为什么在读写锁保护下的写操作和读操作之间是互斥的？这是为了防止读操作读到没有被修改完的值，对吗？

如果写操作还没有进行完，读操作就来读了，那么就只能读到仅修改了一部分的值。这显然破坏了值的完整性，读出来的值也是完全错误的。

所以，一旦你决定了要对一个共享资源进行保护，那就要做到完全的保护。不完全的保护基本上与不保护没有什么区别。

### `sync/atomic.Value`
此类型的值相当于一个容器，可以被用来“原子地”存储和加载任意的值。开箱即用。

它只有两个指针方法——`Store`和`Load`。不过，虽然简单，但还是有一些值得注意的地方的。

1. 一旦atomic.Value类型的值（以下简称原子值）被真正使用，它就不应该再被复制了。只要用它来存储值了，就相当于开始真正使用了。atomic.Value类型属于结构体类型，
而结构体类型属于值类型。所以，复制该类型的值会产生一个完全分离的新值。这个新值相当于被复制的那个值的一个快照。之后，不论后者存储的值怎样改变，都不会影响到前者。
2. 不能用原子值存储`nil`。
3. 我们向原子值存储的第一个值，决定了它今后能且只能存储哪一个类型的值。
4. 尽量不要向原子值中存储引用类型的值。因为这很容易造成安全漏洞。
```go
var box6 atomic.Value
v6 := []int{1, 2, 3}
box6.Store(v6)
v6[1] = 4 // 注意，此处的操作不是并发安全的！
```
切片类型属于引用类型。所以，我在外面改动这个切片值，就等于修改了`box6`中存储的那个值。这相当于绕过了原子值而进行了非并发安全的操作。怎样修补：
```go
store := func(v []int) {
    replica := make([]int, len(v))
    copy(replica, v)
    box6.Store(replica)
}
store(v6)
v6[2] = 5 // 此处的操作是安全的。
```
先为切片值`v6`创建了一个完全的副本。这个副本涉及的数据已经与原值毫不相干了。然后，我再把这个副本存入`box6`。如此一来，无论我再对`v6`的值做怎样的修改，都不会破坏`box6`提供的安全保护。

## sync.WaitGroup
在一些场合下里，我们使用通道的方式看起来都似乎有些蹩脚。比如：声明一个通道，使它的容量与我们手动启用的 goroutine 的数量相同。之后利用这个通道，
让主 goroutine 等待其他 goroutine 的运行结束。更具体地说就是：让其他的 goroutine 在运行结束之前，都向这个通道发送一个元素值，并且，
让主 goroutine 在最后从这个通道中接收元素值，接收的次数需要与其他的 goroutine 的数量相同。

```go
func coordinateWithChan() {
    sign := make(chan struct{}, 2)
    num := int32(0)
    fmt.Printf("The number: %d [with chan struct{}]\n", num)
    max := int32(10)
    go addNum(&num, 1, max, func() {
        sign <- struct{}{}
    })
    go addNum(&num, 2, max, func() {
        sign <- struct{}{}
    })
    <-sign
    <-sign
}
```
`coordinateWithChan`函数中最后的那两行代码了吗？重复的两个接收表达式`<-sign`，很丑陋。
我们可以选用另外一个同步工具，即：`sync`包的`WaitGroup`类型。它比通道更加适合实现这种一对多的 goroutine 协作流程。

`sync.WaitGroup`类型（以下简称`WaitGroup`类型）是开箱即用的，也是并发安全的。

`WaitGroup`类型拥有三个指针方法：`Add`、`Done`和`Wait`。**你可以想象该类型中有一个计数器，它的默认值是`0`。我们可以通过调用该类型值的Add方法来增加，或者减少这个计数器的值**。

**一般情况下，我会用这个方法来记录需要等待的 goroutine 的数量。相对应的，这个类型的`Done`方法，用于对其所属值中计数器的值进行减一操作**。我们可以在需要等待的 goroutine 中，
通过`defer`语句调用它。

而**此类型的`Wait`方法的功能是，阻塞当前的 goroutine，直到其所属值中的计数器归零**。

改造版本：
```go
func coordinateWithWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(2)
	num := int32(0)
	fmt.Printf("The number: %d [with sync.WaitGroup]\n", num)
	max := int32(10)
	go addNum(&num, 3, max, wg.Done)
	go addNum(&num, 4, max, wg.Done)
	wg.Wait()
}
```

### sync.WaitGroup类型值中计数器的值可以小于0吗
不可以。**小于0，会引发一个 panic**。

WaitGroup值是可以被复用的，但需要保证其计数周期的完整性。这里的计数周期指的是这样一个过程：该值中的计数器值由0变为了某个正整数，而后又经过一系列的变化，
最终由某个正整数又变回了0。

如果在一个此类值的Wait方法被执行期间，跨越了两个计数周期，那么就会引发一个 panic。

### 使用注意
- 不要把增加其计数器值的操作和调用其Wait方法的代码，放在不同的 goroutine 中执行。换句话说，要杜绝对同一个WaitGroup值的两种操作的并发执行。

## sync.Once

与`sync.WaitGroup`类型一样，`sync.Once`类型（以下简称Once类型）也属于结构体类型，同样也是开箱即用和并发安全的。由于这个类型中包含了一个`sync.Mutex`类型的字段，
所以，复制该类型的值也会导致功能的失效。

```go
var loadIconsOnce sync.Once
var icons map[string]image.Image
// Concurrency-safe.
func Icon(name string) image.Image {
  loadIconsOnce.Do(loadIcons)
  return icons[name]
}
```
`Once`类型的`Do`方法只接受一个参数，这个参数的类型必须是`func()`，即：无参数声明和结果声明的函数。该方法的功能并不是对每一种参数函数都只执行一次，
而是只**执行“首次被调用时传入的”那个函数，并且之后不会再执行任何参数函数**。

所以，**如果你有多个只需要执行一次的函数，那么就应该为它们中的每一个都分配一个`sync.Once`类型的值**（以下简称`Once`值）。

`Once`类型中还有一个名叫`done`的`uint32`类型的字段。**它的作用是记录其所属值的`Do`方法被调用的次数。不过，该字段的值只可能是0或者1。一旦Do方法的首次调用完成，
它的值就会从0变为1**。

**既然done字段的值不是0就是1，那为什么还要使用需要四个字节的uint32类型呢**？

原因很简单，因为对它的操作必须是“原子”的。`Do`方法在一开始就会通过调用`atomic.LoadUint32`函数来获取该字段的值，并且一旦发现该值为1就会直接返回。
这也初步保证了“Do方法，只会执行首次被调用时传入的函数”。

### Do方法在功能方面的两个特点
- 由于`Do`方法只会在参数函数执行结束之后把`done`字段的值变为1，因此，如果参数函数的执行需要很长时间或者根本就不会结束（比如执行一些守护任务），
那么就有可能会导致相关 goroutine 的同时阻塞
- `Do`方法在参数函数执行结束后，对`done`字段的赋值用的是原子操作，并且，这一操作是被挂在`defer`语句中的。因此，不论参数函数的执行会以怎样的方式结束，`done`字段的值都会变为1。

## context.Context类型
使用`WaitGroup`值的时候，我们最好用**先统一`Add`，再并发`Done`，最后`Wait`**的标准模式来构建协作流程。如果在调用该值的`Wait`方法的同时，为了增大其计数器的值，
而并发地调用该值的`Add`方法，那么就很可能会引发 panic。

但是**如果，我们不能在一开始就确定执行子任务的 goroutine 的数量，那么使用`WaitGroup`值来协调它们和分发子任务的 goroutine，就是有一定风险的**。一个解决方案是：
**分批地启用执行子任务的 goroutine**。

`WaitGroup`值是可以被复用的，但需要保证其计数周期的完整性。尤其是涉及对其`Wait`方法调用的时候，它的下一个计数周期必须要等到，与当前计数周期对应的那个`Wait`方法调用完成之后，
才能够开始。

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

### 使用`context`包中的程序实体，实现一对多的 goroutine 协作流程
用`context`包中的函数和`Context`类型作为实现工具，实现`coordinateWithContext`的函数。这个函数应该具有上面`coordinateWithWaitGroup`函数相同的功能。
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
先后调用了`context.Background`函数和`context.WithCancel`函数，并得到了一个可撤销的`context.Context`类型的值（由变量cxt代表），以及一个`context.CancelFunc`
类型的撤销函数（由变量`cancelFunc`代表）。

注意我给予`addNum`函数的最后一个参数值。它是一个匿名函数，其中只包含了一条`if`语句。这条`if`语句会**原子地**加载`num`变量的值，并判断它是否等于`total`变量的值。

如果两个值相等，那么就调用`cancelFunc`函数。其含义是，如果所有的`addNum`函数都执行完毕，那么就立即通知分发子任务的 goroutine。

**这里分发子任务的 goroutine，即为执行`coordinateWithContext`函数的 goroutine**。它在执行完`for`语句后，会立即调用`cxt`变量的`Done`函数，并试图针对该函数返回的通道，
进行接收操作。

一旦`cancelFunc`函数被调用，针对该通道的接收操作就会马上结束，所以，这样做就可以实现“等待所有的`addNum`函数都执行完毕”的功能。

### context.Context类型
Context类型的值（以下简称Context值）是可以繁衍的，这意味着我们可以通过一个Context值产生出任意个子值。这些子值可以携带其父值的属性和数据，也可以响应我们通过其父值传达的信号。

正因为如此，所有的Context值共同构成了一颗代表了上下文全貌的树形结构。这棵树的树根（或者称上下文根节点）是一个已经在context包中预定义好的Context值，它是**全局唯一**的。
通过调用`context.Background`函数，我们就可以获取到它（在`coordinateWithContext`函数中就是这么做的）。

注意一下，这个**上下文根节点仅仅是一个最基本的支点，它不提供任何额外的功能**。也就是说，它既不可以被撤销（cancel），也不能携带任何数据。

context包中还包含了**四个用于繁衍Context值的函数，即：`WithCancel`、`WithDeadline`、`WithTimeout`和`WithValue`**。

这些函数的第一个参数的类型都是`context.Context`，而名称都为parent。顾名思义，**这个位置上的参数对应的都是它们将会产生的Context值的父值**。

**`WithCancel`函数用于产生一个可撤销的parent的子值**。

在`coordinateWithContext`函数中，通过调用该函数，获得了一个衍生自上下文根节点的Context值，和一个用于触发撤销信号的函数。

`WithDeadline`函数和`WithTimeout`函数则都可以被用来产生一个会**定时撤销**的parent的子值。至于`WithValue`函数，我们可以通过调用它，产生一个会携带额外数据的parent的子值。

### “可撤销的”在context包中代表着什么？“撤销”一个Context值又意味着什么？

这需要从Context类型的声明讲起。这个接口中有两个方法与“撤销”息息相关。`Done`方法会返回一个元素类型为`struct{}`的接收通道。不过，这个接收通道的用途并不是传递元素值，
而是**让调用方去感知“撤销”当前Context值的那个信号**。

一旦当前的Context值被撤销，这里的接收通道就会被立即关闭。我们都知道，对于一个未包含任何元素值的通道来说，它的关闭会使任何针对它的接收操作立即结束。

正因为如此，在`coordinateWithContext`函数中，基于调用表达式`cxt.Done()`的接收操作，才能够起到感知撤销信号的作用。

### 撤销信号是如何在上下文树中传播的

context包的`WithCancel`函数在被调用后会产生两个结果值。第一个结果值就是那个可撤销的Context值，而第二个结果值则是用于触发撤销信号的函数。

在撤销函数被调用之后，对应的Context值会先关闭它内部的接收通道，也就是它的`Done`方法会返回的那个通道。

然后，它会向它的所有子值（或者说子节点）传达撤销信号。这些子值会如法炮制，把撤销信号继续传播下去。最后，这个Context值会断开它与其父值之间的关联。

**通过调用`context.WithValue`函数得到的Context值是不可撤销的**。

### 怎样通过Context值携带数据

**`WithValue`函数在产生新的Context值（以下简称含数据的Context值）的时候需要三个参数，即：父值、键和值**。与“字典对于键的约束”类似，这里键的类型**必须是可判等**的。

原因很简单，当我们从中获取数据的时候，它需要根据给定的键来查找对应的值。不过，这种Context值并不是用字典来存储键和值的，后两者只是被简单地存储在前者的相应字段中而已。

## 临时对象池sync.Pool
 Go 语言标准库中最重要的那几个同步工具，这包括:
 - 互斥锁
 - 读写锁
 - 条件变量
 - 原子操作
 - `sync/atomic.Value`
 - `sync.Once`
 - `sync.WaitGroup`
 - `context.Context`

Go 语言标准库中的还有另一个同步工具：`sync.Pool`。

`sync.Pool`类型可以被称为临时对象池，它的值可以被用来存储临时的对象。与 Go 语言的很多同步工具一样，`sync.Pool`类型也属于结构体类型，它的值在被真正使用之后，就不应该再被复制了。

**临时对象**的意思是：不需要持久使用的某一类值。这类值对于程序来说可有可无，但如果有的话会明显更好。它们的创建和销毁可以在任何时候发生，并且完全不会影响到程序的功能。

**我们可以把临时对象池当作针对某种数据的缓存来用**。

`sync.Pool`类型只有两个方法——`Put`和`Get`。前者用于在当前的池中存放临时对象，它接受一个`interface{}`类型的参数；而后者则被用于从当前的池中获取临时对象，
它会返回一个`interface{}`类型的值。

更具体地说，**这个类型的`Get`方法可能会从当前的池中删除掉任何一个值，然后把这个值作为结果返回。如果此时当前的池中没有任何值，那么这个方法就会使用当前池的`New`字段创建一个新值，
并直接将其返回**。

`sync.Pool`类型的`New`字段代表着创建临时对象的函数。它的类型是没有参数但有唯一结果的函数类型，即：`func() interface{}`。**初始化这个池的时候最好给定它**。

这个函数是`Get`方法最后的临时对象获取手段。`Get`方法如果到了最后，仍然无法获取到一个值，那么就会调用该函数。该函数的结果值并不会被存入当前的临时对象池中，
而是直接返回给`Get`方法的调用方。

**临时对象池中存储的每一个值都应该是独立的、平等的和可重用的**。`sync.Pool`的定位不是做类似连接池的东西，它的用途仅仅是增加对象重用的几率，减少gc的负担。
因为gc带来了编程的方便但同时也增加了运行时开销，使用不当甚至会严重影响程序的性能。因此性能要求高的场景不能任意产生太多的垃圾。如何解决呢？那就是要重用对象了。

一个比较好的例子是`fmt`包，`fmt`包总是需要使用一些`[]byte`之类的对象，golang建立了一个临时对象池，存放着这些对象，如果需要使用一个`[]byte`，就去`Pool`里面拿，
如果拿不到就分配一份。这比起不停生成新的`[]byte`，用完了再等待gc回收来要高效得多。

`sync.Pool`缓存对象的期限是很诡异的，先看一下src/pkg/sync/pool.go里面的一段实现代码：
```go
func init() {
    runtime_registerPoolCleanup(poolCleanup)
}
```

可以看到`pool`包在`init`的时候注册了一个`poolCleanup`函数，它会清除所有的`pool`里面的所有缓存的对象，该函数注册进去之后会在每次gc之前都会调用，
因此**`sync.Pool`缓存的期限只是两次gc之间这段时间**。

## sync.Map
Go 语言自带的字典类型`map`并不是并发安全的。换句话说，在同一时间段内，让不同 goroutine 中的代码，对同一个字典进行读写操作是不安全的。

Go 语言官方终于在 2017 年发布的 Go 1.9 中正式加入了并发安全的字典类型`sync.Map`。

使用`sync.Map`可以显著地减少锁的争用。`sync.Map`本身虽然也用到了锁，但是，它其实在尽可能地避免使用锁。

**使用锁就意味着要把一些并发的操作强制串行化。这往往会降低程序的性能，尤其是在计算机拥有多个 CPU 核心的情况下**。

由于**并发安全字典内部使用的存储介质正是原生字典，又因为它使用的原生字典键类型也是可以包罗万象的`interface{}`，所以，我们绝对不能带着任何实际类型为函数类型、
字典类型或切片类型的键值去操作并发安全字典**。

因为**这些键值的实际类型只有在程序运行期间才能够确定，所以 Go 语言编译器是无法在编译期对它们进行检查的，不正确的键值实际类型肯定会引发 panic**。

**因此，我们在这里首先要做的一件事就是：一定不要违反上述规则。我们应该在每次操作并发安全字典的时候，都去显式地检查键值的实际类型。无论是存、取还是删，都应该如此**。

> **更好的做法是，把针对同一个并发安全字典的这几种操作都集中起来，然后统一地编写检查代码。除此之外，把并发安全字典封装在一个结构体类型中，往往是一个很好的选择**。如果你实在拿不准，那么可以
先通过调用`reflect.TypeOf`函数得到一个键值对应的反射类型值（即：`reflect.Type`类型的值），然后再调用这个值的`Comparable`方法，得到确切的判断结果。

### 并发安全字典如何做到尽量避免使用锁
`sync.Map`类型在内部使用了**大量的原子操作来存取键和值，并使用了两个原生的map作为存储介质**。

其中一个原生map被存在了`sync.Map`的`read`字段中，该字段是`sync/atomic.Value`类型的。简称它为**只读字典**。

**只读字典虽然不会增减其中的键，但却允许变更其中的键所对应的值。**所以，它并不是传统意义上的快照，它的只读特性只是对于其中键的集合而言的。

由read字段的类型可知，`sync.Map`在替换只读字典的时候根本用不着锁。另外，这个只读字典在存储键值对的时候，还在值之上封装了一层。

它先把值转换为了`unsafe.Pointer`类型的值，然后再把后者封装，并储存在其中的原生字典中。如此一来，在变更某个键所对应的值的时候，就也可以使用原子操作了。

`sync.Map`中的另一个原生字典由它的`dirty`字段代表。它存储键值对的方式与`read`字段中的原生字典一致，它的键类型也是`interface{}`，并且同样是把值先做转换和封装后
再进行储存的。称为**脏字典**。

> 脏字典和只读字典如果都存有同一个键值对，那么这里的两个键指的肯定是同一个基本值，对于两个值来说也是如此。正如前文所述，这两个字典在存储键和值的时候都只会存入它们的某个指针，
而不是基本值。

sync.Map在查找指定的键所对应的值的时候，总会先去只读字典中寻找，并不需要锁定互斥锁。只有当确定“只读字典中没有，但脏字典中可能会有这个键”的时候，
它才会在锁的保护下去访问脏字典。

相对应的，sync.Map在存储键值对的时候，只要只读字典中已存有这个键，并且该键值对未被标记为“已删除”，就会把新值存到里面并直接返回，这种情况下也不需要用到锁。

否则，它才会在锁的保护下把键值对存储到脏字典中。这个时候，该键值对的“已删除”标记会被抹去。

只有当一个键值对应该被删除，但却仍然存在于只读字典中的时候，才会被用标记为“已删除”的方式进行逻辑删除，而不会直接被物理删除。这种情况会在重建脏字典以后
的一段时间内出现。不过，过不了多久，它们就会被真正删除掉。在查找和遍历键值对的时候，已被逻辑删除的键值对永远会被无视。

最后，sync.Map会把该键值对中指向值的那个指针置为nil，这是另一种逻辑删除的方式。

除此之外，还有一个细节需要注意，只读字典和脏字典之间是会互相转换的。在脏字典中查找键值对次数足够多的时候，`sync.Map`会把脏字典直接作为只读字典，保
存在它的`read`字段中，然后把代表脏字典的`dirty`字段的值置为`nil`。

在这之后，一旦再有新的键值对存入，它就会依据只读字典去重建脏字典。这个时候，它会把只读字典中已被逻辑删除的键值对过滤掉。理所当然，这些转换操作肯定都需要在锁的
保护下进行。

**`sync.Map` 的只读字典和脏字典中的键值对集合并不是实时同步的，它们在某些时间段内可能会有不同**。

可以看出，在读操作有很多但写操作却很少的情况下，并发安全字典的性能往往会更好。在几个写操作当中，新增键值对的操作对并发安全字典的性能影响是最大的，
其次是删除操作，最后才是修改操作。

如果被操作的键值对已经存在于`sync.Map`的只读字典中，并且没有被逻辑删除，那么修改它并不会使用到锁，对其性能的影响就会很小。

## 竞争检查器
在 `go build`，`go run` 或者 `go test` 命令后面加上 `-race`，就会使编译器创建一个你的应用的“修改”版。

会记录下每一个读或者写共享变量的 `goroutine` 的身份信息。记录下所有的同步事件，比如 `go` 语句，`channel` 操作，
以及对 `(*sync.Mutex).Lock`，`(*sync.WaitGroup).Wait` 等等的调用。

由于需要额外的记录，因此构建时加了竞争检测的程序跑起来会慢一些，且需要更大的内存，即使是这样，这些代价对于很多生产环境的工作来说还是可以接受的。

## Goroutine 调度器
### 先了解并发和并行
#### 并发
一个cpu上能同时执行多项任务，在很短时间内，cpu来回切换任务执行(在某段很短时间内执行程序a，然后又迅速得切换到程序b去执行)，
有时间上的重叠（宏观上是同时的，微观仍是顺序执行）,这样看起来多个任务像是同时执行，这就是并发。

#### 并行
当系统有多个CPU时,每个CPU同一时刻都运行任务，互不抢占自己所在的CPU资源，同时进行，称为并行。

#### 进程

cpu在切换程序的时候，如果不保存上一个程序的状态（也就是我们常说的context--上下文），直接切换下一个程序，就会丢失上一个程序的一系列状态，于是引入了进程这个概念，
用以划分好程序运行时所需要的资源。因此进程就是一个程序运行时候的所需要的基本资源单位（也可以说是程序运行的一个实体）。

#### 线程
cpu切换多个进程的时候，会花费不少的时间，因为切换进程需要切换到内核态，而每次调度需要内核态都需要读取用户态的数据，进程一旦多起来，cpu调度会消耗一大堆资源，因此引入了线
程的概念，线程本身几乎不占有资源，他们共享进程里的资源，内核调度起来不会那么像进程切换那么耗费资源。

#### 协程
协程拥有自己的寄存器上下文和栈。协程调度切换时，将寄存器上下文和栈保存到其他地方，在切回来的时候，恢复先前保存的寄存器上下文和栈。因此，协程能保留上一次调用时
的状态（即所有局部状态的一个特定组合），每次过程重入时，就相当于进入上一次调用的状态，换种说法：进入上一次离开时所处逻辑流的位置。线程和进程的操作是由程序触发系统接口，
最后的执行者是系统；协程的操作执行者则是用户自身程序，goroutine也是协程。

#### 调度器
Go的runtime负责对goroutine进行“调度”。调度本质上就是决定何时哪个goroutine将获得资源开始执行、哪个goroutine应该停止执行让出资源、哪个goroutine应该被唤醒恢复执行等。

操作系统对进程、线程的调度是指操作系统调度器将系统中的多个线程按照一定算法调度到物理CPU上去运行。C、C++等的并发实现就是基于操作系统调度的，即程序负责创建线程，操作系统负责调度。
但是这种支持并发的方式有不少缺陷：
- 对于很多网络服务程序，由于不能大量创建thread，就要在少量thread里做网络多路复用，即：使用epoll/kqueue/IoCompletionPort这套机制，即便有libevent/libev这样的第三方库帮忙，
写起这样的程序也是很不易的
- 一个thread的代价已经比进程小了很多了，但我们依然不能大量创建thread，因为除了每个thread占用的资源不小之外，操作系统调度切换thread的代价也不小；
- 并发单元间通信困难，易错：多个thread之间的通信虽然有多种机制可选，但用起来是相当复杂；

Go采用了**用户层轻量级thread**或者说是**类coroutine**的概念来解决这些问题，Go将之称为**goroutine**。

**goroutine占用的资源非常小(goroutine stack的size默认为2k)，goroutine调度的切换也不用操作系统内核层完成，代价很低**。所有的Go代码都在goroutine中执行，go runtime也一样。
将这些goroutines按照一定算法放到“CPU”上执行的程序就叫做**goroutine调度器**或**goroutine scheduler**。

**一个Go程序对于操作系统来说只是一个用户层程序，对于操作系统而言，它的眼中只有thread，它并不知道什么是Goroutine。goroutine的调度全要靠Go自己完成，
实现Go程序内goroutine之间“公平”的竞争“CPU”资源，这个任务就落到了Go runtime头上**，在一个Go程序中，除了用户代码，剩下的就是go runtime了。

Goroutine的调度问题就变成了**go runtime如何将程序内的众多goroutine按照一定算法调度到“CPU”资源上运行**了。

但是在**操作系统层面，Thread竞争的“CPU”资源是真实的物理CPU**，但在Go程序层面，各个Goroutine要竞争的”CPU”资源是什么呢？Go程序是用户层程序，它本身整体是运行在一个或多个操
作系统线程上的，因此**goroutine们要竞争的所谓“CPU”资源就是操作系统线程**。

Go scheduler的任务：**将goroutines按照一定算法放到不同的操作系统线程中去执行**。这种在语言层面自带调度器的，我们称之为**原生支持并发**。

### G-P-M模型
前面已经知道Go 并发编程模型中的三个主要元素，即：G（goroutine 的缩写）、P（processor 的缩写）和 M（machine 的缩写） M 指代的就是系统级线程。
- G表示goroutine，存储了goroutine的执行stack信息、goroutine状态以及goroutine的任务函数等；G对象是可以重用的。
- P表示逻辑processor，P的数量决定了系统内最大可并行的G的数量（前提：系统的物理cpu核数>=P的数量）；主要用途就是用来执行goroutine的，
它维护了一个goroutine队列，里面存储了所有需要它来执行的goroutine
- M代表着真正的执行计算资源。goroutine就是跑在M之上的。在绑定有效的p后，进入schedule循环；而schedule循环的机制大致是从各种队列、p的本地队列中获取G，切换到G的执行
栈上并执行G的函数，调用goexit做清理工作并回到m，如此反复。M并不保留G状态，这是G可以跨M调度的基础。

![](../imgs/goroutine-scheduler-model.png)

### 抢占式调度
Go并没有时间片的概念。如果某个G没有进行system call调用、没有进行I/O操作、没有阻塞在一个channel操作上，那么M是**如何让G停下来并调度下一个runnable G**的呢？
答案是：G是被抢占调度的。

Go在设计之初并没考虑将goroutine设计成抢占式的。用户负责让各个goroutine交互合作完成任务。一个goroutine只有在涉及到加锁，读写通道或者主动让出CPU等操作时才会触发切换。

垃圾回收器是需要stop the world的。如果垃圾回收器想要运行了，那么它必须先通知其它的goroutine合作停下来，这会造成较长时间的等待时间。考虑一种很极端的情况，所有
的goroutine都停下来了，只有其中一个没有停，那么垃圾回收就会一直等待着没有停的那一个。

抢占式调度可以解决这种问题，在抢占式情况下，如果一个goroutine运行时间过长，它就会被剥夺运行权。Go还只是引入了一些很初级的抢占，只有长时间阻塞于系统调用，或者运行了
较长时间才会被抢占。runtime会在后台有一个检测线程，它会检测这些情况，并通知goroutine执行调度。

Go程序的初始化过程中，runtime开了一条后台线程，运行一个sysmon函数(一般称为监控线程)。这个函数会周期性地做epoll操作，同时它还会检测每个P是否运行了较长时间。
该 M 无需绑定 P 即可运行，该 M 在整个Go程序的运行过程中至关重要。

sysmon每20us~10ms运行一次，sysmon主要完成如下工作：
- 释放闲置超过5分钟的span物理内存；
- 如果超过2分钟没有垃圾回收，强制执行；
- 将长时间未处理的netpoll结果添加到任务队列；
- 向长时间运行的G任务发出抢占调度；
- 收回因syscall长时间阻塞的P；

### channel阻塞或network I/O情况下的调度
如果G被阻塞在某个channel操作或network I/O操作上时，G会被放置到某个wait队列中，而M会尝试运行下一个runnable的G；如果此时没有runnable的G供M运行，那么M将解绑P，
并进入sleep状态。当I/O available或channel操作完成，在wait队列中的G会被唤醒，标记为runnable，放入到某P的队列中，绑定一个M继续执行。

### system call阻塞情况下的调度
如果G被阻塞在某个system call操作上，那么不光G会阻塞，执行该G的M也会解绑P(实质是被sysmon抢走了)，与G一起进入sleep状态。如果此时有idle的M，则P与其绑定继续执行其他G；
如果没有idle M，但仍然有其他G要去执行，那么就会创建一个新M。

当阻塞在syscall上的G完成syscall调用后，G会去尝试获取一个可用的P，如果没有可用的P，那么G会被标记为runnable，之前的那个sleep的M将再次进入sleep。

# 内存管理
## tcmalloc
Golang 的内存管理基于 tcmalloc，什么是 tcmalloc。

tcmalloc是google推出的一种内存分配器，常见的内存分配器还有glibc的ptmalloc和google的jemalloc。相比于ptmalloc，tcmalloc性能更好，特别适用于高并发场景。

### tcmalloc策略
tcmalloc分配的内存主要来自两个地方：全局缓存堆和进程的私有缓存。对于一些小容量的内存申请使用进程的私有缓存，私有缓存不足的时候可以再从全局缓存申请一部分作为私有缓存。
对于大容量的内存申请则需要从全局缓存中进行申请。而**大小容量的边界就是32k**。缓存的组织方式是一个单链表数组，数组的每个元素是一个单链表，链表中的每个元素具有相同的大小。

### 逃逸分析（escape analysis）
对于手动管理内存的语言，比如 C/C++，我们使用`malloc`或者`new`申请的变量会被分配到堆上。但是 Golang 并不是这样，虽然 Golang 语言里面也有`new`。
Golang 编译器决定变量应该分配到什么地方时会进行逃逸分析：
```go
package main

func foo() *int {
	var x int
	return &x
}

func bar() int {
	x := new(int)
	*x = 1
	return *x
}

func main() {}
```
将上面文件保存为`escape.go`，执行下面命令：
```bash
# command-line-arguments
src\escape.go:5:9: &x escapes to heap
src\escape.go:4:6: moved to heap: x
src\escape.go:9:10: bar new(int) does not escape
```
`foo()`中的`x`最后在堆上分配，而`bar()`中的`x`最后分配在了栈上。
在官网 (golang.org) FAQ 上有一个关于变量分配的问题如下：
**如何得知变量是分配在栈（stack）上还是堆（heap）上**？
> 准确地说，你并不需要知道。Golang 中的变量只要被引用就一直会存活，**存储在堆上还是栈上由内部实现决定而和具体的语法没有关系**。
>
> 知道变量的存储位置确实和效率编程有关系。如果可能，Golang 编译器会将函数的局部变量分配到函数栈帧（stack frame）上。然而，**如果编译器不能确保变量在函数`return`之后不再被引用，
> 编译器就会将变量分配到堆上。而且，如果一个局部变量非常大，那么它也应该被分配到堆上而不是栈上**。
>
> 当前情况下，**如果一个变量被取地址，那么它就有可能被分配到堆上**。然而，还要对这些变量做逃逸分析，**如果函数`return`之后，变量不再被引用，则将其分配到栈上**。

### 关键数据结构
- mcache: per-P cache，可以认为是 local cache。
- mcentral: 全局 cache，mcache 不够用的时候向 mcentral 申请。
- mheap: 当 mcentral 也不够用的时候，通过 mheap 向操作系统申请。

多级内存分配器。

#### mcache
每个 Gorontine 的运行都是绑定到一个 P 上面的，**mcache 是每个 P 的 cache**。这么做的好处是**分配内存时不需要加锁**。mcache 结构如下。

#### mcentral
到当 mcache 不够用的时候，会从 mcentral 申请。

虽然在上面我们将 mcentral 和 mheap 作为两个部分来讲，但是作为全局的结构，这两部分是可以定义在一起的。实际上也是这样，mcentral 包含在 mheap 中。

#### 内存分配
给对象 object 分配内存的主要流程：

1. object size > 32K，则使用 mheap 直接分配。
2. object size < 16 byte，使用 mcache 的小对象分配器 tiny 直接分配。 （其实 tiny 就是一个指针，暂且这么说吧。）
3. object size > 16 byte && size <= 32K byte 时，先使用 mcache 中对应的 size class 分配。
4. 如果 mcache 对应的 size class 的 span 已经没有可用的块，则向 mcentral 请求。
5. 如果 mcentral 也没有可用的块，则向 mheap 申请，并切分。
6. 如果 mheap 也没有合适的 span，则向操作系统申请。

## 垃圾回收
Go语言中使用的垃圾回收使用的是**标记清扫算法**。标记清理最典型的做法是三⾊标记。进行垃圾回收时会stop the world。

三色标记算法原理如下：
1. 起初所有对象都是白色。
2. 从根出发扫描所有可达对象，标记为灰色，放入待处理队列。
3. 从队列取出灰色对象，将其引用对象标记为灰色放入队列，自身标记为黑色。
4. 重复 3，直到灰色对象队列为空。此时白色对象即为垃圾，进行回收。

### 何时触发 GC
#### 自动垃圾回收
在堆上分配大于 32K byte 对象的时候进行检测此时是否满足垃圾回收条件，如果满足则进行垃圾回收。
```go
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    ...
    shouldhelpgc := false
    // 分配的对象小于 32K byte
    if size <= maxSmallSize {
        ...
    } else {
        shouldhelpgc = true
        ...
    }
    ...
    // gcShouldStart() 函数进行触发条件检测
    if shouldhelpgc && gcShouldStart(false) {
        // gcStart() 函数进行垃圾回收
        gcStart(gcBackgroundMode, false)
    }
}
```

#### 主动垃圾回收
主动垃圾回收，通过调用 runtime.GC()，这是阻塞式的。
```go
// GC runs a garbage collection and blocks the caller until the
// garbage collection is complete. It may also block the entire
// program.
func GC() {
    gcStart(gcForceBlockMode, false)
}
```

### GC 触发条件
触发条件主要关注下面代码中的中间部分：`forceTrigger || memstats.heap_live >= memstats.gc_trigger`。
`forceTrigger`是 forceGC 的标志；后面半句的意思是当前堆上的活跃对象大于我们初始化时候设置的 GC 触发阈值。在 malloc 以及 free 的时候 heap_live 会一直进行更新。
```go
// gcShouldStart returns true if the exit condition for the _GCoff
// phase has been met. The exit condition should be tested when
// allocating.
//
// If forceTrigger is true, it ignores the current heap size, but
// checks all other conditions. In general this should be false.
func gcShouldStart(forceTrigger bool) bool {
    return gcphase == _GCoff && (forceTrigger || memstats.heap_live >= memstats.gc_trigger) && memstats.enablegc && panicking == 0 && gcpercent >= 0
}

//初始化的时候设置 GC 的触发阈值
func gcinit() {
    _ = setGCPercent(readgogc())
    memstats.gc_trigger = heapminimum
    ...
}
// 启动的时候通过 GOGC 传递百分比 x
// 触发阈值等于 x * defaultHeapMinimum (defaultHeapMinimum 默认是 4M)
func readgogc() int32 {
    p := gogetenv("GOGC")
    if p == "off" {
        return -1
    }
    if n, ok := atoi32(p); ok {
        return n
    }
    return 100
}
```

### golang垃圾回收使用的标记清理
#### STW(stop the world）
在扫描之前执⾏ STW（Stop The World）操作，就是**Runtime把所有的线程全部冻结掉，所有的线程全部冻结掉意味着⽤户逻辑肯定都是暂停的，所有的⽤户对象都不会被修改了**，
这时候去扫描肯定是安全的，对象要么活着要么死着，所以会造成在 STW 操作时所有的线程全部暂停，⽤户逻辑全部停掉，中间暂停时间可能会很⻓，⽤户逻辑对于⽤户的反应就中⽌了。

如何减短这个过程呢， STW过程中有两部分逻辑可以分开处理。我们看⿊⽩对象，扫描完结束以后对象只有⿊⽩对象，⿊⾊对象是接下来程序恢复之后需要使⽤的对象，如果不碰⿊⾊对象只回
收⽩⾊对象的话肯定不会给⽤户逻辑产⽣关联，因为⽩⾊对象肯定不会被⽤户线程引⽤的，所以回收操作实际上可以和⽤户逻辑并发的，因为可以保证回收的所有目标都不会被⽤户线程使⽤，
所以第⼀步回收操作和⽤户逻辑可以并发，因为我们回收的是⽩⾊对象，扫描完以后⽩⾊对象不会被全局变量引⽤、线程栈引⽤。回收⽩⾊对象肯定不会对⽤户线程产⽣竞争，⾸先**回收操作
肯定可以并发的，既然可以和⽤户逻辑并发，这样回收操作不放在 STW时间段⾥⾯缩短 STW 时间**。

#### 写屏障 (write barrier)
**该屏障之前的写操作和之后的写操作相比，先被系统其它组件感知**。

刚把⼀个对象标记为⽩⾊的，⽤户逻辑执⾏了突然引⽤了它，或者说刚刚扫描了 100 个对象正准备回收结果⼜创建了1000个对象在⾥⾯，因为没法结束没办法扫描状态不稳定，像扫描操作就⽐较⿇烦。
于是引⼊了写屏障的技术。

先做⼀次很短暂的STW，为什么需要很短暂的呢，它⾸先要执⾏⼀些简单的状态处理，接下来对内存进⾏扫描，这个时候⽤户逻辑也可以执⾏。⽤户所有新建的对象认为就是⿊⾊的，这次不扫描了下次再说，
新建对象不关⼼了，剩下来处理已经扫描过的对象是不是可能会出问题，已经扫描后的对象可能因为⽤户逻辑造成对象状态发⽣改变，所以**对扫描过后的对象使⽤操作系统写屏障功能⽤来监控⽤户
逻辑这段内存。任何时候这段内存发⽣引⽤改变的时候就会造成写屏障发⽣⼀个信号，垃圾回收器会捕获到这样的信号后就知道这个对象发⽣改变，然后重新扫描这个对象，看看它的引⽤或者被
引⽤是否被改变，这样利⽤状态的重置从⽽实现当对象状态发⽣改变的时候依然可以判断它是活着的还是死的**，这样扫描操作实际上可以做到⼀定程度上的并发，因为它没有办法完全屏蔽STW起
码它当开始启动先拿到⼀个状态，但是它的确可以把扫描时间缩短，现在知道了扫描操作和回收操作都可以⽤户并发。

#### golang回收的本质
实际上把单次暂停时间分散掉了，本来程序执⾏可能是“⽤户逻辑、⼤段GC、⽤户逻辑”，那么分散以后实际上变成了“⽤户逻辑、⼩段 GC、⽤户逻辑、⼩段GC、⽤户逻辑”这样。其实这个很难说 GC 快了。
因为被分散各个地⽅以后可能会频繁的保存⽤户状态，因为垃圾回收之前要保证⽤户状态是稳定的，原来只需要保存⼀次就可以了现在需要保存多次。
