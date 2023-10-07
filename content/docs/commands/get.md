---
title: get
---

# get
```sh
usage: go get [-d] [-m] [-u] [-v] [-insecure] [build flags] [packages]

Get resolves and adds dependencies to the current development module
and then builds and installs them.

The first step is to resolve which dependencies to add.

For each named package or package pattern, get must decide which version of
the corresponding module to use. By default, get chooses the latest tagged
release version, such as v0.4.5 or v1.2.3. If there are no tagged release
versions, get chooses the latest tagged prerelease version, such as
v0.0.1-pre1. If there are no tagged versions at all, get chooses the latest
known commit.

This default version selection can be overridden by adding an @version
suffix to the package argument, as in 'go get golang.org/x/text@v0.3.0'.
For modules stored in source control repositories, the version suffix can
also be a commit hash, branch identifier, or other syntax known to the
source control system, as in 'go get golang.org/x/text@master'.
The version suffix @latest explicitly requests the default behavior
described above.

If a module under consideration is already a dependency of the current
development module, then get will update the required version.
Specifying a version earlier than the current required version is valid and
downgrades the dependency. The version suffix @none indicates that the
dependency should be removed entirely, downgrading or removing modules
depending on it as needed.

Although get defaults to using the latest version of the module containing
a named package, it does not use the latest version of that module's
dependencies. Instead it prefers to use the specific dependency versions
requested by that module. For example, if the latest A requires module
B v1.2.3, while B v1.2.4 and v1.3.1 are also available, then 'go get A'
will use the latest A but then use B v1.2.3, as requested by A. (If there
are competing requirements for a particular module, then 'go get' resolves
those requirements by taking the maximum requested version.)

The -u flag instructs get to update dependencies to use newer minor or
patch releases when available. Continuing the previous example,
'go get -u A' will use the latest A with B v1.3.1 (not B v1.2.3).

The -u=patch flag (not -u patch) instructs get to update dependencies
to use newer patch releases when available. Continuing the previous example,
'go get -u=patch A' will use the latest A with B v1.2.4 (not B v1.2.3).

In general, adding a new dependency may require upgrading
existing dependencies to keep a working build, and 'go get' does
this automatically. Similarly, downgrading one dependency may
require downgrading other dependencies, and 'go get' does
this automatically as well.

The -m flag instructs get to stop here, after resolving, upgrading,
and downgrading modules and updating go.mod. When using -m,
each specified package path must be a module path as well,
not the import path of a package below the module root.

The -insecure flag permits fetching from repositories and resolving
custom domains using insecure schemes such as HTTP. Use with caution.

The second step is to download (if needed), build, and install
the named packages.

If an argument names a module but not a package (because there is no
Go source code in the module's root directory), then the install step
is skipped for that argument, instead of causing a build failure.
For example 'go get golang.org/x/perf' succeeds even though there
is no code corresponding to that import path.

Note that package patterns are allowed and are expanded after resolving
the module versions. For example, 'go get golang.org/x/perf/cmd/...'
adds the latest golang.org/x/perf and then installs the commands in that
latest version.

The -d flag instructs get to download the source code needed to build
the named packages, including downloading necessary dependencies,
but not to build and install them.

With no package arguments, 'go get' applies to the main module,
and to the Go package in the current directory, if any. In particular,
'go get -u' and 'go get -u=patch' update all the dependencies of the
main module. With no package arguments and also without -u,
'go get' is not much more than 'go install', and 'go get -d' not much
more than 'go list'.

For more about modules, see 'go help modules'.

For more about specifying packages, see 'go help packages'.

This text describes the behavior of get using modules to manage source
code and dependencies. If instead the go command is running in GOPATH
mode, the details of get's flags and effects change, as does 'go help get'.
See 'go help modules' and 'go help gopath-get'.

See also: go build, go install, go clean, go mod.
```

使用`go get`命令下载一个包。如`go get github.com/golang/lint/golint`下载了`golint`包，`src`目录下会有`github.com/golang/lint/golint`包目录。
`bin`目录下可以看到`golint`可执行程序。

`go get`本质上可以理解为首先第一步是通过源码工具`clone`代码到`src`下面，然后执行`go install`。

**OPTIONS**
- `-u` 保证每个包是最新版本。