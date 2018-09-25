# 并发编程
## Goroutines
`goroutine`可以简单理解为一个线程，我们程序运行的`main`函数在一个单独的`goroutine`中运行，叫做`main goroutine`。
在代码中可以使用`go`关键字创建`goroutine`。
```go
go f()
```

主函数返回时，所有`goroutine`都会被打断，程序退出。除了从主函数退出或者直接终止程序之外，没有其它的编程方法能够让一个`goroutine`来打断另一个的执行，但是之后可以看到一种方式来实现这个目的，通过`goroutine`之间的通信来让一个`goroutine`请求其它的`goroutine`，使被请求`goroutine`自行结束执行。