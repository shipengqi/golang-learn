## 并发编程
```
Don’t communicate by sharing memory; share memory by communicating.
（不要通过共享内存来通信，而应该通过通信来共享内存。）
```
这是作为 Go 语言的主要创造者之一的 Rob Pike 的至理名言，这也充分体现了 Go 语言最重要的编程理念。

### Goroutines
`goroutine`可以简单理解为一个线程，但是它比线程更小，十几个`goroutine`可能体现在底层就是五六个线程，Go语言内部帮你实现
了这些`goroutine`之间的内存共享。执行`goroutine`只需极少的栈内存(大概是`4~5KB`)，当然会根据相应的数据伸缩。也正因为如此，
可同时运行成千上万个并发任务。`goroutine`比`thread`更易用、更高效、更轻便。我们程序运行的`main`函数在一个单独的`goroutine`中运行，叫做`main goroutine`。
在代码中可以使用`go`关键字创建`goroutine`。
```go
go f()
```

主函数返回时，所有`goroutine`都会被打断，程序退出。除了从主函数退出或者直接终止程序之外，没有其它的编程方法能够让一个`goroutine`来打断另一个的执行，但是之后可以看到一种方式来实
现这个目的，通过`goroutine`之间的通信来让一个`goroutine`请求其它的`goroutine`，使被请求`goroutine`自行结束执行。


#### goroutines泄漏
`goroutines`被永远卡住，就会导致`goroutines`泄漏，例如当使用了无缓存的`channel`，`goroutines`因为`channel`的数据没有被接收而被卡住。
泄漏的`goroutines`不会被自动回收。

#### 什么是主 goroutine，它与我们启用的其他 goroutine 有什么不同
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

#### 怎样才能让主 goroutine 等待其他 goroutine
刚才说过，一旦主 goroutine 中的代码执行完毕，当前的 Go 程序就会结束运行，无论其他的 goroutine 是否已经在运行了。那么，怎样才能做到等其他的 goroutine 运行完毕之后，
再让主 goroutine 结束运行呢？

**使用`time`包**

可以简单粗暴的`time.Sleep(time.Millisecond * 500)`让主 goroutine“小睡”一会儿。在这里传入了“500 毫秒”

问题是我们让主 goroutine“睡眠”多长时间才是合适的呢？如果“睡眠”太短，则很可能不足以让其他的 goroutine 运行完毕，而若“睡眠”太长则纯属浪费时间，这个时间就太难把握了。

**使用通道**。

**使用`sync`包的`sync.WaitGroup`类型**

#### 怎样让启用的多个 goroutine 按照既定的顺序运行
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
### Channels
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

#### close

使用`close`函数关闭`channel`，`channel`关闭后不能再发送数据，但是可以接受已经发送成功的数据，如果`channel`中没有
数据，那么返回一个零值。

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
#### 无缓存channel
无缓存`channel`也叫做同步`channel`，这是因为如果一个`goroutine`基于一个无缓存`channel`发送数据，那么就会阻塞，直到
另一个`goroutine`在相同的`channel`上执行接收操作。同样的，如果一个`goroutine`基于一个无缓存`channel`先执行了接受操作，
也会阻塞，直到另一个`goroutine`在相同的`channel`上执行发送操作。在`channel`成功传输之后，两个`goroutine`之后的语句才会
继续执行。

#### 单向channel

当一个`channel`作为一个函数参数时，它一般总是被专门用于只发送或者只接收。

类型`chan<- int`表示一个只发送`int`的`channel`。相反，类型`<-chan int`表示一个只接收`int`的`channel`。

```go
var uselessChan = make(chan<- int, 1)
```

#### 带缓存channel
```go
ch = make(chan int, 3)
```
带缓存的`channel`内部持有一个元素队列。`make`函数创建`channel`时通过第二个参数指定队列的最大容量。

发送操作会向`channel`的缓存队列`push`元素，接收操作则是`pop`元素，如果队列被塞满了，那么发送操作将阻塞直到另一个`goroutine`执行接收操作而释放了新的队列空间。相反，如果`channel`是空的，接收操作将阻塞直到有另一个`goroutine`执行发送操作而向队列插入元素。

