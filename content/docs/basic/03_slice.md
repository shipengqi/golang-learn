---
title: 切片
weight: 3
---

切片 (slice) 在使用上和数组差不多，区别是切片是**可变长**的，定义的时候不需要指定 size。

切片可以看做是对数组的一层简单的封装，切片的底层数据结构中，包含了一个数组。

切片的结构体：

```go
// src/reflect/value.go
type SliceHeader struct {
	Data uintptr // 指向底层数组
	Len  int     // 当前切片长度
	Cap  int     // 当前切片容量
}
```

注意 `Cap` 也是底层数组的长度。`Data` 是一块连续的内存，可以存储切片 `Cap` 大小的所有元素。

![slice-struct](https://gitee.com/shipengqi/illustrations/raw/main/go/slice-struct.png)

如图，虽然 slice 的 `Len` 是 5，但是底层数组的长度是 10，也就是 `Cap`。

## 初始化

初始化切片有三种方式：

1. 使用 `make`
   ```go
    // len 是切片的初始长度
    // capacity 为可选参数, 指定容量
    s := make([]int, len, capacity)
   ```
2. 使用字面量
   ```go
   arr :=[]int{1,2,3}
   ```
3. 使用下标截取数组或者切片的一部分，这里可以传入三个参数 `[low:high:max]`，`max - low` 是新的切片的容量 cap。
   ```go
   numbers := []int{0,1,2,3,4,5,6,7,8}
   s := numbers[1:4] // [1 2 3]
   s := numbers[4:] // [4 5 6 7 8]
   s := numbers[:3]) // [0 1 2]
   ```

不管使用那种初始化方式，最后都是返回一个 `SliceHeader` 的结构体。

《Go 学习笔记》 第四版 中的示例：

```go
package main

import "fmt"

func main() {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := slice[2:5]
	s2 := s1[2:6:7]

	s2 = append(s2, 100)
	s2 = append(s2, 200)

	s1[2] = 20

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(slice)
}
```

输出：

```
[2 3 20]                
[4 5 6 7 100 200]       
[0 1 2 3 20 5 6 7 100 9]
```

示例中：

- `s1 := slice[2:5]` 得到的 `s1` 的容量为 8，因为没有传入 `max`，容量默认是到底层数组的结尾。
- `s2 := s1[2:6:7]` 得到的 `s2` 的容量为 5（`max - low`）。`s2`，`s1` 和 `slice` 底层数组是同一个，所以 `s2` 中的元素是 `[4,5,6,7]`。

![slice-cut](https://gitee.com/shipengqi/illustrations/raw/main/go/slice-cut.png)

下面的 `s2 = append(s2, 100)` 追加一个元素，容量够用，不需要扩容，但是这个修改会影响所有指向这个底层数组的切片。

![slice-cut-append](https://gitee.com/shipengqi/illustrations/raw/main/go/slice-cut-append.png)

再次追加一个元素 `s2 = append(s2, 200)`，`s2` 的容量不够了，需要扩容，于是 `s2` 申请一块新的连续内存，并将数据拷贝过去，扩容后的容量是原来的 2 倍。
这时候 `s2` 的 `Data` 指向了新的底层数组，已经和 `s1` 这个 `slice` 没有关系了，所以对 `s2` 的修改不会再影响 `s1`。

![slice-cut-append2](https://gitee.com/shipengqi/illustrations/raw/main/go/slice-cut-append2.png)

最后 `s1[2] = 20` 也不会再影响 `s2`。

![slice-cut-append3](https://gitee.com/shipengqi/illustrations/raw/main/go/slice-cut-append3.png)

## 切片是如何扩容的？ 

`append` 是用来向 slice 追加元素的，并**返回一个新的 slice**。

`append` 实际上就是向底层数组添加元素，但是数组的长度是固定的：

当追加元素后切片的大小大于容量，runtime 会对切片进行扩容，这时会申请一块新的连续的内存空间，然后将原数据拷贝到新的内存空间，并且将 `append` 的元素添加到新的底层数组中，并返回这个新的切片。

Go 1.18 后切片的扩容策略：

- 如果当前切片的容量（`oldcap`）小于 256，新切片的容量（`newcap`）为原来的 2 倍.
- 如果当前切片的容量大于 256，计算新切片的容量的公式 `newcap = oldcap+(oldcap+3*256)/4`

## 切片传入函数

Go 是值传递。那么传入一个切片，切片会不会被函数中的操作改变？

**不管传入的是切片还是切片指针，如果改变了底层数组，那么外部切片的底层数组也会被改变**。

示例：

```go
package main

import "fmt"

func appendFunc(s []int) {
   s = append(s, 10, 20, 30)
}

func appendPtrFunc(s *[]int) {
   *s = append(*s, 10, 20, 30)
}

func main() {
   sl := make([]int, 0, 10)

   appendFunc(sl)
   // appendFunc 修改的是 sl 的副本
   // 副本的 struct {
   //    Data uintptr
   //    Len  int
   //    Cap  int
   // }
   // 副本的 len 和 cap 被修改了，但是不会影响外部 slice 的 len 和 cap，所以下面的输出是 []
   fmt.Println(sl) // []
   // appendFunc，虽然没有修改外部 slice 的 len 和 cap，
   // 但是副本 `Data uintptr` 是一个指针的拷贝，和外部 slice 指向的是同一个底层数组
   // 所以底层数组最终是被修改了的，所以下面的输出会包含 10 20 30
   fmt.Println(sl[:10]) // [10 20 30 0 0 0 0 0 0 0]
   // 为什么 sl[:10] 和 sl[:] 的输出不同，是因为 go 的切片的一个优化
   // slice[low:high] 中的 high，最大的取值范围对应着切片的容量（cap），不只是单纯的长度（len）。
   // sl[:10] 可以输出容量范围内的值，并且没有越界。
   // sl[:] 由于 len 为 0，并且没有指定最大索引。high 则会取 len 的值，所以输出为 []
   fmt.Println(sl[:]) // []

   slptr := make([]int, 0, 10)
   appendPtrFunc(&slptr)
   // 这里传入的是切片的指针，会改变外部的 slice slptr
   fmt.Println(slptr) // [10 20 30]
}
```
