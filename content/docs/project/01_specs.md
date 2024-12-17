---
title: 项目规范
weight: 1
---

对于多人协作的项目，每个人的开发习惯都不相同，没有统一的规范，会造成很多问题。比如：代码风格不统一，目录结构杂乱无章，API 定义不统一（URL 和错误码）。

一个好的规范可以提高软件质量，提高开发效率，降低维护成本。

## 选择开源协议

开源项目需要选择一个开源协议，如果不准备开源，就用不到开源协议。

开源许可证，大概有几十种，可分为两大类：

- 宽松式（permissive）许可证：最基本的类型，对用户几乎没有限制，用户可以修改代码后闭源。例如 MIT，Apache 2.0 等。
- Copyleft 许可证：比宽松式许可证的限制要多，修改源码后不可以闭源。例如 GPL，Mozilla（MPL）等。

如何选择自己项目的开源许可证，可以根据下面的图示：

![](https://raw.githubusercontent.com/shipengqi/illustrations/5567f7aabb4fd6d2cbf6ae2d619502b1a0191be4/go/licenses.png)
图片来自于阮一峰的网络日志

## 文档规范

### README

`README.md` 是开发者了解一个项目时阅读的第一个文档，会放在项目的根目录下。主要是用来介绍项目的功能、安装、部署和使用。

    # 项目名称
    
    <!-- 项目描述、Logo 和 Badges -->
    
    ## Overview
    
    <!-- 描述项目的核心功能 -->

    ## Getting started
    
    ### Installation

    <!-- 如何安装 -->
 
    ### Usage
    
    <!-- 用法 -->
    
    ## Contributing
    
    <!-- 如何提交代码 -->


也可以使用快速生成 README 文档的在线工具 [readme.so](https://readme.so/)。

### 项目文档

项目文档一般会放在 `/docs` 目录下。项目文档一般有两类：

- 开发文档：用来说明项目的开发流程，如何搭建开发环境、构建、测试、部署等。
- 用户文档：针对用户的使用文档，一般包括功能介绍文档、安装文档、API 文档、最佳实践、操作指南、常见问题等。

文档最好包含英文和中文 2 个版本。

文档目录结构示例：

```
docs
├── dev                              # 开发文档
│   ├── en-US/                       # 英文版
│   └── zh-CN                        # 中文版
│       ├── contributing.md          
│       └── development.md           
├── guide   
│   ├── en-US/                       # 英文版
│   └── zh-CN                        # 中文版
│       ├── api/                     # API 文档
│       ├── practice/                # 最佳实践，存放一些比较重要的实践文章
│       ├── faq/                     # 常见问题
│       ├── installation/            # 安装文档
│       └── README.md                # Guide 入口文件
```
