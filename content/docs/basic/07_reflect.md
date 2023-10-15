---
title: 反射
weight: 7
---

# 反射

**反射机制，能够在运行时更新变量和检查它们的值、调用它们的方法和它们支持的内在操作，而不需要在编译时就知道
这些变量的具体类型**。弥补了静态语言在动态行为上的一些不足。

## reflect.TypeOf
`reflect.TypeOf` 获取类型信息。
`reflect.TypeOf` 接受任意的 `interface{}` 类型, 并以 `reflect.Type` 形式返回其动态类型：
```go
t := reflect.TypeOf(3)  // a reflect.Type
fmt.Println(t.String()) // "int"
fmt.Println(t)          // "int"

type X int
func main() {
	var a X = 20
	t := reflect.TypeOf(a)
	fmt.Println(t.Name(), t.Kind()) // X int
}
```

上面的代码，**注意区分 `Type` 和 `Kind`，前者表示真实类型（静态类型），后者表示底层类型**。所以在判断类型时，
要选择正确的方式。
```go
type X int
type Y int
func main() {
	var a, b X = 10, 20
	var c Y = 30
	ta, tb, tc := reflect.TypeOf(a), reflect.TypeOf(b), reflect.TypeOf(c)
	fmt.Println(ta == tb, ta == tc) // true false
	fmt.Println(ta.Kind() == tc.Kind()) // true
}
```

### Elem
Elem 方法返回指针，数组，切片，字典或通道的基类型。

```go
fmt.Println(reflect.TypeOf(map[string]int{}).Elem()) // int
```
## reflect.ValueOf
`reflect.ValueOf` 专注于对象实例数据读写。
`reflect.ValueOf` 接受任意的 `interface{}` 类型, 并以 `reflect.Value` 形式返回其动态值：
```go
v := reflect.ValueOf(3) // a reflect.Value
fmt.Println(v)          // "3"
fmt.Printf("%v\n", v)   // "3"
fmt.Println(v.String()) // NOTE: "<int Value>"

x := struct {
    Name string
}{expected}
val := reflect.ValueOf(x)
field := val.Field(0)
fmt.Println(val)            // {Chris}
fmt.Println(field)          //  Chris
fmt.Println(field.String()) // Chris
```

在 Go 中不能对切片使用等号运算符。你可以写一个函数迭代每个元素来检查它们的值。但是一种
比较简单的办法是使用 `reflect.DeepEqual`，它在判断两个变量是否相等时十分有用。

```go
func TestSumAll(t *testing.T) {

    got := SumAll([]int{1,2}, []int{0,9})
    want := []int{3, 9}

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```

**注意**，`reflect.DeepEqual` 不是「类型安全」的，所以有时候会发生比较怪异的行为。比如：
```go
func TestSumAll(t *testing.T) {

    got := SumAll([]int{1,2}, []int{0,9})
    want := "bob"

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```
尝试比较 `slice` 和 `string`。这显然是不合理的，但是却通过了测试。所以使用 `reflect.DeepEqual` 比较简洁但是在使用时需多加小心。