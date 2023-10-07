---
title: log
---

# log
log 模块用于在程序中输出日志。

```go
package main

import "log"

func main() { 
    log.Print("Hello World") // 2019/09/12 13:56:36 Hello World
}
```

## Logger
通过 `New` 函数可以创建多个 `Logger` 实例，函数声明如下：
```go
func New(out io.Writer, prefix string, flag int) *Logger
```

参数：
- `out`：日志输出的 IO 对象，通常是标准输出 `os.Stdout`，`os.Stderr`，或者绑定到文件的 IO。
- `prefix`：日志前缀，可以是任意字符串。
- `flag`：日志包含的通用信息标识位

一条日志的结构：
```
{日志前缀} {标识1} {标识2} ... {标识n} {日志内容}
```

标识通过 `flag` 参数设置，当某个标识被设置，会在日志中进行显示，log 模块中已经提供了如下标识，多个标识通过 `|` 组合：
- Ldate 显示当前日期（当前时区）
- Ltime 显示当前时间（当前时区）
- microseconds 显示当前时间（微秒）
- Llongfile 包含路径的完整文件名
- Lshortfile 不包含路径的文件名
- LUTC Ldata 和 Ltime 使用 UTC 时间
- LstdFlags 标准 Logger 的标识，等价于 Ldate | Ltime

```go
package main

import (
	"log"
    "os"
)

func main() {
    prefix := "[THIS IS THE LOG]"
    logger := log.New(os.Stdout, prefix, log.LstdFlags | log.Lshortfile)
    logger.Print("Hello World") // [THIS IS THE LOG]22019/09/12 12:34:07 log.go:11: Hello World
}
```

## 分类
log 模块中日志输出分为三类，
- Print，输出日志。
- Fatal，在执行完 Print 之后，执行 `os.Exit(1)`。
- Panic。在执行完 Print 之后调用 `panic()` 方法。

除了基础的 `Print` 之外，还有 `Printf` 和 `Println` 方法对输出进格式化，`Fatal` 和 `Panic` 也类似。

## Level
`log` 包没有提供日志分级的功能，需要自己实现：
```go
package main

import (
	"log"
    "os"
)

func main() {
    var (
	logger = log.New(os.Stdout, "INFO: ", log.Lshortfile)
	infof = func(info string) {
		logger.Print(info)
	}
    )
    infof("Hello world")
}
```