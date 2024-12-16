---
title: Go Modules
weight: 8
---

Go 在 1.11 推出了 Go Modules，这是一个新的包管理器，解决了 `GOPATH` 存在的问题。并且 Go 1.13 起不再推荐使用 GOPATH。

## Go Modules 机制

Go Modules 将依赖缓存放在 `$GOPATH/pkg/mod` 目录，并且**同一个依赖的版本，只会缓存一份，供所有项目使用**。

### 启用 Go Modules

Go 1.11 引入了环境变量 `GO111MODULE` 来控制是否启用 Go Modules，`GO111MODULE` 有三个值可选：

- `on` 启用 Go Modules
- `off` 禁用 Go Modules
- `auto`，在 `GOPATH` 下的项目，使用 `GOPATH`，否则启用 Go Modules。

Go 1.16 之前 `GO111MODULE` 的默认值是 `auto`，Go 1.16 起 `GO111MODULE` 的默认值为 `on`。

### 初始化

初始化 Go Modules 项目，首先要开启 Go Modules，然后在项目目录下运行：

```
$ go mod init <project-path>
```

### 下载依赖

下载依赖使用 `go get` 命令，命令格式为 `go get <package[@version]>`。

- `go get golang.org/x/test@latest`，`@latest` 表示选择最新的稳定版本，例如 `v1.2.3`。如果没有稳定版本，选择最新的预发布版本，例如 `v1.2.3-alpha.1`。
如果依赖没有 tag，那么选择最新的 commit。
- `go get golang.org/x/test` 同上。
- `go get golang.org/x/test@v1.2.3` 下载 tag 为 `v1.2.3` 的版本。
- `go get golang.org/x/test@v0` 下载 tag 前缀为 `v0` 的版本。
- `go get golang.org/x/test@master` 下载 master 分支上最新的 commit。
- `go get golang.org/x/test@37s237s` 下载哈希值为 `37s237s` 的 commit，如果该 commit 存在对应的 tag，转换为 tag 并下载。

`go get -u` 更新现有的依赖。

### Go Modules 代理

国内是无法访问 `golang.org` 的，Go 1.13 引入了环境变量 `GOPROXY`，可以用来设置 Go Modules 的代理。

`GOPROXY` 的默认值为 `https://proxy.golang.org,direct`，`GOPROXY` 可以设置多个，用 `,` 分隔。

执行 `go get/install` 时会优先从代理服务器下载依赖。如果从一个代理服务器下载失败，当遇见 `direct` 时，表示回源到依赖的源地址去下载。

#### 设置 GOPROXY

使用 `go env -w GOPROXY=https://goproxy.cn,direct` 命令来设置 `GOPROXY` 的值。

### GOPRIVATE

如果项目有一个私有依赖，设置 `GOPROXY` 也无法访问，可以使用 `GOPRIVATE`。

比如 `GOPRIVATE=corp.example.com,github.com/pookt/demo` 表示前缀可以匹配 `corp.example.com` 或者 `github.com/pookt/demo` 的依赖都会被认为是私有依赖。

`GOPRIVATE` 支持通配符，例如 `*.example.com`。

`GOPRIVATE` 较为特殊，它的值将作为 `GONOPROXY` 和 `GONOSUMDB` 的默认值。所以只使用 `GOPRIVATE` 就足够。

### go.mod 文件

`go.mod` 是 Go Modules 项目所必须的最重要的文件，描述了当前项目的元信息，目前有 5 个关键字：

- `module`：定义当前项目的模块路径。
- `go`：预期的 Go 版本。
- `require`：指定项目的依赖版本，格式为`<依赖的路径> <版本> [// indirect]`。
- `exclude`：排除一个特定的依赖版本。
- `replace`：将一个依赖版本替换为另外一个依赖版本，格式为 `module => newmodule`。

```
module example.com/foobar

go 1.13

require (
    example.com/apple v0.1.2
    example.com/pear v1.2.3
    example.com/watermelon v3.3.10+incompatible
    example.com/banana/v2 v2.3.4 // indirect
    example.com/pineapple v0.0.0-20190924185754-1b0db40df49a
)

exclude example.com/banana v1.2.4
replace example.com/apple v0.1.2 => example.com/rda v0.1.0 
replace example.com/banana => example.com/hugebanana
```

