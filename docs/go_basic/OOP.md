# 面向对象
GO 支持面向对象编程。

## 方法

方法声明：
```go
func (变量名 类型) 方法名() [返回类型]{
   /* 函数体*/
}
```
实例：
```go
/* 定义结构体 */
type Circle struct {
  radius float64
}

func main() {
  var c1 Circle
  c1.radius = 10.00
  fmt.Println("Area of Circle(c1) = ", c1.getArea())
}

// 该 method 属于 Circle 类型对象中的方法
// 这里的 c 叫作方法的接收器，类似 Javascript 的 this
func (c Circle) getArea() float64 {
  // c.radius 即为 Circle 类型对象中的属性
  return 3.14 * c.radius * c.radius
}
```

Go 没有像其它语言那样用`this`或者`self`作为接收器。**Go 可以给任意类型定义方法。

当调用一个函数时，会对其每一个参数值进行拷贝，**如果一个函数需要更新一个变量，或者函数的其中一个参数实在太大我们希望能够避免进行这种默认的拷贝，**
**我们可以传入变量的指针。**
```go
func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}
```
调用指针类型方法`(*Point).ScaleBy`，`()`必须有，否则会被理解为`*(Point.ScaleBy)`。
```go
// 调用指针类型方法
r := &Point{1, 2}
r.ScaleBy(2)

// 简短写法
p := Point{1, 2}

// 编译器会隐式地帮我们用&p去调用ScaleBy这个方法。这种简写方法只适用于“变量”
p.ScaleBy(2)
```

只有类型(`Point`)和指向他们的指针(`*Point`)，才是可能会出现在接收器声明里的两种接收器。此外，为了避免歧义，在声明方法时，如果一个类型名本身是一个指针的话，是不允许其出现在接收器中的:
```go
type P *int
func (P) f() { /* ... */ } // compile error: invalid receiver type
```
注意两点：
1. 不管你的`method`的`receiver`是指针类型还是非指针类型，都是可以通过指针/非指针类型进行调用的，编译器会帮你做类型转换。
2. 在声明一个`method`的`receiver`该是指针还是非指针类型时，你需要考虑两方面的内部，第一方面是这个对象本身是不是特别大，如果声明为非指针变量时，调用会产生一次拷贝；
第二方面是如果你用指针类型作为`receiver`，那么你一定要注意，这种指针类型指向的始终是一块内存地址，就算你对其进行了拷贝。

**方法的接收者类型必须是某个自定义的数据类型，而且不能是接口类型或接口的指针类型**。
- 值方法，就是接收者类型是非指针的自定义数据类型的方法。
- 指针方法，就是接收者类型是指针类型的方法。

### 嵌入结构体扩展类型
```go
import "image/color"

type Point struct{ X, Y float64 }

type ColoredPoint struct {
  Point
  Color color.RGBA
}

red := color.RGBA{255, 0, 0, 255}
blue := color.RGBA{0, 0, 255, 255}
var p = ColoredPoint{Point{1, 1}, red}
var q = ColoredPoint{Point{5, 4}, blue}
fmt.Println(p.Distance(q.Point)) // "5"
p.ScaleBy(2)
q.ScaleBy(2)
fmt.Println(p.Distance(q.Point)) // "10"
```

如果对基于类来实现面向对象的语言比较熟悉的话，可能会倾向于将`Point`看作一个基类，而`ColoredPoint`看作其子类或者继承类。
但这是错误的理解。请注意上面例子中对`Distance`方法的调用。`Distance`有一个参数是`Point`类型，但是这里的`q`虽然貌似是继承了
`Point`类，但`q`并不是，所以尽管`q`有着`Point`这个内嵌类型，我们也必须要显式传入`q.Point`。

