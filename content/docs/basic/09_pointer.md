---
title: 指针
weight: 9
---

**指针和内存地址不能混为一谈**。内存地址是内存中每个字节单元的唯一编号，而指针是一个实体。指针也会分配内存空间，相当于一个保存内存地址的整形变量。

## 指针的限制

### 指针不能参与运算

```go
package main

import "fmt"

func main() {
	a := 1
	b := a
	fmt.Println(b)
	b = &a + 1
}
```

上面的代码编译时会报错：`Invalid operation: &a + 1 (mismatched types *int and untyped int)`。

说明 Go 是不允许对指针进行运算的。

### 不同类型的指针不允许相互转换

```go
package main

func main() {	
	var a int = 100
	var f *float64
	f = &a
}
```

上面的代码编译时会报错：`Cannot use '&a' (type *int) as the type *float64`。

### 不同类型的指针不能比较

因为不同类型的指针之间不能转换，所以也不能赋值。

### 不同类型的指针变量不能相互赋值

同样的由于不同类型的指针之间不能转换，所以也没法使用 `==` 或者 `!=` 进行比较。

## uintptr 类型

`uintptr` 只是一个无符号整型，用于存储内存地址的整形变量。也就是说，和普通的整型一样，是会被 GC 回收的。

## unsafe.Pointer

由于 Go 指针的限制，所以 Go 提供了可以进行**类型转换**的通用指针 `unsafe.Pointer`。

`unsafe.Pointer` 是特别定义的一种指针类型，它指向的对象如果还有用，那么是不会被 GC 回收的。

`unsafe.Pointer` 是各种指针相互转换的桥梁：

- 任何类型的指针 `*T` 可以和 `unsafe.Pointer` 相互转换。
- `uintptr` 可以和 `unsafe.Pointer` 相互转换。

![pointer](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/pointer.png)

指针类型转换示例：

```go
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	v1 := uint(10)
	v2 := int(11)

	fmt.Println(reflect.TypeOf(v1))  // uint
	fmt.Println(reflect.TypeOf(v2))  // int
	fmt.Println(reflect.TypeOf(&v1)) // *uint
	fmt.Println(reflect.TypeOf(&v2)) // *int
	p := &v1
	// 使用 unsafe.Pointer 进行类型转换，将 *int 转为 *uint
	p = (*uint)(unsafe.Pointer(&v2))

	fmt.Println(reflect.TypeOf(p)) // *unit
	fmt.Println(*p)                // 11
}
```
