---
title: 函数
weight: 5
---

# 函数

## 参数传递

Go 使用的是**值传递**，当我们传一个参数值到被调用函数里面时，实际上是传了这个值的一份 copy，（不管是指针，引用类型还是其他类型，
区别无非是拷贝目标对象还是拷贝指针）当在被调用函数中修改参数值的时候，调用函数中相应实参不会发生任何变化，因为数值变化只作用在 copy 上。
但是如果是引用传递，在调用函数时将实际参数的地址传递到函数中，那么在函数中对参数所进行的修改，将影响到实际参数。

注意，如果实参是 `slice`、`map`、`function`、`channel` 等类型（**引用类型**），实参可能会由于函数的间接引用被修改。

没有函数体的函数声明，这表示该函数不是以 Go 实现的。这样的声明定义了函数标识符。

**表面上看，指针参数性能会更好，但是要注意被复制的指针会延长目标对象的生命周期，还可能导致它被分配到堆上，其性能消耗要加上堆内存分配和
垃圾回收的成本。在栈上复制小对象，要比堆上分配内存要快的多**。如果复制成本高，或者需要修改原对象，使用指针更好。

## 可变参数

**变参本质上就是一个切片，只能接受一到多个同类型参数，而且必须在参数列表的最后一个**。比如 `fmt.Printf`，`Printf` 接收一个的必备参数，之
后接收任意个数的后续参数。

在参数列表的最后一个参数类型之前加上省略符号 `...`，表示该函数会接收任意数量的该类型参数。

```go
func sum(vals ...int) int {
 total := 0
 for _, val := range vals {
   total += val
 }
 return total
}

// 调用
fmt.Println(sum())           // "0"
fmt.Println(sum(3))          // "3"
fmt.Println(sum(1, 2, 3, 4)) // "10"

// 还可以使用类似 ES6 的解构赋值的语法
values := []int{1, 2, 3, 4}
fmt.Println(sum(values...)) // "10"
```

## 函数作为值

Go 函数被看作第一类值：函数像其他值一样，拥有类型，可以被赋值给其他变量，传递给函数，从函数返回。

```go
func main(){
 /* 声明函数变量 */
 getSquareRoot := func(x float64) float64 {
  return math.Sqrt(x)
 }

 /* 使用函数 */
 fmt.Println(getSquareRoot(9)) // 3
}
```

## 函数作为参数

声明一个名叫 `operate` 的函数类型，它有两个参数和一个结果，都是 `int` 类型的。

```go
type operate func(x, y int) int
```

编写 `calculate` 函数的签名部分。这个函数除了需要两个 `int` 类型的参数之外，还应该有一个 `operate` 类型的参数。

```go
func calculate(x int, y int, op operate) (int, error) {
    if op == nil {
        return 0, errors.New("invalid operation")
    }
    return op(x, y), nil
}
```

## 闭包

Go 语言支持匿名函数，可作为闭包。

```go
// 返回一个函数
func getSequence() func() int { // func() 是没有参数也没有返回值的函数类型
  i:=0
  // 闭包
   return func() int {
      i+=1
     return i  
   }
}
```

## 关键字 defer

在普通函数或方法前加关键字 `defer`，会使函数或方法延迟执行，直到包含该 `defer` 语句的函数执行完毕时（**无论函数是否出错**），
`defer` 后的函数才会被执行。

Go官方文档中对 `defer` 的执行时机做了阐述，分别是。

- 包裹 `defer` 的函数返回时
- 包裹 `defer` 的函数执行到末尾时
- 所在的 goroutine 发生 panic 时

**注意：** 调用 `os.Exit` 时 `defer` 不会被执行。

`defer` 语句一般被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。因为 `defer` 可以保证让你更任何情况下，
资源都会被释放。

```go
package ioutil
func ReadFile(filename string) ([]byte, error) {
 f, err := os.Open(filename)
 if err != nil {
   return nil, err
 }
 defer f.Close()
 return ReadAll(f)
}

// 互斥锁
var mu sync.Mutex
var m = make(map[string]int)
func lookup(key string) int {
 mu.Lock()
 defer mu.Unlock()
 return m[key]
}

// 记录何时进入和退出函数
func bigSlowOperation() {
 defer trace("bigSlowOperation")() // 运行 trace 函数，记录了进入函数的时间，并返回一个函数值，这个函数值会延迟执行
 extra parentheses
 // ...lots of work…
 time.Sleep(10 * time.Second) // simulate slow
 operation by sleeping
}
func trace(msg string) func() {
 start := time.Now()
 log.Printf("enter %s", msg)
 return func() { 
  log.Printf("exit %s (%s)", msg,time.Since(start)) 
 }
}

// 观察函数的返回值
func double(x int) (result int) { // 有名返回值
  // 由于 defer 在 return 之后执行，所以这里的 result 就是函数最终的返回值
 defer func() { fmt.Printf("double(%d) = %d\n", x,result) }()

 return x + x
}

_ = double(4) // 输出 "double(4) = 8"
```

