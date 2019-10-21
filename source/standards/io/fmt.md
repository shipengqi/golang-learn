---
title: fmt 包
---

# fmt 包
`fmt` 包实现了格式化 I/O 函数，有关格式化输入输出的方法有两大类：`Scan` 和 `Print`。

`print.go` 文件中定义了如下函数：

## Print
```Go
// 普通输出，不带换行符
func Print(a ...interface{}) (n int,  err error)
func Fprint(w io.Writer,  a ...interface{}) (n int,  err error)
func Sprint(a ...interface{}) string

// 输出内容时会加上换行符
func Println(a ...interface{}) (n int,  err error)
func Fprintln(w io.Writer,  a ...interface{}) (n int,  err error)
func Sprintln(a ...interface{}) string

// 按照指定格式化文本输出内容
func Printf(format string,  a ...interface{}) (n int,  err error)
func Fprintf(w io.Writer,  format string,  a ...interface{}) (n int,  err error)
func Sprintf(format string,  a ...interface{}) string
```

如果前缀是 "F", 则指定了 `io.Writer`
如果前缀是 "S", 则是输出到字符串
```
// 输出内容到标准输出 os.Stdout
Print
Printf
Println
// 输出内容到指定的 io.Writer
Fprint
Fprintf
Fprintln
// 输出内容到字符串，并返回
Sprint
Sprintf
Sprintln
```

## Scan
`scan.go` 文件中定义了如下函数：
```go
// 读取内容时不关注换行
func Scan(a ...interface{}) (n int,  err error)
func Fscan(r io.Reader,  a ...interface{}) (n int,  err error)
func Sscan(str string,  a ...interface{}) (n int,  err error)

// 读取到换行时停止，并要求一次提供一行所有条目
func Scanln(a ...interface{}) (n int,  err error)
func Fscanln(r io.Reader,  a ...interface{}) (n int,  err error)
func Sscanln(str string,  a ...interface{}) (n int,  err error) 

// 根据格式化文本读取
func Scanf(format string,  a ...interface{}) (n int,  err error)
func Fscanf(r io.Reader,  format string,  a ...interface{}) (n int,  err error)
func Sscanf(str string,  format string,  a ...interface{}) (n int,  err error)
```

如果前缀是 "F", 则指定了 `io.Reader`
如果前缀是 "S", 则是从字符串读取
```
// 从标准输入os.Stdin读取文本
Scan
Scanf
Scanln
// 从指定的 io.Reader 接口读取文本
Fscan
Fscanf
Fscanln
// 从一个参数字符串读取文本
Sscan
Sscanf
Sscanln
```

## 占位符
**普通占位符**
	
	占位符						说明						举例										输出
	%v		相应值的默认格式。								Printf("%v", site)，Printf("%+v", site)	{studygolang}，{Name:studygolang}
			在打印结构体时，“加号”标记（%+v）会添加字段名
	%#v		相应值的 Go 语法表示							Printf("#v", site)						main.Website{Name:"studygolang"}
	%T		相应值的类型的 Go 语法表示						Printf("%T", site)						main.Website
	%%		字面上的百分号，并非值的占位符					Printf("%%")							%

**布尔占位符**

	占位符						说明						举例										输出
	%t		单词 true 或 false。							Printf("%t", true)						true

**整数占位符**

	占位符						说明						举例									输出
	%b		二进制表示									Printf("%b", 5)						101
	%c		相应Unicode码点所表示的字符					Printf("%c", 0x4E2D)				中
	%d		十进制表示									Printf("%d", 0x12)					18
	%o		八进制表示									Printf("%d", 10)					12
	%q		单引号围绕的字符字面值，由 Go 语法安全地转义		Printf("%q", 0x4E2D)				'中'
	%x		十六进制表示，字母形式为小写 a-f				    Printf("%x", 13)					d
	%X		十六进制表示，字母形式为大写 A-F				    Printf("%x", 13)					D
	%U		Unicode格式：U+1234，等同于 "U+%04X"			Printf("%U", 0x4E2D)				U+4E2D

**浮点数和复数的组成部分（实部和虚部）**

	占位符						说明												举例									输出
	%b		无小数部分的，指数为二的幂的科学计数法，与 strconv.FormatFloat	
			的 'b' 转换格式一致。例如 -123456p-78
	%e		科学计数法，例如 -1234.456e+78									Printf("%e", 10.2)							1.020000e+01
	%E		科学计数法，例如 -1234.456E+78									Printf("%e", 10.2)							1.020000E+01
	%f		有小数点而无指数，例如 123.456									Printf("%f", 10.2)							10.200000
	%g		根据情况选择 %e 或 %f 以产生更紧凑的（无末尾的0）输出				Printf("%g", 10.20)							10.2
	%G		根据情况选择 %E 或 %f 以产生更紧凑的（无末尾的0）输出				Printf("%G", 10.20+2i)						(10.2+2i)

**字符串与字节切片**

	占位符						说明												举例									输出
	%s		输出字符串表示（string 类型或 []byte)							Printf("%s", []byte ("Hello world"))		Hello world
	%5s		指定长度的字符串，这里是以 5 为例，表示最小宽度为 5				    Printf("%5s", []byte ("Hello world"))		Hello
    %-5s	最小宽度为 5（左对齐）
    %.5s	最大宽度为 5
    %5.7s	最小宽度为 5，最大宽度为 7
    %-5.7s	最小宽度为 5，最大宽度为 7（左对齐）
    %5.3s	如果宽度大于 3，则截断
    %05s	如果宽度小于 5，就会在字符串前面补零
	%q		双引号围绕的字符串，由 Go 语法安全地转义							Printf("%q", "Hello world")				    "Hello world"
	%x		十六进制，小写字母，每字节两个字符								Printf("%x", "golang")						676f6c616e67
	%X		十六进制，大写字母，每字节两个字符								Printf("%X", "golang")						676F6C616E67

**指针**

	占位符						说明												举例									输出
	%p		十六进制表示，前缀 0x											Printf("%p", &site)							0x4f57f0
	
**其它标记**

	占位符						说明												举例									输出
	+		总打印数值的正负号；对于%q（%+q）保证只输出 ASCII 编码的字符。			Printf("%+q", "中文")					"\u4e2d\u6587"
	-		在右侧而非左侧填充空格（左对齐该区域）
	#		备用格式：为八进制添加前导 0（%#o），为十六进制添加前导 0x（%#x）或	Printf("%#U", '中')						U+4E2D '中'
			0X（%#X），为 %p（%#p）去掉前导 0x；如果可能的话，%q（%#q）会打印原始
			（即反引号围绕的）字符串；如果是可打印字符，%U（%#U）会写出该字符的
			Unicode 编码形式（如字符 x 会被打印成 U+0078 'x'）。
	' '		（空格）为数值中省略的正负号留出空白（% d）；
			以十六进制（% x, % X）打印字符串或切片时，在字节之间用空格隔开
	0		填充前导的0而非空格；对于数字，这会将填充移到正负号之后
	
示例：
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