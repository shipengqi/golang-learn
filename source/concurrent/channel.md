---
title: channel
---

# channel
```
Don’t communicate by sharing memory; share memory by communicating.
（不要通过共享内存来通信，而应该通过通信来共享内存。）
```
这是作为 Go 语言最重要的编程理念。

通道类型的值是**并发安全**的，这也是 **Go 语言自带的、唯一一个可以满足并发安全性的类型**。

`channels` 是 `goroutine` 之间的通信机制。`goroutine` 通过 `channel` 向另一个 `goroutine` 发送消息
`channel` 和 `goroutine` 结合，可以实现用通信代替共享内存的 `CSP` 模型。

创建 `channel`：
```go
ch := make(chan int)

ch = make(chan int, 3) // buffered channel with capacity 3
```

上面的代码中，`int` 代表这个 `channel` 要发送的数据的类型。第二个参数代表创建一带缓存的 `channel`，容量为 `3`。

**`channel` 的零值是 `nil`。关闭一个 `nil` 的 `channel` 会导致程序 `panic`**。

发送和接收两个操作使用 `<-` 运算符，一个左尖括号紧接着一个减号形象地代表了元素值的传输方向：
```go
// 发送一个值
ch <- x // 将数据 push 到 channel

// 接受一个值
x = <-ch // 取出 channel 的值并复制给变量x

<-ch // 接受的值会被丢弃
```

### close

使用 `close` 函数关闭 `channel`，`channel` 关闭后不能再发送数据，但是可以接受已经发送成功的数据，
如果 `channel` 中没有数据，那么返回一个零值。

注意，**`close` 函数不是一个清理操作，而是一个控制操作**，在确定这个 `channel` 不会继续发送数据时调用。

**因为关闭操作只用于断言不再向 `channel` 发送新的数据，所以只有在 "发送者" 所在的 `goroutine` 才会调用 `close` 函数**，
因此对一个只接收的 `channel` 调用 `close` 将是一个编译错误。

使用 `range` 循环可直接在 `channels` 上面迭代。它依次从 `channel` 接收数据，当 `channel` 被关闭并且没有值可接收时
跳出循环。
```go
naturals := make(chan int)
for x := 0; x < 100; x++ {
    naturals <- x
}
for x := range naturals {
    fmt.Println(x)
}
```

**注意上面的代码会报 `fatal error: all goroutines are asleep - deadlock!`。这个是死锁的错误，因为 `range` 不等到信
道关闭是不会结束读取的。也就是如果 `channel` 没有数据了，那么 `range` 就会阻塞当前 `goroutine`, 直到信道关闭，所以导
致了死锁**。

为了避免这种情况，对于有缓存的信道，显式地关闭信道:
```go
ch := make(chan int, 3)
ch <- 1
ch <- 2
ch <- 3

// 显式地关闭信道
close(ch)

for v := range ch {
    fmt.Println(v)
}
```

### 无缓存 channel
**无缓存 `channel` 也叫做同步 `channel`**，这是因为**如果一个 `goroutine` 基于一个无缓存 `channel` 发送数据，那么就会
阻塞，直到另一个 `goroutine` 在相同的 `channel` 上执行接收操作**。同样的，**如果一个 `goroutine` 基于一个无缓存 `channel` 
先执行了接受操作，也会阻塞，直到另一个 `goroutine` 在相同的 `channel` 上执行发送操作**。在 `channel` 成功传输之后，两个 
`goroutine` 之后的语句才会继续执行。

### 带缓存 channel
```go
ch = make(chan int, 3)
```
带缓存的 `channel` 内部持有一个元素队列。`make` 函数创建 `channel` 时通过第二个参数指定队列的最大容量。

发送操作会向 `channel` 的缓存队列 `push` 元素，接收操作则是 `pop` 元素，如果队列被塞满了，那么发送操作将阻
塞直到另一个 `goroutine` 执行接收操作而释放了新的队列空间。
相反，如果 `channel` 是空的，接收操作将阻塞直到有另一个 `goroutine` 执行发送操作而向队列插入元素。

在大多数情况下，缓冲通道会作为收发双方的中间件。正如前文所述，元素值会先从发送方复制到缓冲通道，之后再由缓冲通道复制给接收方。

但是，当发送操作在执行的时候发现空的通道中，正好有等待的接收操作，那么它会直接把元素值复制给接收方。

### 单向 channel

当一个 `channel` 作为一个函数参数时，它一般总是被专门用于**只发送或者只接收**。

类型 `chan<- int` 表示一个只发送 `int` 的 `channel`。相反，类型 `<-chan int` 表示一个只接收 `int` 的 `channel`。

```go
var uselessChan = make(chan<- int, 1)
```

