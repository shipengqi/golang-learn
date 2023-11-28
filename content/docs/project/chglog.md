---
draft: true
---

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