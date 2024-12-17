---
title: GoReleaser
weight: 12
---

GoReleaser 是用 Go 编写项目的自动发布工具，支持交叉编译，并且支持发布到 Github，Gitlab 和 Gitea。

## 安装

```
go install github.com/goreleaser/goreleaser@latest
```

[更多安装方式](https://goreleaser.com/install/)。

## 使用

生成配置文件 `.goreleaser.yaml`，一般这个文件放在项目的根目录下：

```
goreleaser init
```

下面的命令可以发布一个 "仅限本地" 的 release，一般用来测试 `release` 命令是否可以正常运行。

```
goreleaser release --snapshot --rm-dist
```

修改 `.goreleaser.yaml` 配置后，可以用 `check` 命令检查配置：

```
goreleaser check
```

`--single-target` 只为特定的 `GOOS/GOARCH` 构建二进制文件，这对本地开发很有用：

```
goreleaser build --single-target
```

### 发布一个 release

如果要发布到 Github，需要导出一个环境变量 `GITHUB_TOKEN`，它应该包含一个有效的 GitHub token 与 repo 范围。它将被用来部署发布到你的 GitHub 仓库。[创建一个新的 GitHub 令牌](https://github.com/settings/tokens/new)。

> `write:packages` 权限是 `GITHUB_TOKEN` 需要的最小权限。

GoReleaser 会使用 repo 的最新 [Git 标签](https://git-scm.com/book/en/v2/Git-Basics-Tagging)。

首先需要创建一个 tag 并 push 到 Github：

```yaml
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

> 注意 tag 必须是一个有效的 [semantic version](https://semver.org/)

然后运行：`goreleaser release`。

如果暂时不想创建 tag，可以运行 `goreleaser release --snapshot`，这个命令会不会发布到 Github。

### Dry run

如果想在进行发布之前测试一下，可以通过下面的方式。

#### Build-only Mode

构建项目代码，可以用来验证项目的构建对所有构建目标有没有错误。

```
goreleaser build
```

#### Release Flags

`--skip-publish` 参数可以跳过发布：

```
goreleaser release --skip-publish
```

## build 配置

```yaml
# .goreleaser.yaml
builds:
  # You can have multiple builds defined as a yaml list
  -
    # ID of the build.
    # Defaults to the binary name.
    id: "my-build"

    # Path to main.go file or main package.
    # Notice: when used with `gomod.proxy`, this must be a package.
    #
    # Default is `.`.
    main: ./cmd/my-app

    # Binary name.
    # Can be a path (e.g. `bin/app`) to wrap the binary in a directory.
    # Default is the name of the project directory.
    binary: program

    # Custom flags templates.
    # Default is empty.
    flags:
      - -tags=dev
      - -v

    # Custom asmflags templates.
    # Default is empty.
    asmflags:
      - -D mysymbol
      - all=-trimpath={{.Env.GOPATH}}

    # Custom gcflags templates.
    # Default is empty.
    gcflags:
      - all=-trimpath={{.Env.GOPATH}}
      - ./dontoptimizeme=-N

    # Custom ldflags templates.
    # Default is `-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser`.
    ldflags:
      - -s -w -X main.build={{.Version}}
      - ./usemsan=-msan

    # Custom build tags templates.
    # Default is empty.
    tags:
      - osusergo
      - netgo
      - static_build
      - feature

    # Custom environment variables to be set during the builds.
    #
    # Default: `os.Environ()` merged with what you set the root `env` section.
    env:
      - CGO_ENABLED=0

    # GOOS list to build for.
    # For more info refer to: https://golang.org/doc/install/source#environment
    # Defaults are darwin and linux.
    goos:
      - freebsd
      - windows

    # GOARCH to build for.
    # For more info refer to: https://golang.org/doc/install/source#environment
    # Defaults are 386, amd64 and arm64.
    goarch:
      - amd64
      - arm
      - arm64

    # GOARM to build for when GOARCH is arm.
    # For more info refer to: https://golang.org/doc/install/source#environment
    # Default is only 6.
    goarm:
      - 6
      - 7

    # GOAMD64 to build when GOARCH is amd64.
    # For more info refer to: https://golang.org/doc/install/source#environment
    # Default is only v1.
    goamd64:
      - v2
      - v3

    # GOMIPS and GOMIPS64 to build when GOARCH is mips, mips64, mipsle or mips64le.
    # For more info refer to: https://golang.org/doc/install/source#environment
    # Default is only hardfloat.
    gomips:
      - hardfloat
      - softfloat

    # List of combinations of GOOS + GOARCH + GOARM to ignore.
    # Default is empty.
    ignore:
      - goos: darwin
        goarch: 386
      - goos: linux
        goarch: arm
        goarm: 7
      - goarm: mips64
      - gomips: hardfloat
      - goamd64: v4

    # Optionally override the matrix generation and specify only the final list
    # of targets.
    #
    # Format is `{goos}_{goarch}` with their respective suffixes when
    # applicable: `_{goarm}`, `_{goamd64}`, `_{gomips}`.
    #
    # Special values:
    # - go_118_first_class: evaluates to the first-class targets of go1.18.
    #   Since GoReleaser v1.9.
    # - go_first_class: evaluates to latest stable go first-class targets,
    #   currently same as 1.18. Since GoReleaser v1.9.
    #
    # This overrides `goos`, `goarch`, `goarm`, `gomips`, `goamd64` and
    # `ignores`.
    targets:
      - go_first_class
      - go_118_first_class
      - linux_amd64_v1
      - darwin_arm64
      - linux_arm_6

    # Set a specific go binary to use when building.
    # It is safe to ignore this option in most cases.
    #
    # Default is "go"
    gobinary: "go1.13.4"

    # Sets the command to run to build.
    # Can be useful if you want to build tests, for example,
    # in which case you can set this to "test".
    # It is safe to ignore this option in most cases.
    #
    # Default: build.
    # Since: v1.9.
    command: test

    # Set the modified timestamp on the output binary, typically
    # you would do this to ensure a build was reproducible. Pass
    # empty string to skip modifying the output.
    # Default is empty string.
    mod_timestamp: '{{ .CommitTimestamp }}'

    # Hooks can be used to customize the final binary,
    # for example, to run generators.
    # Those fields allow templates.
    # Default is both hooks empty.
    hooks:
      pre: rice embed-go
      post: ./script.sh {{ .Path }}

    # If true, skip the build.
    # Useful for library projects.
    # Default is false
    skip: false

    # By default, GoReleaser will create your binaries inside
    # `dist/${BuildID}_${BuildTarget}`, which is an unique directory per build
    # target in the matrix.
    # You can set subdirs within that folder using the `binary` property.
    #
    # However, if for some reason you don't want that unique directory to be
    # created, you can set this property.
    # If you do, you are responsible for keeping different builds from
    # overriding each other.
    #
    # Defaults to `false`.
    no_unique_dist_dir: true

    # By default, GoReleaser will check if the main filepath has a main
    # function.
    # This can be used to skip that check, in case you're building tests, for
    # example.
    #
    # Default: false.
    # Since: v1.9.
    no_main_check: true

    # Path to project's (sub)directory containing Go code.
    # This is the working directory for the Go build command(s).
    # If dir does not contain a `go.mod` file, and you are using `gomod.proxy`,
    # produced binaries will be invalid.
    # You would likely want to use `main` instead of this.
    # Default is `.`.
    dir: go

    # Builder allows you to use a different build implementation.
    # This is a GoReleaser Pro feature.
    # Valid options are: `go` and `prebuilt`.
    # Defaults to `go`.
    builder: prebuilt

    # Overrides allows to override some fields for specific targets.
    # This can be specially useful when using CGO.
    # Note: it'll only match if the full target matches.
    #
    # Default: empty.
    # Since: v1.5.
    overrides:
      - goos: darwin
        goarch: arm64
        goamd64: v1
        goarm: ''
        gomips: ''
        ldflags:
          - foo
        tags:
          - bar
        asmflags:
          - foobar
        gcflags:
          - foobaz
        env:
          - CGO_ENABLED=1
```

## 设置自定义 tag

可以使用环境变量强制 build tag 和 previous tag。这在一个 git 提交被多个 git tag 引用的情况下很有用。

```
export GORELEASER_CURRENT_TAG=v1.2.3
export GORELEASER_PREVIOUS_TAG=v1.1.0
goreleaser release
```

## changelog 配置

```yaml
changelog:
  sort: asc
  use: github
  filters:
    exclude:
      - '^test:'
      - '^chore'
      - 'merge conflict'
      - Merge pull request
      - Merge remote-tracking branch
      - Merge branch
      - go mod tidy
  groups:
    - title: 'Dependency Updates'
      regexp: "^.*(feat|fix)\\(deps\\)*:+.*$"
      order: 300
    - title: 'New Features'
      regexp: "^.*feat[(\\w)]*:+.*$"
      order: 100
    - title: 'Bug Fixes'
      regexp: "^.*fix[(\\w)]*:+.*$"
      order: 200
    - title: 'Documentation Updates'
      regexp: "^.*docs[(\\w)]*:+.*$"
      order: 400
    - title: Other work
      order: 9999
```

- `exclude` 下匹配到的文本不会被添加到 CHANGELOG 中。
- `groups` 根据 Commit Message 的 type 分组。 

生成的 CHANGELOG 如下图：

![changlog](https://raw.githubusercontent.com/shipengqi/illustrations/f30c629498b577f1cb6b84a60e98fca4ad1da984/go/changlog.png)
