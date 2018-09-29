# 数据类型
Go 语言的四类数据类型
- 基础类型，数值、字符串和布尔型
- 复合类型，数组和结构体
- 引用类型，指针、切片、字典、函数、通道
- 接口类型

## 基础类型
### 数值类型
#### 整型
- `uint`，无符号 32 或 64 位整型
- `uint8`，无符号 8 位整型 (0 到 255)
- `uint16`，无符号 16 位整型 (0 到 65535)
- `uint32`，无符号 32 位整型 (0 到 4294967295)
- `uint64`，无符号 64 位整型 (0 到 18446744073709551615)
- `int`，有符号 32 或 64 位整型
- `int8`，有符号 8 位整型 (-128 到 127)
- `int16`，有符号 16 位整型 (-32768 到 32767)
- `int32`，有符号 32 位整型 (-2147483648 到 2147483647)
- `int64`，有符号 64 位整型 (-9223372036854775808 到 9223372036854775807)

`int`和`uint`对应的是 CPU 平台机器的字大小。

#### 浮点数
- `float32`，IEEE-754 32位浮点型数，`math.MaxFloat32`表示`float32`能表示的最大数值，大约是`3.4e38`。
- `float64`，IEEE-754 64位浮点型数，`math.MaxFloat64`表示`float64`能表示的最大数值，大约是`1.8e308`。

#### 复数
- `complex64`，对应`float32`浮点数精度。
- `complex128`，对应`float64`浮点数精度。

内置`complex`函数创建复数。`math/cmplx`包提供了复数处理的许多函数。

#### 其他数值类型
- `byte`，`uint8`的别名，通常用于表示一个`Unicode`码点。
- `rune`，`int32`的别名，一般用于强调数值是一个原始的数据而不是一个小的整数。
- `uintptr`，无符号整型，用于存放一个指针，没有指定具体的`bit`大小。

### 布尔类型
布尔类型的值只有两种：`true`和`false`。

### 字符串
字符串就是一串固定长度的字符连接起来的字符序列，不可改变。Go 的字符串是由单个字节连接起来的。Go 的字符串的字节使用`UTF-8`编码标识`Unicode`文本。

#### 字符串操作
- 内置函数`len`可以获取字符串的长度。
- 可以通过`string[index]`获取某个索引位置的字节值，字符串是不可修改的，不能使用`string[index] = "string2"`这种方式
改变字符串。
- `string[i, l]`获取`string`从第`i`个字节位置开始的`l`个字节，返回一个新的字符串。如：
	```go
	s := "hello, world"
	fmt.Println(s[0:5]) // "hello"

	fmt.Println(s[:5]) // "hello"
  fmt.Println(s[7:]) // "world"
  fmt.Println(s[:])  // "hello, world"
	```
- `+`拼接字符串，如`fmt.Println("goodbye" + s[5:])`输出`"goodbye, world"`。
- 使用`==`和`<`进行字符串比较。

一个原生的字符串面值形式是\`...\`，使用反引号代替双引号。在原生的字符串面值中，没有转义操作；全部的内容都是字面的意思，包含退格和换行。

## 复合数据类型
### 数组
数组是一个由固定长度的指定类型元素组成的序列。数组的长度在编译阶段确定。

声明数组：
```go
var 变量名 [SIZE]类型
```

内置函数`len`获取数组长度。通过下标访问元素：
```go
var a [3]int             // array of 3 integers
fmt.Println(a[0])        // print the first element
fmt.Println(a[len(a)-1]) // print the last element, a[2]
```
默认情况下，数组的每个元素都被初始化为元素类型对应的零值。
初始化数组：
```go
var q [3]int = [3]int{1, 2, 3}
var r [3]int = [3]int{1, 2}
fmt.Println(r[2]) // "0"

var balance = []float32{1000.0, 2.0, 3.4, 7.0, 50.0}
mt.Println(len(balance)) // 5
var balance2 = []float32
fmt.Println(len(balance2)) // type []float32 is not an expression

