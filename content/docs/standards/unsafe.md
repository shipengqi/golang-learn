Go 的指针多了一些限制。但这也算是 Go 的成功之处：既可以享受指针带来的便利，又避免了指针的危险性。

限制一：Go的指针不能进行数学运算。

```go
a := 5
p := &a

p++
p = &a + 3
```

上面的代码将不能通过编译，会报编译错误：invalid operation，也就是说不能对指针做数学运算。

限制二：不同类型的指针不能相互转换。

```go
func main() {
  a :=int(100)
  var f *float64
  f = &a
}
```

也会报编译错误：`cannot use &a (type *int) as type *float64 in assignment`

限制三：不同类型的指针不能使用==或!=比较。

只有在两个指针类型相同或者可以相互转换的情况下，才可以对两者进行比较。另外，指针可以通过 == 和 != 直接和 nil 作比较。

限制四：不同类型的指针变量不能相互赋值。

这一点同限制三。

## unsafe

Go 还有非类型安全的指针，这就是 unsafe 包提供的 `unsafe.Pointer`。

unsafe 包用于 Go 编译器，在编译阶段使用。从名字就可以看出来，它是不安全的，官方并不建议使用。

它可以绕过 Go 语言的类型系统，直接操作内存。例如，一般我们不能操作一个结构体的未导出成员，但是通过 unsafe 包就能做到。unsafe 包让我可以直接读写内存，还管你什么导出还是未导出。

Go 语言类型系统是为了安全和效率设计的，有时，安全会导致效率低下。有了 unsafe 包，高阶的程序员就可以利用它绕过类型系统的低效。

## unsafe 实现原理

```go
tyoe ArbitraryType int
type Pointer *ArbitraryType
```

Arbitrary 是任意的意思，也就是说 **Pointer 可以指向任意类型**，实际上它类似于 C 语言里的 `void*`。

三个函数：

```go
func SizeOf(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr
```

Sizeof 返回类型 x 所占据的字节数，但不包含 x 所指向的内容的大小。例如，对于一个指针，函数返回的大小为 8 字节（64位机上），一个 slice 的大小则为 slice header 的大小。

Offsetof 返回结构体成员在内存中的位置离结构体起始处的字节数，所传参数必须是结构体的成员。

Alignof 返回 m，m 是指当类型进行内存对齐时，它分配到的内存地址能整除 m。

注意到以上三个函数返回的结果都是 uintptr 类型，这和 unsafe.Pointer 可以相互转换。三个函数都是在编译期间执行，它们的结果可以直接赋给 const型变量。另外，因为三个函数执行的结果和操作系统、编译器相关，所以是不可移植的。

## 使用 unsafe

### slice 长度

### map 长度
