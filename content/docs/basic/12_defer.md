---
title: defer
weight: 12
draft: true
---

# defer

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