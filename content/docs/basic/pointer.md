---
title: 指针
---
# 指针

**指针和内存地址不能混为一谈**。内存地址是内存中每个字节单元的唯一编号，而指针是一个实体。指针也会分配内存空间，相当于一个
保存内存地址的整形变量。

```go
x := 1
p := &x         // p, of type *int, points to x
fmt.Println(*p) // "1"
*p = 2          // equivalent to x = 2
fmt.Println(x)  // "2"
```

上面的代码，初始化一个变量 `x`，`&` 是取地址操作，`&x` 就是取变量 `x` 的内存地址，那么 `p` 就是一个指针，
类型是 `*int`，`p` 这个指针保存了变量 `x` 的内存地址。接下来 `*p` 表示读取指针指向的变量的值，也就是变量 `x` 的值 1。
`*p`也可以被赋值。

任何类型的指针的零值都是 `nil`。当指针指向同一个变量或者 `nil` 时是相等的。
当一个指针被定义后没有分配到任何变量时，它的值为 `nil`。`nil` 指针也称为空指针。

## 指向指针的指针
```go
var a int
var ptr *int
var pptr **int

a = 3000

/* 指针 ptr 地址 */
ptr = &a

/* 指向指针 ptr 地址 */
pptr = &ptr

/* 获取 pptr 的值 */
fmt.Printf("变量 a = %d\n", a )
fmt.Printf("指针变量 *ptr = %d\n", *ptr )
fmt.Printf("指向指针的指针变量 **pptr = %d\n", **pptr)
```

## 为什么需要指针
相比 Java，Python，Javascript 等引用类型的语言，Golang 拥有类似C语言的指针这个相对古老的特性。但不同于 C 语言，Golang 的指
针是单独的类型，而且也不能对指针做整数运算。从这一点看，Golang 的指针基本就是一种引用。

在学习引用类型语言的时候，总是要先搞清楚，当给一个 `函数/方法` 传参的时候，传进去的是值还是引用。实际上，在大部分引用型语言里，
参数为基本类型时，传进去的大都是值，也就是另外复制了一份参数到当前的函数调用栈。参数为高级类型时，传进去的基本都是引用。

内存管理中的内存区域一般包括 `heap` 和 `stack`，`stack` 主要用来存储当前调用栈用到的简单类型数据：`string`，`boolean`，
`int`，`float` 等。这些类型的内存占用小，容易回收，基本上它们的值和指针占用的空间差不多，因此可以直接复制，`GC` 也比较容易做针对性的
优化。复杂的高级类型占用的内存往往相对较大，存储在 `heap` 中，`GC` 回收频率相对较低，代价也较大，因此传 `引用/指针` 可以避免进行成本较
高的复制操作，并且节省内存，提高程序运行效率。

因此，在下列情况可以考虑使用指针：
1. **需要改变参数的值**
2. **避免复制操作**
3. **节省内存**

而在 Golang 中，具体到高级类型 `struct`，`slice`，`map` 也各有不同。实际上，只有 `struct` 的使用有点复杂，**`slice`，`map`，
`chan`都可以直接使用，不用考虑是值还是指针**。

### `struct`

对于函数（`function`），由函数的参数类型指定，传入的参数的类型不对会报错，例如：
```go
func passValue(s struct){}

func passPointer(s *struct){}
```

对于方法（`method`），接收者（`receiver`）可以是指针，也可以是值，Golang 会在传递参数前自动适配以符合参数的类型。也就是：如果方法的参数
是值，那么按照传值的方式 ，方法内部对 `struct` 的改动无法作用在外部的变量上，例如：
```go
package main

import "fmt"

type MyPoint struct {
    X int
    Y int
}

func printFuncValue(p MyPoint){
    p.X = 1
    p.Y = 1
    fmt.Printf(" -> %v", p)
}

func printFuncPointer(pp *MyPoint){
    pp.X = 1 // 实际上应该写做 (*pp).X，Golang 给了语法糖，减少了麻烦，但是也导致了 * 的不一致
    pp.Y = 1
    fmt.Printf(" -> %v", pp)
}

func (p MyPoint) printMethodValue(){
    p.X += 1
    p.Y += 1
    fmt.Printf(" -> %v", p)
}

// 建议使用指针作为方法（method：printMethodPointer）的接收者（receiver：*MyPoint），一是可以修改接收者的值，二是可以避免大对象的复制
func (pp *MyPoint) printMethodPointer(){
    pp.X += 1
    pp.Y += 1
    fmt.Printf(" -> %v", pp)
}

func main(){
    p := MyPoint{0, 0}
    pp := &MyPoint{0, 0}

    fmt.Printf("\n value to func(value): %v", p)
    printFuncValue(p)
    fmt.Printf(" --> %v", p)
    // Output: value to func(value): {0 0} -> {1 1} --> {0 0}

    //printFuncValue(pp) // cannot use pp (type *MyPoint) as type MyPoint in argument to printFuncValue

    //printFuncPointer(p) // cannot use p (type MyPoint) as type *MyPoint in argument to printFuncPointer

    fmt.Printf("\n pointer to func(pointer): %v", pp)
    printFuncPointer(pp)
    fmt.Printf(" --> %v", pp)
    // Output: pointer to func(pointer): &{0 0} -> &{1 1} --> &{1 1}

    fmt.Printf("\n value to method(value): %v", p)
    p.printMethodValue()
    fmt.Printf(" --> %v", p)
    // Output: value to method(value): {0 0} -> {1 1} --> {0 0}

    fmt.Printf("\n value to method(pointer): %v", p)
    p.printMethodPointer()
    fmt.Printf(" --> %v", p)
    // Output: value to method(pointer): {0 0} -> &{1 1} --> {1 1}

    fmt.Printf("\n pointer to method(value): %v", pp)
    pp.printMethodValue()
    fmt.Printf(" --> %v", pp)
    // Output: pointer to method(value): &{1 1} -> {2 2} --> &{1 1}

    fmt.Printf("\n pointer to method(pointer): %v", pp)
    pp.printMethodPointer()
    fmt.Printf(" --> %v", pp)
    // Output: pointer to method(pointer): &{1 1} -> &{2 2} --> &{2 2}
}
```

### `slice`
**`slice` 实际上相当于对其依附的 `array` 的引用，它不存储数据，只是对 `array` 进行描述。因此，修改 `slice` 中的元素，
改变会体现在 `array` 上，当然也会体现在该 `array` 的所有 `slice` 上**。

### map

**使用 `make(map[string]string)` 返回的本身是个引用，可以直接用来操作**：
```go
map["name"]="Jason"
```

而**如果使用 `map` 的指针，反而会产生错误**：
```go
*map["name"]="Jason"  //  invalid indirect of m["title"] (type string)
(*map)["name"]="Jason"  // invalid indirect of m (type map[string]string)
```

## 哪些值是不可寻址的
1. **不可变的值不可寻址**。常量、基本类型的值字面量、字符串变量的值、函数以及方法的字面量都是如此。
其实这样规定也有安全性方面的考虑。
2. 绝大多数被视为**临时结果的值都是不可寻址的**。算术操作的结果值属于临时结果，针对值字面量的表达式结果值也属于临时结果。
但有一个例外，对切片字面量的索引结果值虽然也属于临时结果，但却是可寻址的。函数的返回值也是临时结果。`++` 和 `--` 并不属
于操作符。
3. **不安全的值不可寻址**，若拿到某值的指针可能会破坏程序的一致性，那么就是不安全的。由于字典的内部机制，对字典的索
引结果值的取址操作都是不安全的。另外，获取由字面量或标识符代表的函数或方法的地址显然也是不安全的。