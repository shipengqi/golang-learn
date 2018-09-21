# 面向对象
GO 支持面向对象编程。

## 方法

方法声明：
```go
func (变量名 类型) 方法() [返回类型]{
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
// 这里的 c 叫作方法的接收器
func (c Circle) getArea() float64 {
  // c.radius 即为 Circle 类型对象中的属性
  return 3.14 * c.radius * c.radius
}
```

Go 没有像其它语言那样用`this`或者`self`作为接收器。Go 可以给任意类型定义方法。

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
2. 在声明一个`method`的`receiver`该是指针还是非指针类型时，你需要考虑两方面的内部，第一方面是这个对象本身是不是特别大，如果声明为非指针变量时，调用会产生一次拷贝；第二方面是如果你用指针类型作为`receiver`，那么你一定要注意，这种指针类型指向的始终是一块内存地址，就算你对其进行了拷贝。