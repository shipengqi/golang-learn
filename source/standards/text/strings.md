---
title: strings
---

# strings

字符串常见操作有：

- 字符串长度；
- 求子串；
- 是否存在某个字符或子串；
- 子串出现的次数（字符串匹配）；
- 字符串分割（切分）为 `[]string`；
- 字符串是否有某个前缀或后缀；
- 字符或子串在字符串中首次出现的位置或最后一次出现的位置；
- 通过某个字符串将 `[]string` 连接起来；
- 字符串重复几次；
- 字符串中子串替换；
- 大小写转换；
- `Trim` 操作；
- ...

## 前缀和后缀

`HasPrefix` 判断字符串 `s` 是否以 `prefix` 开头：

```go
strings.HasPrefix(s, prefix string) bool
```

`HasSuffix` 判断字符串 `s` 是否以 `suffix` 结尾：

```go
strings.HasSuffix(s, suffix string) bool
```
示例：

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	var str string = "This is an example of a string"
	fmt.Printf("T/F? Does the string \"%s\" have prefix %s? ", str, "Th")
	fmt.Printf("%t\n", strings.HasPrefix(str, "Th"))
}
```

输出：
```
T/F? Does the string "This is an example of a string" have prefix Th? true
```
	


## 判断是否包含字符串
`Contains` 判断字符串 `s` 是否包含 `substr`：

```go
strings.Contains(s, substr string) bool
```

## 获取某个子字串在字符串中的位置（索引）

`Index` 返回字符串 `str` 在字符串 `s` 中的索引（`str` 的第一个字符的索引），`-1` 表示字符串 `s` 不包含字符串 `str`：

```go
strings.Index(s, str string) int
```

`LastIndex` 返回字符串 `str` 在字符串 `s` 中最后出现位置的索引（`str` 的第一个字符的索引），`-1` 表示字符串 `s` 不包含字符串 `str`：

```go
strings.LastIndex(s, str string) int
```

示例：

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	var str string = "Hi, I'm Marc, Hi."

	fmt.Printf("The position of \"Marc\" is: ")
	fmt.Printf("%d\n", strings.Index(str, "Marc"))

	fmt.Printf("The position of the first instance of \"Hi\" is: ")
	fmt.Printf("%d\n", strings.Index(str, "Hi"))
	fmt.Printf("The position of the last instance of \"Hi\" is: ")
	fmt.Printf("%d\n", strings.LastIndex(str, "Hi"))

	fmt.Printf("The position of \"Burger\" is: ")
	fmt.Printf("%d\n", strings.Index(str, "Burger"))
}
```

输出：
```
The position of "Marc" is: 8
The position of the first instance of "Hi" is: 0
The position of the last instance of "Hi" is: 14
The position of "Burger" is: -1
```

## 计算字符串出现次数

`Count` 用于计算字符串 `str` 在字符串 `s` 中出现的非重叠次数：

```go
strings.Count(s, str string) int
```

示例：

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	var str string = "Hello, how is it going, Hugo?"
	var manyG = "gggggggggg"

	fmt.Printf("Number of H's in %s is: ", str)
	fmt.Printf("%d\n", strings.Count(str, "H"))

	fmt.Printf("Number of double g's in %s is: ", manyG)
	fmt.Printf("%d\n", strings.Count(manyG, "gg"))
}
```

输出：
```
Number of H's in Hello, how is it going, Hugo? is: 2
Number of double g’s in gggggggggg is: 5
```

## 字符串替换
尽量不使用正则，否则会影响性能。

`Replace` 用于将字符串 `str` 中的前 `n` 个字符串 `old` 替换为字符串 `new`，并返回一个新的字符串，如果 `n = -1` 则替换所
有字符串 `old` 为字符串 `new`：

```go
strings.Replace(str, old, new, n) string
```
示例：

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
    fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
    fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))
}
```

输出：
```
oinky oinky oink
moo moo moo
```


## 重复字符串

`Repeat` 用于重复 `count` 次字符串 `s` 并返回一个新的字符串：

```go
strings.Repeat(s, count int) string
```

