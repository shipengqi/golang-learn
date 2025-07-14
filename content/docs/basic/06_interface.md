---
title: 接口
weight: 6
---

Go 支持接口数据类型，接口是一组方法的集合，任何其他类型只要实现了这些方法就是实现了这个接口，无须显示声明。

**接口只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要**。比如写单元测试，需要 mock 一个类型时，就可以使用接口，mock 的类型和被测试的类型都实现同一个接口即可。

## 原理

`interface` 的底层结构：

```go
type iface struct {
	tab  *itab
	data unsafe.Pointer // 指向实际的数据
}

type itab struct {
	inter  *interfacetype // 表示接口类型，静态类型
	_type  *_type // 表示具体实现了该接口的类型，动态类型
	link   *itab
	hash   uint32
	bad    bool
	inhash bool
	unused [2]byte
	fun    [1]uintptr // 这是一个函数指针数组，用于存储实现了该接口的方法的函数指针，
                      // 当接口调用某个方法时，根据 fun 中的函数指针找到具体的实现
}
```

实际上，`iface` 描述的是非空接口，它包含方法；与之相对的是 `eface`，描述的是**空接口**，不包含任何方法，**Go 里有的类型都 “实现了” 空接口**。

```go
type eface struct {
    _type *_type // 表示空接口所承载的具体的实体类型
    data  unsafe.Pointer // 具体的值
}
```

### 动态类型和静态类型

{{< callout type="info" >}}
**接口变量可以存储任何实现了该接口的变量**。
{{< /callout >}}

Go 中最常见的 `Reader` 和 `Writer` 接口：

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

接下来，就是接口之间的各种转换和赋值了：

```go
var r io.Reader
tty, err := os.OpenFile("/Users/s/Desktop/test", os.O_RDWR, 0)
if err != nil {
    return nil, err
}
r = tty
```

1. `io.Reader` 是 `r` 的静态类型。它的动态类型为 `nil`。
2. `r = tty` 将 `r` 的动态类型变成 `*os.File`。

![go-interface-reflect](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/go-interface-reflect.png)

`*os.File` 其实还实现了 `io.Writer` 接口：

```go
var w io.Writer
w = r.(io.Writer)
```

之所以用断言，而不能直接赋值，是因为 `r` 的**静态类型**是 `io.Reader`，并没有实现 `io.Writer` 接口。断言能否成功，看 `r` 的**动态类型**是否符合要求。

### 空接口类型

```go
var empty interface{}
empty = r
```

![go-empty-interface](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/go-empty-interface.png)

由于 `empty` 是一个空接口，因此所有的类型都实现了它，`w` 可以直接赋给它，不需要执行断言操作。

## 接口赋值

接口在赋值的时候会初始化对应的底层结构，将具体的动态类型转为静态类型：

```go
func convT2I(inter *interfacetype, tab *itab, t *_type, v unsafe.Pointer) (iface, bool) {
   var i iface
   if tab == nil {
      tab = getitab(inter, t, false)
   }
   if tab != nil {
      return i, false
   }
   i.tab = tab
   i.data = v
   return i, true
}
```

## 断言

类型断言也依赖于接口的数据结构，通过检查接口的 `_type` 来判断类型是否于接口的实际类型匹配：

```go
func assertE2I(inter *interfacetype, e eface) (i iface) {
    tab := getitab(inter, e.type, true)
    i.tab = tab
    i.data = e.data
    return
}
```