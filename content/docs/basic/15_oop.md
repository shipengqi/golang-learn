---
title: 面向对象
weight: 15
---

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

Go 没有像其它语言那样用 `this` 或者 `self` 作为接收器。**Go 可以给任意类型定义方法**。

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

只有类型(`Point`)和指向他们的指针(`*Point`)，才是可能会出现在接收器声明里的两种接收器。此外，为了避免歧义，在声明方法时，
如果一个类型名本身是一个指针的话，是不允许其出现在接收器中的:

```go
type P *int
func (P) f() { /* ... */ } // compile error: invalid receiver type
```

### 如何选择 receiver 的类型

1. **不管你的 `method` 的 `receiver` 是指针类型还是非指针类型，都是可以通过指针/非指针类型进行调用的，编译器会帮你做类型转换**。
2. 在声明一个 `method` 的 `receiver` 该是指针还是非指针类型时，你需要考虑：

- 要修改实例状态，用 `*T`，无需修改使用 `T`。
- 大对象建议使用 `*T`，减少复制成本，`T` 调用时会产生一次拷贝。
- 对于引用类型，直接使用 `T`，因为它们本身就是指针包装的。
- 包含 `Mutex` 等并发原语的，使用 `*T`，避免因为复制造成锁操作无效。
- 无法确定时，使用 `*T`。

**方法的接收者类型必须是某个自定义的数据类型，而且不能是接口类型或接口的指针类型**。

- 值方法，就是接收者类型是非指针的自定义数据类型的方法。
- 指针方法，就是接收者类型是指针类型的方法。

#### 实现了 interface 的方法

如果一个类型实现的某个接口的方法，如果接收者是指针类型，那么只能指针赋值：

```go
type I interface {
 Get()
}
type S struct {
}

func (s *S) Get() {
 fmt.Println("get")
}

func main() {
 ss := S{}

 var i I
 //i = ss , 此处编译不过
 //i.Get()

 i = &ss // 必须是指针赋值
 i.Get()
}
```

如果接收者是非指针类型，那么值和指针都可以赋值：

```go
 ss := S{}

 var i I
 i = ss  // 可以赋值
 i.Get()

 i = &ss // 可以赋值
 i.Get()
```

### 方法集

Golang 方法集 ：每个类型都有与之关联的方法集，这会影响到接口实现规则。

```
• 类型 T 方法集包含全部 receiver T 方法。
• 类型 *T 方法集包含全部 receiver T + *T 方法。

• 如类型 S 包含匿名字段 T，则 S 和 *S 方法集包含 T 方法。
• 如类型 S 包含匿名字段 *T，则 S 和 *S 方法集包含 T + *T 方法。
• 不管嵌入 T 或 *T，*S 方法集总是包含 T + *T 方法。
```

**对于结构体嵌套匿名字段的类型是指针还是非指针**，根据实际情况决定。

## 嵌入结构体扩展类型

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

如果对基于类来实现面向对象的语言比较熟悉的话，可能会倾向于将 `Point` 看作一个基类，而 `ColoredPoint` 看作其子类或者继承类。
但这是错误的理解。请注意上面例子中对 `Distance` 方法的调用。`Distance` 有一个参数是 `Point` 类型，但是这里的 `q` 虽然貌
似是继承了`Point` 类，但 `q` 并不是，所以尽管 `q` 有着 `Point` 这个内嵌类型，我们也必须要显式传入 `q.Point`。

### Go 语言是用嵌入字段实现了继承吗

