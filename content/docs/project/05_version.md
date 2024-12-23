---
title: 版本规范
weight: 5
---

Go 官方推荐的版本规范是 [semver](https://semver.org/lang/zh-CN/)（Semantic Versioning），也就是语义化版本。这个规范是 GitHub 起草的一个具有指导意义的、统一的版本号表示规范。

semver 是一种清晰可读的，明确反应版本信息的版本格式：

```
主版本号.次版本号.修订号
```

- 主版本号：做了不兼容的 API 修改。
- 次版本号：向下兼容的新增功能以及修改。
- 修订号： 向下兼容的问题修复。

例如 `v1.2.3`。

semver 还有先行版本号和编译版本号，格式为 `X.Y.Z[-先行版本号][+编译版本号]`。

例如 `v1.2.3-alpha.1+001`，`alpha.1` 就是先行版本号，`001` 是编译版本号。

- 先行版本号，意味着该版本不稳定，可能存在兼容性问题，可以用 `.` 作为分隔符。
- 编译版本号，一般是编译器在编译过程中自动生成的。

> 先行版本号和编译版本号只能是字母、数字，并且不可以有空格。

## 如何确定版本号？

1. 在实际开发的时候，可以使用 `0.1.0` 作为第一个开发版本号，并在后续的每次发行时递增次版本号。
2. 当软件是一个稳定的版本，并且第一次对外发布时，版本号应该是 `1.0.0`。
3. 严格按照 Angular 规范提交代码，版本号可以按照下面的规则来确定：
  - fix 类型的 commit 可以将修订号 `+1`。
  - feat 类型的 commit 可以将次版本号 `+1`。
  - 带有 BREAKING CHANGE 的 commit 可以将主版本号 `+1`。

## 如何处理将要弃用的功能?

弃用已存在的功能，在软件开发中是常规操作，如果要弃用某个功能，要做到两点：

1. 更新用户文档，通知用户。
2. 发布新的次版本，要包含舍弃的功能，直到发布新的主版本，目的是让用户能够平滑的迁移到新的 API。

## 自动生成语义化版本

[gsemver](https://github.com/arnaud-deprez/gsemver) 是一个用 Go 实现的命令行工具，它使用 git commit 来自动生成符合 semver 2.0.0 规范的下一个版本。

### 安装

```
$ go install github.com/arnaud-deprez/gsemver@latest
```

### 使用

下面的命令会根据 git commit 生成下一个 version：

```
gsemver bump
```

### 配置

可以使用配置文件来定义版本的生成规则，一般这个配置文件会放在项目的根目录下。

默认情况下，gsemver 会寻找 `.gsemver.yaml` 或 `$HOME/.gsemver.yaml` 文件，也可以通过命令行参数 `--config`（或 `-c`）选项来指定配置文件。
