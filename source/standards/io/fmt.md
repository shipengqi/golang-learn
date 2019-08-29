---
title: fmt 包
---

# fmt 包
`fmt` 包实现了格式化 I/O 函数，类似于 C 的 `printf` 和 `scanf`。

```go
type user struct {
	name string
}

func main() {
	u := user{"tang"}
	fmt.Printf("% + v\n", u)     // 格式化输出结构               {name: tang}
	fmt.Printf("%#v\n", u)       // 输出值的 Go 语言表示方法       main.user{name: "tang"}
	fmt.Printf("%T\n", u)        // 输出值的类型的 Go 语言表示     main.user
	fmt.Printf("%t\n", true)     // 输出值的 true 或 false   true
	fmt.Printf("%b\n", 1024)     // 二进制表示               10000000000
	fmt.Printf("%c\n", 11111111) // 数值对应的 Unicode 编码字符
	fmt.Printf("%d\n", 10)       // 十进制表示                 10
	fmt.Printf("%o\n", 8)        // 八进制表示                 10
	fmt.Printf("%q\n", 22)       // 转化为十六进制并附上单引号    '\x16'
	fmt.Printf("%x\n", 1223)     // 十六进制表示，用 a-f 表示      4c7
	fmt.Printf("%X\n", 1223)     // 十六进制表示，用 A-F 表示      4c7
	fmt.Printf("%U\n", 1233)     // Unicode 表示
	fmt.Printf("%b\n", 12.34)    // 无小数部分，两位指数的科学计数法 6946802425218990p-49
	fmt.Printf("%e\n", 12.345)   // 科学计数法，e 表示   1.234500e+01
	fmt.Printf("%E\n", 12.34455) // 科学计数法，E 表示   1.234455E+01
	fmt.Printf("%f\n", 12.3456)  // 有小数部分，无指数部分   12.345600
	fmt.Printf("%g\n", 12.3456)  // 根据实际情况采用 %e 或 %f 输出  12.3456
	fmt.Printf("%G\n", 12.3456)  // 根据实际情况采用 %E 或 %f 输出  12.3456
	fmt.Printf("%s\n", "wqdew")  // 直接输出字符串或者 []byte         wqdew
	fmt.Printf("%q\n", "dedede") // 双引号括起来的字符串             "dedede"
	fmt.Printf("%x\n", "abczxc") // 每个字节用两字节十六进制表示，a-f 表示  6162637a7863
	fmt.Printf("%X\n", "asdzxc") // 每个字节用两字节十六进制表示，A-F 表示  6173647A7863
	fmt.Printf("%p\n", 0x123)    // 0x 开头的十六进制数表示
}
```