---
title: select
weight: 11
---

操作系统中的 `select` 是一个系统调用，用于在多个通道上等待 I/O 事件的发生。Go 语言中的 **`select` 也能够让 goroutine 同时等待多个 channel 可读或者可写**，在多个文件或者 Channel 状态改变之前，select 会一直阻塞当前线程或者 goroutine。

`select` 类似于用于通信的 `switch` 语句。每个 `case` 必须是一个通信操作，要么是发送要么是接收。

当条件满足时，`select` 会去通信并执行 `case` 之后的语句，这时候其它通信是不会执行的。如果**多个 `case` 同时满足条件，`select` 会随机地选择一个执行**。如果**没有 `case` 可运行，它将阻塞，直到有 `case` 可运行**。

**如果加入了 `default` 子句，无论涉及通道操作的表达式是否有阻塞，`select` 语句都不会被阻塞**。

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

## 超时

可以利用 `select` 来设置超时，避免 goroutine 阻塞的情况：

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