### cap 和 len
`cap` 函数可以获取 `channel` 内部缓存的容量。
`len` 函数可以获取 `channel` 内部缓存有效元素的个数。

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
1. 对于同一个通道，发送操作之间是互斥的，接收操作之间也是互斥的。Go 语言的运行时系统（以下简称运行时系统）只会执行对同一个通
道的任意个发送操作中的某一个。直到这个元素值被完全复制进该通道之后，其他针对该通道的发送操作才可能被执行。
2. 发送操作和接收操作中对元素值的处理都是不可分割的。发送操作要么还没复制元素值，要么已经复制完毕，绝不会出现只复制了一部分
的情况。接收操作在准备好元素值的副本之后，一定会删除掉通道中的原值，绝不会出现通道中仍有残留的情况。
3. 发送操作在完全完成之前会被阻塞。接收操作也是如此。

**元素值从外界进入通道时会被复制。更具体地说，进入通道的并不是在接收操作符右边的那个元素值，而是它的副本**。

**对于通道中的同一个元素值来说，发送操作和接收操作之间也是互斥的。例如，虽然会出现，正在被复制进通道但还未复制完成的元素值，
但是这时它绝不会被想接收它的一方看到和取走**。

### 发送操作和接收操作在什么时候可能被长时间的阻塞
- 针对**缓冲通道**的情况。如果通道已满，那么对它的所有发送操作都会被阻塞，直到通道中有元素值被接收走。相对的，如果通道已空，
那么对它的所有接收操作都会被阻塞，直到通道中有新的元素值出现。这时，通道会通知最早等待的那个接收操作所在的 goroutine，
并使它再次执行接收操作。
- 对于**非缓冲通道**，情况要简单一些。无论是发送操作还是接收操作，一开始执行就会被阻塞，直到配对的操作也开始执行，才会继续传递。
- **对于值为 `nil` 的通道，不论它的具体类型是什么，对它的发送操作和接收操作都会永久地处于阻塞状态**。它们所属的 goroutine 
中的任何代码，都不再会被执行。注意，由于通道类型是引用类型，所以它的零值就是 `nil`。**当我们只声明该类型的变量但没
有用 `make` 函数对它进行初始化时，该变量的值就会是 `nil`。我们一定不要忘记初始化通道**！

### select 多路复用
`select` 语句是专为通道而设计的，**所以每个 `case` 表达式中都只能包含操作通道的表达式**，比如接收表达式。

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

如果有多个 `channel` 需要接受消息，如果第一个 `channel` 没有消息发过来，那么程序会被阻塞，第二个 `channel` 的消息就也
无法接收了。这时候就需要使用 `select` 多路复用。
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
每一个 `case` 代表一个通信操作，发送或者接收。**如果没有 `case` 可运行，它将阻塞，直到有 `case` 可运行**。
如果多个 `case` 同时满足条件，`select` 会**随机**地选择一个执行。

**为了避免因为发送或者接收导致的阻塞，尤其是当 `channel` 没有准备好写或者读时。`default` 可以设置当其它的操作
都不能够马上被处理时程序需要执行哪些逻辑**。

### 超时
我们可以利用 `select` 来设置超时，避免 `goroutine` 阻塞的情况：
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

#### 使用 select 语句的时候，需要注意的事情
1. 如果加入了默认分支，那么无论涉及通道操作的表达式是否有阻塞，`select` 语句都不会被阻塞。如果那几个表达式都阻塞了，或者
说都没有满足求值的条件，那么默认分支就会被选中并执行。
2. 如果没有加入默认分支，那么一旦所有的 `case` 表达式都没有满足求值条件，那么 `select` 语句就会被阻塞。
直到至少有一个 `case` 表达式满足条件为止。
3. 还记得吗？我们可能会因为通道关闭了，而直接从通道接收到一个其元素类型的零值。所以，在很多时候，我们需要通过接收表达式
的第二个结果值来判断通道是否已经关闭。一旦发现某个通道关闭了，我们就应该及时地屏蔽掉对应的分支或者采取其他措施。这对
于程序逻辑和程序性能都是有好处的。
4. `select` 语句只能对其中的每一个 `case` 表达式各求值一次。所以，如果我们想连续或定时地操作其中的通道的话，就往往需要
通过在 `for` 语句中嵌入 `select` 语句的方式实现。但这时要注意，**简单地在 `select` 语句的分支中使用 `break` 语句，只能结
束当前的 `select` 语句的执行，而并不会对外层的 `for` 语句产生作用。这种错误的用法可能会让这个 `for` 语句无休止地运行下去**。

`break` 退出嵌套循环：
```go
I:
	for i := 0; i < 2; i++ {
		for j := 0; j < 5; j++ {
			if j == 2 {
				break I
			}
			fmt.Println("hello")
		}
		fmt.Println("hi")
	}
```

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

上面的代码 `select` 语句只有一个候选分支，我在其中利用接收表达式的第二个结果值对 `intChan` 通道是否已关闭做了判断，并在
得到肯定结果后，通过 `break` 语句立即结束当前 `select` 语句的执行。
