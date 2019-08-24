---
title: 数组
---
# 数组
数组是一个由固定长度的指定类型元素组成的序列。数组的长度在编译阶段确定。

声明数组：
```go
var 变量名 [SIZE]类型
```

内置函数 `len` 获取数组长度。通过下标访问元素：
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
初始化数组中 `{}` 中的元素个数不能大于 `[]` 中的数字。
如果 `[]` 设置了 `SIZE`，Go 语言会根据元素的个数来设置数组的大小。
上面代码中的**`...`省略号，表示数组的长度是根据初始化值的个数来计算**。

**声明数组 `SIZE` 是必须的，如果没有，那就是切片了。**

## 二维数组
```go
var a = [3][4]int{  
 {0, 1, 2, 3} ,   /*  第一行索引为 0 */
 {4, 5, 6, 7} ,   /*  第二行索引为 1 */
 {8, 9, 10, 11},   /* 第三行索引为 2 */
}
fmt.Printf("a[%d][%d] = %d\n", 2, 3, a[2][3] )
```

`==` 和 `!=` 比较运算符来比较两个数组，只有当两个数组的所有元素都是相等的时候数组才是相等的。

## 数组传入函数
当调用函数时，函数的形参会被赋值，**所以函数参数变量接收的是一个复制的副本，并不是原始调用的变量。** 但是
这种机制，**如果碰到传递一个大数组时，效率较低。这个时候可以显示的传入一个数组指针**（其他语言其实是隐式的传递指针）。
```go
func test(ptr *[32]byte) {
  *ptr = [32]byte{}
}
```