package main

import "fmt"

func fibonacci(n int) int {
	if n < 2 {
		return n
	}

	return fibonacci(n - 2) + fibonacci(n -1)
}

func main() {
	for i := 0; i < 10; i ++ {
		fmt.Printf("%d\t", fibonacci(i))
	}
}