---
title: 代码规范
weight: 3
---

# 代码规范

好的代码规范非常重要，可以提高代码的可读性，减少 bug，提高开发效率。

Go 官方提供的代码规范：

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)

Uber 开源的 Go 编码规范：

- [Uber Go Guide](https://github.com/uber-go/guide)

Go 也提供了一些代码检查工具，例如 `golint`，`goimports`，`go vet` 等，但是这些工具检查的不够全面。

golangci-lint 是一个更加强大的静态代码检查工具。

## golangci-lint

golangci-lint 的运行速度非常快，因为它可以并行的运行 linters，并且重用 Go 的构建缓存，缓存分析结果。

golangci-lint 集成了大量的 [linters](https://golangci-lint.run/usage/linters/)，不需要额外安装，可以直接使用。

### 安装

```
$ go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
# 验证是否安装成功
$ golangci-lint version
```

[更多安装方式](https://golangci-lint.run/usage/install/)。

### 使用

`run` 命令执行代码检查：

```
$ golangci-lint run
```

`linters` 命令打印出 golangci-lint 所支持的 linters：

```
$ golangci-lint linters
```

### 配置

golangci-lint 有两种配置方式：[命令行选项](https://golangci-lint.run/usage/configuration/#command-line-options)和[配置文件](https://golangci-lint.run/usage/configuration/#config-file)。

golangci-lint 会在当前工作目录下的以下路径中查找配置文件：

- `.golangci.yml`
- `.golangci.yaml`
- `.golangci.toml`
- `.golangci.json`

一般会在项目的根目录下创建一个配置文件。配置文件示例：

```yaml
run:
  deadline: 2m

  # Include test files or not.
  # Default: true
  tests: false

linters:
  # Disable all linters.
  # Default: false
  disable-all: true
  # Enable specific linter
  # https://golangci-lint.run/usage/linters/#enabled-by-default
  enable:
    - misspell
    - govet
    - staticcheck
    - errcheck
    - unparam
    - ineffassign
    - nakedret
    - gocyclo
    - dupl
    - goimports
    - revive
    - gosec
    - gosimple
    - typecheck
    - unused

# https://golangci-lint.run/usage/linters
linters-settings:
  gofmt:
    simplify: true
  dupl:
    threshold: 600
```

### 误报

如果出现误报，可以通过下面的方式排出特定的 linter：

以 `staticcheck` 为例：

```yaml
linters-settings:
  staticcheck:
    checks:
      - all
      - '-SA1000' # disable the rule SA1000
      - '-SA1004' # disable the rule SA1004
```

#### 通过文本排除问题

下面的示例，所有在 exclude 中定义的文本的报告都会被排除：

```yaml
issues:
  exclude:
    - "Error return value of .((os\\.)?std(out|err)\\..*|.*Close|.*Flush|os\\.Remove(All)?|.*printf?|os\\.(Un)?Setenv). is not checked"
    - "exported (type|method|function) (.+) should have comment or be unexported"
    - "ST1000: at least one file in a package should have a package comment"
```

下面的示例，来自指定 `linters`，并且包含 `text` 指定的文本的报告会被排除：

```yaml
issues:
  exclude-rules:
    - linters:
        - gomnd
      text: "mnd: Magic number: 9"
```

下面的示例，来自指定 `linters`，并且来自指定的 `source` 的报告会被排除：

```yaml
issues:
  exclude-rules:
    - linters:
        - lll
      source: "^//go:generate "
```

下面的示例，`path` 指定的文件，并且包含 `text` 指定的文本的报告会被排除：

```yaml
issues:
  exclude-rules:
    - path: path/to/a/file.go
      text: "string `example` has (\\d+) occurrences, make it a constant"
```

#### 通过路径排除问题

在下面的示例中，所有匹配 `path-except` 指定路径的文件，并且来自指定 `linters` 的报告会被排除：

```yaml
issues:
  exclude-rules:
    - path: '(.+)_test\.go'
      linters:
        - funlen
        - goconst
```

排除特定路径以外的报告，下面的示例，只检查 test 文件：

```yaml
issues:
  exclude-rules:
    - path-except: '(.+)_test\.go'
      linters:
        - funlen
        - goconst
```

下面的示例，`skip-files` 相关的文件会被排除：

```yaml
run:
  skip-files:
    - path/to/a/file.go
```

下面的示例，`skip-dirs` 相关的目录会被排除：

```yaml
run:
  skip-dirs:
    - path/to/a/dir/
```

#### nolint 指令

使用 `//nolint:all` 可以排除所有问题，如果在行内使用（而不是从行首开始），则只排除这一行的问题。：

```go
var bad_name int //nolint:all
```

排除指定 linters 的问题：

```go
var bad_name int //nolint:golint,unused
```

在行首使用 `nolint`，可以排除整个代码块的问题：

```go
//nolint:all
func allIssuesInThisFunctionAreExcluded() *string {
  // ...
}

//nolint:govet
var (
  a int
  b int
)
```

排除整个文件的问题：

```go
//nolint:unparam
package pkg
```