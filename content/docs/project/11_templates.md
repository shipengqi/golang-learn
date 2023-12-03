---
title: GitHub 模板
weight: 11
---

# GitHub 模板

通过 Issue 和 PR 模板可以自定义和标准化贡献者创建 issue 和 PR 的信息。

## PR Template

PR 模板可以在任意的目录下，如果有多个 PR 模板，需要创建一个 `PULL_REQUEST_TEMPLATE` 目录。

例如，可以在 repo 的根目录下创建 `pull_request_template.md`，也可以在放在 `.github` 目录中 `.github/pull_request_template.md`。 

`pull_request_template.md`:

```
Thank you for contributing to crtctl!

# Please add a summary of your change

# Does your change fix a particular issue?

Fixes #(issue)
```

## Issue Template

Issue 模板存储在 repo 的 `.github/ISSUE_TEMPLATE` 目录中。文件名不区分大小写，扩展名为 `.md`。

`ISSUE_TEMPLATE/bug_report.md`

```markdown
---
name: Bug Report
about: Tell us about a problem you are experiencing
---

**What steps did you take and what happened:**
[A clear and concise description of what the bug is, and what commands you ran.)


**What did you expect to happen:**

**The following information will help us better understand what's going on**:

- Please provide the output and log content of your commands (Pasting long output into a [GitHub gist](https://gist.github.com) or other pastebin is fine.)

**Anything else you would like to add:**
[Miscellaneous information that will assist in solving the issue.]


**Environment:**

- crtctl version (use `crtctl -v`):
- Kubernetes version (use `kubectl version`):
- OS (e.g. from `/etc/os-release`):
```

`ISSUE_TEMPLATE/feature-request.md`:

```markdown
---
name: Feature Request
about: Suggest an idea for this project

---

**Describe the problem/challenge you have**
[A description of the current limitation/problem/challenge that you are experiencing.]


**Describe the solution you'd like**
[A clear and concise description of what you want to happen.]


**Anything else you would like to add:**
[Miscellaneous information that will assist in solving the issue.]


**Environment:**

- crtctl version (use `crtctl -v`):
- Kubernetes version (use `kubectl version`):
- OS (e.g. from `/etc/os-release`):
```

`.github/ISSUE_TEMPLATE` 目录下可以添加配置文件 `config.yml`，这个文件用来定义用户在 repo 中创建 issue 时可以看到哪些 issue 模板。

`ISSUE_TEMPLATE/config.yml`:

```yaml
# blank_issues_enabled 设置为 true，则用户可以选择打开空白 issue，不适用 issue 模板。
blank_issues_enabled: false
# contact_links 将用户引导到外部网站
contact_links:
  - name: GitHub Community Support
    url: https://github.com/orgs/community/discussions
    about: Please ask and answer questions here.
```

更多配置可以查看[官方文档](https://docs.github.com/cn/communities/using-templates-to-encourage-useful-issues-and-pull-requests/about-issue-and-pull-request-templates)。

## Issue 表单

Issue 表单比 Issue 模板的功能更加丰富，可以定义不同的输入类型、验证、默认标签等。Issue 表单使用 `yaml` 文件定义，也是存放在 `.github/ISSUE_TEMPLATE` 目录下。

示例：

```yaml
name: Bug Report
description: Tell us about a problem you are experiencing
labels: [bug, triage]
assignees:
  - shipengqi
body:
  - type: markdown
    attributes:
      value: |
        Thanks for taking the time to fill out this bug report! Please fill the form below.
  - type: textarea
    id: what-happened
    attributes:
      label: What happened?
      description: Also tell us, what did you expect to happen?
    validations:
      required: true
  - type: textarea
    id: reproducible
    attributes:
      label: How can we reproduce this?
      description: Please share a public repository that reproduces the issue, or an example config file.
    validations:
      required: true
  - type: textarea
    id: jaguar-version
    attributes:
      label: jaguar version
      description: "`jaguar --version` output"
      render: bash
    validations:
      required: true
  - type: textarea
    id: os
    attributes:
      label: OS
      description: "e.g. from `/etc/os-release`"
      render: bash
    validations:
      required: true
  - type: checkboxes
    id: search
    attributes:
      label: Search
      options:
        - label: I did search for other open and closed issues before opening this
          required: true
  - type: textarea
    id: ctx
    attributes:
      label: Additional context
      description: Anything else you would like to add
    validations:
      required: false
```

更多配置可以查看[官方文档](https://docs.github.com/zh/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms)。

## 添加安全策略

在 repo 的根目录、`docs` 或 `.github` 文件夹中可以添加一个 `SECURITY.md` 文件。当用户在 repo 中创建一个 issue 时，可以会看到一个指向你项目安全策略的链接。

更多配置可以查看[官方文档](https://docs.github.com/cn/code-security/getting-started/adding-a-security-policy-to-your-repository)。