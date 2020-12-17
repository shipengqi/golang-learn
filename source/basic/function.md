---
title: 函数
---

## 声明函数

`func` 关键字声明函数：

```go
func 函数名(形式参数列表) (返回值列表) {
    函数体
}
```

如果函数返回一个无名变量或者没有返回值，返回值列表的括号可以省略。如果一个函数声明没有返回值列表，那么这个
函数不会返回任何值。

```go
// 两个 int 类型参数 返回一个 int 类型的值
func max(num1, num2 int) int {
   /* 定义局部变量 */
   var result int

   if (num1 > num2) {
      result = num1
   } else {
      result = num2
   }
   return result 
}

// 返回多个类型的值
func swap(x int, y string) (string, int) {
   return y, x
}

// 有名返回值
func Size(rect image.Rectangle) (width, height int, err error)
```

在函数体中，函数的形参作为局部变量，被初始化为调用者提供的值（函数调用必须按照声明顺序为所有参数提供实参）。函数的形参和有
名返回值（也就是对返回值命名）作为函数最外层的局部变量，被存储在相同的词法块中。

## 参数

Go 语言使用的是**值传递**，当我们传一个参数值到被调用函数里面时，实际上是传了这个值的一份 copy，（不管是指针，引用类型还是其他类型，
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

## 错误

Go 中，对于大部分函数而言，永远无法确保能否成功运行（有一部分函数总是能成功的运行。比如 `strings.Contains` 和
`strconv.FormatBool`）。**通常 Go 函数的最后一个返回值用来传递错误信息**。如果导致失败的原因只有一个，返回值可以是一个布尔值，
通常被命名为 `ok`。否则应该返回一个 `error` 类型。

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

### defer 原理

#### 堆上分配

编译器不仅将 defer 关键字都转换成 runtime.deferproc 函数，它还会通过以下三个步骤为所有调用 defer 的函数末尾插入 runtime.deferreturn 的函数调用

runtime.deferproc 负责创建新的延迟调用；
runtime.deferreturn 负责在函数调用结束时执行所有的延迟调用；

runtime.deferproc 会为 defer 创建一个新的 runtime._defer 结构体、设置它的函数指针 fn、程序计数器 pc 和栈指针 sp 并将相关的参数拷贝到相邻的内存空间中：

```go
func deferproc(siz int32, fn *funcval) {
 sp := getcallersp()
 argp := uintptr(unsafe.Pointer(&fn)) + unsafe.Sizeof(fn)
 callerpc := getcallerpc()

 d := newdefer(siz)
 if d._panic != nil {
  throw("deferproc: d.panic != nil after newdefer")
 }
 d.fn = fn
 d.pc = callerpc
 d.sp = sp
 switch siz {
 case 0:
 case sys.PtrSize:
  *(*uintptr)(deferArgs(d)) = *(*uintptr)(unsafe.Pointer(argp))
 default:
  memmove(deferArgs(d), unsafe.Pointer(argp), uintptr(siz))
 }

 return0()
}
```

最后调用的 runtime.return0 是唯一一个不会触发延迟调用的函数，它可以避免递归调用 runtime.deferreturn 函数。

runtime.newdefer 的作用是获得一个 runtime._defer 结构体，有三种方式：

- 从调度器的延迟调用缓存池 sched.deferpool 中取出结构体并将该结构体追加到当前 Goroutine 的缓存池中；
- 从 Goroutine 的延迟调用缓存池 pp.deferpool 中取出结构体；
- 通过 runtime.mallocgc 在堆上创建一个新的结构体；

无论使用哪种方式，只要获取到 runtime._defer 结构体，它都会被追加到所在 Goroutine_defer 链表的最前面。

defer 关键字的插入顺序是从后向前的，而 defer 关键字执行是从前向后的，这也是为什么后调用的 defer 会优先执行。

runtime.deferreturn 会从 Goroutine 的 _defer 链表中取出最前面的 runtime._defer 结构体并调用 runtime.jmpdefer 函数传入需要执行的函数和参数：

