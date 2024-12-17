---
title: select
weight: 11
draft: true
---

`select` 类似于用于通信的 `switch` 语句。每个 `case` 必须是一个通信操作，要么是发送要么是接收。

当条件满足时，`select` 会去通信并执行 `case` 之后的语句，这时候其它通信是不会执行的。
如果多个 `case` 同时满足条件，`select` 会随机地选择一个执行。如果没有 `case` 可运行，它将阻塞，直到有 `case` 可运行。

一个默认的子句应该总是可运行的。

```go
select {
  case communication clause:
      ...
  case communication clause:
      ...
  default: /* 可选 */
   ...
}   
```

`for range` 支持遍历数组，切片，字符串，字典，通道，并返回索引和键值。**`for range` 会复制目标数据。可改用数组指针或者切片**。

`range` 关键字右边的位置上的代码被称为 `range` 表达式。

1. **`range` 表达式只会在 `for` 语句开始执行时被求值一次，无论后边会有多少次迭代**；
2. `range` 表达式的求值结果会被复制，也就是说，被迭代的对象是 `range` 表达式结果值的副本而不是原值。
3. `for range` 在性能比 `for` 稍差，因为 `for range` 会进行值拷贝。

字符串的复制成本很小，切片，字典，通道等引用类型本身是指针的封装，复制成本也很小，无序专门优化。

**如果 `range` 的目标表达式是函数，也只会运行一次**。

```go
numbers1 := []int{1, 2, 3, 4, 5, 6}
for i := range numbers1 {
    if i == 3 {
        numbers1[i] |= i
    }
}
fmt.Println(numbers1)
```

打印的内容会是 `[1 2 3 7 5 6]`，为什么，首先 `i` 是切片的下标，当 `i` 的值等于 3 的时候，与之对应的是切片中的第 4 个元素
值 4。对 4 和 3 进行按位或操作得到的结果是 7。

当 `for` 语句被执行的时候，在 `range` 关键字右边的 `numbers1` 会先被求值。`range` 表达式的结果值可以是数组、数组的指针、
切片、字符串、字典或者允许接收操作的通道中的某一个，并且结果值只能有一个。这里的 `numbers1` 是一个切片,那么迭代变量就可以
有两个，右边的迭代变量代表当次迭代对应的某一个元素值，而左边的迭代变量则代表该元素值在切片中的索引值。
循环控制语句：

- `break`，用于中断当前 `for` 循环或跳出 `switch` 语句
- `continue`，跳过当前循，继续进行下一轮循环。
- `goto`，将控制转移到被标记的语句。通常与条件语句配合使用。可用来实现条件转移， 构成循环，跳出循环体等功能。不推荐
  使用，以免造成流程混乱。

`goto` 实例：

```go
LOOP: for a < 20 {
 if a == 15 {
   /* 跳过迭代 */
   a = a + 1
   goto LOOP
 }
 fmt.Printf("a的值为 : %d\n", a)
 a ++  
}  
```



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

我们可以利用 `select` 来设置超时，避免 goroutine 阻塞的情况：

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
3. 还记得吗？我们可能会因为通道关闭了，而直接从通道接收到一个其元素类型的零值。所以，**在很多时候，我们需要通过接收表达式
   的第二个结果值来判断通道是否已经关闭**。一旦发现某个通道关闭了，我们就应该及时地屏蔽掉对应的分支或者采取其他措施。这对
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
    if !ok { // 使用 ok-idom，判断 channel 是否被关闭
      fmt.Println("The candidate case is closed.")
      break
    }
    fmt.Println("The candidate case is selected.")
}
```

上面的代码 `select` 语句只有一个候选分支，我在其中利用接收表达式的第二个结果值对 `intChan` 通道是否已关闭做了判断，并在
得到肯定结果后，通过 `break` 语句立即结束当前 `select` 语句的执行。
