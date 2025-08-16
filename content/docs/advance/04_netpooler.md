---
title: 网络轮询器
weight: 4
---

Go 语言在网络轮询器中使用 I/O 多路复用模型处理 I/O 操作。

## I/O 模型

### 阻塞 I/O

阻塞 I/O 是最常见的 I/O 模型，在默认情况下，当通过 `read` 或者 `write` 等系统调用读写文件或者网络时，应用程序会被阻塞。

例如，当执行 `read` 系统调用时，应用程序会从用户态陷入内核态，内核会检查文件描述符是否可读；当文件描述符中存在数据时，操作系统内核会将准备好的数据拷贝给应用程序并交回控制权。

操作系统中多数的 I/O 操作都是如上所示的阻塞请求，一旦执行 I/O 操作，应用程序会陷入阻塞等待 I/O 操作的结束。

### 非阻塞 I/O

当进程把一个文件描述符设置成非阻塞时，执行 `read` 和 `write` 等 I/O 操作会立刻返回一个 `EAGAIN` 错误。意味着该文件描述符还在等待缓冲区中的数据。随后，应用程序会不断**轮询**调用 `read` 直到它的返回值大于 0，这时应用程序就可以对读取操作系统缓冲区中的数据并进行操作。

### 多路复用

多路复用是指进程可以同时处理多个文件描述符，一旦某个描述符就绪（一般是读就绪或者写就绪），能够通知程序进行相应的读写操作。

系统调用 select 也可以提供 I/O 多路复用的能力，但是使用它有比较多的限制：

- 监听能力有限：最多只能监听 1024 个文件描述符；
- 内存拷贝开销大：需要维护一个较大的数据结构存储文件描述符，该结构需要拷贝到内核中；
- 时间复杂度 `O(n)`：返回准备就绪的事件个数后，需要遍历所有的文件描述符；

#### epoll

三大核心系统调用：

```c
// 创建 epol 实例
int epoll_create(int size); 

// 管理监控列表
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);

// 等待事件就绪
int epoll_wait(int epfd, struct epoll_event *events, int maxevents, int timeout);
```

### Go 的网络轮询实现

不同的操作系统也都实现了自己的 I/O 多路复用函数，例如：`epoll`、`kqueue` 和 `evport` 等。

Go 为了提高在不同操作系统上的 I/O 操作性能，使用平台特定的函数实现了多个版本的网络轮询模块：

- `src/runtime/netpoll_epoll.go`
- `src/runtime/netpoll_kqueue.go`
- `src/runtime/netpoll_solaris.go`
- `src/runtime/netpoll_windows.go`
- `src/runtime/netpoll_aix.go`
- `src/runtime/netpoll_fake.go`

编译器在编译 Go 程序时，会根据目标平台选择树中特定的分支进行编译。例如，如果目标平台是 Linux，那么就会根据文件中的 `// +build linux` 编译指令选择 `src/runtime/netpoll_epoll.go` 并使用 `epoll` 函数处理用户的 I/O 操作。

`epoll`、`kqueue` 等多路复用模块都要实现以下五个函数，这五个函数构成一个接口：

```go
// 初始化网络轮询器，通过 sync.Once 和 netpollInited 变量保证函数只会调用一次
func netpollinit() 
// 监听文件描述符上的边缘触发事件，创建事件并加入监听
func netpollopen(fd uintptr, pd *pollDesc) int32
// 轮询网络并返回一组已经准备就绪的 goroutine，传入的参数会决定它的行为
// 如果参数小于 0，无限期等待文件描述符就绪；
// 如果参数等于 0，非阻塞地轮询网络；
// 如果参数大于 0，阻塞特定时间轮询网络；
func netpoll(delta int64) gList
// 唤醒网络轮询器，例如：计时器向前修改时间时会通过该函数中断网络轮询器4；
func netpollBreak()
// 判断文件描述符是否被轮询器使用
func netpollIsPollDescriptor(fd uintptr) bool
```

#### 初始化

`runtime.netpollinit`，即 Linux 上的 `epoll`，它主要做了以下几件事情：

1. 是调用 `epollcreate1` 创建一个新的 `epoll` 文件描述符，这个文件描述符会在整个程序的生命周期中使用；
2. 通过 `runtime.nonblockingPipe `创建一个用于通信的管道；
3. 使用 `epollctl` 将用于读取数据的文件描述符打包成 `epollevent` 事件加入监听；

#### goroutine 与 epoll 的协作流程

goroutine 发起网络读操作（如 `conn.Read()`）：

```go
func handleConn(conn net.Conn) {
    buf := make([]byte, 1024)
    n, err := conn.Read(buf) // 非阻塞 I/O，触发 netpoller 介入
    // ...
}
```

- 当 goroutine 调用 `conn.Read()` 时，Go 的 net 包会先尝试**非阻塞读取**（`syscall.Read(fd, buf)`）。
- 如果数据未就绪（返回 `EAGAIN` 错误），Go 会将**文件描述符（fd） 注册到 epoll 监听队列**，并将当前 goroutine 阻塞。

goroutine 挂起：

```go
// 伪代码：Go 运行时内部处理
func Read(fd int, buf []byte) (int, error) {
    for {
        n, err := syscall.Read(fd, buf)
        if err == syscall.EAGAIN {
            // 将 fd 注册到 epoll，并挂起当前 Goroutine
            netpollqueue.Add(fd, currentG)
            gopark() // 挂起 Goroutine，切换执行其他任务
            continue
        }
        return n, err
    }
}

// 将 fd 关联到当前 goroutine
func netpollblock(fd uintptr) {
    pd := getPollDesc(fd)
    pd.setG(getg()) // 关联当前 Goroutine
    gopark(netpollblockcommit, unsafe.Pointer(pd), waitReasonIOWait)
}
```

- `gopark()` 将当前 goroutine 置为 Gwaiting 状态，并让出 CPU 给其他 goroutine。

epoll 监听就绪事件：

- 后台的 网络轮询器线程（通常由 Go 运行时启动）通过 `epoll_wait` 监听所有注册的 fd。
- 当数据到达时，epoll 返回就绪的 fd，轮询器找到关联的 goroutine，将其标记为可运行（Grunnable）。

goroutine 恢复执行：

```go
func netpoll(delay int64) []*g {
    events := epoll_wait(epfd, ...)
    for _, ev := range events {
        pd := (*pollDesc)(ev.data)
        netpollready(&toRun, pd) // 将关联的 Goroutine 加入运行队列
    }
    return toRun
}
```

- 调度器将就绪的 goroutine 放入运行队列，等待 M（OS 线程）执行。
- goroutine 从 `gopark()` 后继续执行，再次调用 `syscall.Read`，此时数据已就绪。