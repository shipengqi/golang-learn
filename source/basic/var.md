---
title: 变量
---
# 变量
`var` 声明变量，必须使用空格隔开：
```go
var 变量名字 类型 = 表达式
```
**类型**或者**表达式**可以省略其中的一个。也就是如果没有类型，可以通过表达式推断出类型，**没有表达式，将会根据类型初始化为对应的零值**。

**零值** 并不是空值，而是一种“变量未填充前”的默认值，通常为 `0`，对应关系：
- 数值类型：`0`
- 布尔类型：`false`
- 字符串: `""`
- 接口或引用类型（包括 `slice`、指针、`map`、`chan` 和函数）：`nil`

注意：
- `map` 的零值是 `nil`， 也就是或如果用 `var testMap map[string]string` 的方式声明，是不能直接通过 unmarshal 或 `map[key]` 操
作，应该使用 `make` 函数。
- `slice` 的零值是 `nil`，不能直接通过下标操作。应该使用 `make` 函数。
- 对于 `struct` 的指针，要注意使用 `var testStruct *testResponse`的方式声明，如果 `testResponse` 内嵌套结构体指针，unmarshal 
会失败，因为指针的零值是 `nil`，应该使用 `testStruct := &testResponse{}` 的方式。

## 声明一组变量
```go
var 变量名字, 变量名字, 变量名字 ... 类型 = 表达式, 表达式, 表达式, ...
```
比如：
```go
// 声明一组 `int` 类型
var i, j, k int                 // int, int, int

// 声明一组不同类型
var b, f, s = true, 2.3, "four" // bool, float64, string

var (
  i int
  pi float32
  prefix string
)
```

## 简短声明
**`:=` 只能在函数内使用，不能提供数据类型**，Go 会自动推断类型：
```go
变量名字 := 表达式
```

```go
var x = 100

func main() {
	fmt.Println(&x, x)
	x := "abc"
	fmt.Println(&x, x)
}
```
上面的代码中 `x := "abc"` 相当于重新定义并初始化了同名的局部变量 `x`，因为**不在同一个作用域**，所以打印出来的结果完全不同。

简短声明，并不总是重新定义比变量，要避免重新定义，首先要在同一个作用域中，至少有一个新的变量被定义：
```go
func main() {
	x := 100
	fmt.Println(&x, x)
	x, y := 200, 300   // 一个新的变量 y，这里的简短声明 x 就是赋值操作
	fmt.Println(&x, x)
}
```

如果重复使用简短声明定义一个变量，会报错：
```go
x := 100
fmt.Println(&x)
x := 200 // 错误， no new variables on left side of :=
```

## 赋值
常见的赋值的方式：
```go
x = 1                       // 命名变量的赋值
*p = true                   // 通过指针间接赋值
person.name = "bob"         // 结构体字段赋值
count[x] = count[x] * scale // 数组、slice 或 map 的元素赋值
count[x] *= scale           // 等价于 count[x] = count[x] * scale，但是省去了对变量表达式的重复计算
x, y = y, x                 // 交换值
f, err = os.Open("foo.txt") // 左边变量的数目必须和右边一致，函数一般会返回一个 error 类型
v, ok = m[key]              // map 查找，返回布尔值类表示操作是否成功
v = m[key]                  // map 查找，也可以返回一个值，失败时返回零值
```

不管是隐式还是显式地赋值，在赋值语句左边的变量和右边最终的求到的值必须有相同的数据类型。这就是**可赋值性**。

进行**多变量赋值**时，首先计算出所有右值，然后再依次赋值：
```go
x, y := 1, 2
x, y = y+3, x+2 // 先计算出 y+3, x+2, 然后赋值
```