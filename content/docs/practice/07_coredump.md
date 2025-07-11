---
title: Go Core Dump 调试
weight: 7
---

Go 也可以开启类似 C++ Core Dump 功能，Core Dump 是程序崩溃时的内存快照。程序崩溃时，可以帮助定位 crash 发生的原因。

## 开启 Core Dump 功能

在 Linux 中，可以通过 `ulimit -c` 查看 Core Dump 功能是否开启：

```
$ ulimit -c
0
```

输出为 0，表示未开启。

使用 `ulimit -c [size]` 来指定 core dump 文件的大小，也就是开启 Core Dump。`ulimit -c unlimited` 表示不限制 core dump 文件的大小。

例如，下面的命令是将 core dump 文件大小设置为 `1MB`：

```
$ ulimit -c 1048576
```

## 如何生成 Core Dump 文件

Go 提供的环境变量 `GOTRACEBACK` 可以用来控制程序崩溃时输出的详细程度。可选值有：

- `none`：不显示任何 goroutine 的堆栈信息。
- `single`：默认选项，显示当前 goroutine 的堆栈信息。
- `all`：显示所有用户创建的 goroutine 的堆栈信息。
- `system`：显示所有 goroutine 的堆栈信息，包括 runtime。
- `crash`：作用和 `system` 一样, 但是会生成 core dump 文件。

可以设置 `export GOTRACEBACK=crash` 来生成 core dump 文件。

编译时要确保使用编译器标志 `-N` 和 `-l` 来构建二进制文件，`-N` 和 `-l` 回禁用编译器优化，因为编译器优化会使调试变得困难。

```
$ go build -gcflags=all="-N -l"
```

## 如何调试 Core Dump 文件

```go
package main

import "math/rand"

func main() {
	var sum int
	for {
		n := rand.Intn(1e6)
		sum += n
		if sum % 42 == 0 {
			panic("panic for GOTRACEBACK")
		}
	}
}
```

上面的示例运行后会直接崩溃：

```
panic: panic for GOTRACEBACK

goroutine 1 [running]:
main.main()
	C:/Code/example.v1/system/coredump/main.go:21 +0x78
```

上面的堆栈信息没有太多有用的信息。

这时就可以使用环境变量 `GOTRACEBACK=crash` 是程序生成 core dump 文件。然后重新运行，现在就会已打印出所有 goroutine，包括 runtime：

```
GOROOT=C:\Program Files\Go #gosetup
GOPATH=C:\Code\gowork #gosetup
"C:\Program Files\Go\bin\go.exe" build -o C:\Users\shipeng\AppData\Local\Temp\GoLand\___1go_build_github_com_shipengqi_example_v1_system_coredump.exe github.com/shipengqi/example.v1/system/coredump #gosetup
C:\Users\shipeng\AppData\Local\Temp\GoLand\___1go_build_github_com_shipengqi_example_v1_system_coredump.exe #gosetup
panic: panic for GOTRACEBACK

goroutine 1 [running]:
panic({0x4408c0, 0x45e5f8})
	C:/Program Files/Go/src/runtime/panic.go:1147 +0x3a8 fp=0xc000047f58 sp=0xc000047e98 pc=0x40ea08
main.main()
	C:/Code/example.v1/system/coredump/main.go:21 +0x78 fp=0xc000047f80 sp=0xc000047f58 pc=0x43be58
runtime.main()
	C:/Program Files/Go/src/runtime/proc.go:255 +0x217 fp=0xc000047fe0 sp=0xc000047f80 pc=0x411437
runtime.goexit()
	C:/Program Files/Go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000047fe8 sp=0xc000047fe0 pc=0x435921

goroutine 2 [force gc (idle)]:
runtime.gopark(0x0, 0x0, 0x0, 0x0, 0x0)
	C:/Program Files/Go/src/runtime/proc.go:366 +0xd6 fp=0xc000043fb0 sp=0xc000043f90 pc=0x4117d6
runtime.goparkunlock(...)
	C:/Program Files/Go/src/runtime/proc.go:372
runtime.forcegchelper()
	C:/Program Files/Go/src/runtime/proc.go:306 +0xb1 fp=0xc000043fe0 sp=0xc000043fb0 pc=0x411671
runtime.goexit()
	C:/Program Files/Go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000043fe8 sp=0xc000043fe0 pc=0x435921
created by runtime.init.7
	C:/Program Files/Go/src/runtime/proc.go:294 +0x25

goroutine 3 [GC sweep wait]:
runtime.gopark(0x0, 0x0, 0x0, 0x0, 0x0)
	C:/Program Files/Go/src/runtime/proc.go:366 +0xd6 fp=0xc000045fb0 sp=0xc000045f90 pc=0x4117d6
runtime.goparkunlock(...)
	C:/Program Files/Go/src/runtime/proc.go:372
runtime.bgsweep()
	C:/Program Files/Go/src/runtime/mgcsweep.go:163 +0x88 fp=0xc000045fe0 sp=0xc000045fb0 pc=0x3fc7e8
runtime.goexit()
	C:/Program Files/Go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000045fe8 sp=0xc000045fe0 pc=0x435921
created by runtime.gcenable
	C:/Program Files/Go/src/runtime/mgc.go:181 +0x55

goroutine 4 [GC scavenge wait]:
runtime.gopark(0x0, 0x0, 0x0, 0x0, 0x0)
	C:/Program Files/Go/src/runtime/proc.go:366 +0xd6 fp=0xc000055f80 sp=0xc000055f60 pc=0x4117d6
runtime.goparkunlock(...)
	C:/Program Files/Go/src/runtime/proc.go:372
runtime.bgscavenge()
	C:/Program Files/Go/src/runtime/mgcscavenge.go:265 +0xcd fp=0xc000055fe0 sp=0xc000055f80 pc=0x3fa8ed
runtime.goexit()
	C:/Program Files/Go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000055fe8 sp=0xc000055fe0 pc=0x435921
created by runtime.gcenable
	C:/Program Files/Go/src/runtime/mgc.go:182 +0x65
```

