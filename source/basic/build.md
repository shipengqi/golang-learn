---
title: Golang 条件编译
---

Golang 支持两种条件编译方式：

- 编译标签( build tag)
- 文件后缀

## 编译标签

编译标签添加的规则：

1. a build tag is evaluated as the OR of space-separated options
2. each option evaluates as the AND of its comma-separated terms
3. each term is an alphanumeric word or, preceded by !, its negation

翻译了就是：

1. 编译标签由**空格分隔**的编译选项(options)以"或"的逻辑关系组成
2. 每个编译选项由**逗号分隔**的条件项以逻辑"与"的关系组成
3. 每个条件项的名字用字母+数字表示，在前面加 `!` 表示否定的意思

例子（编译标签要放在源文件顶部，与 package 空行隔开）

```go
// +build darwin freebsd netbsd openbsd

package testpkg
```

这个将会让这个源文件只能在支持 kqueue 的 BSD 系统里编译

一个源文件里可以有多个编译标签，多个编译标签之间是逻辑"与"的关系

```go
// +build linux darwin
// +build 386
```

这个将限制此源文件只能在 `linux/386` 或者 `darwin/386` 平台下编译.

除了添加系统相关的 tag，还可以自由添加自定义 tag 达到其它目的。
编译方法:
只需要在 `go build` 指令后用 `-tags` 指定编译条件即可

```bash
go build -tags mytag1 mytag2
```

`!` 表示否定的意思：

```go
// +build !linux

package testpkg // correct
```

不会在 linux 平台下编译。

## 文件后缀

这个方法通过改变文件名的后缀来提供条件编译，如果你的源文件包含后缀：`_GOOS.go`，那么这个源文件只会在这个平台下编译，`_GOARCH.go` 也是如此。这两个后缀可以结合在一起使用，但是要注意顺序：`_GOOS_GOARCH.go`， 不能反过来用：`_GOARCH_GOOS.go`.
例子如下：

```bash
mypkg_freebsd_arm.go // only builds on freebsd/arm systems
mypkg_plan9.go       // only builds on plan9
```

## 编译标签和文件后缀的选择

编译标签和文件后缀的功能上有重叠，例如一个文件名：`mypkg_linux.go` 包含了 `// +build linux` 将会出现冗余

通常情况下，如果源文件与平台或者 cpu 架构完全匹配，那么用文件后缀，例如：

```bash
mypkg_linux.go         // only builds on linux systems
mypkg_windows_amd64.go // only builds on windows 64bit platforms
```

相反，如果满足以下任何条件，那么使用编译标签：

- 这个源文件可以在超过一个平台或者超过一个 cpu 架构下可以使用
- 需要去除指定平台
- 有一些自定义的编译条件
