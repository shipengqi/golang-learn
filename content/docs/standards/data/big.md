---
title: big
---

# big
`big` 是 Go 语言提供的进行大数操作的标准库，实现了任意精度算术（大数）。

Go 语言中的 `float64` 类型进行浮点运算，返回结果将精确到 15 位，足以满足大多数的任务。但是当对超出 `int64` 或者 `uint64` 类型这样的大
数进行计算时，如果对精度没有要求，`float32` 或者 `float64` 可以胜任，但如果对精度有严格要求的时候，则不能使用浮点数，在内存中它们只能
被近似的表示。

对于整数的高精度计算 Go 语言中提供了 `big` 包，被包含在 `math` 包下：有用来表示大整数的 `big.Int` 和表示大有理数的 `big.Rat` 类型
（可以表示为 `2/5` 或 `3.1416` 这样的分数，而不是无理数或 `π`）。这些类型可以实现任意位类型的数字，只要内存足够大。缺点是更大的内存
和处理开销使它们使用起来要比内置的数字类型慢很多。

大的整型数字是通过 `big.NewInt(n)` 来构造的，其中 `n` 为 `int64` 类型整数。而大有理数是通过 `big.NewRat(n, d)` 方法构造。`n`（分子）
和 `d`（分母）都是 `int64` 型整数。因为 Go 语言不支持运算符重载，所以所有大数字类型都有像是 `Add()` 和 `Mul()` 这样的方法。它们作用
于作为 receiver 的整数和有理数，大多数情况下它们修改 receiver 并以 receiver 作为返回结果。因为没有必要创建 `big.Int` 类型的临
时变量来存放中间结果，所以运算可以被链式地调用，并节省内存。

示例：

```go
package main

import (
	"fmt"
	"math"
	"math/big"
)

func main() {
	// Here are some calculations with bigInts:
	im := big.NewInt(math.MaxInt64)
	in := im
	io := big.NewInt(1956)
	ip := big.NewInt(1)
	ip.Mul(im, in).Add(ip, im).Div(ip, io)
	fmt.Printf("Big Int: %v\n", ip)
	// Here are some calculations with bigInts:
	rm := big.NewRat(math.MaxInt64, 1956)
	rn := big.NewRat(-1956, math.MaxInt64)
	ro := big.NewRat(19, 56)
	rp := big.NewRat(1111, 2222)
	rq := big.NewRat(1, 1)
	rq.Mul(rm, rn).Add(rq, ro).Mul(rq, rp)
	fmt.Printf("Big Rat: %v\n", rq)
}

/* Output:
Big Int: 43492122561469640008497075573153004
Big Rat: -37/112
*/
```
