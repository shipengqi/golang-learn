# 反射

反射机制，能够在运行时更新变量和检查它们的值、调用它们的方法和它们支持的内在操作，而不需要在编译时就知道这些变量的具体类型。
弥补了静态语言在动态行为上的一些不足。

`reflect.TypeOf`接受任意的`interface{}`类型, 并以`reflect.Type`形式返回其动态类型：
```go
t := reflect.TypeOf(3)  // a reflect.Type
fmt.Println(t.String()) // "int"
fmt.Println(t)          // "int"
```

`reflect.ValueOf`接受任意的`interface{}`类型, 并以`reflect.Value`形式返回其动态值：
```go
v := reflect.ValueOf(3) // a reflect.Value
fmt.Println(v)          // "3"
fmt.Printf("%v\n", v)   // "3"
fmt.Println(v.String()) // NOTE: "<int Value>"
```