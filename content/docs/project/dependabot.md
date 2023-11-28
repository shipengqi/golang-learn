---
draft: true
---

# Dependabot 配置

Dependabot 配置文件 dependabot.yml 使用 YAML 语法。必须将此文件存储在存储库的 `.github` 目录中。 在添加或更新 `dependabot.yml` 文件时，这将立即触发版本更新检查。

`dependabot.yml` 两个必须的字段：`version` 和 `updates`。该文件必须以 `version: 2` 开头。

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

[dependabot.yml 文件的配置选项](https://docs.github.com/zh/code-security/dependabot/dependabot-version-updates/configuration-options-for-the-dependabot.yml-file)