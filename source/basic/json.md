---
title: 序列化
---
# 序列化
Go 对于其他序列化协议如 `Json`，`XML`，`Protocol Buffers`，都有良好的支持，

由标准库中的 `encoding/json`、`encoding/xml`、`encoding/asn1` 等包提供支持，`Protocol Buffers` 的
由 `github.com/golang/protobuf` 包提供支持，并且这类包都有着相似的 API 接口。

GO 中结构体转为 `JSON` 使用 `json.Marshal`，也就是编码操作：
```go
type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
	Actors []string
}

var movies = []Movie{
	{
		Title: "Casablanca", 
		Year: 1942, 
		Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{
		Title: "Cool Hand Luke",
		Year: 1967, 
		Color: true,
		Actors: []string{"Paul Newman"}},
	{
		Title: "Bullitt", 
		Year: 1968, 
		Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}}}	

data, err := json.Marshal(movies)
if err != nil {
    log.Fatalf("JSON marshaling failed: %s", err)
}
fmt.Printf("%s\n", data)
```

`json.MarshalIndent` 格式化输出 `JSON`，例如：
```go
data, err := json.MarshalIndent(movies, "", "    ")
if err != nil {
    log.Fatalf("JSON marshaling failed: %s", err)
}
fmt.Printf("%s\n", data)
```
输出：
```js
[
    {
        "Title": "Casablanca",
        "released": 1942,
        "Actors": [
            "Humphrey Bogart",
            "Ingrid Bergman"
        ]
    },
    {
        "Title": "Cool Hand Luke",
        "released": 1967,
        "color": true,
        "Actors": [
            "Paul Newman"
        ]
    },
    {
        "Title": "Bullitt",
        "released": 1968,
        "color": true,
        "Actors": [
            "Steve McQueen",
            "Jacqueline Bisset"
        ]
    }
]
```

有没有注意到，`Year` 字段名的成员在编码后变成了 `released`，`Color` 变成了小写的 `color`。这是因为结构体的成员 Tag 导致的，
如上面的：
```go
Year   int  `json:"released"`
Color  bool `json:"color,omitempty"`
```

结构体的成员 Tag 可以是任意的字符串面值，但是通常是一系列用空格分隔的 `key:"value"` 键值对序列；因为值中含义双引号字符，
因此成员 Tag 一般用原生字符串面值的形式书写。`json` 开头键名对应的值用于控制 `encoding/json` 包的编码和解码的行为，
并且 `encoding/...` 下面其它的包也遵循这个约定。成员 `Tag` 中 `json` 对应值的第一部分用于指定 JSON 对象的名字，
比如将 Go 语言中的 `TotalCount` 成员对应到 JSON 中的 `total_count` 对象。`Color` 成员的 Tag 还带了一个额外的 `omitempty` 
选项，表示当 Go 语言结构体成员为空或零值时不生成 JSON 对象（这里 `false` 为零值）。果然，`Casablanca` 是一个黑白电影，
并没有输出 `Color` 成员。

**注意，只有导出的结构体成员才会被编码**

解码操作，使用`json.Unmarshal`：
```go
var titles []struct{ Title string }
if err := json.Unmarshal(data, &titles); err != nil {
    log.Fatalf("JSON unmarshaling failed: %s", err)
}
fmt.Println(titles) // "[{Casablanca} {Cool Hand Luke} {Bullitt}]"
```
通过定义合适的Go语言数据结构，我们可以选择性地解码JSON中感兴趣的成员。

基于流式的解码器 `json.Decoder`。针对输出流的  `json.Encoder` 编码对象
