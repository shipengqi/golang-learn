---
title: Go 编译
weight: 1
---

# Go 编译

## 条件编译

Go 支持两种条件编译方式：

- 编译标签（build tag）
- 文件后缀

### 编译标签

编译标签的规则：

1. 空格表示：AND
2. 逗号表示：OR
3. `!` 表示：NOT
4. 换行表示：AND

每个条件项的名字用 "字母+数字" 表示。主要支持以下几种条件：

- 操作系统，例如：`windows`、`linux` 等，对应 `runtime.GOOS` 的值。
- 计算机架构，例如：`amd64`、`386`，对应 `runtime.GOARCH` 的值。
- 编译器，例如：`gccgo`、`gc`，是否开启 CGO,cgo。
- Go 版本，例如：`go1.19`、`go1.20` 等。
- 自定义的标签，例如：编译时通过指定 `-tags` 传入的值。
- `//go:build ignore`，编译时自动忽略该文件

`go:build` 之后必须有空行，否则会被编译器当做普通注释。

```go
//go:build linux,386 darwin,!cgo

package testpkg
```

运算表达式为：`(linux && 386) || (darwin && !cgo)`。

自定义 tag 只需要在 `go build` 指令后用 `-tags` 指定编译条件即可

```bash
go build -tags mytag1 mytag2
```

对于 `-tags`，多个标签既可以用逗号分隔，也可以用空格分隔，但它们都表示"与"的关系。早期 go 版本用空格分隔，后来改成了用逗号分隔，但空格依然可以识别。

`-tags` 也有 `!` 规则，它表示的是没有这个标签。

```go
//go:build !hello
```

```bash
go build -tags=!hello
```

### 文件后缀

这个方法通过改变文件名的后缀来提供条件编译，如果你的源文件包含后缀：`_GOOS.go`，那么这个源文件只会在这个平台下编译，`_GOARCH.go` 也是如此。这两个后缀可以结合在一起使用，但是要注意顺序：`_GOOS_GOARCH.go`， 不能反过来用。
例如：

```bash
mypkg_freebsd_arm.go // only builds on freebsd/arm systems
mypkg_plan9.go       // only builds on plan9
```

文件名必须提供，如果只由后缀的文件名会被编译器忽略：

```
# 这个文件会被编译器忽略
_linux.go
```

### 如何选择编译标签和文件后缀

编译标签和文件后缀的功能上有重叠，例如一个文件名：`mypkg_linux.go` 包含了 `//go:build linux` 将会出现冗余

通常情况下，如果源文件与平台或者 cpu 架构完全匹配，那么使用文件后缀就可以满足，例如：

```bash
mypkg_linux.go         // only builds on linux systems
mypkg_windows_amd64.go // only builds on windows 64bit platforms
```

下面的情况，就可以使用编译标签：

- 这个源文件可以在超过一个平台或者超过一个 cpu 架构
- 需要排除某个平台或架构
- 有一些自定义的编译条件

### +build

`// +build` 功能和 `//go:build` 一样。只不过 `//go:build` 是在 go 1.17 才引入的。与其他现有 Go 指令和编译指示的一致性，例如 `//go:generate`。

## 交叉编译

Go 可以通过设置环境变量来实现交叉编译，用来在一个平台上生成另一个平台的可执行程序。：

```
#  linux amd64
GOOS=linux GOARCH=amd64 go build main.go

# windows amd64
GOOS=windows GOARCH=amd64 go build main.go
```

环境变量 `GOOS` 设置平台, `GOARCH` 设置架构。

## 编译选项