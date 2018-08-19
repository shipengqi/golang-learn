package main

import (
  "fmt"
  "os"
)

func main() {
  // 声明两个string类型的变量
  // 如果变量没有显式初始化，则被隐式地赋予其类型的零值，
  // 数值类型是0，字符串类型是空字符串""
  var s, sep string
  for i := 1; i < len(os.Args); i ++ {
  	s += sep + os.Args[i]
  	sep = " "
  }
  fmt.Println(s)
}