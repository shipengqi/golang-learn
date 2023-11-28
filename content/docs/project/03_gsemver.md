---
title: 版本规范
weight: 3
---

# 版本规范

[gsemver](https://github.com/arnaud-deprez/gsemver) 是一个用 Go（Golang）开发的命令行工具，它使用 git commit 来自动生成符合 semver 2.0.0 规范的下一个版本。

## 安装

```
$ go install github.com/arnaud-deprez/gsemver@latest
```

## 使用

下面的命令会使用 git commit 生成下一个 version：
```
gsemver bump
```

## 配置

你可以使用一个配置文件来定义你自己的规则。默认情况下，会寻找 `.gsemver.yaml` 或 `$HOME/.gsemver.yaml`，可以通过 `--config`（或 `-c`）选项来指定你自己的配置文件。
