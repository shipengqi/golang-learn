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

`// +build` 功能和 `//go:build` 一样。只不过 `//go:build` 是在 go 1.17 才引入的。与其他现有 Go 指令保持一致，例如 `//go:generate`。

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

```sh
go build [-o output] [-i] [build flags] [packages]
```

- `-a` 强制重新编译所有包
- `-n` 把需要执行的编译命令打印出来，但是不执行，这样就可以很容易的知道底层是如何运行的
- `-p n` 指定可以并行可运行的编译数目，默认是 CPU 的数目
- `-o` 指定输出的可执行文件的文件名，可以带路径，例如 `go build -o a/b/c`
- `-i` 安装相应的包，编译并且 `go install`
- `-race` 开启编译的时候自动检测数据竞争的情况，目前只支持 64 位的机器
- `-v` 打印出来我们正在编译的包名
- `-work` 打印出来编译时候的临时文件夹名称，并且如果已经存在的话就不要删除
- `-x` 打印出来执行的命令，其实就是和-n的结果类似，只是这个会执行
- `-ccflags 'arg list'` 传递参数给 5c, 6c, 8c 调用
- `-compiler name` 指定相应的编译器，gccgo 还是 gc
- `-gccgoflags 'arg list'` 传递参数给 gccgo 编译连接调用
- `-gcflags 'arg list'` 编译器参数
- `-installsuffix suffix` 为了和默认的安装包区别开来，采用这个前缀来重新安装那些依赖的包，`-race`的时候默认已经是 `-installsuffix race`,大家可以通过 `-n` 命令来验证
- `-ldflags 'arg list'` 链接器参数
- `-tags 'tag list'` 设置在编译的时候可以适配的那些tag，详细的tag限制参考里面的 Build Constraints

### gcflags

`-gcflags` 参数的格式是

```bash
-gcflags="pattern=arg list"
```

#### pattern

pattern 是选择包的模式，它可以有以下几种定义:

- `main`: 表示 `main` 函数所在的顶级包路径
- `all`: 表示 `GOPATH` 中的所有包。如果是 `go modules` 模式，则表示主模块和它所有的依赖，包括 `test` 文件的依赖
- `std`: 表示 Go 标准库中的所有包
- `...`: `...` 是一个通配符，可以匹配任意字符串(包括空字符串)。
    - `net/...` 表示 net 模块和它的所有子模块
    - `./...` 表示当前主模块和所有子模块
    - 如果 pattern 中包含了 `/` 和 `...`，那么就不会匹配 `vendor` 目录
      例如: `./...` 不会匹配 `./vendor` 目录。可以使用 `./vendor/...` 匹配 `vendor` 目录和它的子模块

`go help packages` 查看模式说明。

#### arg list

空格分隔，如果编译选项中含有空格，可以使用引号包起来。

- `-N`: 禁止编译器优化
- `-l`: 关闭内联 (`inline`)
- `-c`: `int` 编译过程中的并发数，默认是 `1`
- `-B` 禁用越界检查
- `-u` 禁用 unsafe
- `-S` 输出汇编代码
- `-m` 输出优化信息

### ldflags

- `-s` 禁用符号表
- `-w` 禁用 DRAWF 调试信息
- `-X` 设置字符串全局变量值 `-X ver="0.99"`
- `-H` 设置可执行文件格式 `-H windowsgui`

## 减小编译体积

Go 编译器默认编译出来的程序会带有符号表和调试信息，一般来说 release 版本可以去除调试信息以减小二进制体积。

使用 `-w` 和 `-s` 来减少可执行文件的体积。但删除了调试信息后，可执行文件将无法使用 gdb/dlv 调试：

```bash
go build -ldflags="-w -s" ./abc.go
```

### 使用 upx

[upx](https://github.com/upx/upx) 是一个常用的压缩动态库和可执行文件的工具，通常可减少 50-70% 的体积。

下载 [upx](https://github.com/upx/upx/releases) 后解压就可以使用了。

```
# 使用 upx
$ go build -o server main.go && upx -9 server

# 结合编译选项
go build -ldflags="-s -w" -o server main.go && upx -9 server
```

upx 的参数 `-9` 指的是压缩率，1 代表最低压缩率，9 代表最高压缩率。

upx 压缩后的程序和压缩前的程序一样，无需解压仍然能够正常地运行，这种压缩方法称之为**带壳压缩**。

压缩包含两个部分：

- 在程序开头或其他合适的地方插入解压代码
- 将程序的其他部分压缩

执行时，也包含两个部分：

- 首先执行的是程序开头的插入的解压代码，将原来的程序在内存中解压出来
- 再执行解压后的程序。

也就是说，upx 在程序执行时，会有额外的解压动作，不过这个耗时几乎可以忽略。