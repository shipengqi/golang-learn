---
title: Go 命令
---

# Go 命令

```bash
$ go
Go is a tool for managing Go source code.

Usage:

        go command [arguments]

The commands are:

        build       编译指定的源码包以及它们的依赖包
        clean       删除掉执行其它命令时产生的一些文件和目录
        doc         show documentation for package or symbol
        env         打印 Go 的环境信息
        bug         start a bug report
        fix         把指定代码包的所有 Go 语言源码文件中的旧版本代码修正为新版本的代码
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         下载或更新指定的代码包及其依赖包，并对它们进行编译和安装
        install     编译并安装指定的源码包以及它们的依赖包
        list        列出指定的代码包的信息
        mod         Go 的依赖包管理工具
        run         编译并运行 Go 程序
        test        对指定包进行测试
        tool        运行指定的 go 工具
        version     打印 Go 的版本信息
        vet         检查 Go 语言源码中静态错误的工具，报告包中可能出现的错误

Use "go help [command]" for more information about a command.

Additional help topics:

        c           calling between Go and C
        buildmode   build modes
        cache       build and test caching
        filetype    file types
        gopath      GOPATH environment variable
        environment environment variables
        importpath  import path syntax
        packages    package lists
        testflag    testing flags
        testfunc    testing functions

Use "go help [topic]" for more information about that topic.
```

## TODO
- `go get` 和 `go install` 的区别,update blog for makefile
- 添加 `golint` 使用
- [go 命令](https://github.com/hyper0x/go_command_tutorial/blob/master/SUMMARY.md)
- [golint](https://github.com/golang/lint/tree/master/testdata)
  - [uber style guide](https://github.com/xxjwxc/uber_go_guide_cn)
  - https://github.com/golang/go/wiki/CodeReviewComments#variable-names
  - http://docscn.studygolang.com/doc/effective_go.html
  - https://golang.google.cn/doc/effective_go.html
  - https://studygolang.com/articles/3055?fr=sidebar
  - https://www.cnblogs.com/kotagan/p/11364499.html
- vet  
  - https://blog.csdn.net/u012210379/article/details/50443656
  - https://www.jianshu.com/p/19a44cbc69fb  