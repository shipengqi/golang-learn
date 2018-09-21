package main

import "fmt"

func getSequence() func() int {
	i := 0
	return func() int {
		i ++
		return i
	}
}

func main() {
   /* nextNumber 为一个函数，函数 i 为 0 */
   nextNumber := getSequence()  

   /* 调用 nextNumber 函数，i 变量自增 1 并返回 */
   fmt.Println(nextNumber()) // 1
   fmt.Println(nextNumber()) // 2
   fmt.Println(nextNumber()) // 3
   
   /* 创建新的函数 nextNumber1，并查看结果 */
   nextNumber1 := getSequence()  
   fmt.Println(nextNumber1()) // 1
   fmt.Println(nextNumber1()) // 2
}