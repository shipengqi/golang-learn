---
title: 错误处理
weight: 14
---

## Error 和 Panic

Go 的错误设计中，错误应该明确地当成业务的一部分，**任何可以预见的问题都需要做错误处理**。通常 Go 的函数返回两个值，一个是函数的结果，
另一个是错误对象。如果没有发生错误，错误对象就是 `nil`。这种方式，强迫调用者对错误进行处理，以防遗漏任何运行时可能的错误。

**异常则是意料之外的**，例如数组越界，向空 `map` 添加键值对等，Go 遇到异常会自动触发 panic（恐慌），触发 panic 程序会自动退出。
除了程序自动触发异常，一些不允许的错误也可以手动触发异常，例如连接数据库失败主动调用 `panic`。

## 代码分层结构中如何处理错误

一个常见的三层调用：

```go
// controller
if err := mode.ParamCheck(param); err != nil {
    log.Errorf("param=%+v", param)
    return errs.ErrInvalidParam
}

return mode.ListTestName("")

// service
_, err := dao.GetTestName(ctx, settleId)
    if err != nil {
    log.Errorf("GetTestName failed. err: %v", err)
    return errs.ErrDatabase
}

// dao
if err != nil {
    log.Errorf("GetTestDao failed. uery: %s error(%v)", sql, err)
}
```

上面代码的问题：

- 分层结构下，各个层级都打印了日志
- 原生的 error 不包含堆栈信息

### 分层处理

对于上面的问题，可以使用 `github.com/pkg/errors` 来处理错误。这个库主要的方法：

```go
// 需要生成一个新的错误是使用, 包含堆栈信息
func New(message string) error

// 如果有一个现成的 error ，需要对他进行包装处理，可以选择（WithMessage/WithStack/Wrapf）

// 包装 error 同时附加堆栈信息和 message
func Wrapf(err error, format string, args ...interface{}) error

// 只对 error 追加新的 message
func WithMessage(err error, message string) error

// 只对 error 追加堆栈信息
func WithStack(err error) error

// 获得最根本的错误原因
func Cause(err error) error
```

![error-handle](https://github.com/shipengqi/illustrations/blob/63b742cee77624dfc1f8de5946f051e3c8f395be/go/error-handle.png?raw=true)

{{< callout type="info" >}}
调用第三方库或者标准库也考虑使用 `errors.Wrap` 保存堆栈信息。
{{< /callout >}}

{{< callout type="info" >}}
对于 error 级别的日志打印堆栈，`warn` 和 `info` 可以不打印堆栈。
{{< /callout >}}

## ErrGroup

Go 的扩展库 `golang.org/x/sync` 提供了 `errgroup` 包，可以用来控制并发任务，并且可以处理并发任务返回的错误。

使用方式和原理参考 [并发编程/ErrGroup](../../concurrency/12_errorgroup/)。
