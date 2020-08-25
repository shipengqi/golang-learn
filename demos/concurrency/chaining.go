package main

import (
	"fmt"
)

var ngoroutine = 100

func f(left, right chan int) {
	n := 1 + <-right
	fmt.Println("-------", n)
	left <- n
}

func main() {
	leftmost := make(chan int)
	var left, right chan int = nil, leftmost
	for i := 0; i < ngoroutine; i++ {
		left, right = right, make(chan int) // 每次循环就会交换 left 和 right channel
		go f(left, right)
	}
	right <- 0 // 这个 right 是最新创建的一个 channel
	fmt.Println("bang!")
	x := <-leftmost // wait for completion
	fmt.Println(x)  // 100000, ongeveer 1,5 s
}
