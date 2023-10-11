---
title: Go CoreDump 调试
weight: 7
---

# Go CoreDump 调试

Go 也可以开启类似 C++ CoreDump 功能，CoreDump 是异常退出程序的内存快照。程序崩溃时，可以帮助定位 crash 发生的原因。

## 如何生成 CoreDump 文件

`GOTRACEBACK` 可以控制程序崩溃时输出的详细程度。 可选的值：

- `none` 不显示任何 goroutine 栈 trace。
- `single`, 默认选项，显示当前 goroutine 栈 trace。
- `all` 显示所有用户创建的 goroutine 栈 trace。
- `system` 显示所有 goroutine 栈 trace,甚至运行时的 trace。
- `crash` 类似 system, 而且还会生成 core dump。

可以设置 `export GOTRACEBACK=crash` 来生成 core dump。

编译时要确保使用编译器标志 `-N` 和 `-l` 来构建二进制文件,它会禁用编译器优化，编译器优化可能会使调试更加困难。

```
$ go build -gcflags=all="-N -l"
```

如果 coredump 没有生成，可能是 coredump size 配置为 0，如下命令将 coredump 配置为 `1MB` 大小：

```
$ ulimit -c 1048576
```

## 如何调试 CoreDump 文件

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

上面的程序将很快崩溃
```
panic: panic for GOTRACEBACK

goroutine 1 [running]:
main.main()
	C:/Code/example.v1/system/coredump/main.go:21 +0x78
```

无法从上面的 panic 栈 trace 中分辨出崩溃所涉及的值。增加日志或许是一种解决方案，但是我们并不总是知道在何处添加日志。
添加环境变量 `GOTRACEBACK=crash` 再运行它。现在会已打印出所有 goroutine，包括 runtime，因此输出更加详细。 并输出 core dump：

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

需要调试，就可以使用 delve。

安装 delve：

```
$ go install github.com/go-delve/delve/cmd/dlv@latest
```

通过 `dlv core` 命令来调试 coredump。通过 `bt` 命令打印堆栈，并且展示程序造成的 panic。
