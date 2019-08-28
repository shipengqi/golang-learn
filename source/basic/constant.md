---
title: 常量
---
# 常量

`const` 声明常量，运行时不可改变（只读），注意常量的**底层数据类型只能是基础类型（布尔型、数值型和字符串型）**：
```go
const 常量名字 类型 = 表达式
```

"类型"可以省略。也就是如果没有类型，可以通过表达式推导出类型。

比如：
```go
// 声明一个`string`类型
const b string = "abc"
const a = "abc"

// 声明一组不同类型
const c, f, s = true, 2.3, "four" // bool, float64, string

// 批量声明多个常量
const (
  Unknown = 0
  Female = 1
  Male = 2
)

const strSize = len("hello, world")
```
常量表达式的值在**编译期计算**。因此常量表达式中，函数必须是内置函数。如 `unsafe.Sizeof()`，`len()`, `cap()`。

**常量组中，如果不指定类型和初始值，那么就和上一行非空常量右值相同**：
例如：
```go
const (
	a = 1
	b
	c = 2
	d
)

fmt.Println(a, b, c, d) // "1 1 2 2"
```

## iota
**Go 中没有枚举的定义，但是可以使用 `iota`**，`iota` 标识符可以认为是**一个可以被编译器修改的常量**。
在 `const` 声明中，被重置为 `0`，在第一个声明的常量所在的行，`iota` 将会被置为 `0`，然后在每一个有常量声明的
行加 `1`。
```go
const (
	a = iota   // 0
	b          // 1
	c          // 2
	d = "ha"   // "ha", iota += 1
	e          // "ha" ,不指定类型和初始值，那么就和上一行非空常量右值相同,  iota += 1
	f = 100    // 100, iota +=1
	g          // 100,不指定类型和初始值，那么就和上一行非空常量右值相同,  iota +=1
	h = iota   // 7, 中断的 iota 计数必须显示恢复
	i          // 8
)

const (
	i = 1 << iota // 1, 1 << 0
	j = 3 << iota // 6, 3 << 1
	k             // 12, 3 << 2
	l             // 24, 3 << 3
)
```