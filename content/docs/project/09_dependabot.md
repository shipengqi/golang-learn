---
title: GitHub Dependabot
weight: 9
---

# GitHub Dependabot

GitHub Dependabot 的配置文件 `dependabot.yml` 必须存放在代码仓库的 `.github` 目录下。在添加或更新 `dependabot.yml` 文件时，会立即触发版本更新检查。

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "daily"
      time: "08:00"
    labels:
      - "dependencies"
    commit-message:
      prefix: "feat"
      include: "scope"
```

上面的示例，`interval: "daily" time: "08:00"` 表示每天八点会触发版本更新检查。

`dependabot.yml` 文件中两个必须的字段：`version` 和 `updates`。该文件必须以 `version: 2` 开头。

## updates

`updates` 用来配置 Dependabot 如何更新版本或项目的依赖项，常用的选项：

| 选项                | required | 安全更新       | 版本更新       | 说明                   |
|-------------------|----------|------------|------------|----------------------|
| package-ecosystem | yes      | no         | yes        | 要使用的包管理器             |
| directory         | yes      | yes        | yes        | package manifests 位置 |
| schedule.interval | yes      | no         | yes        | 检查更新的频率              |
| allow             | no       | yes        | yes        | 自定义哪些允许更新            |
| assignees         | no       | yes        | yes        | assign PR            |
| labels            | no       | yes        | yes        | 设置 PR 的 label        |
| commit-message         | no       | yes        | yes        | 提交 mesaage 的选项       |
| groups         | no       | no         | yes        | 对某些依赖项的更新分组         |


更多配置：

[dependabot.yml 文件的配置选项](https://docs.github.com/zh/code-security/dependabot/dependabot-version-updates/configuration-options-for-the-dependabot.yml-file)。