---
title: Commit 规范
weight: 2
---

# Commit 规范


# chglog

`chglog` 是一个 CHANGELOG 生成器。

## 安装

```
$ go get github.com/goreleaser/chglog/cmd/chglog@latest
```

## 使用

第一步，初始化一个配置文件 `.chglog.yml`

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

现在，每当你要发布另一个版本时，只需执行 `chglog add --version v#.#.#`（版本必须是 `semver` 格式）。

就是这样！

# git-chglog

`git-chglog` 是一个 CHANGELOG 生成器。

## 安装

```
$ go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
```

## 创建配置文件

```
$ git-chglog --init
```

选项：

- What is the URL of your repository?: https://github.com/shipengqi/crtctl
- What is your favorite style?: github
- Choose the format of your favorite commit message: <type>: <subject> -- feat: Add new feature
- What is your favorite template style?: standard
- Do you include Merge Commit in CHANGELOG?: n
- Do you include Revert Commit in CHANGELOG?: y
- In which directory do you output configuration files and templates?: .chglog

`git-chglog` 的配置文件是一个 yaml 文件，默认路径为 `.chglog/config.yml`。[更多配置](https://github.com/git-chglog/git-chglog#configuration)。


## 使用

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