```go
func deferreturn(arg0 uintptr) {
 gp := getg()
 d := gp._defer
 if d == nil {
  return
 }
 sp := getcallersp()
 ...

 switch d.siz {
 case 0:
 case sys.PtrSize:
  *(*uintptr)(unsafe.Pointer(&arg0)) = *(*uintptr)(deferArgs(d))
 default:
  memmove(unsafe.Pointer(&arg0), deferArgs(d), uintptr(d.siz))
 }
 fn := d.fn
 gp._defer = d.link
 freedefer(d)
 jmpdefer(fn, uintptr(unsafe.Pointer(&arg0)))
}
```

runtime.jmpdefer 是一个用汇编语言实现的运行时函数，它的主要工作是跳转到 defer 所在的代码段并在执行结束之后跳转回 runtime.deferreturn。

#### 栈上分配

在 1.13 中对 defer 关键字进行了优化，当该关键字在函数体中**最多执行一次时**，编译期间的 cmd/compile/internal/gc.state.call 会将结构体分配到栈上并调用 runtime.deferprocStack：

```go
func (s *state) call(n *Node, k callKind) *ssa.Value {
 ...
 var call *ssa.Value
 if k == callDeferStack {
  // 在栈上创建 _defer 结构体
  t := deferstruct(stksize)
  ...

  ACArgs = append(ACArgs, ssa.Param{Type: types.Types[TUINTPTR], Offset: int32(Ctxt.FixedFrameSize())})
  aux := ssa.StaticAuxCall(deferprocStack, ACArgs, ACResults) // 调用 deferprocStack
  arg0 := s.constOffPtrSP(types.Types[TUINTPTR], Ctxt.FixedFrameSize())
  s.store(types.Types[TUINTPTR], arg0, addr)
  call = s.newValue1A(ssa.OpStaticCall, types.TypeMem, aux, s.mem())
  call.AuxInt = stksize
 } else {
  ...
 }
 s.vars[&memVar] = call
 ...
}
```

因为在编译期间我们已经创建了 runtime._defer 结构体，所以 runtime.deferprocStack 函数在运行期间我们只需要设置以为未在编译期间初始化的值并将栈上的结构体追加到函数的链表上：

```go
func deferprocStack(d *_defer) {
 gp := getg()
 d.started = false
 d.heap = false // 栈上分配的 _defer
 d.openDefer = false
 d.sp = getcallersp()
 d.pc = getcallerpc()
 d.framepc = 0
 d.varp = 0
 *(*uintptr)(unsafe.Pointer(&d._panic)) = 0
 *(*uintptr)(unsafe.Pointer(&d.fd)) = 0
 *(*uintptr)(unsafe.Pointer(&d.link)) = uintptr(unsafe.Pointer(gp._defer))
 *(*uintptr)(unsafe.Pointer(&gp._defer)) = uintptr(unsafe.Pointer(d))

 return0()
}
```

除了分配位置的不同，栈上分配和堆上分配的 runtime._defer 并没有本质的不同，而该方法可以适用于绝大多数的场景，与堆上分配的 runtime._defer 相比，该方法可以将 defer 关键字的额外开销降低 ~30%。

#### 开放编码

在 1.14 中通过开发编码（Open Coded）实现 defer 关键字，该设计使用代码内联优化 defer 关键的额外开销并引入函数数据 funcdata 管理 panic 的调用3，该优化可以将 defer 的调用开销从 1.13 版本的 ~35ns 降低至 ~6ns 左右：

开发编码只会在满足以下的条件时启用：

- 函数的 defer 数量少于或者等于 8 个；
- 函数的 defer 关键字不能在循环中执行；
- 函数的 return 语句与 defer 语句的乘积小于或者等于 15 个；

一旦确定使用开放编码，就会在编译期间初始化延迟比特和延迟记录。

编译期间判断 defer 关键字、return 语句的个数确定是否开启开放编码优化；
通过 deferBits 和 cmd/compile/internal/gc.openDeferInfo 存储 defer 关键字的相关信息；
如果 defer 关键字的执行可以在编译期间确定，会在函数返回前直接插入相应的代码，否则会由运行时的 runtime.deferreturn 处理；

## panic 和 recover 原理

