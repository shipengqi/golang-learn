---
title: Go 环境配置
---

# Go 环境配置
## 安装
Windows 下安装，官网 [下载安装包](https://golang.org/dl/)，直接安装。
默认情况下 `.msi` 文件会安装在 `c:\Go` 目录下。安装完成后默认会将 `c:\Go\bin` 目录添加到 `PATH` 环境变量中。
并添加环境变量 `GOROOT`，值为 Go 安装根目录 `C:\Go\`。重启命令窗口生效。

打开 CMD 输入 `go` 命令，验证是否安装成功。否则检查环境变量 `Path` 和 `GOROOT`。

## 工作区
### GOROOT
环境变量 `GOROOT` 用来指定 Go 的安装目录，Go 的标准库也在这个位置。目录结构与 `GOPATH` 类似。

### GOPATH
我们安装好 Go 之后，**必须配置一个环境变量 `GOPATH`**，这个 `GOPATH` 路径是用来指定当前工作目录的。
**不能和 Go 的安装目录（`GOROOT`）一样**。

工作区的目录结构：
```bash
GOPATH/
    src/ # 源码目录
    bin/ # 存放编译后的可执行程序
    pkg/ # 存放编译后的包的目标文件
```

`GOPATH` 允许多个目录，当有多个目录时，请注意分隔符，多个目录的时候 Windows 是分号 `;`，Linux 系统是冒号 `:`，
当有多个 `GOPATH` 时，默认会将 `go get` 的内容放在第一个目录下。


## Go Module
golang 1.11 已经支持 Go Module。这是官方提倡的新的包管理，乃至项目管理机制，可以不再需要 `GOPATH` 的存在。

### Module 机制
Go Module 不同于以往基于 `GOPATH` 和 Vendor 的项目构建，其主要是通过 `$GOPATH/pkg/mod` 下的缓存包来对项目进行构建。
 Go Module 可以通过 `GO111MODULE` 来控制是否启用，`GO111MODULE` 有三种类型:
- `on` 所有的构建，都使用 Module 机制
- `off` 所有的构建，都不使用 Module 机制，而是使用 `GOPATH` 和 Vendor
- `auto` 在 GOPATH 下的项目，不使用 Module 机制，不在 `GOPATH` 下的项目使用

### 和 dep 的区别
- dep 是解析所有的包引用，然后在 `$GOPATH/pkg/dep` 下进行缓存，再在项目下生成 vendor，然后基于 vendor 来构建项目，
无法脱离 `GOPATH`。
- mod 是解析所有的包引用，然后在 `$GOPATH/pkg/mod` 下进行缓存，直接基于缓存包来构建项目，所以可以脱离 `GOPATH`

### 准备环境
- golang 1.11 的环境需要开启 `GO11MODULE` ，并且**确保项目目录不在 `GOPATH` 中**。
```sh
export GO111MODULE=on
```
- golang 1.12 只需要确保实验目录不在 `GOPATH` 中。
- 配置代理 `export GOPROXY=https://goproxy.io`。（如果拉取包失败，会报  `cannot find module for path xxx` 的错误）

### 迁移到 Go Module
```sh
# clone 项目, 不要在 `GOPATH` 中, 比如之前的项目的结构是 `GOPATH/src/cdf-mannager`
git clone https://github.com/xxx/cdf-mannager

# 删除 vender
cd cdf-mannager
rm -rf vender

# init
go mod init cdf-mannager

# 下载依赖 也可以不执行这一步， go run 或 go build 会自动下载
go mod download
```

Go 会把 `Gopkg.lock` 或者 `glide.lock` 中的依赖项写入到 `go.mod` 文件中。`go.mod` 文件的内容像下面这样：
```
module cdf-manager

require (
        github.com/fsnotify/fsnotify v1.4.7
        github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7
        github.com/gin-gonic/gin v0.0.0-20180814085852-b869fe1415e4
        github.com/golang/protobuf v0.0.0-20170601230230-5a0f697c9ed9
        github.com/hashicorp/hcl v1.0.0
        github.com/inconshreveable/mousetrap v0.0.0-20141017200713-76626ae9c91c
        github.com/json-iterator/go v0.0.0-20170829155851-36b14963da70
        github.com/lexkong/log v0.0.0-20180607165131-972f9cd951fc
        github.com/magiconair/properties v1.8.0
        github.com/mattn/go-isatty v0.0.0-20170307163044-57fdcb988a5c
        github.com/mitchellh/mapstructure v1.1.2
        github.com/pelletier/go-toml v1.2.0
        github.com/satori/go.uuid v0.0.0-20180103152354-f58768cc1a7a
        github.com/spf13/afero v1.1.2
        github.com/spf13/cast v1.3.0
        github.com/spf13/cobra v0.0.0-20180427134550-ef82de70bb3f
        github.com/spf13/jwalterweatherman v1.0.0
        github.com/spf13/pflag v1.0.3
        github.com/spf13/viper v0.0.0-20181207100336-6d33b5a963d9
        github.com/ugorji/go v1.1.2-0.20180831062425-e253f1f20942
        github.com/willf/pad v0.0.0-20160331131008-b3d780601022
        golang.org/x/sys v0.0.0-20190116161447-11f53e031339
        golang.org/x/text v0.3.0
        gopkg.in/go-playground/validator.v8 v8.0.0-20160718134125-5f57d2222ad7
        gopkg.in/yaml.v2 v2.2.2
)

```

## 开发工具
常用 IDE：
- LiteIDE
- Sublime
- GoLand
- VS Code

