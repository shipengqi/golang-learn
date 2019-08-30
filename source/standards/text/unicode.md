---
title: unicode
---

# unicode
go 对 unicode 的支持包含三个包 :

- `unicode`
- `unicode/utf8`
- `unicode/utf16`

unicode 包包含基本的字符判断函数。utf8 包主要负责 `rune` 和 `byte` 之间的转换。utf16 包负责 `rune` 和 `uint16` 数组之间的转换。

## unicode 包

unicode 包含了对 `rune` 的判断。这个包把所有 unicode 涉及到的编码进行了分类，使用结构

```go
type RangeTable struct {
	R16         []Range16
	R32         []Range32
	LatinOffset int
}
```

来表示这个功能的字符集。这些字符集都集中列表在 `table.go` 这个源码里面。

比如控制字符集：

```golang
var _Pc = &RangeTable{
	R16: []Range16{
		{0x005f, 0x203f, 8160},
		{0x2040, 0x2054, 20},
		{0xfe33, 0xfe34, 1},
		{0xfe4d, 0xfe4f, 1},
		{0xff3f, 0xff3f, 1},
	},
}
```

回到包的函数，我们看到有下面这些判断函数：

```
func IsControl(r rune) bool  // 是否控制字符
func IsDigit(r rune) bool  // 是否阿拉伯数字字符，即 0-9
func IsGraphic(r rune) bool // 是否图形字符
func IsLetter(r rune) bool // 是否字母
func IsLower(r rune) bool // 是否小写字符
func IsMark(r rune) bool // 是否符号字符
func IsNumber(r rune) bool // 是否数字字符，比如罗马数字Ⅷ也是数字字符
func IsOneOf(ranges []*RangeTable, r rune) bool // 是否是 RangeTable 中的一个
func IsPrint(r rune) bool // 是否可打印字符
func IsPunct(r rune) bool // 是否标点符号
func IsSpace(r rune) bool // 是否空格
func IsSymbol(r rune) bool // 是否符号字符
func IsTitle(r rune) bool // 是否 title case
func IsUpper(r rune) bool // 是否大写字符
```

例子：

```go
func main() {
	single := '\u0015'
	fmt.Println(unicode.IsControl(single))  //true
	single = '\ufe35'
	fmt.Println(unicode.IsControl(single)) // false

	digit := rune('1')
	fmt.Println(unicode.IsDigit(digit)) //true
	fmt.Println(unicode.IsNumber(digit)) //true
	letter := rune(' Ⅷ ')
	fmt.Println(unicode.IsDigit(letter)) //false
	fmt.Println(unicode.IsNumber(letter)) //true
}
```

## utf8 包

utf8 里面的函数就有一些字节和字符的转换。

判断是否符合 utf8 编码的函数：
```go
func Valid(p []byte) bool
func ValidRune(r rune) bool
func ValidString(s string) bool
```


判断 rune 的长度的函数：
- `func RuneLen(r rune) int`

判断字节串或者字符串的 rune 数
- `func RuneCount(p []byte) int`
- `func RuneCountInString(s string) (n int)`

编码和解码 rune 到 byte
- `func DecodeRune(p []byte) (r rune, size int)`
- `func EncodeRune(p []byte, r rune) int`

## 2.5.3 utf16 包 ##

utf16 的包的函数就比较少了。

将 int16 和 rune 进行转换
- `func Decode(s []uint16) []rune`
- `func Encode(s []rune) []uint16`