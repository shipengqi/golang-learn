package main

import "fmt"

type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}

// Array 实现 Interface 接口
type Array []int

func (arr Array) Len() int {
    return len(arr)
}

func (arr Array) Less(i, j int) bool {
    return arr[i] < arr[j]
}

func (arr Array) Swap(i, j int) {
    arr[i], arr[j] = arr[j], arr[i]
}

// 匿名接口(anonymous interface)
type reverse struct {
    Interface
}

// 重写(override)
func (r reverse) Less(i, j int) bool {
    return r.Interface.Less(j, i)
}

// 构造 reverse Interface
func Reverse(data Interface) Interface {
    return &reverse{data}
}

func main() {
    arr := Array{1, 2, 3}
    rarr := Reverse(arr)
    fmt.Println(arr.Less(0,1))
    fmt.Println(rarr.Less(0,1))
}