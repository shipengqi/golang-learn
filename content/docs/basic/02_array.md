---
title: Array
weight: 2
---

# Array

数组是一个由固定长度，并且相同类型的元素组成的数据结构。计算机会为数组分配一块连续的内存来保存其中的元素，我们可以利用数组中元素的索引快速访问特定元素。

**存储元素类型相同、但是大小不同的数组类型在 Go 语言看来也是完全不同的，只有两个条件都相同才是同一类型**。

Go 是值传递，所以函数参数变量接收的是一个值的副本。这种机制，**在传递一个大数组时，效率较低。这个时候可以显示的传入一个数组指针**（其他语言其实是隐式的指针传递）。

```go
func test(ptr *[32]byte) {
  *ptr = [32]byte{}
}
```

## 初始化

```go
arr1 := [3]int{1, 2, 3}
arr2 := [...]int{1, 2, 3} // `...` 省略号，表示数组的长度是根据初始化值的个数来计算
```

数组的长度在编译阶段确定，初始化之后大小就无法改变。

```go
// NewArray returns a new fixed-length array Type.
func NewArray(elem *Type, bound int64) *Type {
	if bound < 0 {
		base.Fatalf("NewArray: invalid bound %v", bound)
	}
	t := newType(TARRAY)
	t.extra = &Array{Elem: elem, Bound: bound}
	if elem.HasShape() {
		t.SetHasShape(true)
	}
	return t
}
```

编译期间由 [cmd/compile/internal/types.NewArray](https://github.com/golang/go/blob/2744155d369ca838be57d1eba90c3c6bfc4a3b30/src/cmd/compile/internal/types/type.go#L542) 生成数组类型。

参数 `elem *Type` 数组类型和 `bound int64` 数组的大小构成了 `Array` 类型。当前**数组是否应该在堆栈中初始化在编译期就确定了**。

对于一个由字面量组成的数组，根据数组元素数量的不同，编译器会在负责初始化字面量的 [cmd/compile/internal/walk.anylit](https://github.com/golang/go/blob/2744155d369ca838be57d1eba90c3c6bfc4a3b30/src/cmd/compile/internal/walk/complit.go#L527) 函数中做两种不同的优化：

- 当元素数量小于或者等于 4 个时，会直接将数组中的元素放置在栈上；
- 当元素数量大于 4 个时，会将数组中的元素放置到静态区并在运行时取出；

当元素数量小于或者等于 4 个时，[cmd/compile/internal/walk.fixedlit](https://github.com/golang/go/blob/2744155d369ca838be57d1eba90c3c6bfc4a3b30/src/cmd/compile/internal/walk/complit.go#L192) 会负责在函数编译之前将 `[3]{1, 2, 3}` 转换成更加原始的语句：

```go

```