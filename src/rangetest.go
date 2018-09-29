package main

import "fmt"

func main() {
	x := []int{100, 101, 102}

	for key, value := range x {
		fmt.Println(key, value)
	}

	for _, value := range x {
		fmt.Println(value)
	}

	for key := range x {
		fmt.Println(key)
	}
}