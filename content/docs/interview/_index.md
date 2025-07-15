---
title: 常见面试题
weight: 6
draft: true
---

## 基础

1. Go 断言时会发生内存拷贝么

不涉及拷贝的情况：

当接口动态值是指针或引用类型（如 `*struct`、`slice`、`map`、`chan` 等）时，断言操作仅检查类型信息，不会拷贝底层数据。

涉及拷贝的情况：

若接口存储的是值类型（如 `struct`、`int` 等），断言为具体类型时会拷贝一份新值。

2. `defer` 语句中的变量快照可能会失效

```go
func main() {
    a := 1
    defer func() {
        fmt.Println("defferred", a)
    }()
    a = 2
    fmt.Println("normal", a)
}
```

输出：

```bash
normal 2
defferred 2
```

`defer+引用类型`：

```go
func main() {
    a := []int{0}
    defer func() {
        fmt.Println("defferred", a)
    }()
    a[0] = 1
    fmt.Println("normal", a)
}
```

输出：

```bash
normal [1]
defferred [1]
```

3. `init` 函数在什么时候执行？

在 Go 执行之前自动调用，运行在 `main` 函数之前。

- 先执行导入包的 `init`
- 再执行当前包的 `init`
- 执行 `main`

同一个文件中的 `init` 按声明顺序执行，同一个包中的多个文件中的 `init` 按文件名排序执行。

{{< callout type="info" >}}
如果一个 `init` 函数中启动了一个 goroutine，那么这个 goroutine 和 `main.mian` 是并发执行的。
{{< /callout >}}

4. Go 中 `copy` 函数是深拷贝还是浅拷贝？

首先 `copy` 函数只能用于切片。

`copy` 函数的行为既不是完全的深拷贝，也不是完全的浅拷贝，而是针对不同数据类型表现出不同的复制行为。

- 对于基本类型切片是值拷贝，例如 `[]int{1,2,3}`
- 对于引用类型切片是浅拷贝，例如 `[][]int{{1,2,3},{4,5,6}}`，仅复制指针或引用，新旧切片共享底层数据。

5. 为什么数组不能使用 `copy` 函数？

- **数组的长度是类型的一部分**（如 `[3]int` 和 `[4]int` 是不同类型）。
- **赋值即拷贝**：数组赋值会自动触发深拷贝，无需额外操作。

6. 为什么 `map` 的值不能寻址？

因为 Go 的 `map` 不是线程安全的，直接通过指针访问 `map` 中的值，可能出现数据经侦的问题。为了避免这个问题，就直接禁止了。

```go
func main() {
    m := map[string]int{
        "a": 1,
        "b": 2,
    }
    fmt.Println(&m["a"]) // 编译时就会报错，因为 map 不能直接寻址
}
```

7. 引用类型和指针有什么不同？

- **引用类型是指在内存中存储的是数据的地址（引用）**，而不是直接存储数据。例如切片，字典，通道等。
- 指针永存存储变量的内存地址。

8. 字符串和 `[]byte` 转换会发生内存拷贝么？

会。因为字符串是不可变的，而 `[]byte` 是可变的。为了保证这两者的内存安全，转换时会发生内存拷贝。

9. 翻转含有中文、数字、英文的字符串

将字符串转成 `[]rune`，然后进行翻转，最后转成字符串。

```go
func main() {
    s := "你好，世界！"
    runes := []rune(s)
    i := 0
    j := len(runes) - 1
    for i < j {
        runes[i], runes[j] = runes[j], runes[i]
        i++
        j--
    }
    fmt.Println(string(runes))
}
```

{{< callout type="info" >}}
为什么不使用 `[]byte`，对于 `ascii` 字符（英文、数字）来说，一个字符用一个字节表示，所以 `[]byte` 是足够的。但是对于中文通常需要 3 个字节表示一个字符。而 `rune` 类型就是 Go 用来表示 Unicode 字符的类型，一个 `rune` 类型的值占 4 个字节。
{{< /callout >}}

10. 未初始化的 `map` 可以操作么？

未初始化的 `map` 不能执行**插入**和**更新**操作，否则会导致错误 `panic: assignment to entry in nil map`。

未初始化的 `map` 可以执行**查询**操作，返回零值。删除也不会报错。

```go
package main

import "fmt"

func main() {
	var m map[string]int  // 声明一个 nil map
	fmt.Println(m["key"]) // 输出 0，因为 m 是 nil，零值是 0

	delete(m, "key")      // 删除操作不会报错
	fmt.Println(m)        // 输出 nil
}
```

{{< callout type="info" >}}
Go 对 `delete` 操作做了特殊处理，如果 `map` 是 `nil`，就什么也不做。
{{< /callout >}}

