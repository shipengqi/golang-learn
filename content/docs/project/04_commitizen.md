---
title: Commit 规范
weight: 4
---

# Commit 规范

多人协作开发一个项目时，每个开发的 Commit Message 五花八门，时间久了，提交的历史变得很难看，而且有的 Commit Message 可能过于简单，可读性较差。

一个好的 Commit 规范可以使 Commit Message 的可读性更好，并且可以实现自动化。

一个好的 Commit Message 应该满足以下要求：

- 清晰地描述 commit 的变更内容。
- 可以基于这些 Commit Message 进行过滤查找，比如只查找某个版本新增的功能：`git log --oneline --grep "^feat|^fix"`。
- 可以基于规范化的 Commit Message 生成 Change Log。
- 可以依据某些类型的 Commit Message 触发构建或者发布流程，比如当类型为 feat、fix 时触发 CI 流程。
- 确定语义化版本的版本号。比如 fix 类型可以映射为 PATCH 版本，feat 类型可以映射为 MINOR 版本。带有 BREAKING CHANGE 的 commit，可以映射为 MAJOR 版本。

目前，开源社区有多种 Commit 规范，例如 jQuery、Angular 等。[Angular 规范](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit?pli=1#heading=h.uyo6cb12dt6w)是使用最广泛的，格式清晰易读。

## Angular 规范

Angular 规范中，Commit Message 包含三个部分：Header、Body 和 Footer。格式如下：

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

### Header

Header 包括提价类型（type，必需的）、作用域（scope，可选的）和主题（subject）。Header 是必需的。

**type**包括：

- `feat`：增加了新功能
- `fix`：修复问题
- `pref`：优化性能
- `test`：测试代码修改
- `refactor`：代码重构
- `style`：不影响代码含义的修改，比如空格、格式化、缺失的分号等
- `docs`：对文档进行了修改
- `build`：对构建系统或者外部依赖项进行了修改
- `ci`：对 CI/CD 配置文件或脚本进行了修改
- `chore`：其他类型

**scope**：

用来说明 commit 的影响范围的，不同项目会有不同的 scope，项目初期，可以设置一些粒度比较大的 scope，比如可以按组件名或者功能来设置 scope。后续，如果项目有变动或者有新功能，
可以再用追加的方式添加新的 scope。

scope 不适合设置太具体的值。太具体的话，一方面会导致项目有太多的 scope，难以维护。另一方面，开发者也难以确定 commit 属于哪个具体的 scope，导致错放 scope。

**subject**：

对本次 commit 的简短描述。

- 必须以动词开头、使用现在时。
- 第一个字母必须是小写。
- 末尾不能添加句号。

### Body

对本次 commit 的更详细的描述。Body 的要求和 Header 的 subject 是一样的。Body 是可选的。

- 应该包含本次 commit 的动机以及和之前行为的对比。

### Footer

Footer 也是可选的。通常用来说明不兼容的改动和关闭的 Issue 列表。格式如下：

```
BREAKING CHANGE: <breaking change summary>
<BLANK LINE>
<breaking change description + migration instructions>
<BLANK LINE>
<BLANK LINE>
Closes #<issue number>
```

示例：

```
BREAKING CHANGE: isolate scope bindings definition has changed and
    the inject option for the directive controller injection was removed.
    
    To migrate the code follow the example below:
    
    Before:
    
    scope: {
      myAttr: 'attribute',
      myBind: 'bind',
      myExpression: 'expression',
      myEval: 'evaluate',
      myAccessor: 'accessor'
    }
    
    After:
    
    scope: {
      myAttr: '@',
      myBind: '@',
      myExpression: '&',
      // myEval - usually not useful, but in cases where the expression is assignable, you can use '='
      myAccessor: '=' // in directive's template change myAccessor() to myAccessor
    }
    
    The removed `inject` wasn't generaly useful for directives so there should be no code using it.
    
    
Closes #123, #245, #992    
```

## 自动生成规范化的 Commit Message

可以使用一些开源的工具，来自动化地生成规范化的 Commit Message：

- [commitizen](https://github.com/commitizen/cz-cli)，Javascript 实现，需要安装 Node.js。
- [commitizen-go](https://github.com/lintingzhen/commitizen-go)，Go 版本的简化版的 commitizen，下载二进制文件就可以直接使用。

上面两个命令都可以进入交互模式，并根据提示生成 Commit Message，然后提交。

### commitizen

#### 安装

```
$ npm install -g commitizen
```

## 自动生成 CHANGELOG

### goreleaser/chglog

[chglog](https://github.com/goreleaser/chglog) 是 goreleaser 开源的一个 CHANGELOG 生成器。

#### 安装

```
$ go get github.com/goreleaser/chglog/cmd/chglog@latest
```

#### 使用

第一步，初始化一个配置文件 `.chglog.yml`，一般放在项目的根目录下：

```
$ chglog config
```

根据需要修改配置文件：

```yaml
conventional-commits: false
deb:
  distribution: []
  urgency: ""
debug: false
owner: ""
package-name: ""
```

下一步，执行 `chglog init`：

```yaml
- semver: 0.0.1
  date: 2019-10-18T16:05:33-07:00
  packager: dj gilcrease <example@example.com>
  changes:
    - commit: 2c499787328348f09ae1e8f03757c6483b9a938a
      note: |-
        oops i forgot to use Conventional Commits style message

        This should NOT break anything even if I am asking to build the changelog using Conventional Commits style message
    - commit: 3ec1e9a60d07cc060cee727c97ffc8aac5713943
      note: |-
        feat: added file two feature

        BREAKING CHANGE: this is a backwards incompatible change
    - commit: 2cc00abc77d401a541d18c26e5c7fbef1effd3ed
      note: |-
        feat: added the fileone feature

        * This is a test repo
        * so ya!
```

然后执行 `chglog format --template repo > CHANGELOG.md` 来生成 `CHANGELOG.md` 文件。

现在，每当要发布另一个版本时，只需执行 `chglog add --version v#.#.#`（版本必须是 `semver` 格式）。


### git-chglog

[git-chglog](https://github.com/git-chglog/git-chglog) 也是一个 CHANGELOG 生成器。

#### 安装

```
$ go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
```

#### 创建配置文件

```
$ git-chglog --init
```

选项：

- What is the URL of your repository?: <repo address>
- What is your favorite style?: github
- Choose the format of your favorite commit message: <type>: <subject> -- feat: Add new feature
- What is your favorite template style?: standard
- Do you include Merge Commit in CHANGELOG?: n
- Do you include Revert Commit in CHANGELOG?: y
- In which directory do you output configuration files and templates?: .chglog

git-chglog 的配置文件是一个 yaml 文件，默认路径为 `.chglog/config.yml`。[更多配置](https://github.com/git-chglog/git-chglog#configuration)。


#### 使用

使用 `-o`（`--output`）输出 changelog 文件：

```bash
$ git-chglog -o CHANGELOG/CHANGELOG-v0.1.0.md
```

```bash
$ git-chglog

  If <tag query> is not specified, it corresponds to all tags.
  This is the simplest example.

$ git-chglog 1.0.0..2.0.0

  The above is a command to generate CHANGELOG including commit of 1.0.0 to 2.0.0.

$ git-chglog 1.0.0

  The above is a command to generate CHANGELOG including commit of only 1.0.0.

$ git-chglog $(git describe --tags $(git rev-list --tags --max-count=1))

  The above is a command to generate CHANGELOG with the commit included in the latest tag.

$ git-chglog --output CHANGELOG.md

  The above is a command to output to CHANGELOG.md instead of standard output.

$ git-chglog --config custom/dir/config.yml

  The above is a command that uses a configuration file placed other than ".chglog/config.yml".
```