q := [...]int{1, 2, 3}
fmt.Printf("%T\n", q) // "[3]int"
```
初始化数组中`{}`中的元素个数不能大于`[]`中的数字。
如果`[]`中的`SIZE`，Go 语言会根据元素的个数来设置数组的大小。
上面代码中的`...`省略号，表示数组的长度是根据初始化值的个数来计算。

声明数组`SIZE`是必须的，如果没有，那就是切片了。

#### 二维数组
```go
a = [3][4]int{  
 {0, 1, 2, 3} ,   /*  第一行索引为 0 */
 {4, 5, 6, 7} ,   /*  第二行索引为 1 */
 {8, 9, 10, 11},   /* 第三行索引为 2 */
}
fmt.Printf("a[%d][%d] = %d\n", 2, 3, a[2][3] )
```

`==`和`!=`比较运算符来比较两个数组，只有当两个数组的所有元素都是相等的时候数组才是相等的。

#### 数组传入函数
当调用函数时，函数的形参会被赋值，**所以函数参数变量接收的是一个复制的副本，并不是原始调用的变量。** 但是
这种机制，如果碰到传递一个大数组时，效率较低。这个时候可以显示的传入一个数组指针（其他语言其实是隐式的传递指针）。
```go
func test(ptr *[32]byte) {
  *ptr = [32]byte{}
}
```

### slice
slice的语法和数组很像，由于数组长度是固定的，所以使用`slice`相比数组会更灵活，`slice`是动态的，长度可以增加也可以减少。
还有一点与数组不同，切片不需要说明长度。

定义切片，和定义数组的区别就是不需要指定`SIZE`：
```go
var 变量名 []类型
```
一个`slice`由三个部分构成：指针、长度和容量。长度不能超过容量。
一个切片在未初始化之前默认为`nil`，长度为`0`。

初始化切片：
```go
// 直接初始化切片，[] 表示是切片类型，{1,2,3}初始化值依次是1,2,3.其 cap=len=3
s :=[]int {1,2,3 }

// 初始化切片s,是数组arr的引用
s := arr[:]

// 将arr中从下标startIndex到endIndex-1 下的元素创建为一个新的切片
s := arr[startIndex:endIndex] 

// 缺省endIndex时将表示一直到arr的最后一个元素
s := arr[startIndex:] 

// 缺省startIndex时将表示从arr的第一个元素开始
s := arr[:endIndex]

// 使用 make 函数来创建切片
// len 是数组的长度并且也是切片的初始长度
// capacity 为可选参数, 指定容量
s := make([]int, len, capacity)
```

#### len() 和 cap()
- `len`获取切片长度。
- `cap`计算切片的最大容量

#### append() 和 copy()
- `append`向切片追加新元素
- `copy`拷贝切片

#### 截取切片
```go
/* 创建切片 */
numbers := []int{0,1,2,3,4,5,6,7,8}   

/* 打印子切片从索引1(包含) 到索引4(不包含)*/
fmt.Println("numbers[1:4] ==", numbers[1:4]) // numbers[1:4] == [1 2 3]

/* 默认下限为 0*/
fmt.Println("numbers[:3] ==", numbers[:3]) // numbers[:3] == [0 1 2]

/* 默认上限为 len(s)*/
fmt.Println("numbers[4:] ==", numbers[4:]) // numbers[4:] == [4 5 6 7 8]

numbers1 := make([]int,0,5)

/* 打印子切片从索引  0(包含) 到索引 2(不包含) */
number2 := numbers[:2]
fmt.Printf("len=%d cap=%d slice=%v\n",len(number2),cap(number2),number2) // len=2 cap=9 slice=[0 1]
/* 打印子切片从索引 2(包含) 到索引 5(不包含) */
number3 := numbers[2:5]
fmt.Printf("len=%d cap=%d slice=%v\n",len(number3),cap(number3),number3) // len=3 cap=7 slice=[2 3 4]
```

### map
`map`是一个无序的`key/value`对的集合，使用`hash`表实现的。
定义`map`，使用`map`关键字：
```go
/* 声明变量，默认 map 是 nil */
var 变量名 map[键类型]值类型

