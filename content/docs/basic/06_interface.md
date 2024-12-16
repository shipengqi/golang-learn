---
title: 接口
weight: 6
draft: true
---

Go 支持接口数据类型，接口是一组方法的集合，任何其他类型只要实现了这些方法就是实现了这个接口，无须显示声明。

**接口只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要**。比如写单元测试，需要 mock 一个类型时，就可以使用接口，mock 的类型和被测试的类型都实现同一个接口即可。

## 原理

接口的结构体 `iface` 和 `eface`：

```go
// src/runtime/runtime2.go#L204
type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type eface struct {
	_type *_type
	data  unsafe.Pointer
}
```

- `iface` 表示了包含一组方法的接口。
- `eface` 表示空接口，也就是 `interface{}`。空接口在使用中很常见，所以在实现时专门定义了一个类型。

### 结构体嵌入接口类型

Go 语言的结构体还可以嵌入接口类型。

```go
type Interface interface {             
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}

// Array 实现 Interface 接口
type Array []int

func (arr Array) Len() int {
    return len(arr)
}

func (arr Array) Less(i, j int) bool {
    return arr[i] < arr[j]
}

func (arr Array) Swap(i, j int) {
    arr[i], arr[j] = arr[j], arr[i]
}

// 匿名接口(anonymous interface)
type reverse struct {
    Interface
}

// 重写(override)
func (r reverse) Less(i, j int) bool {
    return r.Interface.Less(j, i)
}

// 构造 reverse Interface
func Reverse(data Interface) Interface {
    return &reverse{data}
}

func main() {
    arr := Array{1, 2, 3}
    rarr := Reverse(arr)
    fmt.Println(arr.Less(0,1))
    fmt.Println(rarr.Less(0,1))
}
```

`reverse` 结构体内嵌了一个名为 `Interface` 的 `interface`，并且实现 `Less` 函数，但是
却没有实现 `Len`, `Swap` 函数。

为什么这么设计？

通过这种方法可以让 **`reverse` 实现 `Interface` 这个接口类型，并且仅实现某个指定的方法，而不需要实现这个接口下的所有方法**。

对比一下传统的组合匿名结构体实现重写的写法：

```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}

type Array []int

func (arr Array) Len() int {
    return len(arr)
}

func (arr Array) Less(i, j int) bool {
    return arr[i] < arr[j]
}

func (arr Array) Swap(i, j int) {
    arr[i], arr[j] = arr[j], arr[i]
}

// 匿名struct
type reverse struct {
    Array
}

// 重写
func (r reverse) Less(i, j int) bool {
    return r.Array.Less(j, i)
}

// 构造 reverse Interface
func Reverse(data Array) Interface {
    return &reverse{data}
}

func main() {
    arr := Array{1, 2, 3}
    rarr := Reverse(arr)
    fmt.Println(arr.Less(0, 1))
    fmt.Println(rarr.Less(0, 1))
}
```

匿名接口的优点，**匿名接口的方式不依赖具体实现，可以对任意实现了该接口的类型进行重写**。
