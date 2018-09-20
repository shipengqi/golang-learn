package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// 这里的 json:"released" json:"color,omitempty" 是成员的Tag
// 这个 Tag 会导致编码后，成员的名字变成可 Tag 的名字
// 比如这里的 Year 在编码后会变为 released，Color 会变为 color
// 注意这个 tag json:"color,omitempty" 有两个值，第一个就是接送对象的名字
// 第二个 omitempty 选项表示当Go 结构体成员为空或零值时不生成JSON对象（false 为零值）。
// 比如下面的Casablanca Color成员为false 编码后json对象没有color属性
// 可以运行这个程序查看结果
type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
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

func main() {
	data, err := json.MarshalIndent(movies, "", "    ")
	if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)
}