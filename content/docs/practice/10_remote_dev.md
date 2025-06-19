---
title: Go 远程开发调试
weight: 10
---

VS Code 是一款开源的代码编辑器，功能强大，支持远程开发调试。

## 搭建环境

要实现 Go 远程开发调试，需要先安装 [Go for Visual Studio Code](https://marketplace.visualstudio.com/items?itemName=golang.Go) 插件。

VS Code 的 Remote 功能由三个插件组成，分别适用于三种不同的场景：

- [Remote - SSH](https://code.visualstudio.com/docs/remote/ssh)：利用 SSH 连接远程主机进行开发。
- [Remote - Container](https://code.visualstudio.com/docs/devcontainers/containers)：连接当前机器上的容器进行开发。
- [Remote - WSL](https://code.visualstudio.com/docs/remote/wsl)：连接子系统（Windows Subsystem for Linux）进行开发。

SSH 模式的原理：

![architecture-vscode-ssh](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/architecture-vscode-ssh.png)
图片来自于 Visual Studio Code 官网

### 连接远程机器

安装插件 [Remote SSH](https://code.visualstudio.com/docs/remote/ssh)。

> 服务器需要支持 SSH 连接。

安装后，点击左下角的 `Open a Remote Window`，选择 `Connect to Host`。

点击 `Add New SSH Host` 配置你的远程机器，或者选择已经配置好的 Hosts。

也可以使用快捷键 `F1` 或者 `ctrl+shift+p` 打开 commands，输入 `Open SSH Configuration File` 直接编辑配置文件：

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

1. 连接 host。
2. 选择 platform：Linux, Windows, macOS。
3. 输入密码建立连接。
4. 点击 `Open Folder` 就可以打开远程机器上的代码目录了。
5. VS Code 会提示远程机器需要安装 Go 扩展，选择安装。

左侧边栏的 Remote Explorer，可以快速打开远程机器上的代码目录：

![vscode-remote-usage](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/vscode-remote-usage.png)

### 配置免密登录

使用快捷键 `F1` 或者 `ctrl+shift+p` 打开 commands，输入 `Open SSH Configuration File` 编辑配置文件：

```
Host shccdfrh75vm7.hpeswlab.net
    HostName shccdfrh75vm7.hpeswlab.net
    User root
    IdentityFile <absolute-path>/.ssh/id_rsa
```

如果没有秘钥，可以使用 `ssh-keygen -t rsa` 命令生成。

将 SSH 公钥添加到远程机器：

```
$ ssh-copy-id username@remote-host
```

如果 `ssh-copy-id` 命令不存在，就手动将 `<absolute-path>/.ssh/id_rsa.pub` 的内容，追加到远程机器的 `~/.ssh/authorized_keys` 文件后面。

## 远程开发

连接到远程主机后，就可以进行远程开发了。可以像本地开发一样查看，修改文件。

![vscode-ssh-dev](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/vscode-ssh-dev.png)

## 远程调试

Go 远程调试本地机器和远程机器都需要[安装 "delve"](https://github.com/derekparker/delve/blob/master/Documentation/installation/README.md)：

```
$ go install github.com/go-delve/delve/cmd/dlv@latest
```

安装完成后需要配置调试工具，点击侧边栏中的 "Run and Debug"，点击 "create a launch.json file" 会在 `.vscode` 目录下创建一个运行配置文件 `launch.json`。

下面是一个调试 Go 程序的 `launch.json` 示例：

```json
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug helm list -A",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/helm",
      "args": ["list", "-A"],
      "env": {
        "HELM_DRIVER": "configmap"
      }
    },
    {
      "name": "Launch test function",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}",
      "args": [
        "-test.run",
        "MyTestFunction"
      ]
    },
    {
      "name": "Launch executable",
      "type": "go",
      "request": "launch",
      "mode": "exec",
      "program": "absolute-path-to-the-executable"
    },
    {
      "name": "Launch test package",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}"
    },
    {
      "name": "Attach to local process",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "processId": 12784
    }
  ]
}
```

常用属性：

- `type`：调试器类型。`node` 用于内置的 Node 调试器，`php` 和 `go` 用于 PHP 和 Go 扩展。
- `request`：值可以是 `launch`，`attach`。当需要对一个已经运行的的程序 debug 时才使用 `attach`，其他时候使用 `launch`。
- `mode`：值可以是 `auto`，`debug`，`remote`，`test`，`exec`。 对于 `attach` 只有 `local`，`remote`。
- `program`：启动调试器时要运行的可执行文件或文件。
- `args`： 传递给调试程序的参数。
- `env`：环境变量（空值可用于 "取消定义 "变量），`env` 中的值会覆盖 `envFile` 中的值。
- `envFile`：包含环境变量的 dotenv 文件的路径。
- `cwd`：当前工作目录，用于查找依赖文件和其他文件。
- `port`：连接到运行进程时的端口。
- `stopOnEntry`：程序启动时立即中断。
- `console`：使用哪种控制台，例如内部控制台、集成终端或外部终端。
- `showLog`：是否在调试控制台打印日志, 一般为 `true`。
- `buildFlags`：构建程序时需要传递给 Go 编译器的 Flags，例如 `-tags=your_tag`。
- `remotePath`：`mode` 为 `remote` 时, 需要指定调试文件所在服务器的绝对路径。
- `processId`：进程 id。
- `host`：目标服务器地址。
- `port`：目标端口。

常用的变量：

- `${workspaceFolder}` 调试工作空间下的根目录下的所有文件。
- `${file}` 调试当前文件。
- `${fileDirname}` 调试当前文件所在目录下的所有文件。

更多的属性和变量可以查看 [VS Code Debugging 文档](https://code.visualstudio.com/docs/editor/debugging)。

配置好 `launch.json` 后，在代码上打上断点，打开侧边栏的 "Run and Debug"，选择运行的配置，就可以开始调试了。