#### replace

`replace` 是用来将一个依赖版本替换为另外一个依赖版本，格式为 `module => newmodule`。

- `newmodule` 可以是本地相对路径，例如 `github.com/gin-gonic/gin => ./gin`。
- `newmodule` 也可以是本地绝对路径，例如 `github.com/gin-gonic/gin => /home/root/gin`。
- `newmodule` 可以是网络路径，例如 `golang.org/x/text v0.3.2 => github.com/golang/text v0.3.2`。

#### 依赖的导入路径说明

上面示例中 `example.com/banana/v2 v2.3.4`，`example.com/banana/v2` 的导入路径有 `/v2` 为什么其他依赖的导入路径没有 `/v0` 或者 `/v1`。

因为 Go modules 在主版本号为 `v0` 和 `v1` 的情况下省略了版本号，不需要在模块导入路径包含主版本的信息。而在主版本号为 `v2` 及以上则需要在导入路径末尾加上主版本号。

#### v0.0.0-xxx 是什么版本

Go 拉去的依赖如果没有 tag，那么选择最新的 commit。例如上面示例中的 `example.com/pineapple v0.0.0-20190924185754-1b0db40df49a`。

`v0.0.0` 是因为 `example.com/pineapple` 这个依赖不存在 tag，`20190924185754` 最新一次 commit 的 commit 时间，`1b0db40df49a` 是 commit 的哈希值。

#### indirect 

上面示例中的 `example.com/banana/v2 v2.3.4 // indirect`。`indirect` 表示该依赖为间接依赖。

通常上 `go.mod` 中出现的都应该是直接依赖，但是下面的两种情况会在 `go.mod` 中添加间接依赖：

- 当前项目的某个直接依赖没有使用 Go Modules。
- 当前项目的某个直接依赖的 `go.mod` 文件中缺失某个依赖，那么这个缺失的依赖会被添加在当前项目的 `go.mod` 文件中，作为间接依赖。  

#### incompatible

上面示例中的 `example.com/watermelon v3.3.10+incompatible`。`incompatible` 表示该依赖的路径跟版本不符合规范，`v3.3.10` 版本按照规范，引用路径应该为 `example.com/watermelon/v3`。
所以 Go 会在版本后加上 `+incompatible`。


### go.sum 文件

`go.sum` 列出了当前项目所有直接或间接依赖的版本，记录每个依赖的哈希值，目的是为了保证项目所依赖的版本不会被篡改。


## go mod 命令

`go mod` 常用的几个子命令：

- `init`：初始化 `go.mod` 文件
- `tidy`：自动添加项目依赖，并移除无用的依赖
- `download`：下载依赖到本地缓存。
- `graph`：查看现有的依赖结构
- `why`：查看为什么需要一个依赖


### 迁移回 vendor 模式

`go mod vendor` 可以将 Go Modules 迁移回到模式。

这个命令并只是单纯地把 `go.sum` 中的所有依赖下载到 `vendor` 目录里。

再使用 `go build -mod=vendor` 来构建项目，因为在 Go Modules 模式下 `go build` 是屏蔽 vendor 机制的。

注意发布时需要带上 vendor 目录。

## 其他

### 设置 HTTP Proxy 却仍然无法下载依赖

通常如果设置了 HTTP Proxy，`go get/install` 会使用指定的代理去下载依赖，例如：

```bash
# windows
set http_proxy=http://[user]:[pass]@[proxy_ip]:[proxy_port]/
set https_proxy=http://[user]:[pass]@[proxy_ip]:[proxy_port]/

# linux
export http_proxy=http://[user]:[pass]@[proxy_ip]:[proxy_port]/
export https_proxy=http://[user]:[pass]@[proxy_ip]:[proxy_port]/
```

但是，如果拉取的依赖是使用 Git 作为源控制管理器，那么还需要配置 Git 的 Proxy，否则还是无法下载依赖：

```bash
git config --global http.proxy http://[user]:[pass]@[proxy_ip]:[proxy_port]/
git config --global https.proxy http://[user]:[pass]@[proxy_ip]:[proxy_port]/
```

### 清理缓存

`go clean -modcache` 可以用来清理所有缓存的依赖。
