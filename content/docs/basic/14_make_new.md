---
title: make 和 new
weight: 14
---

Go 中初始化一个结构，有两个关键字：`make` 和 `new`。虽然都是用于初始化结构，但是有很大的不同。

1. `new` 是根据传入的类型分配一片内存空间，并返回指向这片内存空间的指针。
   - 任何类型都可以使用 `new` 来初始化。
   - **内存里存的值是对应类型的零值**，这就意味着，使用 `new` 初始化切片、map 和 channel 时，得到是 `nil`。
2. `make` 是用来初始化内置的数据结构，也就是切片、map 和 channel。
   ```go
   // sl 是一个结构体 reflect.SliceHeader；
   sl := make([]int, 0, 100)
   // m 是一个指向 runtime.hmap 结构体的指针
   m := make(map[int]bool, 10)
   // ch 是一个指向 runtime.hchan 结构体的指针
   ch := make(chan int, 5)
   ```