在大多数情况下，缓冲通道会作为收发双方的中间件。正如前文所述，元素值会先从发送方复制到缓冲通道，之后再由缓冲通道复制给接收方。

但是，当发送操作在执行的时候发现空的通道中，正好有等待的接收操作，那么它会直接把元素值复制给接收方。

#### cap 和 len
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

#### 通道的发送和接收操作的特性
1. 对于同一个通道，发送操作之间是互斥的，接收操作之间也是互斥的。，Go 语言的运行时系统（以下简称运行时系统）只会执行对同一个通道的任意个发送操作中的某一个。直到这个元素值
被完全复制进该通道之后，其他针对该通道的发送操作才可能被执行。
2. 发送操作和接收操作中对元素值的处理都是不可分割的。发送操作要么还没复制元素值，要么已经复制完毕，绝不会出现只复制了一部分的情况。接收操作在准备好元素值的副本
之后，一定会删除掉通道中的原值，绝不会出现通道中仍有残留的情况。
3. 发送操作在完全完成之前会被阻塞。接收操作也是如此。

**元素值从外界进入通道时会被复制。更具体地说，进入通道的并不是在接收操作符右边的那个元素值，而是它的副本**。

**对于通道中的同一个元素值来说，发送操作和接收操作之间也是互斥的。例如，虽然会出现，正在被复制进通道但还未复制完成的元素值，但是这时它绝不会被想接收它的一方看到和取走**。

#### 发送操作和接收操作在什么时候可能被长时间的阻塞
- 针对**缓冲通道**的情况。如果通道已满，那么对它的所有发送操作都会被阻塞，直到通道中有元素值被接收走。相对的，如果通道已空，那么对它的所有接收操作都会被阻塞，直到通道中有新的元素值出现。这时，通道会通知最早等待的那个接收操作所在的 goroutine，并使它再次执行接收操作。
- 对于**非缓冲通道**，情况要简单一些。无论是发送操作还是接收操作，一开始执行就会被阻塞，直到配对的操作也开始执行，才会继续传递。
- **对于值为nil的通道，不论它的具体类型是什么，对它的发送操作和接收操作都会永久地处于阻塞状态**。它们所属的 goroutine 中的任何代码，都不再会被执行。注意，由于通道类型是引用类型，所以它的零值就是nil。**当我们只声明该类型的变量但没有用make函数对它进行初始化时，该变量的值就会是nil。我们一定不要忘记初始化通道**！

### select 多路复用
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

#### 超时
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

#### 使用select语句的时候，需要注意的事情
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

### 共享变量
无论任何时候，只要有两个以上`goroutine`并发访问同一变量，且至少其中的一个是写操作的时候就会发生数据竞争。
避免数据竞争的三种方式：
1. 不去写变量。读取不可能出现数据竞争。
2. 避免从多个`goroutine`访问变量，尽量把变量限定在了一个单独的`goroutine`中。(**不要使用共享数据来通信，使用通信来共享数据**)
3. 互斥锁

#### 互斥锁
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
##### sync.Mutex
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

##### 读写锁
如果有多个`goroutine`读取变量，那么是并发安全的，这个时候使用`sync.Mutex`加锁就没有必要。可以使用`sync.RWMutex`读写锁（多读单写锁）。
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


#### sync.Once 初始化
```go
var loadIconsOnce sync.Once
var icons map[string]image.Image
// Concurrency-safe.
func Icon(name string) image.Image {
  loadIconsOnce.Do(loadIcons)
  return icons[name]
}
```
`Do`这个唯一的方法需要接收初始化函数作为其参数。

#### 竞争检查器
在`go build`，`go run`或者`go test`命令后面加上`-race`，就会使编译器创建一个你的应用的“修改”版。

会记录下每一个读或者写共享变量的`goroutine`的身份信息。记录下所有的同步事件，比如`go`语句，`channel`操作，
以及对`(*sync.Mutex).Lock`，`(*sync.WaitGroup).Wait`等等的调用。

由于需要额外的记录，因此构建时加了竞争检测的程序跑起来会慢一些，且需要更大的内存，即使是这样，这些代价对于很多生产环境的工作来说还是可以接受的。