示例：

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
    fmt.Println("ba" + strings.Repeat("na", 2))
}
```

输出：
```
banana
```

## 大小写转换

`ToLower` 将字符串中的 Unicode 字符全部转换为小写字符：

```go
strings.ToLower(s) string
```

`ToUpper` 将字符串中的 Unicode 字符全部转换为大写字符：

```go
strings.ToUpper(s) string
```

示例：

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	var orig string = "Hey, how are you George?"
	var lower string
	var upper string

	fmt.Printf("The original string is: %s\n", orig)
	lower = strings.ToLower(orig)
	fmt.Printf("The lowercase string is: %s\n", lower)
	upper = strings.ToUpper(orig)
	fmt.Printf("The uppercase string is: %s\n", upper)
}
```

输出：
```
The original string is: Hey, how are you George?
The lowercase string is: hey, how are you george?
The uppercase string is: HEY, HOW ARE YOU GEORGE?
```


## 修改字符串
`Trim` 系列函数可以删除字符串首尾的连续多余字符，包括：

```go
// 删除字符串首尾的字符
func Trim(s string, cutset string) string

// 删除字符串首的字符
func TrimLeft(s string, cutset string) string

// 删除字符串尾部的字符
func TrimRight(s string, cutset string) string

// 删除字符串首尾的空格
func TrimSpace(s string) string
```
示例：

```go
s := "cutjjjcut"
// 将字符串 s 首尾的 `cut` 去除掉
newStr := strings.Trim(s, "cut")

fmt.Println(newStr)
```

输出：
```
jjj
```

## JOIN
`Join` 函数将字符串数组（或 `slice`）连接起来：
```go
func Join(a []string, sep string) string
```

示例：

```go
fmt.Println(strings.Join([]string{"name=xxx", "age=xx"}, "&"))
```

输出：
```
name=xxx&age=xx
```

## 分割字符串

### Fields
```go
// 用一个或多个连续的空格分隔字符串 s，返回子字符串的数组（slice）
func Fields(s string) []string
```

示例：
```go
fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   ")) // Fields are: ["foo" "bar" "baz"]
```

`Fields` 使用一个或多个空格分隔，也就是说返回的字符串中不会包含空格字符串。

> 如果字符串 s 只包含空格，则返回空列表 (`[]string` 的长度为 `0`）

### Split 和 SplitAfter、 SplitN 和 SplitAfterN
```go
func Split(s, sep string) []string { return genSplit(s, sep, 0, -1) }
func SplitAfter(s, sep string) []string { return genSplit(s, sep, len(sep), -1) }
func SplitN(s, sep string, n int) []string { return genSplit(s, sep, 0, n) }
func SplitAfterN(s, sep string, n int) []string { return genSplit(s, sep, len(sep), n) }
```

它们都调用了 `genSplit` 函数。这四个函数都是通过 `sep` 进行分割，返回 `[]string`。

- 如果 `sep` 为空，相当于分成一个个的 UTF-8 字符，如 `Split("abc","")`，得到的是 `[a b c]`。
- `Split(s, sep)` 和 `SplitN(s, sep, -1)` 等价。
- `SplitAfter(s, sep)` 和 `SplitAfterN(s, sep, -1)` 等价。

#### Split 和 SplitAfter 的区别
```go
fmt.Printf("%q\n", strings.Split("foo,bar,baz", ","))  // ["foo" "bar" "baz"]
fmt.Printf("%q\n", strings.SplitAfter("foo,bar,baz", ",")) // ["foo," "bar," "baz"]
```
从输出可以看出，`SplitAfter` 会保留 `sep`。

#### SplitN 和 SplitAfterN
这两个函数通过最后一个参数 `n` 控制返回的结果中的 `slice` 中的元素个数：
- 当 `n < 0` 时，返回所有的子字符串
- 当 `n == 0` 时，返回的结果是 `nil`
- 当 `n > 0` 时，表示返回的 `slice` 中最多只有 `n` 个元素，其中，最后一个元素不会分割

```go
fmt.Printf("%q\n", strings.SplitN("foo,bar,baz", ",", 2))                // ["foo" "bar,baz"]
fmt.Printf("%q\n", strings.Split("a,b,c", ","))                          // ["a" "b" "c"]
fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))   // ["" "man " "plan " "canal panama"]
fmt.Printf("%q\n", strings.Split(" xyz ", ""))                           // [" " "x" "y" "z" " "]
fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))              // [""]
```
## 从字符串中读取内容

函数 `strings.NewReader(str)` 用于生成一个 `Reader` 并读取字符串中的内容，然后返回指向该 `Reader` 的指针，从其它类型读取内容的函数还有：

- `Read()` 从 []byte 中读取内容。
- `ReadByte()` 和 `ReadRune()` 从字符串中读取下一个 byte 或者 rune。
