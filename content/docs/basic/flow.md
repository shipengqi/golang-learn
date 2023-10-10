---
title: 控制语句
---

## select

`select` 类似于用于通信的 `switch` 语句。每个 `case` 必须是一个通信操作，要么是发送要么是接收。

当条件满足时，`select` 会去通信并执行 `case` 之后的语句，这时候其它通信是不会执行的。
如果多个 `case` 同时满足条件，`select` 会随机地选择一个执行。如果没有 `case` 可运行，它将阻塞，直到有 `case` 可运行。

一个默认的子句应该总是可运行的。

```go
select {
  case communication clause:
      ...
  case communication clause:
      ...
  default: /* 可选 */
   ...
}   
```

`for range` 支持遍历数组，切片，字符串，字典，通道，并返回索引和键值。**`for range` 会复制目标数据。可改用数组指针或者切片**。

`range` 关键字右边的位置上的代码被称为 `range` 表达式。

1. **`range` 表达式只会在 `for` 语句开始执行时被求值一次，无论后边会有多少次迭代**；
2. `range` 表达式的求值结果会被复制，也就是说，被迭代的对象是 `range` 表达式结果值的副本而不是原值。
3. `for range` 在性能比 `for` 稍差，因为 `for range` 会进行值拷贝。

字符串的复制成本很小，切片，字典，通道等引用类型本身是指针的封装，复制成本也很小，无序专门优化。

**如果 `range` 的目标表达式是函数，也只会运行一次**。

```go
numbers1 := []int{1, 2, 3, 4, 5, 6}
for i := range numbers1 {
    if i == 3 {
        numbers1[i] |= i
    }
}
fmt.Println(numbers1)
```

打印的内容会是 `[1 2 3 7 5 6]`，为什么，首先 `i` 是切片的下标，当 `i` 的值等于 3 的时候，与之对应的是切片中的第 4 个元素
值 4。对 4 和 3 进行按位或操作得到的结果是 7。

当 `for` 语句被执行的时候，在 `range` 关键字右边的 `numbers1` 会先被求值。`range` 表达式的结果值可以是数组、数组的指针、
切片、字符串、字典或者允许接收操作的通道中的某一个，并且结果值只能有一个。这里的 `numbers1` 是一个切片,那么迭代变量就可以
有两个，右边的迭代变量代表当次迭代对应的某一个元素值，而左边的迭代变量则代表该元素值在切片中的索引值。
循环控制语句：

- `break`，用于中断当前 `for` 循环或跳出 `switch` 语句
- `continue`，跳过当前循，继续进行下一轮循环。
- `goto`，将控制转移到被标记的语句。通常与条件语句配合使用。可用来实现条件转移， 构成循环，跳出循环体等功能。不推荐
使用，以免造成流程混乱。

`goto` 实例：

```go
LOOP: for a < 20 {
 if a == 15 {
   /* 跳过迭代 */
   a = a + 1
   goto LOOP
 }
 fmt.Printf("a的值为 : %d\n", a)
 a ++  
}  
```

### 使用带有 `range` 子句的 `for` 语句遍历字符串值的时候应该注意

带有 `range` 子句的 `for` 语句会先把被遍历的字符串值拆成一个**字节序列**（注意是字节序列），然后再试图找出这个字节序列中
包含的每一个 UTF-8 编码值，或者说每一个 Unicode 字符。

这样的 `for` 语句可以为两个迭代变量赋值。如果存在两个迭代变量，那么赋给第一个变量的值就将会是当前字节序列中的某个 UTF-8 编码
值的第一个字节所对应的那个索引值。而赋给第二个变量的值则是这个 UTF-8 编码值代表的那个 Unicode 字符，其类型会是 `rune`。

```go
str := "Go 爱好者 "
for i, c := range str {
    fmt.Printf("%d: %q [% x]\n", i, c, []byte(string(c)))
}
```

完整的打印内容如下：

```bash
0: 'G' [47]
1: 'o' [6f]
2: '爱' [e7 88 b1]
5: '好' [e5 a5 bd]
8: '者' [e8 80 85]
```

**注意了，'爱'是由三个字节共同表达的，所以第四个 Unicode 字符'好'对应的索引值并不是 3，而是 2 加 3 后得到的 5**。

<https://studygolang.com/articles/25094>
<https://studygolang.com/articles/9701>
<https://talkgo.org/discuss/2019-01-10-anlayze-range/>

<https://blog.csdn.net/qq_25870633/article/details/83339538>
<https://zhuanlan.zhihu.com/p/91044663>
<https://www.jianshu.com/p/86a99efeece5>
<https://blog.csdn.net/u011957758/article/details/82230316>
<https://www.cnblogs.com/howo/p/10507934.html>
