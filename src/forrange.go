package main

import (
  "fmt"
  "os"
)

func main() {
  // 声明两个string类型的变量
  // 如果变量没有显式初始化，则被隐式地赋予其类型的零值，
  // 数值类型是0，字符串类型是空字符串""
  // 这个例子不需要索引，但 range 的语法要求,  要处理元素,  必须处理索引
  // 这里的 _ 是空标识符，因为Go不允许有无用的变量，
  // 空标识符可用于任何语法需要变量名但程序逻辑不需要的时候
  var s, sep string
  for _, arg := range os.Args[1:] { 
  	s += sep + arg
  	sep = " "
  }
  fmt.Println(s)
}