---
title: GitHub Actions
weight: 9
---

# 基于 GitHub Actions 的 CI/CD

GitHub Actions 是 GitHub 为托管在 github.com 站点的项目提供的持续集成服务。

在构建持续集成任务时，需要完成很多操作，比如克隆代码、编译代码、运行单元测试、构建和发布镜像等。GitHub 把这些操作称为 Actions。

Actions 是可以共享的，开发者可以将 Actions 上传到 GitHub 的 [Actions 市场](https://github.com/marketplace?type=actions)。如果需要某个 Action，直接引用即可。
整个持续集成过程，就变成了一个 Actions 的组合。

Action 其实是一个独立的脚本，可以将 Action 存放在 GitHub 代码仓库中，通过 `<userName>/<repoName>` 的语法引用 Action。例如，`actions/checkout@v2` 表示 `https://github.com/actions/checkout` 这个仓库，`tag` 是 `v2`。

## GitHub Actions 术语

- `workflow`：一个 `.yml` 文件对应一个 workflow，也就是一次持续集成。一个 GitHub 仓库可以包含多个 workflow，只要是在 `.github/workflow` 目录下的
  `.yml` 文件都会被 GitHub 执行。
- `job`：一个 workflow 由一个或多个 `job` 构成，每个 `job` 代表一个持续集成任务。
- `step`：每个 `job` 由多个 `step` 构成，一步步完成。
- `action`：每个 `step` 可以依次执行一个或多个命令（action）。
- `on`：一个 workflow 的触发条件，决定了当前的 workflow 在什么时候被执行。

## workflow 

GitHub Actions 配置文件存放在代码仓库的 `.github/workflows` 目录下，文件后缀为 `.yml`、`.yaml`。GitHub 只要发现 `.github/workflows` 目录里面有 `.yml` 文件，就会自动运行该文件。

### 基础配置

- `name` 是 workflow 的名称。如果省略该字段，默认为当前 workflow 的文件名。
- `on` 指定触发 workflow 的条件。
  - `on: push`，意思是，`push` 事件触发 workflow。也可以是事件的数组，例如: `on: [push, pull_request]`。[更多触发事件](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows)。
  - `on.<push|pull_request>.<tags|branches>`，指定触发事件时，我们可以限定分支或标签。
    ```yaml
    # 只有 master 分支发生 push 事件时，才会触发 workflow。
    on:
      push:
        branches:
          - master
    ```
- `jobs.<job_id>.name` 表示要执行的一项或多项任务。`jobs` 字段里面，需要写出每一项任务的 `job_id`，具体名称自定义。`job_id` 里面的 `name` 字段是任务的说明。
  ```yaml
  # jobs 字段包含两项任务，job_id 分别是 my_first_job 和 my_second_job。
  jobs:
    my_first_job:
      name: My first job
    my_second_job:
      name: My second job
  ```
- `jobs.<job_id>.runs-on` `runs-on` 字段指定运行所需要的虚拟机环境，它是必填字段。可用的虚拟：
  - ubuntu-latest、ubuntu-18.04 或 ubuntu-16.04。
  - windows-latest、windows-2019 或 windows-2016。
  - macOS-latest 或 macOS-10.14。
- `jobs.<job_id>.steps` 指定每个 Job 的运行步骤，可以包含一个或多个步骤。每个步骤都可以指定下面三个字段。
  - `jobs.<job_id>.steps.name`：步骤名称。
  - `jobs.<job_id>.steps.run`：该步骤运行的命令或者 action。
  - `jobs.<job_id>.steps.env`：该步骤所需的环境变量。
  ```yaml
  name: Hello
  on: push
  jobs:
  my-job:
    name: My Job
    runs-on: ubuntu-latest
    steps:
    - name: Print a greeting
      env:
        GITHUB_TOKEN: {{ secrets.PAT }}
      run: |
        echo hello
  ```
- `jobs.<job_id>.uses` 可以引用别人已经创建的 actions。引用格式为 `username/repo@verison`，例如 `uses: actions/setup-go@v3`。
- `jobs.<job_id>.with` 设置 action 的参数。每个参数都是一个 `key/value`。
  ```yaml
  jobs:
    my_first_job:
    steps:
      - name: Set up Node
      - uses: actions/setup-node@v3
        with:
          node-version: '14'
  ```
- `jobs.<job_id>.run` 执行的命令。可以有多个命令，例如：
  ```yaml
  - name: Build
    run: |
      go mod tidy
      go build -v -o crtctl .
  ```

### 设置 job 的依赖关系

`needs` 字段可以指定当前任务的依赖关系，即运行顺序。

```yaml
jobs:
  job1:
  job2:
    needs: job1
  job3:
    needs: [job1, job2]
```

上面的示例，job1 必须先于 job2 成功完成，而 job3 等待 job1 和 job2 成功完成后才能运行。

不要求依赖的 job 是否成功：

```yaml
jobs:
  job1:
  job2:
    needs: job1
  job3:
    if: ${{ always() }}
    needs: [job1, job2]
```

上面的示例，job3 使用 `always()` 条件表达式，确保始终在 job1 和 job2 完成（无论是否成功）后运行。

### 使用构建矩阵

如果想在多个系统或者多个语言版本上测试构建，就需要设置构建矩阵。例如，在多个操作系统、多个 Go 版本下跑测试，可以使用如下 workflow 配置：

```yaml
name: Go Test
on: [push, pull_request]

jobs:
build:
    name: Test with go ${{ matrix.go_version }} on ${{ matrix.os }}
    runs-on: ${{ matrix.os }}
strategy:
  matrix:
    go_version: [1.15, 1.16]
    os: [ubuntu-latest, macOS-latest]
steps:
  - name: Set up Go ${{ matrix.go_version }}
    uses: actions/setup-go@v2
    with:
      go-version: ${{ matrix.go_version }}
    id: go
```

`strategy.matrix` 配置了该工作流程运行的环境矩阵，会在 4 台不同配置的服务器上执行该 workflow：`ubuntu-latest.1.15`、`ubuntu-latest.1.16`、
`macOS-latest.1.15`、`macOS-latest.1.16`。

### 使用 Secrets

在构建过程中，如果有用到 token 等敏感数据，此时就可以使用 secrets。我们在对应项目中选择 `Settings-> Secrets`，就可以创建 secret。

例如在 Secrets 中创建一个名为 `MySecrets` 的 secret，然后在 workflow 中引用：

```yaml
name: Go Test
on: [push, pull_request]
jobs:
  helloci-build:
    name: Test with go
    runs-on: [ubuntu-latest]
    environment:
      name: helloci
    steps:
      - name: use secrets
        env:
          super_secret: ${{ secrets.MySecrets }}
```

secret name 不区分大小写，所以如果新建 secret 的名字是 `name`，使用时用 `secrets.name` 或者 `secrets.Name` 都是可以的。

[更过 workflow 配置](https://docs.github.com/cn/actions/using-workflows/workflow-syntax-for-github-actions)。

## 常用 actions

### 静态代码检查

[golangci-lint-action](https://github.com/golangci/golangci-lint-action) 是 golangci-lint 官方提供的 action。

action 默认会读取项目根目录下的 `.golangci.yml` 配置文件。可以使用 `--config` 指定配置文件： `args: --config=/my/path/.golangci.yml`。

```yaml
name: golangci-lint
on:
  push:
    tags:
      - v*
    branches:
      - main
    paths-ignore:
      - 'docs/**'
      - 'README.md'
  pull_request:
    paths-ignore:
      - 'docs/**'
      - 'README.md'
permissions:
  contents: read

jobs:
  golangci:
    strategy:
      matrix:
        go: [ '1.20', '1.21' ]
        os: [ ubuntu-latest, windows-latest ]
    permissions:
      contents: read  # for actions/checkout to fetch code
      pull-requests: read  # for golangci/golangci-lint-action to fetch pull requests
    name: lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: stable # get the latest stable version from the go-versions repository manifest.
          cache: false
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          args: --timeout=10m
```

### 自动发布

[goreleaser-action](https://github.com/goreleaser/goreleaser-action) GoReleaser 官方提供和的 action。

action 默认读取项目根目录下的 `.goreleaser.yaml` 配置文件。可以使用 `--config` 指定配置文件： `args: --config=/my/path/.goreleaser.yml`。

```yaml
name: goreleaser

on:
  pull_request:
  push:

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Set up Go
        uses: actions/setup-go@v4
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          # either 'goreleaser' (default) or 'goreleaser-pro'
          distribution: goreleaser
          version: latest
          args: release --clean --rm-dist --debug
        env:
          GITHUB_TOKEN: ${{ secrets.PAT }}
          # Your GoReleaser Pro key, if you are using the 'goreleaser-pro' distribution
          # GORELEASER_KEY: ${{ secrets.GORELEASER_KEY }}
```

### 使用 Artifact 存储文件

在构建过程中，可能会输出一些构建产物，比如日志文件、测试结果等。可以使用 GitHub Actions Artifact 来存储。使用 [action/upload-artifact](https://github.com/actions/upload-artifact)
和 [download-artifact](https://github.com/actions/download-artifact) 进行构建参数的相关操作。

```yaml
steps:
  - run: npm ci
  - run: npm test
  - name: Upload Test Coverage File
    uses: actions/upload-artifact@v1.0.0
    with:
      name: coverage-output
      path: coverage
```

执行成功后，我们就能在对应 action 面板看到生成的 Artifact。

### 使用缓存加快 workflow

为了使 workflow 更快、更高效，可以为依赖项及其他经常重复使用的文件创建和使用缓存。例如：npm，go mod。要缓存 job 的依赖项可以使用 [cache](https://github.com/actions/cache) 。

cache 会根据 `key` 尝试还原缓存。当找到缓存时，会将缓存的文件还原到你配置的 `path`。

如果找到缓存，cache 会在 job 成功完成时会使用你提供的 `key` 自动创建一个新缓存。并包含 `path` 指定的文件。

可以选择提供在 `key` 与现有缓存不匹配时要使用的 `restore-keys` 列表。 从另一个分支还原缓存时，`restore-keys` 列表非常有用，因为 `restore-keys`
可以部分匹配缓存 `key`。

```yaml
# Look for a CLI that's made for this PR
- name: Fetch built CLI
  id: cli-cache
  uses: actions/cache@v2
  with:
    path: ./_output/linux/amd64/bin/crtctl
    # The cache key a combination of the current PR number and the commit SHA
    key: crtctl-${{ github.event.pull_request.number }}-${{ github.sha }}
```

#### 输入参数

- `key`：必须。缓存的 key。 它可以是变量、上下文值、静态字符串和函数的任何组合。 密钥最大长度为 512 个字符，密钥长度超过最大长度将导致操作失败。
- `path`：必须。运行器上用于缓存或还原的路径。可以指定单个路径，也可以在单独的行上添加多个路径。 例如：
  ```yaml
  - name: Cache Gradle packages
    uses: actions/cache@v3
    with:
      path: |
        ~/.gradle/caches
        ~/.gradle/wrapper
  ```
- `restore-keys`：可选的。备用的缓存 key 字符串，每个 key 放置在一个新行上。如果 key 没有命中缓存，则按照提供的顺序依次使用这些还原键来查找和还原缓存。例如：
  ```yaml
  restore-keys: |
    npm-feature-${{ hashFiles('package-lock.json') }}
    npm-feature-
    npm-
  ```

#### 输出参数

- `cache-hit`：布尔值，是否命中缓存。

```yaml
- if: ${{ steps.cache-npm.outputs.cache-hit != 'true' }}
  name: List the state of node modules
  continue-on-error: true
  run: npm list
```

#### 缓存匹配过程

1. 当 `key` 匹配现有缓存时，被称为缓存命中，并且操作会将缓存的文件还原到 `path` 目录。
2. 当 `key` 不匹配现有缓存时，则被称为缓存失误，在作业成功完成时会自动创建一个新缓存。发生缓存失误时，该操作还会搜索指定的 `restore-keys` 以查找任何匹配项： 
  - 如果提供 `restore-keys`，cache 操作将按顺序搜索与 `restore-keys` 列表匹配的任何缓存。 
    - 当存在精确匹配时，该操作会将缓存中的文件还原到 `path` 目录。
    - 如果没有精确匹配，操作将会搜索恢复键值的部分匹配。 当操作找到部分匹配时，最近的缓存将还原到 `path` 目录。
  - cache 操作完成，作业中的下一个步骤运行。
  - 如果作业成功完成，则操作将自动创建一个包含 `path` 目录内容的新缓存。

[匹配缓存键的详细过程](https://docs.github.com/cn/actions/using-workflows/caching-dependencies-to-speed-up-workflows#matching-a-cache-key) 。

#### 使用限制和收回政策

GitHub 将删除 7 天内未被访问的任何缓存条目。 可以存储的缓存数没有限制，但存储库中所有缓存的总大小限制为 10 GB。

如果超过此限制，GitHub 将保存新缓存，但会开始收回缓存，直到总大小小于存储库限制。

### 自动打 Label

使用 [actions/labeler](https://github.com/marketplace/actions/labeler) 来实现自动打 Label。

#### 使用

创建 `.github/labeler.yml` 文件，该文件包含标签列表和需要匹配的 [minimatch](https://github.com/isaacs/minimatch) globs，以应用标签。

`.github/labeler.yml` 文件中，key 就是 label 的名字，值是文件路径。

Workflow 示例：

```yaml
on:
  pull_request_target:
    types: [opened, reopened, synchronize, ready_for_review]

jobs:
  # Automatically labels PRs based on file globs in the change.
  triage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/labeler@v3
        with:
          repo-token: "${{ secrets.GITHUB_TOKEN }}"
          configuration-path: .github/labels.yml
```

输入参数：

- `repo-token`：GITHUB_TOKEN，需要 `contents:read` 和 `pull-requests:write` 权限。
- `configuration-path`：Label 配置文件路径。
- `sync-labels`：当匹配的文件被还原或不再被 PR 改变时，是否要删除标签。

### 在一个 PR 创建或打开时为自动 assign reviewer

使用 [auto-assign-action](https://github.com/marketplace/actions/auto-assign-action) 来实现自动 assign。

创建配置文件，例如：`.github/auto_assign.yml`。在文件中添加 reviewers/assignees。

```yaml
# Set to true to add reviewers to pull requests
addReviewers: true

# Set to true to add assignees to pull requests
addAssignees: false

# Set addAssignees to 'author' to set the PR creator as the assignee.
# addAssignees: author

# A list of reviewers to be added to pull requests (GitHub user name)
reviewers:
  - reviewerA
  - reviewerB
  - reviewerC

# A number of reviewers added to the pull request
# Set 0 to add all the reviewers (default: 0)
numberOfReviewers: 0
# A list of assignees, overrides reviewers if set
# assignees:
#   - assigneeA

# A number of assignees to add to the pull request
# Set to 0 to add all of the assignees.
# Uses numberOfReviewers if unset.
# numberOfAssignees: 2

# Set to true to add reviewers from different groups to pull requests
useReviewGroups: true

# A list of reviewers, split into different groups, to be added to pull requests (GitHub user name)
reviewGroups:
  groupA:
    - reviewerA
    - reviewerB
    - reviewerC
  groupB:
    - reviewerD
    - reviewerE
    - reviewerF

# Set to true to add assignees from different groups to pull requests
useAssigneeGroups: false
# A list of assignees, split into different froups, to be added to pull requests (GitHub user name)
# assigneeGroups:
#   groupA:
#     - assigneeA
#     - assigneeB
#     - assigneeC
#   groupB:
#     - assigneeD
#     - assigneeE
#     - assigneeF

# A list of keywords to be skipped the process that add reviewers if pull requests include it
# skipKeywords:
#   - wip

# The action will only run for non-draft PRs. If you want to run for all PRs, you need to enable it to run on drafts.
# runOnDraft: true
```


Workflow 示例：

```yaml
name: "Auto Assign Author"

# pull_request_target means that this will run on pull requests, but in
# the context of the base repo. This should mean PRs from forks are supported.
on:
  pull_request_target:
    types: [opened, reopened, ready_for_review]

jobs:
  # Automatically assigns reviewers and owner
  add-reviews:
    runs-on: ubuntu-latest
    steps:
      - name: Set the author of a PR as the assignee
        uses: kentaro-m/auto-assign-action@v1.2.4
        with:
          configuration-path: ".github/auto_assignees.yml"
          repo-token: "${{ secrets.GITHUB_TOKEN }}"
```

### 关闭不活跃的 Issue 和 PR

使用 [close-stale-issues](https://github.com/marketplace/actions/close-stale-issues) 来自动关闭长时间不活跃的 PR 和 issues。

配置必须在默认分支上，默认值将会：

- 在 60 天没有活跃的 issue 和 PR 上添加一个 "Stale" 标签，并添加 comments。
- 添加 "Stale" 标签 7 天后关闭 issue 和 PR。
- 如果 issue 和 PR 发生更新/评论，"Stale" 标签将被删除，计时器会重启。

需要的权限：

```yaml
permissions:
  contents: write # only for delete-branch option
  issues: write
  pull-requests: write
```

#### 示例

```yaml
name: "Close stale issues and PRs"
on:
  schedule:
    # First of every month
    - cron: "30 1 * * *"

jobs:
  stale:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/stale@v3
        with:
          repo-token: ${{ secrets.GITHUB_TOKEN }}
          stale-issue-message: "This issue is stale because it has been open 30 days with no activity. Remove stale label or comment or this will be closed in 5 days. If a Velero team member has requested log or more information, please provide the output of the shared commands."
          close-issue-message: "This issue was closed because it has been stalled for 5 days with no activity."
          days-before-issue-stale: 30
          days-before-issue-close: 5
          # Disable stale PRs for now; they can remain open.
          days-before-pr-stale: -1
          days-before-pr-close: -1
          # Only issues made after Oct 01 2022.
          start-date: "2022-10-01T00:00:00"
          # Only make issues stale if they have these labels. Comma separated.
          only-labels: "Needs info,Duplicate"
```

### 使用 Gitleaks 进行静态代码分析

[Gitleaks](https://github.com/marketplace/actions/gitleaks) 是一款 SAST 工具，用于检测和防止 git 仓库中的密码、API 密钥和令牌等硬编码秘密。

```yaml
name: gitleaks
on:
  pull_request:
  push:
  workflow_dispatch:
  schedule:
    - cron: "0 4 * * *" # run once a day at 4 AM
jobs:
  scan:
    name: gitleaks
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GITLEAKS_LICENSE: ${{ secrets.GITLEAKS_LICENSE}} # Only required for Organizations, not personal accounts.
```

常用配置：

```
# 定义了如何检测 secrets
[[rules]]
# 规则 id
id = "ignore-testdata"
# 为单条规则加入一个允许列表，以减少误报，或忽略已知的 secret 的提交。
[rules.allowlist]
paths = ['''.*/testdata/*''']
# 全局的允许列表
[allowlist]
```

更多配置 [Configuration](https://github.com/gitleaks/gitleaks#configuration)。

### 使用 Grype 扫描容器镜像和文件系统漏洞

[Grype](https://github.com/anchore/grype) 是一款针对容器镜像和文件系统的漏洞扫描程序。如果发现了漏洞，还可选择以可配置的严重程度失败。

```yaml
name: "grype"
on:
  push:
    branches: ['main']
    tags: ['v*']
  pull_request:
jobs:
  scan-source:
    name: scan-source
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      actions: read
      contents: read
    steps:
      - uses: actions/checkout@v3
      - uses: anchore/scan-action@v3
        with:
          path: "."
          fail-build: true
```

#### `grype` Configuration

默认配置文件的搜索路径:

- `.grype.yaml`
- `.grype/config.yaml`
- `~/.grype.yaml`
- `<XDG_CONFIG_HOME>/grype/config.yaml`

也可以使用  `--config`/`-c` 来指定配置文件。

常用配置：

```yaml
# 扫描时，如果发现严重性达到或超过设置的值，则返回代码为 1。默认为未设置
fail-on-severity: high
# 如果使用 SBOM 输入，则在软件包没有 CPE 时自动生成 CPE
add-cpes-if-none: true
# 输出格式 (允许的值: table, json, cyclonedx)
output: table
# 要从扫描中排除的文件
exclude:
  - "**/testdata/**"
# 如果看到 Grype 报告误报或任何其他不想看到的漏洞匹配，可以配置 "忽略规则"，Grype 会忽略匹配结果
ignore:
  - fix-state: unknown # 允许的值: "fixed", "not-fixed", "wont-fix", or "unknown"
    vulnerability: "CVE-2008-4318" # vulnerability ID
```

[更多 Grype 配置](https://github.com/anchore/grype)。

### 使用 CodeQL 进行安全性代码分析

GitHub CodeQL Action 是一个用于安全性代码分析的 GitHub Actions，使用 CodeQL 查询语言来搜索项目中的代码漏洞和安全问题。扫描完成后，CodeQL Action 会生成报告，扫描查询结果。

CodeQL 可以在 `Security -> Overview -> Code scanning alerts -> Set up code scanning` 找到官方给的 CodeQL Workflow Template。选择 `Set up this workflow` 就可以用 template 了。

也可以自己在 workflow 中加上 CodeQL Action：

```yaml
name: "codeql"
on:
  push:
    branches: [ main ]
jobs:
  analyze:
    name: analyze
    runs-on: ubuntu-latest
    permissions:
      actions: read
      contents: read
      security-events: write
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: stable
      - name: initialize codeql
        uses: github/codeql-action/init@v2
        with:
          languages: go # javascript, csharp, python, cpp, java
      - name: build package
        run: go build ./cmd
      # build package C/C++, C#, Java, Go, Swift  可以直接使用 CodeQL 的 autobuild 作替代
      # - name: auto build package  
      #   uses: github/codeql-action/autobuild@v2
      - uses: github/codeql-action/analyze@v2
```

### 自动提交 action 运行期间产生的文件

[git-auto-commit-action](https://github.com/stefanzweifel/git-auto-commit-action) 用于检测工作流运行期间更改的文件，并将其提交和推送回 GitHub 仓库。默认情况下，提交会以 "GitHub Action" 的名义进行，并由上次提交的用户共同撰写。

`CONTRIBUTING.md`，ChangeLog 之类的改动可以使用该 action 来实现自动提交。

```yaml
name: Format

on: push

jobs:
  format-code:
    runs-on: ubuntu-latest

    permissions:
      # Give the default GITHUB_TOKEN write permission to commit and push the
      # added or changed files to the repository.
      contents: write

    steps:
      - uses: actions/checkout@v3

      # Other steps that change files in the repository

      # Commit all changed files back to the repository
      - uses: stefanzweifel/git-auto-commit-action@v4
```

### 扫描 PR 中的依赖关系

[dependency-review-action](https://github.com/actions/dependency-review-action) 可以用来扫描 PR 中的依赖关系更改，如果引入了任何漏洞或无效许可证，则会引发错误。

```yaml
name: 'Dependency Review'
on: [pull_request]

permissions:
  contents: read

jobs:
  dependency-review:
    runs-on: ubuntu-latest
    steps:
      - name: 'Checkout Repository'
        uses: actions/checkout@v3
      - name: 'Dependency Review'
        uses: actions/dependency-review-action@v3
```

### 如何在 Action 中访问 GitHub

#### 使用 GitHub Access token

1. 首先需要生成一个 Access Token，[创建 token](https://github.com/settings/tokens/new)。
2. 在 repo 的 Settings 页面中添加 Secret，例如，我的 secret 命名为 PAT。

在 Action 中使用：

```
    steps:
      - name: release
        run: |
          GITHUB_TOKEN=${{ secrets.PAT }} make release
```

通过 Access Token 的方式 clone repo：

```
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          repository: shipengqi/crtctl
          token: ${{ secrets.PAT }}
          path: crtctl
```

上面的方式用的是 HTTPS 的方式。通过 `git remote -v` 查看可以看到 remote 的地址。

#### 使用 SSH

1. 首先需要一个 GitHub 中已经配置好的 ssh 的 public key。
2. 在 repo 的 Settings 页面中添加 Secret，例如，我的 secret 命名为 SSH_KEY。

在 Action 中配置 ssh：

```
- name: Install SSH Key
  uses: shimataro/ssh-key-action@v2
  with:
    key: ${{ secrets.SSH_KEY }} 
    known_hosts: 'just-a-placeholder-so-we-dont-get-errors'
```

之后就可以在 Action 的后续步骤中像在本地一样使用 SSH 的方式来 clone repo 和提交代码了。