panic 能够改变程序的控制流，函数调用panic 时会立刻停止执行函数的其他代码，并在执行结束后在当前 Goroutine 中递归执行调用方的延迟函数调用 defer；
recover 可以中止 panic 造成的程序崩溃。它是一个只能在 defer 中发挥作用的函数，在其他作用域中调用不会发挥任何作用；

- panic 只会触发当前 Goroutine 的延迟函数调用；
- recover 只有在 defer 函数中调用才会生效；
- panic 允许在 defer 中嵌套多次调用；

defer 关键字对应的 runtime.deferproc 会将延迟调用函数与调用方所在 Goroutine 进行关联。所以当程序发生崩溃时只会调用当前 Goroutine 的延迟调用函数也是非常合理的。

多个 Goroutine 之间没有太多的关联，一个 Goroutine 在 panic 时也不应该执行其他 Goroutine 的延迟函数。

recover 只有在发生 panic 之后调用才会生效。需要在 defer 中使用 recover 关键字。

多次调用 panic 也不会影响 defer 函数的正常执行。所以使用 defer 进行收尾的工作一般来说都是安全的。

数据结构 runtime._panic

```go
type _panic struct {
 argp      unsafe.Pointer
 arg       interface{}
 link      *_panic
 recovered bool
 aborted   bool

 pc        uintptr
 sp        unsafe.Pointer
 goexit    bool
}
```

runtime.gopanic，该函数的执行过程包含以下几个步骤：

1. 创建新的 runtime._panic 结构并添加到所.在 Goroutine_panic 链表的最前面；
2. 在循环中不断从当前 Goroutine 的 _defer .中链表获取 runtime._defer 并调用 runtime.reflectcall 运行延迟调用函数；
3. 调用 runtime.fatalpanic 中止整个程序；

### 崩溃恢复

编译器会将关键字 recover 转换成 runtime.gorecover：

```go
func gorecover(argp uintptr) interface{} {
 p := gp._panic
 if p != nil && !p.recovered && argp == uintptr(p.argp) {
  p.recovered = true
  return p.arg
 }
 return nil
}
```

如果当前 Goroutine 没有调用 panic，那么该函数会直接返回 nil，这也是崩溃恢复在非 defer 中调用会失效的原因。

在正常情况下，它会修改 runtime._panic 结构体的 recovered 字段，runtime.gorecover 函数本身不包含恢复程序的逻辑，程序的恢复也是由 runtime.gopanic 函数负责的：

```go
func gopanic(e interface{}) {
 ...

 for {
  // 执行延迟调用函数，可能会设置 p.recovered = true
  ...

  pc := d.pc
  sp := unsafe.Pointer(d.sp)

  ...
  if p.recovered {
   gp._panic = p.link
   for gp._panic != nil && gp._panic.aborted {
    gp._panic = gp._panic.link
   }
   if gp._panic == nil {
    gp.sig = 0
   }
   gp.sigcode0 = uintptr(sp)
   gp.sigcode1 = pc
   mcall(recovery)
   throw("recovery failed")
  }
 }
 ...
}
```


编译器会负责做转换关键字的工作；
将 panic 和 recover 分别转换成 runtime.gopanic 和 runtime.gorecover；
将 defer 转换成 deferproc 函数；
在调用 defer 的函数末尾调用 deferreturn 函数；
在运行过程中遇到 gopanic 方法时，会从 Goroutine 的链表依次取出 _defer 结构体并执行；
如果调用延迟执行函数时遇到了 gorecover 就会将 _panic.recovered 标记成 true 并返回 panic 的参数；
在这次调用结束之后，gopanic 会从 _defer 结构体中取出程序计数器 pc 和栈指针 sp 并调用 recovery 函数进行恢复程序；
recovery 会根据传入的 pc 和 sp 跳转回 deferproc；
编译器自动生成的代码会发现 deferproc 的返回值不为 0，这时会跳回 deferreturn 并恢复到正常的执行流程；
如果没有遇到 gorecover 就会依次遍历所有的 _defer 结构，并在最后调用 fatalpanic 中止程序、打印 panic 的参数并返回错误码 2；