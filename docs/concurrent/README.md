## Goroutines
`goroutine`可以简单理解为一个线程，但是它比线程更小，十几个`goroutine`可能体现在底层就是五六个线程，Go语言内部帮你实现
了这些`goroutine`之间的内存共享。执行`goroutine`只需极少的栈内存(大概是`4~5KB`)，当然会根据相应的数据伸缩。也正因为如此，
可同时运行成千上万个并发任务。`goroutine`比`thread`更易用、更高效、更轻便。我们程序运行的`main`函数在一个单独的`goroutine`中运行，叫做`main goroutine`。
在代码中可以使用`go`关键字创建`goroutine`。
```go
go f()
```

主函数返回时，所有`goroutine`都会被打断，程序退出。除了从主函数退出或者直接终止程序之外，没有其它的编程方法能够让一个`goroutine`来打断另一个的执行，但是之后可以看到一种方式来实现这个目的，通过`goroutine`之间的通信来让一个`goroutine`请求其它的`goroutine`，使被请求`goroutine`自行结束执行。

## Channels
`channels`是`goroutine`之间的通信机制。`goroutine`通过`channel`向另一个`goroutine`发送消息。`channel`和`goroutine`结合，
可以实现用通信代替共享内存的`CSP`模型（Go的口头禅**不要使用共享数据来通信，使用通信来共享数据**）。

创建`channel`：
```go
ch := make(chan int)

ch = make(chan int, 3) // buffered channel with capacity 3
```

上面的代码中，`int`代表这个`channel`要发送的数据的类型。第二个参数代表创建一带缓存的`channel`，容量为`3`。
`channel`的零值是`nil`。

发送和接收两个操作使用`<-`运算符：
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

### 带缓存channel
```go
ch = make(chan int, 3)
```
带缓存的`channel`内部持有一个元素队列。`make`函数创建`channel`时通过第二个参数指定队列的最大容量。

发送操作会向`channel`的缓存队列`push`元素，接收操作则是`pop`元素，如果队列被塞满了，那么发送操作将阻塞直到另一个`goroutine`执行接收操作而释放了新的队列空间。相反，如果`channel`是空的，接收操作将阻塞直到有另一个`goroutine`执行发送操作而向队列插入元素。

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

### goroutines泄漏
`goroutines`被永远卡住，就会导致`goroutines`泄漏，例如当使用了无缓存的`channel`，`goroutines`因为`channel`的数据没有被接收而被卡住。
泄漏的`goroutines`不会被自动回收。


## select 多路复用
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
每一个`case`代表一个通信操作，发送或者接收。如果没有`case`可运行，它将阻塞，直到有`case`可运行。
如果多个`case`同时满足条件，`select`会随机地选择一个执行。

为了避免因为发送或者接收导致的阻塞，尤其是当`channel`没有准备好写或者读时。`default`可以设置当其它的操作都不能够马上被处理时程序需要执行哪些逻辑。

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

#### 读写锁
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


### sync.Once 初始化
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

### 竞争检查器
在`go build`，`go ru`n或者`go test`命令后面加上`-race`，就会使编译器创建一个你的应用的“修改”版。

会记录下每一个读或者写共享变量的`goroutine`的身份信息。记录下所有的同步事件，比如`go`语句，`channel`操作，
以及对`(*sync.Mutex).Lock`，`(*sync.WaitGroup).Wait`等等的调用。

由于需要额外的记录，因此构建时加了竞争检测的程序跑起来会慢一些，且需要更大的内存，即使是这样，这些代价对于很多生产环境的工作来说还是可以接受的。