11. 自定义类型与 `[]byte` 互转

可以使用 `encoding/binary` 来实现：

```go
type MyType struct {
    A int32
    B float64
}

func MyTypeSliceToByteSlice(data []MyType) ([]byte, error) {
    buf := new(bytes.Buffer)
    for _, value := range data {
        if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
            return nil, err
        }
    }
    return buf.Bytes(), nil
}


func ByteSliceToMyTypeSlice(data []byte) ([]MyType, error) {
    buf := bytes.NewReader(data)
    var result []MyType
    for buf.Len() > 0 {
        var value MyType
        if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
            return nil, err
        }
        result = append(result, value)
    }
    return result, nil
}
```

建议使用 `encoding/gob`，这个包可以用来序列化和反序列化数据类型，使用更加简单：

```go
import (
    "bytes"
    "encoding/gob"
)

func MyTypeSliceToByteSliceGob(data []MyType) ([]byte, error) {
    buf := new(bytes.Buffer)
    encoder := gob.NewEncoder(buf)
    if err := encoder.Encode(data); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func ByteSliceToMyTypeSliceGob(data []byte) ([]MyType, error) {
    buf := bytes.NewBuffer(data)
    decoder := gob.NewDecoder(buf)
    var result []MyType
    if err := decoder.Decode(&result); err != nil {
        return nil, err
    }
    return result, nil
}
```

12. `struct` 是否可以比较？

可以，**前提是 `struct` 中的所有字段都是可以比较的**。

13. 如何比较包含不可以比较字段的 `struct`？

- 逐个字段比较
- 将 `struct` 序列化成字符串 `json.Marshal`，然后比较字符串。
- `reflect.DeepEqual`，但是性能不高。

14. 如何顺序读取 `map`？

提取出 key 并排序，来实现顺序访问？

```go
package main

import (
	"fmt"
	"sort"
)

func main() {
	m := map[string]int{
		"a": 1,
		"c": 3,
		"b": 2,
	}

	// 提取键
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// 排序键
	sort.Strings(keys)

	// 按顺序访问 map
	for _, k := range keys {
		fmt.Println(k, m[k])
	}
}
```

或者第三方包 `orderedmap`。

15. `switch` 如何强制执行下一个 `case` 块？

使用 `fallthrough` 关键字。其作用是**强制执行下一个 `case` 块，而不考虑下一个 `case` 块的条件判断**。


```go
package main

import "fmt"

func main() {
    num := 1
    switch num {
    case 1:
        fmt.Println("This is case 1")
        fallthrough
    case 2:
        fmt.Println("This is case 2")
    case 3:
        fmt.Println("This is case 3")
    default:
        fmt.Println("This is the default case")
    }
}
```

当 `num` 为 1 时，输出：

```bash
This is case 1
This is case 2
```

16. 为什么常量、字符串、字典不可寻址？

- 常量在编译期间已确定值，直接替换到汇编指令中，不会分配内存空间。
- 字符串由于内容无法修改，为了防止通过指针修改其内容，也是不可寻址的。
- 字段元素不可寻址，是因为字典内部的元素地址可能发生变化，比如扩容，因此禁止对元素寻址。

```go
s := "hello"
// 以下操作都是非法的：
// &s[0]      // 错误：无法获取字符串元素的地址
// &s         // 可以获取字符串变量的地址，但不是字符串内容的地址
```

17. 两个 nil 可能不相等么？

当一个接口变量的实际值为 `nil`，但是动态类型不是空，这个接口变量与 `nil` 不相等。

```go
package main

import "fmt"

func main() {
   var err1 error                 // 声明一个 error 类型的接口变量，初始为 nil
   var err2 error = (*MyError)(nil) // 将一个具体类型的 nil 指针赋值给接口变量

   fmt.Println(err1 == nil) // 输出: true
   fmt.Println(err2 == nil) // 输出: false
   fmt.Println(err1 == err2) // 输出: false
}

type MyError struct{}

func (e *MyError) Error() string {
   return "MyError"
}
```

18. `float` 类型可以作为 `map` 的 key 么？

可以，但是不建议。因为 `float` 本身有精度问题，可能导致两个表面上看起来不一样的浮点数，实际上相同。

```go
package main

import "fmt"

func main() {
    m := map[float64]int{}
    m[1.0] = 1
    m[1.100000000000000000000001] = 2
    m[1.100000000000000000000002] = 3
    for k, v := range m {
       fmt.Println("k=", k, "v=", v)
    }
}
```

输出：

```bash
k= 1 v= 1
k= 1.1 v= 3
```
