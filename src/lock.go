package main

import "fmt"

func main() {
	var x, y int

	go func() {
		x = 1
		fmt.Println(y)
	}()

	go func() {
		y = 1
		fmt.Println(x)
	}()
}