package main

import "fmt"


func f() {}

var g = "g"

func main() {
    f := "f"
    fmt.Println(f) // "f"; local var f shadows package-level func f
    fmt.Println(g) // "g"; package-level var
    // fmt.Println(h)  // compile error: undefined: h
}