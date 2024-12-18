---
title: 项目管理
weight: 8
---

Go 项目通常使用 Makefile 作为项目管理工具。

通常 Go 项目的 Makefile 应该包括：格式化代码、静态代码检查、单元测试、代码构建、文件清理、帮助等功能。

学习 Makefile 的语法，推荐学习[《跟我一起写 Makefile》 (PDF 重制版)](https://github.com/seisman/how-to-write-makefile)。

## Makefile 结构

随着项目越来越大，需要管理的功能就会越来越多，如果全部放在一个 Makefile 中，会导致 Makefile 过大，难以维护，可读性差。所以设计 Makefile 结构时，最好采用**分层**的设计。

项目根目录下的 Makefile 来聚合子目录下的 Makefile 命令。将复杂的 shell 命令封装在 shell 脚本中，供 Makefile 直接调用，而一些简单的命令则可以直接集成在 Makefile 中。

![makefile](https://gitee.com/shipengqi/illustrations/raw/main/go/makefile.png)

示例；

```makefile
.PHONY: all
all: modules lint test build

# ==============================================================================
# Includes

include hack/include/common.mk # make sure include common.mk at the first include line
include hack/include/go.mk
include hack/include/release.mk

# ==============================================================================
# Targets

## build: build binary file.
.PHONY: build
build: modules
	@$(MAKE) go.build

## tag: generate release tag.
.PHONY: tag
tag:
	@$(MAKE) release.tag
	
## modules: add missing and remove unused modules.
.PHONY: modules
modules:
	@go mod tidy
	
## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: run unit test and get test coverage.
.PHONY: test
test:
	@$(MAKE) test.cover		
```

## Makefile 技巧

### 使用通配符和函数增强扩展性

```makefile
.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS))

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.golangci-lint
install.golangci-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: install.gsemver
install.gsemver:
	@go install github.com/arnaud-deprez/gsemver@latest

.PHONY: install.releaser
install.releaser:
	@go install github.com/goreleaser/goreleaser@latest

.PHONY: install.ginkgo
install.ginkgo:
	@go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

上面的示例 `tools.install.%` 和 `tools.verify.%` 都使用了通配符 `%`，在执行 `make tools.verify.ginkgo`, `make tools.verify.releaser` 这些命令时，都可以匹配到 `tools.verify.%` 这个规则。

如果不使用通配符，那么就要为这些 tools 分别去定义规则，例如 `tools.verify.ginkgo`。

上面的 `$*` 是自动变量，表示匹配到的值，例如 ginkgo、releaser。

`addprefix` 是一个函数，作用是给文件添加一个前缀. 

### 带层级的命名方式

使用带层级的命名方式，例如 `tools.verify.ginkgo`，实现**目标分组管理**。当 Makefile 有大量目标时，通过分组，可以更好地管理这些目标。可以通过组名识别出该目标的功能类别。还可以减小目标重名的概率。

### 定义环境变量

可以用一个特定的 Makefile 文件来定义环境变量，例如上面第一个实例中的 `common.mk`，然后在入口 Makefile 中第一个引入。

这些环境变量就对所有的 Makefile 文件生效，修改时也只要修改一处，避免重复工作。