/* 使用 make 函数 */
变量名 := make(map[键类型]值类型)

/* 字面值的语法创建 */
变量名 := map[键类型]值类型{
  key1: value1,
	key2: value2,
	...
}
```
一个`map`在未初始化之前默认为`nil`。
通过索引下标`key`来访问`map`中对应的`value`
```go
age, ok := ages["bob"]
if !ok { /* "bob" is not a key in this map; age == 0. */ }
```
`ok`表示操作结果，是一个布尔值。这叫做`ok-idiom`模式，就是在多返回值中返回一个`ok`布尔值，表示是否操作
成功。

#### delet()
`delete`函数删除`map`元素。
```go
delete(mapName, key)
```

#### 遍历
可以使用`for range`遍历`map`：
```go
for key, value := range mapName {
	fmt.Println(mapName[key])
}
```
**`Map`的迭代顺序是不确定的。可以先使用`sort`包排序**。

### 结构体
结构体是由一系列具有相同类型或不同类型的数据构成的数据集合。
结构体定义需要使用`type`和`struct`语句, `struct`语句定义一个新的数据类型, `type` 语句定义了结构体的名称：
```go
// 定义了结构体类型
type struct_variable_type struct {
   member definition;
   member definition;
   ...
   member definition;
}

// 声明
variable_name := structure_variable_type{value1, value2...valuen}
// 或
variable_name := structure_variable_type{ key1: value1, key2: value2..., keyn: valuen}
```

用点号`.`操作符访问结构体成员, 实例：
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
`.`点操作符也可以和指向结构体的指针一起工作:
```go
var employeeOfTheMonth *Employee = &dilbert
employeeOfTheMonth.Position += " (proactive team player)"
```

一个结构体可能同时包含导出和未导出的成员, 如果结构体成员名字是以大写字母开头的，那么该成员就是导出的。
未导出的成员, 不允许在外部包修改。

通常一行对应一个结构体成员，成员的名字在前类型在后，不过如果相邻的成员类型如果相同的话可以被合并到一行:
```go
type Employee struct {
	ID            int
	Name, Address string
	Salary        int
}
```

一个命名为 S 的结构体类型将不能再包含 S 类型的成员：因为一个聚合的值不能包含它自身。（该限制同样适应于数组。）
但是S类型的结构体可以包含 *S 指针类型的成员，这可以让我们创建递归的数据结构，比如链表和树结构等：
```go
type tree struct {
	value       int
	left, right *tree
}
```


#### 结构体字面值

结构体字面值可以指定每个成员的值:
```go
type Point struct{ X, Y int }

p := Point{1, 2}
```

#### 结构体比较

两个结构体将可以使用`==`或`!=`运算符进行比较。
```go
type Point struct{ X, Y int }

p := Point{1, 2}
q := Point{2, 1}
fmt.Println(p.X == q.X && p.Y == q.Y) // "false"
fmt.Println(p == q)                   // "false"
```

#### 结构体嵌入 匿名成员
Go 语言提供的不同寻常的结构体嵌入机制让一个命名的结构体包含另一个结构体类型的匿名成员，
这样就可以通过简单的点运算符`x.f`来访问匿名成员链中嵌套的`x.d.e.f`成员。
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

Go 语言有一个特性可以只声明一个成员对应的数据类型而定义成员的名字；这类成员就叫匿名成员。
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

上面的代码中，`Circle`和`Wheel`各自都有一个匿名成员。我们可以说`Point`类型被嵌入到了`Circle`结构体，同时`Circle`类型被嵌入到了`Wheel`结构体。
但是结构体字面值并没有简短表示匿名成员的语法，所以下面的代码，会编译失败：
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

不能同时包含两个类型相同的匿名成员，这会导致名字冲突。