---
title: Go 远程开发调试
weight: 10
draft: true
---

# Go 远程开发调试

基于 VS Code

## 连接远程机器

首先要安装扩展 [Remote SSH](https://code.visualstudio.com/docs/remote/ssh)。

> 服务器需要支持 SSH 连接。

安装后，点击左下角的 `Open a Remote Window`，选择 `Connect to Host`。

点击 `Add New SSH Host` 配置你的远程机器，或者选择已经配置好的 Hosts。

也可以使用快捷键 `ctrl+shift+p` 打开 commands，输入 `Open SSH Configuration File` 直接编辑配置文件：

```
# Read more about SSH config files: https://linux.die.net/man/5/ssh_config
Host shcCDFrh75vm8.hpeswlab.net
    HostName shcCDFrh75vm8.hpeswlab.net
    Port 22
    User root

Host shccdfrh75vm7.hpeswlab.net
    HostName shccdfrh75vm7.hpeswlab.net
    User root
```

配置好之后：

1. 连接 host，
2. 选择 platform：Linux, Windows, macOS
3. 输入密码建立连接
4. 点击 `Open Folder` 就可以打开远程机器上的代码目录了。
5. VS Code 会提示远程机器需要安装 Go 扩展，选择安装。

### 配置免密登录

使用快捷键 `ctrl+shift+p` 打开 commands，输入 `Open SSH Configuration File` 编辑配置文件：

```
Host shccdfrh75vm7.hpeswlab.net
    HostName shccdfrh75vm7.hpeswlab.net
    User root
    IdentityFile <absolute-path>/.ssh/id_rsa
```

如果没有秘钥使用 `ssh-keygen -t rsa` 命令生成。

将 SSH 公钥添加到远程机器：

```
$ ssh-copy-id username@remote-host
```

如果 `ssh-copy-id` 不存在，就手动将 `<absolute-path>/.ssh/id_rsa.pub` 的内容，追加到远程机器的 `~/.ssh/authorized_keys` 文件后面。

## 远程开发

## 远程调试