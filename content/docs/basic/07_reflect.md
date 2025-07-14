---
title: 反射
weight: 7
---

Go 提供了一种机制在运行时更新变量和检查它们的值、调用它们的方法，但是在**编译时并不知道这些变量的具体类型**，这称为反射机制。

使用反射的常见场景有以下两种：

- **不能明确接口调用哪个函数**，需要根据传入的参数在运行时决定。
- **不能明确传入函数的参数类型**，需要在运行时处理任意对象。

不推荐使用反射的理由：

- 代码可读性差。
- Go 在编译过程中，编译器能提前发现一些类型错误，但是对于反射代码是无能为力的。所以包含反射相关的代码，很可能会运行很久，才会出错，这时候经常是直接 panic，可能会造成严重的后果。
- 反射对性能影响比较大。所以，对于一个项目中处于运行效率关键位置的代码，尽量避免使用反射特性。

## types 和 interface

Go 中，每个变量都有一个静态类型，在编译阶段就确定了的，比如 `int, float64, []int` 等等。注意，**这个类型是声明时候的类型，不是底层数据类型**。

例如：

```go
type MyInt int

var i int
var j MyInt
```

尽管 `i`，`j` 的底层类型都是 `int`，但它们是**不同的静态类型**，除非进行类型转换，否则，`i` 和 `j` 不能同时出现在等号两侧。`j` 的静态类型就是 `MyInt`。

反射主要与 `interface{}` 类型相关。

## reflect 包

反射的相关函数在 `reflect` 包中。`reflect` 包里定义了一个接口和一个结构体，即 `reflect.Type` 和 `reflect.Value`

- `reflect.Type` 是一个接口，主要提供关于类型相关的信息，所以它和 `_type` 关联比较紧密；
- `reflect.Value` 是一个一个结构体变量，包含类型信息以及实际值。结合 `_type` 和 `data` 两者，因此可以用来获取甚至改变类型的值。
- `func TypeOf(i interface{}) Type` 函数可以用来获取 `reflect.Type`。
- `func ValueOf(i interface{}) Value` 函数可以用来获取 `reflect.Value`。

### TypeOf

`reflect.TypeOf` 获取值的类型信息。

`reflect.TypeOf` 接受任意的 `interface{}` 类型, 并以返回其动态类型 `reflect.Type`：

```go
t := reflect.TypeOf(3)  // a reflect.Type
fmt.Println(t.String()) // "int"
fmt.Println(t)          // "int"

type X int
func main() {
	var a X = 20
	t := reflect.TypeOf(a)
	fmt.Println(t.Name(), t.Kind()) // X int
}
```

上面的代码，**注意区分 `Type` 和 `Kind`，前者表示真实类型（静态类型），后者表示底层类型（动态类型）**。所以在判断类型时，要选择正确的方式。

```go
type X int
type Y int
func main() {
	var a, b X = 10, 20
	var c Y = 30
	ta, tb, tc := reflect.TypeOf(a), reflect.TypeOf(b), reflect.TypeOf(c)
	fmt.Println(ta == tb, ta == tc) // true false
	fmt.Println(ta.Kind() == tc.Kind()) // true
}
```

#### 原理

`TypeOf` 源码：

```go
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}

func toType(t *rtype) Type {
	if t == nil {
		return nil
	}
	return t
}
```

- `emptyInterface` 和上面提到的 `eface`是一回事，并且在不同的源码包：前者在 `reflect` 包，后者在 `runtime` 包
- `eface.typ` 就是动态类型。
- `toType` 只是做了一个类型转换。


### ValueOf

`reflect.ValueOf` 专注于对象实例数据读写。

`reflect.ValueOf` 接受任意的 `interface{}` 类型, 并以 `reflect.Value` 形式返回其动态值：

```go
v := reflect.ValueOf(3) // a reflect.Value
fmt.Println(v)          // "3"
fmt.Printf("%v\n", v)   // "3"
fmt.Println(v.String()) // <int Value>

type Person struct {
	Name string
	Age  int
}

p := Person{"Alice", 30}
v := reflect.ValueOf(p)
fmt.Println(v.Field(0).String()) // 输出: Alice
fmt.Println(v.Field(1).Int())    // 输出: 30
fmt.Println(v.String()) // <main.Person Value>
```

### DeepEqual

```go
func DeepEqual(x, y interface{}) bool
```

`reflect.DeepEqual` 函数的参数是两个 `interface`，也就是可以输入任意类型，输出 `true` 或者 `flase` **表示输入的两个变量是否是“深度”相等**。

**如果是不同的类型，即使是底层类型相同，相应的值也相同，那么两者也不是“深度”相等**。

```go
type MyInt int
type YourInt int

func main() {
	m := MyInt(1)
	y := YourInt(1)

	fmt.Println(reflect.DeepEqual(m, y)) // false
}
```

`m, y` 底层都是 `int`，而且值都是 1，但是两者静态类型不同，前者是 `MyInt`，后者是 `YourInt`，因此两者不是“深度”相等。

`DeepEqual` 的比较情形：

| 类型 | 深度相等情形 |
|------|--------------|
| `Array` | 相同索引处的元素“深度”相等 |
| `Struct` | 相应字段，包含导出和不导出，“深度”相等 |
| `Func` | 只有两者都是 `nil` 时 |
| `Interfac`e | 两者存储的具体值“深度”相等 |
| `Map` | 1、都为 `nil`；2、非空、长度相等，指向同一个 `map` 实体对象，或者相应的 key 指向的 value “深度”相等 |
| `Pointer` | 1、使用 `==` 比较的结果相等；2、指向的实体“深度”相等 |
| `Slice` | 1、都为 `nil`；2、非空、长度相等，首元素指向同一个底层数组的相同元素，即 `&x[0] == &y[0]` 或者 相同索引处的元素“深度”相等 |
| `numbers`, `bools`, `strings`, and `channels` | 使用 `==` 比较的结果为真 |


一般情况下，`DeepEqual` 的实现只需要递归地调用 `==` 就可以比较两个变量是否是真的“深度”相等。

有一些异常情况：比如 `func` 类型是不可比较的类型，只有在两个 `func` 类型都是 `nil` 的情况下，才是“深度”相等；`float` 类型，由于精度的原因，也是不能使用 `==` 比较的；包含 `func` 类型或者 `float` 类型的 `struct`，`interface`， `array` 等。

## 反射的三大定律

1. 反射是一种检测存储在 `interface` 中的类型和值机制。这可以通过 `TypeOf` 函数和 `ValueOf` 函数得到。
2. 第二条实际上和第一条是相反的机制，它将 `ValueOf` 的返回值通过 `Interface()` 函数反向转变成 `interface` 变量。
3. 如果需要操作一个反射变量，那么它必须是可设置的。