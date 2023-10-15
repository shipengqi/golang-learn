---
title: range
weight: 10
---

# range

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