#### Go 语言是用嵌入字段实现了继承吗
Go 语言中**没有继承的概念，它所做的是通过嵌入字段的方式实现了类型之间的组合**。具体原因和理念请见[Why is there no type inheritance?](https://golang.org/doc/faq#inheritance)。

简单来说，面向对象编程中的继承，其实是通过牺牲一定的代码简洁性来换取可扩展性，而且这种可扩展性是通过侵入的方式来实现的。类型之间的组合采用的是非声明的方式，
我们不需要显式地声明某个类型实现了某个接口，或者一个类型继承了另一个类型。

同时，类型组合也是非侵入式的，它不会破坏类型的封装或加重类型之间的耦合。我们要做的只是把类型当做字段嵌入进来，然后坐享其成地使用嵌入字段所拥有的一切。
如果嵌入字段有哪里不合心意，我们还可以用“包装”或“屏蔽”的方式去调整和优化。

另外，类型间的组合也是灵活的，我们总是可以通过嵌入字段的方式把一个类型的属性和能力“嫁接”给另一个类型。

这时候，被嵌入类型也就自然而然地实现了嵌入字段所实现的接口。再者，组合要比继承更加简洁和清晰，Go 语言可以轻而易举地通过嵌入多个字段来实现功能强大的类型，
却不会有多重继承那样复杂的层次结构和可观的管理成本。

### 封装
一个对象的变量或者方法如果对调用方是不可见的话，一般就被定义为“封装”。通过首字母大小写来定义是否从包中导出。
封装一个对象，必须定义为一个`struct`：
```go
type IntSet struct {
  words []uint64
}
```

优点：
- 调用方不能直接修改对象的变量值
- 隐藏实现的细节，防止调用方依赖那些可能变化的具体实现，这样使设计包的程序员在不破坏对外的api情况下能得到更大的自由。
- 阻止了外部调用方对对象内部的值任意地进行修改。

### `String`方法
在 Go 语言中，**我们可以通过为一个类型编写名为`String`的方法，来自定义该类型的字符串表示形式。这个`String`方法不需要任何参数声明，但需要有一个`string`类型的结果声明**。
```go
type AnimalCategory struct {
    kingdom string // 界。
    phylum string // 门。
    class  string // 纲。
    order  string // 目。
    family string // 科。
    genus  string // 属。
    species string // 种。
}

func (ac AnimalCategory) String() string {
    return fmt.Sprintf("%s%s%s%s%s%s%s",ac.kingdom, ac.phylum, ac.class, ac.order,ac.family, ac.genus, ac.species)
}

category := AnimalCategory{species: "cat"}
fmt.Printf("The animal category: %s\n", category)

```
正因为如此，我在调用`fmt.Printf`函数时，使用占位符`%s`和`category`值本身就可以打印出后者的字符串表示形式，而**无需显式地调用它的`String`方法**。

`fmt.Printf`函数会自己去寻找它。此时的打印内容会是`The animal category: cat`。显而易见，`category`的`String`方法成功地引用了当前值的所有字段。

## 接口

Go 支持接口数据类型，接口类型是一种抽象的类型。接口类型具体描述了一系列方法的集合，任何其他类型只要实现了这些方法就是实现了这个接口。
接口只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要。

接口的零值就是它的类型和值的部分都是`nil`。

简单的说，`interface`是一组`method`的组合，我们通过`interface`来定义对象的一组行为。

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
`interface {}`被称为空接口类型，它没有任何方法，类似 Javascrit 的`Object`。所有的类型都实现了空`interface`，
空`interface`在我们需要存储任意类型的数值的时候相当有用，因为它可以存储任意类型的数值。
```go
// 定义a为空接口
var a interface{}
var i int = 5
s := "Hello world"
// a可以存储任意类型的数值
a = i
a = s
```
一个函数把`interface{}`作为参数，那么他可以接受任意类型的值作为参数，如果一个函数返回`interface{}`,那么也就可以返回任意类型的值。

`interface{}`可以存储任意类型，那么怎么判断存储了什么类型？

#### 类型断言

Go语言里面有一个语法，可以直接判断是否是该类型的变量： `value, ok = element.(T)`，这里`value`就是变量的值，`ok`是一个`bool`类型，`element`是`interface`变量，`T`是断言的类型。
```go
// comma-ok
for index, element := range list {
	if value, ok := element.(int); ok {
		fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
	} else if value, ok := element.(string); ok {
		fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
	} else if value, ok := element.(Person); ok {
		fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
	} else {
		fmt.Printf("list[%d] is of a different type\n", index)
	}
}


// 或者 使用 switch
for index, element := range list{
	switch value := element.(type) {
		case int:
			fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
		case string:
			fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
		case Person:
			fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
		default:
			fmt.Println("list[%d] is of a different type", index)
	}
}
```

**注意，`element.(type)`语法不能在`switch`外的任何逻辑里面使用，如果你要在`switch`外面判断一个类型就使用`comma-ok`**。

### error 接口
Go 内置了错误接口。
```go
type error interface {
  Error() string
}
```
创建一个`error`最简单的方法就是调用`errors.New`函数。

`error`包：
```go
package errors

func New(text string) error { return &errorString{text} }

type errorString struct { text string }

func (e *errorString) Error() string { return e.text }
```

`fmt.Errorf`封装了`errors.New`函数，它会处理字符串格式化。

### 接口的实际用途
```go
package main

import (
    "fmt"
)

//定义interface
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

上面的代码`fmt.Printf("Vowels are %c", v.FindVowels())`是可以直接使用`fmt.Printf("Vowels are %c", name.FindVowels())`的，
那么我们定义的变量V没有没有了意义。看下面的代码：
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

这里作为一个js开发，理解不了接口的作用，因为js是弱类型语言，go是强类型语言，像上面的数组`employees`，js可以直接塞入实现了`CalculateSalary`方法的类。因为js数组的元素没有类型限制，
可以塞入不同的类型。但是go不可以，所以这个时候体现出了接口的作用，`Contract`和`Permanent`是不一样的结构体类型，但是可以定义一个`SalaryCalculator`接口类型的数组，
就可以在`totalExpense`中调用元素的`CalculateSalary`方法。否则就要分别调用`pemp1.SalaryCalculator()`，`pemp2.SalaryCalculator()`，`cemp1.SalaryCalculator()`。