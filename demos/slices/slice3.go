package main

import "fmt"

func myAppend(s []int) []int {
	// 这里 s 虽然改变了，但并不会影响外层函数的 s
	s = append(s, 100)
	return s
}

func myAppendPtr(s *[]int) {
	// 会改变外层 s 本身
	*s = append(*s, 100)
	return
}

// 当直接用切片作为函数参数时，可以改变切片的元素，不能改变切片本身；想要改变切片本身，可以将改变后的切片返回，函数调用者接收改变后的切片或者将切片指针作为函数参数。
func main() {
	s := []int{1, 1, 1}
	newS := myAppend(s)

	fmt.Println(s)    // [1 1 1]
	fmt.Println(newS) // [1 1 1 100]

	s = newS

	myAppendPtr(&s)
	fmt.Println(s) // [1 1 1 100 100]
}