同级目录下会成一个文件名前缀是 `core` 的文件，然后就可以使用 delve 调试。

### 调试 

调试需要先安装 delve：

```
$ go install github.com/go-delve/delve/cmd/dlv@latest
```

然后执行命令 `dlv core <可执行文件> <core 文件>` 会进入交互模式：

```
$ dlv core main core.27507
Type 'help' for list of commands.
(dlv)
```

输入 `goroutines` 可以查看所有 goroutines 信息：

```
(dlv) goroutines
* goroutine 1 - User: ./main.go:11 main.main (0x47023e) (thread 27507)
  goroutine 2 - User: /usr/local/go/src/runtime/proc.go:399 runtime.gopark (0x439ffc) [force gc (idle)]
  goroutine 3 - User: /usr/local/go/src/runtime/proc.go:399 runtime.gopark (0x439ffc) [GC sweep wait]
  goroutine 4 - User: /usr/local/go/src/runtime/proc.go:399 runtime.gopark (0x439ffc) [GC scavenge wait]
[4 goroutines]
```

Goroutine 1 是 main goroutine，也是导致崩溃的 goroutine，输入 `goroutine 1` 切换到 goroutine 1 的栈帧：

```
(dlv) goroutine 1
Switched from 1 to 1 (thread 27507)
(dlv) 
```

执行 `bt` 查看详细的栈帧信息：

```
(dlv) bt
 0  0x0000000000465021 in runtime.raise
    at /usr/local/go/src/runtime/sys_linux_amd64.s:154
 1  0x000000000044c525 in runtime.dieFromSignal
    at /usr/local/go/src/runtime/signal_unix.go:903
 2  0x000000000044cbb5 in runtime.sigfwdgo
    at /usr/local/go/src/runtime/signal_unix.go:1108
 3  0x000000000044b485 in runtime.sigtrampgo
    at /usr/local/go/src/runtime/signal_unix.go:432
 4  0x0000000000465306 in runtime.sigtramp
    at /usr/local/go/src/runtime/sys_linux_amd64.s:352
 5  0x0000000000465400 in runtime.sigreturn__sigaction
    at /usr/local/go/src/runtime/sys_linux_amd64.s:471
 6  0x0000000000000001 in ???
    at ?:-1
 7  0x000000000044c712 in runtime.crash
    at /usr/local/go/src/runtime/signal_unix.go:985
 8  0x000000000043785e in runtime.fatalpanic
    at /usr/local/go/src/runtime/panic.go:1202
 9  0x0000000000436fb9 in runtime.gopanic
    at /usr/local/go/src/runtime/panic.go:1017
10  0x000000000047023e in main.main
    at ./main.go:11
11  0x0000000000439b87 in runtime.main
    at /usr/local/go/src/runtime/proc.go:267
12  0x0000000000463821 in runtime.goexit
    at /usr/local/go/src/runtime/asm_amd64.s:1650
(dlv) 
```

上面的输出中：

```
10  0x000000000047023e in main.main
    at ./main.go:11
```

可以定位到导致崩溃的代码在 `main.go`，然后输入 `frame 10` 进入具体的代码中：

```
(dlv) frame 10
> runtime.raise() /usr/local/go/src/runtime/sys_linux_amd64.s:154 (PC: 0x465021)
Warning: debugging optimized function
Frame 10: ./main.go:11 (PC: 47023e)
Warning: listing may not match stale executable
     6:         var sum int
     7:         for {
     8:                 n := rand.Intn(1e6)
     9:                 sum += n
    10:                 if sum % 42 == 0 {
=>  11:                         panic("panic for GOTRACEBACK")
    12:                 }
    13:         }
    14: }
(dlv) 
```

可以定位到第 11 行代码导致的 panic。
