package main

import "fmt"

func main() {
  x := 100
  // p := *int(&x)  // invalid indirect of int(&x) (type int)  cannot convert &x (type *int) to type int
	
	p := (*int)(&x)  // 应该用括号扩住*int，否则会被解析为*(int(&x))
	fmt.Println(p)
}