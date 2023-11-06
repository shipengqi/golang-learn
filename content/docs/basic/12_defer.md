---
title: defer
weight: 12
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