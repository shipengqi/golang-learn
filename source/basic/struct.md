---
title: 结构体
---
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

### 如果被嵌入类型和嵌入类型有同名的方法，那么调用哪一个的方法
**只要名称相同，无论这两个方法的签名是否一致，被嵌入类型的方法都会“屏蔽”掉嵌入字段的同名方法**。

类似的，由于我们同样可以像访问被嵌入类型的字段那样，直接访问嵌入字段的字段，所以**如果这两个结构体类型里存在同名的字段，
那么嵌入字段中的那个字段一定会被“屏蔽”**。

正因为嵌入字段的字段和方法都可以“嫁接”到被嵌入类型上，所以即使在两个同名的成员一个是字段，另一个是方法的情况下，这种“屏蔽”现象依然会存在。

**不过，即使被屏蔽了，我们仍然可以通过链式的选择表达式，选择到嵌入字段的字段或方法**。

嵌入字段本身也有嵌入字段的情况，这种情况下，“屏蔽”现象会以嵌入的层级为依据，嵌入层级越深的字段或方法越可能被“屏蔽”。