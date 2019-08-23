---
title: 使用文档主题
---
# 使用文档主题

本文档使用 [Hexo Doc Theme](https://zalando-incubator.github.io/hexo-theme-doc/index.html) 搭建。

## Quick Start
1. 获取源码
```sh
$ git clone git@github.com:zalando-incubator/hexo-theme-doc-seed.git
```

2. `hexo-theme-doc-seed` 的以下文件拷贝到项目 root 目录下，例如 `kubernetes-learn`：
- `source` 目录
- `_data` 目录
- `images` 目录
- `package.json`
- `_config.yaml`
- `.zappr.yaml`

3. 安装依赖
```sh
$ yarn
```

4. 修改 `package.json`，否则 `hexo server` 或者 `hexo s` 可能会找不到命令。
```js
  "hexo": {
    "version": "3.9.0"
  },
  "scripts": {
    "start": "hexo s -p 8082"
  },
```

5. 修改 `_config.yml`
```yml
theme: ../node_modules/hexo-theme-doc

# 如果你的网站存放在子目录中，例如 http://yoursite.com/blog
# 则 url 设为 http://yoursite.com/blog 并把 root 设为 /blog/
url: http://www.shipengqi.top/kubernetes-learn
root: /kubernetes-learn/

# deploy
deploy:
- type: git
  repo: git@github.com:shipengqi/kubernetes-learn.git
  branch: gh-pages
```

6. 启动开发服务，访问 http://localhost:8082 。
```sh
$ yarn start
```

## Index
`source` 目录下创建 `index.md` 文件。这个 `index.md` 文件就是文档首页。

## 添加文档
`source` 目录下创建 `markdown` 文件，例如：

```md
---
title: Lorem Ipsum
---

# Lorem Ipsum

Lorem ipsum
```
也可以创建文档子目录，例如 `source/usage`。

## Sidebar
`source` 目录下的 `_data` 目录下的 `navigation.yaml` 设置 `sidebar` 和其他的一些配置。

```yml
logo:
  text: My Documentation
  type: link
  path: index.html

main:
- text: PROJECTS
  type: label
- text: My Awesome Projects
  type: link
  path: projects/my-awesome-project.html
  children:
  - text: My Awesome Projects Page 1
    type: link
    path: projects/my-awesome-project-page-1.html
```

- **logo**: navigation Logo
- **main**: left sidebar

对于每个导航项，必须定义一个 `type`，并根据类型定义 `text` 和 `path` 等其他属性。
每个导航项，也可以定义一个 `children`，这个属性可以嵌套导航项。

### type
`type` 有两种类型：
- **label**: 导航项的标签
- **link**: 导航项下级文档 link

**`link` 类型的导航项的 `path` 的值是文件的路径，但注意扩展名为`.html`**。

## Favicon
```yml
theme_config:
  favicon: images/favicon.ico
```