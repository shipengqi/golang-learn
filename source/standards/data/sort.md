---
title: sort
---

# sort
`sort` 包中实现了几种基本的排序算法：插入排序、归并排序、堆排序和快速排序。但是这四种排序方法是不公开的，只用于 `sort` 包
内部使用。所以在对数据集合排序时不必考虑应当选择哪一种排序方法，只要实现了 `sort.Interface` 定义的三个方法就可以对数据集合进
行排序。`sort` 包会根据实际数据自动选择高效的排序算法。

```go
type Interface interface {
	// Len 为集合内元素的总数
	Len() int
	// 如果 index 为 i 的元素小于 index 为 j 的元素，则返回 true，否则 false
	Less(i, j int) bool
	// Swap 交换索引为 i 和 j 的元素
	Swap(i, j int)
}
```


为了方便对常用数据类型的操作，`sort` 包原生支持 `[]int`、`[]float64` 和 `[]string` 三种内建数据类型切片的排序操作。
即不必实现 `sort.Interface` 接口的三个方法。

## 数据集合排序

对数据集合（包括自定义数据类型的集合）排序需要实现 `sort.Interface` 接口的三个方法：

数据集合实现了这三个方法后，即可调用该包的 `Sort()` 方法进行排序。
`Sort()` 方法定义如下：
```go
func Sort(data Interface)
```
`Sort()` 方法惟一的参数就是待排序的数据集合。

该包还提供了一个方法可以判断数据集合是否已经排好顺序，该方法的内部实现依赖于我们自己实现的 `Len()` 和 `Less()` 方法：
```go
func IsSorted(data Interface) bool {
    n := data.Len()
    for i := n - 1; i > 0; i-- {
        if data.Less(i, i-1) {
            return false
        }
    }
    return true
}
```
使用 `sort` 包对学生成绩排序的示例：

```go
package main

import (
	"fmt"
	"sort"
)

// 学生成绩结构体
type StuScore struct {
    name  string    // 姓名
    score int   // 成绩
}

type StuScores []StuScore

func (s StuScores) Len() int {
	return len(s)
}

func (s StuScores) Less(i, j int) bool {
	return s[i].score < s[j].score
}

func (s StuScores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func main() {
    students := StuScores{
        {"alan", 95},
        {"hikerell", 91},
        {"acmfly", 96},
        {"leao", 90},
	}

	// 打印未排序的 students 数据
    fmt.Println("Default:\n\t", students)
    // StuScores 已经实现了 sort.Interface 接口 , 所以可以调用 Sort 函数进行排序
	sort.Sort(students)
	// 判断是否已经排好顺序
	fmt.Println("IS Sorted?\n\t", sort.IsSorted(students))
	// 打印排序后的 students 数据
    fmt.Println("Sorted:\n\t",students)
}
```
输出：
```sh
Default:
     [{alan 95} {hikerell 91} {acmfly 96} {leao 90}]
IS Sorted?
     true
Sorted:
     [{leao 90} {hikerell 91} {alan 95} {acmfly 96}]
```

## Reverse

上面的代码实现的是升序排序，如果要实现降序排序修改 `Less()` 函数：
```go
// 将小于号修改为大于号
func (s StuScores) Less(i, j int) bool {
	return s[i].score > s[j].score
}
```
此外，`sort`包提供了 `Reverse()` 方法，可以允许将数据按 `Less()` 定义的排序方式逆序排序，而不必修改 `Less()` 代码。
```go
func Reverse(data Interface) Interface
```

`Reverse()` 返回的一个 `sort.Interface` 接口类型，整个 `Reverse()` 的内部实现比较有趣：
```go
// 定义了一个 reverse 结构类型，嵌入 Interface 接口
type reverse struct {
    Interface
}

// reverse 结构类型的 Less() 方法拥有嵌入的 Less() 方法相反的行为
func (r reverse) Less(i, j int) bool {
    return r.Interface.Less(j, i)
}

// 返回新的实现 Interface 接口的数据类型
func Reverse(data Interface) Interface {
    return &reverse{data}
}
```
了解内部原理后，可以在学生成绩排序示例中使用 `Reverse()` 来实现成绩升序排序：
```go
sort.Sort(sort.Reverse(students))
fmt.Println(students)
```

## Search
```go
func Search(n int, f func(int) bool) int
```

`Search()` 函数一个常用的使用方式是搜索元素 x 是否在已经升序排好的切片 s 中：

```go
x := 11
s := []int{3, 6, 8, 11, 45} // 已经升序排序的集合
pos := sort.Search(len(s), func(i int) bool { return s[i] >= x })
if pos < len(s) && s[pos] == x {
    fmt.Println(x, " 在 s 中的位置为：", pos)
} else {
    fmt.Println("s 不包含元素 ", x)
}
```

官方文档给出的小程序：

```go
func GuessingGame() {
	var s string
	fmt.Printf("Pick an integer from 0 to 100.\n")
	answer := sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d? ", i)
		fmt.Scanf("%s", &s)
		return s != "" && s[0] == 'y'
	})
	fmt.Printf("Your number is %d.\n", answer)
}
```

## 已经支持的内部数据类型排序
`sort`包原生支持 `[]int`、`[]float64` 和 `[]string` 三种内建数据类型切片的排序操作。

### IntSlice 类型和 []int

`[]int` 切片排序内部实现及使用方法与 `[]float64` 和 `[]string` 类似。

`sort`包定义了一个 `IntSlice` 类型，并且实现了 `sort.Interface` 接口：

