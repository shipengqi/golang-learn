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

## 如何调试 CoreDump 文件

可以使用 delve 来进行调试：

```
$ go install github.com/go-delve/delve/cmd/dlv@latest
```