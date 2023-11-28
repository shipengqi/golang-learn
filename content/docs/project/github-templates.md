---
draft: true
---

# 模板

[Issue 和 PR 模板的官方文档](https://docs.github.com/cn/communities/using-templates-to-encourage-useful-issues-and-pull-requests/about-issue-and-pull-request-templates) 。

通过 Issue 和 PR 模板可以自定义和标准化贡献者创建 issue 和 PR 的信息。

## PR Template

PR 模板可以在任意的目录下，如果要包含多个 PR 模板，需要创建子目录 `PULL_REQUEST_TEMPLATE`。

例如可以在 repo 的 root 目录下创建 `pull_request_template.md`，也可以在隐藏目录中 `.github/pull_request_template.md`。 

`pull_request_template.md`:

```
Thank you for contributing to crtctl!

# Please add a summary of your change

# Does your change fix a particular issue?

Fixes #(issue)
```

## Issue Template

当 repo 中创建了 issue 模板以后，贡献者在打开一个 issue 时可以选择一个合适的模板。issue 模板存储在 repo 的默认分支中的 `.github/ISSUE_TEMPLATE` 目录中。
文件名不区分大小写，且需要 `.md` 扩展名。

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

可通过在 `.github/ISSUE_TEMPLATE` 目录下添加 `config.yml` 文件来自定义用户在 repo 中创建 issue 是看到的 issue 模板。

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

## 添加安全策略

[添加安全策略官方文档](https://docs.github.com/cn/code-security/getting-started/adding-a-security-policy-to-your-repository) 。

为了给人们报告你项目中的安全漏洞的指示，你可以在你的 repo 的根目录、`docs` 或 `.github` 文件夹中添加一个 `SECURITY.md` 文件。当有人在
你的 repo 中创建一个 issue 时，他们会看到一个指向你项目安全策略的链接。
