---
title: range
weight: 10
draft: true
---

为什么下面的 `range` 遍历只执行了 3 次？

```go
func main() {
    nums := []int{1, 2, 3}
    for i := range nums {
        nums = append(nums, i)
        fmt.Println(nums)
    }
    // 输出：
    // [1 2 3 0]
    // [1 2 3 0 1]
    // [1 2 3 0 1 2]
}
```

因为 `range` 遍历的是切片的副本，而不是切片本身。在遍历的过程中，切片的长度会发生变化，但是遍历的次数在遍历开始时就已经确定了。


## range 变量作用域问题

```go
package main

import "sync"
import "time"

func main() {
    wg := sync.WaitGroup{}
    values := []int{1, 2, 3}
    for _, v := range values {
        go func() {
            fmt.Println(v)
        }()
    }
    <-time.Tick(time.Second)
}
```

Go 1.22 之前的版本会输出：

```bash
3
3
3
```

因为在 `range` 循环中获取返回变量的地址都完全相同（虽然有多个值），变量 v 的作用域覆盖了整个循环体，每次循环只是更新了 v 的值，而没有创建新的变量。

Go 1.22 版本修复了这个问题，**循环时的变量作用域被修改为每次循环都创建一个独立的变量**。