上面的例子中我们知道 `defer` 函数可以观察函数返回值，`defer` 函数还可以修改函数的返回值：

```go
func triple(x int) (result int) {
 defer func() { result += x }()
 return double(x)
}
fmt.Println(triple(4)) // "12"
```

### defer 的性能

相比直接用 CALL 汇编指令调用函数，`defer` 要花费更大代价，包括注册，调用操作，额为的缓存开销。

```go
func call () {
 m.Lock()
 m.Unlock()
}

func deferCall()  {
 m.Lock()
 defer m.Unlock()
}

func BenchmarkCall(b *testing.B)  {
 for i := 0; i < b.N; i ++ {
  call()
 }
}


func BenchmarkDeferCall(b *testing.B)  {
 for i := 0; i < b.N; i ++ {
  deferCall()
 }
}
```

```sh
$ go test -bench=.
goos: windows
goarch: amd64
pkg: github.com/shipengqi/golang-learn/demos/defers
BenchmarkCall-8         92349604                12.9 ns/op
BenchmarkDeferCall-8    34305316                36.3 ns/op
PASS
ok      github.com/shipengqi/golang-learn/demos/defers  2.571s

```

性能相差三倍，尽量避免使用 `defer`。

### 什么时候不应该使用 defer

比如处理日志文件，不恰当的 `defer` 会导致关闭文件延时。

```go
func main() {
    for i := 0; i < 100; i ++ {
        f, err := os.Open(fmt.Sprintf("%d.log", i))
        if err != nil {
            continue
        }
        defer f.Close()
        // something
    }
}

```

上面的 `defer` 导致所有的 `f` 都是在 `main` 函数退出时才调用，白白消耗了资源。所以应该直接调用 `Close` 函数，
将文件操作封装到一个函数中，在该函数中调用 `Close` 函数。

### 如果一个函数中有多条 defer 语句，那么那几个 defer 函数调用的执行顺序是怎样的

在同一个函数中，**`defer` 函数调用的执行顺序与它们分别所属的 `defer` 语句的出现顺序（更严谨地说，是执行顺序）完全相反**。

在 `defer` 语句每次执行的时候，Go 语言会把它携带的 `defer` 函数及其参数值另行存储到一个队列中。

这个队列与该 `defer` 语句所属的函数是对应的，并且，它是先进后出（FILO）的，相当于一个栈。

在需要执行某个函数中的 `defer` 函数调用的时候，Go 语言会先拿到对应的队列，然后从该队列中一个一个地取出 `defer` 函数及
其参数值，并逐个执行调用。

## 传入函数的那些参数值后来怎么样了

```go
package main

import "fmt"

func main() {
 array1 := [3]string{"a", "b", "c"}
 fmt.Printf("The array: %v\n", array1)
 array2 := modifyArray(array1)
 fmt.Printf("The modified array: %v\n", array2)
 fmt.Printf("The original array: %v\n", array1)
}

func modifyArray(a [3]string) [3]string {
 a[1] = "x"
 return a
}
```

在 `main` 函数中声明了一个数组 `array1`，然后把它传给了函数 `modify`，`modify` 对参数值稍作修改后将其作为结果值返回。`main`
 函数中的代码拿到这个结果之后打印了它（即 `array2`），以及原来的数组 `array1`。关键问题是，原数组会因 `modify` 函数对参数
 值的修改而改变吗？

答案是：原数组不会改变。为什么呢？原因是，**所有传给函数的参数值都会被复制，函数在其内部使用的并不是参数值的原值，
而是它的副本**。

由于数组是值类型，所以每一次复制都会拷贝它，以及它的所有元素值。

注意，**对于引用类型，比如：切片、字典、通道，像上面那样复制它们的值，只会拷贝它们本身而已，并不会拷贝它们引用的底层数据。
也就是说，这时只是浅表复制，而不是深层复制**。

以切片值为例，如此复制的时候，只是拷贝了它指向底层数组中某一个元素的指针，以及它的长度值和容量值，而它的底层数组并不会被拷贝。