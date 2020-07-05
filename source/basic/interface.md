---
title: 接口
---


Go 支持接口数据类型，接口类型是一种抽象的类型。接口类型具体描述了一系列方法的集合，任何其他类型只要实现了这些方法就是实
现了这个接口，无须显示声明。**接口只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要**。

**一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口**。

接口的零值就是它的类型和值的部分都是 `nil`。

简单的说，`interface` 是一组 `method` 的组合，我们通过 `interface` 来定义对象的一组行为。

定义接口：

```go
type 接口名 interface {
  方法名1 [返回类型]
  方法名2 [返回类型]
  方法名3 [返回类型]
  ...
}

/* 定义结构体 */
type struct_name struct {
   /* variables */
}

/* 实现接口方法 */
func (struct_name_variable struct_name) 方法名1() [返回类型] {
   /* 方法实现 */
}
...
func (struct_name_variable struct_name) 方法名2() [返回类型] {
   /* 方法实现*/
}
```

实例：

```go
type Phone interface {
  call()
}

type NokiaPhone struct {
}

func (nokiaPhone NokiaPhone) call() {
  fmt.Println("I am Nokia, I can call you!")
}

type IPhone struct {
}

func (iPhone IPhone) call() {
  fmt.Println("I am iPhone, I can call you!")
}

func main() {
  var phone Phone

  phone = new(NokiaPhone)
  phone.call()

  phone = new(IPhone)
  phone.call()
}
```

接口类型也可以通过组合已有的接口来定义：

```go
type Reader interface {
  Read(p []byte) (n int, err error)
}
type Closer interface {
  Close() error
}


type ReadWriteCloser interface {
  Reader
  Writer
  Closer
}

// 混合
type ReadWriter interface {
  Read(p []byte) (n int, err error)
  Writer
}
```

### 空接口类型

`interface {}` 被称为空接口类型，它没有任何方法。所有的类型都实现了空 `interface`，
空 `interface` 在我们需要存储任意类型的数值的时候相当有用，因为它可以存储任意类型的数值。

```go
// 定义a为空接口
var a interface{}
var i int = 5
s := "Hello world"
// a可以存储任意类型的数值
a = i
a = s
```

一个函数把 `interface{}` 作为参数，那么他可以接受任意类型的值作为参数，如果一个函数返回 `interface{}`,
那么也就可以返回任意类型的值。

`interface{}` 可以存储任意类型，那么怎么判断存储了什么类型？

### error 接口

Go 内置了错误接口。

```go
type error interface {
  Error() string
}
```

创建一个 `error` 最简单的方法就是调用 `errors.New` 函数。

`error`包：

```go
package errors

func New(text string) error { return &errorString{text} }

type errorString struct { text string }

func (e *errorString) Error() string { return e.text }
```

`fmt.Errorf` 封装了 `errors.New` 函数，它会处理字符串格式化。**当我们想通过模板化的方式生成错误信息，并得到错误值时，
可以使用`fmt.Errorf`函数。该函数所做的其实就是先调用 `fmt.Sprintf` 函数，得到确切的错误信息；再调用 `errors.New` 函数，
得到包含该错误信息的 `error` 类型值，最后返回该值**。

实际上，`error` 类型值的 `Error` 方法就相当于其他类型值的 `String` 方法。

### 接口的实际用途

```go
package main

import (
    "fmt"
)

//定义 interface
type VowelsFinder interface {
    FindVowels() []rune
}

type MyString string

//实现接口
func (ms MyString) FindVowels() []rune {
    var vowels []rune
    for _, rune := range ms {
        if rune == 'a' || rune == 'e' || rune == 'i' || rune == 'o' || rune == 'u' {
            vowels = append(vowels, rune)
        }
    }
    return vowels
}

func main() {
    name := MyString("Sam Anderson") // 类型转换
    var v VowelsFinder // 定义一个接口类型的变量
    v = name
    fmt.Printf("Vowels are %c", v.FindVowels())

}
```

上面的代码 `fmt.Printf("Vowels are %c", v.FindVowels())` 是可以直接使用 `fmt.Printf("Vowels are %c", name.FindVowels())`
的，那么我们定义的变量 `V` 没有没有了意义。看下面的代码：

```go
package main

import (
 "fmt"
)

// 薪资计算器接口
type SalaryCalculator interface {
 CalculateSalary() int
}
// 普通挖掘机员工
type Contract struct {
 empId  int
 basicpay int
}
// 有蓝翔技校证的员工
type Permanent struct {
 empId  int
 basicpay int
 jj int // 奖金
}

func (p Permanent) CalculateSalary() int {
 return p.basicpay + p.jj
}

func (c Contract) CalculateSalary() int {
 return c.basicpay
}
// 总开支
func totalExpense(s []SalaryCalculator) {
 expense := 0
 for _, v := range s {
  expense = expense + v.CalculateSalary()
 }
 fmt.Printf("总开支 $%d", expense)
}

func main() {
 pemp1 := Permanent{1,3000,10000}
 pemp2 := Permanent{2, 3000, 20000}
 cemp1 := Contract{3, 3000}
 employees := []SalaryCalculator{pemp1, pemp2, cemp1}
 totalExpense(employees)
}
```

这个时候体现出了接口的作用，`Contract` 和 `Permanent` 是不一样的结构体类型，但是可以定义一个 `SalaryCalculator` 接口类
型的数组，就可以在 `totalExpense` 中调用元素的 `CalculateSalary` 方法。
