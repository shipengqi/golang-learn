---
title: mod
---

# mod

Golang 在 1.11 推出了 Go Module。这是官方提倡的新的包管理，乃至项目管理机制，解决了 `GOPATH` 的问题，相当于弃用了 `GOPATH`。

## Go Module 机制
Go Module 不同于基于 `GOPATH` 和 Vendor 的项目构建，其主要是通过 `$GOPATH/pkg/mod` 下缓存的模块来对项目进行构建。
**同一个模块版本的数据只缓存一份，所有其他模块共享使用**。

可以使用 `go clean -modcache` 清理所有已缓存的模块版本数据。

## GO111MODULE
Go Module 目前是可选的，可以通过环境变量 `GO111MODULE` 来控制是否启用，`GO111MODULE` 有三种类型:
- `on` 所有的构建，都使用 Module 机制
- `off` 所有的构建，都不使用 Module 机制，而是使用 `GOPATH` 和 Vendor
- `auto` 在 `GOPATH` 下的项目，不使用 Module 机制，不在 `GOPATH` 下的项目使用
 

## GOPROXY
`GOPROXY` 用于设置 Go Module 代理。使 Go 在后续拉取模块版本时能够脱离传统的 VCS 方式从镜像站点快速拉取。它的值是一个以 `,` 
分割的 Go module proxy 列表。Golang 1.13 以后它有一个默认的值 `GOPROXY=https://proxy.golang.org,direct`，
但是 `proxy.golang.org` 在中国是无法访问的，可以执行 `go env -w GOPROXY=https://goproxy.cn,direct` 来替换这个值。

- `off`，当 `GOPROXY=off` 时禁止 Go 在后续操作中使用 Go module proxy。
- `direct`，值列表中的 `direct` 用于指示 Go 回源到模块版本的源地址去抓取(如 GitHub)。当值列表中上一个 Go module proxy 返
回 404 或 410 错误时，Go 自动尝试列表中的下一个 proxy，当遇见 `direct` 时回源源地址，遇见 EOF 时终止并抛
出 “invalid version: unknown revision...” 的错误。


## go.mod
`go.mod` 是 Go moduels 项目所必须的最重要的文件，描述了当前项目（也就是当前模块）的元信息，每一行都以一个动词开头，目前有 5 个动词:
- `module`：定义当前项目的模块路径。
- `go`：设置预期的 Go 版本。
- `require`：设置特定的模块版本。
- `exclude`：从使用中排除一个特定的模块版本。
- `replace`：将一个模块版本替换为另外一个模块版本。

```go
module example.com/foobar

go 1.13

require (
    example.com/apple v0.1.2
    example.com/banana v1.2.3
    example.com/banana/v2 v2.3.4
    example.com/pineapple v0.0.0-20190924185754-1b0db40df49a
)

exclude example.com/banana v1.2.4
replace example.com/apple v0.1.2 => example.com/rda v0.1.0 
replace example.com/banana => example.com/hugebanana
```

## go.sum

`go.sum` 类似于 dep 的 `Gopkg.lock`。列出了当前项目直接或间接依赖的所有模块版本，并写明了那些模块版本的 SHA-256 哈希值以备 Go 在今
后的操作中保证项目所依赖的那些模块版本不会被篡改。
```go
k8s.io/client-go v0.0.0-20190620085101-78d2af792bab h1:E8Fecph0qbNsAbijJJQryKu4Oi9QTp5cVpjTE+nqg6g=
k8s.io/client-go v0.0.0-20190620085101-78d2af792bab/go.mod h1:E95RaSlHr79aHaX0aGSwcPNfygDiPKOVXdmivCIZT0k=
```
上面示例中一个模块路径有两种，前者为 Go module 打包整个模块包文件 zip 后再进行 hash 值，而后者为针对 `go.mod` 的 hash 值。
他们两者，要不就是同时存在，要不就是只存在 `go.mod` hash。

当 Go 认为肯定用不到某个模块版本的时候就会省略它的 zip hash，就会出现不存在 zip hash，只存在 `go.mod` hash 的情况。

## Go Checksum Database
Go Checksum Database 用于保护 Go 从任何源拉到 Go 模块版本不会被篡改。详细可以查看 `go help module-auth`。

