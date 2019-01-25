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
上面代码中的**`...`省略号，表示数组的长度是根据初始化值的个数来计算**。

**声明数组`SIZE`是必须的，如果没有，那就是切片了。**

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
这种机制，**如果碰到传递一个大数组时，效率较低。这个时候可以显示的传入一个数组指针**（其他语言其实是隐式的传递指针）。
```go
func test(ptr *[32]byte) {
  *ptr = [32]byte{}
}
```

### slice
slice的语法和数组很像，由于数组长度是固定的，所以使用`slice`相比数组会更灵活，`slice`是动态的，长度可以增加也可以减少。
还有一点与数组不同，切片不需要说明长度。

**定义切片，和定义数组的区别就是不需要指定`SIZE`**：
```go
var 变量名 []类型
```
一个`slice`由三个部分构成：指针、长度和容量。长度不能超过容量。
一个切片在未初始化之前默认为`nil`，长度为`0`。

初始化切片：
```go
// 直接初始化切片，[] 表示是切片类型，{1,2,3}初始化值依次是1,2,3.其 cap=len=3
s :=[]int {1,2,3}

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

#### 怎样估算切片容量的增长

一旦一个切片无法容纳更多的元素，Go 语言就会想办法扩容。但它并不会改变原来的切片，而是会生成一个容量更大的切片，然后将把原有的元素和新元素一并拷贝到新切片中。
般的情况下，你可以简单地认为新切片的容量（以下简称新容量）将会是原切片容量（以下简称原容量）的 2 倍。

但是，当原切片的长度（以下简称原长度）大于或等于1024时，Go 语言将会以原容量的1.25倍作为新容量的基准（以下新容量基准）。新容量基准会被调整（不断地与1.25相乘），
直到结果不小于原长度与要追加的元素数量之和（以下简称新长度）。最终，新容量往往会比新长度大一些，当然，相等也是可能的。

一个切片的底层数组永远不会被替换。为什么？虽然在扩容的时候 Go 语言一定会生成新的底层数组，但是它也同时生成了新的切片。它是把新的切片作为了新底层数组的窗口，
而没有对原切片及其底层数组做任何改动。

在无需扩容时，append函数返回的是指向原底层数组的新切片，而在需要扩容时，`append`函数返回的是指向新底层数组的新切片。

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

使用`map`过程中需要注意的几点：
- `map`是无序的，每次打印出来的`map`都会不一样，它不能通过`index`获取，而必须通过`key`获取
- `map`的长度是不固定的，也就是和`slice`一样，也是一种引用类型
- 内置的`len`函数同样适用于`map`，返回`map`拥有的`key`的数量
- `map`的值可以很方便的修改，通过`numbers["one"]=11`可以很容易的把`key`为`one`的字典值改为11
- `map`和其他基本型别不同，它不是`thread-safe`，在多个`go-routine`存取时，必须使用`mutex lock`机制

#### delete()
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

#### map的键类型不能是哪些类型
map的键和元素的最大不同在于，前者的类型是受限的，而后者却可以是任意类型的。

**map的键类型不可以是函数类型、字典类型和切片类型**。

为什么？

Go 语言规范规定，在键类型的值之间必须可以施加操作符`==`和`!=`。换句话说，键类型的值必须要支持判等操作。由于函数类型、字典类型和切片类型的值并不支持判等操作，所以字典的键类型不能是
这些类型。

另外，如果键的类型是接口类型的，那么键值的实际类型也不能是上述三种类型，否则在程序运行过程中会引发 panic（即运行时恐慌）。
```go
var badMap2 = map[interface{}]int{
"1":   1,
[]int{2}: 2, // 这里会引发 panic。
3:    3,
}
```

#### 优先考虑哪些类型作为字典的键类型
求哈希和判等操作的速度越快，对应的类型就越适合作为键类型。

对于所有的基本类型、指针类型，以及数组类型、结构体类型和接口类型，Go 语言都有一套算法与之对应。这套算法中就包含了哈希和判等。以求哈希的操作为例，宽度越小的类型速度通常越快。
对于布尔类型、整数类型、浮点数类型、复数类型和指针类型来说都是如此。对于字符串类型，由于它的宽度是不定的，所以要看它的值的具体长度，长度越短求哈希越快。

类型的宽度是指它的单个值需要占用的字节数。比如，bool、int8和uint8类型的一个值需要占用的字节数都是1，因此这些类型的宽度就都是1。


#### 在值为nil的字典上执行读写操作会成功吗
当我们仅声明而不初始化一个字典类型的变量的时候，它的值会是`nil`。

**除了添加键 - 元素对，我们在一个值为nil的字典上做任何操作都不会引起错误**。当我们试图在一个值为nil的字典中添加键 - 元素对的时候，Go 语言的运行时系统就会立即抛出一个 panic。
可以先使用`make`函数初始化。

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

##### 如果被嵌入类型和嵌入类型有同名的方法，那么调用哪一个的方法
**只要名称相同，无论这两个方法的签名是否一致，被嵌入类型的方法都会“屏蔽”掉嵌入字段的同名方法**。

类似的，由于我们同样可以像访问被嵌入类型的字段那样，直接访问嵌入字段的字段，所以如果这两个结构体类型里存在同名的字段，那么嵌入字段中的那个字段一定会被“屏蔽”。

正因为嵌入字段的字段和方法都可以“嫁接”到被嵌入类型上，所以即使在两个同名的成员一个是字段，另一个是方法的情况下，这种“屏蔽”现象依然会存在。

**不过，即使被屏蔽了，我们仍然可以通过链式的选择表达式，选择到嵌入字段的字段或方法**。

嵌入字段本身也有嵌入字段的情况，这种情况下，“屏蔽”现象会以嵌入的层级为依据，嵌入层级越深的字段或方法越可能被“屏蔽”。

## 类型转换
Go 强制使用显示类型转换。这样可以确定语句和表达式的明确含义。

```go
a := 100
b := byte(a)
c := a + int(b) 混合类型表达式，类型必须保持一致
```
在 Go 中，非布尔值不能当做`true/false`使用，这点和我常用的js不同：
```go
x := 100

