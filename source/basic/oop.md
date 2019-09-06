---
title: 面向对象
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
  - 包含 `Mutex` 等同步字段，使用 `*T`，避免因为复制造成锁操作无效。
  - 无法确定时，使用 `*T`。

**方法的接收者类型必须是某个自定义的数据类型，而且不能是接口类型或接口的指针类型**。
- 值方法，就是接收者类型是非指针的自定义数据类型的方法。
- 指针方法，就是接收者类型是指针类型的方法。

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