## GOSUMDB
`GOSUMDB` 是一个 Go checksum database 的值。当它等于 `off` 时表示禁止 Go 在后续操作中校验模块版本。

- 默认值 `sum.golang.org` 中国无法访问，可以将 `GOPROXY` 设置为 `goproxy.cn`。`goproxy.cn` 支持代理 `sum.golang.org`。

## go mod 命令
```sh
Go mod provides access to operations on modules.

Note that support for modules is built into all the go commands,
not just 'go mod'. For example, day-to-day adding, removing, upgrading,
and downgrading of dependencies should be done using 'go get'.
See 'go help modules' for an overview of module functionality.

Usage:

        go mod <command> [arguments]

The commands are:

        download    下载 go.mod 文件中指明的所有依赖到本地缓存
        edit        编辑 go.mod 文件
        graph       查看现有的依赖结构
        init        在当前目录生成 go.mod 文件
        tidy        添加依赖的模块，并移除无用的模块
        vendor      导出现有的所有依赖
        verify      校验一个模块是否被篡改过
        why         解释为什么需要一个模块

Use "go help mod <command>" for more information about a command.
```


## 关于私有 module
如果项目依赖了私有模块，`GOPROXY` 访问不到，可以使用 `GOPRIVATE`。

比如 `GOPRIVATE=*.corp.example.com` 表示所有模块路径以 `corp.example.com` 的下一级域名 (如 `team1.corp.example.com`) 为前缀的
模块版本都将不经过 Go module proxy 和 Go checksum database （**注意不包括 `corp.example.com` 本身**）。

`GOPRIVATE` 较为特殊，它的值将作为 `GONOPROXY` 和 `GONOSUMDB` 的默认值。所以只使用 `GOPRIVATE` 就足够。

## 迁移项目到 Go Module  
### 准备环境
1. 开启 `GO11MODULE`：`go env -w GO111MODULE=on`，**确保项目目录不在 `GOPATH` 中**。
2. 配置代理 `export GOPROXY=https://goproxy.cn,direct`。

### 迁移
```sh
# clone 项目, 不要在 `GOPATH` 中, 之前的项目的结构是 `GOPATH/src/cdf-mannager`
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

**如果是一个新项目，或者删除了 `Gopkg.lock` 文件，可以直接运行：**
```sh
go mod init cdf-mannager

# 拉取必须模块 移除不用的模块
go mod tidy
```

接下来就可以运行 `go run main.go` 了。

## 迁移到 vendor
如果不想使用 go mod 的缓存方式，可以使用 `go mod vendor` 回到使用的 vendor 目录进行包管理的方式。

这个命令并只是单纯地把 `go.sum` 中的所有依赖下载到 vendor 目录里。

再使用 `go build -mod=vendor` 来构建项目，因为在 go modules 模式下 `go build` 是屏蔽 vendor 机制的:

发布时需要带上 vendor 目录。

## 添加新依赖包
添加新依赖包有下面几种方式：
1. 直接修改 `go.mod` 文件，然后执行 `go mod download`。
2. 使用 `go get packagename@vx.x.x`，会自动更新 `go.mod` 文件的。
3. `go run`、`go build` 也会自动下载依赖。

`go get` 拉取新的依赖：


## 依赖包冲突问题
迁移后遇到了下面的报错：
```sh
../gowork/pkg/mod/github.com/gin-gonic/gin@v0.0.0-20180814085852-b869fe1415e4/binding/msgpack.go:12:2: unknown import path "github.com/ugorji/go/codec": ambiguous import: found github.com/ugorji/go/codec in multiple modules:
	github.com/ugorji/go v0.0.0-20170215201144-c88ee250d022 (/root/gowork/pkg/mod/github.com/ugorji/go@v0.0.0-20170215201144-c88ee250d022/codec)
	github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8 (/root/gowork/pkg/mod/github.com/ugorji/go/codec@v0.0.0-20181204163529-d75b2dcb6bc8)
