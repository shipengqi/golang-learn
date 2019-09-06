---
title: 错误
---
# 错误

## error 类型
`error` 类型是内置的接口类型。`error` 类型可能是 `nil` 或者 `non-nil`，`nil` 表示成功。

### 错误处理
当函数调用返回错误时，最常用的处理方式是**传播错误**，如。
```go
resp, err := http.Get(url)
if err != nil{ // 将这个HTTP错误返回给调用者
    return nil, err
}


doc, err := html.Parse(resp.Body)
resp.Body.Close()
if err != nil {
	// fmt.Errorf函数使用fmt.Sprintf格式化错误信息并返回
	// 使用该函数前缀添加额外的上下文信息到原始错误信息。
  return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
}
```
由于错误信息经常是以链式组合在一起的，所以错误信息中应避免大写和换行符。

编写错误信息时，我们要确保错误信息对问题细节的描述是详尽的。尤其是要注意错误信息表达的一致性，即相同的函数或同包内
的同一组函数返回的错误在构成和处理方式上是相似的。

根据不同的场景，我们可能要对错误做些特殊处理，比如错误重试机制，或者打印错误日志，或者直接忽略错误。

### 文件结尾错误
`io` 包在任何由文件结束引起的读取失败都返回同一个错误 `io.EOF`：
```go
in := bufio.NewReader(os.Stdin)
for {
    r, _, err := in.ReadRune()
    if err == io.EOF {
        break // finished reading
    }
    if err != nil {
        return fmt.Errorf("read failed:%v", err)
    }
    // ...use r…
}
```

## Panic 异常
Go 运行时错误会引起 `painc` 异常。
一般而言，当 `panic` 异常发生时，程序会中断运行，并立即执行在该 `goroutine` 中被延迟的函数（`defer` 机制）。随后，程序崩溃
并输出日志信息。

由于 `panic` 会引起程序的崩溃，因此 `panic` 一般用于严重错误，如程序内部的逻辑不一致。但是对于大部分漏洞，我们应该使
用 Go 提供的错误机制，而不是 `panic`，尽量避免程序的崩溃。

### panic 函数
`panic` 函数接受任何值作为参数。当某些不应该发生的场景发生时，我们就应该调用 `panic`。

### panic 详情中都有什么
```sh
panic: runtime error: index out of range

goroutine 1 [running]:
main.main()
/Users/haolin/GeekTime/Golang_Puzzlers/src/puzzlers/article19/q0/demo47.go:5 +0x3d
exit status 2
```
第一行是 `panic: runtime error: index out of range`。其中的 `runtime error` 的含义是，这是一个 `runtime` 代码包中
抛出的` panic`。

`goroutine 1 [running]`，它表示有一个 ID 为1的 `goroutine` 在此 `panic` 被引发的时候正在运行。这里的 ID 其实并不重要。

`main.main()` 表明了这个 `goroutine` 包装的 go 函数就是命令源码文件中的那个`main`函数，也就是说这里的 `goroutine` 正
是**主 goroutine**。

再下面的一行，指出的就是这个 `goroutine` 中的哪一行代码在此 panic 被引发时正在执行。含了此行代码在其所属的源码文件中的行数，
以及这个源码文件的绝对路径。

`+0x3d` 代表的是：此行代码相对于其所属函数的入口程序计数偏移量。用处并不大。

`exit status 2` 表明我的这个程序是以退出状态码2结束运行的。**在大多数操作系统中，只要退出状态码不是 0，都意味着程序运行的非正
常结束**。在 Go 语言中，**因 panic 导致程序结束运行的退出状态码一般都会是 2**。


### 从 panic 被引发到程序终止运行的大致过程是什么

此行代码所属函数的执行随即终止。紧接着，控制权并不会在此有片刻停留，它又会立即转移至再上一级的调用代码处。控制权如此一级一
级地沿着调用栈的反方向传播至顶端，
也就是我们编写的最外层函数那里。

这里的最外层函数指的是go函数，对于主 goroutine 来说就是 `main` 函数。但是控制权也不会停留在那里，而是被 Go 语言运行时系统收回。

随后，程序崩溃并终止运行，承载程序这次运行的进程也会随之死亡并消失。与此同时，在这个控制权传播的过程中，panic 详情会被逐
渐地积累和完善，并会在程序终止之前被打印出来。

### 怎样让 panic 包含一个值，以及应该让它包含什么样的值
其实很简单，在调用 `panic` 函数时，把某个值作为参数传给该函数就可以了。`panic` 函数的唯一一个参数是空接口
（也就是`interface{}`）类型的，所以从语法上讲，它可以接受任何类型的值。

但是，我们**最好传入 `error` 类型的错误值，或者其他的可以被有效序列化的值。这里的“有效序列化”指的是，可以更易读地去表示
形式转换**。

## Recover 捕获异常
一般情况下，我们不能因为某个处理函数引发的 `panic` 异常，杀掉整个进程，可以使用 `recover` 函数恢复 `panic` 异常。

`panic` 时会调用 `recover`，但是 `recover` 不能滥用，可能会引起资源泄漏或者其他问题。我们可以将 `panic value` 设置成特
殊类型，来标识某个 `panic` 是否应该被恢复。**`recover` 只能在 `defer` 修饰的函数中使用**:
```go
func soleTitle(doc *html.Node) (title string, err error) {
	type bailout struct{}
	defer func() {
		switch p := recover(); p {
            case nil:       // no panic
            case bailout{}: // "expected" panic
                err = fmt.Errorf("multiple title elements")
            default:
                panic(p) // unexpected panic; carry on panicking
		}
	}()
    panic(bailout{}) 
}
```

上面的代码，`deferred` 函数调用 `recover`，并检查 `panic value`。当 `panic value` 是 `bailout{}` 类型时，`deferred` 函数生
成一个 `error` 返回给调用者。
当 `panic value` 是其他 `non-nil` 值时，表示发生了未知的 `panic` 异常。

### 正确调用 recover 函数
```go
package main

import (
    "fmt"
    "errors"
)

func main() {
    fmt.Println("Enter function main.")
    // 引发 panic。
    panic(errors.New("something wrong"))
    p := recover()
    fmt.Printf("panic: %s\n", p)
    fmt.Println("Exit function main.")
}
```
上面的代码，`recover` 函数调用并不会起到任何作用，甚至都没有机会执行。因为 panic 一旦发生，控制权就会讯速地沿着调用栈的反方向
传播。所以，**在 panic 函数调用之后的代码，根本就没有执行的机会**。

先调用 `recover` 函数，再调用 `panic` 函数会怎么样呢？
如果在我们调用 `recover` 函数时未发生 panic，那么该函数就不会做任何事情，并且只会返回一个 `nil`。

**`defer` 语句调用 `recover` 函数才是正确的打开方式**。

无论函数结束执行的原因是什么，其中的 `defer` 函数调用都会在它即将结束执行的那一刻执行。即使导致它执行结束的原因是一
个 panic 也会是这样。

要注意，我们要**尽量把 `defer` 语句写在函数体的开始处，因为在引发 panic 的语句之后的所有语句，都不会有任何执行机会**。
