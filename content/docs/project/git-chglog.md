---
draft: true
---

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