```

通过 `go mod graph` 可以查看具体依赖路径：
```sh
github.com/spf13/viper@v1.3.2 github.com/ugorji/go/codec@v0.0.0-20181204163529-d75b2dcb6bc8
github.com/gin-gonic/gin@v1.3.1-0.20190120102704-f38a3fe65f10 github.com/ugorji/go@v1.1.1
```

可以看到 `viper` 和 `gin` 分别依赖了 `github.com/ugorji/go` 和 `github.com/ugorji/go/codec`。

应该是 `go` 把这两个 `path` 当成不同的模块引入导致的冲突。[workaround](https://github.com/ugorji/go/issues/279)。

## 管理 Go 的环境变量
- Golang 1.13 新增了 `go env -w` 用于写入环境变量，写入到 `$HOME/.config/go/env` （`os.UserConfigDir` 返回的路径）文件中。
- `go env -w` 不会覆盖系统环境变量。
- 建议删除 Go 相关的系统环境变量，使用 `go env -w` 配置。

## 控制包的版本
`go get` 进行包管理时：
- 拉取最新的版本(优先择取 tag)：`go get golang.org/x/text@latest`
- 拉取 master 分支的最新 commit：`go get golang.org/x/text@master`
- 拉取 tag 为 `v0.3.2` 的 commit：`go get golang.org/x/text@v0.3.2`
- 拉取 hash 为 342b231 的 commit，最终会被转换为 `v0.3.2`：`go get golang.org/x/text@342b2e`。因为 Go modules 会与 tag 进
行对比，若发现对应的 commit 与 tag 有关联，则进行转换。
- 用 `go get -u` 更新现有的依赖，`go get -u all` 更新所有模块。

## 发布 module
### 语义化版本
Golang 官方推荐的最佳实践叫做 semver（Semantic Versioning），也就是语义化版本。

就是一种清晰可读的，明确反应版本信息的版本格式。
```
版本格式：主版本号.次版本号.修订号
```

- 主版本号：做了不兼容的 API 修改
- 次版本号：向下兼容的新增功能
- 修订号： 向下兼容的问题修正。

形如 `vX.Y.Z`。

#### 语义化版本的问题
如果你使用和发布的包没有版本 tag 或者处于 1.x 版本，那么可能体会不到什么区别，主要的区别体现在 `v2.x` 以及更高版本的包上。

go module 的谦容性规则：**如果旧软件包和新软件包具有相同的导入路径，则新软件包必须向后兼容旧软件包**
也就是说如果导入路径不同，就无需保持兼容。

`v2.x` 表示发生了重大变化，无法保证向后兼容，这时就需要在包的导入路径的末尾附加版本信息：
```go
module my-module/v2

require (
  some/pkg/v2 v2.0.0
  some/pkg/v2/mod1 v2.0.0
  my/pkg/v3 v3.0.1
)
```
格式总结为 `pkgpath/vN`，其中 N 是大于 1 的主要版本号。代码里导入时也需要附带上这个版本信息，如 `import "some/my-module/v2"`。

### go.sum
npm 的 `package-lock.json` 会记录所有库的准确版本，来源以及校验和，发布时不需要带上它，因为内容过于详细会对版本控制以及变更记录
等带来负面影响。

`go.sum` 也有类似的作用，会记录当前 module 所有的顶层和间接依赖，以及这些依赖的校验和，从而提供一个可以 100% 复现的构建过程并对构建对
象提供安全性的保证。同时还会保留过去使用的包的版本信息，以便日后可能的版本回退，这一点也与普通的锁文件不同。

准确地说，`go.sum` 是一个构建状态跟踪文件。

所以应该把 **`go.sum` 和 `go.mod` 一同添加进版本控制工具的跟踪列表，同时需要随着你的模块一起发布**。

### 包版本
当发布一个 `v2.x` 版本的库时，需要进行以下操作：
1. 将 `module my-module` 改成 `module my-module/v2`
2. 将源代码中使用了 v2+ 版本包的 `import "my-module"` 改为 `import "my-module/v2"`
3. 仔细检查你的代码中所有 `my-module` 包的版本是否统一，修改那些不兼容的问题
4. 在 changelog 中仔细列出所有 breaking changes

官方推荐将 `v2.x` 版本放在在一个新分支来避免混淆，如：
```sh
git checkout -b v2 

go mod edit --module=github.com/example/my-module/v2

git commit go.mod -m "upgrade to v2"

git tag v2.0.0

git push --tags origin v2
```