Go 语言中**没有继承的概念，它所做的是通过嵌入字段的方式实现了类型之间的组合**。
具体原因和理念请见 [Why is there no type inheritance?](https://golang.org/doc/faq#inheritance)。

简单来说，面向对象编程中的继承，其实是通过牺牲一定的代码简洁性来换取可扩展性，而且这种可扩展性是通过侵入的方式来实现的。
类型之间的组合采用的是非声明的方式，我们不需要显式地声明某个类型实现了某个接口，或者一个类型继承了另一个类型。

同时，类型组合也是非侵入式的，它不会破坏类型的封装或加重类型之间的耦合。我们要做的只是把类型当做字段嵌入进来，然后坐
享其成地使用嵌入字段所拥有的一切。如果嵌入字段有哪里不合心意，我们还可以用“包装”或“屏蔽”的方式去调整和优化。

另外，类型间的组合也是灵活的，我们总是可以通过嵌入字段的方式把一个类型的属性和能力“嫁接”给另一个类型。

这时候，被嵌入类型也就自然而然地实现了嵌入字段所实现的接口。再者，组合要比继承更加简洁和清晰，Go 语言可以轻而易举地通过嵌入
多个字段来实现功能强大的类型，却不会有多重继承那样复杂的层次结构和可观的管理成本。

## 封装

一个对象的变量或者方法如果对调用方是不可见的话，一般就被定义为“封装”。通过首字母大小写来定义是否从包中导出。
封装一个对象，必须定义为一个 `struct`：

```go
type IntSet struct {
  words []uint64
}
```

优点：

- 调用方不能直接修改对象的变量值
- 隐藏实现的细节，防止调用方依赖那些可能变化的具体实现，这样使设计包的程序员在不破坏对外的api情况下能得到更大的自由。
- 阻止了外部调用方对对象内部的值任意地进行修改。

## `String` 方法

在 Go 语言中，**我们可以通过为一个类型编写名为 `String` 的方法，来自定义该类型的字符串表示形式。这个 `String` 方法不需
要任何参数声明，但需要有一个 `string` 类型的结果声明**。

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

正因为如此，我在调用 `fmt.Printf` 函数时，使用占位符 `%s` 和 `category` 值本身就可以打印出后者的字符串表示形式，
而**无需显式地调用它的 `String` 方法**。

`fmt.Printf` 函数会自己去寻找它。此时的打印内容会是 `The animal category: cat`。显而易见，`category` 的 `String` 方法成
功地引用了当前值的所有字段。

当你广泛使用一个自定义类型时，最好为它定义 `String()` 方法。

**不要在 `String()` 方法里面调用涉及 `String()` 方法的方法，它会导致意料之外的错误**，比如：

```go
type TT float64

func (t TT) String() string {
    return fmt.Sprintf("%v", t)
}
t.String()
```

它导致了一个无限递归调用（`TT.String()` 调用 `fmt.Sprintf`，而 `fmt.Sprintf` 又会反过来调用 `TT.String()`...），很快就会导
致内存溢出。


# 结构体

结构体是由一系列具有相同类型或不同类型的数据构成的数据集合。
结构体定义需要使用 `type` 和 `struct` 语句, `struct` 语句定义一个新的数据类型, `type` 语句定义了结构体的名称：

```go
// 定义了结构体类型
type struct_variable_type struct {
   member definition;
   member definition;
   ...
   member definition;
}

variable_name := structure_variable_type{value1, value2...valuen}
// 或
variable_name := structure_variable_type{ key1: value1, key2: value2..., keyn: valuen}
```

用点号 `.` 操作符访问结构体成员, 实例：

```go
type Books struct {
 title string
 author string
 subject string
 book_id int
}


func main() {
 var Book1 Books        /* 声明 Book1 为 Books 类型 */

 /* book 1 描述 */
 Book1.title = "Go 语言"
 Book1.author = "www.runoob.com"
 Book1.subject = "Go 语言教程"
 Book1.book_id = 6495407

  /* 打印 Book1 信息 */
 fmt.Printf( "Book 1 title : %s\n", Book1.title)
 fmt.Printf( "Book 1 author : %s\n", Book1.author)
 fmt.Printf( "Book 1 subject : %s\n", Book1.subject)
 fmt.Printf( "Book 1 book_id : %d\n", Book1.book_id)
}
```

`.` 点操作符也可以和指向结构体的指针一起工作:

```go
var employeeOfTheMonth *Employee = &dilbert
employeeOfTheMonth.Position += " (proactive team player)"
```

**一个结构体可能同时包含导出和未导出的成员, 如果结构体成员名字是以大写字母开头的，那么该成员就是导出的。
未导出的成员, 不允许在外部包修改。**

通常一行对应一个结构体成员，成员的名字在前类型在后，不过如果相邻的成员类型如果相同的话可以被合并到一行:

```go
type Employee struct {
 ID            int
 Name, Address string
 Salary        int
}
```

一个命名为 S 的结构体类型将不能再包含 S 类型的成员：因为一个聚合的值不能包含它自身。（该限制同样适应于数组。）
但是S类型的结构体可以包含 `*S` 指针类型的成员，这可以让我们创建递归的数据结构，比如链表和树结构等：

```go
type tree struct {
 value       int
 left, right *tree
}
```

## 结构体的零值

```go
type Person struct {
  AgeYears int
  Name string
  Friends []Person
}

var p Person // Person{0, "", nil}
```

变量 `p` 只声明但没有赋值，所以 `p` 的所有字段都有对应的零值。

**注意如果声明结构体指针使用 `var p *Person` 的方式，那么 `p` 只是一个 `nil` 指针，建议使用 `p := &Person{}` 的方式声明，
`p` 的值是 `&Person{0, "", nil}`，避免 json unmarshal 出错**。

## 结构体字面值

结构体字面值可以指定每个成员的值:

```go
type Point struct{ X, Y int }

p := Point{1, 2}
```

## 结构体比较

两个结构体将可以使用 `==` 或 `!=` 运算符进行比较。

```go
type Point struct{ X, Y int }

p := Point{1, 2}
q := Point{2, 1}
fmt.Println(p.X == q.X && p.Y == q.Y) // "false"
fmt.Println(p == q)                   // "false"
```

## 结构体嵌入 匿名成员

Go 语言提供的不同寻常的结构体嵌入机制让一个命名的结构体包含另一个结构体类型的匿名成员，
这样就可以通过简单的点运算符 `x.f` 来访问匿名成员链中嵌套的 `x.d.e.f` 成员。

```go
type Point struct {
    X, Y int
}

type Circle struct {
    Center Point
    Radius int
}

type Wheel struct {
    Circle Circle
    Spokes int
}
```

上面的代码，会使访问每个成员变得繁琐：

```go
var w Wheel
w.Circle.Center.X = 8
w.Circle.Center.Y = 8
w.Circle.Radius = 5
w.Spokes = 20
```

Go 语言有一个特性可以**只声明一个成员对应的数据类型而定义成员的名字；这类成员就叫匿名成员**。Go 语言规范规定，
如果一个字段的声明中只有字段的类型名而没有字段的名称，那么它就是一个嵌入字段，也可以被称为匿名字段。
匿名成员的数据类型必须是命名的类型或指向一个命名的类型的指针。

```go
type Point struct {
    X, Y int
}


type Circle struct {
    Point
    Radius int
}

type Wheel struct {
    Circle
    Spokes int
}

var w Wheel
w.X = 8            // equivalent to w.Circle.Point.X = 8
w.Y = 8            // equivalent to w.Circle.Point.Y = 8
w.Radius = 5       // equivalent to w.Circle.Radius = 5
w.Spokes = 20
```

上面的代码中，`Circle` 和 `Wheel` 各自都有一个匿名成员。我们可以说 `Point` 类型被嵌入到了 `Circle` 结构体，
同时 `Circle` 类型被嵌入到了 `Wheel` 结构体。但是**结构体字面值并没有简短表示匿名成员的语法**，所以下面的代码，
会编译失败：

```go
w = Wheel{8, 8, 5, 20}                       // compile error: unknown fields
w = Wheel{X: 8, Y: 8, Radius: 5, Spokes: 20} // compile error: unknown fields

// 正确的语法
w = Wheel{Circle{Point{8, 8}, 5}, 20}

w = Wheel{
    Circle: Circle{
        Point:  Point{X: 8, Y: 8},
        Radius: 5,
    },
    Spokes: 20, // NOTE: trailing comma necessary here (and at Radius)
}
```

**不能同时包含两个类型相同的匿名成员，这会导致名字冲突**。

### 嵌入接口类型

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

### 如果被嵌入类型和嵌入类型有同名的方法，那么调用哪一个的方法

**只要名称相同，无论这两个方法的签名是否一致，被嵌入类型的方法都会“屏蔽”掉嵌入字段的同名方法**。

类似的，由于我们同样可以像访问被嵌入类型的字段那样，直接访问嵌入字段的字段，所以**如果这两个结构体类型里存在同名的字段，
那么嵌入字段中的那个字段一定会被“屏蔽”**。

正因为嵌入字段的字段和方法都可以“嫁接”到被嵌入类型上，所以即使在两个同名的成员一个是字段，另一个是方法的情况下，这种“屏蔽”现象依然会存在。

**不过，即使被屏蔽了，我们仍然可以通过链式的选择表达式，选择到嵌入字段的字段或方法**。

嵌入字段本身也有嵌入字段的情况，这种情况下，“屏蔽”现象会以嵌入的层级为依据，嵌入层级越深的字段或方法越可能被“屏蔽”。