if x { // 错误 x 不是布尔值

}
```

如果要转换为指针类型，或者单向`channel`，或者函数，要给类型加上`()`，避免编译器分析错误，如：
```go
x := 100
(*int)(&x) // *int 加括号，否则会被解析为*(int(&x))

(<- channel int)(c)
(func())(f)
(func()int)(f) // 有返回值的函数其实可以不加括号，但是加括号的话，语义清晰
```

## 零值
“零值”，所指并非是空值，而是一种“变量未填充前”的默认值，通常为`0`：
```
int     0
int8    0
int32   0
int64   0
uint    0x0
rune    0 //rune的实际类型是 int32
byte    0x0 // byte的实际类型是 uint8
float32 0 //长度为 4 byte
float64 0 //长度为 8 byte
bool    false
string  ""
```

## container包
Go 语言的链表实现在标准库的`container/list`代码包中。

这个代码包中有两个公开的程序实体——`List`和`Element`，`List`实现了一个双向链表（以下简称链表），而`Element`则代表了链表中元素的结构。

List的四种方法:
- `MoveBefore`方法和`MoveAfter`方法，它们分别用于把给定的元素移动到另一个元素的前面和后面。
- `MoveToFront`方法和`MoveToBack`方法，分别用于把给定的元素移动到链表的最前端和最后端。


```go
func (l *List) MoveBefore(e, mark *Element)
func (l *List) MoveAfter(e, mark *Element)

func (l *List) MoveToFront(e *Element)
func (l *List) MoveToBack(e *Element)
```
“给定的元素”都是*Element类型。

如果我们自己生成这样的值，然后把它作为“给定的元素”传给链表的方法，那么会发生什么？链表会接受它吗？

不会接受，这些方法将不会对链表做出任何改动。因为我们自己生成的Element值并不在链表中，所以也就谈不上“在链表中移动元素”。

- InsertBefore和InsertAfter方法分别用于在指定的元素之前和之后插入新元素
- PushFront和PushBack方法则分别用于在链表的最前端和最后端插入新元素。