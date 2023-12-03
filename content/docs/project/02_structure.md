---
title: 项目的目录结构
weight: 2
---

# 项目的目录结构

一个好的目录结构设计应该是易维护、易扩展的。至少要满足以下几个要求：

- **命名清晰**：目录命名要清晰、简洁，能清晰地表达出该目录实现的功能，并且目录名最好用**单数**。单数足以说明这个目录的功能，避免单复混用。
- **功能明确**：一个目录所要实现的功能应该是明确的、并且在整个项目目录中具有很高的辨识度。当需要新增一个功能时，能够非常清楚地知道把这个功能放在哪个目录下。
- **全面性**：目录结构应该尽可能全面地包含研发过程中需要的功能，例如文档、脚本、源码管理、API 实现、工具、第三方包、测试、编译产物等。
- **可预测性**：项目规模一定是从小到大的，所以一个好的目录结构应该能够在项目变大时，仍然保持之前的目录结构。
- **可扩展性**：每个目录下存放了同类的功能，在项目变大时，这些目录应该可以存放更多同类功能。

根据项目的功能，目录结构可以分为两种：

- 平铺式目录结构
- 结构化目录结构

## 平铺式目录结构

当一个项目是一个工具库时，适合使用平铺式目录结构。项目的代码都存放在项目的根目录下，可以**减少项目引用路径的长度**。例如 `github.com/golang/glog`：

```
$ ls glog/
glog_file.go glog_flags.go glog.go  glog_test.go  go.mod go.sum LICENSE  README
```

## 结构化目录结构

当一个项目是一个应用时，适合使用结构化目录结构。目前 Go 社区比较推荐的结构化目录结构是 [project-layout](https://github.com/golang-standards/project-layout)。

下面是一套结合 project-layout 总结出的目录结构：

```
├── api                       # 存放不同类型的 API 定义文件
│   └── swagger               # Swagger API 文档
├── cmd                       # cmd 下可以包含多个组件目录，组件目录下存放各个组件的 main 包
│   └── apiserver            
│       └── apiserver.go
├── chart                     # helm chart 文件
├── conf                      # 项目部署的配置文件
├── docs                      # 项目文档
│   ├── dev
│   │   ├── en-US
│   │   └── zh-CN
│   ├── guide
│   │   ├── en-US
│   │   └── zh-CN
│   └── README.md
├── examples                  # 项目使用示例
├── go.mod
├── go.sum
├── hack                      # 项目构建，持续集成相关的文件           
│   ├── include               # 存放 makefile 文件，实现入口 Makefile 文件中的各个功能
│   ├── scripts               # 存放 Shell 脚本
│   ├── docker                # 包含多个组件目录，组件目录下存放各个组件的 Dockerfile，Docker Compose 文件等
│   │   └── apiserver
│   │       └── Dockerfile
│   │
├── internal                  # internal 下可以包含多个组件目录，组件目录下存放各个组件的业务代码
│   ├── apiserver             # 组件的业务逻辑代码
│   │   ├── apiserver.go      # 组件应用的入口文件
│   │   ├── config            # 根据 options 创建组件应用的配置
│   │   ├── controller        # HTTP API 的实现，包含请求参数的解析、校验、返回响应，具体的业务逻辑在 service 目录下
│   │   │   └── v1            # API 的 v1 版本 
│   │   │       └── user
│   │   ├── options           # 组件的命令行选项，可以 internal/pkg/options 中的命令行选项
│   │   ├── service           # 具体的业务逻辑
│   │   │   └── v1            # v1 版本 
│   │   │       └── user
│   │   ├── store             # 数据库操作的代码，可以创建多个目录，对应不同的数据库
│   │   │   ├── mysql         
│   │   │   │   ├── mysql.go
│   │   │   │   └── user
│   │   │   └── fake          
│   │   │
│   ├── pkg                   # 仅项目内可用的工具包
│   │   ├── code              # 项目内共享的错误码
│   │   ├── options           # 项目内共享的命令行选项
│   │   └── util
├── LICENSE
├── Makefile                  # Makefile 入口文件
├── pkg                       # 全局可用的工具包，可以被外部引用
│   └── util
├── README.md
├── test                      # 存放测试代码
│   ├── testdata              # 测试数据
│   └── e2e                   # e2e 测试代码
```