```go
    type IntSlice []int
    func (p IntSlice) Len() int           { return len(p) }
    func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
    func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
    // IntSlice 类型定义了 Sort() 方法，包装了 sort.Sort() 函数
    func (p IntSlice) Sort() { Sort(p) }
    // IntSlice 类型定义了 Search() 方法，包装了 SearchInts() 函数
    func (p IntSlice) Search(x int) int { return SearchInts(p, x) }
```
并且提供的 `sort.Ints()` 方法使用了该 `IntSlice` 类型：
```go
    func Ints(a []int) { Sort(IntSlice(a)) }
```

所以，对 `[]int` 切片排序更常使用 `sort.Ints()`，而不是直接使用 `IntSlice` 类型：

```go
s := []int{5, 2, 6, 3, 1, 4} // 未排序的切片数据
sort.Ints(s)
fmt.Println(s) // 将会输出[1 2 3 4 5 6]
```
如果要使用降序排序，显然要用前面提到的 Reverse() 方法：

```go
s := []int{5, 2, 6, 3, 1, 4} // 未排序的切片数据
sort.Sort(sort.Reverse(sort.IntSlice(s)))
fmt.Println(s) // 将会输出[6 5 4 3 2 1]
```

如果要查找整数 `x` 在切片 `a` 中的位置，相对于前面提到的 `Search()` 方法，`sort` 包提供了 `SearchInts()`:

```go
func SearchInts(a []int, x int) int
```
注意，`SearchInts()` 的使用条件为：**切片 `a` 已经升序排序**
以下是一个错误使用的例子：

```go
s := []int{5, 2, 6, 3, 1, 4} // 未排序的切片数据
fmt.Println(sort.SearchInts(s, 2)) // 将会输出 0 而不是 1
```

### Float64Slice 类型及 []float64

实现与 Ints 类似：

```go
type Float64Slice []float64

func (p Float64Slice) Len() int           { return len(p) }
func (p Float64Slice) Less(i, j int) bool { return p[i] < p[j] || isNaN(p[i]) && !isNaN(p[j]) }
func (p Float64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Float64Slice) Sort() { Sort(p) }
func (p Float64Slice) Search(x float64) int { return SearchFloat64s(p, x) }
```
与 `Sort()`、`IsSorted()`、`Search()` 相对应的三个方法：

```go
func Float64s(a []float64)
func Float64sAreSorted(a []float64) bool
func SearchFloat64s(a []float64, x float64) int
```

在上面 `Float64Slice` 类型定义的 `Less` 方法中，有一个内部函数 `isNaN()`。
`isNaN()` 与 `math` 包中 `IsNaN()` 实现完全相同，`sort` 包之所以不使用 `math.IsNaN()`，完全是基于包依赖性的考虑，
`sort` 包的实现不依赖与其他任何包。

### StringSlice 类型及 []string

两个 `string` 对象之间的大小比较是基于“字典序”的。

```go
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p StringSlice) Sort() { Sort(p) }
func (p StringSlice) Search(x string) int { return SearchStrings(p, x) }
```

## []interface 排序与查找

只要实现了 `sort.Interface` 接口，即可通过 `sort` 包内的函数完成排序，查找等操作。但是这种用法对于其它数据类型的 `slice`
不友好，可能我们需要为大量的 `struct` 定义一个单独的 `[]struct` 类型，再为其实现 `sort.Interface` 接口，例如：
```go
type Person struct {
    Name string
    Age  int
}
type Persons []Person

func (p Persons) Len() int {
    panic("implement me")
}

func (p Persons) Less(i, j int) bool {
    panic("implement me")
}

func (p Persons) Swap(i, j int) {
    panic("implement me")
}
```

`sort` 包提供了以下函数：

```go
func Slice(slice interface{}, less func(i, j int) bool)
func SliceStable(slice interface{}, less func(i, j int) bool)
func SliceIsSorted(slice interface{}, less func(i, j int) bool) bool
func Search(n int, f func(int) bool) int
```
排序相关的三个函数都接收 `[]interface`，并且需要传入一个比较函数，用于为程序比较两个变量的大小，因为
函数签名和作用域的原因，这个函数只能是 `匿名函数`。

### sort.Slice
利用 `sort.Slice` 函数，而不用提供一个特定的 `sort.Interface` 的实现，而是 `Less(i，j int)` 作为一个比较回调函数，可以简单
地传递给 `sort.Slice` 进行排序。**不建议使用，因为在 `sort.Slice` 中使用了 `reflect`**。

```go
people := []struct {
    Name string
    Age  int
}{
    {"Gopher", 7},
    {"Alice", 55},
    {"Vera", 24},
    {"Bob", 75},
}

sort.Slice(people, func(i, j int) bool { return people[i].Age < people[j].Age }) // 按年龄升序排序
fmt.Println("Sort by age:", people)

// Output:
// Sort by age: [{Gopher 7} {Vera 24} {Alice 55} {Bob 75}]
```

### sort.Search

该函数判断 `[]interface` 是否存在指定元素，举个栗子：

- 升序 slice

> sort 包为 `[]int`,`[]float64`,`[]string` 提供的 Search 函数其实也是调用的该函数，因为该函数是使用的二分查找法，所以要
求 slice 为升序排序状态。并且判断条件必须为 `>=`，这也是官方库提供的三个查找相关函数的的写法。

```go
a := []int{2, 3, 4, 200, 100, 21, 234, 56}
x := 21

sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })   // 升序排序
index := sort.Search(len(a), func(i int) bool { return a[i] >= x }) // 查找元素

if index < len(a) && a[index] == x {
    fmt.Printf("found %d at index %d in %v\n", x, index, a)
} else {
    fmt.Printf("%d not found in %v,index:%d\n", x, a, index)
}

// Output:
// found 21 at index 3 in [2 3 4 21 56 100 200 234]
```

